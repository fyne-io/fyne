package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// TableKernx represents a 'kern' or 'kerx' kerning table.
// It supports both Microsoft and Apple formats.
type TableKernx []KernSubtable

func parseTableKerx(data []byte, numGlyphs int) (TableKernx, error) {
	if len(data) < 8 {
		return nil, errors.New("invalid kerx table (EOF)")
	}
	// version := binary.BigEndian.Uint16(data)
	// padding
	nTables := binary.BigEndian.Uint32(data[4:])

	// "sanitize" before allocating
	if len(data) < int(nTables)*12 {
		return nil, errors.New("invalid kerx table (EOF)")
	}
	currentOffset := 8
	out := make(TableKernx, nTables)
	for i := range out {
		if len(data) < currentOffset {
			return nil, errors.New("invalid kerx table (EOF)")
		}
		var (
			size int
			err  error
		)
		out[i], size, err = parseKerxSubtable(data[currentOffset:], numGlyphs)
		if err != nil {
			return nil, err
		}
		currentOffset += size
	}
	return out, nil
}

// unified coverage flags (from 'kerx')
const (
	kerxBackwards   = 1 << 12
	kerxVariation   = 1 << 13
	kerxCrossStream = 1 << 14
	kerxVertical    = 1 << 15
)

// KernSubtable contains kerning information.
// Some formats provides an easy lookup method: see SimpleKerns.
// Others require a state machine to interpret it.
type KernSubtable struct {
	Data       interface{ isKernSubtable() }
	coverage   uint16 // high bit of the Coverage field
	IsExtended bool   // `true` for AAT `kerx` subtables
	TupleCount int    // 0 for scalar values
}

// IsHorizontal returns true if the subtable has horizontal kerning values.
func (k KernSubtable) IsHorizontal() bool {
	return k.coverage&kerxVertical == 0
}

// IsBackwards returns true if state-table based should process the glyphs backwards.
func (k KernSubtable) IsBackwards() bool {
	return k.coverage&kerxBackwards != 0
}

// IsCrossStream returns true if the subtable has cross-stream kerning values.
func (k KernSubtable) IsCrossStream() bool {
	return k.coverage&kerxCrossStream != 0
}

// IsVariation returns true if the subtable has variation kerning values.
func (k KernSubtable) IsVariation() bool {
	return k.coverage&kerxVariation != 0
}

func parseKerxSubtable(data []byte, numGlyphs int) (out KernSubtable, _ int, err error) {
	out.IsExtended = true
	const kerxSubtableHeaderLength = 12
	if len(data) < kerxSubtableHeaderLength {
		return out, 0, errors.New("invalid kerx subtable (EOF)")
	}
	length := int(binary.BigEndian.Uint32(data))
	if len(data) < int(length) {
		return out, 0, errors.New("invalid kerx subtable (EOF)")
	}

	coverage := binary.BigEndian.Uint32(data[4:])
	out.TupleCount = int(binary.BigEndian.Uint32(data[8:]))

	out.coverage = uint16(coverage >> 16) // high bit

	data = data[:length]
	const formatMask = 0x000000FF
	switch f := coverage & formatMask; f {
	case 0:
		out.Data, err = parseKernxSubtable0(data, kerxSubtableHeaderLength, true, out.TupleCount)
	case 1:
		out.Data, err = parseKernxSubtable1(data, kerxSubtableHeaderLength, true, numGlyphs, out.TupleCount)
	case 2:
		out.Data, err = parseKernxSubtable2(data, kerxSubtableHeaderLength, true, numGlyphs, out.TupleCount)
	case 4:
		out.Data, err = parseKerxSubtable4(data, numGlyphs)
	case 6:
		out.Data, err = parseKerxSubtable6(data, numGlyphs, out.TupleCount)
	default:
		return out, 0, fmt.Errorf("unsupported kerx subtable format: %d", f)
	}
	return out, length, err
}

type Kern0 []KerningPair

func (Kern0) isKernSubtable() {}

// data starts at the subtable header with length `headerLength`
// `extended` is true for 'kerx', false for 'kern'
func parseKernxSubtable0(data []byte, headerLength int, extended bool, tupleCount int) (Kern0, error) {
	binSearchHeaderLength := 8
	if extended {
		binSearchHeaderLength = 16
	}
	if len(data) < headerLength+binSearchHeaderLength {
		return nil, errors.New("invalid kern/x subtable format 0 (EOF)")
	}
	var nPairs int
	if extended {
		nPairs = int(binary.BigEndian.Uint32(data[headerLength:]))
	} else {
		nPairs = int(binary.BigEndian.Uint16(data[headerLength:]))
	}
	out, err := parseKerningPairs(data[headerLength+binSearchHeaderLength:], nPairs)
	if err != nil {
		return nil, err
	}

	if tupleCount != 0 { // interpret values as offset
		for i, pair := range out {
			if len(data) < int(uint16(pair.Value))+2 {
				return nil, errors.New("invalid kern/x subtable format 0 (EOF)")
			}
			out[i].Value = int16(binary.BigEndian.Uint16(data[pair.Value:]))
		}
	}
	return out, err
}

func (k Kern0) KernPair(left, right GID) int16 {
	key := uint32(left)<<16 | uint32(right)
	low, high := 0, len(k)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := k[mid].key()
		if key < p {
			high = mid
		} else if key > p {
			low = mid + 1
		} else {
			return k[mid].Value
		}
	}
	return 0
}

// Kernx1 state entry flags
const (
	Kerx1Push        = 0x8000 // If set, push this glyph on the kerning stack.
	Kerx1DontAdvance = 0x4000 // If set, don't advance to the next glyph before going to the new state.
	Kerx1Reset       = 0x2000 // If set, reset the kerning data (clear the stack)
	Kern1Offset      = 0x3FFF // Byte offset from beginning of subtable to the  value table for the glyphs on the kerning stack.
)

type Kern1 struct {
	Values  []int16 // After successful parsing, may be safely indexed by AATStateEntry.AsKernxIndex() from `Machine`
	Machine AATStateTable
}

func (Kern1) isKernSubtable() {}

// data starts at the subtable header
// tupleCount is optionnal
func parseKernxSubtable1(data []byte, headerLength int, extended bool, numGlyphs int, tupleCount int) (out Kern1, err error) {
	if len(data) < headerLength {
		return out, errors.New("invalid kern/x subtable format 1 (EOF)")
	}
	data = data[headerLength:]

	// we need the offset to the data following the stateTable
	var valuesOffset, extraDataSize int
	if extended {
		if len(data) < aatExtStateHeaderSize+4 {
			return out, errors.New("invalid kerx subtable format 1 (EOF)")
		}
		valuesOffset = int(binary.BigEndian.Uint32(data[aatExtStateHeaderSize:]))
		extraDataSize = 2
	} else {
		if len(data) < aatStateHeaderSize+2 {
			return out, errors.New("invalid kern subtable format 1 (EOF)")
		}
		valuesOffset = int(binary.BigEndian.Uint16(data[aatStateHeaderSize:]))
		extraDataSize = 0
	}
	if len(data) < valuesOffset {
		return out, errors.New("invalid kern/x subtable format 1 (EOF)")
	}
	out.Machine, err = parseStateTable(data, extraDataSize, extended, numGlyphs)
	if err != nil {
		return out, err
	}

	// find the maximum index need in the values array
	var maxi uint16
	for i := range out.Machine.entries {
		entry := &out.Machine.entries[i]
		if !extended { // start by resolving offset -> index
			offset := int(entry.Flags & Kern1Offset)
			if offset == 0 || offset < valuesOffset {
				binary.BigEndian.PutUint16(entry.data[:], 0xFFFF)
			} else {
				index := uint16((offset - valuesOffset) / 2)
				binary.BigEndian.PutUint16(entry.data[:], index)
			}
		}
		if index := entry.AsKernxIndex(); index != 0xFFFF && index > maxi {
			maxi = index
		}
	}

	if tupleCount == 0 {
		tupleCount = 1
	}
	nbUint16Min := tupleCount * int(maxi+1)
	if len(data) < valuesOffset+2*nbUint16Min {
		return out, errors.New("invalid kern/x subtable format 1 (EOF)")
	}
	data = data[valuesOffset:]
	/* From Apple 'kern' spec:
	 * "Each pops one glyph from the kerning stack and applies the kerning value to it.
	 * The end of the list is marked by an odd value... */

	out.Values = make([]int16, 0, nbUint16Min)
	for len(data) >= 2 { // gracefully handle missing odd value
		v := int16(binary.BigEndian.Uint16(data))
		out.Values = append(out.Values, v)
		data = data[2:]
		if len(out.Values) >= nbUint16Min && v&1 != 0 {
			break
		}
	}
	return out, nil
}

type Kern2 struct {
	// Values are pre-multiplied by the number of bytes in one row and
	// offset by the offset of the array from the start of the subtable.
	left               Class
	right              Class // Values are pre-multiplied by 2
	tableData          []byte
	kerningArrayOffset int  // start of the actual kerning data in `kernings`
	hasTuples          bool // if true, the kerning value is actually an offset into `tableData`
}

func (Kern2) isKernSubtable() {}

func (k Kern2) KernPair(left, right GID) int16 {
	l, _ := k.left.ClassID(left)
	r, _ := k.right.ClassID(right)
	index := int(l) + int(r)
	if len(k.tableData) < index+2 || index < k.kerningArrayOffset {
		return 0
	}
	kernVal := binary.BigEndian.Uint16(k.tableData[index:])
	if k.hasTuples && int(kernVal)+2 <= len(k.tableData) {
		kernVal = binary.BigEndian.Uint16(k.tableData[kernVal:])
	}
	return int16(kernVal)
}

// data starts at the subtable header
func parseKernxSubtable2(data []byte, headerLength int, extended bool, numGlyphs int, tupleCount int) (out Kern2, err error) {
	subHeaderLength := 8
	if extended {
		subHeaderLength = 16
	}
	if len(data) < headerLength+subHeaderLength {
		return out, errors.New("invalid kern/x subtable format 2 (EOF)")
	}

	out.hasTuples = tupleCount != 0

	var leftOffset, rightOffset, arrayOffset uint32
	if extended {
		// out.rowWidth = binary.BigEndian.Uint32(data[headerLength:])
		leftOffset = binary.BigEndian.Uint32(data[headerLength+4:])
		rightOffset = binary.BigEndian.Uint32(data[headerLength+8:])
		arrayOffset = binary.BigEndian.Uint32(data[headerLength+12:])
	} else {
		// out.rowWidth = uint32(binary.BigEndian.Uint16(data[headerLength:]))
		leftOffset = uint32(binary.BigEndian.Uint16(data[headerLength+2:]))
		rightOffset = uint32(binary.BigEndian.Uint16(data[headerLength+4:]))
		arrayOffset = uint32(binary.BigEndian.Uint16(data[headerLength+6:]))
	}

	if len(data) < int(arrayOffset) {
		return out, errors.New("invalid kerx subtable format 2 (EOF)")
	}

	if extended {
		out.left, err = parseAATLookupTable(data, leftOffset, numGlyphs, false)
		if err != nil {
			return out, err
		}
		out.right, err = parseAATLookupTable(data, rightOffset, numGlyphs, false)
		if err != nil {
			return out, err
		}
	} else {
		out.left, err = parseClassFormat1(data[leftOffset:], 2)
		if err != nil {
			return out, fmt.Errorf("invalid kern subtable format 2: %s", err)
		}
		out.right, err = parseClassFormat1(data[rightOffset:], 2)
		if err != nil {
			return out, fmt.Errorf("invalid kern subtable format 2: %s", err)
		}
	}
	out.tableData = data                      // since the class already has the offset, just store the raw slice
	out.kerningArrayOffset = int(arrayOffset) // store it to check for invalid offset values
	return out, err
}

type KerxAnchor interface {
	isKernAnchor()
}

type KerxAnchorControl struct {
	Mark, Current uint16
}

func (KerxAnchorControl) isKernAnchor() {}

type KerxAnchorAnchor struct {
	Mark, Current uint16
}

func (KerxAnchorAnchor) isKernAnchor() {}

type KerxAnchorCoordinates struct {
	MarkX, MarkY, CurrentX, CurrentY int16
}

func (KerxAnchorCoordinates) isKernAnchor() {}

type Kerx4 struct {
	Anchors []KerxAnchor
	Machine AATStateTable
	flags   uint32
}

func (Kerx4) isKernSubtable() {}

// ActionType returns 0, 1 or 2 .
func (k Kerx4) ActionType() uint8 {
	const ActionType = 0xC0000000 // A two-bit field containing the action type.
	return uint8(k.flags & ActionType >> 30)
}

// data starts at the subtable header
func parseKerxSubtable4(data []byte, numGlyphs int) (out Kerx4, err error) {
	// we need the offset to the data following the stateTable
	if len(data) < 12+aatExtStateHeaderSize+4 {
		return out, errors.New("invalid kerx subtable format 4 (EOF)")
	}
	data = data[12:]
	out.flags = binary.BigEndian.Uint32(data[aatExtStateHeaderSize:])
	out.Machine, err = parseStateTable(data, 2, true, numGlyphs)
	if err != nil {
		return out, err
	}

	// find the maximum index need in the actions array
	var maxi uint16
	for _, entry := range out.Machine.entries {
		if index := entry.AsKernxIndex(); index != 0xFFFF && index > maxi {
			maxi = index
		}
	}

	const Offset = 0x00FFFFFF // Masks the offset in bytes from the beginning of the subtable to the beginning of the control point table.

	controlOffset := int(out.flags & Offset)
	switch actionType := out.ActionType(); actionType {
	case 0:
		if len(data) < controlOffset+4*int(maxi+1) {
			return out, errors.New("invalid kerx subtable format 4 (EOF)")
		}
		out.Anchors = make([]KerxAnchor, int(maxi+1))
		for i := range out.Anchors {
			anchor := KerxAnchorControl{
				Mark:    binary.BigEndian.Uint16(data[controlOffset+4*i:]),
				Current: binary.BigEndian.Uint16(data[controlOffset+4*i+2:]),
			}
			out.Anchors[i] = anchor
		}
	case 1:
		if len(data) < controlOffset+4*int(maxi+1) {
			return out, errors.New("invalid kerx subtable format 4 (EOF)")
		}
		out.Anchors = make([]KerxAnchor, int(maxi+1))
		for i := range out.Anchors {
			anchor := KerxAnchorAnchor{
				Mark:    binary.BigEndian.Uint16(data[controlOffset+4*i:]),
				Current: binary.BigEndian.Uint16(data[controlOffset+4*i+2:]),
			}
			out.Anchors[i] = anchor
		}
	case 2:
		if len(data) < controlOffset+8*int(maxi+1) {
			return out, errors.New("invalid kerx subtable format 4 (EOF)")
		}
		out.Anchors = make([]KerxAnchor, int(maxi+1))
		for i := range out.Anchors {
			anchor := KerxAnchorCoordinates{
				MarkX:    int16(binary.BigEndian.Uint16(data[controlOffset+8*i:])),
				MarkY:    int16(binary.BigEndian.Uint16(data[controlOffset+8*i+2:])),
				CurrentX: int16(binary.BigEndian.Uint16(data[controlOffset+8*i+4:])),
				CurrentY: int16(binary.BigEndian.Uint16(data[controlOffset+8*i+6:])),
			}
			out.Anchors[i] = anchor
		}
	default:
		return out, fmt.Errorf("invalid kerx subtable format 4 action type: %d", actionType)
	}

	return out, nil
}

type Kerx6 struct {
	row    Class // Values are pre-multiplied by `columnCount`
	column Class
	// with rowCount * columnCount
	// for tuples the values are the first element of the tuple
	kernings              []int16
	rowCount, columnCount uint16
}

func (Kerx6) isKernSubtable() {}

func (k Kerx6) KernPair(left, right GID) int16 {
	l, _ := k.row.ClassID(left)
	r, _ := k.column.ClassID(right)
	index := int(l) + int(r)
	if len(k.kernings) < index {
		return 0
	}
	return k.kernings[index]
}

// data starts at the subtable header
func parseKerxSubtable6(data []byte, numGlyphs, tupleCount int) (out Kerx6, err error) {
	if len(data) < 12+20 {
		return out, errors.New("invalid kerx subtable format 2 (EOF)")
	}
	flags := binary.BigEndian.Uint32(data[12:])
	out.rowCount = binary.BigEndian.Uint16(data[12+4:])
	out.columnCount = binary.BigEndian.Uint16(data[12+6:])
	rowTableOffset := binary.BigEndian.Uint32(data[12+8:])
	colTableOffset := binary.BigEndian.Uint32(data[12+12:])
	arrayOffset := int(binary.BigEndian.Uint32(data[12+16:]))

	isLong := flags&1 != 0

	out.row, err = parseAATLookupTable(data, rowTableOffset, numGlyphs, isLong)
	if err != nil {
		return out, err
	}
	out.column, err = parseAATLookupTable(data, colTableOffset, numGlyphs, isLong)
	if err != nil {
		return out, err
	}
	if len(data) < arrayOffset {
		return out, errors.New("invalid kerx subtable format 6 (EOF)")
	}
	length := int(out.rowCount) * int(out.columnCount)
	var tmp []uint32
	if isLong {
		if len(data) < arrayOffset+length*4 {
			return out, errors.New("invalid kerx subtable format 6 (EOF)")
		}
		tmp = parseUint32s(data[arrayOffset:], length)
	} else {
		if len(data) < arrayOffset+length*2 {
			return out, errors.New("invalid kerx subtable format 6 (EOF)")
		}
		tmp = make([]uint32, length)
		for i := range tmp {
			tmp[i] = uint32(binary.BigEndian.Uint16(data[arrayOffset+2*i:]))
		}
	}

	out.kernings = make([]int16, len(tmp))
	if tupleCount != 0 { // interpret kern values as offset
		if len(data) < 12+24 {
			return out, errors.New("invalid kerx subtable format 2 (EOF)")
		}
		// If the tupleCount is 1 or more, then the kerning array contains offsets from the beginning
		// of the kerningVectors table to a tupleCount-dimensional vector of FUnits controlling the kerning.
		kerningVectorsOffet := int(binary.BigEndian.Uint32(data[12+20:]))
		for i, v := range tmp {
			if len(data) < kerningVectorsOffet+int(v)+2 {
				return out, errors.New("invalid kerx subtable format 2 (EOF)")
			}
			out.kernings[i] = int16(binary.BigEndian.Uint16(data[kerningVectorsOffet+int(v):]))
		}
	} else {
		// a kerning value greater than an int16 should not happen
		for i, v := range tmp {
			out.kernings[i] = int16(v)
		}
	}
	return out, err
}
