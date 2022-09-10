package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// TableGSUB is the Glyph Substitution (GSUB) table.
// It provides data for substition of glyphs for appropriate rendering of scripts,
// such as cursively-connecting forms in Arabic script,
// or for advanced typographic effects, such as ligatures.
type TableGSUB struct {
	Lookups []LookupGSUB
	TableLayout
}

func parseTableGSUB(data []byte) (out TableGSUB, err error) {
	tableLayout, lookups, err := parseTableLayout(data)
	if err != nil {
		return out, err
	}
	out = TableGSUB{
		TableLayout: tableLayout,
		Lookups:     make([]LookupGSUB, len(lookups)),
	}
	for i, l := range lookups {
		out.Lookups[i], err = l.parseGSUB(uint16(len(lookups)))
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

// GSUBType identifies the kind of lookup format, for GSUB tables.
type GSUBType uint16

const (
	GSUBSingle    GSUBType = 1 + iota // Single (format 1.1 1.2)	Replace one glyph with one glyph
	GSUBMultiple                      // Multiple (format 2.1)	Replace one glyph with more than one glyph
	GSUBAlternate                     // Alternate (format 3.1)	Replace one glyph with one of many glyphs
	GSUBLigature                      // Ligature (format 4.1)	Replace multiple glyphs with one glyph
	GSUBContext                       // Context (format 5.1 5.2 5.3)	Replace one or more glyphs in context
	GSUBChaining                      // Chaining Context (format 6.1 6.2 6.3)	Replace one or more glyphs in chained context
	// Extension Substitution (format 7.1) Extension mechanism for other substitutions
	// The table pointed at by this lookup is returned instead.
	gsubExtension
	GSUBReverse // Reverse chaining context single (format 8.1)
)

// GSUBSubtable is one of the subtables of a
// GSUB lookup.
type GSUBSubtable struct {
	// For SubChaining - Format 3, its the coverage of the first input.
	Coverage Coverage
	Data     interface{ Type() GSUBType }
}

// LookupGSUB is a lookup for GSUB tables.
type LookupGSUB struct {
	// After successful parsing, all subtables have the `GSUBType`.
	Subtables []GSUBSubtable
	LookupOptions
	Type GSUBType
}

// interpret the lookup as a GSUB lookup
// lookupLength is used to sanitize nested lookups
func (header lookup) parseGSUB(lookupListLength uint16) (out LookupGSUB, err error) {
	out.Type = GSUBType(header.kind)
	out.LookupOptions = header.LookupOptions

	out.Subtables = make([]GSUBSubtable, len(header.subtableOffsets))
	for i, offset := range header.subtableOffsets {
		out.Subtables[i], err = parseGSUBSubtable(header.data, int(offset), out.Type, lookupListLength)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func parseGSUBSubtable(data []byte, offset int, kind GSUBType, lookupListLength uint16) (out GSUBSubtable, err error) {
	// read the format and coverage
	if offset+4 >= len(data) {
		return out, fmt.Errorf("invalid lookup subtable offset %d", offset)
	}
	format := binary.BigEndian.Uint16(data[offset:])

	// almost all table have a coverage offset, right after the format; special case the others
	// see below for the coverage
	if kind == gsubExtension || (kind == GSUBChaining || kind == GSUBContext) && format == 3 {
		out.Coverage = CoverageList{}
	} else {
		covOffset := uint32(binary.BigEndian.Uint16(data[offset+2:])) // relative to the subtable
		out.Coverage, err = parseCoverage(data[offset:], covOffset)
		if err != nil {
			return out, fmt.Errorf("invalid GSUB table (format %d-%d): %s", kind, format, err)
		}
	}

	// read the actual lookup
	switch kind {
	case GSUBSingle:
		out.Data, err = parseSingleSub(format, data[offset:])
	case GSUBMultiple:
		out.Data, err = parseMultipleSub(data[offset:], out.Coverage)
	case GSUBAlternate:
		out.Data, err = parseAlternateSub(data[offset:], out.Coverage)
	case GSUBLigature:
		out.Data, err = parseLigatureSub(data[offset:], out.Coverage)
	case GSUBContext:
		out.Data, err = parseSequenceContextSub(format, data[offset:], lookupListLength, &out.Coverage)
	case GSUBChaining:
		out.Data, err = parseChainedSequenceContextSub(format, data[offset:], lookupListLength, &out.Coverage)
	case gsubExtension:
		out, err = parseExtensionSub(data[offset:], lookupListLength)
	case GSUBReverse:
		out.Data, err = parseReverseChainedSequenceContextSub(data[offset:], out.Coverage)
	default:
		return out, fmt.Errorf("unsupported gsub lookup type %d", kind)
	}
	return out, err
}

func (GSUBSingle1) Type() GSUBType { return GSUBSingle }
func (GSUBSingle2) Type() GSUBType { return GSUBSingle }

// data starts at the subtable (but format has already been read)
func parseSingleSub(format uint16, data []byte) (out interface{ Type() GSUBType }, err error) {
	switch format {
	case 1:
		return parseSingleSub1(data)
	case 2:
		return parseSingleSub2(data)
	default:
		return nil, fmt.Errorf("unsupported single substitution format: %d", format)
	}
}

// GSUBSingle1 is a Single Substitution Format 1, expressed as a delta
// from the coverage.
type GSUBSingle1 int16

// data is at the begining of the subtable
func parseSingleSub1(data []byte) (GSUBSingle1, error) {
	if len(data) < 6 {
		return 0, errors.New("invalid single subsitution table (format 1)")
	}
	// format and coverage already read
	delta := GSUBSingle1(binary.BigEndian.Uint16(data[4:]))
	return delta, nil
}

// GSUBSingle2 is a Single Substitution Format 2, expressed as substitutes
type GSUBSingle2 []GID

// data is at the begining of the subtable
func parseSingleSub2(data []byte) (GSUBSingle2, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid single subsitution table (format 2)")
	}
	// format and coverage already read
	glyphCount := binary.BigEndian.Uint16(data[4:])
	if len(data) < 6+int(glyphCount)*2 {
		return nil, errors.New("invalid single subsitution table (format 2)")
	}
	out := make(GSUBSingle2, glyphCount)
	for i := range out {
		out[i] = GID(binary.BigEndian.Uint16(data[6+2*i:]))
	}
	return out, nil
}

type GSUBMultiple1 [][]GID

func (GSUBMultiple1) Type() GSUBType { return GSUBMultiple }

// data starts at the subtable (but format has already been read)
func parseMultipleSub(data []byte, cov Coverage) (GSUBMultiple1, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid multiple subsitution table")
	}

	// format and coverage already processed
	count := binary.BigEndian.Uint16(data[4:])

	// check length conformance
	if cov.Size() != int(count) {
		return nil, errors.New("invalid multiple subsitution table")
	}

	if 6+int(count)*2 > len(data) {
		return nil, fmt.Errorf("invalid multiple subsitution table")
	}

	out := make(GSUBMultiple1, count)
	var err error
	for i := range out {
		offset := binary.BigEndian.Uint16(data[6+2*i:])
		if int(offset) > len(data) {
			return out, errors.New("invalid multiple subsitution table")
		}
		out[i], err = parseMultipleSet(data[offset:])
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func parseMultipleSet(data []byte) ([]GID, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid multiple subsitution table")
	}
	count := binary.BigEndian.Uint16(data)
	if 2+int(count)*2 > len(data) {
		return nil, fmt.Errorf("invalid multiple subsitution table")
	}
	out := make([]GID, count)
	for i := range out {
		out[i] = GID(binary.BigEndian.Uint16(data[2+2*i:]))
	}
	return out, nil
}

type GSUBAlternate1 [][]GID

func (GSUBAlternate1) Type() GSUBType { return GSUBAlternate }

// data starts at the subtable (but format has already been read)
func parseAlternateSub(data []byte, cov Coverage) (GSUBAlternate1, error) {
	out, err := parseMultipleSub(data, cov)
	if err != nil {
		return nil, errors.New("invalid alternate substitution table")
	}
	return GSUBAlternate1(out), nil
}

// GSUBLigature1 stores one ligature set per glyph in the coverage.
type GSUBLigature1 [][]LigatureGlyph

func (GSUBLigature1) Type() GSUBType { return GSUBLigature }

func parseLigatureSub(data []byte, cov Coverage) (GSUBLigature1, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid ligature subsitution table")
	}

	// format and coverage already processed
	count := binary.BigEndian.Uint16(data[4:])

	// check length conformance
	if cov.Size() != int(count) {
		return nil, errors.New("invalid ligature subsitution table")
	}

	if 6+int(count)*2 > len(data) {
		return nil, fmt.Errorf("invalid ligature subsitution table")
	}
	out := make([][]LigatureGlyph, count)
	var err error
	for i := range out {
		ligSetOffset := binary.BigEndian.Uint16(data[6+2*i:])
		if int(ligSetOffset) > len(data) {
			return out, errors.New("invalid ligature subsitution table")
		}
		out[i], err = parseLigatureSet(data[ligSetOffset:])
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

type LigatureGlyph struct {
	// Glyphs composing the ligature, starting after the
	// implicit first glyph, given in the coverage of the
	// SubstitutionLigature table
	Components []uint16
	// Output ligature glyph
	Glyph GID
}

// Matches tests if the ligature should be applied on `glyphsFromSecond`,
// which starts from the second glyph.
func (l LigatureGlyph) Matches(glyphsFromSecond []GID) bool {
	if len(glyphsFromSecond) != len(l.Components) {
		return false
	}
	for i, g := range glyphsFromSecond {
		if g != GID(l.Components[i]) {
			return false
		}
	}
	return true
}

// data is at the begining of the ligature set table
func parseLigatureSet(data []byte) ([]LigatureGlyph, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid ligature set table")
	}
	count := binary.BigEndian.Uint16(data)
	out := make([]LigatureGlyph, count)
	var err error
	for i := range out {
		ligOffset := binary.BigEndian.Uint16(data[2+2*i:])
		if int(ligOffset)+4 > len(data) {
			return nil, errors.New("invalid ligature set table")
		}
		out[i].Glyph = GID(binary.BigEndian.Uint16(data[ligOffset:]))
		ligCount := binary.BigEndian.Uint16(data[ligOffset+2:])
		if ligCount == 0 {
			return nil, errors.New("invalid ligature set table")
		}
		out[i].Components, err = parseUint16s(data[ligOffset+4:], int(ligCount)-1)
		if err != nil {
			return nil, fmt.Errorf("invalid ligature set table: %s", err)
		}
	}
	return out, nil
}

type (
	GSUBContext1 LookupContext1
	GSUBContext2 LookupContext2
	GSUBContext3 LookupContext3
)

func (GSUBContext1) Type() GSUBType { return GSUBContext }
func (GSUBContext2) Type() GSUBType { return GSUBContext }
func (GSUBContext3) Type() GSUBType { return GSUBContext }

// lookupLength is used to sanitize lookup indexes.
// cov is used for ContextFormat3
func parseSequenceContextSub(format uint16, data []byte, lookupLength uint16, cov *Coverage) (interface{ Type() GSUBType }, error) {
	switch format {
	case 1:
		out, err := parseSequenceContext1(data, lookupLength)
		return GSUBContext1(out), err
	case 2:
		out, err := parseSequenceContext2(data, lookupLength)
		return GSUBContext2(out), err
	case 3:
		out, err := parseSequenceContext3(data, lookupLength)
		if len(out.Coverages) != 0 {
			*cov = out.Coverages[0]
		}
		return GSUBContext3(out), err
	default:
		return nil, fmt.Errorf("unsupported sequence context format %d", format)
	}
}

type (
	GSUBChainedContext1 LookupChainedContext1
	GSUBChainedContext2 LookupChainedContext2
	GSUBChainedContext3 LookupChainedContext3
)

func (GSUBChainedContext1) Type() GSUBType { return GSUBChaining }
func (GSUBChainedContext2) Type() GSUBType { return GSUBChaining }
func (GSUBChainedContext3) Type() GSUBType { return GSUBChaining }

// lookupLength is used to sanitize lookup indexes.
// cov is used for ContextFormat3
func parseChainedSequenceContextSub(format uint16, data []byte, lookupLength uint16, cov *Coverage) (interface{ Type() GSUBType }, error) {
	switch format {
	case 1:
		out, err := parseChainedSequenceContext1(data, lookupLength)
		return GSUBChainedContext1(out), err
	case 2:
		out, err := parseChainedSequenceContext2(data, lookupLength)
		return GSUBChainedContext2(out), err
	case 3:
		out, err := parseChainedSequenceContext3(data, lookupLength)
		if len(out.Input) != 0 {
			*cov = out.Input[0]
		}
		return GSUBChainedContext3(out), err
	default:
		return nil, fmt.Errorf("unsupported sequence context format %d", format)
	}
}

// returns the extension subtable instead
func parseExtensionSub(data []byte, lookupListLength uint16) (GSUBSubtable, error) {
	if len(data) < 8 {
		return GSUBSubtable{}, errors.New("invalid extension substitution table")
	}
	extensionType := GSUBType(binary.BigEndian.Uint16(data[2:]))
	offset := binary.BigEndian.Uint32(data[4:])

	if extensionType == gsubExtension {
		return GSUBSubtable{}, errors.New("invalid extension substitution table")
	}

	return parseGSUBSubtable(data, int(offset), extensionType, lookupListLength)
}

type GSUBReverseChainedContext1 struct {
	Backtrack   []Coverage
	Lookahead   []Coverage
	Substitutes []GID
}

func (GSUBReverseChainedContext1) Type() GSUBType { return GSUBReverse }

func parseReverseChainedSequenceContextSub(data []byte, cov Coverage) (out GSUBReverseChainedContext1, err error) {
	if len(data) < 6 {
		return out, errors.New("invalid reversed chained sequence context format 3 table")
	}
	covCount := binary.BigEndian.Uint16(data[4:])
	out.Backtrack = make([]Coverage, covCount)
	for i := range out.Backtrack {
		covOffset := uint32(binary.BigEndian.Uint16(data[6+2*i:]))
		out.Backtrack[i], err = parseCoverage(data, covOffset)
		if err != nil {
			return out, err
		}
	}
	endBacktrack := 6 + 2*int(covCount)

	if len(data) < endBacktrack+2 {
		return out, errors.New("invalid reversed chained sequence context format 3 table")
	}
	covCount = binary.BigEndian.Uint16(data[endBacktrack:])
	out.Lookahead = make([]Coverage, covCount)
	for i := range out.Lookahead {
		covOffset := uint32(binary.BigEndian.Uint16(data[endBacktrack+2+2*i:]))
		out.Lookahead[i], err = parseCoverage(data, covOffset)
		if err != nil {
			return out, err
		}
	}
	endLookahead := endBacktrack + 2 + 2*int(covCount)

	if len(data) < endBacktrack+2 {
		return out, errors.New("invalid reversed chained sequence context format 3 table")
	}
	glyphCount := binary.BigEndian.Uint16(data[endLookahead:])

	if cov.Size() != int(glyphCount) {
		return out, errors.New("invalid reversed chained sequence context format 3 table")
	}

	if len(data) < endLookahead+2+2*int(glyphCount) {
		return out, errors.New("invalid reversed chained sequence context format 3 table")
	}
	out.Substitutes = make([]GID, glyphCount)
	for i := range out.Substitutes {
		out.Substitutes[i] = GID(binary.BigEndian.Uint16(data[endLookahead+2+2*i:]))
	}
	return out, err
}
