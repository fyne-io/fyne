package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	_ SimpleKerns = Kern0{}
	_ SimpleKerns = Kern2{}
	_ SimpleKerns = Kerx6{}
)

// SimpleKerns store a compact form of the kerning
// values. It is not implemented by complex AAT kerning subtables.
type SimpleKerns interface {
	// KernPair return the kern value for the given pair, or zero.
	// The value is expressed in glyph units and
	// is negative when glyphs should be closer.
	KernPair(left, right GID) int16
	// // Size returns the number of kerning pairs
	// Size() int
}

// assume non overlapping kerns, otherwise the return value is undefined
type kernUnions []SimpleKerns

func (ks kernUnions) KernPair(left, right GID) int16 {
	for _, k := range ks {
		out := k.KernPair(left, right)
		if out != 0 {
			return out
		}
	}
	return 0
}

// there are several formats for the 'kern' table, due to the
// differents specs from Apple and Microsoft. The concepts are similar,
// but the bit sizes of the various fields differ.
// We apply the following logic:
//   - read the first uint16 -> it's always the major version
//   - if it's 0, we have a Miscrosoft table
//   - if it's 1, we have an Apple table. We read the next uint16,
//     to differentiate between the old and the new Apple format.
func parseKernTable(input []byte, numGlyphs int) (TableKernx, error) {
	if len(input) < 4 {
		return nil, errors.New("invalid kern table (EOF)")
	}

	var (
		numTables            uint32
		subtableHeaderLength int
	)

	major := binary.BigEndian.Uint16(input)
	switch major {
	case 0:
		numTables = uint32(binary.BigEndian.Uint16(input[2:]))
		subtableHeaderLength = 6
		input = input[4:]
	case 1:
		subtableHeaderLength = 8
		nextUint16 := binary.BigEndian.Uint16(input[2:])
		if nextUint16 == 0 {
			// either new format or old format with 0 subtables, the later being invalid (or at least useless)
			if len(input) < 8 {
				return nil, errors.New("invalid kern table version 1 (EOF)")
			}
			numTables = binary.BigEndian.Uint32(input[4:])
			input = input[8:]
		} else {
			// old format
			numTables = uint32(nextUint16)
			input = input[4:]
		}

	default:
		return nil, fmt.Errorf("unsupported kern table version: %d", major)
	}

	out := make([]KernSubtable, numTables)
	var (
		err    error
		nbRead int
	)
	for i := range out {
		if len(input) < nbRead {
			return nil, errors.New("invalid kern table EOF)")
		}
		input = input[nbRead:]
		out[i], nbRead, err = parseKernSubtable(input, subtableHeaderLength, numGlyphs)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// also returns the length of the subtable
func parseKernSubtable(input []byte, subtableHeaderLength, numGlyphs int) (out KernSubtable, length int, err error) {
	out.IsExtended = false
	if len(input) < subtableHeaderLength {
		return out, 0, errors.New("invalid kern subtable (EOF)")
	}
	var format byte
	if subtableHeaderLength == 6 { // OT format
		length = int(binary.BigEndian.Uint16(input[2:]))
		coverage := binary.BigEndian.Uint16(input[4:])
		// synthesize a coverage flag following kerx conventions
		const (
			Horizontal  = 0x01
			CrossStream = 0x04
		)
		if coverage&Horizontal == 0 { // vertical
			out.coverage |= kerxVertical
		}
		if coverage&CrossStream != 0 {
			out.coverage |= kerxCrossStream
		}
		format = byte(coverage >> 8)
	} else { // AAT format
		length = int(binary.BigEndian.Uint32(input))
		out.coverage = binary.BigEndian.Uint16(input[4:])
		format = byte(out.coverage) // low bit
	}

	switch format {
	case 0:
		out.Data, err = parseKernxSubtable0(input, subtableHeaderLength, false, 0)
	case 1:
		out.Data, err = parseKernxSubtable1(input, subtableHeaderLength, false, numGlyphs, 0)
	case 2:
		out.Data, err = parseKernxSubtable2(input, subtableHeaderLength, false, numGlyphs, 0)
	case 3:
		out.Data, err = parseKernSubtable3(input)
	default:
		return out, 0, fmt.Errorf("invalid kern subtable format %d", format)
	}

	return out, length, err
}

type KerningPair struct {
	Left, Right GID
	// Note: For 'kerx' table version 4 with tuples, this is
	// the first element of the kerning tuple.
	Value int16
}

func (kp KerningPair) key() uint32 { return uint32(kp.Left)<<16 | uint32(kp.Right) }

func parseKerningPairs(data []byte, count int) ([]KerningPair, error) {
	const entrySize = 6
	if len(data) < entrySize*count {
		return nil, errors.New("invalid kerning pairs array (EOF)")
	}
	out := make([]KerningPair, count)
	for i := range out {
		out[i].Left = GID(binary.BigEndian.Uint16(data[entrySize*i:]))
		out[i].Right = GID(binary.BigEndian.Uint16(data[entrySize*i+2:]))
		out[i].Value = int16(binary.BigEndian.Uint16(data[entrySize*i+4:]))
	}
	return out, nil
}

// Kern3 is the Apple kerning subtable format 3
type Kern3 struct {
	leftClass, rightClass []uint8   // length glyphCount
	kernIndex             [][]uint8 // size length(leftClass) x length(rightClass)
	kernValues            []int16
}

func (Kern3) isKernSubtable() {}

func (ks Kern3) KernPair(left, right GID) int16 {
	if int(left) >= len(ks.leftClass) || int(right) >= len(ks.rightClass) { // should not happend
		return 0
	}

	index := ks.kernIndex[ks.leftClass[left]][ks.rightClass[right]] // sanitized during parsing
	return ks.kernValues[index]                                     // sanitized during parsing
}

func parseKernSubtable3(data []byte) (out Kern3, err error) {
	// apple 'kern' header
	if len(data) < 8+6 {
		return out, errors.New("invalid kern subtable format 3 (EOF)")
	}
	glyphCount := int(binary.BigEndian.Uint16(data[8:]))
	kernValueCount, leftClassCount, rightClassCount := data[10], data[11], data[12]
	// flags is ignored
	if len(data) < 8+6+2*int(kernValueCount)+2*glyphCount+int(leftClassCount)*int(rightClassCount) {
		return out, errors.New("invalid kern subtable format 3 (EOF)")
	}
	data = data[8+6:]
	out.kernValues = make([]int16, kernValueCount)
	for i := range out.kernValues {
		out.kernValues[i] = int16(binary.BigEndian.Uint16(data[2*i:]))
	}
	data = data[2*kernValueCount:]

	out.leftClass = data[:glyphCount]
	out.rightClass = data[glyphCount : 2*glyphCount]
	data = data[2*glyphCount:]

	out.kernIndex = make([][]uint8, leftClassCount)
	for i := range out.kernIndex {
		out.kernIndex[i] = data[i*int(rightClassCount) : (i+1)*int(rightClassCount)]

		// sanitize index values
		for _, index := range out.kernIndex[i] {
			if index >= kernValueCount {
				return out, errors.New("invalid kern subtable format 3 index value")
			}
		}
	}

	// sanitize class values
	for i := range out.leftClass {
		if out.leftClass[i] >= leftClassCount {
			return out, errors.New("invalid kern subtable format 3 class value")
		}
		if out.rightClass[i] >= rightClassCount {
			return out, errors.New("invalid kern subtable format 3 class value")
		}
	}

	return out, nil
}
