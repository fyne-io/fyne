package harfbuzz

import (
	"fmt"

	"github.com/go-text/typesetting/opentype/api/font"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-layout-gpos-table.hh Copyright © 2007,2008,2009,2010  Red Hat, Inc.; 2010,2012,2013  Google, Inc.  Behdad Esfahbod

var _ layoutLookup = lookupGPOS{}

// implements layoutLookup
type lookupGPOS font.GPOSLookup

func (l lookupGPOS) Props() uint32 { return l.LookupOptions.Props() }

func (l lookupGPOS) collectCoverage(dst *setDigest) {
	for _, table := range l.Subtables {
		dst.collectCoverage(table.Cov())
	}
}

func (l lookupGPOS) dispatchSubtables(ctx *getSubtablesContext) {
	for _, table := range l.Subtables {
		*ctx = append(*ctx, newGPOSApplicable(table))
	}
}

func (l lookupGPOS) dispatchApply(ctx *otApplyContext) bool {
	for _, table := range l.Subtables {
		if ctx.applyGPOS(table) {
			return true
		}
	}
	return false
}

func (lookupGPOS) isReverse() bool { return false }

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

func applyRecurseGPOS(c *otApplyContext, lookupIndex uint16) bool {
	gpos := c.font.face.GPOS
	l := lookupGPOS(gpos.Lookups[lookupIndex])
	return c.applyRecurseLookup(lookupIndex, l)
}

// return `true` is the positionning found a match and was applied
func (c *otApplyContext) applyGPOS(table tables.GPOSLookup) bool {
	buffer := c.buffer
	glyphID := buffer.cur(0).Glyph
	glyphPos := buffer.curPos(0)
	index, ok := table.Cov().Index(gID(glyphID))
	if !ok {
		return false
	}

	if debugMode >= 2 {
		fmt.Printf("\tAPPLY - type %T at index %d\n", table, c.buffer.idx)
	}

	switch data := table.(type) {
	case tables.SinglePos:
		switch inner := data.Data.(type) {
		case tables.SinglePosData1:
			c.applyGPOSValueRecord(inner.ValueFormat, inner.ValueRecord, glyphPos)
		case tables.SinglePosData2:
			c.applyGPOSValueRecord(inner.ValueFormat, inner.ValueRecords[index], glyphPos)
		}
		buffer.idx++
	case tables.PairPos:
		skippyIter := &c.iterInput
		skippyIter.reset(buffer.idx, 1)
		if !skippyIter.next() {
			return false
		}
		switch inner := data.Data.(type) {
		case tables.PairPosData1:
			set := inner.PairSets[index]
			record := set.FindGlyph(gID(buffer.Info[skippyIter.idx].Glyph))
			if record == nil {
				return false
			}
			c.applyGPOSPair(inner.ValueFormat1, inner.ValueFormat2, record.ValueRecord1, record.ValueRecord2, skippyIter.idx)
		case tables.PairPosData2:
			class1, _ := inner.ClassDef1.Class(gID(glyphID))
			class2, _ := inner.ClassDef2.Class(gID(buffer.Info[skippyIter.idx].Glyph))
			vals := inner.Class1Records[class1][class2]
			c.applyGPOSPair(inner.ValueFormat1, inner.ValueFormat2, vals.ValueRecord1, vals.ValueRecord2, skippyIter.idx)
		}

	case tables.CursivePos:
		return c.applyGPOSCursive(data, index)
	case tables.MarkBasePos:
		return c.applyGPOSMarkToBase(data, index)
	case tables.MarkLigPos:
		return c.applyGPOSMarkToLigature(data, index)
	case tables.MarkMarkPos:
		return c.applyGPOSMarkToMark(data, index)

	case tables.ContextualPos:
		switch inner := data.Data.(type) {
		case tables.ContextualPos1:
			return c.applyLookupContext1(tables.SequenceContextFormat1(inner), index)
		case tables.ContextualPos2:
			return c.applyLookupContext2(tables.SequenceContextFormat2(inner), index, glyphID)
		case tables.ContextualPos3:
			return c.applyLookupContext3(tables.SequenceContextFormat3(inner), index)
		}

	case tables.ChainedContextualPos:
		switch inner := data.Data.(type) {
		case tables.ChainedContextualPos1:
			return c.applyLookupChainedContext1(tables.ChainedSequenceContextFormat1(inner), index)
		case tables.ChainedContextualPos2:
			return c.applyLookupChainedContext2(tables.ChainedSequenceContextFormat2(inner), index, glyphID)
		case tables.ChainedContextualPos3:
			return c.applyLookupChainedContext3(tables.ChainedSequenceContextFormat3(inner), index)
		}
	}
	return true
}

func (c *otApplyContext) applyGPOSValueRecord(format tables.ValueFormat, v tables.ValueRecord, glyphPos *GlyphPosition) bool {
	var ret bool
	if format == 0 {
		return ret
	}

	font := c.font
	horizontal := c.direction.isHorizontal()

	if format&tables.XPlacement != 0 {
		glyphPos.XOffset += font.emScaleX(v.XPlacement)
		ret = ret || v.XPlacement != 0
	}
	if format&tables.YPlacement != 0 {
		glyphPos.YOffset += font.emScaleY(v.YPlacement)
		ret = ret || v.YPlacement != 0
	}
	if format&tables.XAdvance != 0 {
		if horizontal {
			glyphPos.XAdvance += font.emScaleX(v.XAdvance)
			ret = ret || v.XAdvance != 0
		}
	}
	/* YAdvance values grow downward but font-space grows upward, hence negation */
	if format&tables.YAdvance != 0 {
		if !horizontal {
			glyphPos.YAdvance -= font.emScaleY(v.YAdvance)
			ret = ret || v.YAdvance != 0
		}
	}

	if format&tables.Devices == 0 {
		return ret
	}

	useXDevice := font.face.XPpem != 0 || len(font.varCoords()) != 0
	useYDevice := font.face.YPpem != 0 || len(font.varCoords()) != 0

	if !useXDevice && !useYDevice {
		return ret
	}

	if format&tables.XPlaDevice != 0 && useXDevice {
		glyphPos.XOffset += font.getXDelta(c.varStore, v.XPlaDevice)
		ret = ret || v.XPlaDevice != nil
	}
	if format&tables.YPlaDevice != 0 && useYDevice {
		glyphPos.YOffset += font.getYDelta(c.varStore, v.YPlaDevice)
		ret = ret || v.YPlaDevice != nil
	}
	if format&tables.XAdvDevice != 0 && horizontal && useXDevice {
		glyphPos.XAdvance += font.getXDelta(c.varStore, v.XAdvDevice)
		ret = ret || v.XAdvDevice != nil
	}
	if format&tables.YAdvDevice != 0 && !horizontal && useYDevice {
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

func (c *otApplyContext) applyGPOSPair(f1, f2 tables.ValueFormat, v1, v2 tables.ValueRecord, pos int) {
	buffer := c.buffer

	ap1 := c.applyGPOSValueRecord(f1, v1, buffer.curPos(0))
	ap2 := c.applyGPOSValueRecord(f2, v2, &buffer.Pos[pos])

	if ap1 || ap2 {
		buffer.unsafeToBreak(buffer.idx, pos+1)
	}
	buffer.idx = pos
	if f2 != 0 {
		buffer.idx++
	}
}

func (c *otApplyContext) applyGPOSCursive(data tables.CursivePos, covIndex int) bool {
	buffer := c.buffer

	thisRecord := data.EntryExits[covIndex]
	if thisRecord.EntryAnchor == nil {
		return false
	}

	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	if !skippyIter.prev() {
		return false
	}

	prevIndex, ok := data.Cov().Index(gID(buffer.Info[skippyIter.idx].Glyph))
	if !ok {
		return false
	}
	prevRecord := data.EntryExits[prevIndex]
	if prevRecord.ExitAnchor == nil {
		return false
	}

	i := skippyIter.idx
	j := buffer.idx

	buffer.unsafeToBreak(i, j)
	exitX, exitY := c.getAnchor(prevRecord.ExitAnchor, buffer.Info[i].Glyph)
	entryX, entryY := c.getAnchor(thisRecord.EntryAnchor, buffer.Info[j].Glyph)

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
	if uint16(c.lookupProps)&otRightToLeft == 0 {
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
func (c *otApplyContext) getAnchor(anchor tables.Anchor, glyph GID) (x, y float32) {
	font := c.font
	switch anchor := anchor.(type) {
	case tables.AnchorFormat1:
		return font.emFscaleX(anchor.XCoordinate), font.emFscaleY(anchor.YCoordinate)
	case tables.AnchorFormat2:
		xPpem, yPpem := font.face.XPpem, font.face.YPpem
		var cx, cy Position
		ret := xPpem != 0 || yPpem != 0
		if ret {
			cx, cy, ret = font.getGlyphContourPointForOrigin(glyph, anchor.AnchorPoint, LeftToRight)
		}
		if ret && xPpem != 0 {
			x = float32(cx)
		} else {
			x = font.emFscaleX(anchor.XCoordinate)
		}
		if ret && yPpem != 0 {
			y = float32(cy)
		} else {
			y = font.emFscaleY(anchor.YCoordinate)
		}
		return x, y
	case tables.AnchorFormat3:
		x, y = font.emFscaleX(anchor.XCoordinate), font.emFscaleY(anchor.YCoordinate)
		if font.face.XPpem != 0 || len(font.varCoords()) != 0 {
			x += float32(font.getXDelta(c.varStore, anchor.XDevice))
		}
		if font.face.YPpem != 0 || len(font.varCoords()) != 0 {
			y += float32(font.getYDelta(c.varStore, anchor.YDevice))
		}
		return x, y
	default:
		panic("exhaustive switch")
	}
}

func (c *otApplyContext) applyGPOSMarks(marks tables.MarkArray, markIndex, glyphIndex int, anchors [][]tables.Anchor, glyphPos int) bool {
	buffer := c.buffer
	markClass := marks.MarkRecords[markIndex].MarkClass
	markAnchor := marks.MarkAnchors[markIndex]

	glyphAnchor := anchors[glyphIndex][markClass]
	// If this subtable doesn't have an anchor for this base and this class,
	// return false such that the subsequent subtables have a chance at it.
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

func (c *otApplyContext) applyGPOSMarkToBase(data tables.MarkBasePos, markIndex int) bool {
	buffer := c.buffer

	// now we search backwards for a non-mark glyph
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	skippyIter.matcher.lookupProps = uint32(otIgnoreMarks)
	for {
		if !skippyIter.prev() {
			return false
		}
		// We only want to attach to the first of a MultipleSubst sequence.
		// https://github.com/harfbuzz/harfbuzz/issues/740
		// Reject others...
		// ...but stop if we find a mark in the MultipleSubst sequence:
		// https://github.com/harfbuzz/harfbuzz/issues/1020
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

	baseIndex, ok := data.BaseCoverage.Index(gID(buffer.Info[skippyIter.idx].Glyph))
	if !ok {
		return false
	}

	return c.applyGPOSMarks(data.MarkArray, markIndex, baseIndex, data.BaseArray.BaseAnchors, skippyIter.idx)
}

func (c *otApplyContext) applyGPOSMarkToLigature(data tables.MarkLigPos, markIndex int) bool {
	buffer := c.buffer

	// now we search backwards for a non-mark glyph
	skippyIter := &c.iterInput
	skippyIter.reset(buffer.idx, 1)
	skippyIter.matcher.lookupProps = uint32(otIgnoreMarks)
	if !skippyIter.prev() {
		return false
	}

	j := skippyIter.idx
	ligIndex, ok := data.LigatureCoverage.Index(gID(buffer.Info[j].Glyph))
	if !ok {
		return false
	}

	ligAttach := data.LigatureArray.LigatureAttachs[ligIndex].ComponentAnchors

	// Find component to attach to
	compCount := len(ligAttach)
	if compCount == 0 {
		return false
	}

	// We must now check whether the ligature ID of the current mark glyph
	// is identical to the ligature ID of the found ligature.  If yes, we
	// can directly use the component index.  If not, we attach the mark
	// glyph to the last component of the ligature.
	ligID := buffer.Info[j].getLigID()
	markID := buffer.cur(0).getLigID()
	markComp := buffer.cur(0).getLigComp()
	compIndex := compCount - 1
	if ligID != 0 && ligID == markID && markComp > 0 {
		compIndex = min(compCount, int(buffer.cur(0).getLigComp())) - 1
	}

	return c.applyGPOSMarks(data.MarkArray, markIndex, compIndex, ligAttach, skippyIter.idx)
}

func (c *otApplyContext) applyGPOSMarkToMark(data tables.MarkMarkPos, mark1Index int) bool {
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
	mark2Index, ok := data.Mark2Coverage.Index(gID(buffer.Info[j].Glyph))
	if !ok {
		return false
	}

	return c.applyGPOSMarks(data.Mark1Array, mark1Index, mark2Index, data.Mark2Array.Mark2Anchors, j)
}
