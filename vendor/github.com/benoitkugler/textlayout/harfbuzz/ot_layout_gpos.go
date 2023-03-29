package harfbuzz

import (
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
	tt "github.com/benoitkugler/textlayout/fonts/truetype"
)

// ported from harfbuzz/src/hb-ot-layout-gpos-table.hh Copyright © 2007,2008,2009,2010  Red Hat, Inc.; 2010,2012,2013  Google, Inc.  Behdad Esfahbod

// attach_type_t
const (
	attachTypeNone = 0x00

	/* Each attachment should be either a mark or a cursive; can't be both. */
	attachTypeMark    = 0x01
	attachTypeCursive = 0x02
)

func positionStartGPOS(buffer *Buffer) {
	for i := range buffer.Pos {
		buffer.Pos[i].attachChain = 0
		buffer.Pos[i].attachType = 0
	}
}

func propagateAttachmentOffsets(pos []GlyphPosition, i int, direction Direction) {
	/* Adjusts offsets of attached glyphs (both cursive and mark) to accumulate
	 * offset of glyph they are attached to. */
	chain, type_ := pos[i].attachChain, pos[i].attachType
	if chain == 0 {
		return
	}

	pos[i].attachChain = 0

	j := i + int(chain)

	if j >= len(pos) {
		return
	}

	propagateAttachmentOffsets(pos, j, direction)

	//   assert (!!(type_ & attachTypeMark) ^ !!(type_ & attachTypeCursive));

	if (type_ & attachTypeCursive) != 0 {
		if direction.isHorizontal() {
			pos[i].YOffset += pos[j].YOffset
		} else {
			pos[i].XOffset += pos[j].XOffset
		}
	} else /*if (type_ & attachTypeMark)*/ {
		pos[i].XOffset += pos[j].XOffset
		pos[i].YOffset += pos[j].YOffset

		// assert (j < i);
		if direction.isForward() {
			for _, p := range pos[j:i] {
				pos[i].XOffset -= p.XAdvance
				pos[i].YOffset -= p.YAdvance
			}
		} else {
			for _, p := range pos[j+1 : i+1] {
				pos[i].XOffset += p.XAdvance
				pos[i].YOffset += p.YAdvance
			}
		}
	}
}

func positionFinishOffsetsGPOS(buffer *Buffer) {
	pos := buffer.Pos
	direction := buffer.Props.Direction

	/* Handle attachments */
	if buffer.scratchFlags&bsfHasGPOSAttachment != 0 {

		if debugMode >= 2 {
			fmt.Println("POSITION - handling attachments")
		}

		for i := range pos {
			propagateAttachmentOffsets(pos, i, direction)
		}
	}
}

var _ layoutLookup = lookupGPOS{}

// implements layoutLookup
type lookupGPOS tt.LookupGPOS

func (l lookupGPOS) collectCoverage(dst *setDigest) {
	for _, table := range l.Subtables {
		dst.collectCoverage(table.Coverage)
	}
}

func (l lookupGPOS) dispatchSubtables(ctx *getSubtablesContext) {
	for _, table := range l.Subtables {
		*ctx = append(*ctx, newGPOSApplicable(table))
	}
}

func (l lookupGPOS) dispatchApply(ctx *otApplyContext) bool {
	for _, table := range l.Subtables {
		if gposSubtable(table).apply(ctx) {
			return true
		}
	}
	return false
}

func (lookupGPOS) isReverse() bool { return false }

func applyRecurseGPOS(c *otApplyContext, lookupIndex uint16) bool {
	gpos := c.font.otTables.GPOS
	l := lookupGPOS(gpos.Lookups[lookupIndex])
	return c.applyRecurseLookup(lookupIndex, l)
}

//  implements `hb_apply_func_t`
type gposSubtable tt.GPOSSubtable

// return `true` is the positionning found a match and was applied
func (table gposSubtable) apply(c *otApplyContext) bool {
	buffer := c.buffer
	glyphID := buffer.cur(0).Glyph
	glyphPos := buffer.curPos(0)
	index, ok := table.Coverage.Index(glyphID)
	if !ok {
		return false
	}

	if debugMode >= 2 {
		fmt.Printf("\tAPPLY - type %T at index %d\n", table.Data, c.buffer.idx)
	}

	switch data := table.Data.(type) {
	case tt.GPOSSingle1:
		c.applyGPOSValueRecord(data.Format, data.Value, glyphPos)
		buffer.idx++
	case tt.GPOSSingle2:
		c.applyGPOSValueRecord(data.Format, data.Values[index], glyphPos)
		buffer.idx++
	case tt.GPOSPair1:
		skippyIter := &c.iterInput
		skippyIter.reset(buffer.idx, 1)
		if !skippyIter.next() {
			return false
		}
		set := data.Values[index]
		record := set.FindGlyph(buffer.Info[skippyIter.idx].Glyph)
		if record == nil {
			return false
		}
		c.applyGPOSPair(data.Formats, record.Pos, skippyIter.idx)
	case tt.GPOSPair2:
		skippyIter := &c.iterInput
		skippyIter.reset(buffer.idx, 1)
		if !skippyIter.next() {
			return false
		}
		class1, _ := data.First.ClassID(glyphID)
		class2, _ := data.Second.ClassID(buffer.Info[skippyIter.idx].Glyph)
		vals := data.Values[class1][class2]
		c.applyGPOSPair(data.Formats, vals, skippyIter.idx)
	case tt.GPOSCursive1:
		return c.applyGPOSCursive(data, index, table.Coverage)
	case tt.GPOSMarkToBase1:
		return c.applyGPOSMarkToBase(data, index)
	case tt.GPOSMarkToLigature1:
		return c.applyGPOSMarkToLigature(data, index)
	case tt.GPOSMarkToMark1:
		return c.applyGPOSMarkToMark(data, index)
	case tt.GPOSContext1:
		return c.applyLookupContext1(tt.LookupContext1(data), index)
	case tt.GPOSContext2:
		return c.applyLookupContext2(tt.LookupContext2(data), index, glyphID)
	case tt.GPOSContext3:
		return c.applyLookupContext3(tt.LookupContext3(data), index)
	case tt.GPOSChainedContext1:
		return c.applyLookupChainedContext1(tt.LookupChainedContext1(data), index)
	case tt.GPOSChainedContext2:
		return c.applyLookupChainedContext2(tt.LookupChainedContext2(data), index, glyphID)
	case tt.GPOSChainedContext3:
		return c.applyLookupChainedContext3(tt.LookupChainedContext3(data), index)
	}
	return true
}

func (c *otApplyContext) applyGPOSValueRecord(format tt.GPOSValueFormat, v tt.GPOSValueRecord, glyphPos *GlyphPosition) bool {
	var ret bool
	if format == 0 {
		return ret
	}

	font := c.font
	horizontal := c.direction.isHorizontal()

	if format&tt.XPlacement != 0 {
		glyphPos.XOffset += font.emScaleX(v.XPlacement)
		ret = ret || v.XPlacement != 0
	}
	if format&tt.YPlacement != 0 {
		glyphPos.YOffset += font.emScaleY(v.YPlacement)
		ret = ret || v.YPlacement != 0
	}
	if format&tt.XAdvance != 0 {
		if horizontal {
			glyphPos.XAdvance += font.emScaleX(v.XAdvance)
			ret = ret || v.XAdvance != 0
		}
	}
	/* YAdvance values grow downward but font-space grows upward, hence negation */
	if format&tt.YAdvance != 0 {
		if !horizontal {
			glyphPos.YAdvance -= font.emScaleY(v.YAdvance)
			ret = ret || v.YAdvance != 0
		}
	}

	if format&tt.Devices == 0 {
		return ret
	}

	useXDevice := font.XPpem != 0 || len(font.varCoords()) != 0
	useYDevice := font.YPpem != 0 || len(font.varCoords()) != 0

	if !useXDevice && !useYDevice {
		return ret
	}

	if format&tt.XPlaDevice != 0 && useXDevice {
		glyphPos.XOffset += font.getXDelta(c.varStore, v.XPlaDevice)
		ret = ret || v.XPlaDevice != nil
	}
	if format&tt.YPlaDevice != 0 && useYDevice {
		glyphPos.YOffset += font.getYDelta(c.varStore, v.YPlaDevice)
		ret = ret || v.YPlaDevice != nil
	}
	if format&tt.XAdvDevice != 0 && horizontal && useXDevice {
		glyphPos.XAdvance += font.getXDelta(c.varStore, v.XAdvDevice)
		ret = ret || v.XAdvDevice != nil
	}
	if format&tt.YAdvDevice != 0 && !horizontal && useYDevice {
		/* YAdvance values grow downward but font-space grows upward, hence negation */
		glyphPos.YAdvance -= font.getYDelta(c.varStore, v.YAdvDevice)
		ret = ret || v.YAdvDevice != nil
	}
	return ret
}

func reverseCursiveMinorOffset(pos []GlyphPosition, i int, direction Direction, newParent int) {
	chain, type_ := pos[i].attachChain, pos[i].attachType
	if chain == 0 || type_&attachTypeCursive == 0 {
		return
	}

	pos[i].attachChain = 0

	j := i + int(chain)

	// stop if we see new parent in the chain
	if j == newParent {
		return
	}
	reverseCursiveMinorOffset(pos, j, direction, newParent)

	if direction.isHorizontal() {
		pos[j].YOffset = -pos[i].YOffset
	} else {
		pos[j].XOffset = -pos[i].XOffset
	}

	pos[j].attachChain = -chain
	pos[j].attachType = type_
}

func (c *otApplyContext) applyGPOSPair(formats [2]tt.GPOSValueFormat, values [2]tt.GPOSValueRecord, pos int) {
	buffer := c.buffer

	ap1 := c.applyGPOSValueRecord(formats[0], values[0], buffer.curPos(0))
	ap2 := c.applyGPOSValueRecord(formats[1], values[1], &buffer.Pos[pos])

	if ap1 || ap2 {
		buffer.unsafeToBreak(buffer.idx, pos+1)
	}
	buffer.idx = pos
	if formats[1] != 0 {
		buffer.idx++
	}
}

func (c *otApplyContext) applyGPOSCursive(data tt.GPOSCursive1, covIndex int, cov tt.Coverage) bool {
	buffer := c.buffer

	thisRecord := data[covIndex]
	if thisRecord[0] == nil {
		return false
	}

	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	if !skippyIter.prev() {
		return false
	}

	prevIndex, ok := cov.Index(buffer.Info[skippyIter.idx].Glyph)
	if !ok {
		return false
	}
	prevRecord := data[prevIndex]
	if prevRecord[1] == nil {
		return false
	}

	i := skippyIter.idx
	j := buffer.idx

	buffer.unsafeToBreak(i, j)
	exitX, exitY := c.getAnchor(prevRecord[1], buffer.Info[i].Glyph)
	entryX, entryY := c.getAnchor(thisRecord[0], buffer.Info[j].Glyph)

	pos := buffer.Pos

	var d Position
	/* Main-direction adjustment */
	switch c.direction {
	case LeftToRight:
		pos[i].XAdvance = roundf(exitX) + pos[i].XOffset

		d = roundf(entryX) + pos[j].XOffset
		pos[j].XAdvance -= d
		pos[j].XOffset -= d
	case RightToLeft:
		d = roundf(exitX) + pos[i].XOffset
		pos[i].XAdvance -= d
		pos[i].XOffset -= d

		pos[j].XAdvance = roundf(entryX) + pos[j].XOffset
	case TopToBottom:
		pos[i].YAdvance = roundf(exitY) + pos[i].YOffset

		d = roundf(entryY) + pos[j].YOffset
		pos[j].YAdvance -= d
		pos[j].YOffset -= d
	case BottomToTop:
		d = roundf(exitY) + pos[i].YOffset
		pos[i].YAdvance -= d
		pos[i].YOffset -= d

		pos[j].YAdvance = roundf(entryY)
	}

	/* Cross-direction adjustment */

	/* We attach child to parent (think graph theory and rooted trees whereas
	 * the root stays on baseline and each node aligns itself against its
	 * parent.
	 *
	 * Optimize things for the case of RightToLeft, as that's most common in
	 * Arabic. */
	child := i
	parent := j
	xOffset := Position(entryX - exitX)
	yOffset := Position(entryY - exitY)
	if uint16(c.lookupProps)&tt.RightToLeft == 0 {
		k := child
		child = parent
		parent = k
		xOffset = -xOffset
		yOffset = -yOffset
	}

	/* If child was already connected to someone else, walk through its old
	 * chain and reverse the link direction, such that the whole tree of its
	 * previous connection now attaches to new parent.  Watch out for case
	 * where new parent is on the path from old chain...
	 */
	reverseCursiveMinorOffset(pos, child, c.direction, parent)

	pos[child].attachType = attachTypeCursive
	pos[child].attachChain = int16(parent - child)
	buffer.scratchFlags |= bsfHasGPOSAttachment
	if c.direction.isHorizontal() {
		pos[child].YOffset = yOffset
	} else {
		pos[child].XOffset = xOffset
	}

	/* If parent was attached to child, break them free.
	 * https://github.com/harfbuzz/harfbuzz/issues/2469 */
	if pos[parent].attachChain == -pos[child].attachChain {
		pos[parent].attachChain = 0
	}

	buffer.idx++
	return true
}

// panic if anchor is nil
func (c *otApplyContext) getAnchor(anchor tt.GPOSAnchor, glyph fonts.GID) (x, y float32) {
	font := c.font
	switch anchor := anchor.(type) {
	case tt.GPOSAnchorFormat1:
		return font.emFscaleX(anchor.X), font.emFscaleY(anchor.Y)
	case tt.GPOSAnchorFormat2:
		xPpem, yPpem := font.XPpem, font.YPpem
		var cx, cy Position
		ret := xPpem != 0 || yPpem != 0
		if ret {
			cx, cy, ret = font.getGlyphContourPointForOrigin(glyph, anchor.AnchorPoint, LeftToRight)
		}
		if ret && xPpem != 0 {
			x = float32(cx)
		} else {
			x = font.emFscaleX(anchor.X)
		}
		if ret && yPpem != 0 {
			y = float32(cy)
		} else {
			y = font.emFscaleY(anchor.Y)
		}
		return x, y
	case tt.GPOSAnchorFormat3:
		x, y = font.emFscaleX(anchor.X), font.emFscaleY(anchor.Y)
		if font.XPpem != 0 || len(font.varCoords()) != 0 {
			x += float32(font.getXDelta(c.varStore, anchor.XDevice))
		}
		if font.YPpem != 0 || len(font.varCoords()) != 0 {
			y += float32(font.getYDelta(c.varStore, anchor.YDevice))
		}
		return x, y
	default:
		panic("exhaustive switch")
	}
}

func (c *otApplyContext) applyGPOSMarks(marks []tt.GPOSMark, markIndex, glyphIndex int, anchors [][]tt.GPOSAnchor, glyphPos int) bool {
	buffer := c.buffer
	record := &marks[markIndex]
	markClass := record.ClassValue
	markAnchor := record.Anchor

	glyphAnchor := anchors[glyphIndex][markClass]
	/* If this subtable doesn't have an anchor for this base and this class,
	 * return false such that the subsequent subtables have a chance at it. */
	if glyphAnchor == nil {
		return false
	}

	buffer.unsafeToBreak(glyphPos, buffer.idx)
	markX, markY := c.getAnchor(markAnchor, buffer.cur(0).Glyph)
	baseX, baseY := c.getAnchor(glyphAnchor, buffer.Info[glyphPos].Glyph)

	o := buffer.curPos(0)
	o.XOffset = roundf(baseX - markX)
	o.YOffset = roundf(baseY - markY)
	o.attachType = attachTypeMark
	o.attachChain = int16(glyphPos - buffer.idx)
	buffer.scratchFlags |= bsfHasGPOSAttachment

	buffer.idx++
	return true
}

func (c *otApplyContext) applyGPOSMarkToBase(data tt.GPOSMarkToBase1, markIndex int) bool {
	buffer := c.buffer

	// now we search backwards for a non-mark glyph
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	skippyIter.matcher.lookupProps = uint32(tt.IgnoreMarks)
	for {
		if !skippyIter.prev() {
			return false
		}
		/* We only want to attach to the first of a MultipleSubst sequence.
		 * https://github.com/harfbuzz/harfbuzz/issues/740
		 * Reject others...
		 * ...but stop if we find a mark in the MultipleSubst sequence:
		 * https://github.com/harfbuzz/harfbuzz/issues/1020 */
		if !buffer.Info[skippyIter.idx].multiplied() || buffer.Info[skippyIter.idx].getLigComp() == 0 ||
			skippyIter.idx == 0 || buffer.Info[skippyIter.idx-1].isMark() ||
			buffer.Info[skippyIter.idx].getLigID() != buffer.Info[skippyIter.idx-1].getLigID() ||
			buffer.Info[skippyIter.idx].getLigComp() != buffer.Info[skippyIter.idx-1].getLigComp()+1 {
			break
		}
		skippyIter.reject()
	}

	/* Checking that matched glyph is actually a base glyph by GDEF is too strong; disabled */
	//if (!_hb_glyph_info_is_base_glyph (&buffer.Info[skippyIter.idx])) { return false; }

	baseIndex, ok := data.BaseCoverage.Index(buffer.Info[skippyIter.idx].Glyph)
	if !ok {
		return false
	}

	return c.applyGPOSMarks(data.Marks, markIndex, baseIndex, data.Bases, skippyIter.idx)
}

func (c *otApplyContext) applyGPOSMarkToLigature(data tt.GPOSMarkToLigature1, markIndex int) bool {
	buffer := c.buffer

	// now we search backwards for a non-mark glyph
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	skippyIter.matcher.lookupProps = uint32(tt.IgnoreMarks)
	if !skippyIter.prev() {
		return false
	}

	j := skippyIter.idx
	ligIndex, ok := data.LigatureCoverage.Index(buffer.Info[j].Glyph)
	if !ok {
		return false
	}

	ligAttach := data.Ligatures[ligIndex]

	/* Find component to attach to */
	compCount := len(ligAttach)
	if compCount == 0 {
		return false
	}

	/* We must now check whether the ligature ID of the current mark glyph
	 * is identical to the ligature ID of the found ligature.  If yes, we
	 * can directly use the component index.  If not, we attach the mark
	 * glyph to the last component of the ligature. */
	ligID := buffer.Info[j].getLigID()
	markID := buffer.cur(0).getLigID()
	markComp := buffer.cur(0).getLigComp()
	compIndex := compCount - 1
	if ligID != 0 && ligID == markID && markComp > 0 {
		compIndex = min(compCount, int(buffer.cur(0).getLigComp())) - 1
	}

	return c.applyGPOSMarks(data.Marks, markIndex, compIndex, ligAttach, skippyIter.idx)
}

func (c *otApplyContext) applyGPOSMarkToMark(data tt.GPOSMarkToMark1, mark1Index int) bool {
	buffer := c.buffer

	// now we search backwards for a suitable mark glyph until a non-mark glyph
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	skippyIter.matcher.lookupProps = c.lookupProps &^ uint32(ignoreFlags)
	if !skippyIter.prev() {
		return false
	}

	if !buffer.Info[skippyIter.idx].isMark() {
		return false
	}

	j := skippyIter.idx

	id1 := buffer.cur(0).getLigID()
	id2 := buffer.Info[j].getLigID()
	comp1 := buffer.cur(0).getLigComp()
	comp2 := buffer.Info[j].getLigComp()

	if id1 == id2 {
		if id1 == 0 { /* Marks belonging to the same base. */
			goto good
		} else if comp1 == comp2 { /* Marks belonging to the same ligature component. */
			goto good
		}
	} else {
		/* If ligature ids don't match, it may be the case that one of the marks
		* itself is a ligature.  In which case match. */
		if (id1 > 0 && comp1 == 0) || (id2 > 0 && comp2 == 0) {
			goto good
		}
	}

	/* Didn't match. */
	return false

good:
	mark2Index, ok := data.Mark2Coverage.Index(buffer.Info[j].Glyph)
	if !ok {
		return false
	}

	return c.applyGPOSMarks(data.Marks1, mark1Index, mark2Index, data.Marks2, j)
}
