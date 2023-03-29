package harfbuzz

import (
	"fmt"

	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-shape-complex-khmer.cc Copyright © 2011,2012  Google, Inc. Behdad Esfahbod

var _ otComplexShaper = (*complexShaperKhmer)(nil)

// Khmer shaper
type complexShaperKhmer struct {
	plan khmerShapePlan
}

var khmerFeatures = [...]otMapFeature{
	/*
	* Basic features.
	* These features are applied in order, one at a time, after reordering.
	 */
	{loader.NewTag('p', 'r', 'e', 'f'), ffManualJoiners},
	{loader.NewTag('b', 'l', 'w', 'f'), ffManualJoiners},
	{loader.NewTag('a', 'b', 'v', 'f'), ffManualJoiners},
	{loader.NewTag('p', 's', 't', 'f'), ffManualJoiners},
	{loader.NewTag('c', 'f', 'a', 'r'), ffManualJoiners},
	/*
	* Other features.
	* These features are applied all at once after clearing syllables.
	 */
	{loader.NewTag('p', 'r', 'e', 's'), ffGlobalManualJoiners},
	{loader.NewTag('a', 'b', 'v', 's'), ffGlobalManualJoiners},
	{loader.NewTag('b', 'l', 'w', 's'), ffGlobalManualJoiners},
	{loader.NewTag('p', 's', 't', 's'), ffGlobalManualJoiners},
}

// Must be in the same order as the khmerFeatures array.
const (
	khmerPref = iota
	khmerBlwf
	khmerAbvf
	khmerPstf
	khmerCfar

	khmerPres
	khmerAbvs
	khmerBlws
	khmerPsts

	khmerNumFeatures
	khmerBasicFeatures = khmerPres /* Don't forget to update this! */
)

func (cs *complexShaperKhmer) collectFeatures(plan *otShapePlanner) {
	map_ := &plan.map_

	/* Do this before any lookups have been applied. */
	map_.addGSUBPause(setupSyllablesKhmer)
	map_.addGSUBPause(cs.reorderKhmer)

	/* Testing suggests that Uniscribe does NOT pause between basic
	* features.  Test with KhmerUI.ttf and the following three
	* sequences:
	*
	*   U+1789,U+17BC
	*   U+1789,U+17D2,U+1789
	*   U+1789,U+17D2,U+1789,U+17BC
	*
	* https://github.com/harfbuzz/harfbuzz/issues/974
	 */
	map_.enableFeature(loader.NewTag('l', 'o', 'c', 'l'))
	map_.enableFeature(loader.NewTag('c', 'c', 'm', 'p'))

	i := 0
	for ; i < khmerBasicFeatures; i++ {
		map_.addFeatureExt(khmerFeatures[i].tag, khmerFeatures[i].flags, 1)
	}

	map_.addGSUBPause(clearSyllables)

	for ; i < khmerNumFeatures; i++ {
		map_.addFeatureExt(khmerFeatures[i].tag, khmerFeatures[i].flags, 1)
	}
}

func (complexShaperKhmer) overrideFeatures(plan *otShapePlanner) {
	map_ := &plan.map_

	/* Khmer spec has 'clig' as part of required shaping features:
	* "Apply feature 'clig' to form ligatures that are desired for
	* typographical correctness.", hence in overrides... */
	map_.enableFeature(loader.NewTag('c', 'l', 'i', 'g'))

	/* Uniscribe does not apply 'kern' in Khmer. */
	if UniscribeBugCompatible {
		map_.disableFeature(loader.NewTag('k', 'e', 'r', 'n'))
	}

	map_.disableFeature(loader.NewTag('l', 'i', 'g', 'a'))
}

type khmerShapePlan struct {
	viramaGlyph GID
	maskArray   [khmerNumFeatures]GlyphMask
}

func (cs *complexShaperKhmer) dataCreate(plan *otShapePlan) {
	var khmerPlan khmerShapePlan

	khmerPlan.viramaGlyph = ^GID(0)

	for i := range khmerPlan.maskArray {
		if khmerFeatures[i].flags&ffGLOBAL == 0 {
			khmerPlan.maskArray[i] = plan.map_.getMask1(khmerFeatures[i].tag)
		}
	}

	cs.plan = khmerPlan
}

func (cs *complexShaperKhmer) setupMasks(_ *otShapePlan, buffer *Buffer, _ *Font) {
	/* We cannot setup masks here.  We save information about characters
	* and setup masks later on in a pause-callback. */

	info := buffer.Info
	for i := range info {
		setKhmerProperties(&info[i])
	}
}

/* Note: This enum is duplicated in the -machine.rl source file.
 * Not sure how to avoid duplication. */
const (
	otRobatic = 20
	otXgroup  = 21
	otYgroup  = 22
)

func setKhmerProperties(info *GlyphInfo) {
	u := info.codepoint
	type_ := indicGetCategories(u)
	cat := uint8(type_ & 0xFF)
	pos := uint8(type_ >> 8)

	/*
	 * Re-assign category
	 *
	 * These categories are experimentally extracted from what Uniscribe allows.
	 */
	switch u {
	case 0x179A:
		cat = otRa

	case 0x17CC, 0x17C9, 0x17CA:
		cat = otRobatic

	case 0x17C6, 0x17CB, 0x17CD, 0x17CE, 0x17CF, 0x17D0, 0x17D1:
		cat = otXgroup

	case 0x17C7, 0x17C8, 0x17DD, 0x17D3: /* Just guessing. Uniscribe doesn't categorize it. */
		cat = otYgroup
	}

	/*
	 * Re-assign position.
	 */
	if cat == otM {
		switch pos {
		case posPreC:
			cat = otVPre
		case posBelowC:
			cat = otVBlw
		case posAboveC:
			cat = otVAbv
		case posPostC:
			cat = otVPst
		}
	}

	info.complexCategory = cat
}

func setupSyllablesKhmer(_ *otShapePlan, _ *Font, buffer *Buffer) {
	findSyllablesKhmer(buffer)
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		buffer.unsafeToBreak(start, end)
	}
}

func foundSyllableKhmer(syllableType uint8, ts, te int, info []GlyphInfo, syllableSerial *uint8) {
	for i := ts; i < te; i++ {
		info[i].syllable = (*syllableSerial << 4) | syllableType
	}
	*syllableSerial++
	if *syllableSerial == 16 {
		*syllableSerial = 1
	}
}

/* Rules from:
 * https://docs.microsoft.com/en-us/typography/script-development/devanagari */
func (khmerPlan *khmerShapePlan) reorderConsonantSyllable(buffer *Buffer, start, end int) {
	info := buffer.Info

	/* Setup masks. */
	{
		/* Post-base */
		mask := khmerPlan.maskArray[khmerBlwf] |
			khmerPlan.maskArray[khmerAbvf] |
			khmerPlan.maskArray[khmerPstf]
		for i := start + 1; i < end; i++ {
			info[i].Mask |= mask
		}
	}

	numCoengs := 0
	for i := start + 1; i < end; i++ {
		/* """
		 * When a COENG + (Cons | IndV) combination are found (and subscript count
		 * is less than two) the character combination is handled according to the
		 * subscript type of the character following the COENG.
		 *
		 * ...
		 *
		 * Subscript Type 2 - The COENG + RO characters are reordered to immediately
		 * before the base glyph. Then the COENG + RO characters are assigned to have
		 * the 'pref' OpenType feature applied to them.
		 * """
		 */
		if info[i].complexCategory == otCoeng && numCoengs <= 2 && i+1 < end {
			numCoengs++

			if info[i+1].complexCategory == otRa {
				for j := 0; j < 2; j++ {
					info[i+j].Mask |= khmerPlan.maskArray[khmerPref]
				}

				/* Move the Coeng,Ro sequence to the start. */
				buffer.mergeClusters(start, i+2)
				t0 := info[i]
				t1 := info[i+1]
				copy(info[start+2:], info[start:i])
				info[start] = t0
				info[start+1] = t1

				/* Mark the subsequent stuff with 'cfar'.  Used in Khmer.
				 * Read the feature spec.
				 * This allows distinguishing the following cases with MS Khmer fonts:
				 * U+1784,U+17D2,U+179A,U+17D2,U+1782
				 * U+1784,U+17D2,U+1782,U+17D2,U+179A
				 */
				if khmerPlan.maskArray[khmerCfar] != 0 {
					for j := i + 2; j < end; j++ {
						info[j].Mask |= khmerPlan.maskArray[khmerCfar]
					}
				}

				numCoengs = 2 /* Done. */
			}
		} else if info[i].complexCategory == otVPre { /* Reorder left matra piece. */
			/* Move to the start. */
			buffer.mergeClusters(start, i+1)
			t := info[i]
			copy(info[start+1:], info[start:i])
			info[start] = t
		}
	}
}

func (cs *complexShaperKhmer) reorderSyllableKhmer(buffer *Buffer, start, end int) {
	syllableType := buffer.Info[start].syllable & 0x0F
	switch syllableType {
	case khmerBrokenCluster, /* We already inserted dotted-circles, so just call the consonant_syllable. */
		khmerConsonantSyllable:
		cs.plan.reorderConsonantSyllable(buffer, start, end)
	}
}

func (cs *complexShaperKhmer) reorderKhmer(_ *otShapePlan, font *Font, buffer *Buffer) {
	if debugMode >= 1 {
		fmt.Println("KHMER - start reordering khmer")
	}

	syllabicInsertDottedCircles(font, buffer, khmerBrokenCluster, otDOTTEDCIRCLE, otRepha, -1)
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		cs.reorderSyllableKhmer(buffer, start, end)
	}

	if debugMode >= 1 {
		fmt.Println("KHMER - end reordering khmer")
	}
}

func (complexShaperKhmer) decompose(c *otNormalizeContext, ab rune) (rune, rune, bool) {
	switch ab {
	/*
	 * Decompose split matras that don't have Unicode decompositions.
	 */

	/* Khmer */
	case 0x17BE:
		return 0x17C1, 0x17BE, true
	case 0x17BF:
		return 0x17C1, 0x17BF, true
	case 0x17C0:
		return 0x17C1, 0x17C0, true
	case 0x17C4:
		return 0x17C1, 0x17C4, true
	case 0x17C5:
		return 0x17C1, 0x17C5, true
	}

	return uni.decompose(ab)
}

func (complexShaperKhmer) compose(_ *otNormalizeContext, a, b rune) (rune, bool) {
	/* Avoid recomposing split matras. */
	if uni.generalCategory(a).isMark() {
		return 0, false
	}

	return uni.compose(a, b)
}

func (complexShaperKhmer) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksNone, false
}

func (complexShaperKhmer) normalizationPreference() normalizationMode {
	return nmComposedDiacriticsNoShortCircuit
}

func (complexShaperKhmer) gposTag() tables.Tag                         { return 0 }
func (complexShaperKhmer) preprocessText(*otShapePlan, *Buffer, *Font) {}
func (complexShaperKhmer) postprocessGlyphs(*otShapePlan, *Buffer, *Font) {
}
func (complexShaperKhmer) reorderMarks(*otShapePlan, *Buffer, int, int) {}
