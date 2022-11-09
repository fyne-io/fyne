package harfbuzz

import (
	"fmt"
	"math"

	"github.com/benoitkugler/textlayout/fonts"
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-ot-layout-gsubgpos.hh Copyright © 2007,2008,2009,2010  Red Hat, Inc. 2010,2012  Google, Inc.  Behdad Esfahbod

// GSUB or GPOS lookup
type layoutLookup interface {
	// accumulate the subtables coverage into the diggest
	collectCoverage(*setDigest)
	// walk the subtables to add them to the context
	dispatchSubtables(*getSubtablesContext)

	// walk the subtables and apply the sub/pos
	dispatchApply(ctx *otApplyContext) bool

	Props() uint32
	isReverse() bool
}

/*
 * GSUB/GPOS Common
 */

const ignoreFlags = tt.IgnoreBaseGlyphs | tt.IgnoreLigatures | tt.IgnoreMarks

// use a digest to speedup match
type otLayoutLookupAccelerator struct {
	lookup    layoutLookup
	subtables getSubtablesContext
	digest    setDigest
}

func (ac *otLayoutLookupAccelerator) init(lookup layoutLookup) {
	ac.lookup = lookup
	ac.digest = setDigest{}
	lookup.collectCoverage(&ac.digest)
	ac.subtables = nil
	lookup.dispatchSubtables(&ac.subtables)
}

// apply the subtables and stops at the first success.
func (ac *otLayoutLookupAccelerator) apply(c *otApplyContext) bool {
	for _, table := range ac.subtables {
		if table.apply(c) {
			return true
		}
	}
	return false
}

// represents one layout subtable, with its own coverage
type applicable struct {
	obj interface{ apply(c *otApplyContext) bool }

	digest setDigest
}

func newGSUBApplicable(table tt.GSUBSubtable) applicable {
	ap := applicable{obj: gsubSubtable(table)}
	ap.digest.collectCoverage(table.Coverage)
	return ap
}

func newGPOSApplicable(table tt.GPOSSubtable) applicable {
	ap := applicable{obj: gposSubtable(table)}
	ap.digest.collectCoverage(table.Coverage)
	return ap
}

func (ap applicable) apply(c *otApplyContext) bool {
	return ap.digest.mayHave(c.buffer.cur(0).Glyph) && ap.obj.apply(c)
}

type getSubtablesContext []applicable

// one for GSUB, one for GPOS (known at compile time)
type otProxyMeta struct {
	recurseFunc recurseFunc
	tableIndex  int
	inplace     bool
}

var (
	proxyGSUB = otProxyMeta{tableIndex: 0, inplace: false, recurseFunc: applyRecurseGSUB}
	proxyGPOS = otProxyMeta{tableIndex: 1, inplace: true, recurseFunc: applyRecurseGPOS}
)

type otProxy struct {
	otProxyMeta
	accels []otLayoutLookupAccelerator
}

type wouldApplyContext struct {
	face        fonts.FaceMetrics
	glyphs      []fonts.GID
	indices     []uint16 // see get1N
	zeroContext bool
}

// `value` interpretation is dictated by the context
type matcherFunc = func(gid fonts.GID, value uint16) bool

// interprets `value` as a Glyph
func matchGlyph(gid fonts.GID, value uint16) bool { return gid == fonts.GID(value) }

// interprets `value` as a Class
func matchClass(class tt.Class) matcherFunc {
	return func(gid fonts.GID, value uint16) bool {
		c, _ := class.ClassID(gid)
		return uint16(c) == value
	}
}

// interprets `value` as an index in coverage array
func matchCoverage(covs []tt.Coverage) matcherFunc {
	return func(gid fonts.GID, value uint16) bool {
		_, covered := covs[value].Index(gid)
		return covered
	}
}

const (
	no = iota
	yes
	maybe
)

type otApplyContextMatcher struct {
	matchFunc   matcherFunc
	lookupProps uint32
	mask        GlyphMask
	ignoreZWNJ  bool
	ignoreZWJ   bool
	syllable    uint8
}

func (m otApplyContextMatcher) mayMatch(info *GlyphInfo, glyphData []uint16) uint8 {
	if info.Mask&m.mask == 0 || (m.syllable != 0 && m.syllable != info.syllable) {
		return no
	}

	if m.matchFunc != nil {
		if m.matchFunc(info.Glyph, glyphData[0]) {
			return yes
		}
		return no
	}

	return maybe
}

func (m otApplyContextMatcher) maySkip(c *otApplyContext, info *GlyphInfo) uint8 {
	if !c.checkGlyphProperty(info, m.lookupProps) {
		return yes
	}

	if info.isDefaultIgnorableAndNotHidden() && (m.ignoreZWNJ || !info.isZwnj()) &&
		(m.ignoreZWJ || !info.isZwj()) {
		return maybe
	}

	return no
}

type skippingIterator struct {
	c       *otApplyContext
	matcher otApplyContextMatcher

	matchGlyphDataArray []uint16
	matchGlyphDataStart int // start as index in match_glyph_data_array

	idx      int
	numItems int
	end      int
}

func (it *skippingIterator) init(c *otApplyContext, contextMatch bool) {
	it.c = c
	it.setMatchFunc(nil, nil)
	it.matcher.matchFunc = nil
	it.matcher.lookupProps = c.lookupProps
	/* Ignore ZWNJ if we are matching GPOS, or matching GSUB context and asked to. */
	it.matcher.ignoreZWNJ = c.tableIndex == 1 || (contextMatch && c.autoZWNJ)
	/* Ignore ZWJ if we are matching context, or asked to. */
	it.matcher.ignoreZWJ = contextMatch || c.autoZWJ
	if contextMatch {
		it.matcher.mask = math.MaxUint32
	} else {
		it.matcher.mask = c.lookupMask
	}
}

// 	 void set_lookup_props (uint lookupProps)
// 	 {
// 	   matcher.set_lookup_props (lookupProps);
// 	 }

func (it *skippingIterator) setMatchFunc(matchFunc matcherFunc, glyphData []uint16) {
	it.matcher.matchFunc = matchFunc
	it.matchGlyphDataArray = glyphData
	it.matchGlyphDataStart = 0
}

func (it *skippingIterator) reset(startIndex, numItems int) {
	it.idx = startIndex
	it.numItems = numItems
	it.end = len(it.c.buffer.Info)
	if startIndex == it.c.buffer.idx {
		it.matcher.syllable = it.c.buffer.cur(0).syllable
	} else {
		it.matcher.syllable = 0
	}
}

func (it *skippingIterator) reject() {
	it.numItems++
	if len(it.matchGlyphDataArray) != 0 {
		it.matchGlyphDataStart--
	}
}

func (it *skippingIterator) maySkip(info *GlyphInfo) uint8 { return it.matcher.maySkip(it.c, info) }

func (it *skippingIterator) next() bool {
	for it.idx+it.numItems < it.end {
		it.idx++
		info := &it.c.buffer.Info[it.idx]

		skip := it.matcher.maySkip(it.c, info)
		if skip == yes {
			continue
		}

		match := it.matcher.mayMatch(info, it.matchGlyphDataArray[it.matchGlyphDataStart:])
		if match == yes || (match == maybe && skip == no) {
			it.numItems--
			if len(it.matchGlyphDataArray) != 0 {
				it.matchGlyphDataStart++
			}
			return true
		}

		if skip == no {
			return false
		}
	}
	return false
}

func (it *skippingIterator) prev() bool {
	L := len(it.c.buffer.outInfo)
	//    assert (num_items > 0);
	for it.idx > it.numItems-1 {
		it.idx--
		var info *GlyphInfo
		if it.idx < L {
			info = &it.c.buffer.outInfo[it.idx]
		} else {
			// we are in "position mode" : outInfo is not used anymore
			// in the C implementation, outInfo and info now are sharing the same storage
			info = &it.c.buffer.Info[it.idx]
		}

		skip := it.matcher.maySkip(it.c, info)
		if skip == yes {
			continue
		}

		match := it.matcher.mayMatch(info, it.matchGlyphDataArray[it.matchGlyphDataStart:])
		if match == yes || (match == maybe && skip == no) {
			it.numItems--
			if len(it.matchGlyphDataArray) != 0 {
				it.matchGlyphDataStart++
			}
			return true
		}

		if skip == no {
			return false
		}
	}
	return false
}

type recurseFunc = func(c *otApplyContext, lookupIndex uint16) bool

type otApplyContext struct {
	face   fonts.FaceMetrics
	font   *Font
	buffer *Buffer

	recurseFunc recurseFunc
	gdef        tt.TableGDEF
	varStore    tt.VariationStore
	indices     []uint16 // see get1N()

	iterContext skippingIterator
	iterInput   skippingIterator

	nestingLevelLeft int
	tableIndex       int
	lookupMask       GlyphMask
	lookupProps      uint32
	randomState      uint32
	lookupIndex      uint16
	direction        Direction

	hasGlyphClasses bool
	autoZWNJ        bool
	autoZWJ         bool
	random          bool
}

func newOtApplyContext(tableIndex int, font *Font, buffer *Buffer) otApplyContext {
	var out otApplyContext
	out.font = font
	out.face = font.face
	out.buffer = buffer
	out.gdef = font.otTables.GDEF
	out.varStore = out.gdef.VariationStore
	out.direction = buffer.Props.Direction
	out.lookupMask = 1
	out.tableIndex = tableIndex
	out.lookupIndex = math.MaxUint16
	out.nestingLevelLeft = maxNestingLevel
	out.hasGlyphClasses = out.gdef.Class != nil
	out.autoZWNJ = true
	out.autoZWJ = true
	out.randomState = 1

	out.initIters()
	return out
}

func (c *otApplyContext) initIters() {
	c.iterInput.init(c, false)
	c.iterContext.init(c, true)
}

func (c *otApplyContext) setLookupMask(mask GlyphMask) {
	c.lookupMask = mask
	c.initIters()
}

func (c *otApplyContext) setAutoZWNJ(autoZwnj bool) {
	c.autoZWNJ = autoZwnj
	c.initIters()
}

func (c *otApplyContext) setAutoZWJ(autoZwj bool) {
	c.autoZWJ = autoZwj
	c.initIters()
}

func (c *otApplyContext) setLookupProps(lookupProps uint32) {
	c.lookupProps = lookupProps
	c.initIters()
}

func (c *otApplyContext) applyRecurseLookup(lookupIndex uint16, l layoutLookup) bool {
	savedLookupProps := c.lookupProps
	savedLookupIndex := c.lookupIndex

	c.lookupIndex = lookupIndex
	c.setLookupProps(l.Props())

	ret := l.dispatchApply(c)

	c.lookupIndex = savedLookupIndex
	c.setLookupProps(savedLookupProps)
	return ret
}

func (c *otApplyContext) substituteLookup(accel *otLayoutLookupAccelerator) {
	c.applyString(proxyGSUB, accel)
}

func (c *otApplyContext) checkGlyphProperty(info *GlyphInfo, matchProps uint32) bool {
	glyphProps := info.glyphProps

	/* Not covered, if, for example, glyph class is ligature and
	 * matchProps includes LookupFlags::IgnoreLigatures */
	if (glyphProps & uint16(matchProps) & ignoreFlags) != 0 {
		return false
	}

	if glyphProps&tt.Mark != 0 {
		return c.matchPropertiesMark(info.Glyph, glyphProps, matchProps)
	}

	return true
}

func (c *otApplyContext) matchPropertiesMark(glyph fonts.GID, glyphProps uint16, matchProps uint32) bool {
	/* If using mark filtering sets, the high uint16 of
	 * matchProps has the set index. */
	if tt.LookupFlag(matchProps)&tt.UseMarkFilteringSet != 0 {
		_, has := c.gdef.MarkGlyphSet[matchProps>>16].Index(glyph)
		return has
	}

	/* The second byte of matchProps has the meaning
	 * "ignore marks of attachment type different than
	 * the attachment type specified." */
	if tt.LookupFlag(matchProps)&tt.MarkAttachmentType != 0 {
		return uint16(matchProps)&tt.MarkAttachmentType == (glyphProps & tt.MarkAttachmentType)
	}

	return true
}

func (c *otApplyContext) setGlyphProps(glyphIndex fonts.GID) {
	c.setGlyphPropsExt(glyphIndex, 0, false, false)
}

func (c *otApplyContext) setGlyphPropsExt(glyphIndex fonts.GID, classGuess uint16, ligature, component bool) {
	addIn := c.buffer.cur(0).glyphProps & preserve
	addIn |= substituted
	if ligature {
		addIn |= ligated
		/* In the only place that the MULTIPLIED bit is used, Uniscribe
		* seems to only care about the "last" transformation between
		* Ligature and Multiple substitutions.  Ie. if you ligate, expand,
		* and ligate again, it forgives the multiplication and acts as
		* if only ligation happened.  As such, clear MULTIPLIED bit.
		 */
		addIn &= ^multiplied
	}
	if component {
		addIn |= multiplied
	}
	if c.hasGlyphClasses {
		c.buffer.cur(0).glyphProps = addIn | c.gdef.GetGlyphProps(glyphIndex)
	} else if classGuess != 0 {
		c.buffer.cur(0).glyphProps = addIn | classGuess
	}
}

func (c *otApplyContext) replaceGlyph(glyphIndex fonts.GID) {
	c.setGlyphProps(glyphIndex)
	c.buffer.replaceGlyphIndex(glyphIndex)
}

func (c *otApplyContext) randomNumber() uint32 {
	/* http://www.cplusplus.com/reference/random/minstd_rand/ */
	c.randomState = c.randomState * 48271 % 2147483647
	return c.randomState
}

func (c *otApplyContext) applyRuleSet(ruleSet []tt.SequenceRule, match matcherFunc) bool {
	for _, rule := range ruleSet {
		// the first which match is applied
		applied := c.contextApplyLookup(rule.Input, rule.Lookups, match)
		if applied {
			return true
		}
	}
	return false
}

func (c *otApplyContext) applyChainRuleSet(ruleSet []tt.ChainedSequenceRule, match [3]matcherFunc) bool {
	for i, rule := range ruleSet {

		if debugMode >= 2 {
			fmt.Println("APPLY - chain rule number", i)
		}

		b := c.chainContextApplyLookup(rule.Backtrack, rule.Input, rule.Lookahead, rule.Lookups, match)
		if b { // stop at the first application
			return true
		}
	}
	return false
}

//  `input` starts with second glyph (`inputCount` = len(input)+1)
func (c *otApplyContext) contextApplyLookup(input []uint16, lookupRecord []tt.SequenceLookup, lookupContext matcherFunc) bool {
	matchLength := 0
	var matchPositions [maxContextLength]int
	hasMatch, matchLength, _ := c.matchInput(input, lookupContext, &matchPositions)
	if !hasMatch {
		return false
	}
	c.buffer.unsafeToBreak(c.buffer.idx, c.buffer.idx+matchLength)
	c.applyLookup(len(input)+1, &matchPositions, lookupRecord, matchLength)
	return true
}

//  `input` starts with second glyph (`inputCount` = len(input)+1)
// lookupsContexts : backtrack, input, lookahead
func (c *otApplyContext) chainContextApplyLookup(backtrack, input, lookahead []uint16,
	lookupRecord []tt.SequenceLookup, lookupContexts [3]matcherFunc) bool {
	var matchPositions [maxContextLength]int

	hasMatch, matchLength, _ := c.matchInput(input, lookupContexts[1], &matchPositions)
	if !hasMatch {
		return false
	}

	hasMatch, startIndex := c.matchBacktrack(backtrack, lookupContexts[0])
	if !hasMatch {
		return false
	}

	hasMatch, endIndex := c.matchLookahead(lookahead, lookupContexts[2], matchLength)
	if !hasMatch {
		return false
	}

	c.buffer.unsafeToBreakFromOutbuffer(startIndex, endIndex)
	c.applyLookup(len(input)+1, &matchPositions, lookupRecord, matchLength)
	return true
}

func (c *wouldApplyContext) wouldApplyLookupContext1(data tt.LookupContext1, index int) bool {
	if index >= len(data) { // index is not sanitized in tt.Parse
		return false
	}
	ruleSet := data[index]
	return c.wouldApplyRuleSet(ruleSet, matchGlyph)
}

func (c *wouldApplyContext) wouldApplyLookupContext2(data tt.LookupContext2, index int, glyphID fonts.GID) bool {
	class, _ := data.Class.ClassID(glyphID)
	ruleSet := data.SequenceSets[class]
	return c.wouldApplyRuleSet(ruleSet, matchClass(data.Class))
}

func (c *wouldApplyContext) wouldApplyLookupContext3(data tt.LookupContext3, index int) bool {
	covIndices := get1N(&c.indices, 1, len(data.Coverages))
	return c.wouldMatchInput(covIndices, matchCoverage(data.Coverages))
}

func (c *wouldApplyContext) wouldApplyRuleSet(ruleSet []tt.SequenceRule, match matcherFunc) bool {
	for _, rule := range ruleSet {
		if c.wouldMatchInput(rule.Input, match) {
			return true
		}
	}
	return false
}

func (c *wouldApplyContext) wouldApplyChainRuleSet(ruleSet []tt.ChainedSequenceRule, inputMatch matcherFunc) bool {
	for _, rule := range ruleSet {
		if c.wouldApplyChainLookup(rule.Backtrack, rule.Input, rule.Lookahead, inputMatch) {
			return true
		}
	}
	return false
}

func (c *wouldApplyContext) wouldApplyLookupChainedContext1(data tt.LookupChainedContext1, index int) bool {
	if index >= len(data) { // index is not sanitized in tt.Parse
		return false
	}
	ruleSet := data[index]
	return c.wouldApplyChainRuleSet(ruleSet, matchGlyph)
}

func (c *wouldApplyContext) wouldApplyLookupChainedContext2(data tt.LookupChainedContext2, index int, glyphID fonts.GID) bool {
	class, _ := data.InputClass.ClassID(glyphID)
	ruleSet := data.SequenceSets[class]
	return c.wouldApplyChainRuleSet(ruleSet, matchClass(data.InputClass))
}

func (c *wouldApplyContext) wouldApplyLookupChainedContext3(data tt.LookupChainedContext3, index int) bool {
	lB, lI, lL := len(data.Backtrack), len(data.Input), len(data.Lookahead)
	return c.wouldApplyChainLookup(get1N(&c.indices, 0, lB), get1N(&c.indices, 1, lI), get1N(&c.indices, 0, lL),
		matchCoverage(data.Input))
}

// `input` starts with second glyph (`inputCount` = len(input)+1)
// only the input lookupsContext is needed
func (c *wouldApplyContext) wouldApplyChainLookup(backtrack, input, lookahead []uint16, inputLookupContext matcherFunc) bool {
	contextOk := true
	if c.zeroContext {
		contextOk = len(backtrack) == 0 && len(lookahead) == 0
	}
	return contextOk && c.wouldMatchInput(input, inputLookupContext)
}

// `input` starts with second glyph (`count` = len(input)+1)
func (c *wouldApplyContext) wouldMatchInput(input []uint16, matchFunc matcherFunc) bool {
	if len(c.glyphs) != len(input)+1 {
		return false
	}

	for i, glyph := range input {
		if !matchFunc(c.glyphs[i+1], glyph) {
			return false
		}
	}

	return true
}

// `input` starts with second glyph (`inputCount` = len(input)+1)
func (c *otApplyContext) matchInput(input []uint16, matchFunc matcherFunc,
	matchPositions *[maxContextLength]int) (bool, int, uint8) {
	count := len(input) + 1
	if count > maxContextLength {
		return false, 0, 0
	}
	buffer := c.buffer
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, count-1)
	skippyIter.setMatchFunc(matchFunc, input)

	/*
	* This is perhaps the trickiest part of OpenType...  Remarks:
	*
	* - If all components of the ligature were marks, we call this a mark ligature.
	*
	* - If there is no GDEF, and the ligature is NOT a mark ligature, we categorize
	*   it as a ligature glyph.
	*
	* - Ligatures cannot be formed across glyphs attached to different components
	*   of previous ligatures.  Eg. the sequence is LAM,SHADDA,LAM,FATHA,HEH, and
	*   LAM,LAM,HEH form a ligature, leaving SHADDA,FATHA next to eachother.
	*   However, it would be wrong to ligate that SHADDA,FATHA sequence.
	*   There are a couple of exceptions to this:
	*
	*   o If a ligature tries ligating with marks that belong to it itself, go ahead,
	*     assuming that the font designer knows what they are doing (otherwise it can
	*     break Indic stuff when a matra wants to ligate with a conjunct,
	*
	*   o If two marks want to ligate and they belong to different components of the
	*     same ligature glyph, and said ligature glyph is to be ignored according to
	*     mark-filtering rules, then allow.
	*     https://github.com/harfbuzz/harfbuzz/issues/545
	 */

	totalComponentCount := buffer.cur(0).getLigNumComps()

	firstLigID := buffer.cur(0).getLigID()
	firstLigComp := buffer.cur(0).getLigComp()

	const (
		ligbaseNotChecked = iota
		ligbaseMayNotSkip
		ligbaseMaySkip
	)
	ligbase := ligbaseNotChecked
	matchPositions[0] = buffer.idx
	for i := 1; i < count; i++ {
		if !skippyIter.next() {
			return false, 0, 0
		}

		matchPositions[i] = skippyIter.idx

		thisLigID := buffer.Info[skippyIter.idx].getLigID()
		thisLigComp := buffer.Info[skippyIter.idx].getLigComp()
		if firstLigID != 0 && firstLigComp != 0 {
			/* If first component was attached to a previous ligature component,
			* all subsequent components should be attached to the same ligature
			* component, otherwise we shouldn't ligate them... */
			if firstLigID != thisLigID || firstLigComp != thisLigComp {
				/* ...unless, we are attached to a base ligature and that base
				 * ligature is ignorable. */
				if ligbase == ligbaseNotChecked {
					found := false
					out := buffer.outInfo
					j := len(out)
					for j != 0 && out[j-1].getLigID() == firstLigID {
						if out[j-1].getLigComp() == 0 {
							j--
							found = true
							break
						}
						j--
					}

					if found && skippyIter.maySkip(&out[j]) == yes {
						ligbase = ligbaseMaySkip
					} else {
						ligbase = ligbaseMayNotSkip
					}
				}

				if ligbase == ligbaseMayNotSkip {
					return false, 0, 0
				}
			}
		} else {
			/* If first component was NOT attached to a previous ligature component,
			* all subsequent components should also NOT be attached to any ligature
			* component, unless they are attached to the first component itself! */
			if thisLigID != 0 && thisLigComp != 0 && (thisLigID != firstLigID) {
				return false, 0, 0
			}
		}

		totalComponentCount += buffer.Info[skippyIter.idx].getLigNumComps()
	}

	endOffset := skippyIter.idx - buffer.idx + 1

	return true, endOffset, totalComponentCount
}

// `count` and `matchPositions` include the first glyph
func (c *otApplyContext) ligateInput(count int, matchPositions [maxContextLength]int,
	matchLength int, ligGlyph fonts.GID, totalComponentCount uint8) {
	buffer := c.buffer

	buffer.mergeClusters(buffer.idx, buffer.idx+matchLength)

	/* - If a base and one or more marks ligate, consider that as a base, NOT
	*   ligature, such that all following marks can still attach to it.
	*   https://github.com/harfbuzz/harfbuzz/issues/1109
	*
	* - If all components of the ligature were marks, we call this a mark ligature.
	*   If it *is* a mark ligature, we don't allocate a new ligature id, and leave
	*   the ligature to keep its old ligature id.  This will allow it to attach to
	*   a base ligature in GPOS.  Eg. if the sequence is: LAM,LAM,SHADDA,FATHA,HEH,
	*   and LAM,LAM,HEH for a ligature, they will leave SHADDA and FATHA with a
	*   ligature id and component value of 2.  Then if SHADDA,FATHA form a ligature
	*   later, we don't want them to lose their ligature id/component, otherwise
	*   GPOS will fail to correctly position the mark ligature on top of the
	*   LAM,LAM,HEH ligature.  See:
	*     https://bugzilla.gnome.org/show_bug.cgi?id=676343
	*
	* - If a ligature is formed of components that some of which are also ligatures
	*   themselves, and those ligature components had marks attached to *their*
	*   components, we have to attach the marks to the new ligature component
	*   positions!  Now *that*'s tricky!  And these marks may be following the
	*   last component of the whole sequence, so we should loop forward looking
	*   for them and update them.
	*
	*   Eg. the sequence is LAM,LAM,SHADDA,FATHA,HEH, and the font first forms a
	*   'calt' ligature of LAM,HEH, leaving the SHADDA and FATHA with a ligature
	*   id and component == 1.  Now, during 'liga', the LAM and the LAM-HEH ligature
	*   form a LAM-LAM-HEH ligature.  We need to reassign the SHADDA and FATHA to
	*   the new ligature with a component value of 2.
	*
	*   This in fact happened to a font...  See:
	*   https://bugzilla.gnome.org/show_bug.cgi?id=437633
	 */

	isBaseLigature := buffer.Info[matchPositions[0]].isBaseGlyph()
	isMarkLigature := buffer.Info[matchPositions[0]].isMark()
	for i := 1; i < count; i++ {
		if !buffer.Info[matchPositions[i]].isMark() {
			isBaseLigature = false
			isMarkLigature = false
			break
		}
	}
	isLigature := !isBaseLigature && !isMarkLigature

	klass, ligID := uint16(0), uint8(0)
	if isLigature {
		klass = tt.Ligature
		ligID = buffer.allocateLigID()
	}
	lastLigID := buffer.cur(0).getLigID()
	lastNumComponents := buffer.cur(0).getLigNumComps()
	componentsSoFar := lastNumComponents

	if isLigature {
		buffer.cur(0).setLigPropsForLigature(ligID, totalComponentCount)
		if buffer.cur(0).unicode.generalCategory() == nonSpacingMark {
			buffer.cur(0).setGeneralCategory(otherLetter)
		}
	}

	// ReplaceGlyph_with_ligature
	c.setGlyphPropsExt(ligGlyph, klass, true, false)
	buffer.replaceGlyphIndex(ligGlyph)

	for i := 1; i < count; i++ {
		for buffer.idx < matchPositions[i] {
			if isLigature {
				thisComp := buffer.cur(0).getLigComp()
				if thisComp == 0 {
					thisComp = lastNumComponents
				}
				newLigComp := componentsSoFar - lastNumComponents +
					min8(thisComp, lastNumComponents)
				buffer.cur(0).setLigPropsForMark(ligID, newLigComp)
			}
			buffer.nextGlyph()
		}

		lastLigID = buffer.cur(0).getLigID()
		lastNumComponents = buffer.cur(0).getLigNumComps()
		componentsSoFar += lastNumComponents

		/* Skip the base glyph */
		buffer.skipGlyph()
	}

	if !isMarkLigature && lastLigID != 0 {
		/* Re-adjust components for any marks following. */
		for i := buffer.idx; i < len(buffer.Info); i++ {
			if lastLigID != buffer.Info[i].getLigID() {
				break
			}

			thisComp := buffer.Info[i].getLigComp()
			if thisComp == 0 {
				break
			}

			newLigComp := componentsSoFar - lastNumComponents +
				min8(thisComp, lastNumComponents)
			buffer.Info[i].setLigPropsForMark(ligID, newLigComp)
		}
	}
}

func (c *otApplyContext) recurse(subLookupIndex uint16) bool {
	if c.nestingLevelLeft == 0 || c.recurseFunc == nil || c.buffer.maxOps <= 0 {
		if c.buffer.maxOps <= 0 {
			c.buffer.maxOps--
			return false
		}
		c.buffer.maxOps--
	}

	c.nestingLevelLeft--
	ret := c.recurseFunc(c, subLookupIndex)
	c.nestingLevelLeft++
	return ret
}

// `count` and `matchPositions` include the first glyph
// `lookupRecord` is in design order
func (c *otApplyContext) applyLookup(count int, matchPositions *[maxContextLength]int,
	lookupRecord []tt.SequenceLookup, matchLength int) {
	buffer := c.buffer
	var end int

	/* All positions are distance from beginning of *output* buffer.
	* Adjust. */
	{
		bl := buffer.backtrackLen()
		end = bl + matchLength

		delta := bl - buffer.idx
		/* Convert positions to new indexing. */
		for j := 0; j < count; j++ {
			matchPositions[j] += delta
		}
	}

	for _, lk := range lookupRecord {
		idx := int(lk.InputIndex)
		if idx >= count { // invalid, ignored
			continue
		}

		/* Don't recurse to ourself at same position.
		 * Note that this test is too naive, it doesn't catch longer loops. */
		if idx == 0 && lk.LookupIndex == c.lookupIndex {
			continue
		}

		buffer.moveTo(matchPositions[idx])

		if buffer.maxOps <= 0 {
			break
		}

		origLen := buffer.backtrackLen() + buffer.lookaheadLen()

		if debugMode >= 2 {
			fmt.Printf("\t\tAPPLY nested lookup %d\n", lk.LookupIndex)
		}

		if !c.recurse(lk.LookupIndex) {
			continue
		}

		newLen := buffer.backtrackLen() + buffer.lookaheadLen()
		delta := newLen - origLen

		if delta == 0 {
			continue
		}

		/* Recursed lookup changed buffer len. Adjust.
		 *
		 * TODO:
		 *
		 * Right now, if buffer length increased by n, we assume n new glyphs
		 * were added right after the current position, and if buffer length
		 * was decreased by n, we assume n match positions after the current
		 * one where removed.  The former (buffer length increased) case is
		 * fine, but the decrease case can be improved in at least two ways,
		 * both of which are significant:
		 *
		 *   - If recursed-to lookup is MultipleSubst and buffer length
		 *     decreased, then it's current match position that was deleted,
		 *     NOT the one after it.
		 *
		 *   - If buffer length was decreased by n, it does not necessarily
		 *     mean that n match positions where removed, as there might
		 *     have been marks and default-ignorables in the sequence.  We
		 *     should instead drop match positions between current-position
		 *     and current-position + n instead.
		 *
		 * It should be possible to construct tests for both of these cases.
		 */

		end += delta
		if end <= int(matchPositions[idx]) {
			/* End might end up being smaller than matchPositions[idx] if the recursed
			* lookup ended up removing many items, more than we have had matched.
			* Just never rewind end back and get out of here.
			* https://bugs.chromium.org/p/chromium/issues/detail?id=659496 */
			end = matchPositions[idx]
			/* There can't be any further changes. */
			break
		}

		next := idx + 1 /* next now is the position after the recursed lookup. */

		if delta > 0 {
			if delta+count > maxContextLength {
				break
			}
		} else {
			/* NOTE: delta is negative. */
			delta = max(delta, int(next)-int(count))
			next -= delta
		}

		/* Shift! */
		copy(matchPositions[next+delta:], matchPositions[next:count])
		next += delta
		count += delta

		/* Fill in new entries. */
		for j := idx + 1; j < next; j++ {
			matchPositions[j] = matchPositions[j-1] + 1
		}

		/* And fixup the rest. */
		for ; next < count; next++ {
			matchPositions[next] += delta
		}

	}

	buffer.moveTo(end)
}

func (c *otApplyContext) matchBacktrack(backtrack []uint16, matchFunc matcherFunc) (bool, int) {
	skippyIter := &c.iterContext
	skippyIter.reset(c.buffer.backtrackLen(), len(backtrack))
	skippyIter.setMatchFunc(matchFunc, backtrack)

	for i := 0; i < len(backtrack); i++ {
		if !skippyIter.prev() {
			return false, 0
		}
	}

	return true, skippyIter.idx
}

func (c *otApplyContext) matchLookahead(lookahead []uint16, matchFunc matcherFunc, offset int) (bool, int) {
	skippyIter := &c.iterContext
	skippyIter.reset(c.buffer.idx+offset-1, len(lookahead))
	skippyIter.setMatchFunc(matchFunc, lookahead)

	for i := 0; i < len(lookahead); i++ {
		if !skippyIter.next() {
			return false, 0
		}
	}

	return true, skippyIter.idx + 1
}

func (c *otApplyContext) applyLookupContext1(data tt.LookupContext1, index int) bool {
	if index >= len(data) { // index is not sanitized in tt.Parse
		return false
	}
	ruleSet := data[index]
	return c.applyRuleSet(ruleSet, matchGlyph)
}

func (c *otApplyContext) applyLookupContext2(data tt.LookupContext2, index int, glyphID fonts.GID) bool {
	class, _ := data.Class.ClassID(glyphID)
	ruleSet := data.SequenceSets[class]
	return c.applyRuleSet(ruleSet, matchClass(data.Class))
}

// return a slice containing [start, start+1, ..., end-1],
// using `indices` as an internal buffer to avoid allocations
// these indices are used to refer to coverage
func get1N(indices *[]uint16, start, end int) []uint16 {
	if end > cap(*indices) {
		*indices = make([]uint16, end)
		for i := range *indices {
			(*indices)[i] = uint16(i)
		}
	}
	return (*indices)[start:end]
}

func (c *otApplyContext) applyLookupContext3(data tt.LookupContext3, index int) bool {
	covIndices := get1N(&c.indices, 1, len(data.Coverages))
	return c.contextApplyLookup(covIndices, data.SequenceLookups, matchCoverage(data.Coverages))
}

func (c *otApplyContext) applyLookupChainedContext1(data tt.LookupChainedContext1, index int) bool {
	if index >= len(data) { // index is not sanitized in tt.Parse
		return false
	}
	ruleSet := data[index]
	return c.applyChainRuleSet(ruleSet, [3]matcherFunc{matchGlyph, matchGlyph, matchGlyph})
}

func (c *otApplyContext) applyLookupChainedContext2(data tt.LookupChainedContext2, index int, glyphID fonts.GID) bool {
	class, _ := data.InputClass.ClassID(glyphID)
	ruleSet := data.SequenceSets[class]
	return c.applyChainRuleSet(ruleSet, [3]matcherFunc{
		matchClass(data.BacktrackClass), matchClass(data.InputClass), matchClass(data.LookaheadClass),
	})
}

func (c *otApplyContext) applyLookupChainedContext3(data tt.LookupChainedContext3, index int) bool {
	lB, lI, lL := len(data.Backtrack), len(data.Input), len(data.Lookahead)
	return c.chainContextApplyLookup(get1N(&c.indices, 0, lB), get1N(&c.indices, 1, lI), get1N(&c.indices, 0, lL),
		data.SequenceLookups, [3]matcherFunc{
			matchCoverage(data.Backtrack), matchCoverage(data.Input), matchCoverage(data.Lookahead),
		})
}
