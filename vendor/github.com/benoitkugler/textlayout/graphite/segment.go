package graphite

const maxSegGrowthFactor = 64

type charInfo struct {
	before int // slot index before us, comes before
	after  int // slot index after us, comes after
	base   int // index into input string corresponding to this charinfo
	// featureIndex int  // index into features list in the segment −> Always 0
	char        rune  // Unicode character from character stream
	breakWeight int16 // breakweight coming from lb table
	flags       uint8 // 0,1 segment split.
}

func (ch *charInfo) addFlags(val uint8) { ch.flags |= val }

// Segment represents a line of text.
// It is used internally during shaping and
// returned as the result of the operation.
type Segment struct {
	// Start of the segment (may be nil for empty segments)
	First *Slot
	last  *Slot // last slot in the segment

	face *GraphiteFace
	silf *passes // selected subtable

	feats FeaturesValue // applied values

	// character info, one per input character
	charinfo []charInfo

	freeSlots  *Slot // linked list of free slots
	collisions []slotCollision

	// Advance of the whole segment
	Advance Position

	// Number of slots (output characters).
	// Since slots may be added or deleted during shaping,
	// it may differ from the number of characters ot the text input.
	// It could be directly computed by walking the linked list but is cached
	// for performance reasons.
	NumGlyphs int

	passBits uint32 // if bit set then skip pass
	flags    uint8  // General purpose flags
	dir      int8   // text direction

}

func (seg *Segment) currdir() bool { return ((seg.dir>>reverseBit)^seg.dir)&1 != 0 }

const (
	initCollisions = 1 + iota
	hasCollisions
)

func (seg *Segment) hasCollisionInfo() bool {
	return (seg.flags&hasCollisions) != 0 && seg.collisions != nil
}

func (seg *Segment) mergePassBits(val uint32) { seg.passBits &= val }

func (seg *Segment) processRunes(text []rune) {
	for slotID, r := range text {
		gid, _ := seg.face.cmap.Lookup(r)
		if gid == 0 {
			gid = seg.silf.findPdseudoGlyph(r)
		}
		seg.appendSlot(slotID, r, gid)
	}
	// the initial segment has one slot per input character
	seg.NumGlyphs = len(text)
}

func (seg *Segment) newSlot() *Slot {
	sl := new(Slot)
	sl.userAttrs = make([]int16, seg.silf.userAttibutes)
	return sl
}

func (seg *Segment) newJustify() *slotJustify {
	return new(slotJustify)
}

func (seg *Segment) appendSlot(index int, cid rune, gid GID) {
	sl := seg.newSlot()

	info := &seg.charinfo[index]
	info.char = cid
	// info.featureIndex = featureID
	info.base = index
	glyph := seg.face.getGlyph(gid)
	if glyph != nil {
		info.breakWeight = glyph.attrs.get(uint16(seg.silf.attrBreakWeight))
	}

	sl.setGlyph(seg, gid)
	sl.original, sl.Before, sl.After = index, index, index
	if seg.last != nil {
		seg.last.Next = sl
	}
	sl.prev = seg.last
	seg.last = sl
	if seg.First == nil {
		seg.First = sl
	}

	if aPassBits := uint16(seg.silf.attrSkipPasses); glyph != nil && aPassBits != 0 {
		m := uint32(glyph.attrs.get(aPassBits))
		if len(seg.silf.passes) > 16 {
			m |= uint32(glyph.attrs.get(aPassBits+1)) << 16
		}
		seg.mergePassBits(m)
	}
}

func (seg *Segment) freeSlot(aSlot *Slot) {
	if aSlot == nil {
		return
	}
	if seg.last == aSlot {
		seg.last = aSlot.prev
	}
	if seg.First == aSlot {
		seg.First = aSlot.Next
	}
	if aSlot.parent != nil {
		aSlot.parent.removeChild(aSlot)
	}
	for aSlot.child != nil {
		if aSlot.child.parent == aSlot {
			aSlot.child.parent = nil
			aSlot.removeChild(aSlot.child)
		} else {
			aSlot.child = nil
		}
	}

	// update next pointer
	aSlot.Next = seg.freeSlots
	seg.freeSlots = aSlot
}

const reverseBit = 6

// reverse the slots but keep diacritics in their same position after their bases
func (seg *Segment) reverseSlots() {
	seg.dir = seg.dir ^ 1<<reverseBit // invert the reverse flag
	if seg.First == seg.last {
		return
	} // skip 0 or 1 glyph runs

	var (
		curr                  = seg.First
		t, tlast, tfirst, out *Slot
	)

	for curr != nil && seg.getSlotBidiClass(curr) == 16 {
		curr = curr.Next
	}
	if curr == nil {
		return
	}
	tfirst = curr.prev
	tlast = curr

	for curr != nil {
		if seg.getSlotBidiClass(curr) == 16 {
			d := curr.Next
			for d != nil && seg.getSlotBidiClass(d) == 16 {
				d = d.Next
			}
			if d != nil {
				d = d.prev
			} else {
				d = seg.last
			}
			p := out.Next // one after the diacritics. out can't be null
			if p != nil {
				p.prev = d
			} else {
				tlast = d
			}
			t = d.Next
			d.Next = p
			curr.prev = out
			out.Next = curr
		} else { // will always fire first time round the loop
			if out != nil {
				out.prev = curr
			}
			t = curr.Next
			curr.Next = out
			out = curr
		}
		curr = t
	}
	out.prev = tfirst
	if tfirst != nil {
		tfirst.Next = out
	} else {
		seg.First = out
	}
	seg.last = tlast
}

func (seg *Segment) positionSlots(font *FontOptions, iStart, iEnd *Slot, isRTL, isFinal bool) Position {
	var (
		currpos    Position
		clusterMin float32
		bbox       rect
		reorder    = (seg.currdir() != isRTL)
	)

	if reorder {
		seg.reverseSlots()
		iStart, iEnd = iEnd, iStart
	}
	if iStart == nil {
		iStart = seg.First
	}
	if iEnd == nil {
		iEnd = seg.last
	}

	if iStart == nil || iEnd == nil { // only true for empty segments
		return currpos
	}

	if isRTL {
		for s, end := iEnd, iStart.prev; s != nil && s != end; s = s.prev {
			if s.isBase() {
				clusterMin = currpos.X
				currpos = s.finalise(seg, font, currpos, &bbox, 0, &clusterMin, isRTL, isFinal, 0)
			}
		}
	} else {
		for s, end := iStart, iEnd.Next; s != nil && s != end; s = s.Next {
			if s.isBase() {
				clusterMin = currpos.X
				currpos = s.finalise(seg, font, currpos, &bbox, 0, &clusterMin, isRTL, isFinal, 0)
			}
		}
	}
	if reorder {
		seg.reverseSlots()
	}

	return currpos
}

func (seg *Segment) doMirror(aMirror byte) {
	for s := seg.First; s != nil; s = s.Next {
		g := GID(seg.face.getGlyphAttr(s.glyphID, uint16(aMirror)))
		if g != 0 && (seg.dir&4 == 0 || seg.face.getGlyphAttr(s.glyphID, uint16(aMirror)+1) == 0) {
			s.setGlyph(seg, g)
		}
	}
}

func (seg *Segment) getSlotBidiClass(s *Slot) int8 {
	if res := s.bidiCls; res != -1 {
		return res
	}
	res := int8(seg.face.getGlyphAttr(s.glyphID, uint16(seg.silf.attrDirectionality)))
	s.bidiCls = res
	return res
}

// check the bounds and return nil if needed
func (seg *Segment) getCharInfo(index int) *charInfo {
	if index < len(seg.charinfo) {
		return &seg.charinfo[index]
	}
	return nil
}

// check the bounds and return nil if needed
func (seg *Segment) getCollisionInfo(s *Slot) *slotCollision {
	if s.index < len(seg.collisions) {
		return &seg.collisions[s.index]
	}
	return nil
}

func (seg *Segment) getFeature(findex uint8) int32 {
	// findex reference the font feat table
	if int(findex) < len(seg.face.feat) {
		featRef := seg.face.feat[findex].id
		if featVal := seg.feats.FindFeature(featRef); featVal != nil {
			return int32(featVal.Value)
		}
	}
	return 0
}

func (seg *Segment) setFeature(findex uint8, val int16) {
	// findex reference the font feat table
	if int(findex) < len(seg.face.feat) {
		featRef := seg.face.feat[findex].id
		if featVal := seg.feats.FindFeature(featRef); featVal != nil {
			featVal.Value = val
		}
	}
}

func (seg *Segment) getGlyphMetric(iSlot *Slot, metric, attrLevel uint8, rtl bool) int32 {
	if attrLevel > 0 {
		is := iSlot.findRoot()
		return is.clusterMetric(seg, metric, attrLevel, rtl)
	}
	return seg.face.getGlyphMetric(iSlot.glyphID, metric)
}

func (seg *Segment) finalise(font *FontOptions, reverse bool) {
	if seg.First == nil || seg.last == nil {
		return
	}
	seg.Advance = seg.positionSlots(font, seg.First, seg.last, seg.silf.isRTL, true)
	// associateChars(0, seg.numCharinfo);
	if reverse && seg.currdir() != (seg.dir&1 != 0) {
		seg.reverseSlots()
	}
	seg.linkClusters(seg.First, seg.last)
}

func (seg *Segment) linkClusters(s, end *Slot) {
	end = end.Next

	for ; s != end && !s.isBase(); s = s.Next {
	}
	ls := s

	if seg.dir&1 != 0 {
		for ; s != end; s = s.Next {
			if !s.isBase() {
				continue
			}

			s.sibling = ls
			ls = s
		}
	} else {
		for ; s != end; s = s.Next {
			if !s.isBase() {
				continue
			}

			ls.sibling = s
			ls = s
		}
	}
}

func (seg *Segment) associateChars(offset, numChars int) {
	subSlice := seg.charinfo[offset : offset+numChars]
	for i := range subSlice {
		subSlice[i].before = -1
		subSlice[i].after = -1
	}
	for s, i := seg.First, 0; s != nil; s, i = s.Next, i+1 {
		j := s.Before
		if j < 0 {
			continue
		}

		for after := s.After; j <= after; j++ {
			c := seg.getCharInfo(j)
			if c.before == -1 || i < c.before {
				c.before = i
			}
			if c.after < i {
				c.after = i
			}
		}
		s.index = i
	}
	for s := seg.First; s != nil; s = s.Next {
		var a int
		for a = s.After + 1; a < offset+numChars && seg.getCharInfo(a).after < 0; a++ {
			seg.getCharInfo(a).after = s.index
		}
		a--
		s.After = a

		for a = s.Before - 1; a >= offset && seg.getCharInfo(a).before < 0; a-- {
			seg.getCharInfo(a).before = s.index
		}
		a++
		s.Before = a
	}
}

func (seg *Segment) initCollisions() bool {
	seg.collisions = seg.collisions[:0]
	seg.collisions = append(seg.collisions, make([]slotCollision, seg.NumGlyphs)...)

	for p := seg.First; p != nil; p = p.Next {
		if p.index < seg.NumGlyphs {
			seg.collisions[p.index].init(seg, p)
		} else {
			return false
		}
	}
	return true
}
