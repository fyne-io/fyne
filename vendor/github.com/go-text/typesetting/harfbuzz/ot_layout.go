package harfbuzz

import (
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from src/hb-ot-layout.cc, hb-ot-layout.hh
// Copyright © 1998-2004  David Turner and Werner Lemberg
// Copyright © 2006  2007,2008,2009  Red Hat, Inc. 2012,2013  Google, Inc. Behdad Esfahbod

const (
	// This bit relates only to the correct processing of
	// the cursive attachment lookup type (GPOS lookup type 3).
	// When this bit is set, the last glyph in a given sequence to
	// which the cursive attachment lookup is applied, will be positioned on the baseline.
	otRightToLeft      uint16 = 1 << iota
	otIgnoreBaseGlyphs        // If set, skips over base glyphs
	otIgnoreLigatures         // If set, skips over ligatures
	otIgnoreMarks             // If set, skips over all combining marks
	// If set, indicates that the lookup table structure
	// is followed by a MarkFilteringSet field.
	// The layout engine skips over all mark glyphs not in the mark filtering set indicated.
	_
	_ uint16 = 0x00E0 // For future use (Set to zero)
	// If not zero, skips over all marks of attachment
	// type different from specified.
	otMarkAttachmentType uint16 = 0xFF00
)

//  /**
//   * SECTION:hb-ot-layout
//   * @title: hb-ot-layout
//   * @short_description: OpenType Layout
//   * @include: hb-ot.h
//   *
//   * Functions for querying OpenType Layout features in the font face.
//   **/

const maxNestingLevel = 6

func (c *otApplyContext) applyString(proxy otProxyMeta, accel *otLayoutLookupAccelerator) {
	buffer := c.buffer
	lookup := accel.lookup

	if len(buffer.Info) == 0 || c.lookupMask == 0 {
		return
	}
	c.setLookupProps(lookup.Props())
	if !lookup.isReverse() {
		// in/out forward substitution/positioning
		if !proxy.inplace {
			buffer.clearOutput()
		}
		buffer.idx = 0

		c.applyForward(accel)
		if !proxy.inplace {
			buffer.swapBuffers()
		}
	} else {
		/* in-place backward substitution/positioning */
		// assert (!buffer->have_output);

		buffer.idx = len(buffer.Info) - 1

		c.applyBackward(accel)
	}
}

func (c *otApplyContext) applyForward(accel *otLayoutLookupAccelerator) bool {
	ret := false
	buffer := c.buffer
	for buffer.idx < len(buffer.Info) {
		applied := false
		if accel.digest.mayHave(gID(buffer.cur(0).Glyph)) &&
			(buffer.cur(0).Mask&c.lookupMask) != 0 &&
			c.checkGlyphProperty(buffer.cur(0), c.lookupProps) {
			applied = accel.apply(c)
		}

		if applied {
			ret = true
		} else {
			buffer.nextGlyph()
		}
	}
	return ret
}

func (c *otApplyContext) applyBackward(accel *otLayoutLookupAccelerator) bool {
	ret := false
	buffer := c.buffer
	for do := true; do; do = buffer.idx >= 0 {
		if accel.digest.mayHave(gID(buffer.cur(0).Glyph)) &&
			(buffer.cur(0).Mask&c.lookupMask != 0) &&
			c.checkGlyphProperty(buffer.cur(0), c.lookupProps) {
			applied := accel.apply(c)
			ret = ret || applied
		}

		// the reverse lookup doesn't "advance" cursor (for good reason).
		buffer.idx--

	}
	return ret
}

/*
 * kern
 */

// tests whether a face includes any state-machine kerning in the 'kern' table.
//
// Does NOT examine the GPOS table.
func hasMachineKerning(kern font.Kernx) bool {
	for _, subtable := range kern {
		if _, isType1 := subtable.Data.(font.Kern1); isType1 {
			return true
		}
	}
	return false
}

// tests whether a face has any cross-stream kerning (i.e., kerns
// that make adjustments perpendicular to the direction of the text
// flow: Y adjustments in horizontal text or X adjustments in
// vertical text) in the 'kern' table.
//
// Does NOT examine the GPOS table.
func hasCrossKerning(kern font.Kernx) bool {
	for _, subtable := range kern {
		if subtable.IsCrossStream() {
			return true
		}
	}
	return false
}

func (sp *otShapePlan) otLayoutKern(font *Font, buffer *Buffer) {
	kern := font.face.Kern
	c := newAatApplyContext(sp, font, buffer)
	c.applyKernx(kern)
}

var otTagLatinScript = loader.NewTag('l', 'a', 't', 'n')

// SelectScript selects an OpenType script from the `scriptTags` array,
// returning its index in the Scripts slice and the script tag.
//
// If `table` does not have any of the requested scripts, then `DFLT`,
// `dflt`, and `latn` tags are tried in that order. If the table still does not
// have any of these scripts, NoScriptIndex is returned.
//
// An additional boolean if returned : it is `true` if one of the requested scripts is selected, or `false` if a fallback
// script is selected or if no scripts are selected.
func SelectScript(table *font.Layout, scriptTags []tables.Tag) (int, tables.Tag, bool) {
	for _, tag := range scriptTags {
		if scriptIndex := table.FindScript(tag); scriptIndex != -1 {
			return scriptIndex, tag, true
		}
	}

	// try finding 'DFLT'
	if scriptIndex := table.FindScript(tagDefaultScript); scriptIndex != -1 {
		return scriptIndex, tagDefaultScript, false
	}

	// try with 'dflt'; MS site has had typos and many fonts use it now :(
	if scriptIndex := table.FindScript(tagDefaultLanguage); scriptIndex != -1 {
		return scriptIndex, tagDefaultLanguage, false
	}

	// try with 'latn'; some old fonts put their features there even though
	// they're really trying to support Thai, for example :(
	if scriptIndex := table.FindScript(otTagLatinScript); scriptIndex != -1 {
		return scriptIndex, otTagLatinScript, false
	}

	return NoScriptIndex, NoScriptIndex, false
}

// SelectLanguage fetches the index of the first language tag from `languageTags` in the specified layout table,
// underneath `scriptIndex`.
// It not found, the `dflt` language tag is searched.
// Return `true` if the requested language tag is found, `false` otherwise.
// If `scriptIndex` is `NoScriptIndex` or if no language is found, `DefaultLanguageIndex` is returned.
func SelectLanguage(table *font.Layout, scriptIndex int, languageTags []tables.Tag) (int, bool) {
	if scriptIndex == NoScriptIndex {
		return DefaultLanguageIndex, false
	}

	s := table.Scripts[scriptIndex]

	for _, lang := range languageTags {
		if languageIndex := s.FindLanguage(lang); languageIndex != -1 {
			return languageIndex, true
		}
	}

	// try finding 'dflt'
	if languageIndex := s.FindLanguage(tagDefaultLanguage); languageIndex != -1 {
		return languageIndex, false
	}

	return DefaultLanguageIndex, false
}

func findFeature(g *font.Layout, featureTag tables.Tag) uint16 {
	if index, ok := g.FindFeatureIndex(featureTag); ok {
		return index
	}
	return NoFeatureIndex
}

// Fetches the index of a given feature tag in the specified face's GSUB table
// or GPOS table, underneath the specified script and language.
// Return `NoFeatureIndex` it the the feature is not found.
func FindFeatureForLang(table *font.Layout, scriptIndex, languageIndex int, featureTag tables.Tag) uint16 {
	if scriptIndex == NoScriptIndex {
		return NoFeatureIndex
	}

	l := table.Scripts[scriptIndex].GetLangSys(uint16(languageIndex))
	for _, fIndex := range l.FeatureIndices {
		if featureTag == table.Features[fIndex].Tag {
			return fIndex
		}
	}

	return NoFeatureIndex
}

// Fetches the tag of a requested feature index in the given layout table,
// underneath the specified script and language. Returns -1 if no feature is requested.
func getRequiredFeature(g *font.Layout, scriptIndex, languageIndex int) (uint16, tables.Tag) {
	if scriptIndex == NoScriptIndex || languageIndex == DefaultLanguageIndex {
		return NoFeatureIndex, 0
	}

	l := g.Scripts[scriptIndex].LangSys[languageIndex]
	if l.RequiredFeatureIndex == 0xFFFF {
		return NoFeatureIndex, 0
	}
	index := l.RequiredFeatureIndex
	return index, g.Features[index].Tag
}

// getFeatureLookupsWithVar fetches a list of all lookups enumerated for the specified feature, in
// the given table, enabled at the specified variations index.
// it returns the basic feature if `variationsIndex == noVariationsIndex`
func getFeatureLookupsWithVar(table *font.Layout, featureIndex uint16, variationsIndex int) []uint16 {
	if featureIndex == NoFeatureIndex {
		return nil
	}

	if variationsIndex == noVariationsIndex { // just fetch the feature
		return table.Features[featureIndex].LookupListIndices
	}

	// hook the variations
	subs := table.FeatureVariations[variationsIndex].Substitutions.Substitutions
	for _, sub := range subs {
		if sub.FeatureIndex == featureIndex {
			return sub.AlternateFeature.LookupListIndices
		}
	}
	return nil
}

// tests whether a specified lookup index in the specified face would
// trigger a substitution on the given glyph sequence.
// zeroContext indicating whether substitutions should be context-free.
func otLayoutLookupWouldSubstitute(font *Font, lookupIndex uint16, glyphs []GID, zeroContext bool) bool {
	gsub := font.face.GSUB
	if int(lookupIndex) >= len(gsub.Lookups) {
		return false
	}
	c := wouldApplyContext{glyphs, nil, zeroContext}

	l := lookupGSUB(gsub.Lookups[lookupIndex])
	return l.wouldApply(&c, &font.gsubAccels[lookupIndex])
}

// Called before substitution lookups are performed, to ensure that glyph
// class and other properties are set on the glyphs in the buffer.
func layoutSubstituteStart(font *Font, buffer *Buffer) {
	gdef := font.face.GDEF
	hasClass := gdef.GlyphClassDef != nil
	for i := range buffer.Info {
		if hasClass {
			buffer.Info[i].glyphProps = gdef.GlyphProps(gID(buffer.Info[i].Glyph))
		}
		buffer.Info[i].ligProps = 0
		buffer.Info[i].syllable = 0
	}
}

func otLayoutDeleteGlyphsInplace(buffer *Buffer, filter func(*GlyphInfo) bool) {
	// Merge clusters and delete filtered glyphs.
	var (
		j    int
		info = buffer.Info
		pos  = buffer.Pos
	)
	for i := range info {
		if filter(&info[i]) {
			/* Merge clusters.
			* Same logic as buffer.delete_glyph(), but for in-place removal. */

			cluster := info[i].Cluster
			if i+1 < len(buffer.Info) && cluster == info[i+1].Cluster {
				/* Cluster survives; do nothing. */
				continue
			}

			if j != 0 {
				/* Merge cluster backward. */
				if cluster < info[j-1].Cluster {
					mask := info[i].Mask
					oldCluster := info[j-1].Cluster
					for k := j; k != 0 && info[k-1].Cluster == oldCluster; k-- {
						info[k-1].setCluster(cluster, mask)
					}
				}
				continue
			}

			if i+1 < len(buffer.Info) {
				/* Merge cluster forward. */
				buffer.mergeClusters(i, i+2)
			}

			continue
		}

		if j != i {
			info[j] = info[i]
			pos[j] = pos[i]
		}
		j++
	}
	buffer.Info = buffer.Info[:j]
	buffer.Pos = buffer.Pos[:j]
}

// Called before positioning lookups are performed, to ensure that glyph
// attachment types and glyph-attachment chains are set for the glyphs in the buffer.
func otLayoutPositionStart(_ *Font, buffer *Buffer) {
	positionStartGPOS(buffer)
}

// Called after positioning lookups are performed, to finish glyph offsets.
func otLayoutPositionFinishOffsets(_ *Font, buffer *Buffer) {
	positionFinishOffsetsGPOS(buffer)
}

func clearSyllables(_ *otShapePlan, _ *Font, buffer *Buffer) {
	info := buffer.Info
	for i := range info {
		info[i].syllable = 0
	}
}

func glyphInfoSubstituted(info *GlyphInfo) bool {
	return (info.glyphProps & substituted) != 0
}

func clearSubstitutionFlags(_ *otShapePlan, _ *Font, buffer *Buffer) {
	info := buffer.Info
	for i := range info {
		info[i].glyphProps &= ^substituted
	}
}

func reverseGraphemes(b *Buffer) {
	b.reverseGroups(func(_, gi2 *GlyphInfo) bool { return gi2.isContinuation() }, b.ClusterLevel == MonotoneGraphemes)
}
