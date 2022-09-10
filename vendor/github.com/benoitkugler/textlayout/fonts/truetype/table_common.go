package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

// Class group glyph indices.
// Conceptually it is a map[GID]uint32, but it may
// be implemented more efficiently.
type Class interface {
	// ClassID returns the class ID for the provided glyph. Returns (0, false)
	// for glyphs not covered by this class.
	ClassID(GID) (uint32, bool)

	// GlyphSize returns the number of glyphs covered.
	GlyphSize() int

	// Extent returns the maximum class ID + 1. This is the length
	// required for an array to be indexed by the class values.
	Extent() int
}

// parseClass parse `buf`, starting at `offset`.
func parseClass(buf []byte, offset uint16) (Class, error) {
	if len(buf) < int(offset)+2 {
		return nil, errors.New("invalid class table (EOF)")
	}
	buf = buf[offset:]
	switch format := binary.BigEndian.Uint16(buf); format {
	case 1:
		return parseClassFormat1(buf[2:], 2)
	case 2:
		return parseClassLookupFormat2(buf)
	default:
		return nil, fmt.Errorf("unsupported class definition format %d", format)
	}
}

type classFormat1 struct {
	classIDs   []uint32 // array of target class IDs. gi is the index into that array (minus StartGlyph).
	startGlyph GID
}

func (c classFormat1) ClassID(gi GID) (uint32, bool) {
	if gi < c.startGlyph || gi >= c.startGlyph+GID(len(c.classIDs)) {
		return 0, false
	}
	return c.classIDs[gi-c.startGlyph], true
}

func (c classFormat1) GlyphSize() int { return len(c.classIDs) }

func (c classFormat1) Extent() int {
	max := uint32(0)
	for _, cid := range c.classIDs {
		if cid >= max {
			max = cid
		}
	}
	return int(max) + 1
}

// parseClassFormat1 parses a class table, with format 1.
// For compatibility reasons, it expects `buf` to start at the first glyph,
// not at the class format.
// `valueByteSize` is 1, 2 or 4
func parseClassFormat1(data []byte, valueByteSize int) (out classFormat1, err error) {
	// ClassDefFormat 1: startGlyphID, glyphCount, []classValueArray
	const headerSize = 4 // excluding classFormat
	if len(data) < headerSize {
		return out, errors.New("invalid class format 1 (EOF)")
	}

	out.startGlyph = GID(binary.BigEndian.Uint16(data))
	count := int(binary.BigEndian.Uint16(data[2:]))
	if len(data) < 4+count*valueByteSize {
		return out, errors.New("invalid class format 1 (EOF)")
	}
	data = data[4:]
	out.classIDs = make([]uint32, count)
	switch valueByteSize {
	case 1:
		for i, b := range data[0:count] {
			out.classIDs[i] = uint32(b)
		}
	case 2:
		for i := range out.classIDs {
			out.classIDs[i] = uint32(binary.BigEndian.Uint16(data[i*2:]))
		}
	case 4:
		for i := range out.classIDs {
			out.classIDs[i] = binary.BigEndian.Uint32(data[i*4:])
		}
	default:
		panic("invalid byte size")
	}

	return out, nil
}

type classRangeRecord struct {
	start, end    gid
	targetClassID uint32
}

type classFormat2 []classRangeRecord

// 'adapted' from golang/x/image/font/sfnt
func (c classFormat2) ClassID(gi_ GID) (uint32, bool) {
	gi := uint16(gi_)
	num := len(c)
	if num == 0 {
		return 0, false
	}

	// classRange is an array of startGlyphID, endGlyphID and target class ID.
	// Ranges are non-overlapping.
	// E.g. 130, 135, 1   137, 137, 5   etc

	idx := sort.Search(num, func(i int) bool { return gi <= c[i].start })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if class := c[idx]; gi == c[idx].start {
			return class.targetClassID, true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if class := c[idx]; gi >= class.start && gi <= class.end {
			return class.targetClassID, true
		}
	}

	return 0, false
}

func (c classFormat2) GlyphSize() int {
	out := 0
	for _, class := range c {
		out += int(class.end - class.start + 1)
	}
	return out
}

func (c classFormat2) Extent() int {
	max := uint32(0)
	for _, r := range c {
		if r.targetClassID >= max {
			max = r.targetClassID
		}
	}
	return int(max) + 1
}

// ClassDefFormat 2: classFormat, classRangeCount, []classRangeRecords
func parseClassLookupFormat2(buf []byte) (classFormat2, error) {
	const headerSize = 4 // including classFormat
	if len(buf) < headerSize {
		return nil, errors.New("invalid class format 2 (EOF)")
	}

	num := int(binary.BigEndian.Uint16(buf[2:]))
	if len(buf) < headerSize+num*6 {
		return nil, errors.New("invalid class format 2 (EOF)")
	}

	out := make(classFormat2, num)
	for i := range out {
		out[i].start = gid(binary.BigEndian.Uint16(buf[headerSize+i*6:]))
		out[i].end = gid(binary.BigEndian.Uint16(buf[headerSize+i*6+2:]))
		out[i].targetClassID = uint32(binary.BigEndian.Uint16(buf[headerSize+i*6+4:]))
	}
	return out, nil
}

// Coverage specifies all the glyphs affected by a substitution or
// positioning operation described in a subtable.
// Conceptually is it a []GlyphIndex, but it may be implemented for efficiently.
// See the concrete types `CoverageList` and `CoverageRanges`.
type Coverage interface {
	// Index returns the index of the provided glyph, or
	// `false` if the glyph is not covered by this lookup.
	// Note: this method is injective: two distincts, covered glyphs are mapped
	// to distincts tables.
	Index(GID) (int, bool)

	// Size return the number of glyphs covered. For non empty Coverages, it is also
	// 1 + (maximum index returned)
	Size() int
}

// if l[i] = gi then gi has coverage index of i
func parseCoverage(buf []byte, offset uint32) (Coverage, error) {
	if len(buf) < int(offset)+2 { // format and count
		return nil, errors.New("invalid coverage table")
	}
	buf = buf[offset:]
	switch format := binary.BigEndian.Uint16(buf); format {
	case 1:
		// Coverage Format 1: coverageFormat, glyphCount, []glyphArray
		return fetchCoverageList(buf[2:])
	case 2:
		// Coverage Format 2: coverageFormat, rangeCount, []rangeRecords{startGlyphID, endGlyphID, startCoverageIndex}
		return fetchCoverageRange(buf[2:])
	default:
		return nil, fmt.Errorf("unsupported coverage format %d", format)
	}
}

// CoverageList is a coverage with format 1.
// The glyphs are sorted in ascending order.
type CoverageList []GID

func (cl CoverageList) Index(gi GID) (int, bool) {
	num := len(cl)
	idx := sort.Search(num, func(i int) bool { return gi <= cl[i] })
	if idx < num && cl[idx] == gi {
		return idx, true
	}
	return 0, false
}

func (cl CoverageList) Size() int { return len(cl) }

// func (cl coverageList) maxIndex() int { return len(cl) - 1 }

func fetchCoverageList(buf []byte) (CoverageList, error) {
	const headerSize, entrySize = 2, 2
	if len(buf) < headerSize {
		return nil, errInvalidGPOSKern
	}

	num := int(binary.BigEndian.Uint16(buf))
	if len(buf) < headerSize+num*entrySize {
		return nil, errInvalidGPOSKern
	}

	out := make(CoverageList, num)
	for i := range out {
		out[i] = GID(binary.BigEndian.Uint16(buf[headerSize+2*i:]))
	}
	return out, nil
}

// CoverageRange store a range of indexes, starting from StartCoverage.
// For example, for the glyphs 12,13,14,15, and the indexes 7,8,9,10,
// the CoverageRange would be {12, 15, 7}.
type CoverageRange struct {
	Start, End    GID
	StartCoverage int
}

// CoverageRanges is a coverage with format 2.
// Ranges are non-overlapping.
// The following GlyphIDs/index pairs are stored as follows:
//	 glyphs: 130, 131, 132, 133, 134, 135, 137
//	 indexes: 0, 1, 2, 3, 4, 5, 6
//   ranges: {130, 135, 0}    {137, 137, 6}
// StartCoverage is used to calculate the index without counting
// the length of the preceeding ranges
type CoverageRanges []CoverageRange

func (cr CoverageRanges) Index(gi GID) (int, bool) {
	num := len(cr)
	if num == 0 {
		return 0, false
	}

	idx := sort.Search(num, func(i int) bool { return gi <= cr[i].Start })
	// idx either points to a matching start, or to the next range (or idx==num)
	// e.g. with the range example from above: 130 points to 130-135 range, 133 points to 137-137 range

	// check if gi is the start of a range, but only if sort.Search returned a valid result
	if idx < num {
		if rang := cr[idx]; gi == rang.Start {
			return int(rang.StartCoverage), true
		}
	}
	// check if gi is in previous range
	if idx > 0 {
		idx--
		if rang := cr[idx]; gi >= rang.Start && gi <= rang.End {
			return rang.StartCoverage + int(gi-rang.Start), true
		}
	}

	return 0, false
}

func (cr CoverageRanges) Size() int {
	size := 0
	for _, r := range cr {
		size += int(r.End - r.Start + 1)
	}
	return size
}

// func (cr coverageRanges) maxIndex() int {
// 	lastRange := cr[len(cr)-1]
// 	return lastRange.startCoverage + int(lastRange.end-lastRange.start)
// }

func fetchCoverageRange(buf []byte) (CoverageRanges, error) {
	const headerSize, entrySize = 2, 6
	if len(buf) < headerSize {
		return nil, errInvalidGPOSKern
	}

	num := int(binary.BigEndian.Uint16(buf))
	if len(buf) < headerSize+num*entrySize {
		return nil, errInvalidGPOSKern
	}

	out := make(CoverageRanges, num)
	for i := range out {
		out[i].Start = GID(binary.BigEndian.Uint16(buf[headerSize+i*entrySize:]))
		out[i].End = GID(binary.BigEndian.Uint16(buf[headerSize+i*entrySize+2:]))
		out[i].StartCoverage = int(binary.BigEndian.Uint16(buf[headerSize+i*entrySize+4:]))
	}
	return out, nil
}

// DeviceTable is either an GPOSDeviceHinting for standard fonts,
// or a GPOSDeviceVariation for variable fonts.
type DeviceTable interface {
	isDevice()
}

func (DeviceHinting) isDevice()   {}
func (DeviceVariation) isDevice() {}

type DeviceHinting struct {
	// with length endSize - startSize + 1
	Values []int8
	// correction range, in ppem
	StartSize, EndSize uint16
}

// GetDelta returns the hint for the given `ppem`, scaled by `scale`.
// It returns 0 for out of range `ppem` values.
func (dev DeviceHinting) GetDelta(ppem uint16, scale int32) int32 {
	if ppem == 0 {
		return 0
	}

	if ppem < dev.StartSize || ppem > dev.EndSize {
		return 0
	}

	pixels := dev.Values[ppem-dev.StartSize]

	return int32(pixels) * (scale / int32(ppem))
}

type DeviceVariation VariationStoreIndex

func parseDeviceTable(data []byte, offset uint16) (DeviceTable, error) {
	if len(data) < int(offset)+6 {
		return nil, errors.New("invalid positionning device subtable (EOF)")
	}
	first := binary.BigEndian.Uint16(data[offset:])
	second := binary.BigEndian.Uint16(data[offset+2:])
	format := binary.BigEndian.Uint16(data[offset+4:])

	switch format {
	case 1, 2, 3:
		var out DeviceHinting

		out.StartSize, out.EndSize = first, second
		if out.EndSize < out.StartSize {
			return nil, errors.New("invalid positionning device subtable")
		}

		nbPerUint16 := 16 / (1 << format) // 8, 4 or 2
		outLength := int(out.EndSize - out.StartSize + 1)
		var count int
		if outLength%nbPerUint16 == 0 {
			count = outLength / nbPerUint16
		} else {
			// add padding
			count = outLength/nbPerUint16 + 1
		}
		uint16s, err := parseUint16s(data[offset+6:], count)
		if err != nil {
			return nil, err
		}
		out.Values = make([]int8, count*nbPerUint16) // handle rounding error by reslicing after
		switch format {
		case 1:
			for i, u := range uint16s {
				uint16As2Bits(out.Values[i*8:], u)
			}
		case 2:
			for i, u := range uint16s {
				uint16As4Bits(out.Values[i*4:], u)
			}
		case 3:
			for i, u := range uint16s {
				uint16As8Bits(out.Values[i*2:], u)
			}
		}
		out.Values = out.Values[:outLength]
		return out, nil
	case 0x8000:
		return DeviceVariation{DeltaSetOuter: first, DeltaSetInner: second}, nil
	default:
		return nil, fmt.Errorf("unsupported positionning device subtable: %d", format)
	}
}

// write 8 elements
func uint16As2Bits(dst []int8, u uint16) {
	const mask = 0xFE // 11111110
	dst[0] = int8((0-uint8(u>>15&1))&mask | uint8(u>>14&1))
	dst[1] = int8((0-uint8(u>>13&1))&mask | uint8(u>>12&1))
	dst[2] = int8((0-uint8(u>>11&1))&mask | uint8(u>>10&1))
	dst[3] = int8((0-uint8(u>>9&1))&mask | uint8(u>>8&1))
	dst[4] = int8((0-uint8(u>>7&1))&mask | uint8(u>>6&1))
	dst[5] = int8((0-uint8(u>>5&1))&mask | uint8(u>>4&1))
	dst[6] = int8((0-uint8(u>>3&1))&mask | uint8(u>>2&1))
	dst[7] = int8((0-uint8(u>>1&1))&mask | uint8(u>>0&1))
}

// write 4 elements
func uint16As4Bits(dst []int8, u uint16) {
	const mask = 0xF8 // 11111000

	dst[0] = int8((0-uint8(u>>15&1))&mask | uint8(u>>12&0x07))
	dst[1] = int8((0-uint8(u>>11&1))&mask | uint8(u>>8&0x07))
	dst[2] = int8((0-uint8(u>>7&1))&mask | uint8(u>>4&0x07))
	dst[3] = int8((0-uint8(u>>3&1))&mask | uint8(u>>0&0x07))
}

// write 2 elements
func uint16As8Bits(dst []int8, u uint16) {
	dst[0] = int8(u >> 8)
	dst[1] = int8(u)
}
