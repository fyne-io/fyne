package truetype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type LookupFlag = uint16

const (
	// This bit relates only to the correct processing of
	// the cursive attachment lookup type (GPOS lookup type 3).
	// When this bit is set, the last glyph in a given sequence to
	// which the cursive attachment lookup is applied, will be positioned on the baseline.
	RightToLeft      LookupFlag = 1 << iota
	IgnoreBaseGlyphs            // If set, skips over base glyphs
	IgnoreLigatures             // If set, skips over ligatures
	IgnoreMarks                 // If set, skips over all combining marks
	// If set, indicates that the lookup table structure
	// is followed by a MarkFilteringSet field.
	// The layout engine skips over all mark glyphs not in the mark filtering set indicated.
	UseMarkFilteringSet
	Reserved LookupFlag = 0x00E0 // For future use (Set to zero)
	// If not zero, skips over all marks of attachment
	// type different from specified.
	MarkAttachmentType LookupFlag = 0xFF00
)

// TableLayout represents the common layout table used by GPOS and GSUB.
// The Features field contains all the features for this layout. However,
// the script and language determines which feature is used.
//
// See https://www.microsoft.com/typography/otspec/chapter2.htm#organization
// See https://www.microsoft.com/typography/otspec/GPOS.htm
// See https://www.microsoft.com/typography/otspec/GSUB.htm
type TableLayout struct {
	Scripts           []Script
	Features          []FeatureRecord
	FeatureVariations []FeatureVariation
	header            layoutHeader11
}

// FindScript looks for `script` and return its index into the Scripts slice,
// or -1 if the tag is not found.
func (t *TableLayout) FindScript(script Tag) int {
	// Scripts is sorted: binary search
	low, high := 0, len(t.Scripts)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := t.Scripts[mid].Tag
		if script < p {
			high = mid
		} else if script > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// FindVariationIndex returns the first feature variation matching
// the specified variation coordinates, as an index in the
// `FeatureVariations` field.
// It returns `-1` if not found.
func (t *TableLayout) FindVariationIndex(coords []float32) int {
	for i, record := range t.FeatureVariations {
		if record.evaluate(coords) {
			return i
		}
	}
	return -1
}

// FindFeatureIndex fetches the index for a given feature tag in the specified face's GSUB table
// or GPOS table.
// Returns false if not found
func (t *TableLayout) FindFeatureIndex(featureTag Tag) (uint16, bool) {
	for i, feat := range t.Features { // i fits in uint16
		if featureTag == feat.Tag {
			return uint16(i), true
		}
	}
	return 0, false
}

// Script represents a single script (i.e "latn" (Latin), "cyrl" (Cyrillic), etc).
type Script struct {
	DefaultLanguage *LangSys
	Languages       []LangSys
	Tag             Tag
}

// FindLanguage looks for `language` and return its index into the Languages slice,
// or -1 if the tag is not found.
func (t Script) FindLanguage(language Tag) int {
	// Languages is sorted: binary search
	low, high := 0, len(t.Languages)
	for low < high {
		mid := low + (high-low)/2 // avoid overflow when computing mid
		p := t.Languages[mid].Tag
		if language < p {
			high = mid
		} else if language > p {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

// GetLangSys return the language at `index`. It `index` is out of range (for example with 0xFFFF),
// it returns `DefaultLanguage` (which may be empty)
func (t Script) GetLangSys(index uint16) LangSys {
	if int(index) >= len(t.Languages) {
		if t.DefaultLanguage != nil {
			return *t.DefaultLanguage
		}
		return LangSys{}
	}
	return t.Languages[index]
}

// FeatureRecord associate a tag with a feature
type FeatureRecord struct {
	Feature
	Tag Tag
}

// Feature represents a glyph substitution or glyph positioning features.
type Feature struct {
	LookupIndices []uint16
	paramsOffet   uint16
}

type LookupOptions struct {
	Flag LookupFlag // Lookup qualifiers.
	// Index (base 0) into GDEF mark glyph sets structure,
	// meaningfull only if UseMarkFilteringSet is set.
	MarkFilteringSet uint16 // TODO: sanitize with gdef
}

// Props returns a 32-bit integer where the lower 16-bit is `Flag` and
// the higher 16-bit is `MarkFilteringSet` if the lookup uses one.
func (l LookupOptions) Props() uint32 {
	flag := uint32(l.Flag)
	if l.Flag&UseMarkFilteringSet != 0 {
		flag |= uint32(l.MarkFilteringSet) << 16
	}
	return flag
}

// lookup represents a feature lookup table, common to GSUB and GPOS, before resolving
// the specialized lookup format.
type lookup struct {
	subtableOffsets []uint16 // Array of offsets to lookup subtables, from beginning of Lookup table
	data            []byte   // input data of the lookup table
	LookupOptions
	kind uint16
}

// versionHeader is the beginning of on-disk format of the GPOS/GSUB version header.
// See https://www.microsoft.com/typography/otspec/GPOS.htm
// See https://www.microsoft.com/typography/otspec/GSUB.htm
type versionHeader struct {
	Major uint16 // Major version of the GPOS/GSUB table.
	Minor uint16 // Minor version of the GPOS/GSUB table.
}

// layoutHeader10 is the on-disk format of GPOS/GSUB version header when major=1 and minor=0.
type layoutHeader10 struct {
	ScriptListOffset  uint16 // offset to ScriptList table, from beginning of GPOS/GSUB table.
	FeatureListOffset uint16 // offset to FeatureList table, from beginning of GPOS/GSUB table.
	LookupListOffset  uint16 // offset to LookupList table, from beginning of GPOS/GSUB table.
}

// layoutHeader11 is the on-disk format of GPOS/GSUB version header when major=1 and minor=1.
type layoutHeader11 struct {
	layoutHeader10
	FeatureVariationsOffset uint32 // offset to FeatureVariations table, from beginning of GPOS/GSUB table (may be NULL).
}

// tagOffsetRecord is a on-disk format of a Tag and Offset record, commonly used thoughout this table.
type tagOffsetRecord struct {
	Tag    Tag    // 4-byte script tag identifier
	Offset uint16 // Offset to object from beginning of list
}

type (
	scriptRecord  = tagOffsetRecord
	featureRecord = tagOffsetRecord
	lookupRecord  = tagOffsetRecord
	langSysRecord = tagOffsetRecord
)

// LangSys represents the language system for a script.
type LangSys struct {
	// Features contains the index of the features for this language,
	// relative to the Features slice of the table
	Features []uint16
	// Index of a feature required for this language system.
	// If no required features, default to 0xFFFF
	RequiredFeatureIndex uint16
	Tag                  Tag
}

// parseLangSys parses a single Language System table. b expected to be the beginning of Script table.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#langSysTbl
func (t *TableLayout) parseLangSys(b []byte, record langSysRecord) (LangSys, error) {
	var out LangSys
	if int(record.Offset) >= len(b) {
		return out, io.ErrUnexpectedEOF
	}

	r := bytes.NewReader(b[record.Offset:])

	var lang struct {
		LookupOrder          uint16 // = NULL (reserved for an offset to a reordering table)
		RequiredFeatureIndex uint16 // Index of a feature required for this language system; if no required features = 0xFFFF
		FeatureIndexCount    uint16 // Number of feature index values for this language system — excludes the required feature
		// featureIndices[featureIndexCount] uint16 // Array of indices into the FeatureList, in arbitrary order
	}

	if err := binary.Read(r, binary.BigEndian, &lang); err != nil {
		return out, fmt.Errorf("reading langSysTable: %s", err)
	}

	featureIndices := make([]uint16, lang.FeatureIndexCount)
	if err := binary.Read(r, binary.BigEndian, &featureIndices); err != nil {
		return out, fmt.Errorf("reading langSysTable featureIndices[%d]: %s", lang.FeatureIndexCount, err)
	}

	if req := lang.RequiredFeatureIndex; req != 0xFFFF && int(req) >= len(t.Features) {
		return out, fmt.Errorf("invalid required feature indice %d", req)
	}

	return LangSys{
		Tag:                  record.Tag,
		RequiredFeatureIndex: lang.RequiredFeatureIndex,
		Features:             featureIndices,
	}, nil
}

// parseScript parses a single Script table. b expected to be the beginning of ScriptList.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#sTbl_lsRec
func (t *TableLayout) parseScript(b []byte, record scriptRecord) (Script, error) {
	if int(record.Offset) >= len(b) {
		return Script{}, io.ErrUnexpectedEOF
	}

	b = b[record.Offset:]
	r := bytes.NewReader(b)

	var script struct {
		DefaultLangSys uint16 // Offset to default LangSys table, from beginning of Script table — may be NULL
		LangSysCount   uint16 // Number of LangSysRecords for this script — excluding the default LangSys
		// langSysRecords[langSysCount] langSysRecord // Array of LangSysRecords, listed alphabetically by LangSys tag
	}
	if err := binary.Read(r, binary.BigEndian, &script); err != nil {
		return Script{}, fmt.Errorf("reading scriptTable: %s", err)
	}

	var defaultLang *LangSys
	var langs []LangSys

	if script.DefaultLangSys > 0 {
		def, err := t.parseLangSys(b, langSysRecord{Offset: script.DefaultLangSys})
		if err != nil {
			return Script{}, err
		}
		defaultLang = &def
	}

	for i := 0; i < int(script.LangSysCount); i++ {
		var langRecord langSysRecord
		if err := binary.Read(r, binary.BigEndian, &langRecord); err != nil {
			return Script{}, fmt.Errorf("reading langSysRecord[%d]: %s", i, err)
		}

		if langRecord.Offset == script.DefaultLangSys {
			// Don't process the same language twice
			continue
		}

		lang, err := t.parseLangSys(b, langRecord)
		if err != nil {
			return Script{}, err
		}

		langs = append(langs, lang)
	}

	return Script{
		Tag:             record.Tag,
		DefaultLanguage: defaultLang,
		Languages:       langs,
	}, nil
}

// parseScriptList parses the ScriptList.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#slTbl_sRec
func (t *TableLayout) parseScriptList(buf []byte) error {
	offset := int(t.header.ScriptListOffset)
	if offset >= len(buf) {
		return io.ErrUnexpectedEOF
	}

	b := buf[offset:]
	r := bytes.NewReader(b)

	var count uint16
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return fmt.Errorf("reading scriptCount: %s", err)
	}

	t.Scripts = make([]Script, count)
	for i := 0; i < int(count); i++ {
		var record scriptRecord
		if err := binary.Read(r, binary.BigEndian, &record); err != nil {
			return fmt.Errorf("reading scriptRecord[%d]: %s", i, err)
		}

		script, err := t.parseScript(b, record)
		if err != nil {
			return err
		}

		t.Scripts[i] = script
	}

	return nil
}

// parseFeature parses a single Feature table. b expected to be the beginning of the feature
// See https://www.microsoft.com/typography/otspec/chapter2.htm#featTbl
func parseFeature(b []byte) (Feature, error) {
	r := bytes.NewReader(b)

	var feature struct {
		FeatureParams    uint16 // = NULL (reserved for offset to FeatureParams)
		LookupIndexCount uint16 // Number of LookupList indices for this feature
		// lookupListIndices [lookupIndexCount]uint16 // Array of indices into the LookupList — zero-based (first lookup is LookupListIndex = 0)}
	}
	if err := binary.Read(r, binary.BigEndian, &feature); err != nil {
		return Feature{}, fmt.Errorf("reading featureTable: %s", err)
	}
	lookupIndices := make([]uint16, feature.LookupIndexCount)
	if err := binary.Read(r, binary.BigEndian, &lookupIndices); err != nil {
		return Feature{}, fmt.Errorf("reading featureTable: %s", err)
	}

	return Feature{paramsOffet: feature.FeatureParams, LookupIndices: lookupIndices}, nil
}

// parseFeatureList parses the FeatureList.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#flTbl
func (t *TableLayout) parseFeatureList(buf []byte) error {
	offset := int(t.header.FeatureListOffset)
	if offset >= len(buf) {
		return io.ErrUnexpectedEOF
	}

	b := buf[offset:]
	r := bytes.NewReader(b)

	var count uint16
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return fmt.Errorf("reading featureCount: %s", err)
	}

	t.Features = make([]FeatureRecord, count)
	for i := 0; i < int(count); i++ {
		var record featureRecord
		if err := binary.Read(r, binary.BigEndian, &record); err != nil {
			return fmt.Errorf("reading featureRecord[%d]: %s", i, err)
		}

		if len(b) < int(record.Offset) {
			return io.ErrUnexpectedEOF
		}
		feature, err := parseFeature(b[record.Offset:])
		if err != nil {
			return err
		}

		t.Features[i] = FeatureRecord{Tag: record.Tag, Feature: feature}
	}

	return nil
}

// parseLookup parses a single Lookup table. b expected to be the beginning of LookupList.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#featTbl
func (t *TableLayout) parseLookup(b []byte, lookupTableOffset uint16) (lookup, error) {
	if int(lookupTableOffset) >= len(b) {
		return lookup{}, io.ErrUnexpectedEOF
	}

	b = b[lookupTableOffset:]
	const tableHeaderSize = 6
	if len(b) < tableHeaderSize {
		return lookup{}, io.ErrUnexpectedEOF
	}

	type_ := binary.BigEndian.Uint16(b)
	flag := LookupFlag(binary.BigEndian.Uint16(b[2:]))
	subTableCount := binary.BigEndian.Uint16(b[4:])

	endTable := tableHeaderSize + 2*int(subTableCount)
	if len(b) < endTable {
		return lookup{}, io.ErrUnexpectedEOF
	}

	subtableOffsets := make([]uint16, subTableCount)
	for i := range subtableOffsets {
		subtableOffsets[i] = binary.BigEndian.Uint16(b[tableHeaderSize+2*i:])
	}

	out := lookup{
		kind:            type_,
		subtableOffsets: subtableOffsets,
		data:            b,
	}
	out.LookupOptions.Flag = flag

	if flag&UseMarkFilteringSet != 0 {
		if len(b) < endTable+2 {
			return lookup{}, io.ErrUnexpectedEOF
		}
		out.LookupOptions.MarkFilteringSet = binary.BigEndian.Uint16(b[endTable:])
	}

	return out, nil
}

// parseLookupList parses the LookupList.
// See https://www.microsoft.com/typography/otspec/chapter2.htm#lulTbl
func (t *TableLayout) parseLookupList(buf []byte) ([]lookup, error) {
	offset := int(t.header.LookupListOffset)
	if offset >= len(buf) {
		return nil, io.ErrUnexpectedEOF
	}

	b := buf[offset:]
	r := bytes.NewReader(b)

	var count uint16
	if err := binary.Read(r, binary.BigEndian, &count); err != nil {
		return nil, fmt.Errorf("reading lookupCount: %s", err)
	}

	lookups := make([]lookup, count)
	for i := 0; i < int(count); i++ {
		var lookupTableOffset uint16
		if err := binary.Read(r, binary.BigEndian, &lookupTableOffset); err != nil {
			return nil, fmt.Errorf("reading lookupRecord[%d]: %s", i, err)
		}

		l, err := t.parseLookup(b, lookupTableOffset)
		if err != nil {
			return nil, err
		}

		lookups[i] = l
	}

	return lookups, nil
}

type FeatureVariation struct {
	ConditionSet         []ConditionFormat1
	FeatureSubstitutions []FeatureSubstitution
}

// returns `true` if the feature is concerned by the `coords`
func (fv FeatureVariation) evaluate(coords []float32) bool {
	for _, c := range fv.ConditionSet {
		if !c.evaluate(coords) {
			return false
		}
	}
	return true
}

// parseFeatureVariationList parses the FeatureVariationList.
// See https://docs.microsoft.com/fr-fr/typography/opentype/spec/chapter2#featurevariations-table
func (t *TableLayout) parseFeatureVariationList(buf []byte) (err error) {
	if t.header.FeatureVariationsOffset == 0 {
		return nil
	}

	offset := int(t.header.FeatureVariationsOffset)
	if offset >= len(buf) {
		return io.ErrUnexpectedEOF
	}

	b := buf[offset:]
	r := bytes.NewReader(b)
	var header struct {
		versionHeader
		Count uint32
	}
	if err = binary.Read(r, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("reading FeatureVariation header: %s", err)
	}
	if len(b) < int(header.Count)*4 {
		return io.ErrUnexpectedEOF
	}

	t.FeatureVariations = make([]FeatureVariation, header.Count)
	for i := 0; i < int(header.Count); i++ {
		var record struct {
			ConditionSetOffset             uint32 // Offset to a condition set table, from beginning of FeatureVariations table.
			FeatureTableSubstitutionOffset uint32 // Offset to a feature table substitution table, from beginning of the FeatureVariations table.
		}
		if err = binary.Read(r, binary.BigEndian, &record); err != nil {
			return fmt.Errorf("reading featureVariationtRecord[%d]: %s", i, err)
		}

		if len(b) < int(record.ConditionSetOffset) || len(b) < int(record.FeatureTableSubstitutionOffset) {
			return io.ErrUnexpectedEOF
		}

		t.FeatureVariations[i].ConditionSet, err = parseConditionSet(b[record.ConditionSetOffset:])
		if err != nil {
			return err
		}

		t.FeatureVariations[i].FeatureSubstitutions, err = parseFeatureSubstitution(b[record.FeatureTableSubstitutionOffset:])
		if err != nil {
			return err
		}
	}

	return nil
}

// buf is at the begining of the condition set table
func parseConditionSet(buf []byte) ([]ConditionFormat1, error) {
	if len(buf) < 2 {
		return nil, io.ErrUnexpectedEOF
	}
	count := binary.BigEndian.Uint16(buf)
	if len(buf) < 2+int(count)*4 {
		return nil, io.ErrUnexpectedEOF
	}
	out := make([]ConditionFormat1, count)
	var err error
	for i := range out {
		offset := binary.BigEndian.Uint32(buf[2+4*i:])
		if len(buf) < int(offset) {
			return nil, io.ErrUnexpectedEOF
		}
		out[i], err = parseCondition(buf[offset:])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

type ConditionFormat1 struct {
	Axis uint16 // Index (zero-based) for the variation axis within the 'fvar' table.
	// Minimum and maximum values of the font variation instances
	// that satisfy this condition.
	Min, Max float32
}

// returns `true` if `coords` match the condition `c`
func (c ConditionFormat1) evaluate(coords []float32) bool {
	var coord float32
	if int(c.Axis) < len(coords) {
		coord = coords[c.Axis]
	}
	return c.Min <= coord && coord <= c.Max
}

// buf is at the begining of the condition
func parseCondition(buf []byte) (ConditionFormat1, error) {
	var out ConditionFormat1
	if len(buf) < 2 {
		return out, io.ErrUnexpectedEOF
	}
	format := binary.BigEndian.Uint16(buf)
	switch format {
	case 1:
		if len(buf) < 8 {
			return out, io.ErrUnexpectedEOF
		}
		out.Axis = binary.BigEndian.Uint16(buf[2:])
		out.Min = fixed214ToFloat(binary.BigEndian.Uint16(buf[4:]))
		out.Max = fixed214ToFloat(binary.BigEndian.Uint16(buf[6:]))
	default:
		return out, fmt.Errorf("invalid or unsupported condition format")
	}
	return out, nil
}

type FeatureSubstitution struct {
	AlternateFeature Feature
	FeatureIndex     uint16 // The feature table index to match.
}

// buf is as the begining of the table
func parseFeatureSubstitution(buf []byte) ([]FeatureSubstitution, error) {
	if len(buf) < 6 {
		return nil, io.ErrUnexpectedEOF
	}
	count := binary.BigEndian.Uint16(buf[4:])
	if len(buf) < 6+6*int(count) {
		return nil, io.ErrUnexpectedEOF
	}
	out := make([]FeatureSubstitution, count)
	for i := range out {
		out[i].FeatureIndex = binary.BigEndian.Uint16(buf[6+i*6:])
		alternateFeatureOffset := binary.BigEndian.Uint32(buf[6+i*6+2:])
		if len(buf) < int(alternateFeatureOffset) {
			return nil, io.ErrUnexpectedEOF
		}
		var err error
		out[i].AlternateFeature, err = parseFeature(buf[alternateFeatureOffset:])
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// parseTableLayout parses a common Layout Table used by GPOS and GSUB.
func parseTableLayout(buf []byte) (TableLayout, []lookup, error) {
	var t TableLayout

	r := bytes.NewReader(buf)
	var version versionHeader
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return t, nil, fmt.Errorf("reading layout version header: %s", err)
	}

	if version.Major != 1 {
		return t, nil, fmt.Errorf("unsupported layout major version: %d", version.Major)
	}

	switch version.Minor {
	case 0:
		if err := binary.Read(r, binary.BigEndian, &t.header.layoutHeader10); err != nil {
			return t, nil, fmt.Errorf("reading layout header: %s", err)
		}
	case 1:
		if err := binary.Read(r, binary.BigEndian, &t.header); err != nil {
			return t, nil, fmt.Errorf("reading layout header: %s", err)
		}
	default:
		return t, nil, fmt.Errorf("unsupported layout minor version: %d", version.Minor)
	}

	lookups, err := t.parseLookupList(buf)
	if err != nil {
		return t, nil, err
	}

	if err = t.parseFeatureList(buf); err != nil {
		return t, nil, err
	}

	if err = t.parseScriptList(buf); err != nil {
		return t, nil, err
	}

	if err = t.parseFeatureVariationList(buf); err != nil {
		return t, nil, err
	}

	err = t.sanitize(len(lookups))

	return t, lookups, err
}

// check that all indices are valid
func (t *TableLayout) sanitize(lookupCount int) error {
	// features
	for _, feat := range t.Features {
		for _, ind := range feat.LookupIndices {
			if int(ind) >= lookupCount {
				return fmt.Errorf("invalid lookup indice %d in features", ind)
			}
		}
	}

	// langSys
	for _, script := range t.Scripts {
		for _, lang := range script.Languages {
			for _, ind := range lang.Features {
				if int(ind) >= len(t.Features) {
					return fmt.Errorf("invalid feature indice %d in scripts", ind)
				}
			}
		}
		if lang := script.DefaultLanguage; lang != nil {
			for _, ind := range lang.Features {
				if int(ind) >= len(t.Features) {
					return fmt.Errorf("invalid feature indice %d in scripts", ind)
				}
			}
		}
	}

	// variable features
	for _, varFeat := range t.FeatureVariations {
		for _, subs := range varFeat.FeatureSubstitutions {
			if int(subs.FeatureIndex) >= len(t.Features) {
				return fmt.Errorf("invalid feature indice %d in feature variations", subs.FeatureIndex)
			}
		}
	}
	return nil
}

// shared by GSUB and GPOS

// SequenceLookup is used to specify an action (a nested lookup)
// to be applied to a glyph at a particular sequence position within the input sequence.
type SequenceLookup struct {
	InputIndex  uint16 // Index (zero-based) into the input glyph sequence
	LookupIndex uint16 // Index (zero-based) into the LookupList
}

// SequenceRule is used in Context format 1 and 2
type SequenceRule struct {
	// Starts with the second glyph
	// For format1, it is interpreted as GlyphIndex, for format 2, as ClassID
	Input   []uint16
	Lookups []SequenceLookup
}

type LookupContext1 [][]SequenceRule

func parseSequenceContext1(data []byte, lookupLength uint16) ([][]SequenceRule, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid sequence context format 1 table")
	}
	count := binary.BigEndian.Uint16(data[4:])
	if len(data) < 6+int(count)*2 {
		return nil, errors.New("invalid sequence context format 1 table")
	}

	// we dont check count against coverage since
	// "The seqRuleSetCount should match the number of glyphs in the Coverage table.
	// If these differ, the extra coverage glyphs or extra sequence rule sets are ignored."

	out := make([][]SequenceRule, int(count))
	var err error
	for i := range out {
		seqOffset := binary.BigEndian.Uint16(data[6+2*i:])
		if len(data) < int(seqOffset) {
			return nil, errors.New("invalid sequence context format 1 table")
		}
		out[i], err = parseSequenceRuleSet(data[seqOffset:], lookupLength)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// data starts at the sequenceRuleSet table
func parseSequenceRuleSet(data []byte, lookupLength uint16) ([]SequenceRule, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid sequence rule set table (EOF)")
	}
	count := binary.BigEndian.Uint16(data)
	if len(data) < 2+int(count)*2 {
		return nil, errors.New("invalid sequence rule set table (EOF)")
	}

	out := make([]SequenceRule, int(count))
	var err error
	for i := range out {
		ruleOffset := binary.BigEndian.Uint16(data[2+2*i:])
		if len(data) < int(ruleOffset) {
			return nil, errors.New("invalid sequence rule set table")
		}
		out[i], err = parseSequenceRule(data[ruleOffset:], lookupLength)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// data starts at the beginning of the list, and has always been checked for length
// `inputLength` and `lookupListLength` are used to sanitize the index access
func parseSequenceLookups(data []byte, out []SequenceLookup, inputLength, lookupListLength uint16) error {
	for i := range out {
		inputIndex := binary.BigEndian.Uint16(data[4*i:])
		if inputIndex >= inputLength {
			return fmt.Errorf("invalid sequence lookup table (input index %d for %d)", inputIndex, inputLength)
		}
		out[i].InputIndex = inputIndex
		lookupIndex := binary.BigEndian.Uint16(data[4*i+2:])
		if lookupIndex >= lookupListLength {
			return fmt.Errorf("invalid sequence lookup table (lookup index %d for %d)", lookupIndex, lookupListLength)
		}
		out[i].LookupIndex = lookupIndex
	}
	return nil
}

// data starts at the sequenceRule
func parseSequenceRule(data []byte, lookupLength uint16) (out SequenceRule, err error) {
	if len(data) < 4 {
		return out, errors.New("invalid sequence rule table header (EOF)")
	}
	glyphCount := binary.BigEndian.Uint16(data)
	lookupCount := int(binary.BigEndian.Uint16(data[2:]))
	if glyphCount == 0 {
		return out, errors.New("invalid sequence rule table (no input)")
	}
	startLookups := 4 + 2*int(glyphCount-1)
	if len(data) < startLookups+4*lookupCount {
		return out, errors.New("invalid sequence rule table length (EOF)")
	}

	out.Input, _ = parseUint16s(data[4:], int(glyphCount-1)) // length already checked

	out.Lookups = make([]SequenceLookup, lookupCount)
	err = parseSequenceLookups(data[startLookups:], out.Lookups, glyphCount, lookupLength)
	return out, err
}

type LookupContext2 struct {
	Class        Class
	SequenceSets [][]SequenceRule
}

func parseSequenceContext2(data []byte, lookupLength uint16) (out LookupContext2, err error) {
	if len(data) < 8 {
		return out, errors.New("invalid sequence context format 2 table (EOF)")
	}
	classDefOffset := binary.BigEndian.Uint16(data[4:])
	seqNumber := int(binary.BigEndian.Uint16(data[6:]))

	out.Class, err = parseClass(data, classDefOffset)
	if err != nil {
		return out, fmt.Errorf("invalid sequence context format 2 table: %s", err)
	}

	if len(data) < 8+2*seqNumber {
		return out, errors.New("invalid sequence context format 2 table (EOF)")
	}
	out.SequenceSets = make([][]SequenceRule, seqNumber)
	for i := range out.SequenceSets {
		sequenceOffset := binary.BigEndian.Uint16(data[8+2*i:])

		// "If no patterns are defined that begin with a particular class,
		// then the offset for that class value can be set to NULL."
		if sequenceOffset == 0 {
			continue
		}

		if len(data) < int(sequenceOffset) {
			return out, errors.New("invalid sequence context format 2 table (EOF)")
		}
		out.SequenceSets[i], err = parseSequenceRuleSet(data[sequenceOffset:], lookupLength)
		if err != nil {
			return out, err
		}
	}

	if needed := out.Class.Extent(); seqNumber < needed {
		// gracefully add empty sequence; needed is less than 0xFFFF + 1
		out.SequenceSets = append(out.SequenceSets, make([]SequenceRule, needed-seqNumber))
	}

	return out, nil
}

type LookupContext3 struct {
	Coverages       []Coverage
	SequenceLookups []SequenceLookup
}

func parseSequenceContext3(data []byte, lookupLength uint16) (out LookupContext3, err error) {
	if len(data) < 6 {
		return out, errors.New("invalid sequence context format 3 table")
	}
	covCount := binary.BigEndian.Uint16(data[2:])
	lookupCount := int(binary.BigEndian.Uint16(data[4:]))
	startLookups := 6 + 2*int(covCount)
	if len(data) < startLookups+4*lookupCount {
		return out, errors.New("invalid sequence context format 3 table")
	}

	out.Coverages = make([]Coverage, covCount)
	for i := range out.Coverages {
		covOffset := binary.BigEndian.Uint16(data[6+2*i:])
		out.Coverages[i], err = parseCoverage(data, uint32(covOffset))
		if err != nil {
			return out, err
		}
	}
	out.SequenceLookups = make([]SequenceLookup, lookupCount)
	err = parseSequenceLookups(data[startLookups:], out.SequenceLookups, covCount, lookupLength)
	return out, err
}

type ChainedSequenceRule struct {
	SequenceRule
	Backtrack []uint16
	Lookahead []uint16
}

type LookupChainedContext1 [][]ChainedSequenceRule

func parseChainedSequenceContext1(data []byte, lookupLength uint16) ([][]ChainedSequenceRule, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid sequence context format 1 table")
	}
	count := binary.BigEndian.Uint16(data[4:])
	if len(data) < 6+int(count)*2 {
		return nil, errors.New("invalid sequence context format 1 table")
	}

	// we dont check count against coverage since
	// "The seqRuleSetCount should match the number of glyphs in the Coverage table.
	// If these differ, the extra coverage glyphs or extra sequence rule sets are ignored."

	out := make([][]ChainedSequenceRule, int(count))
	var err error
	for i := range out {
		seqOffset := binary.BigEndian.Uint16(data[6+2*i:])
		if len(data) < int(seqOffset) {
			return nil, errors.New("invalid sequence context format 1 table")
		}
		out[i], err = parseChainedSequenceRuleSet(data[seqOffset:], lookupLength)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// data starts at the chainedSequenceRuleSet table
func parseChainedSequenceRuleSet(data []byte, lookupLength uint16) ([]ChainedSequenceRule, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid sequence rule set table")
	}
	count := binary.BigEndian.Uint16(data)
	if len(data) < 6+int(count)*2 {
		return nil, errors.New("invalid sequence rule set table")
	}

	out := make([]ChainedSequenceRule, int(count))
	var err error
	for i := range out {
		ruleOffset := binary.BigEndian.Uint16(data[2+2*i:])
		if len(data) < int(ruleOffset) {
			return nil, errors.New("invalid sequence rule set table")
		}
		out[i], err = parseChainedSequenceRule(data[ruleOffset:], lookupLength)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// data starts at the chainedSequenceRule
func parseChainedSequenceRule(data []byte, lookupLength uint16) (out ChainedSequenceRule, err error) {
	if len(data) < 2 {
		return out, errors.New("invalid chained sequence rule table header (EOF)")
	}
	backtrackGlyphCount := int(binary.BigEndian.Uint16(data))
	out.Backtrack, err = parseUint16s(data[2:], backtrackGlyphCount)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence rule table length: %s", err)
	}
	data = data[2+2*backtrackGlyphCount:]

	if len(data) < 2 {
		return out, errors.New("invalid chained sequence rule table header (EOF)")
	}
	glyphCount := binary.BigEndian.Uint16(data)
	if glyphCount == 0 {
		return out, errors.New("invalid chained sequence rule table (no input)")
	}
	out.Input, err = parseUint16s(data[2:], int(glyphCount)-1)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence rule table length: %s", err)
	}
	data = data[2+2*int(glyphCount-1):]

	if len(data) < 2 {
		return out, errors.New("invalid chained sequence rule table header (EOF)")
	}
	lookaheadGlyphCount := int(binary.BigEndian.Uint16(data))
	out.Lookahead, err = parseUint16s(data[2:], lookaheadGlyphCount)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence rule table length: %s", err)
	}
	data = data[2+2*lookaheadGlyphCount:]

	if len(data) < 2 {
		return out, errors.New("invalid chained sequence rule table header (EOF)")
	}
	lookupCount := int(binary.BigEndian.Uint16(data))
	if len(data) < 2+4*lookupCount {
		return out, errors.New("invalid chained sequence rule table length (EOF)")
	}
	out.Lookups = make([]SequenceLookup, lookupCount)
	err = parseSequenceLookups(data[2:], out.Lookups, glyphCount, lookupLength)
	return out, err
}

type LookupChainedContext2 struct {
	BacktrackClass Class
	InputClass     Class
	LookaheadClass Class
	SequenceSets   [][]ChainedSequenceRule
}

func parseChainedSequenceContext2(data []byte, lookupLength uint16) (out LookupChainedContext2, err error) {
	if len(data) < 12 {
		return out, errors.New("invalid chained sequence context format 2 table (EOF)")
	}
	backtrackDefOffset := binary.BigEndian.Uint16(data[4:])
	inputDefOffset := binary.BigEndian.Uint16(data[6:])
	lookaheadDefOffset := binary.BigEndian.Uint16(data[8:])
	seqNumber := int(binary.BigEndian.Uint16(data[10:]))

	out.BacktrackClass, err = parseClass(data, backtrackDefOffset)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence context format 2 table: %s", err)
	}
	out.LookaheadClass, err = parseClass(data, lookaheadDefOffset)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence context format 2 table: %s", err)
	}
	out.InputClass, err = parseClass(data, inputDefOffset)
	if err != nil {
		return out, fmt.Errorf("invalid chained sequence context format 2 table: %s", err)
	}

	if len(data) < 8+2*seqNumber {
		return out, errors.New("invalid chained sequence context format 2 table (EOF)")
	}
	out.SequenceSets = make([][]ChainedSequenceRule, seqNumber)
	for i := range out.SequenceSets {
		sequenceOffset := binary.BigEndian.Uint16(data[12+2*i:])

		// "If no patterns are defined that begin with a particular class,
		// then the offset for that class value can be set to NULL."
		if sequenceOffset == 0 {
			continue
		}

		if len(data) < int(sequenceOffset) {
			return out, errors.New("invalid chained sequence context format 2 table (EOF)")
		}
		out.SequenceSets[i], err = parseChainedSequenceRuleSet(data[sequenceOffset:], lookupLength)
		if err != nil {
			return out, err
		}
	}

	if needed := out.InputClass.Extent(); seqNumber < needed {
		// gracefully add empty sequence; needed is less than 0xFFFF + 1
		out.SequenceSets = append(out.SequenceSets, make([]ChainedSequenceRule, needed-seqNumber))
	}

	return out, nil
}

type LookupChainedContext3 struct {
	Backtrack       []Coverage
	Input           []Coverage
	Lookahead       []Coverage
	SequenceLookups []SequenceLookup
}

func parseChainedSequenceContext3(data []byte, lookupLength uint16) (out LookupChainedContext3, err error) {
	if len(data) < 4 {
		return out, errors.New("invalid chained sequence context format 3 table")
	}
	covCount := binary.BigEndian.Uint16(data[2:])
	out.Backtrack = make([]Coverage, covCount)
	for i := range out.Backtrack {
		covOffset := binary.BigEndian.Uint16(data[4+2*i:])
		out.Backtrack[i], err = parseCoverage(data, uint32(covOffset))
		if err != nil {
			return out, err
		}
	}
	endBacktrack := 4 + 2*int(covCount)

	if len(data) < endBacktrack+2 {
		return out, errors.New("invalid chained sequence context format 3 table")
	}
	inputCount := binary.BigEndian.Uint16(data[endBacktrack:])
	out.Input = make([]Coverage, inputCount)
	for i := range out.Input {
		covOffset := binary.BigEndian.Uint16(data[endBacktrack+2+2*i:])
		out.Input[i], err = parseCoverage(data, uint32(covOffset))
		if err != nil {
			return out, err
		}
	}
	endInput := endBacktrack + 2 + 2*int(inputCount)

	if len(data) < endInput+2 {
		return out, errors.New("invalid chained sequence context format 3 table")
	}
	covCount = binary.BigEndian.Uint16(data[endInput:])
	out.Lookahead = make([]Coverage, covCount)
	for i := range out.Lookahead {
		covOffset := binary.BigEndian.Uint16(data[endInput+2+2*i:])
		out.Lookahead[i], err = parseCoverage(data, uint32(covOffset))
		if err != nil {
			return out, err
		}
	}
	endLookahead := endInput + 2 + 2*int(covCount)

	if len(data) < endLookahead+2 {
		return out, errors.New("invalid chained sequence context format 3 table")
	}
	lookupCount := int(binary.BigEndian.Uint16(data[endLookahead:]))
	if len(data) < endLookahead+2+4*lookupCount {
		return out, errors.New("invalid chained sequence context format 3 table")
	}
	out.SequenceLookups = make([]SequenceLookup, lookupCount)
	err = parseSequenceLookups(data[endLookahead+2:], out.SequenceLookups, inputCount, lookupLength)
	return out, err
}
