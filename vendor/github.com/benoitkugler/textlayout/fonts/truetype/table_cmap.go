package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
	"golang.org/x/text/encoding/charmap"
)

type (
	Cmap     = fonts.Cmap
	CmapIter = fonts.CmapIter
)

// TableCmap defines the mapping of character codes to the glyph index values used in the font.
// It may contain more than one subtable, in order to support more than one character encoding scheme.
type TableCmap struct {
	Cmaps            []CmapSubtable
	unicodeVariation unicodeVariations
}

// FindSubtable returns the cmap for the given platform and encoding, or nil if not found.
func (t *TableCmap) FindSubtable(id CmapID) Cmap {
	key := id.key()
	// binary search
	for i, j := 0, len(t.Cmaps); i < j; {
		h := i + (j-i)/2
		entryKey := t.Cmaps[h].ID.key()
		if key < entryKey {
			j = h
		} else if entryKey < key {
			i = h + 1
		} else {
			return t.Cmaps[h].Cmap
		}
	}
	return nil
}

// BestEncoding returns the widest encoding supported. For valid fonts,
// the returned cmap won't be nil.
func (t TableCmap) BestEncoding() (Cmap, fonts.CmapEncoding) {
	// direct adaption from harfbuzz/src/hb-ot-cmap-table.hh

	// Prefer symbol if available.
	if subtable := t.FindSubtable(CmapID{PlatformMicrosoft, PEMicrosoftSymbolCs}); subtable != nil {
		return subtable, fonts.EncSymbol
	}

	/* 32-bit subtables. */
	if cmap := t.FindSubtable(CmapID{PlatformMicrosoft, PEMicrosoftUcs4}); cmap != nil {
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, PEUnicodeFull13}); cmap != nil {
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, PEUnicodeFull}); cmap != nil {
		return cmap, fonts.EncUnicode
	}

	/* 16-bit subtables. */
	if cmap := t.FindSubtable(CmapID{PlatformMicrosoft, PEMicrosoftUnicodeCs}); cmap != nil {
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, PEUnicodeBMP}); cmap != nil {
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, 2}); cmap != nil { // deprecated
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, 1}); cmap != nil { // deprecated
		return cmap, fonts.EncUnicode
	}
	if cmap := t.FindSubtable(CmapID{PlatformUnicode, 0}); cmap != nil { // deprecated
		return cmap, fonts.EncUnicode
	}

	if len(t.Cmaps) != 0 {
		return t.Cmaps[0].Cmap, fonts.EncOther
	}
	return nil, fonts.EncOther
}

type unicodeVariations []variationSelector

func (t unicodeVariations) getGlyphVariant(r, selector rune) (GID, uint8) {
	// binary search
	for i, j := 0, len(t); i < j; {
		h := i + (j-i)/2
		entryKey := t[h].varSelector
		if selector < entryKey {
			j = h
		} else if entryKey < selector {
			i = h + 1
		} else {
			return t[h].getGlyph(r)
		}
	}
	return 0, variantNotFound
}

type cmap0 = fonts.CmapSimple

type cmap4 []cmapEntry16

type cmap4Iter struct {
	data cmap4
	pos1 int // into data
	pos2 int // either into data[pos1].indexes or an offset between start and end
}

func (it *cmap4Iter) Next() bool {
	return it.pos1 < len(it.data)
}

func (it *cmap4Iter) Char() (r rune, gy GID) {
	entry := it.data[it.pos1]
	if entry.indexes == nil {
		r = rune(it.pos2 + int(entry.start))
		gy = GID(uint16(it.pos2) + entry.start + entry.delta)
		if uint16(it.pos2) == entry.end-entry.start {
			// we have read the last glyph in this part
			it.pos2 = 0
			it.pos1++
		} else {
			it.pos2++
		}
	} else { // pos2 is the array index
		r = rune(it.pos2) + rune(entry.start)
		gy = GID(entry.indexes[it.pos2])
		if gy != 0 {
			gy += GID(entry.delta)
		}
		if it.pos2 == len(entry.indexes)-1 {
			// we have read the last glyph in this part
			it.pos2 = 0
			it.pos1++
		} else {
			it.pos2++
		}
	}

	return r, gy
}

func (s cmap4) Iter() CmapIter {
	return &cmap4Iter{data: s}
}

func (s cmap4) Lookup(r rune) (GID, bool) {
	if uint32(r) > 0xffff {
		return 0, false
	}
	// binary search
	c := uint16(r)
	for i, j := 0, len(s); i < j; {
		h := i + (j-i)/2
		entry := s[h]
		if c < entry.start {
			j = h
		} else if entry.end < c {
			i = h + 1
		} else if entry.indexes == nil {
			return GID(c + entry.delta), true
		} else {
			glyph := entry.indexes[c-entry.start]
			if glyph == 0 {
				return 0, false
			}
			return GID(glyph + entry.delta), true
		}
	}
	return 0, false
}

type cmap6or10 struct {
	entries   []uint16
	firstCode rune
}

type cmap6Iter struct {
	data cmap6or10
	pos  int // index into data.entries
}

func (it *cmap6Iter) Next() bool {
	return it.pos < len(it.data.entries)
}

func (it *cmap6Iter) Char() (rune, GID) {
	entry := it.data.entries[it.pos]
	r := rune(it.pos) + it.data.firstCode
	gy := GID(entry)
	it.pos++
	return r, gy
}

func (s cmap6or10) Iter() CmapIter {
	return &cmap6Iter{data: s}
}

func (s cmap6or10) Lookup(r rune) (GID, bool) {
	if r < s.firstCode {
		return 0, false
	}
	c := int(r - s.firstCode)
	if c >= len(s.entries) {
		return 0, false
	}
	return GID(s.entries[c]), true
}

type cmap12 []cmapEntry32

type cmap12Iter struct {
	data cmap12
	pos1 int // into data
	pos2 int // offset from start
}

func (it *cmap12Iter) Next() bool { return it.pos1 < len(it.data) }

func (it *cmap12Iter) Char() (r rune, gy GID) {
	entry := it.data[it.pos1]
	r = rune(it.pos2 + int(entry.start))
	gy = GID(it.pos2 + int(entry.value))
	if uint32(it.pos2) == entry.end-entry.start {
		// we have read the last glyph in this part
		it.pos2 = 0
		it.pos1++
	} else {
		it.pos2++
	}

	return r, gy
}

func (s cmap12) Iter() CmapIter { return &cmap12Iter{data: s} }

func (s cmap12) Lookup(r rune) (GID, bool) {
	c := uint32(r)
	// binary search
	for i, j := 0, len(s); i < j; {
		h := i + (j-i)/2
		entry := s[h]
		if c < entry.start {
			j = h
		} else if entry.end < c {
			i = h + 1
		} else {
			return GID(c - entry.start + entry.value), true
		}
	}
	return 0, false
}

type cmap13 []cmapEntry32

type cmap13Iter struct {
	data cmap13
	pos1 int // into data
	pos2 int // offset from start
}

func (it *cmap13Iter) Next() bool {
	return it.pos1 < len(it.data)
}

func (it *cmap13Iter) Char() (r rune, gy GID) {
	entry := it.data[it.pos1]
	r = rune(it.pos2 + int(entry.start))
	gy = GID(entry.value)
	if uint32(it.pos2) == entry.end-entry.start {
		// we have read the last glyph in this part
		it.pos2 = 0
		it.pos1++
	} else {
		it.pos2++
	}

	return r, gy
}

func (s cmap13) Iter() CmapIter { return &cmap13Iter{data: s} }

func (s cmap13) Lookup(r rune) (GID, bool) {
	c := uint32(r)
	// binary search
	for i, j := 0, len(s); i < j; {
		h := i + (j-i)/2
		entry := s[h]
		if c < entry.start {
			j = h
		} else if entry.end < c {
			i = h + 1
		} else {
			return GID(entry.value), true
		}
	}
	return 0, false
}

// CmapID groups the platform and encoding of a Cmap subtable.
type CmapID struct {
	Platform PlatformID
	Encoding PlatformEncodingID
}

type CmapSubtable struct {
	Cmap Cmap
	ID   CmapID
}

func (c CmapID) key() uint32 { return uint32(c.Platform)<<16 | uint32(c.Encoding) }

// IsSymbolic returns `true` for the special case of a symbolic cmap, for which
// the codepoints are not interpreted as Unicode.
func (c CmapID) IsSymbolic() bool {
	return c.Platform == PlatformMicrosoft && c.Encoding == PEMicrosoftSymbolCs
}

// https://www.microsoft.com/typography/OTSPEC/cmap.htm
// direct adaption from golang.org/x/image/font/sfnt
func parseTableCmap(input []byte) (out TableCmap, err error) {
	const headerSize, entrySize = 4, 8
	if len(input) < headerSize {
		return out, errors.New("invalid 'cmap' table (EOF)")
	}
	// version is skipped
	numSubtables := int(binary.BigEndian.Uint16(input[2:]))

	if len(input) < headerSize+entrySize*numSubtables {
		return out, errors.New("invalid 'cmap' table (EOF)")
	}

	for i := 0; i < numSubtables; i++ {
		bufSubtable := input[headerSize+entrySize*i:]

		var cmap CmapSubtable
		cmap.ID.Platform = PlatformID(binary.BigEndian.Uint16(bufSubtable))
		cmap.ID.Encoding = PlatformEncodingID(binary.BigEndian.Uint16(bufSubtable[2:]))

		offset := binary.BigEndian.Uint32(bufSubtable[4:])
		if len(input) < int(offset)+2 { // format
			return out, errors.New("invalid cmap subtable (EOF)")
		}
		format := binary.BigEndian.Uint16(input[offset:])

		if format == 14 { // special case for variation selector
			if cmap.ID.Platform != PlatformUnicode && cmap.ID.Platform != 5 {
				return out, errors.New("invalid cmap subtable (EOF)")
			}
			out.unicodeVariation, err = parseCmapFormat14(input, offset)
			if err != nil {
				return out, err
			}
		} else if format == 2 { // for now, we just ignore these formats
			continue
		} else {
			cmap.Cmap, err = parseCmapSubtable(format, input, uint32(offset))
			if err != nil {
				return out, err
			}
			out.Cmaps = append(out.Cmaps, cmap)
		}
	}

	if len(out.Cmaps) == 0 {
		return out, errors.New("empty 'cmap' table")
	}

	return out, nil
}

// format 14 has already been handled
func parseCmapSubtable(format uint16, input []byte, offset uint32) (Cmap, error) {
	switch format {
	case 0:
		return parseCmapFormat0(input, offset)
	case 2:
		// parseCmapFormat2(input, offset)
		return nil, fmt.Errorf("unsupported cmap subtable format: %d", format)
	case 4:
		return parseCmapFormat4(input, offset)
	case 6:
		return parseCmapFormat6(input, offset)
	case 10:
		return parseCmapFormat10(input, offset)
	case 12:
		return parseCmapFormat12(input, offset)
	case 13:
		return parseCmapFormat13(input, offset)
	default:
		return nil, fmt.Errorf("unsupported cmap subtable format: %d", format)
	}
}

func parseCmapFormat0(input []byte, offset uint32) (cmap0, error) {
	if len(input) < int(offset)+6+256 {
		return nil, errors.New("invalid cmap subtable format 0 (EOF)")
	}

	chars := cmap0{}
	for x, index := range input[offset+6 : offset+6+256] {
		r := charmap.Macintosh.DecodeByte(byte(x))
		// The source rune r is not representable in the Macintosh-Roman encoding.
		if r != 0 {
			chars[r] = GID(index)
		}
	}
	return chars, nil
}

type cmap2 struct {
	subHeaders      []cmapFormat2SubHeader
	glyphIndexArray []uint16 // may safely by slice by [header.rangeIndex:header.rangeIndex+header.entryCount]
	subHeaderKeys   [256]uint16
	language        uint16
}

func parseCmapFormat2(input []byte, offset uint32) (out cmap2, err error) {
	const headerSize = 6 + 2*256
	if len(input) < int(offset)+headerSize {
		return out, errors.New("invalid cmap subtable format 2 (EOF)")
	}
	input = input[offset:]

	out.language = binary.BigEndian.Uint16(input[4:])
	var maxIndex uint16
	// find the maximum index possible
	for i := range out.subHeaderKeys {
		out.subHeaderKeys[i] = binary.BigEndian.Uint16(input[6+2*i:]) / 8 // in the file, it is index * 8
		if out.subHeaderKeys[i] > maxIndex {
			maxIndex = out.subHeaderKeys[i]
		}
	}

	// parse the subHeaders
	startGlyphIndexArray := headerSize + 8*int(maxIndex+1)
	if len(input) < startGlyphIndexArray {
		return out, errors.New("invalid cmap subtable format 2 (EOF)")
	}
	out.subHeaders = make([]cmapFormat2SubHeader, maxIndex+1)
	for i := range out.subHeaders {
		pos := headerSize + 8*i
		subHeader := &out.subHeaders[i]
		subHeader.firstCode = binary.BigEndian.Uint16(input[pos:])
		subHeader.entryCount = binary.BigEndian.Uint16(input[pos+2:])
		subHeader.idDelta = int16(binary.BigEndian.Uint16(input[pos+4:]))
		idRangeOffset := binary.BigEndian.Uint16(input[pos+6:])
		// convert the offset into an index for convenience
		startSliceOffset := pos + 6 + int(idRangeOffset)
		if startSliceOffset < startGlyphIndexArray || len(input) < startSliceOffset+2*int(subHeader.entryCount) {
			return out, fmt.Errorf("invalid cmap subtable format 2: invalid idRangeOffset %d", idRangeOffset)
		}
		subHeader.rangeIndex = (startSliceOffset - startGlyphIndexArray) / 2
	}

	// the first subHeader has special values :
	// "For the one-byte case with k = 0, the structure subHeaders[0]
	// will show firstCode = 0, entryCount = 256, and idDelta = 0."
	if out.subHeaders[0] != (cmapFormat2SubHeader{0, 256, 0, 0}) {
		return out, fmt.Errorf("invalid cmap format2 first subHeader record: %v", out.subHeaders[0])
	}

	// load the glyphIndexArray
	out.glyphIndexArray, err = parseUint16s(input[startGlyphIndexArray:], len(input[startGlyphIndexArray:])/2)

	return out, err
}

type cmapFormat2SubHeader struct {
	firstCode  uint16
	entryCount uint16
	idDelta    int16
	rangeIndex int
}

func parseCmapFormat4(input []byte, offset uint32) (cmap4, error) {
	const headerSize = 14
	if len(input) < int(offset)+headerSize {
		return nil, errors.New("invalid cmap subtable format 4 (EOF)")
	}
	input = input[offset:]

	segCount := int(binary.BigEndian.Uint16(input[6:]))
	if segCount&1 != 0 {
		return nil, errors.New("invalid cmap subtable format 4 (odd segment count)")
	}
	segCount /= 2

	input = input[headerSize:]
	eLength := 8*segCount + 2 // 2 is for the reservedPad field
	if len(input) < eLength {
		return nil, fmt.Errorf("invalid cmap subtable format 4: EOF for %d segments", segCount)
	}
	glyphIDArray := input[eLength:]

	entries := make(cmap4, segCount)
	for i := range entries {
		cm := cmapEntry16{
			end:   binary.BigEndian.Uint16(input[2*i:]),
			start: binary.BigEndian.Uint16(input[2+2*(segCount+i):]),
			delta: binary.BigEndian.Uint16(input[2+2*(2*segCount+i):]),
		}
		idRangeOffset := int(binary.BigEndian.Uint16(input[2+2*(3*segCount+i):]))

		// some fonts use 0xFFFF for idRangeOff for the last segment
		if cm.start != 0xFFFF && idRangeOffset != 0 {
			// we resolve the indexes
			cm.indexes = make([]gid, cm.end-cm.start+1)
			indexStart := idRangeOffset/2 + i - segCount
			if len(glyphIDArray) < 2*(indexStart+len(cm.indexes)) {
				return nil, errors.New("invalid cmap subtable format 4 glyphs array length")
			}
			for j := range cm.indexes {
				index := indexStart + j
				cm.indexes[j] = gid(binary.BigEndian.Uint16(glyphIDArray[2*index:]))
			}
		}

		entries[i] = cm
	}
	return entries, nil
}

func parseCmapFormat6(input []byte, offset uint32) (out cmap6or10, err error) {
	const headerSize = 10
	if len(input) < int(offset)+headerSize {
		return out, errors.New("invalid cmap subtable format 6 (EOF)")
	}
	input = input[offset:]

	out.firstCode = rune(binary.BigEndian.Uint16(input[6:]))
	entryCount := int(binary.BigEndian.Uint16(input[8:]))

	out.entries, err = parseUint16s(input[headerSize:], entryCount)
	return out, err
}

func parseCmapFormat10(input []byte, offset uint32) (out cmap6or10, err error) {
	const headerSize = 20
	if len(input) < int(offset)+headerSize {
		return out, errors.New("invalid cmap subtable format 6 (EOF)")
	}
	input = input[offset:]

	out.firstCode = rune(binary.BigEndian.Uint32(input[12:]))
	entryCount := int(binary.BigEndian.Uint32(input[16:]))

	out.entries, err = parseUint16s(input[headerSize:], entryCount)
	return out, err
}

func parseCmapFormat12(input []byte, offset uint32) (cmap12, error) {
	return parseCmapFormat12or13(input, offset)
}

func parseCmapFormat13(input []byte, offset uint32) (cmap13, error) {
	return parseCmapFormat12or13(input, offset)
}

func parseCmapFormat12or13(input []byte, offset uint32) ([]cmapEntry32, error) {
	const headerSize = 16
	if len(input) < int(offset)+headerSize {
		return nil, errors.New("invalid cmap subtable format 12 (EOF)")
	}
	input = input[offset:]
	// length := binary.BigEndian.Uint32(bufHeader[4:])
	numGroups := int(binary.BigEndian.Uint32(input[12:]))

	if len(input) < headerSize+12*numGroups {
		return nil, errors.New("invalid cmap subtable format 12 (EOF)")
	}

	entries := make([]cmapEntry32, numGroups)
	for i := range entries {
		entries[i] = cmapEntry32{
			start: binary.BigEndian.Uint32(input[headerSize+0+12*i:]),
			end:   binary.BigEndian.Uint32(input[headerSize+4+12*i:]),
			value: binary.BigEndian.Uint32(input[headerSize+8+12*i:]),
		}
	}
	return entries, nil
}

// if indexes is nil, delta is used
type cmapEntry16 struct {
	// we prefere not to keep a link to a buffer (via an offset)
	// and eagerly resolve it
	indexes    []gid // length end - start + 1
	end, start uint16
	delta      uint16 // arithmetic modulo 0xFFFF
}

type cmapEntry32 struct {
	start, end, value uint32
}

func parseCmapFormat14(data []byte, offset uint32) (unicodeVariations, error) {
	if len(data) < int(offset)+10 {
		return nil, errors.New("invalid cmap subtable format 14 (EOF)")
	}
	data = data[offset:]
	count := binary.BigEndian.Uint32(data[6:])

	if len(data) < 10+int(count)*11 {
		return nil, errors.New("invalid cmap subtable format 14 (EOF)")
	}
	out := make(unicodeVariations, count)
	var err error
	for i := range out {
		out[i].varSelector = parseUint24(data[10+11*i:])

		offsetDefault := binary.BigEndian.Uint32(data[10+11*i+3:])
		if offsetDefault != 0 {
			out[i].defaultUVS, err = parseUnicodeRanges(data, offsetDefault)
			if err != nil {
				return nil, err
			}
		}

		offsetNonDefault := binary.BigEndian.Uint32(data[10+11*i+7:])
		if offsetNonDefault != 0 {
			out[i].nonDefaultUVS, err = parseUVSMappings(data, offsetNonDefault)
			if err != nil {
				return nil, err
			}
		}
	}

	return out, nil
}

type variationSelector struct {
	defaultUVS    []unicodeRange
	nonDefaultUVS []uvsMapping
	varSelector   rune
}

const (
	variantNotFound = iota
	variantUseDefault
	variantFound
)

func (vs variationSelector) getGlyph(r rune) (GID, uint8) {
	// binary search
	for i, j := 0, len(vs.defaultUVS); i < j; {
		h := i + (j-i)/2
		entry := vs.defaultUVS[h]
		if r < entry.start {
			j = h
		} else if entry.start+rune(entry.additionalCount) < r {
			i = h + 1
		} else {
			return 0, variantUseDefault
		}
	}

	for i, j := 0, len(vs.nonDefaultUVS); i < j; {
		h := i + (j-i)/2
		entry := vs.nonDefaultUVS[h].unicode
		if r < entry {
			j = h
		} else if entry < r {
			i = h + 1
		} else {
			return GID(vs.nonDefaultUVS[h].glyphID), variantFound
		}
	}

	return 0, variantNotFound
}

type unicodeRange struct {
	start           rune
	additionalCount uint8 // 0 for a singleton range
}

func parseUnicodeRanges(data []byte, offset uint32) ([]unicodeRange, error) {
	if len(data) < int(offset)+4 {
		return nil, errors.New("invalid unicode ranges (EOF)")
	}
	count := binary.BigEndian.Uint32(data[offset:])
	if len(data) < int(offset)+4+4*int(count) {
		return nil, errors.New("invalid unicode ranges (EOF)")
	}
	data = data[offset+4:]
	out := make([]unicodeRange, count)
	for i := range out {
		out[i].start = parseUint24(data[4*i:])
		out[i].additionalCount = data[4*i+3]
	}
	return out, nil
}

type uvsMapping struct {
	unicode rune
	glyphID gid
}

func parseUVSMappings(data []byte, offset uint32) ([]uvsMapping, error) {
	if len(data) < int(offset)+4 {
		return nil, errors.New("invalid UVS mappings (EOF)")
	}
	count := binary.BigEndian.Uint32(data[offset:])
	if len(data) < int(offset)+4+5*int(count) {
		return nil, errors.New("invalid UVS mappings (EOF)")
	}
	data = data[offset+4:]
	out := make([]uvsMapping, count)
	for i := range out {
		out[i].unicode = parseUint24(data[5*i:])
		out[i].glyphID = gid(binary.BigEndian.Uint16(data[5*i+3:]))
	}
	return out, nil
}

// same as binary.BigEndian.Uint32, but for 24 bit uint
func parseUint24(b []byte) rune {
	_ = b[2] // BCE
	return rune(b[0])<<16 | rune(b[1])<<8 | rune(b[2])
}
