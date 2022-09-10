package truetype

// parser of Apple AAT layout tables
// We dont support the deprecated 'mort' tables

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TableMorx []MorxChain

func parseTableMorx(data []byte, numGlyphs int) (TableMorx, error) {
	if len(data) < 8 {
		return nil, errors.New("invalid morx table (EOF)")
	}
	version := binary.BigEndian.Uint16(data)
	// unused
	nChains := binary.BigEndian.Uint32(data[4:])

	// "sanitize" before allocating
	if len(data) < int(nChains)*12 {
		return nil, errors.New("invalid morx table (EOF)")
	}
	currentOffset := 8
	out := make(TableMorx, nChains)
	for i := range out {
		if len(data) < currentOffset {
			return nil, errors.New("invalid morx table (EOF)")
		}
		var (
			size int
			err  error
		)
		out[i], size, err = parseMorxChain(version, data[currentOffset:], numGlyphs)
		if err != nil {
			return nil, err
		}
		currentOffset += size
	}
	return out, nil
}

type MorxChain struct {
	Features     []AATFeature
	Subtables    []MortxSubtable
	DefaultFlags uint32
}

func parseMorxChain(version uint16, data []byte, numGlyphs int) (out MorxChain, size int, err error) {
	switch version {
	case 1:
		return out, 0, fmt.Errorf("deprecated mort tables are not supported")
	case 2, 3:
		return parseMorxChain23(data, numGlyphs)
	default:
		return out, 0, fmt.Errorf("unsupported morx version %d", version)
	}
}

func parseMorxChain23(data []byte, numGlyphs int) (out MorxChain, size int, err error) {
	if len(data) < 16 {
		return out, 0, errors.New("invalid morx table (EOF)")
	}
	out.DefaultFlags = binary.BigEndian.Uint32(data)
	size = int(binary.BigEndian.Uint32(data[4:]))
	nFeatures := binary.BigEndian.Uint32(data[8:])
	nSubtables := binary.BigEndian.Uint32(data[12:])

	if len(data) < 12*int(nFeatures) {
		return out, 0, errors.New("invalid morx table (EOF)")
	}
	out.Features = make([]AATFeature, nFeatures)
	for i := range out.Features {
		out.Features[i].Type = binary.BigEndian.Uint16(data[16+12*i:])
		out.Features[i].Setting = binary.BigEndian.Uint16(data[16+12*i+2:])
		out.Features[i].EnableFlags = binary.BigEndian.Uint32(data[16+12*i+4:])
		out.Features[i].DisableFlags = binary.BigEndian.Uint32(data[16+12*i+8:])
	}

	// "sanitize" before allocating
	currentOffset := 16 + 12*int(nFeatures)
	if len(data) < currentOffset+12*int(nSubtables) { // at least
		return out, 0, errors.New("invalid morx table (EOF)")
	}
	out.Subtables = make([]MortxSubtable, nSubtables)
	var subtableLength int
	for i := range out.Subtables {
		if len(data) < currentOffset {
			return out, 0, errors.New("invalid morx table (EOF)")
		}
		out.Subtables[i], subtableLength, err = parseMorxSubtable(data[currentOffset:], numGlyphs)
		if err != nil {
			return out, 0, err
		}
		currentOffset += subtableLength
	}
	return out, size, nil
}

type AATFeature struct {
	Type, Setting uint16
	EnableFlags   uint32 // Flags for the settings that this feature and setting enables.
	DisableFlags  uint32 // Complement of flags for the settings that this feature and setting disable.
}

// MorxSubtableType indicates the kind of 'morx' subtable.
// See the constants.
type MorxSubtableType uint8

const (
	MorxRearrangement MorxSubtableType = iota
	MorxContextual
	MorxLigature
	_ // reserved
	MorxNonContextual
	MorxInsertion
)

type MortxSubtable struct {
	Data     interface{ Type() MorxSubtableType }
	Coverage uint8  // high byte of the coverage flag
	Flags    uint32 // Mask identifying which subtable this is.
}

// also returns the length of the subtable (in bytes)
func parseMorxSubtable(data []byte, numGlyphs int) (out MortxSubtable, length int, err error) {
	if len(data) < 12 {
		return out, 0, errors.New("invalid morx subtable (EOF)")
	}
	length = int(binary.BigEndian.Uint32(data))
	if len(data) < int(length) {
		return out, 0, errors.New("invalid morx subtable (EOF)")
	}
	out.Coverage = data[4]            // high order byte
	kind := MorxSubtableType(data[7]) // low order
	out.Flags = binary.BigEndian.Uint32(data[8:])
	data = data[12:length]
	switch kind {
	case MorxRearrangement:
		out.Data, err = parseRearrangementSubtable(data, numGlyphs)
	case MorxContextual:
		out.Data, err = parseContextualSubtable(data, numGlyphs)
	case MorxLigature:
		out.Data, err = parseLigatureSubtable(data, numGlyphs)
	case MorxNonContextual:
		out.Data, err = parseNonContextualSubtable(data, numGlyphs)
	case MorxInsertion:
		out.Data, err = parseInsertionSubtable(data, numGlyphs)
	default:
		return out, 0, fmt.Errorf("invalid morx subtable type: %d", kind)
	}
	return out, length, err
}

// MorxRearrangementSubtable is a 'morx' subtable format 0.
type MorxRearrangementSubtable AATStateTable

func (MorxRearrangementSubtable) Type() MorxSubtableType { return MorxRearrangement }

// MorxRearrangement flags
const (
	/* If set, make the current glyph the first
	* glyph to be rearranged. */
	MRMarkFirst = 0x8000
	/* If set, don't advance to the next glyph
	* before going to the new state. This means
	* that the glyph index doesn't change, even
	* if the glyph at that index has changed. */
	MRDontAdvance = 0x4000
	/* If set, make the current glyph the last
	* glyph to be rearranged. */
	MRMarkLast = 0x2000
	/* These bits are reserved and should be set to 0. */
	MRReserved = 0x1FF0
	/* The type of rearrangement specified. */
	MRVerb = 0x000F
)

func parseRearrangementSubtable(data []byte, numGlyphs int) (MorxRearrangementSubtable, error) {
	s, err := parseStateTable(data, 0, true, numGlyphs)
	return MorxRearrangementSubtable(s), err
}

type MorxContextualSubtable struct {
	Substitutions []Class
	Machine       AATStateTable
}

func (MorxContextualSubtable) Type() MorxSubtableType { return MorxContextual }

// MorxContextualSubtable flags
const (
	MCSetMark = 0x8000 /* If set, make the current glyph the marked glyph. */
	/* If set, don't advance to the next glyph before
	* going to the new state. */
	MCDontAdvance = 0x4000
	MCReserved    = 0x3FFF /* These bits are reserved and should be set to 0. */
)

func parseContextualSubtable(data []byte, numGlyphs int) (out MorxContextualSubtable, err error) {
	// we need the offset to the data following the stateTable
	if len(data) < aatExtStateHeaderSize+4 {
		return out, errors.New("invalid morx contextual subtable (EOF)")
	}
	subsOffset := binary.BigEndian.Uint32(data[aatExtStateHeaderSize:])
	if len(data) < int(subsOffset) {
		return out, errors.New("invalid morx contextual subtable (EOF)")
	}
	out.Machine, err = parseStateTable(data[:subsOffset], 4, true, numGlyphs)
	if err != nil {
		return out, err
	}

	// find the maximum index need in the substitution array
	var maxi uint16
	for _, entry := range out.Machine.entries {
		markIndex, currentIndex := entry.AsMorxContextual()
		if markIndex != 0xFFFF && markIndex > maxi {
			maxi = markIndex
		}
		if currentIndex != 0xFFFF && currentIndex > maxi {
			maxi = currentIndex
		}
	}

	// know look for these lookup tables
	data = data[subsOffset:]
	if len(data) < 4*(int(maxi)+1) {
		return out, errors.New("invalid morx contextual subtable (EOF)")
	}
	out.Substitutions = make([]Class, maxi+1)
	for i := range out.Substitutions {
		offset := binary.BigEndian.Uint32(data[i*4:])
		out.Substitutions[i], err = parseAATLookupTable(data, offset, numGlyphs, false)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

type MorxLigatureSubtable struct {
	LigatureAction []uint32
	Component      []uint16
	Ligatures      []GID
	Machine        AATStateTable
}

func (MorxLigatureSubtable) Type() MorxSubtableType { return MorxLigature }

// MorxLigatureSubtable flags
const (
	/* Push this glyph onto the component stack for
	* eventual processing. */
	MLSetComponent = 0x8000
	/* Leave the glyph pointer at this glyph for the
	next iteration. */
	MLDontAdvance   = 0x4000
	MLPerformAction = 0x2000 // Use the ligActionIndex to process a ligature group.
	/* Byte offset from beginning of subtable to the
	 * ligature action list. This value must be a
	 * multiple of 4. */
	MLOffset = 0x3FFF

	/* This is the last action in the list. This also
	* implies storage. */
	MLActionLast = 1 << 31
	/* Store the ligature at the current cumulated index
	* in the ligature table in place of the marked
	* (i.e. currently-popped) glyph. */
	MLActionStore = 1 << 30
	/* A 30-bit value which is sign-extended to 32-bits
	* and added to the glyph ID, resulting in an index
	* into the component table. */
	MLActionOffset = 0x3FFFFFFF
)

func parseLigatureSubtable(data []byte, numGlyphs int) (out MorxLigatureSubtable, err error) {
	if len(data) < aatExtStateHeaderSize+12 {
		return out, errors.New("invalid morx ligature subtable (EOF)")
	}
	ligActionOffset := int(binary.BigEndian.Uint32(data[aatExtStateHeaderSize:]))
	componentOffset := int(binary.BigEndian.Uint32(data[aatExtStateHeaderSize+4:]))
	ligatureOffset := int(binary.BigEndian.Uint32(data[aatExtStateHeaderSize+8:]))
	// we need the offset to the data following the stateTable
	// for now, we assume the offsets are actually sorted
	if ligActionOffset > componentOffset || componentOffset > ligatureOffset || len(data) < int(ligatureOffset) {
		return out, errors.New("invalid morx ligature subtable (EOF)")
	}
	out.Machine, err = parseStateTable(data[:ligActionOffset], 2, true, numGlyphs)
	if err != nil {
		return out, err
	}

	// fetch the maximum start index
	maxIndex := -1
	for _, entry := range out.Machine.entries {
		if entry.Flags&MLPerformAction == 0 {
			continue
		}
		if index := int(entry.AsMorxLigature()); index > maxIndex {
			maxIndex = index
		}
	}
	// fetch the action table, up to the last entry
	if len(data) < ligActionOffset+4*int(maxIndex+1) {
		return out, errors.New("invalid morx ligature subtable (EOF)")
	}
	actionData := data[ligActionOffset:]
	for len(actionData) >= 4 { // stop gracefully if the last action was not found
		action := binary.BigEndian.Uint32(actionData)
		// data is truncated to the end of the table,
		// so the memory allocation is bounded by the table size
		out.LigatureAction = append(out.LigatureAction, action)
		actionData = actionData[4:]
		// dont break before maxIndex
		if len(out.LigatureAction) > maxIndex && action&MLActionLast != 0 {
			break
		}
	}

	componentCount := (ligatureOffset - componentOffset) / 2
	out.Component = make([]uint16, componentCount)
	for i := range out.Component {
		out.Component[i] = binary.BigEndian.Uint16(data[componentOffset+2*i:])
	}
	ligatureCount := (len(data) - ligatureOffset) / 2
	out.Ligatures = make([]GID, ligatureCount)
	for i := range out.Ligatures {
		out.Ligatures[i] = GID(binary.BigEndian.Uint16(data[ligatureOffset+2*i:]))
	}
	return out, nil
}

type MorxNonContextualSubtable struct {
	Class // the lookup value is interpreted as a GlyphIndex
}

func (MorxNonContextualSubtable) Type() MorxSubtableType { return MorxNonContextual }

func parseNonContextualSubtable(data []byte, numGlyphs int) (MorxNonContextualSubtable, error) {
	c, err := parseAATLookupTable(data, 0, numGlyphs, false)
	return MorxNonContextualSubtable{c}, err
}

type MorxInsertionSubtable struct {
	// After successul parsing, this array may be safely
	// indexed by the indexes and counts from Machine entries.
	Insertions []GID
	Machine    AATStateTable
}

func (MorxInsertionSubtable) Type() MorxSubtableType { return MorxInsertion }

// MorxInsertionSubtable flags
const (
	// If set, mark the current glyph.
	MISetMark = 0x8000
	// If set, don't advance to the next glyph before
	// going to the new state.  This does not mean
	// that the glyph pointed to is the same one as
	// before. If you've made insertions immediately
	// downstream of the current glyph, the next glyph
	// processed would in fact be the first one
	// inserted.
	MIDontAdvance = 0x4000
	// If set, and the currentInsertList is nonzero,
	// then the specified glyph list will be inserted
	// as a kashida-like insertion, either before or
	// after the current glyph (depending on the state
	// of the currentInsertBefore flag). If clear, and
	// the currentInsertList is nonzero, then the
	// specified glyph list will be inserted as a
	// split-vowel-like insertion, either before or
	// after the current glyph (depending on the state
	// of the currentInsertBefore flag).
	MICurrentIsKashidaLike = 0x2000
	// If set, and the markedInsertList is nonzero,
	// then the specified glyph list will be inserted
	// as a kashida-like insertion, either before or
	// after the marked glyph (depending on the state
	// of the markedInsertBefore flag). If clear, and
	// the markedInsertList is nonzero, then the
	// specified glyph list will be inserted as a
	// split-vowel-like insertion, either before or
	// after the marked glyph (depending on the state
	// of the markedInsertBefore flag).
	MIMarkedIsKashidaLike = 0x1000
	// If set, specifies that insertions are to be made
	// to the left of the current glyph. If clear,
	// they're made to the right of the current glyph.
	MICurrentInsertBefore = 0x0800
	// If set, specifies that insertions are to be
	// made to the left of the marked glyph. If clear,
	// they're made to the right of the marked glyph.
	MIMarkedInsertBefore = 0x0400
	// This 5-bit field is treated as a count of the
	// number of glyphs to insert at the current
	// position. Since zero means no insertions, the
	// largest number of insertions at any given
	// current location is 31 glyphs.
	MICurrentInsertCount = 0x3E0
	// This 5-bit field is treated as a count of the
	// number of glyphs to insert at the marked
	// position. Since zero means no insertions, the
	// largest number of insertions at any given
	// marked location is 31 glyphs.
	MIMarkedInsertCount = 0x001F
)

func parseInsertionSubtable(data []byte, numGlyphs int) (out MorxInsertionSubtable, err error) {
	// we need the offset to the data following the stateTable
	if len(data) < aatExtStateHeaderSize+4 {
		return out, errors.New("invalid morx insertion subtable (EOF)")
	}
	insertionOffset := binary.BigEndian.Uint32(data[aatExtStateHeaderSize:])
	if len(data) < int(insertionOffset) {
		return out, errors.New("invalid morx insertion subtable (EOF)")
	}
	out.Machine, err = parseStateTable(data[:insertionOffset], 4, true, numGlyphs)
	if err != nil {
		return out, err
	}

	// find the maximum index needed in the insertions array,
	// taking into account the number of insertions
	var maxi uint16
	for _, entry := range out.Machine.entries {
		currentIndex, markedIndex := entry.AsMorxInsertion()
		if currentIndex != 0xFFFF {
			indexEnd := currentIndex + (entry.Flags&MICurrentInsertCount)>>5
			if indexEnd > maxi {
				maxi = indexEnd
			}
		}
		if markedIndex != 0xFFFF {
			indexEnd := markedIndex + entry.Flags&MIMarkedInsertCount
			if indexEnd > maxi {
				maxi = indexEnd
			}
		}
	}

	// know look for these lookup tables
	data = data[insertionOffset:]
	if len(data) < 2*int(maxi) {
		return out, errors.New("invalid morx insertion subtable (EOF)")
	}
	out.Insertions = make([]GID, maxi)
	for i := range out.Insertions {
		out.Insertions[i] = GID(binary.BigEndian.Uint16(data[i*2:]))
	}

	return out, nil
}
