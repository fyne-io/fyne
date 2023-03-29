// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"fmt"
)

// Kerx is the extended kerning table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6kerx.html
type Kerx struct {
	version uint16         // The version number of the extended kerning table (currently 2, 3, or 4).
	padding uint16         // Unused; set to zero.
	nTables uint32         // The number of subtables included in the extended kerning table.
	Tables  []KerxSubtable `arrayCount:"ComputedField-nTables"`
}

// extended versions

// binarygen: argument=valuesCount int
type KerxSubtable struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	Coverage   uint16 // Circumstances under which this table is used.
	padding    byte   // unused
	version    kerxSTVersion
	TupleCount uint32   // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	Data       KerxData `unionField:"version" arguments:"tupleCount=.TupleCount, valuesCount=valuesCount"`
}

// check and return the subtable length
func (ks *KerxSubtable) parseEnd(src []byte, _ int) (int, error) {
	if L := len(src); L < int(ks.length) {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", ks.length, L)
	}
	return int(ks.length), nil
}

type kerxSTVersion byte

const (
	kerxSTVersion0 kerxSTVersion = iota
	kerxSTVersion1
	kerxSTVersion2
	_
	kerxSTVersion4
	_
	kerxSTVersion6
)

type KerxData interface {
	isKerxData()
}

func (KerxData0) isKerxData() {}
func (KerxData1) isKerxData() {}
func (KerxData2) isKerxData() {}
func (KerxData4) isKerxData() {}
func (KerxData6) isKerxData() {}

// binarygen: argument=tupleCount int
type KerxData0 struct {
	nPairs        uint32
	searchRange   uint32
	entrySelector uint32
	rangeShift    uint32
	Pairs         []Kernx0Record `arrayCount:"ComputedField-nPairs"`
}

// resolve offset for variable fonts
func (kd *KerxData0) parseEnd(src []byte, tupleCount int) (int, error) {
	if tupleCount != 0 { // interpret values as offset
		for i, pair := range kd.Pairs {
			if L, E := len(src), int(uint16(pair.Value))+2; L < E {
				return 0, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
			}
			kd.Pairs[i].Value = int16(binary.BigEndian.Uint16(src[pair.Value:]))
		}
	}
	return len(src), nil
}

type Kernx0Record struct {
	Left, Right GlyphID
	Value       int16
}

// Kernx1 state entry flags
const (
	Kerx1Push        = 0x8000 // If set, push this glyph on the kerning stack.
	Kerx1DontAdvance = 0x4000 // If set, don't advance to the next glyph before going to the new state.
	Kerx1Reset       = 0x2000 // If set, reset the kerning data (clear the stack)
	Kern1Offset      = 0x3FFF // Byte offset from beginning of subtable to the  value table for the glyphs on the kerning stack.
)

// binarygen: argument=tupleCount int
// binarygen: argument=valuesCount int
type KerxData1 struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount, entryDataSize=2"`
	valueTable       Offset32
	Values           []int16 `isOpaque:""`
}

// From Apple 'kern' spec:
// Each pops one glyph from the kerning stack and applies the kerning value to it.
// The end of the list is marked by an odd value...
func parseKernx1Values(src []byte, entries []AATStateEntry, valueTableOffset, tupleCount int) ([]int16, error) {
	// find the maximum index need in the values array
	var maxi uint16
	for _, entry := range entries {
		if index := entry.AsKernxIndex(); index != 0xFFFF && index > maxi {
			maxi = index
		}
	}

	if tupleCount == 0 {
		tupleCount = 1
	}
	nbUint16Min := tupleCount * int(maxi+1)
	if L, E := len(src), valueTableOffset+2*nbUint16Min; L < E {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}

	src = src[valueTableOffset:]
	out := make([]int16, 0, nbUint16Min)
	for len(src) >= 2 { // gracefully handle missing odd value
		v := int16(binary.BigEndian.Uint16(src))
		out = append(out, v)
		src = src[2:]
		if len(out) >= nbUint16Min && v&1 != 0 {
			break
		}
	}
	return out, nil
}

func (kx *KerxData1) parseValues(src []byte, tupleCount, _ int) error {
	var err error
	kx.Values, err = parseKernx1Values(src, kx.Entries, int(kx.valueTable), tupleCount)
	return err
}

type KerxData2 struct {
	rowWidth     uint32    // The number of bytes in each row of the kerning value array
	Left         AATLookup `offsetSize:"Offset32" offsetRelativeTo:"Parent"` // Offset from beginning of this subtable to the left-hand offset table.
	Right        AATLookup `offsetSize:"Offset32" offsetRelativeTo:"Parent"` // Offset from beginning of this subtable to right-hand offset table.
	KerningStart Offset32  // Offset from beginning of this subtable to the start of the kerning array.
	KerningData  []byte    `subsliceStart:"AtStart" arrayCount:"ToEnd"` // indexed by Left + Right
}

// binarygen: argument=valuesCount int
type KerxData4 struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount,entryDataSize=2"`
	Flags            uint32
	Anchors          KerxAnchors `isOpaque:""`
}

func (kd KerxData4) nAnchors() int {
	// find the maximum index need in the actions array
	var maxi uint16
	for _, entry := range kd.Entries {
		if index := entry.AsKernxIndex(); index != 0xFFFF && index > maxi {
			maxi = index
		}
	}
	return int(maxi) + 1
}

func (kd *KerxData4) parseAnchors(src []byte, _ int) error {
	nAnchors := kd.nAnchors()
	const Offset = 0x00FFFFFF // Masks the offset in bytes from the beginning of the subtable to the beginning of the control point table.
	controlOffset := int(kd.Flags & Offset)
	if L := len(src); L < controlOffset {
		return fmt.Errorf("EOF: expected length: %d, got %d", controlOffset, L)
	}
	var err error
	switch kd.ActionType() {
	case 0:
		kd.Anchors, _, err = ParseKerxAnchorControls(src[controlOffset:], nAnchors)
	case 1:
		kd.Anchors, _, err = ParseKerxAnchorAnchors(src[controlOffset:], nAnchors)
	case 2:
		kd.Anchors, _, err = ParseKerxAnchorCoordinates(src[controlOffset:], nAnchors)
	default:
		return fmt.Errorf("invalid Kerx4 anchor format %d", kd.ActionType())
	}
	return err
}

// ActionType returns 0, 1 or 2, according to the anchor format :
//   - 0 : KerxAnchorControls
//   - 1 : KerxAnchorAnchors
//   - 2 : KerxAnchorCoordinates
func (kd KerxData4) ActionType() uint8 {
	const ActionType = 0xC0000000 // A two-bit field containing the action type.
	return uint8((kd.Flags & ActionType) >> 30)
}

type KerxAnchors interface {
	isKerxAnchors()
}

func (KerxAnchorControls) isKerxAnchors()    {}
func (KerxAnchorAnchors) isKerxAnchors()     {}
func (KerxAnchorCoordinates) isKerxAnchors() {}

type KerxAnchorControls struct {
	Anchors []KAControl
}
type KerxAnchorAnchors struct {
	Anchors []KAAnchor
}
type KerxAnchorCoordinates struct {
	Anchors []KACoordinates
}

type KAControl struct {
	Mark, Current uint16
}

type KAAnchor struct {
	Mark, Current uint16
}

type KACoordinates struct {
	MarkX, MarkY, CurrentX, CurrentY int16
}

// binarygen: argument=tupleCount int
// binarygen: argument=valuesCount int
type KerxData6 struct {
	flags                  uint32         // Flags for this subtable. See below.
	rowCount               uint16         // The number of rows in the kerning value array
	columnCount            uint16         // The number of columns in the kerning value array
	rowIndexTableOffset    uint32         // Offset from beginning of this subtable to the row index lookup table.
	columnIndexTableOffset uint32         // Offset from beginning of this subtable to column index offset table.
	kerningArrayOffset     uint32         // Offset from beginning of this subtable to the start of the kerning array.
	kerningVectorOffset    uint32         // Offset from beginning of this subtable to the start of the kerning vectors. This value is only present if the tupleCount for this subtable is 1 or more.
	Row                    AatLookupMixed `isOpaque:"" offsetRelativeTo:"Parent"` // Values are pre-multiplied by `columnCount`
	Column                 AatLookupMixed `isOpaque:"" offsetRelativeTo:"Parent"`
	// with rowCount * columnCount
	// for tuples the values are estParseKerx (Not yet run).the first element of the tuple
	Kernings []int16 `isOpaque:""  offsetRelativeTo:"Parent"`
}

func (kd *KerxData6) parseRow(_, parentSrc []byte, _, valuesCount int) error {
	isExtended := kd.flags&1 != 0
	if L := len(parentSrc); L < int(kd.rowIndexTableOffset) {
		return fmt.Errorf("EOF: expected length: %d, got %d", kd.rowIndexTableOffset, L)
	}
	var err error
	if isExtended {
		kd.Row, _, err = ParseAATLookupExt(parentSrc[kd.rowIndexTableOffset:], valuesCount)
	} else {
		kd.Row, _, err = ParseAATLookup(parentSrc[kd.rowIndexTableOffset:], valuesCount)
	}
	return err
}

func (kd *KerxData6) parseColumn(_, parentSrc []byte, _, valuesCount int) error {
	isExtended := kd.flags&1 != 0
	if L := len(parentSrc); L < int(kd.columnIndexTableOffset) {
		return fmt.Errorf("EOF: expected length: %d, got %d", kd.columnIndexTableOffset, L)
	}
	var err error
	if isExtended {
		kd.Column, _, err = ParseAATLookupExt(parentSrc[kd.columnIndexTableOffset:], valuesCount)
	} else {
		kd.Column, _, err = ParseAATLookup(parentSrc[kd.columnIndexTableOffset:], valuesCount)
	}
	return err
}

func (kd *KerxData6) parseKernings(_, parentSrc []byte, tupleCount, _ int) error {
	isExtended := kd.flags&1 != 0

	length := int(kd.rowCount) * int(kd.columnCount)
	var tmp []uint32
	if isExtended {
		if L, E := len(parentSrc), int(kd.kerningArrayOffset)+length*4; L < E {
			return fmt.Errorf("EOF: expected length: %d, got %d", E, L)
		}
		tmp = make([]uint32, length)
		for i := range tmp {
			tmp[i] = binary.BigEndian.Uint32(parentSrc[int(kd.kerningArrayOffset)+4*i:])
		}
	} else {
		if L, E := len(parentSrc), int(kd.kerningArrayOffset)+length*2; L < E {
			return fmt.Errorf("EOF: expected length: %d, got %d", E, L)
		}
		tmp = make([]uint32, length)
		for i := range tmp {
			tmp[i] = uint32(binary.BigEndian.Uint16(parentSrc[int(kd.kerningArrayOffset)+2*i:]))
		}
	}

	kd.Kernings = make([]int16, len(tmp))
	if tupleCount != 0 { // interpret kern values as offset
		// If the tupleCount is 1 or more, then the kerning array contains offsets from the beginning
		// of the kerningVectors table to a tupleCount-dimensional vector of FUnits controlling the kerning.
		for i, v := range tmp {
			kerningOffset := int(kd.kerningVectorOffset) + int(v)
			if L := len(parentSrc); L < kerningOffset+2 {
				return fmt.Errorf("EOF: expected length: %d, got %d", kerningOffset+2, L)
			}
			kd.Kernings[i] = int16(binary.BigEndian.Uint16(parentSrc[kerningOffset:]))
		}
	} else {
		// a kerning value greater than an int16 should not happen
		for i, v := range tmp {
			kd.Kernings[i] = int16(v)
		}
	}
	return nil
}

//lint:ignore U1000 this type is required so that the code generator add a ParseAATLookupExt function
type dummy struct {
	A AATLookupExt
}
