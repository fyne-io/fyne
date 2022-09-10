package harfbuzz

import (
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
)

// ported from harfbuzz/src/hb-ot-shape-normalize.cc Copyright © 2011,2012  Google, Inc. Behdad Esfahbod

/*
 * HIGHLEVEL DESIGN:
 *
 * This file exports one main function: otShapeNormalize().
 *
 * This function closely reflects the Unicode Normalization Algorithm,
 * yet it's different.
 *
 * Each shaper specifies whether it prefers decomposed (NFD) or composed (NFC).
 * The logic however tries to use whatever the font can support.
 *
 * In general what happens is that: each grapheme is decomposed in a chain
 * of 1:2 decompositions, marks reordered, and then recomposed if desired,
 * so far it's like Unicode Normalization.  However, the decomposition and
 * recomposition only happens if the font supports the resulting characters.
 *
 * The goals are:
 *
 *   - Try to render all canonically equivalent strings similarly.  To really
 *     achieve this we have to always do the full decomposition and then
 *     selectively recompose from there.  It's kinda too expensive though, so
 *     we skip some cases.  For example, if composed is desired, we simply
 *     don't touch 1-character clusters that are supported by the font, even
 *     though their NFC may be different.
 *
 *   - When a font has a precomposed character for a sequence but the 'ccmp'
 *     feature in the font is not adequate, use the precomposed character
 *     which typically has better mark positioning.
 *
 *   - When a font does not support a combining mark, but supports it precomposed
 *     with previous base, use that.  This needs the itemizer to have this
 *     knowledge too.  We need to provide assistance to the itemizer.
 *
 *   - When a font does not support a character but supports its canonical
 *     decomposition, well, use the decomposition.
 *
 *   - The complex shapers can customize the compose and decompose functions to
 *     offload some of their requirements to the normalizer.  For example, the
 *     Indic shaper may want to disallow recomposing of two matras.
 */

const shapeComplexMaxCombiningMarks = 32

type normalizationMode uint8

const (
	nmNone normalizationMode = iota
	nmDecomposed
	nmComposedDiacritics               // never composes base-to-base
	nmComposedDiacriticsNoShortCircuit // always fully decomposes and then recompose back

	nmAuto    // see below for logic.
	nmDefault = nmAuto
)

type otNormalizeContext struct {
	plan   *otShapePlan
	buffer *Buffer
	font   *Font
	// hb_unicode_funcs_t *unicode;
	decompose func(c *otNormalizeContext, ab rune) (a, b rune, ok bool)
	compose   func(c *otNormalizeContext, a, b rune) (ab rune, ok bool)
}

func setGlyph(info *GlyphInfo, font *Font) {
	info.Glyph, _ = font.face.NominalGlyph(info.codepoint)
}

func outputChar(buffer *Buffer, unichar rune, glyph fonts.GID) {
	buffer.cur(0).Glyph = glyph
	buffer.outputRune(unichar) // this is very confusing indeed.
	buffer.prev().setUnicodeProps(buffer)
}

func nextChar(buffer *Buffer, glyph fonts.GID) {
	buffer.cur(0).Glyph = glyph
	buffer.nextGlyph()
}

// returns 0 if didn't decompose, number of resulting characters otherwise.
func decompose(c *otNormalizeContext, shortest bool, ab rune) int {
	var aGlyph, bGlyph fonts.GID
	buffer := c.buffer
	font := c.font
	a, b, ok := c.decompose(c, ab)
	if !ok {
		return 0
	}
	bGlyph, ok = font.face.NominalGlyph(b)
	if b != 0 && !ok {
		return 0
	}

	aGlyph, hasA := font.face.NominalGlyph(a)
	if shortest && hasA {
		/// output a and b
		outputChar(buffer, a, aGlyph)
		if b != 0 {
			outputChar(buffer, b, bGlyph)
			return 2
		}
		return 1
	}

	if ret := decompose(c, shortest, a); ret != 0 {
		if b != 0 {
			outputChar(buffer, b, bGlyph)
			return ret + 1
		}
		return ret
	}

	if hasA {
		outputChar(buffer, a, aGlyph)
		if b != 0 {
			outputChar(buffer, b, bGlyph)
			return 2
		}
		return 1
	}

	return 0
}

func (c *otNormalizeContext) decomposeCurrentCharacter(shortest bool) {
	buffer := c.buffer
	u := buffer.cur(0).codepoint
	glyph, ok := c.font.nominalGlyph(u, c.buffer.NotFound)

	if shortest && ok {
		nextChar(buffer, glyph)
		return
	}

	if decompose(c, shortest, u) != 0 {
		buffer.skipGlyph()
		return
	}

	if !shortest && ok {
		nextChar(buffer, glyph)
		return
	}

	if buffer.cur(0).isUnicodeSpace() {
		spaceType := uni.spaceFallbackType(u)
		if spaceGlyph, ok := c.font.face.NominalGlyph(0x0020); spaceType != notSpace && ok {
			buffer.cur(0).setUnicodeSpaceFallbackType(spaceType)
			nextChar(buffer, spaceGlyph)
			buffer.scratchFlags |= bsfHasSpaceFallback
			return
		}
	}

	if u == 0x2011 {
		/* U+2011 is the only sensible character that is a no-break version of another character
		 * and not a space. The space ones are handled already.  Handle this lone one. */
		if otherGlyph, ok := c.font.face.NominalGlyph(0x2010); ok {
			nextChar(buffer, otherGlyph)
			return
		}
	}

	nextChar(buffer, glyph)
}

func (c *otNormalizeContext) handleVariationSelectorCluster(end int) {
	buffer := c.buffer
	if debugMode >= 1 {
		fmt.Printf("NORMALIZE - variation selector cluster at index %d\n", buffer.idx)
	}
	font := c.font
	for buffer.idx < end-1 {
		if uni.isVariationSelector(buffer.cur(+1).codepoint) {
			var ok bool
			buffer.cur(0).Glyph, ok = font.face.(FaceOpentype).VariationGlyph(buffer.cur(0).codepoint, buffer.cur(+1).codepoint)
			if ok {
				r := buffer.cur(0).codepoint
				buffer.replaceGlyphs(2, []rune{r}, nil)
			} else {
				// Just pass on the two characters separately, let GSUB do its magic.
				setGlyph(buffer.cur(0), font)
				buffer.nextGlyph()
				setGlyph(buffer.cur(0), font)
				buffer.nextGlyph()
			}
			// skip any further variation selectors.
			for buffer.idx < end && uni.isVariationSelector(buffer.cur(0).codepoint) {
				setGlyph(buffer.cur(0), font)
				buffer.nextGlyph()
			}
		} else {
			setGlyph(buffer.cur(0), font)
			buffer.nextGlyph()
		}
	}
	if buffer.idx < end {
		setGlyph(buffer.cur(0), font)
		buffer.nextGlyph()
	}
}

func (c *otNormalizeContext) decomposeMultiCharCluster(end int, shortCircuit bool) {
	buffer := c.buffer
	if debugMode >= 1 {
		fmt.Printf("NORMALIZE - decompose multi char cluster at index %d\n", buffer.idx)
	}

	for i := buffer.idx; i < end; i++ {
		if uni.isVariationSelector(buffer.Info[i].codepoint) {
			c.handleVariationSelectorCluster(end)
			return
		}
	}
	for buffer.idx < end {
		c.decomposeCurrentCharacter(shortCircuit)
	}
}

func compareCombiningClass(pa, pb *GlyphInfo) int {
	a := pa.getModifiedCombiningClass()
	b := pb.getModifiedCombiningClass()
	if a < b {
		return -1
	} else if a == b {
		return 0
	}
	return 1
}

func otShapeNormalize(plan *otShapePlan, buffer *Buffer, font *Font) {
	if len(buffer.Info) == 0 {
		return
	}

	mode := plan.shaper.normalizationPreference()
	if mode == nmAuto {
		if plan.hasGposMark {
			// https://github.com/harfbuzz/harfbuzz/issues/653#issuecomment-423905920
			mode = nmComposedDiacritics
		} else {
			mode = nmComposedDiacritics
		}
	}
	c := otNormalizeContext{
		plan,
		buffer,
		font,
		plan.shaper.decompose,
		plan.shaper.compose,
	}

	alwaysShortCircuit := mode == nmNone
	mightShortCircuit := alwaysShortCircuit ||
		(mode != nmDecomposed &&
			mode != nmComposedDiacriticsNoShortCircuit)

	/* We do a fairly straightforward yet custom normalization process in three
	* separate rounds: decompose, reorder, recompose (if desired). Currently
	* this makes two buffer swaps.  We can make it faster by moving the last
	* two rounds into the inner loop for the first round, but it's more readable
	* this way. */

	/* First round, decompose */

	allSimple := true
	buffer.clearOutput()
	count := len(buffer.Info)
	buffer.idx = 0
	var end int
	for do := true; do; do = buffer.idx < count {
		for end = buffer.idx + 1; end < count; end++ {
			if buffer.Info[end].isUnicodeMark() {
				break
			}
		}

		if end < count {
			end-- // leave one base for the marks to cluster with.
		}
		// from idx to end are simple clusters.
		if mightShortCircuit {
			var (
				i  int
				ok bool
			)
			for i = buffer.idx; i < end; i++ {
				buffer.Info[i].Glyph, ok = font.face.NominalGlyph(buffer.Info[i].codepoint)
				if !ok {
					break
				}
			}
			buffer.nextGlyphs(i - buffer.idx)
		}
		for buffer.idx < end {
			c.decomposeCurrentCharacter(mightShortCircuit)
		}

		if buffer.idx == count {
			break
		}

		allSimple = false

		// find all the marks now.
		for end = buffer.idx + 1; end < count; end++ {
			if !buffer.Info[end].isUnicodeMark() {
				break
			}
		}

		// idx to end is one non-simple cluster.
		c.decomposeMultiCharCluster(end, alwaysShortCircuit)
	}

	buffer.swapBuffers()
	/* Second round, reorder (inplace) */

	if !allSimple {
		if debugMode >= 1 {
			fmt.Println("NORMALIZE - start reorder")
		}
		count = len(buffer.Info)
		for i := 0; i < count; i++ {
			if buffer.Info[i].getModifiedCombiningClass() == 0 {
				continue
			}

			var end int
			for end = i + 1; end < count; end++ {
				if buffer.Info[end].getModifiedCombiningClass() == 0 {
					break
				}
			}

			// we are going to do a O(n^2).  Only do this if the sequence is short.
			if end-i > shapeComplexMaxCombiningMarks {
				i = end
				continue
			}

			buffer.sort(i, end, compareCombiningClass)

			plan.shaper.reorderMarks(plan, buffer, i, end)

			i = end
		}
		if debugMode >= 1 {
			fmt.Println("NORMALIZE - end reorder")
		}
	}

	if buffer.scratchFlags&bsfHasCGJ != 0 {
		/* For all CGJ, check if it prevented any reordering at all.
		 * If it did NOT, then make it skippable.
		 * https://github.com/harfbuzz/harfbuzz/issues/554 */
		for i := 1; i+1 < len(buffer.Info); i++ {
			if buffer.Info[i].codepoint == 0x034F /*CGJ*/ &&
				(buffer.Info[i+1].getModifiedCombiningClass() == 0 || buffer.Info[i-1].getModifiedCombiningClass() <= buffer.Info[i+1].getModifiedCombiningClass()) {
				buffer.Info[i].unhide()
			}
		}
	}

	/* Third round, recompose */

	if !allSimple &&
		(mode == nmComposedDiacritics ||
			mode == nmComposedDiacriticsNoShortCircuit) {

		if debugMode >= 1 {
			fmt.Println("NORMALIZE - recompose")
		}

		/* As noted in the comment earlier, we don't try to combine
		 * ccc=0 chars with their previous Starter. */

		buffer.clearOutput()
		count = len(buffer.Info)
		starter := 0
		buffer.nextGlyph()
		for buffer.idx < count {
			/* We don't try to compose a non-mark character with it's preceding starter.
			* This is both an optimization to avoid trying to compose every two neighboring
			* glyphs in most scripts AND a desired feature for Hangul.  Apparently Hangul
			* fonts are not designed to mix-and-match pre-composed syllables and Jamo. */
			if buffer.cur(0).isUnicodeMark() {
				/* If there's anything between the starter and this char, they should have CCC
				* smaller than this character's. */
				if starter == len(buffer.outInfo)-1 ||
					buffer.prev().getModifiedCombiningClass() < buffer.cur(0).getModifiedCombiningClass() {
					/* And compose. */
					composed, ok := c.compose(&c, buffer.outInfo[starter].codepoint, buffer.cur(0).codepoint)
					if ok { // And the font has glyph for the composite.
						glyph, ok := font.face.NominalGlyph(composed) /* Composes. */
						if ok {
							buffer.nextGlyph() /* Copy to out-buffer. */
							buffer.mergeOutClusters(starter, len(buffer.outInfo))
							buffer.outInfo = buffer.outInfo[:len(buffer.outInfo)-1] // remove the second composable.
							/* Modify starter and carry on. */
							buffer.outInfo[starter].codepoint = composed
							buffer.outInfo[starter].Glyph = glyph
							buffer.outInfo[starter].setUnicodeProps(buffer)
							continue
						}
					}
				}
			}

			/* Blocked, or doesn't compose. */
			buffer.nextGlyph()

			if buffer.prev().getModifiedCombiningClass() == 0 {
				starter = len(buffer.outInfo) - 1
			}
		}
		buffer.swapBuffers()
	}
}
