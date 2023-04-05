// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
)

// The following are types shared by GSUB and GPOS tables

// Coverage specifies all the glyphs affected by a substitution or
// positioning operation described in a subtable.
// Conceptually is it a []GlyphIndex, with an Index method,
// but it may be implemented for efficiently.
// See https://learn.microsoft.com/typography/opentype/spec/chapter2#lookup-table
type Coverage interface {
	isCov()

	// Index returns the index of the provided glyph, or
	// `false` if the glyph is not covered by this lookup.
	// Note: this method is injective: two distincts, covered glyphs are mapped
	// to distincts indices.
	Index(GlyphID) (int, bool)

	// Len return the number of glyphs covered.
	// It is 0 for empty coverages.
	// For non empty Coverages, it is also 1 + (maximum index returned)
	Len() int
}

func (Coverage1) isCov() {}
func (Coverage2) isCov() {}

type Coverage1 struct {
	format uint16    `unionTag:"1"`
	Glyphs []GlyphID `arrayCount:"FirstUint16"`
}

type Coverage2 struct {
	format uint16        `unionTag:"2"`
	Ranges []RangeRecord `arrayCount:"FirstUint16"`
}

type RangeRecord struct {
	StartGlyphID       GlyphID // First glyph ID in the range
	EndGlyphID         GlyphID // Last glyph ID in the range
	StartCoverageIndex uint16  // Coverage Index of first glyph ID in range
}

// ClassDef stores a value for a set of GlyphIDs.
// Conceptually it is a map[GlyphID]uint16, but it may
// be implemented more efficiently.
type ClassDef interface {
	isClassDef()
	Class(gi GlyphID) (uint16, bool)

	// Extent returns the maximum class ID + 1. This is the length
	// required for an array to be indexed by the class values.
	Extent() int
}

func (ClassDef1) isClassDef() {}
func (ClassDef2) isClassDef() {}

type ClassDef1 struct {
	format          uint16   `unionTag:"1"`
	StartGlyphID    GlyphID  // First glyph ID of the classValueArray
	ClassValueArray []uint16 `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type ClassDef2 struct {
	format            uint16             `unionTag:"2"`
	ClassRangeRecords []ClassRangeRecord `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type ClassRangeRecord struct {
	StartGlyphID GlyphID // First glyph ID in the range
	EndGlyphID   GlyphID // Last glyph ID in the range
	Class        uint16  // Applied to all glyphs in the range
}

// Lookups

type SequenceLookupRecord struct {
	SequenceIndex   uint16 // Index (zero-based) into the input glyph sequence
	LookupListIndex uint16 // Index (zero-based) into the LookupList
}

type SequenceContextFormat1 struct {
	format     uint16            `unionTag:"1"`
	coverage   Coverage          `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of SequenceContextFormat1 table
	SeqRuleSet []SequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[seqRuleSetCount]	Array of offsets to SequenceRuleSet tables, from beginning of SequenceContextFormat1 table (offsets may be NULL)
}

func (sc *SequenceContextFormat1) sanitize(lookupCount uint16) error {
	for _, set := range sc.SeqRuleSet {
		for _, rule := range set.SeqRule {
			if err := rule.sanitize(lookupCount); err != nil {
				return err
			}
		}
	}
	return nil
}

type SequenceRuleSet struct {
	SeqRule []SequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type SequenceRule struct {
	glyphCount       uint16                 // Number of glyphs in the input glyph sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	InputSequence    []GlyphID              `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Array of input glyph IDs—starting with the second glyph
	SeqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of Sequence lookup records
}

func (sr *SequenceRule) sanitize(lookupCount uint16) error {
	for _, rec := range sr.SeqLookupRecords {
		if rec.SequenceIndex >= sr.glyphCount {
			return fmt.Errorf("invalid sequence lookup table (input index %d >= %d)", rec.SequenceIndex, sr.glyphCount)
		}
		if rec.LookupListIndex >= lookupCount {
			return fmt.Errorf("invalid sequence lookup table (lookup index %d >= %d)", rec.LookupListIndex, lookupCount)
		}
	}
	return nil
}

type SequenceContextFormat2 struct {
	format          uint16                 `unionTag:"2"`
	coverage        Coverage               `offsetSize:"Offset16"`                            //	Offset to Coverage table, from beginning of SequenceContextFormat2 table
	ClassDef        ClassDef               `offsetSize:"Offset16"`                            //	Offset to ClassDef table, from beginning of SequenceContextFormat2 table
	ClassSeqRuleSet []ClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[classSeqRuleSetCount]	Array of offsets to ClassSequenceRuleSet tables, from beginning of SequenceContextFormat2 table (may be NULL)
}

// ClassSequenceRuleSet has the same binary format as SequenceRuleSet,
// and using the same type simplifies later processing.
type ClassSequenceRuleSet = SequenceRuleSet

type SequenceContextFormat3 struct {
	format           uint16                 `unionTag:"3"`
	glyphCount       uint16                 // Number of glyphs in the input sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	Coverages        []Coverage             `arrayCount:"ComputedField-glyphCount" offsetsArray:"Offset16"` //[glyphCount]	Array of offsets to Coverage tables, from beginning of SequenceContextFormat3 subtable
	SeqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

type ChainedSequenceContextFormat1 struct {
	format            uint16                   `unionTag:"1"`
	coverage          Coverage                 `offsetSize:"Offset16"`                             //	Offset to Coverage table, from beginning of ChainSequenceContextFormat1 table
	ChainedSeqRuleSet []ChainedSequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[chainedSeqRuleSetCount]	Array of offsets to ChainedSeqRuleSet tables, from beginning of ChainedSequenceContextFormat1 table (may be NULL)
}

type ChainedSequenceRuleSet struct {
	ChainedSeqRules []ChainedSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type ChainedSequenceRule struct {
	BacktrackSequence []GlyphID              `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of backtrack glyph IDs
	inputGlyphCount   uint16                 //	Number of glyphs in the input sequence
	InputSequence     []GlyphID              `arrayCount:"ComputedField-inputGlyphCount-1"` //[inputGlyphCount - 1]	Array of input glyph IDs—start with second glyph
	LookaheadSequence []GlyphID              `arrayCount:"FirstUint16"`                     //[lookaheadGlyphCount]	Array of lookahead glyph IDs
	SeqLookupRecords  []SequenceLookupRecord `arrayCount:"FirstUint16"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

type ChainedSequenceContextFormat2 struct {
	format                 uint16                        `unionTag:"2"`
	coverage               Coverage                      `offsetSize:"Offset16"`                            // Offset to Coverage table, from beginning of ChainedSequenceContextFormat2 table
	BacktrackClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing backtrack sequence context, from beginning of ChainedSequenceContextFormat2 table
	InputClassDef          ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing input sequence context, from beginning of ChainedSequenceContextFormat2 table
	LookaheadClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing lookahead sequence context, from beginning of ChainedSequenceContextFormat2 table
	ChainedClassSeqRuleSet []ChainedClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[chainedClassSeqRuleSetCount]	Array of offsets to ChainedClassSequenceRuleSet tables, from beginning of ChainedSequenceContextFormat2 table (may be NULL)
}

// ChainedClassSequenceRuleSet has the same binary format as ChainedSequenceRuleSet,
// and using the same type simplifies later processing.
type ChainedClassSequenceRuleSet = ChainedSequenceRuleSet

type ChainedSequenceContextFormat3 struct {
	format             uint16                 `unionTag:"3"`
	BacktrackCoverages []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[backtrackGlyphCount]	Array of offsets to coverage tables for the backtrack sequence
	InputCoverages     []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[inputGlyphCount]	Array of offsets to coverage tables for the input sequence
	LookaheadCoverages []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables for the lookahead sequence
	SeqLookupRecords   []SequenceLookupRecord `arrayCount:"FirstUint16"`                         //[seqLookupCount]	Array of SequenceLookupRecords
}

type Extension struct {
	substFormat         uint16   //	Format identifier. Set to 1.
	ExtensionLookupType uint16   //	Lookup type of subtable referenced by extensionOffset (that is, the extension subtable).
	ExtensionOffset     Offset32 //	Offset to the extension subtable, of lookup type extensionLookupType, relative to the start of the ExtensionSubstFormat1 subtable.
	RawData             []byte   `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

// GSUB is the Glyph Substitution (GSUB) table.
// It provides data for substition of glyphs for appropriate rendering of scripts,
// such as cursively-connecting forms in Arabic script,
// or for advanced typographic effects, such as ligatures.
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/gsub
type GSUB Layout

// GSUBLookup is one lookup subtable data
type GSUBLookup interface {
	isGSUBLookup()

	// Coverage returns the coverage of the lookup subtable.
	// For ContextualSubs3 and ChainedContextualSubs3, its the coverage of the first input.
	Cov() Coverage
}

func (SingleSubs) isGSUBLookup()             {}
func (MultipleSubs) isGSUBLookup()           {}
func (AlternateSubs) isGSUBLookup()          {}
func (LigatureSubs) isGSUBLookup()           {}
func (ContextualSubs) isGSUBLookup()         {}
func (ChainedContextualSubs) isGSUBLookup()  {}
func (ExtensionSubs) isGSUBLookup()          {}
func (ReverseChainSingleSubs) isGSUBLookup() {}

func (ms MultipleSubs) Sanitize() error {
	if exp, got := ms.Coverage.Len(), len(ms.Sequences); exp != got {
		return fmt.Errorf("GSUB: invalid MultipleSubs sequences count (%d != %d)", exp, got)
	}
	return nil
}

func (ls LigatureSubs) Sanitize() error {
	if exp, got := ls.Coverage.Len(), len(ls.LigatureSets); exp != got {
		return fmt.Errorf("GSUB: invalid LigatureSubs sets count (%d != %d)", exp, got)
	}
	return nil
}

func (cs ContextualSubs) Sanitize(lookupCount uint16) error {
	if f1, isFormat1 := cs.Data.(ContextualSubs1); isFormat1 {
		return (*SequenceContextFormat1)(&f1).sanitize(lookupCount)
	}
	return nil
}

func (rs ReverseChainSingleSubs) Sanitize() error {
	if exp, got := rs.coverage.Len(), len(rs.SubstituteGlyphIDs); exp != got {
		return fmt.Errorf("GSUB: invalid ReverseChainSingleSubs glyphs count (%d != %d)", exp, got)
	}
	return nil
}

func (ext ExtensionSubs) Resolve() (GSUBLookup, error) {
	if L, E := len(ext.RawData), int(ext.ExtensionOffset); L < E {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}
	lk, err := parseGSUBLookup(ext.RawData[ext.ExtensionOffset:], ext.ExtensionLookupType)
	if err != nil {
		return nil, err
	}
	if _, isExt := lk.(ExtensionSubs); isExt {
		return nil, errors.New("invalid extension substitution table")
	}
	return lk, nil
}

func parseGSUBLookup(src []byte, lookupType uint16) (out GSUBLookup, err error) {
	switch lookupType {
	case 1: // Single (format 1.1 1.2)	Replace one glyph with one glyph
		out, _, err = ParseSingleSubs(src)
	case 2: // Multiple (format 2.1)	Replace one glyph with more than one glyph
		out, _, err = ParseMultipleSubs(src)
	case 3: // Alternate (format 3.1)	Replace one glyph with one of many glyphs
		out, _, err = ParseAlternateSubs(src)
	case 4: // Ligature (format 4.1)	Replace multiple glyphs with one glyph
		out, _, err = ParseLigatureSubs(src)
	case 5: // Context (format 5.1 5.2 5.3)	Replace one or more glyphs in context
		out, _, err = ParseContextualSubs(src)
	case 6: // Chaining Context (format 6.1 6.2 6.3)	Replace one or more glyphs in chained context
		out, _, err = ParseChainedContextualSubs(src)
	case 7: // Extension Substitution (format 7.1) Extension mechanism for other substitutions
		out, _, err = ParseExtensionSubs(src)
	case 8: // Reverse chaining context single (format 8.1)
		out, _, err = ParseReverseChainSingleSubs(src)
	default:
		err = fmt.Errorf("invalid GSUB Loopkup type %d", lookupType)
	}
	return out, err
}

// AsGSUBLookups returns the GSUB lookup subtables.
func (lk Lookup) AsGSUBLookups() ([]GSUBLookup, error) {
	var err error
	out := make([]GSUBLookup, len(lk.subtableOffsets))
	for i, offset := range lk.subtableOffsets {
		if L := len(lk.rawData); L < int(offset) {
			return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		out[i], err = parseGSUBLookup(lk.rawData[offset:], lk.lookupType)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ------------------------ GPOS common data structures ------------------------

// GPOS is the Glyph Positioning (GPOS) table.
// It provides precise control over glyph placement
// for sophisticated text layout and rendering in each script
// and language system that a font supports.
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/gpos
type GPOS Layout

type GPOSLookup interface {
	isGPOSLookup()

	// Coverage returns the coverage of the lookup subtable.
	// For ContextualPos3 and ChainedContextualPos3, its the coverage of the first input.
	Cov() Coverage
}

func (SinglePos) isGPOSLookup()            {}
func (PairPos) isGPOSLookup()              {}
func (CursivePos) isGPOSLookup()           {}
func (MarkBasePos) isGPOSLookup()          {}
func (MarkLigPos) isGPOSLookup()           {}
func (MarkMarkPos) isGPOSLookup()          {}
func (ContextualPos) isGPOSLookup()        {}
func (ChainedContextualPos) isGPOSLookup() {}
func (ExtensionPos) isGPOSLookup()         {}

func (sp *SinglePos) Sanitize() error {
	if f2, isFormat2 := sp.Data.(SinglePosData2); isFormat2 {
		if exp, got := f2.coverage.Len(), len(f2.ValueRecords); exp != got {
			return fmt.Errorf("GPOS: invalid SinglePos values count (%d != %d)", exp, got)
		}
	}
	return nil
}

func (pp *PairPos) Sanitize() error {
	if f1, isFormat1 := pp.Data.(PairPosData1); isFormat1 {
		// there are fonts with to much PairSets : accept it
		if exp, got := f1.coverage.Len(), len(f1.PairSets); exp > got {
			return fmt.Errorf("GPOS: invalid PairPos1 sets count (%d > %d)", exp, got)
		}
	} else if f2, isFormat2 := pp.Data.(PairPosData2); isFormat2 {
		if exp, got := f2.ClassDef1.Extent(), int(f2.class1Count); exp != got {
			return fmt.Errorf("GPOS: invalid PairPos2 class1 count (%d != %d)", exp, got)
		}
		if exp, got := f2.ClassDef2.Extent(), int(f2.class2Count); exp != got {
			return fmt.Errorf("GPOS: invalid PairPos2 class2 count (%d != %d)", exp, got)
		}
	}
	return nil
}

func (mp *MarkBasePos) Sanitize() error {
	if exp, got := mp.markCoverage.Len(), len(mp.MarkArray.MarkRecords); exp != got {
		return fmt.Errorf("GPOS: invalid MarkBasePos marks count (%d != %d)", exp, got)
	}
	if exp, got := mp.BaseCoverage.Len(), len(mp.BaseArray.BaseAnchors); exp != got {
		return fmt.Errorf("GPOS: invalid MarkBasePos marks count (%d != %d)", exp, got)
	}

	return nil
}

func (mp *MarkLigPos) Sanitize() error {
	if exp, got := mp.MarkCoverage.Len(), len(mp.MarkArray.MarkAnchors); exp != got {
		return fmt.Errorf("GPOS: invalid MarkBasePos marks count (%d != %d)", exp, got)
	}
	if exp, got := mp.LigatureCoverage.Len(), len(mp.LigatureArray.LigatureAttachs); exp != got {
		return fmt.Errorf("GPOS: invalid MarkBasePos marks count (%d != %d)", exp, got)
	}

	return nil
}

func (cs *ContextualPos) Sanitize(lookupCount uint16) error {
	if f1, isFormat1 := cs.Data.(ContextualPos1); isFormat1 {
		return (*SequenceContextFormat1)(&f1).sanitize(lookupCount)
	}
	return nil
}

func (ext ExtensionPos) Resolve() (GPOSLookup, error) {
	if L, E := len(ext.RawData), int(ext.ExtensionOffset); L < E {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}
	lk, err := parseGPOSLookup(ext.RawData[ext.ExtensionOffset:], ext.ExtensionLookupType)
	if err != nil {
		return nil, err
	}
	if _, isExt := lk.(ExtensionPos); isExt {
		return nil, errors.New("invalid extension positioning table")
	}
	return lk, nil
}

func parseGPOSLookup(src []byte, lookupType uint16) (out GPOSLookup, err error) {
	switch lookupType {
	case 1: // Single adjustment	Adjust position of a single glyph
		out, _, err = ParseSinglePos(src)
	case 2: // Pair adjustment	Adjust position of a pair of glyphs
		out, _, err = ParsePairPos(src)
	case 3: // Cursive attachment	Attach cursive glyphs
		out, _, err = ParseCursivePos(src)
	case 4: // MarkToBase attachment	Attach a combining mark to a base glyph
		out, _, err = ParseMarkBasePos(src)
	case 5: // MarkToLigature attachment	Attach a combining mark to a ligature
		out, _, err = ParseMarkLigPos(src)
	case 6: // MarkToMark attachment	Attach a combining mark to another mark
		out, _, err = ParseMarkMarkPos(src)
	case 7: // Context positioning	Position one or more glyphs in context
		out, _, err = ParseContextualPos(src)
	case 8: // Chained Context positioning	Position one or more glyphs in chained context
		out, _, err = ParseChainedContextualPos(src)
	case 9: // Extension positioning	Extension mechanism for other positionings
		out, _, err = ParseExtensionPos(src)
	default:
		err = fmt.Errorf("invalid GPOS Loopkup type %d", lookupType)
	}
	return out, err
}

// AsGPOSLookups returns the GPOS lookup subtables
func (lk Lookup) AsGPOSLookups() ([]GPOSLookup, error) {
	var err error
	out := make([]GPOSLookup, len(lk.subtableOffsets))
	for i, offset := range lk.subtableOffsets {
		if L := len(lk.rawData); L < int(offset) {
			return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		out[i], err = parseGPOSLookup(lk.rawData[offset:], lk.lookupType)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ValueFormat is a mask indicating which field
// are set in a GPOS [ValueRecord].
// It is often shared between many records.
type ValueFormat uint16

// number of fields present
func (f ValueFormat) size() int { return bits.OnesCount16(uint16(f)) }

const (
	XPlacement ValueFormat = 1 << iota // Includes horizontal adjustment for placement
	YPlacement                         // Includes vertical adjustment for placement
	XAdvance                           // Includes horizontal adjustment for advance
	YAdvance                           // Includes vertical adjustment for advance
	XPlaDevice                         // Includes horizontal Device table for placement
	YPlaDevice                         // Includes vertical Device table for placement
	XAdvDevice                         // Includes horizontal Device table for advance
	YAdvDevice                         // Includes vertical Device table for advance

	//  Mask for having any Device table
	Devices = XPlaDevice | YPlaDevice | XAdvDevice | YAdvDevice
)

// ValueRecord has optional fields
type ValueRecord struct {
	XPlacement int16       // Horizontal adjustment for placement, in design units.
	YPlacement int16       // Vertical adjustment for placement, in design units.
	XAdvance   int16       // Horizontal adjustment for advance, in design units — only used for horizontal layout.
	YAdvance   int16       // Vertical adjustment for advance, in design units — only used for vertical layout.
	XPlaDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for horizontal placement, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	YPlaDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for vertical placement, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	XAdvDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for horizontal advance, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup subtable) — may be NULL.
	YAdvDevice DeviceTable // Offset to Device table (non-variable font) / VariationIndex table (variable font) for vertical advance, from beginning of the immediate parent table (SinglePos or PairPosFormat2 lookup subtable, PairSet table within a PairPosFormat1 lookup su
}

// [data] must start at the immediate parent table, [offset] indicating
// the start of the record in it.
// Returns [offset] + the number of bytes read from [offset]
// Note that a [format] with value 0, is supported, resulting in a no-op
func parseValueRecord(format ValueFormat, data []byte, offset int) (out ValueRecord, _ int, err error) {
	if L := len(data); L < offset {
		return out, 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
	}

	size := format.size() // number of fields present
	if size == 0 {        // return early
		return out, offset, nil
	}
	// start by parsing the list of values
	values, err := ParseUint16s(data[offset:], size)
	if err != nil {
		return out, 0, fmt.Errorf("invalid value record: %s", err)
	}
	// follow the order
	cursor := 0
	if format&XPlacement != 0 {
		out.XPlacement = int16(values[cursor])
		cursor++
	}
	if format&YPlacement != 0 {
		out.YPlacement = int16(values[cursor])
		cursor++
	}
	if format&XAdvance != 0 {
		out.XAdvance = int16(values[cursor])
		cursor++
	}
	if format&YAdvance != 0 {
		out.YAdvance = int16(values[cursor])
		cursor++
	}
	if format&XPlaDevice != 0 {
		if devOffset := values[cursor]; devOffset != 0 {
			out.XPlaDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		cursor++
	}
	if format&YPlaDevice != 0 {
		if devOffset := values[cursor]; devOffset != 0 {
			out.YPlaDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		cursor++
	}
	if format&XAdvDevice != 0 {
		if devOffset := values[cursor]; devOffset != 0 {
			out.XAdvDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		cursor++
	}
	if format&YAdvDevice != 0 {
		if devOffset := values[cursor]; devOffset != 0 {
			out.YAdvDevice, err = parseDeviceTable(data, devOffset)
			if err != nil {
				return out, 0, err
			}
		}
		cursor++ // useless actually
	}
	return out, offset + 2*size, err
}

// DeviceTable is either an DeviceHinting for standard fonts,
// or a DeviceVariation for variable fonts.
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

type DeviceVariation VariationStoreIndex

func parseDeviceTable(src []byte, offset uint16) (DeviceTable, error) {
	if L := len(src); L < int(offset)+6 {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", offset+6, L)
	}
	var header DeviceTableHeader
	header.mustParse(src[offset:])

	switch format := header.deltaFormat; format {
	case 1, 2, 3:
		var out DeviceHinting

		out.StartSize, out.EndSize = header.first, header.second
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
		uint16s, err := ParseUint16s(src[offset+6:], count)
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
		return DeviceVariation{DeltaSetOuter: header.first, DeltaSetInner: header.second}, nil
	default:
		return nil, fmt.Errorf("unsupported positionning device subtable: %d", format)
	}
}

type PairValueRecord struct {
	SecondGlyph  GlyphID     //	Glyph ID of second glyph in the pair (first glyph is listed in the Coverage table).
	ValueRecord1 ValueRecord //	Positioning data for the first glyph in the pair.
	ValueRecord2 ValueRecord //	Positioning data for the second glyph in the pair.
}

type Class1Record []Class2Record //[class2Count]	Array of Class2 records, ordered by classes in classDef2.

type Class2Record struct {
	ValueRecord1 ValueRecord //	Positioning for first glyph — empty if valueFormat1 = 0.
	ValueRecord2 ValueRecord //	Positioning for second glyph — empty if valueFormat2 = 0.
}

func (AnchorFormat1) isAnchor() {}
func (AnchorFormat2) isAnchor() {}
func (AnchorFormat3) isAnchor() {}

type AnchorFormat1 struct {
	anchorFormat uint16 `unionTag:"1"`
	XCoordinate  int16  // Horizontal value, in design units
	YCoordinate  int16  // Vertical value, in design units
}

type AnchorFormat2 struct {
	anchorFormat uint16 `unionTag:"2"`
	XCoordinate  int16  // Horizontal value, in design units
	YCoordinate  int16  // Vertical value, in design units
	AnchorPoint  uint16 // Index to glyph contour point
}

type AnchorFormat3 struct {
	anchorFormat  uint16      `unionTag:"3"`
	XCoordinate   int16       // Horizontal value, in design units
	YCoordinate   int16       // Vertical value, in design units
	xDeviceOffset Offset16    // Offset to Device table (non-variable font) / VariationIndex table (variable font) for X coordinate, from beginning of Anchor table (may be NULL)
	yDeviceOffset Offset16    // Offset to Device table (non-variable font) / VariationIndex table (variable font) for Y coordinate, from beginning of Anchor table (may be NULL)
	XDevice       DeviceTable `isOpaque:""` // Offset to Device table (non-variable font) / VariationIndex table (variable font) for X coordinate, from beginning of Anchor table (may be NULL)
	YDevice       DeviceTable `isOpaque:""` // Offset to Device table (non-variable font) / VariationIndex table (variable font) for Y coordinate, from beginning of Anchor table (may be NULL)
}

func (af *AnchorFormat3) parseXDevice(src []byte) error {
	if af.xDeviceOffset == 0 {
		return nil
	}
	var err error
	af.XDevice, err = parseDeviceTable(src, uint16(af.xDeviceOffset))
	return err
}

func (af *AnchorFormat3) parseYDevice(src []byte) error {
	if af.yDeviceOffset == 0 {
		return nil
	}
	var err error
	af.YDevice, err = parseDeviceTable(src, uint16(af.yDeviceOffset))
	return err
}

type MarkArray struct {
	MarkRecords []MarkRecord `arrayCount:"FirstUint16"` //[markCount]	Array of MarkRecords, ordered by corresponding glyphs in the associated mark Coverage table.
	MarkAnchors []Anchor     `isOpaque:""`              // with same length as MarkRecords
}

func (ma *MarkArray) parseMarkAnchors(src []byte) error {
	ma.MarkAnchors = make([]Anchor, len(ma.MarkRecords))
	var err error
	for i, rec := range ma.MarkRecords {
		if L := len(src); L < int(rec.markAnchorOffset) {
			return fmt.Errorf("EOF: expected length: %d, got %d", rec.markAnchorOffset, L)
		}
		ma.MarkAnchors[i], _, err = ParseAnchor(src[rec.markAnchorOffset:])
		if err != nil {
			return err
		}
	}
	return nil
}

type MarkRecord struct {
	MarkClass        uint16   // Class defined for the associated mark.
	markAnchorOffset Offset16 // Offset to Anchor table, from beginning of MarkArray table.
}

// ------------------------------ parsing helpers ------------------------------

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

// ParseUint16s interprets data as a (big endian) uint16 slice.
// It returns an error if [data] is not long enough for the given [count].
func ParseUint16s(src []byte, count int) ([]uint16, error) {
	if L := len(src); L < 2*count {
		return nil, fmt.Errorf("EOF: expected length: %d, got %d", 2*count, L)
	}
	out := make([]uint16, count)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(src[2*i:])
	}
	return out, nil
}
