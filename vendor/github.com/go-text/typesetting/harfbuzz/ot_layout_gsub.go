package harfbuzz

import (
	"fmt"
	"math/bits"

	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/tables"
)

const maxContextLength = 64

var _ layoutLookup = lookupGSUB{}

// implements layoutLookup
type lookupGSUB font.GSUBLookup

func (l lookupGSUB) Props() uint32 { return l.LookupOptions.Props() }

func (l lookupGSUB) collectCoverage(dst *setDigest) {
	for _, table := range l.Subtables {
		dst.collectCoverage(table.Cov())
	}
}

func (l lookupGSUB) dispatchSubtables(ctx *getSubtablesContext) {
	for _, table := range l.Subtables {
		*ctx = append(*ctx, newGSUBApplicable(table))
	}
}

func (l lookupGSUB) dispatchApply(ctx *otApplyContext) bool {
	for _, table := range l.Subtables {
		if ctx.applyGSUB(table) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) wouldApply(ctx *wouldApplyContext, accel *otLayoutLookupAccelerator) bool {
	if len(ctx.glyphs) == 0 {
		return false
	}
	if !accel.digest.mayHave(gID(ctx.glyphs[0])) {
		return false
	}
	// dispatch on subtables
	for _, table := range l.Subtables {
		if ctx.wouldApplyGSUB(table) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) isReverse() bool {
	if len(l.Subtables) == 0 {
		return false
	}
	_, is := l.Subtables[0].(tables.ReverseChainSingleSubs)
	return is
}

func applyRecurseGSUB(c *otApplyContext, lookupIndex uint16) bool {
	gsub := c.font.face.GSUB
	l := lookupGSUB(gsub.Lookups[lookupIndex])
	return c.applyRecurseLookup(lookupIndex, l)
}

// matchesLigature tests if the ligature should be applied on `glyphsFromSecond`,
// which starts from the second glyph.
func matchesLigature(l tables.Ligature, glyphsFromSecond []GID) bool {
	if len(glyphsFromSecond) != len(l.ComponentGlyphIDs) {
		return false
	}
	for i, g := range glyphsFromSecond {
		if g != GID(l.ComponentGlyphIDs[i]) {
			return false
		}
	}
	return true
}

// return `true` is we should apply this lookup to the glyphs in `c`,
// which are assumed to be non empty
func (c *wouldApplyContext) wouldApplyGSUB(table tables.GSUBLookup) bool {
	index, ok := table.Cov().Index(gID(c.glyphs[0]))
	switch data := table.(type) {
	case tables.SingleSubs, tables.MultipleSubs, tables.AlternateSubs, tables.ReverseChainSingleSubs:
		return len(c.glyphs) == 1 && ok

	case tables.LigatureSubs:
		if !ok {
			return false
		}
		ligatureSet := data.LigatureSets[index].Ligatures
		glyphsFromSecond := c.glyphs[1:]
		for _, ligature := range ligatureSet {
			if matchesLigature(ligature, glyphsFromSecond) {
				return true
			}
		}
		return false

	case tables.ContextualSubs:
		switch inner := data.Data.(type) {
		case tables.ContextualSubs1:
			return c.wouldApplyLookupContext1(tables.SequenceContextFormat1(inner), index)
		case tables.ContextualSubs2:
			return c.wouldApplyLookupContext2(tables.SequenceContextFormat2(inner), index, c.glyphs[0])
		case tables.ContextualSubs3:
			return c.wouldApplyLookupContext3(tables.SequenceContextFormat3(inner), index)
		}

	case tables.ChainedContextualSubs:
		switch inner := data.Data.(type) {
		case tables.ChainedContextualSubs1:
			return c.wouldApplyLookupChainedContext1(tables.ChainedSequenceContextFormat1(inner), index)
		case tables.ChainedContextualSubs2:
			return c.wouldApplyLookupChainedContext2(tables.ChainedSequenceContextFormat2(inner), index, c.glyphs[0])
		case tables.ChainedContextualSubs3:
			return c.wouldApplyLookupChainedContext3(tables.ChainedSequenceContextFormat3(inner), index)
		}

	}
	return false
}

// return `true` is the subsitution found a match and was applied
func (c *otApplyContext) applyGSUB(table tables.GSUBLookup) bool {
	glyph := c.buffer.cur(0)
	glyphID := glyph.Glyph
	index, ok := table.Cov().Index(gID(glyphID))
	if !ok {
		return false
	}

	if debugMode >= 2 {
		fmt.Printf("\tAPPLY - type %T at index %d\n", table, c.buffer.idx)
	}

	switch data := table.(type) {
	case tables.SingleSubs:
		switch inner := data.Data.(type) {
		case tables.SingleSubstData1:
			/* According to the Adobe Annotated OpenType Suite, result is always
			* limited to 16bit. */
			glyphID = GID(uint16(int(glyphID) + int(inner.DeltaGlyphID)))
			c.replaceGlyph(glyphID)
		case tables.SingleSubstData2:
			if index >= len(inner.SubstituteGlyphIDs) { // index is not sanitized in tables.Parse
				return false
			}
			c.replaceGlyph(GID(inner.SubstituteGlyphIDs[index]))
		}

	case tables.MultipleSubs:
		c.applySubsSequence(data.Sequences[index].SubstituteGlyphIDs)

	case tables.AlternateSubs:
		alternates := data.AlternateSets[index].AlternateGlyphIDs
		return c.applySubsAlternate(alternates)

	case tables.LigatureSubs:
		ligatureSet := data.LigatureSets[index].Ligatures
		return c.applySubsLigature(ligatureSet)

	case tables.ContextualSubs:
		switch inner := data.Data.(type) {
		case tables.ContextualSubs1:
			return c.applyLookupContext1(tables.SequenceContextFormat1(inner), index)
		case tables.ContextualSubs2:
			return c.applyLookupContext2(tables.SequenceContextFormat2(inner), index, glyphID)
		case tables.ContextualSubs3:
			return c.applyLookupContext3(tables.SequenceContextFormat3(inner), index)
		}

	case tables.ChainedContextualSubs:
		switch inner := data.Data.(type) {
		case tables.ChainedContextualSubs1:
			return c.applyLookupChainedContext1(tables.ChainedSequenceContextFormat1(inner), index)
		case tables.ChainedContextualSubs2:
			return c.applyLookupChainedContext2(tables.ChainedSequenceContextFormat2(inner), index, glyphID)
		case tables.ChainedContextualSubs3:
			return c.applyLookupChainedContext3(tables.ChainedSequenceContextFormat3(inner), index)
		}

	case tables.ReverseChainSingleSubs:
		if c.nestingLevelLeft != maxNestingLevel {
			return false // no chaining to this type
		}
		lB, lL := len(data.BacktrackCoverages), len(data.LookaheadCoverages)
		hasMatch, startIndex := c.matchBacktrack(get1N(&c.indices, 0, lB), matchCoverage(data.BacktrackCoverages))
		if !hasMatch {
			return false
		}

		hasMatch, endIndex := c.matchLookahead(get1N(&c.indices, 0, lL), matchCoverage(data.LookaheadCoverages), 1)
		if !hasMatch {
			return false
		}

		c.buffer.unsafeToBreakFromOutbuffer(startIndex, endIndex)
		c.setGlyphProps(GID(data.SubstituteGlyphIDs[index]))
		c.buffer.cur(0).Glyph = GID(data.SubstituteGlyphIDs[index])
		// Note: We DON'T decrease buffer.idx.  The main loop does it
		// for us.  This is useful for preventing surprises if someone
		// calls us through a Context lookup.
	}

	return true
}

func (c *otApplyContext) applySubsSequence(seq []gID) {
	/* Special-case to make it in-place and not consider this
	 * as a "multiplied" substitution. */
	switch len(seq) {
	case 1:
		c.replaceGlyph(GID(seq[0]))
	case 0:
		/* Spec disallows this, but Uniscribe allows it.
		 * https://github.com/harfbuzz/harfbuzz/issues/253 */
		c.buffer.deleteGlyph()
	default:
		var klass uint16
		if c.buffer.cur(0).isLigature() {
			klass = tables.GPBaseGlyph
		}
		ligID := c.buffer.cur(0).getLigID()
		for i, g := range seq {
			// If is attached to a ligature, don't disturb that.
			// https://github.com/harfbuzz/harfbuzz/issues/3069
			if ligID == 0 {
				c.buffer.cur(0).setLigPropsForMark(0, uint8(i))
			}
			c.setGlyphPropsExt(GID(g), klass, false, true)
			c.buffer.outputGlyphIndex(GID(g))
		}
		c.buffer.skipGlyph()
	}
}

func (c *otApplyContext) applySubsAlternate(alternates []gID) bool {
	count := uint32(len(alternates))
	if count == 0 {
		return false
	}

	glyphMask := c.buffer.cur(0).Mask
	lookupMask := c.lookupMask

	/* Note: This breaks badly if two features enabled this lookup together. */

	shift := bits.TrailingZeros32(lookupMask)
	altIndex := (lookupMask & glyphMask) >> shift

	/* If altIndex is MAX_VALUE, randomize feature if it is the rand feature. */
	if altIndex == otMapMaxValue && c.random {
		// Maybe we can do better than unsafe-to-break all; but since we are
		// changing random state, it would be hard to track that.  Good 'nough.
		c.buffer.unsafeToBreak(0, len(c.buffer.Info))
		altIndex = c.randomNumber()%count + 1
	}

	if altIndex > count || altIndex == 0 {
		return false
	}

	c.replaceGlyph(GID(alternates[altIndex-1]))
	return true
}

func (c *otApplyContext) applySubsLigature(ligatureSet []tables.Ligature) bool {
	for _, lig := range ligatureSet {
		count := len(lig.ComponentGlyphIDs) + 1

		// Special-case to make it in-place and not consider this
		// as a "ligated" substitution.
		if count == 1 {
			c.replaceGlyph(GID(lig.LigatureGlyph))
			return true
		}

		var matchPositions [maxContextLength]int

		ok, matchLength, totalComponentCount := c.matchInput(lig.ComponentGlyphIDs, matchGlyph, &matchPositions)
		if !ok {
			continue
		}
		c.ligateInput(count, matchPositions, matchLength, lig.LigatureGlyph, totalComponentCount)

		return true
	}
	return false
}
