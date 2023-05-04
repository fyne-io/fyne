package harfbuzz

import (
	"fmt"
	"sort"

	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
	ucd "github.com/go-text/typesetting/unicodedata"
)

// ported from harfbuzz/src/hb-ot-shape-complex-arabic.cc, hb-ot-shape-complex-arabic-fallback.hh Copyright © 2010,2012  Google, Inc. Behdad Esfahbod

var _ otComplexShaper = (*complexShaperArabic)(nil)

const flagArabicHasStch = bsfComplex0

/* See:
 * https://github.com/harfbuzz/harfbuzz/commit/6e6f82b6f3dde0fc6c3c7d991d9ec6cfff57823d#commitcomment-14248516 */
func isWord(genCat generalCategory) bool {
	const mask = 1<<unassigned |
		1<<privateUse |
		/*1 <<  LowercaseLetter |*/
		1<<modifierLetter |
		1<<otherLetter |
		/*1 <<  TitlecaseLetter |*/
		/*1 <<  UppercaseLetter |*/
		1<<spacingMark |
		1<<enclosingMark |
		1<<nonSpacingMark |
		1<<decimalNumber |
		1<<letterNumber |
		1<<otherNumber |
		1<<currencySymbol |
		1<<modifierSymbol |
		1<<mathSymbol |
		1<<otherSymbol
	return (1<<genCat)&mask != 0
}

/*
 * Joining types:
 */

// index into arabicStateTable
const (
	joiningTypeU = iota
	joiningTypeL
	joiningTypeR
	joiningTypeD
	joiningGroupAlaph
	joiningGroupDalathRish
	numStateMachineCols
	joiningTypeT
	joiningTypeC = joiningTypeD
)

func getJoiningType(u rune, genCat generalCategory) uint8 {
	if jType, ok := ucd.ArabicJoinings[u]; ok {
		switch jType {
		case ucd.U:
			return joiningTypeU
		case ucd.L:
			return joiningTypeL
		case ucd.R:
			return joiningTypeR
		case ucd.D:
			return joiningTypeD
		case ucd.Alaph:
			return joiningGroupAlaph
		case ucd.DalathRish:
			return joiningGroupDalathRish
		case ucd.T:
			return joiningTypeT
		case ucd.C:
			return joiningTypeC
		}
	}

	const mask = 1<<nonSpacingMark | 1<<enclosingMark | 1<<format
	if 1<<genCat&mask != 0 {
		return joiningTypeT
	}
	return joiningTypeU
}

func featureIsSyriac(tag loader.Tag) bool {
	return '2' <= byte(tag) && byte(tag) <= '3'
}

var arabicFeatures = [...]loader.Tag{
	loader.NewTag('i', 's', 'o', 'l'),
	loader.NewTag('f', 'i', 'n', 'a'),
	loader.NewTag('f', 'i', 'n', '2'),
	loader.NewTag('f', 'i', 'n', '3'),
	loader.NewTag('m', 'e', 'd', 'i'),
	loader.NewTag('m', 'e', 'd', '2'),
	loader.NewTag('i', 'n', 'i', 't'),
	0,
}

/* Same order as the feature array */
const (
	arabIsol = iota
	arabFina
	arabFin2
	araFin3
	arabMedi
	arabMed2
	arabInit

	arabNone

	/* We abuse the same byte for other things... */
	arabStchFixed
	arabStchRepeating
)

var arabicStateTable = [...][numStateMachineCols]struct {
	prevAction uint8
	currAction uint8
	nextState  uint16
}{
	/*   jt_U,          jt_L,          jt_R,          jt_D,          jg_ALAPH,      jg_DALATH_RISH */

	/* State 0: prev was U, not willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabNone, arabIsol, 1}, {arabNone, arabIsol, 2}, {arabNone, arabIsol, 1}, {arabNone, arabIsol, 6}},

	/* State 1: prev was R or ISOL/ALAPH, not willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabNone, arabIsol, 1}, {arabNone, arabIsol, 2}, {arabNone, arabFin2, 5}, {arabNone, arabIsol, 6}},

	/* State 2: prev was D/L in ISOL form, willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabInit, arabFina, 1}, {arabInit, arabFina, 3}, {arabInit, arabFina, 4}, {arabInit, arabFina, 6}},

	/* State 3: prev was D in FINA form, willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabMedi, arabFina, 1}, {arabMedi, arabFina, 3}, {arabMedi, arabFina, 4}, {arabMedi, arabFina, 6}},

	/* State 4: prev was FINA ALAPH, not willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabMed2, arabIsol, 1}, {arabMed2, arabIsol, 2}, {arabMed2, arabFin2, 5}, {arabMed2, arabIsol, 6}},

	/* State 5: prev was FIN2/FIN3 ALAPH, not willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabIsol, arabIsol, 1}, {arabIsol, arabIsol, 2}, {arabIsol, arabFin2, 5}, {arabIsol, arabIsol, 6}},

	/* State 6: prev was DALATH/RISH, not willing to join. */
	{{arabNone, arabNone, 0}, {arabNone, arabIsol, 2}, {arabNone, arabIsol, 1}, {arabNone, arabIsol, 2}, {arabNone, araFin3, 5}, {arabNone, arabIsol, 6}},
}

type complexShaperArabic struct {
	complexShaperNil

	plan arabicShapePlan
}

func (complexShaperArabic) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksByGdefLate, true
}

func (complexShaperArabic) normalizationPreference() normalizationMode {
	return nmDefault
}

func (cs *complexShaperArabic) collectFeatures(plan *otShapePlanner) {
	map_ := &plan.map_

	/* We apply features according to the Arabic spec, with pauses
	* in between most.
	*
	* The pause between init/medi/... and rlig is required.  See eg:
	* https://bugzilla.mozilla.org/show_bug.cgi?id=644184
	*
	* The pauses between init/medi/... themselves are not necessarily
	* needed as only one of those features is applied to any character.
	* The only difference it makes is when fonts have contextual
	* substitutions.  We now follow the order of the spec, which makes
	* for better experience if that's what Uniscribe is doing.
	*
	* At least for Arabic, looks like Uniscribe has a pause between
	* rlig and calt.  Otherwise the IranNastaliq's ALLAH ligature won't
	* work.  However, testing shows that rlig and calt are applied
	* together for Mongolian in Uniscribe.  As such, we only add a
	* pause for Arabic, not other scripts.
	*
	* A pause after calt is required to make KFGQPC Uthmanic Script HAFS
	* work correctly.  See https://github.com/harfbuzz/harfbuzz/issues/505
	 */

	map_.enableFeature(loader.NewTag('s', 't', 'c', 'h'))
	map_.addGSUBPause(recordStch)

	map_.enableFeature(loader.NewTag('c', 'c', 'm', 'p'))
	map_.enableFeature(loader.NewTag('l', 'o', 'c', 'l'))

	map_.addGSUBPause(nil)

	for _, arabFeat := range arabicFeatures {
		hasFallback := plan.props.Script == language.Arabic && !featureIsSyriac(arabFeat)
		fl := ffNone
		if hasFallback {
			fl = ffHasFallback
		}
		map_.addFeatureExt(arabFeat, fl, 1)
		map_.addGSUBPause(nil)
	}

	/* Unicode says a ZWNJ means "don't ligate". In Arabic script
	* however, it says a ZWJ should also mean "don't ligate". So we run
	* the main ligating features as MANUAL_ZWJ. */

	map_.enableFeatureExt(loader.NewTag('r', 'l', 'i', 'g'), ffManualZWJ|ffHasFallback, 1)

	if plan.props.Script == language.Arabic {
		map_.addGSUBPause(arabicFallbackShape)
	}
	/* No pause after rclt.  See 98460779bae19e4d64d29461ff154b3527bf8420. */
	map_.enableFeatureExt(loader.NewTag('r', 'c', 'l', 't'), ffManualZWJ, 1)
	map_.enableFeatureExt(loader.NewTag('c', 'a', 'l', 't'), ffManualZWJ, 1)
	map_.addGSUBPause(nil)

	/* The spec includes 'cswh'.  Earlier versions of Windows
	* used to enable this by default, but testing suggests
	* that Windows 8 and later do not enable it by default,
	* and spec now says 'Off by default'.
	* We disabled this in ae23c24c32.
	* Note that IranNastaliq uses this feature extensively
	* to fixup broken glyph sequences.  Oh well...
	* Test case: U+0643,U+0640,U+0631. */
	//map_.enable_feature (newTag('c','s','w','h'));
	map_.enableFeature(loader.NewTag('m', 's', 'e', 't'))
}

type arabicShapePlan struct {
	fallbackPlan *arabicFallbackPlan
	/* The "+ 1" in the next array is to accommodate for the "NONE" command,
	* which is not an OpenType feature, but this simplifies the code by not
	* having to do a "if (... < NONE) ..." and just rely on the fact that
	* maskArray[NONE] == 0. */
	maskArray  [len(arabicFeatures) + 1]GlyphMask
	doFallback bool
	hasStch    bool
}

func newArabicPlan(plan *otShapePlan) arabicShapePlan {
	var arabicPlan arabicShapePlan

	arabicPlan.doFallback = plan.props.Script == language.Arabic
	arabicPlan.hasStch = plan.map_.getMask1(loader.NewTag('s', 't', 'c', 'h')) != 0
	for i, arabFeat := range arabicFeatures {
		arabicPlan.maskArray[i] = plan.map_.getMask1(arabFeat)
		arabicPlan.doFallback = arabicPlan.doFallback &&
			(featureIsSyriac(arabFeat) || plan.map_.needsFallback(arabFeat))
	}
	return arabicPlan
}

func (cs *complexShaperArabic) dataCreate(plan *otShapePlan) {
	cs.plan = newArabicPlan(plan)
}

func arabicJoining(buffer *Buffer) {
	info := buffer.Info
	prev, state := -1, uint16(0)

	// check pre-context
	for _, u := range buffer.context[0] {
		thisType := getJoiningType(u, uni.generalCategory(u))

		if thisType == joiningTypeT {
			continue
		}

		entry := &arabicStateTable[state][thisType]
		state = entry.nextState
		break
	}

	for i := 0; i < len(info); i++ {
		thisType := getJoiningType(info[i].codepoint, info[i].unicode.generalCategory())

		if thisType == joiningTypeT {
			info[i].complexAux = arabNone
			continue
		}

		entry := &arabicStateTable[state][thisType]

		if entry.prevAction != arabNone && prev != -1 {
			info[prev].complexAux = entry.prevAction
			buffer.unsafeToBreak(prev, i+1)
		}

		info[i].complexAux = entry.currAction

		prev = i
		state = entry.nextState
	}

	for _, u := range buffer.context[1] {
		thisType := getJoiningType(u, uni.generalCategory(u))

		if thisType == joiningTypeT {
			continue
		}

		entry := &arabicStateTable[state][thisType]
		if entry.prevAction != arabNone && prev != -1 {
			info[prev].complexAux = entry.prevAction
		}
		break
	}
}

func mongolianVariationSelectors(buffer *Buffer) {
	// copy complexAux from base to Mongolian variation selectors.
	info := buffer.Info
	for i := 1; i < len(info); i++ {
		if cp := info[i].codepoint; 0x180B <= cp && cp <= 0x180D || cp == 0x180F {
			info[i].complexAux = info[i-1].complexAux
		}
	}
}

func (arabicPlan arabicShapePlan) setupMasks(buffer *Buffer, script language.Script) {
	arabicJoining(buffer)
	if script == language.Mongolian {
		mongolianVariationSelectors(buffer)
	}

	info := buffer.Info
	for i := range info {
		info[i].Mask |= arabicPlan.maskArray[info[i].complexAux]
	}
}

func (cs *complexShaperArabic) setupMasks(plan *otShapePlan, buffer *Buffer, _ *Font) {
	cs.plan.setupMasks(buffer, plan.props.Script)
}

func arabicFallbackShape(plan *otShapePlan, font *Font, buffer *Buffer) {
	arabicPlan := plan.shaper.(*complexShaperArabic).plan

	if !arabicPlan.doFallback {
		return
	}

	fallbackPlan := arabicPlan.fallbackPlan
	if fallbackPlan == nil {
		// this sucks. We need a font to build the fallback plan...
		fallbackPlan = newArabicFallbackPlan(plan, font)
	}

	fallbackPlan.shape(font, buffer)
}

//  /*
//   * Stretch feature: "stch".
//   * See example here:
//   * https://docs.microsoft.com/en-us/typography/script-development/syriac
//   * We implement this in a generic way, such that the Arabic subtending
//   * marks can use it as well.
//   */

func recordStch(plan *otShapePlan, _ *Font, buffer *Buffer) {
	arabicPlan := plan.shaper.(*complexShaperArabic).plan
	if !arabicPlan.hasStch {
		return
	}

	/* 'stch' feature was just applied.  Look for anything that multiplied,
	* and record it for stch treatment later.  Note that rtlm, frac, etc
	* are applied before stch, but we assume that they didn't result in
	* anything multiplying into 5 pieces, so it's safe-ish... */

	info := buffer.Info
	for i := range info {
		if info[i].multiplied() {
			comp := info[i].getLigComp()
			if comp%2 != 0 {
				info[i].complexAux = arabStchRepeating
			} else {
				info[i].complexAux = arabStchFixed
			}
			buffer.scratchFlags |= flagArabicHasStch
		}
	}
}

func inRange(sa uint8) bool {
	return arabStchFixed <= sa && sa <= arabStchRepeating
}

func (cs *complexShaperArabic) postprocessGlyphs(plan *otShapePlan, buffer *Buffer, font *Font) {
	if buffer.scratchFlags&flagArabicHasStch == 0 {
		return
	}

	/* The Arabic shaper currently always processes in RTL mode, so we should
	* stretch / position the stretched pieces to the left / preceding glyphs. */

	/* We do a two pass implementation:
	* First pass calculates the exact number of extra glyphs we need,
	* We then enlarge buffer to have that much room,
	* Second pass applies the stretch, copying things to the end of buffer. */

	sign := Position(+1)
	if font.XScale < 0 {
		sign = -1
	}
	const (
		MEASURE = iota
		CUT
	)
	var (
		originCount       = len(buffer.Info) // before enlarging
		extraGlyphsNeeded = 0                // Set during MEASURE, used during CUT
	)
	for step := MEASURE; step <= CUT; step++ {
		info := buffer.Info
		pos := buffer.Pos
		j := len(info) // enlarged after MEASURE
		for i := originCount; i != 0; i-- {
			if sa := info[i-1].complexAux; !inRange(sa) {
				if step == CUT {
					j--
					info[j] = info[i-1]
					pos[j] = pos[i-1]
				}
				continue
			}

			/* Yay, justification! */
			var (
				wTotal     Position // Total to be filled
				wFixed     Position // Sum of fixed tiles
				wRepeating Position // Sum of repeating tiles
				nFixed     = 0
				nRepeating = 0
			)
			end := i
			for i != 0 && inRange(info[i-1].complexAux) {
				i--
				width := font.GlyphHAdvance(info[i].Glyph)
				if info[i].complexAux == arabStchFixed {
					wFixed += width
					nFixed++
				} else {
					wRepeating += width
					nRepeating++
				}
			}
			start := i
			context := i
			for context != 0 && !inRange(info[context-1].complexAux) &&
				((&info[context-1]).isDefaultIgnorable() ||
					isWord((&info[context-1]).unicode.generalCategory())) {
				context--
				wTotal += pos[context].XAdvance
			}
			i++ // Don't touch i again.

			if debugMode >= 1 {
				fmt.Printf("ARABIC - step %d: stretch at (%d,%d,%d)\n", step+1, context, start, end)
				fmt.Printf("ARABIC - rest of word:    count=%d width %d\n", start-context, wTotal)
				fmt.Printf("ARABIC - fixed tiles:     count=%d width=%d\n", nFixed, wFixed)
				fmt.Printf("ARABIC - repeating tiles: count=%d width=%d\n", nRepeating, wRepeating)
			}

			// number of additional times to repeat each repeating tile.
			var nCopies int

			wRemaining := wTotal - wFixed
			if sign*wRemaining > sign*wRepeating && sign*wRepeating > 0 {
				nCopies = int((sign*wRemaining)/(sign*wRepeating) - 1)
			}

			// see if we can improve the fit by adding an extra repeat and squeezing them together a bit.
			var extraRepeatOverlap Position
			shortfall := sign*wRemaining - sign*wRepeating*(Position(nCopies)+1)
			if shortfall > 0 && nRepeating > 0 {
				nCopies++
				excess := (Position(nCopies)+1)*sign*wRepeating - sign*wRemaining
				if excess > 0 {
					extraRepeatOverlap = excess / Position(nCopies*nRepeating)
				}
			}

			if step == MEASURE {
				extraGlyphsNeeded += nCopies * nRepeating
				if debugMode >= 1 {
					fmt.Printf("ARABIC - will add extra %d copies of repeating tiles\n", nCopies)
				}
			} else {
				buffer.unsafeToBreak(context, end)
				var xOffset Position
				for k := end; k > start; k-- {
					width := font.GlyphHAdvance(info[k-1].Glyph)

					repeat := 1
					if info[k-1].complexAux == arabStchRepeating {
						repeat += nCopies
					}

					if debugMode >= 1 {
						fmt.Printf("ARABIC - appending %d copies of glyph %d; j=%d\n", repeat, info[k-1].codepoint, j)
					}
					for n := 0; n < repeat; n++ {
						xOffset -= width
						if n > 0 {
							xOffset += extraRepeatOverlap
						}
						pos[k-1].XOffset = xOffset
						// append copy.
						j--
						info[j] = info[k-1]
						pos[j] = pos[k-1]
					}
				}
			}
		}

		if step == MEASURE { // enlarge
			buffer.Info = append(buffer.Info, make([]GlyphInfo, extraGlyphsNeeded)...)
			buffer.Pos = append(buffer.Pos, make([]GlyphPosition, extraGlyphsNeeded)...)
		}
	}
}

// https://www.unicode.org/reports/tr53/
var modifierCombiningMarks = [...]rune{
	0x0654, /* ARABIC HAMZA ABOVE */
	0x0655, /* ARABIC HAMZA BELOW */
	0x0658, /* ARABIC MARK NOON GHUNNA */
	0x06DC, /* ARABIC SMALL HIGH SEEN */
	0x06E3, /* ARABIC SMALL LOW SEEN */
	0x06E7, /* ARABIC SMALL HIGH YEH */
	0x06E8, /* ARABIC SMALL HIGH NOON */
	0x08D3, /* ARABIC SMALL LOW WAW */
	0x08F3, /* ARABIC SMALL HIGH WAW */
}

func infoIsMcm(info *GlyphInfo) bool {
	u := info.codepoint
	for i := 0; i < len(modifierCombiningMarks); i++ {
		if u == modifierCombiningMarks[i] {
			return true
		}
	}
	return false
}

func (cs *complexShaperArabic) reorderMarks(_ *otShapePlan, buffer *Buffer, start, end int) {
	info := buffer.Info

	if debugMode >= 1 {
		fmt.Printf("ARABIC - Reordering marks from %d to %d\n", start, end)
	}

	i := start
	for cc := uint8(220); cc <= 230; cc += 10 {
		if debugMode >= 1 {
			fmt.Printf("ARABIC - Looking for %d's starting at %d\n", cc, i)
		}
		for i < end && info[i].getModifiedCombiningClass() < cc {
			i++
		}
		if debugMode >= 1 {
			fmt.Printf("ARABIC - Looking for %d's stopped at %d\n", cc, i)
		}

		if i == end {
			break
		}

		if info[i].getModifiedCombiningClass() > cc {
			continue
		}

		j := i
		for j < end && info[j].getModifiedCombiningClass() == cc && infoIsMcm(&info[j]) {
			j++
		}

		if i == j {
			continue
		}

		if debugMode >= 1 {
			fmt.Printf("ARABIC - Found %d's from %d to %d", cc, i, j)
			// shift it!
			fmt.Printf("ARABIC - Shifting %d's: %d %d", cc, i, j)
		}

		var temp [shapeComplexMaxCombiningMarks]GlyphInfo
		//  assert (j - i <= len (temp));
		buffer.mergeClusters(start, j)
		copy(temp[:j-i], info[i:])
		copy(info[start+j-i:], info[start:i])
		copy(info[start:], temp[:j-i])

		/* Renumber CC such that the reordered sequence is still sorted.
		 * 22 and 26 are chosen because they are smaller than all Arabic categories,
		 * and are folded back to 220/230 respectively during fallback mark positioning.
		 *
		 * We do this because the CGJ-handling logic in the normalizer relies on
		 * mark sequences having an increasing order even after this reordering.
		 * https://github.com/harfbuzz/harfbuzz/issues/554
		 * This, however, does break some obscure sequences, where the normalizer
		 * might compose a sequence that it should not.  For example, in the seequence
		 * ALEF, HAMZAH, MADDAH, we should NOT try to compose ALEF+MADDAH, but with this
		 * renumbering, we will. */
		newStart := start + j - i
		newCc := mcc26
		if cc == 220 {
			newCc = mcc26
		}
		for start < newStart {
			info[start].setModifiedCombiningClass(newCc)
			start++
		}

		i = j
	}
}

/* Features ordered the same as the entries in ucd.ArabicShaping rows,
 * followed by rlig.  Don't change. */
var arabicFallbackFeatures = [...]loader.Tag{
	loader.NewTag('i', 's', 'o', 'l'),
	loader.NewTag('f', 'i', 'n', 'a'),
	loader.NewTag('i', 'n', 'i', 't'),
	loader.NewTag('m', 'e', 'd', 'i'),
	loader.NewTag('r', 'l', 'i', 'g'),
}

// used to sort both array at the same time
type jointGlyphs struct {
	glyphs, substitutes []gID
}

func (a jointGlyphs) Len() int { return len(a.glyphs) }
func (a jointGlyphs) Swap(i, j int) {
	a.glyphs[i], a.glyphs[j] = a.glyphs[j], a.glyphs[i]
	a.substitutes[i], a.substitutes[j] = a.substitutes[j], a.substitutes[i]
}
func (a jointGlyphs) Less(i, j int) bool { return a.glyphs[i] < a.glyphs[j] }

func arabicFallbackSynthesizeLookupSingle(ft *Font, featureIndex int) *lookupGSUB {
	var glyphs, substitutes []gID

	// populate arrays
	for u := rune(ucd.FirstArabicShape); u <= ucd.LastArabicShape; u++ {
		s := rune(ucd.ArabicShaping[u-ucd.FirstArabicShape][featureIndex])
		uGlyph, hasU := ft.face.NominalGlyph(u)
		sGlyph, hasS := ft.face.NominalGlyph(s)

		if s == 0 || !hasU || !hasS || uGlyph == sGlyph || uGlyph > 0xFFFF || sGlyph > 0xFFFF {
			continue
		}

		glyphs = append(glyphs, gID(uGlyph))
		substitutes = append(substitutes, gID(sGlyph))
	}

	if len(glyphs) == 0 {
		return nil
	}

	sort.Stable(jointGlyphs{glyphs: glyphs, substitutes: substitutes})

	return &lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables: []tables.GSUBLookup{
			tables.SingleSubs{Data: tables.SingleSubstData2{
				Coverage:           tables.Coverage1{Glyphs: glyphs},
				SubstituteGlyphIDs: substitutes,
			}},
		},
	}
}

// used to sort both array at the same time
type glyphsIndirections struct {
	glyphs       []gID
	indirections []int
}

func (a glyphsIndirections) Len() int { return len(a.glyphs) }
func (a glyphsIndirections) Swap(i, j int) {
	a.glyphs[i], a.glyphs[j] = a.glyphs[j], a.glyphs[i]
	a.indirections[i], a.indirections[j] = a.indirections[j], a.indirections[i]
}
func (a glyphsIndirections) Less(i, j int) bool { return a.glyphs[i] < a.glyphs[j] }

func arabicFallbackSynthesizeLookupLigature(ft *Font) *lookupGSUB {
	var (
		firstGlyphs            []gID
		firstGlyphsIndirection []int // original index into ArabicLigatures
	)

	/* Populate arrays */

	// sort out the first-glyphs
	for firstGlyphIdx, lig := range ucd.ArabicLigatures {
		firstGlyph, ok := ft.face.NominalGlyph(lig.First)
		if !ok {
			continue
		}
		firstGlyphs = append(firstGlyphs, gID(firstGlyph))
		firstGlyphsIndirection = append(firstGlyphsIndirection, firstGlyphIdx)
	}

	if len(firstGlyphs) == 0 {
		return nil
	}

	sort.Stable(glyphsIndirections{glyphs: firstGlyphs, indirections: firstGlyphsIndirection})

	var out tables.LigatureSubs
	out.Coverage = tables.Coverage1{Glyphs: firstGlyphs}

	// now that the first-glyphs are sorted, walk again, populate ligatures.
	for _, firstGlyphIdx := range firstGlyphsIndirection {
		ligs := ucd.ArabicLigatures[firstGlyphIdx].Ligatures
		var ligatureSet tables.LigatureSet
		for _, v := range ligs {
			secondU, ligatureU := v[0], v[1]
			secondGlyph, hasSecond := ft.face.NominalGlyph(secondU)
			ligatureGlyph, hasLigature := ft.face.NominalGlyph(ligatureU)
			if secondU == 0 || !hasSecond || !hasLigature {
				continue
			}
			ligatureSet.Ligatures = append(ligatureSet.Ligatures, tables.Ligature{
				LigatureGlyph:     gID(ligatureGlyph),
				ComponentGlyphIDs: []uint16{uint16(secondGlyph)}, // ligatures are 2-component
			})
		}
		out.LigatureSets = append(out.LigatureSets, ligatureSet)
	}

	return &lookupGSUB{
		LookupOptions: font.LookupOptions{Flag: otIgnoreMarks},
		Subtables: []tables.GSUBLookup{
			out,
		},
	}
}

func arabicFallbackSynthesizeLookup(font *Font, featureIndex int) *lookupGSUB {
	if featureIndex < 4 {
		return arabicFallbackSynthesizeLookupSingle(font, featureIndex)
	}
	return arabicFallbackSynthesizeLookupLigature(font)
}

const arabicFallbackMaxLookups = 5

type arabicFallbackPlan struct {
	accelArray [arabicFallbackMaxLookups]otLayoutLookupAccelerator
	numLookups int
	maskArray  [arabicFallbackMaxLookups]GlyphMask
}

func (fbPlan *arabicFallbackPlan) initWin1256(plan *otShapePlan, font *Font) bool {
	// does this font look like it's Windows-1256-encoded?
	g1, _ := font.face.NominalGlyph(0x0627) /* ALEF */
	g2, _ := font.face.NominalGlyph(0x0644) /* LAM */
	g3, _ := font.face.NominalGlyph(0x0649) /* ALEF MAKSURA */
	g4, _ := font.face.NominalGlyph(0x064A) /* YEH */
	g5, _ := font.face.NominalGlyph(0x0652) /* SUKUN */
	if !(g1 == 199 && g2 == 225 && g3 == 236 && g4 == 237 && g5 == 250) {
		return false
	}

	var j int
	for _, man := range arabicWin1256GsubLookups {
		fbPlan.maskArray[j] = plan.map_.getMask1(man.tag)
		if fbPlan.maskArray[j] != 0 {
			if man.lookup != nil {
				fbPlan.accelArray[j].init(*man.lookup)
				j++
			}
		}
	}

	fbPlan.numLookups = j

	return j > 0
}

func (fbPlan *arabicFallbackPlan) initUnicode(plan *otShapePlan, font *Font) bool {
	var j int
	for i, feat := range arabicFallbackFeatures {
		fbPlan.maskArray[j] = plan.map_.getMask1(feat)
		if fbPlan.maskArray[j] != 0 {
			lk := arabicFallbackSynthesizeLookup(font, i)
			if lk != nil {
				fbPlan.accelArray[j].init(*lk)
				j++
			}
		}
	}

	fbPlan.numLookups = j

	return j > 0
}

func newArabicFallbackPlan(plan *otShapePlan, font *Font) *arabicFallbackPlan {
	var fbPlan arabicFallbackPlan

	/* Try synthesizing GSUB table using Unicode Arabic Presentation Forms,
	* in case the font has cmap entries for the presentation-forms characters. */
	if fbPlan.initUnicode(plan, font) {
		return &fbPlan
	}

	/* See if this looks like a Windows-1256-encoded font. If it does, use a
	* hand-coded GSUB table. */
	if fbPlan.initWin1256(plan, font) {
		return &fbPlan
	}

	return &arabicFallbackPlan{}
}

func (fbPlan *arabicFallbackPlan) shape(font *Font, buffer *Buffer) {
	c := newOtApplyContext(0, font, buffer)
	for i := 0; i < fbPlan.numLookups; i++ {
		if fbPlan.accelArray[i].lookup != nil {
			c.setLookupMask(fbPlan.maskArray[i])
			c.substituteLookup(&fbPlan.accelArray[i])
		}
	}
}
