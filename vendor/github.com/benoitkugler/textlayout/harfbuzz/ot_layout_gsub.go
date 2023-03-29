package harfbuzz

import (
	"fmt"
	"math/bits"

	"github.com/benoitkugler/textlayout/fonts"
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

const maxContextLength = 64

var _ layoutLookup = lookupGSUB{}

// implements layoutLookup
type lookupGSUB tt.LookupGSUB

func (l lookupGSUB) collectCoverage(dst *setDigest) {
	for _, table := range l.Subtables {
		dst.collectCoverage(table.Coverage)
	}
}

func (l lookupGSUB) dispatchSubtables(ctx *getSubtablesContext) {
	for _, table := range l.Subtables {
		*ctx = append(*ctx, newGSUBApplicable(table))
	}
}

func (l lookupGSUB) dispatchApply(ctx *otApplyContext) bool {
	for _, table := range l.Subtables {
		if gsubSubtable(table).apply(ctx) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) wouldApply(ctx *wouldApplyContext, accel *otLayoutLookupAccelerator) bool {
	if len(ctx.glyphs) == 0 {
		return false
	}
	if !accel.digest.mayHave(ctx.glyphs[0]) {
		return false
	}
	// dispatch on subtables
	for _, table := range l.Subtables {
		if gsubSubtable(table).wouldApply(ctx) {
			return true
		}
	}
	return false
}

func (l lookupGSUB) isReverse() bool { return l.Type == tt.GSUBReverse }

func applyRecurseGSUB(c *otApplyContext, lookupIndex uint16) bool {
	gsub := c.font.otTables.GSUB
	l := lookupGSUB(gsub.Lookups[lookupIndex])
	return c.applyRecurseLookup(lookupIndex, l)
}

//  implements `hb_apply_func_t`
type gsubSubtable tt.GSUBSubtable

// return `true` is we should apply this lookup to the glyphs in `c`,
// which are assumed to be non empty
func (table gsubSubtable) wouldApply(c *wouldApplyContext) bool {
	index, ok := table.Coverage.Index(c.glyphs[0])
	switch data := table.Data.(type) {
	case tt.GSUBSingle1, tt.GSUBSingle2, tt.GSUBMultiple1, tt.GSUBAlternate1, tt.GSUBReverseChainedContext1:
		return len(c.glyphs) == 1 && ok

	case tt.GSUBLigature1:
		if !ok {
			return false
		}
		ligatureSet := data[index]
		glyphsFromSecond := c.glyphs[1:]
		for _, ligature := range ligatureSet {
			if ligature.Matches(glyphsFromSecond) {
				return true
			}
		}
		return false

	case tt.GSUBContext1:
		return c.wouldApplyLookupContext1(tt.LookupContext1(data), index)
	case tt.GSUBContext2:
		return c.wouldApplyLookupContext2(tt.LookupContext2(data), index, c.glyphs[0])
	case tt.GSUBContext3:
		return c.wouldApplyLookupContext3(tt.LookupContext3(data), index)
	case tt.GSUBChainedContext1:
		return c.wouldApplyLookupChainedContext1(tt.LookupChainedContext1(data), index)
	case tt.GSUBChainedContext2:
		return c.wouldApplyLookupChainedContext2(tt.LookupChainedContext2(data), index, c.glyphs[0])
	case tt.GSUBChainedContext3:
		return c.wouldApplyLookupChainedContext3(tt.LookupChainedContext3(data), index)
	}
	return false
}

// return `true` is the subsitution found a match and was applied
func (table gsubSubtable) apply(c *otApplyContext) bool {
	glyph := c.buffer.cur(0)
	glyphID := glyph.Glyph
	index, ok := table.Coverage.Index(glyphID)
	if !ok {
		return false
	}

	if debugMode >= 2 {
		fmt.Printf("\tAPPLY - type %T at index %d\n", table.Data, c.buffer.idx)
	}

	switch data := table.Data.(type) {
	case tt.GSUBSingle1:
		/* According to the Adobe Annotated OpenType Suite, result is always
		* limited to 16bit. */
		glyphID = fonts.GID(uint16(int(glyphID) + int(data)))
		c.replaceGlyph(glyphID)
	case tt.GSUBSingle2:
		if index >= len(data) { // index is not sanitized in tt.Parse
			return false
		}
		c.replaceGlyph(data[index])

	case tt.GSUBMultiple1:
		c.applySubsSequence(data[index])

	case tt.GSUBAlternate1:
		alternates := data[index]
		return c.applySubsAlternate(alternates)

	case tt.GSUBLigature1:
		ligatureSet := data[index]
		return c.applySubsLigature(ligatureSet)

	case tt.GSUBContext1:
		return c.applyLookupContext1(tt.LookupContext1(data), index)
	case tt.GSUBContext2:
		return c.applyLookupContext2(tt.LookupContext2(data), index, glyphID)
	case tt.GSUBContext3:
		return c.applyLookupContext3(tt.LookupContext3(data), index)
	case tt.GSUBChainedContext1:
		return c.applyLookupChainedContext1(tt.LookupChainedContext1(data), index)
	case tt.GSUBChainedContext2:
		return c.applyLookupChainedContext2(tt.LookupChainedContext2(data), index, glyphID)
	case tt.GSUBChainedContext3:
		return c.applyLookupChainedContext3(tt.LookupChainedContext3(data), index)

	case tt.GSUBReverseChainedContext1:
		if c.nestingLevelLeft != maxNestingLevel {
			return false // no chaining to this type
		}
		lB, lL := len(data.Backtrack), len(data.Lookahead)
		hasMatch, startIndex := c.matchBacktrack(get1N(&c.indices, 0, lB), matchCoverage(data.Backtrack))
		if !hasMatch {
			return false
		}

		hasMatch, endIndex := c.matchLookahead(get1N(&c.indices, 0, lL), matchCoverage(data.Lookahead), 1)
		if !hasMatch {
			return false
		}

		c.buffer.unsafeToBreakFromOutbuffer(startIndex, endIndex)
		c.setGlyphProps(data.Substitutes[index])
		c.buffer.cur(0).Glyph = data.Substitutes[index]
		/* Note: We DON'T decrease buffer.idx.  The main loop does it
		 * for us.  This is useful for preventing surprises if someone
		 * calls us through a Context lookup. */

	}

	return true
}

func (c *otApplyContext) applySubsSequence(seq []fonts.GID) {
	/* Special-case to make it in-place and not consider this
	 * as a "multiplied" substitution. */
	switch len(seq) {
	case 1:
		c.replaceGlyph(seq[0])
	case 0:
		/* Spec disallows this, but Uniscribe allows it.
		 * https://github.com/harfbuzz/harfbuzz/issues/253 */
		c.buffer.deleteGlyph()
	default:
		var klass uint16
		if c.buffer.cur(0).isLigature() {
			klass = tt.BaseGlyph
		}
		ligID := c.buffer.cur(0).getLigID()
		for i, g := range seq {
			/* If is attached to a ligature, don't disturb that.
			 * https://github.com/harfbuzz/harfbuzz/issues/3069 */
			if ligID == 0 {
				c.buffer.cur(0).setLigPropsForMark(0, uint8(i))
			}
			c.setGlyphPropsExt(g, klass, false, true)
			c.buffer.outputGlyphIndex(g)
		}
		c.buffer.skipGlyph()
	}
}

func (c *otApplyContext) applySubsAlternate(alternates []fonts.GID) bool {
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

	c.replaceGlyph(alternates[altIndex-1])
	return true
}

func (c *otApplyContext) applySubsLigature(ligatureSet []tt.LigatureGlyph) bool {
	for _, lig := range ligatureSet {
		count := len(lig.Components) + 1

		/* Special-case to make it in-place and not consider this
		 * as a "ligated" substitution. */
		if count == 1 {
			c.replaceGlyph(lig.Glyph)
			return true
		}

		var matchPositions [maxContextLength]int

		ok, matchLength, totalComponentCount := c.matchInput(lig.Components, matchGlyph, &matchPositions)
		if !ok {
			continue
		}
		c.ligateInput(count, matchPositions, matchLength, lig.Glyph, totalComponentCount)

		return true
	}
	return false
}
