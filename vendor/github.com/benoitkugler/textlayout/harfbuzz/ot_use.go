package harfbuzz

import (
	"fmt"

	tt "github.com/benoitkugler/textlayout/fonts/truetype"
	ucd "github.com/benoitkugler/textlayout/unicodedata"
)

// ported from harfbuzz/src/hb-ot-shape-complex-use.cc Copyright © 2015  Mozilla Foundation. Google, Inc. Jonathan Kew, Behdad Esfahbod

/*
 * Universal Shaping Engine.
 * https://docs.microsoft.com/en-us/typography/script-development/use
 */

var _ otComplexShaper = (*complexShaperUSE)(nil)

/*
 * Basic features.
 * These features are applied all at once, before reordering.
 */
var useBasicFeatures = [...]tt.Tag{
	tt.NewTag('r', 'k', 'r', 'f'),
	tt.NewTag('a', 'b', 'v', 'f'),
	tt.NewTag('b', 'l', 'w', 'f'),
	tt.NewTag('h', 'a', 'l', 'f'),
	tt.NewTag('p', 's', 't', 'f'),
	tt.NewTag('v', 'a', 't', 'u'),
	tt.NewTag('c', 'j', 'c', 't'),
}

var useTopographicalFeatures = [...]tt.Tag{
	tt.NewTag('i', 's', 'o', 'l'),
	tt.NewTag('i', 'n', 'i', 't'),
	tt.NewTag('m', 'e', 'd', 'i'),
	tt.NewTag('f', 'i', 'n', 'a'),
}

/* Same order as useTopographicalFeatures. */
const (
	joiningFormIsol = iota
	joiningFormInit
	joiningFormMedi
	joiningFormFina
	joiningFormNone
)

/*
 * Other features.
 * These features are applied all at once, after reordering and
 * clearing syllables.
 */
var useOtherFeatures = [...]tt.Tag{
	tt.NewTag('a', 'b', 'v', 's'),
	tt.NewTag('b', 'l', 'w', 's'),
	tt.NewTag('h', 'a', 'l', 'n'),
	tt.NewTag('p', 'r', 'e', 's'),
	tt.NewTag('p', 's', 't', 's'),
}

type useShapePlan struct {
	arabicPlan *arabicShapePlan
	rphfMask   GlyphMask
}

type complexShaperUSE struct {
	complexShaperNil

	plan useShapePlan
}

func (cs *complexShaperUSE) collectFeatures(plan *otShapePlanner) {
	map_ := &plan.map_

	/* Do this before any lookups have been applied. */
	map_.addGSUBPause(cs.setupSyllablesUse)

	/* "Default glyph pre-processing group" */
	map_.enableFeature(tt.NewTag('l', 'o', 'c', 'l'))
	map_.enableFeature(tt.NewTag('c', 'c', 'm', 'p'))
	map_.enableFeature(tt.NewTag('n', 'u', 'k', 't'))
	map_.enableFeatureExt(tt.NewTag('a', 'k', 'h', 'n'), ffManualZWJ, 1)

	/* "Reordering group" */
	map_.addGSUBPause(clearSubstitutionFlags)
	map_.addFeatureExt(tt.NewTag('r', 'p', 'h', 'f'), ffManualZWJ, 1)
	map_.addGSUBPause(cs.recordRphfUse)
	map_.addGSUBPause(clearSubstitutionFlags)
	map_.enableFeatureExt(tt.NewTag('p', 'r', 'e', 'f'), ffManualZWJ, 1)
	map_.addGSUBPause(recordPrefUse)

	/* "Orthographic unit shaping group" */
	for _, basicFeat := range useBasicFeatures {
		map_.enableFeatureExt(basicFeat, ffManualZWJ, 1)
	}

	map_.addGSUBPause(reorderUse)
	map_.addGSUBPause(clearSyllables)

	/* "Topographical features" */
	for _, topoFeat := range useTopographicalFeatures {
		map_.addFeature(topoFeat)
	}
	map_.addGSUBPause(nil)

	/* "Standard typographic presentation" */
	for _, otherFeat := range useOtherFeatures {
		map_.enableFeatureExt(otherFeat, ffManualZWJ, 1)
	}
}

func (cs *complexShaperUSE) dataCreate(plan *otShapePlan) {
	var usePlan useShapePlan

	usePlan.rphfMask = plan.map_.getMask1(tt.NewTag('r', 'p', 'h', 'f'))

	if ucd.HasArabicJoining(plan.props.Script) {
		pl := newArabicPlan(plan)
		usePlan.arabicPlan = &pl
	}

	cs.plan = usePlan
}

func (cs *complexShaperUSE) setupMasks(plan *otShapePlan, buffer *Buffer, _ *Font) {
	usePlan := cs.plan
	/* Do this before allocating complexCategory. */
	if usePlan.arabicPlan != nil {
		usePlan.arabicPlan.setupMasks(buffer, plan.props.Script)
	}

	/* We cannot setup masks here.  We save information about characters
	* and setup masks later on in a pause-callback. */

	info := buffer.Info
	for i := range info {
		info[i].complexCategory = getUSECategory(info[i].codepoint)
	}
}

func (cs *complexShaperUSE) setupRphfMask(buffer *Buffer) {
	usePlan := cs.plan

	mask := usePlan.rphfMask
	if mask == 0 {
		return
	}

	info := buffer.Info
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		limit := 1
		if info[start].complexCategory != useSyllableMachine_ex_R {
			limit = min(3, end-start)
		}
		for i := start; i < start+limit; i++ {
			info[i].Mask |= mask
		}
	}
}

func (cs *complexShaperUSE) setupTopographicalMasks(plan *otShapePlan, buffer *Buffer) {
	if cs.plan.arabicPlan != nil {
		return
	}
	var (
		masks    [4]GlyphMask
		allMasks uint32
	)
	for i := range masks {
		masks[i] = plan.map_.getMask1(useTopographicalFeatures[i])
		if masks[i] == plan.map_.globalMask {
			masks[i] = 0
		}
		allMasks |= masks[i]
	}
	if allMasks == 0 {
		return
	}
	otherMasks := ^allMasks

	lastStart := 0
	lastForm := joiningFormNone
	info := buffer.Info
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		syllableType := info[start].syllable & 0x0F
		switch syllableType {
		case useSymbolCluster, useHieroglyphCluster, useNonCluster:
			// these don't join.  Nothing to do.
			lastForm = joiningFormNone

		case useViramaTerminatedCluster, useSakotTerminatedCluster, useStandardCluster, useNumberJoinerTerminatedCluster, useNumeralCluster, useBrokenCluster:
			join := lastForm == joiningFormFina || lastForm == joiningFormIsol
			if join {
				// fixup previous syllable's form.
				if lastForm == joiningFormFina {
					lastForm = joiningFormMedi
				} else {
					lastForm = joiningFormInit
				}
				for i := lastStart; i < start; i++ {
					info[i].Mask = (info[i].Mask & otherMasks) | masks[lastForm]
				}
			}

			// form for this syllable.
			lastForm = joiningFormIsol
			if join {
				lastForm = joiningFormFina
			}
			for i := start; i < end; i++ {
				info[i].Mask = (info[i].Mask & otherMasks) | masks[lastForm]
			}
		}

		lastStart = start
	}
}

func (cs *complexShaperUSE) setupSyllablesUse(plan *otShapePlan, _ *Font, buffer *Buffer) {
	findSyllablesUse(buffer)
	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		buffer.unsafeToBreak(start, end)
	}
	cs.setupRphfMask(buffer)
	cs.setupTopographicalMasks(plan, buffer)
}

func (cs *complexShaperUSE) recordRphfUse(plan *otShapePlan, _ *Font, buffer *Buffer) {
	usePlan := cs.plan

	mask := usePlan.rphfMask
	if mask == 0 {
		return
	}
	info := buffer.Info

	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		// mark a substituted repha as USE(R).
		for i := start; i < end && (info[i].Mask&mask) != 0; i++ {
			if glyphInfoSubstituted(&info[i]) {
				info[i].complexCategory = useSyllableMachine_ex_R
				break
			}
		}
	}
}

func recordPrefUse(_ *otShapePlan, _ *Font, buffer *Buffer) {
	info := buffer.Info

	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		// mark a substituted pref as VPre, as they behave the same way.
		for i := start; i < end; i++ {
			if glyphInfoSubstituted(&info[i]) {
				info[i].complexCategory = useSyllableMachine_ex_VPre
				break
			}
		}
	}
}

func isHalantUse(info *GlyphInfo) bool {
	return (info.complexCategory == useSyllableMachine_ex_H || info.complexCategory == useSyllableMachine_ex_HVM) &&
		!info.ligated()
}

func reorderSyllableUse(buffer *Buffer, start, end int) {
	syllableType := (buffer.Info[start].syllable & 0x0F)
	/* Only a few syllable types need reordering. */
	const mask = 1<<useViramaTerminatedCluster |
		1<<useSakotTerminatedCluster |
		1<<useStandardCluster |
		1<<useBrokenCluster
	if 1<<syllableType&mask == 0 {
		return
	}

	info := buffer.Info

	const postBaseFlags64 int64 = (1<<useSyllableMachine_ex_FAbv |
		1<<useSyllableMachine_ex_FBlw |
		1<<useSyllableMachine_ex_FPst |
		1<<useSyllableMachine_ex_MAbv |
		1<<useSyllableMachine_ex_MBlw |
		1<<useSyllableMachine_ex_MPst |
		1<<useSyllableMachine_ex_MPre |
		1<<useSyllableMachine_ex_VAbv |
		1<<useSyllableMachine_ex_VBlw |
		1<<useSyllableMachine_ex_VPst |
		1<<useSyllableMachine_ex_VPre |
		1<<useSyllableMachine_ex_VMAbv |
		1<<useSyllableMachine_ex_VMBlw |
		1<<useSyllableMachine_ex_VMPst |
		1<<useSyllableMachine_ex_VMPre)

	/* Move things forward. */
	if info[start].complexCategory == useSyllableMachine_ex_R && end-start > 1 {
		/* Got a repha.  Reorder it towards the end, but before the first post-base
		 * glyph. */
		for i := start + 1; i < end; i++ {
			isPostBaseGlyph := (int64(1<<(info[i].complexCategory))&postBaseFlags64) != 0 ||
				isHalantUse(&info[i])
			if isPostBaseGlyph || i == end-1 {
				/* If we hit a post-base glyph, move before it; otherwise move to the
				 * end. Shift things in between backward. */

				if isPostBaseGlyph {
					i--
				}

				buffer.mergeClusters(start, i+1)
				t := info[start]
				copy(info[start:i], info[start+1:])
				info[i] = t

				break
			}
		}
	}

	/* Move things back. */
	j := start
	for i := start; i < end; i++ {
		flag := 1 << (info[i].complexCategory)
		if isHalantUse(&info[i]) {
			/* If we hit a halant, move after it; otherwise move to the beginning, and
			* shift things in between forward. */
			j = i + 1
		} else if flag&(1<<useSyllableMachine_ex_VPre|1<<useSyllableMachine_ex_VMPre) != 0 &&
			/* Only move the first component of a MultipleSubst. */
			info[i].getLigComp() == 0 && j < i {
			buffer.mergeClusters(j, i+1)
			t := info[i]
			copy(info[j+1:], info[j:i])
			info[j] = t
		}
	}
}

func reorderUse(_ *otShapePlan, font *Font, buffer *Buffer) {
	if debugMode >= 1 {
		fmt.Println("USE - start reordering USE")
	}
	syllabicInsertDottedCircles(font, buffer, useBrokenCluster,
		useSyllableMachine_ex_B, useSyllableMachine_ex_R, -1)

	iter, count := buffer.syllableIterator()
	for start, end := iter.next(); start < count; start, end = iter.next() {
		reorderSyllableUse(buffer, start, end)
	}
	if debugMode >= 1 {
		fmt.Println("USE - end reordering USE")
	}
}

func (cs *complexShaperUSE) preprocessText(_ *otShapePlan, buffer *Buffer, _ *Font) {
	preprocessTextVowelConstraints(buffer)
}

func (cs *complexShaperUSE) compose(_ *otNormalizeContext, a, b rune) (rune, bool) {
	// avoid recomposing split matras.
	if uni.generalCategory(a).isMark() {
		return 0, false
	}

	return uni.compose(a, b)
}

func (complexShaperUSE) marksBehavior() (zeroWidthMarks, bool) {
	return zeroWidthMarksByGdefEarly, false
}

func (complexShaperUSE) normalizationPreference() normalizationMode {
	return nmComposedDiacriticsNoShortCircuit
}
