// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package cff

// code is adapted from golang.org/x/image/font/sfnt

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/go-text/typesetting/opentype/api"
	ps "github.com/go-text/typesetting/opentype/api/font/cff/interpreter"
	"github.com/go-text/typesetting/opentype/tables"
)

var (
	errUnsupportedCFFVersion       = errors.New("unsupported CFF version")
	errUnsupportedCFFFDSelectTable = errors.New("unsupported FD Select version")
)

// Font represents a parsed CFF font.
type Font struct {
	userStrings userStrings
	fdSelect    fdSelect // only valid for CIDFonts
	charset     []uint16 // indexed by glyph ID

	cidFontName string

	// Charstrings contains the actual glyph definition.
	// It has a length of numGlyphs and is indexed by glyph ID.
	// See `LoadGlyph` for a way to intepret the glyph data.
	Charstrings [][]byte

	fontName    []byte // name from the Name INDEX
	globalSubrs [][]byte

	// array of length 1 for non CIDFonts
	// For CIDFonts, it can be safely indexed by `fdSelect` output
	localSubrs [][][]byte
}

// Parse parses a .cff font file.
// Although CFF enables multiple font or CIDFont programs to be bundled together in a
// single file, embedded CFF font file in PDF or in TrueType/OpenType fonts
// shall consist of exactly one font or CIDFont. Thus, this function
// returns an error if the file contains more than one font.
func Parse(file []byte) (*Font, error) {
	// read 4 bytes to check if its a supported CFF file
	if L := len(file); L < 4 {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", 4, L)
	}
	if file[0] != 1 || file[1] != 0 || file[2] != 4 {
		return nil, errUnsupportedCFFVersion
	}
	p := cffParser{src: file, offset: 4}
	out, err := p.parse()
	if err != nil {
		return nil, err
	}

	if len(out) > 1 {
		return nil, errors.New("only one font is allowed CFF table")
	}

	return &out[0], nil
}

// GlyphName returns the name of the glyph or an empty string if not found.
func (f *Font) GlyphName(glyph api.GID) string {
	if f.fdSelect != nil || int(glyph) >= len(f.charset) {
		return ""
	}
	out, _ := f.userStrings.getString(f.charset[glyph])
	return out
}

// since SID = 0 means .notdef, we use a reserved value
// to mean unset
const unsetSID = uint16(0xFFFF)

type userStrings [][]byte

// return either the predefined string or the user defined one
func (u userStrings) getString(sid uint16) (string, error) {
	if sid == unsetSID {
		return "", nil
	}
	if sid < 391 {
		return stdStrings[sid], nil
	}
	sid -= 391
	if int(sid) >= len(u) {
		return "", fmt.Errorf("invalid glyph index %d", sid)
	}
	return string(u[sid]), nil
}

// Compact Font Format (CFF) fonts are written in PostScript, a stack-based
// programming language.
//
// A fundamental concept is a DICT, or a key-value map, expressed in reverse
// Polish notation. For example, this sequence of operations:
//   - push the number 379
//   - version operator
//   - push the number 392
//   - Notice operator
//   - etc
//   - push the number 100
//   - push the number 0
//   - push the number 500
//   - push the number 800
//   - FontBBox operator
//   - etc
//
// defines a DICT that maps "version" to the String ID (SID) 379, "Notice" to
// the SID 392, "FontBBox" to the four numbers [100, 0, 500, 800], etc.
//
// The first 391 String IDs (starting at 0) are predefined as per the CFF spec
// Appendix A, in 5176.CFF.pdf referenced below. For example, 379 means
// "001.000". String ID 392 is not predefined, and is mapped by a separate
// structure, the "String INDEX", inside the CFF data. (String ID 391 is also
// not predefined. Specifically for go-opentype-testdata/data/toys/CFFTest.otf, 391 means
// "uni4E2D", as this font contains a glyph for U+4E2D).
//
// The actual glyph vectors are similarly encoded (in PostScript), in a format
// called Type 2 Charstrings. The wire encoding is similar to but not exactly
// the same as CFF's. For example, the byte 0x05 means FontBBox for CFF DICTs,
// but means rlineto (relative line-to) for Type 2 Charstrings. See
// 5176.CFF.pdf Appendix H and 5177.Type2.pdf Appendix A in the PDF files
// referenced below.
//
// The relevant specifications are:
//   - http://wwwimages.adobe.com/content/dam/Adobe/en/devnet/font/pdfs/5176.CFF.pdf
//   - http://wwwimages.adobe.com/content/dam/Adobe/en/devnet/font/pdfs/5177.Type2.pdf
type cffParser struct {
	src    []byte // whole input
	offset int    // current position
}

func (p *cffParser) parse() ([]Font, error) {
	// header was checked prior to this call

	// Parse the Name INDEX.
	fontNames, err := p.parseNames()
	if err != nil {
		return nil, err
	}

	topDicts, err := p.parseTopDicts()
	if err != nil {
		return nil, err
	}
	// 5176.CFF.pdf section 8 "Top DICT INDEX" says that the count here
	// should match the count of the Name INDEX
	if len(topDicts) != len(fontNames) {
		return nil, fmt.Errorf("top DICT length doest not match Names (%d, %d)", len(topDicts),
			len(fontNames))
	}

	// parse the String INDEX.
	strs, err := p.parseUserStrings()
	if err != nil {
		return nil, err
	}

	out := make([]Font, len(topDicts))

	// use the strings to fetch the PSInfo
	for i, topDict := range topDicts {
		out[i].fontName = fontNames[i]
		out[i].userStrings = strs

		// skip PSInfo, and cidFontName

		out[i].cidFontName, err = strs.getString(topDict.cidFontName)
		if err != nil {
			return nil, err
		}
	}

	// Parse the Global Subrs [Subroutines] INDEX,
	// shared among all fonts.
	globalSubrs, err := p.parseIndex()
	if err != nil {
		return nil, err
	}

	for i, topDict := range topDicts {
		out[i].globalSubrs = globalSubrs

		// Parse the CharStrings INDEX, whose location was found in the Top DICT.
		if err = p.seek(topDict.charStringsOffset); err != nil {
			return nil, err
		}
		out[i].Charstrings, err = p.parseIndex()
		if err != nil {
			return nil, err
		}
		numGlyphs := uint16(len(out[i].Charstrings))

		out[i].charset, err = p.parseCharset(topDict.charsetOffset, numGlyphs)
		if err != nil {
			return nil, err
		}

		// skip encoding

		if !topDict.isCIDFont {
			// Parse the Private DICT, whose location was found in the Top DICT.
			var localSubrs [][]byte
			localSubrs, err = p.parsePrivateDICT(topDict.privateDictOffset, topDict.privateDictLength)
			if err != nil {
				return nil, err
			}
			out[i].localSubrs = [][][]byte{localSubrs}
		} else {
			// Parse the Font Dict Select data, whose location was found in the Top
			// DICT.
			out[i].fdSelect, err = p.parseFDSelect(topDict.fdSelect, numGlyphs)
			if err != nil {
				return nil, err
			}
			indexExtent := out[i].fdSelect.extent()

			// Parse the Font Dicts. Each one contains its own Private DICT.
			if err = p.seek(topDict.fdArray); err != nil {
				return nil, err
			}
			topDicts, err := p.parseTopDicts()
			if err != nil {
				return nil, err
			}
			if len(topDicts) < indexExtent {
				return nil, fmt.Errorf("invalid number of font dicts: %d (for %d)",
					len(topDicts), indexExtent)
			}
			multiSubrs := make([][][]byte, len(topDicts))
			for i, topDict := range topDicts {
				multiSubrs[i], err = p.parsePrivateDICT(topDict.privateDictOffset, topDict.privateDictLength)
				if err != nil {
					return nil, err
				}
			}
			out[i].localSubrs = multiSubrs
		}
	}

	return out, nil
}

func (p *cffParser) parseTopDicts() ([]topDictData, error) {
	// Parse the Top DICT INDEX.
	instructions, err := p.parseIndex()
	if err != nil {
		return nil, err
	}

	out := make([]topDictData, len(instructions)) // guarded by uint16 max size
	var psi ps.Machine
	for i, buf := range instructions {
		topDict := &out[i]

		// set default value before parsing
		topDict.underlinePosition = -100
		topDict.underlineThickness = 50
		topDict.version = unsetSID
		topDict.notice = unsetSID
		topDict.fullName = unsetSID
		topDict.familyName = unsetSID
		topDict.weight = unsetSID
		topDict.cidFontName = unsetSID

		if err = psi.Run(buf, nil, nil, topDict); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// parse the general form of an index
func (p *cffParser) parseIndex() ([][]byte, error) {
	count, offSize, err := p.parseIndexHeader()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	out := make([][]byte, count)

	stringsLocations := make([]uint32, int(count)+1)
	if err := p.parseIndexLocations(stringsLocations, offSize); err != nil {
		return nil, err
	}

	for i := range out {
		length := stringsLocations[i+1] - stringsLocations[i]
		buf, err := p.read(int(length))
		if err != nil {
			return nil, err
		}
		out[i] = buf
	}
	return out, nil
}

// parse the Name INDEX
func (p *cffParser) parseNames() ([][]byte, error) {
	return p.parseIndex()
}

// parse the String INDEX
func (p *cffParser) parseUserStrings() (userStrings, error) {
	index, err := p.parseIndex()
	return userStrings(index), err
}

// Parse the charset data, whose location was found in the Top DICT.
func (p *cffParser) parseCharset(charsetOffset int32, numGlyphs uint16) ([]uint16, error) {
	// Predefined charset may have offset of 0 to 2 // Table 22
	var charset []uint16
	switch charsetOffset {
	case 0: // ISOAdobe
		charset = charsetISOAdobe[:]
	case 1: // Expert
		charset = charsetExpert[:]
	case 2: // ExpertSubset
		charset = charsetExpertSubset[:]
	default: // custom
		if err := p.seek(charsetOffset); err != nil {
			return nil, err
		}
		buf, err := p.read(1)
		if err != nil {
			return nil, err
		}
		charset = make([]uint16, numGlyphs)
		switch buf[0] { // format
		case 0:
			buf, err = p.read(2 * (int(numGlyphs) - 1)) // ".notdef" is omited, and has an implicit SID of 0
			if err != nil {
				return nil, err
			}
			for i := uint16(1); i < numGlyphs; i++ {
				charset[i] = binary.BigEndian.Uint16(buf[2*i-2:])
			}
		case 1:
			for i := uint16(1); i < numGlyphs; {
				buf, err = p.read(3)
				if err != nil {
					return nil, err
				}
				first, nLeft := binary.BigEndian.Uint16(buf), uint16(buf[2])
				for j := uint16(0); j <= nLeft && i < numGlyphs; j++ {
					charset[i] = first + j
					i++
				}
			}
		case 2:
			for i := uint16(1); i < numGlyphs; {
				buf, err = p.read(4)
				if err != nil {
					return nil, err
				}
				first, nLeft := binary.BigEndian.Uint16(buf), binary.BigEndian.Uint16(buf[2:])
				for j := uint16(0); j <= nLeft && i < numGlyphs; j++ {
					charset[i] = first + j
					i++
				}
			}
		default:
			return nil, fmt.Errorf("invalid custom charset format %d", buf[0])
		}
	}
	return charset, nil
}

// fdSelect holds a CFF font's Font Dict Select data.
type fdSelect interface {
	fontDictIndex(glyph tables.GlyphID) (byte, error)
	// return the maximum index + 1 (it's the length of an array
	// which can be safely indexed by the indexes)
	extent() int
}

type fdSelect0 []byte

func (fds fdSelect0) fontDictIndex(glyph tables.GlyphID) (byte, error) {
	if int(glyph) >= len(fds) {
		return 0, errors.New("invalid glyph index")
	}
	return fds[glyph], nil
}

func (fds fdSelect0) extent() int {
	max := -1
	for _, b := range fds {
		if int(b) > max {
			max = int(b)
		}
	}
	return max + 1
}

type range3 struct {
	first tables.GlyphID
	fd    byte
}

type fdSelect3 struct {
	ranges   []range3
	sentinel tables.GlyphID // = numGlyphs
}

func (fds fdSelect3) fontDictIndex(x tables.GlyphID) (byte, error) {
	lo, hi := 0, len(fds.ranges)
	for lo < hi {
		i := (lo + hi) / 2
		r := fds.ranges[i]
		xlo := r.first
		if x < xlo {
			hi = i
			continue
		}
		xhi := fds.sentinel
		if i < len(fds.ranges)-1 {
			xhi = fds.ranges[i+1].first
		}
		if xhi <= x {
			lo = i + 1
			continue
		}
		return r.fd, nil
	}
	return 0, errors.New("invalid glyph index")
}

func (fds fdSelect3) extent() int {
	max := -1
	for _, b := range fds.ranges {
		if int(b.fd) > max {
			max = int(b.fd)
		}
	}
	return max + 1
}

// parseFDSelect parses the Font Dict Select data as per 5176.CFF.pdf section
// 19 "FDSelect".
func (p *cffParser) parseFDSelect(offset int32, numGlyphs uint16) (fdSelect, error) {
	if err := p.seek(offset); err != nil {
		return nil, err
	}
	buf, err := p.read(1)
	if err != nil {
		return nil, err
	}
	switch buf[0] { // format
	case 0:
		if len(p.src) < p.offset+int(numGlyphs) {
			return nil, errors.New("invalid FDSelect data")
		}
		return fdSelect0(p.src[p.offset : p.offset+int(numGlyphs)]), nil
	case 3:
		buf, err = p.read(2)
		if err != nil {
			return nil, err
		}
		numRanges := binary.BigEndian.Uint16(buf)
		if len(p.src) < p.offset+3*int(numRanges)+2 {
			return nil, errors.New("invalid FDSelect data")
		}
		out := fdSelect3{
			sentinel: tables.GlyphID(numGlyphs),
			ranges:   make([]range3, numRanges),
		}
		for i := range out.ranges {
			// 	buf holds the range [xlo, xhi).
			out.ranges[i].first = tables.GlyphID(binary.BigEndian.Uint16(p.src[p.offset+3*i:]))
			out.ranges[i].fd = p.src[p.offset+3*i+2]
		}
		return out, nil
	}
	return nil, errUnsupportedCFFFDSelectTable
}

// Parse Private DICT and the Local Subrs [Subroutines] INDEX
func (p *cffParser) parsePrivateDICT(offset, length int32) ([][]byte, error) {
	if length == 0 {
		return nil, nil
	}
	if err := p.seek(offset); err != nil {
		return nil, err
	}
	buf, err := p.read(int(length))
	if err != nil {
		return nil, err
	}
	var (
		psi  ps.Machine
		priv privateDict
	)
	if err = psi.Run(buf, nil, nil, &priv); err != nil {
		return nil, err
	}

	if priv.subrsOffset == 0 {
		return nil, nil
	}

	// "The local subrs offset is relative to the beginning of the Private DICT data"
	if err = p.seek(offset + priv.subrsOffset); err != nil {
		return nil, errors.New("invalid local subroutines offset")
	}
	subrs, err := p.parseIndex()
	if err != nil {
		return nil, err
	}
	return subrs, nil
}

// read returns the n bytes from p.offset and advances p.offset by n.
func (p *cffParser) read(n int) ([]byte, error) {
	if n < 0 || len(p.src) < p.offset+n {
		return nil, errors.New("invalid CFF font file (EOF)")
	}
	out := p.src[p.offset : p.offset+n]
	p.offset += n
	return out, nil
}

func (p *cffParser) seek(offset int32) error {
	if offset < 0 || len(p.src) < int(offset) {
		return errors.New("invalid CFF font file (EOF)")
	}
	p.offset = int(offset)
	return nil
}

func bigEndian(b []byte) uint32 {
	switch len(b) {
	case 1:
		return uint32(b[0])
	case 2:
		return uint32(b[0])<<8 | uint32(b[1])
	case 3:
		return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
	case 4:
		return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	}
	panic("unreachable")
}

func (p *cffParser) parseIndexHeader() (count uint16, offSize int32, err error) {
	buf, err := p.read(2)
	if err != nil {
		return 0, 0, err
	}
	count = binary.BigEndian.Uint16(buf)
	// 5176.CFF.pdf section 5 "INDEX Data" says that "An empty INDEX is
	// represented by a count field with a 0 value and no additional fields.
	// Thus, the total size of an empty INDEX is 2 bytes".
	if count == 0 {
		return count, 0, nil
	}
	buf, err = p.read(1)
	if err != nil {
		return 0, 0, err
	}
	offSize = int32(buf[0])
	if offSize < 1 || 4 < offSize {
		return 0, 0, fmt.Errorf("invalid offset size %d", offSize)
	}
	return count, offSize, nil
}

func (p *cffParser) parseIndexLocations(dst []uint32, offSize int32) error {
	if len(dst) == 0 {
		return nil
	}
	buf, err := p.read(len(dst) * int(offSize))
	if err != nil {
		return err
	}

	prev := uint32(0)
	for i := range dst {
		loc := bigEndian(buf[:offSize])
		buf = buf[offSize:]

		// Locations are off by 1 byte. 5176.CFF.pdf section 5 "INDEX Data"
		// says that "Offsets in the offset array are relative to the byte that
		// precedes the object data... This ensures that every object has a
		// corresponding offset which is always nonzero".
		if loc == 0 {
			return errors.New("invalid CFF index locations (0)")
		}
		loc--

		// In the same paragraph, "Therefore the first element of the offset
		// array is always 1" before correcting for the off-by-1.
		if i == 0 {
			if loc != 0 {
				return errors.New("invalid CFF index locations (not 0)")
			}
		} else if loc < prev { // Check that locations are increasing
			return errors.New("invalid CFF index locations (not increasing)")
		}

		// Check that locations are in bounds.
		if uint32(len(p.src)-p.offset) < loc {
			return errors.New("invalid CFF index locations (out of bounds)")
		}

		dst[i] = uint32(p.offset) + loc
		prev = loc
	}
	return nil
}

// topDictData contains fields specific to the Top DICT context.
type topDictData struct {
	// SIDs, to be decoded using the string index
	version, notice, fullName, familyName, weight      uint16
	isFixedPitch                                       bool
	italicAngle, underlinePosition, underlineThickness float32
	charsetOffset                                      int32
	encodingOffset                                     int32
	charStringsOffset                                  int32
	fdArray                                            int32
	fdSelect                                           int32
	isCIDFont                                          bool
	cidFontName                                        uint16
	privateDictOffset                                  int32
	privateDictLength                                  int32
}

func (topDict *topDictData) Context() ps.Context { return ps.TopDict }

func (topDict *topDictData) Apply(state *ps.Machine, op ps.Operator) error {
	ops := topDictOperators[0]
	if op.IsEscaped {
		ops = topDictOperators[1]
	}
	if int(op.Operator) >= len(ops) {
		return fmt.Errorf("invalid operator %s in Top Dict", op)
	}
	opFunc := ops[op.Operator]
	if opFunc.run == nil {
		return fmt.Errorf("invalid operator %s in Top Dict", op)
	}
	if state.ArgStack.Top < opFunc.numPop {
		return fmt.Errorf("invalid number of arguments for operator %s in Top Dict", op)
	}
	err := opFunc.run(topDict, state)
	if err != nil {
		return err
	}
	err = state.ArgStack.PopN(opFunc.numPop)
	return err
}

// The Top DICT operators are defined by 5176.CFF.pdf Table 9 "Top DICT
// Operator Entries" and Table 10 "CIDFont Operator Extensions".
type topDictOperator struct {
	// run is the function that implements the operator. Nil means that we
	// ignore the operator, other than popping its arguments off the stack.
	run func(*topDictData, *ps.Machine) error

	// numPop is the number of stack values to pop. -1 means "array" and -2
	// means "delta" as per 5176.CFF.pdf Table 6 "Operand Types".
	numPop int32
}

func topDictNoOp(*topDictData, *ps.Machine) error { return nil }

var topDictOperators = [2][]topDictOperator{
	// 1-byte operators.
	{
		0: {func(t *topDictData, s *ps.Machine) error {
			t.version = s.ArgStack.Uint16()
			return nil
		}, +1 /*version*/},
		1: {func(t *topDictData, s *ps.Machine) error {
			t.notice = s.ArgStack.Uint16()
			return nil
		}, +1 /*Notice*/},
		2: {func(t *topDictData, s *ps.Machine) error {
			t.fullName = s.ArgStack.Uint16()
			return nil
		}, +1 /*FullName*/},
		3: {func(t *topDictData, s *ps.Machine) error {
			t.familyName = s.ArgStack.Uint16()
			return nil
		}, +1 /*FamilyName*/},
		4: {func(t *topDictData, s *ps.Machine) error {
			t.weight = s.ArgStack.Uint16()
			return nil
		}, +1 /*Weight*/},
		5:  {topDictNoOp, -1 /*FontBBox*/},
		13: {topDictNoOp, +1 /*UniqueID*/},
		14: {topDictNoOp, -1 /*XUID*/},
		15: {func(t *topDictData, s *ps.Machine) error {
			t.charsetOffset = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +1 /*charset*/},
		16: {func(t *topDictData, s *ps.Machine) error {
			t.encodingOffset = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +1 /*Encoding*/},
		17: {func(t *topDictData, s *ps.Machine) error {
			t.charStringsOffset = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +1 /*CharStrings*/},
		18: {func(t *topDictData, s *ps.Machine) error {
			t.privateDictLength = s.ArgStack.Vals[s.ArgStack.Top-2]
			t.privateDictOffset = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +2 /*Private*/},
	},
	// 2-byte operators. The first byte is the escape byte.
	{
		0: {topDictNoOp, +1 /*Copyright*/},
		1: {func(t *topDictData, s *ps.Machine) error {
			t.isFixedPitch = s.ArgStack.Vals[s.ArgStack.Top-1] == 1
			return nil
		}, +1 /*isFixedPitch*/},
		2: {func(t *topDictData, s *ps.Machine) error {
			t.italicAngle = s.ArgStack.Float()
			return nil
		}, +1 /*ItalicAngle*/},
		3: {func(t *topDictData, s *ps.Machine) error {
			t.underlinePosition = s.ArgStack.Float()
			return nil
		}, +1 /*UnderlinePosition*/},
		4: {func(t *topDictData, s *ps.Machine) error {
			t.underlineThickness = s.ArgStack.Float()
			return nil
		}, +1 /*UnderlineThickness*/},
		5: {topDictNoOp, +1 /*PaintType*/},
		6: {func(_ *topDictData, i *ps.Machine) error {
			if version := i.ArgStack.Vals[i.ArgStack.Top-1]; version != 2 {
				return fmt.Errorf("charstring type %d not supported", version)
			}
			return nil
		}, +1 /*CharstringType*/},
		7:  {topDictNoOp, -1 /*FontMatrix*/},
		8:  {topDictNoOp, +1 /*StrokeWidth*/},
		20: {topDictNoOp, +1 /*SyntheticBase*/},
		21: {topDictNoOp, +1 /*PostScript*/},
		22: {topDictNoOp, +1 /*BaseFontName*/},
		23: {topDictNoOp, -2 /*BaseFontBlend*/},
		30: {func(t *topDictData, _ *ps.Machine) error {
			t.isCIDFont = true
			return nil
		}, +3 /*ROS*/},
		31: {topDictNoOp, +1 /*CIDFontVersion*/},
		32: {topDictNoOp, +1 /*CIDFontRevision*/},
		33: {topDictNoOp, +1 /*CIDFontType*/},
		34: {topDictNoOp, +1 /*CIDCount*/},
		35: {topDictNoOp, +1 /*UIDBase*/},
		36: {func(t *topDictData, s *ps.Machine) error {
			t.fdArray = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +1 /*FDArray*/},
		37: {func(t *topDictData, s *ps.Machine) error {
			t.fdSelect = s.ArgStack.Vals[s.ArgStack.Top-1]
			return nil
		}, +1 /*FDSelect*/},
		38: {func(t *topDictData, s *ps.Machine) error {
			t.cidFontName = s.ArgStack.Uint16()
			return nil
		}, +1 /*FontName*/},
	},
}

// privateDict contains fields specific to the Private DICT context.
type privateDict struct {
	subrsOffset                  int32
	defaultWidthX, nominalWidthX int32
}

func (privateDict) Context() ps.Context { return ps.PrivateDict }

// The Private DICT operators are defined by 5176.CFF.pdf Table 23 "Private
// DICT Operators".
func (priv *privateDict) Apply(state *ps.Machine, op ps.Operator) error {
	if !op.IsEscaped { // 1-byte operators.
		switch op.Operator {
		case 6, 7, 8, 9: // "BlueValues" "OtherBlues" "FamilyBlues" "FamilyOtherBlues"
			return state.ArgStack.PopN(-2)
		case 10, 11: // "StdHW" "StdVW"
			return state.ArgStack.PopN(1)
		case 20: // "defaultWidthX"
			if state.ArgStack.Top < 1 {
				return errors.New("invalid stack size for 'defaultWidthX' in private Dict charstring")
			}
			priv.defaultWidthX = state.ArgStack.Vals[state.ArgStack.Top-1]
			return state.ArgStack.PopN(1)
		case 21: // "nominalWidthX"
			if state.ArgStack.Top < 1 {
				return errors.New("invalid stack size for 'nominalWidthX' in private Dict charstring")
			}
			priv.nominalWidthX = state.ArgStack.Vals[state.ArgStack.Top-1]
			return state.ArgStack.PopN(1)
		case 19: // "Subrs" pop 1
			if state.ArgStack.Top < 1 {
				return errors.New("invalid stack size for 'subrs' in private Dict charstring")
			}
			priv.subrsOffset = state.ArgStack.Vals[state.ArgStack.Top-1]
			return state.ArgStack.PopN(1)
		}
	} else { // 2-byte operators. The first byte is the escape byte.
		switch op.Operator {
		case 9, 10, 11, 14, 17, 18, 19: // "BlueScale" "BlueShift" "BlueFuzz" "ForceBold" "LanguageGroup" "ExpansionFactor" "initialRandomSeed"
			return state.ArgStack.PopN(1)
		case 12, 13: //  "StemSnapH"  "StemSnapV"
			return state.ArgStack.PopN(-2)
		}
	}
	return errors.New("invalid operand in private Dict charstring")
}
