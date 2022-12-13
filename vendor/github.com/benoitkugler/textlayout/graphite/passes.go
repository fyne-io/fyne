package graphite

import (
	"errors"
	"fmt"
	"sort"
)

type passtype uint8

const (
	ptUNKNOWN passtype = iota
	ptLINEBREAK
	ptSUBSTITUTE
	ptPOSITIONING
	ptJUSTIFICATION
)

// compute the columns from the ranges
func (pass *silfPass) computeColumns() ([]uint16, error) {
	if len(pass.ranges) == 0 {
		return nil, nil
	}
	numGlyphs := pass.ranges[len(pass.ranges)-1].LastId + 1
	cols := make([]uint16, numGlyphs)
	for i := range cols {
		cols[i] = 0xFFFF
	}
	for _, range_ := range pass.ranges {
		ci := range_.FirstId
		ciEnd := range_.LastId + 1
		col := range_.ColId

		if ci >= ciEnd || ciEnd > numGlyphs || col >= pass.NumColumns {
			return nil, fmt.Errorf("invalid pass range: %v", range_)
		}

		// A glyph must only belong to one column at a time
		for ci != ciEnd && cols[ci] == 0xffff {
			cols[ci] = col
			ci++
		}

		if ci != ciEnd {
			// we exit early, meaning a column was already attributed to a glyph
			return nil, errors.New("invalid pass range")
		}
	}
	return cols, nil
}

// load the code for the rules
func (pass *silfPass) computeRuleTable(context codeContext) ([]rule, error) {
	var err error
	out := make([]rule, pass.NumRules)
	for i := range pass.ruleSortKeys {
		r := rule{
			sortKey:    pass.ruleSortKeys[i],
			preContext: pass.rulePreContext[i],
		}

		if r.preContext > pass.maxRulePreContext || r.preContext < pass.minRulePreContext {
			return nil, fmt.Errorf("invalid rule preContext %d for [%d ... %d]", r.preContext, pass.minRulePreContext, pass.maxRulePreContext)
		}

		r.action, err = newCode(false, pass.actions[i], r.preContext, r.sortKey, context, false)
		if err != nil {
			return nil, fmt.Errorf("invalid rule action code: %s", err)
		}
		r.constraint, err = newCode(true, pass.ruleConstraints[i], r.preContext, r.sortKey, context, false)
		if err != nil {
			return nil, fmt.Errorf("invalid rule constraint code: %s", err)
		}

		out[i] = r
	}
	return out, nil
}

// performs the equivalent of --a in C
func decrease(a *uint8) uint8 {
	*a -= 1
	return *a
}

// encode the actions to apply to the input string
// it is directly obtained from the font file
type pass struct {
	// assign column to a subset of the glyph indices (GID . column; column < NumColumns)
	constraint *code // optional
	columns    []uint16
	// all the possible rules of the pass
	// their are activated conditionnaly on the input
	ruleTable []rule

	successStates [][]uint16 // (state index - numSuccess) . rule numbers (index into `rules`)
	startStates   []uint16
	transitions   [][]uint16 // each sub array has length NumColums

	collisionThreshold float32
	isReverseDirection bool
	collisionLoops     uint8
	kerningColls       uint8

	numStates                    uint16
	maxPreContext, minPreContext uint16
	maxRuleLoop                  uint8
}

// sanitizes and interprets one pass subtable
func newPass(tablePass *silfPass, context codeContext) (out pass, err error) {
	out.isReverseDirection = (tablePass.Flags>>5)&0x1 != 0
	out.collisionLoops = tablePass.Flags & 0x7
	out.kerningColls = (tablePass.Flags >> 3) & 0x3
	out.collisionThreshold = float32(tablePass.collisionThreshold)
	if out.collisionThreshold == 0 {
		out.collisionThreshold = 10 // default value
	}

	out.maxPreContext, out.minPreContext = uint16(tablePass.maxRulePreContext), uint16(tablePass.minRulePreContext)
	out.startStates = tablePass.startStates
	out.numStates = tablePass.NumRows
	out.transitions = tablePass.stateTransitions
	out.maxRuleLoop = tablePass.MaxRuleLoop
	out.successStates = tablePass.ruleMap

	if err = tablePass.sanitize(); err != nil {
		return out, fmt.Errorf("invalid silf pass subtable: %s", err)
	}

	out.columns, err = tablePass.computeColumns()
	if err != nil {
		return out, fmt.Errorf("invalid silf pass columns: %s", err)
	}

	out.ruleTable, err = tablePass.computeRuleTable(context)
	if err != nil {
		return out, fmt.Errorf("invalid silf pass rules: %s", err)
	}

	// sort the rules entries
	for _, l := range out.successStates {
		sort.Slice(l, func(i, j int) bool { return compareRuleIndex(out.ruleTable, l[i], l[j]) })
	}

	if len(tablePass.passConstraint) != 0 {
		context.Pt = ptUNKNOWN

		// if numRules == 0, which happens for instance in the Awami font
		// the "natural" value for tablePass.rulePreContext[0], tablePass.ruleSortKeys[0]
		// if the next field in the font file, that is tablePass.collisionThreshold
		preContext, ruleLength := tablePass.collisionThreshold, uint16(tablePass.collisionThreshold)
		if tablePass.NumRules != 0 {
			preContext, ruleLength = tablePass.rulePreContext[0], tablePass.ruleSortKeys[0]
		}

		constraint, err := newCode(true, tablePass.passConstraint, preContext, ruleLength, context, false)
		if err != nil {
			return out, fmt.Errorf("invalid silf pass constraint: %s", err)
		}
		out.constraint = &constraint
	}

	return out, nil
}

func (pass *pass) testPassConstraint(m *machine) (bool, error) {
	if pass.constraint == nil {
		return true, nil
	}

	m.map_.reset(m.map_.segment.First, 0)
	m.map_.pushSlot(m.map_.segment.First)
	ret, _, err := m.run(pass.constraint, 1)

	if debugMode >= 2 {
		tr.setCurrentPassConstraint(ret != 0 && err == nil)
	}

	return ret != 0 && err == nil, err
}

func (pa *pass) findAndDoRule(slot *Slot, m *machine, fsm *finiteStateMachine) (*Slot, error) {
	if rules := pa.runFSM(fsm, slot); len(rules) != 0 {
		// Search for the first rule which passes the constraint
		var (
			i int
			r uint16
		)
		for ; i < len(rules); i++ {
			r = rules[i]

			ok, err := pa.testConstraint(&fsm.ruleTable[r], m)
			if err != nil {
				return slot, fmt.Errorf("finding rule: %s", err)
			}
			if ok {
				break
			}
		}

		if debugMode >= 2 {
			tr.startDumpRule(fsm, i)
		}

		if i < len(rules) {
			r := rules[i]
			rule := &fsm.ruleTable[r]
			var (
				adv int32
				err error
			)
			adv, slot, err = pa.doAction(&rule.action, m)

			if debugMode >= 2 {
				tr.dumpRuleOutput(fsm, r, slot)
			}

			if err != nil {
				return slot, fmt.Errorf("applying rule: %s", err)
			}

			if rule.action.delete {
				slot = fsm.slots.collectGarbage(slot)
			}
			slot = pa.adjustSlot(adv, slot, &fsm.slots)

			if debugMode >= 2 {
				tr.dumpRuleCursor(slot)
			}

			return slot, nil
		}

		if debugMode >= 2 {
			tr.dumpRuleCursor(slot.Next)
		}
	}

	slot = slot.Next
	return slot, nil
}

// select the rules IDs to apply (may be empty)
func (pass *pass) runFSM(fsm *finiteStateMachine, slot *Slot) []uint16 {
	slot = fsm.reset(slot, pass.maxPreContext, pass.ruleTable)
	if fsm.slots.preContext < uint16(pass.minPreContext) {
		return nil
	}
	state := pass.startStates[pass.maxPreContext-fsm.slots.preContext]
	var freeSlots uint8 = maxSlots
	successStart := pass.numStates - uint16(len(pass.successStates)) // order checked in silfPassHeader.sanitize
	for do := true; do; do = state != 0 && slot != nil {
		fsm.slots.pushSlot(slot)
		if int(slot.glyphID) >= len(pass.columns) || pass.columns[slot.glyphID] == 0xffff ||
			decrease(&freeSlots) == 0 || int(state) >= len(pass.transitions) {
			if freeSlots == 0 {
				return nil
			}
			return fsm.rules
		}
		transitions := pass.transitions[state]
		state = transitions[pass.columns[slot.glyphID]]
		if state >= successStart {
			fsm.accumulateRules(pass.successStates[state-successStart])
		}

		slot = slot.Next
	}

	fsm.slots.pushSlot(slot)
	return fsm.rules
}

func (pass *pass) testConstraint(r *rule, m *machine) (bool, error) {
	currContext := m.map_.preContext
	rulePreContext := uint16(r.preContext)
	if currContext < rulePreContext || int(r.sortKey+currContext-rulePreContext) > m.map_.size {
		return false, nil
	}

	map_ := int(1 + currContext - rulePreContext)
	if m.map_.slots[map_+int(r.sortKey)-1] == nil {
		return false, nil
	}

	if len(r.constraint.instrs) == 0 {
		return true, nil
	}

	for n := r.sortKey; n != 0 && map_ != 0; n, map_ = n-1, map_+1 {
		if m.map_.slots[map_] == nil {
			continue
		}
		var (
			ret int32
			err error
		)
		ret, map_, err = m.run(&r.constraint, map_)
		if err != nil {
			return false, err
		}
		if ret == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (pass *pass) doAction(code *code, m *machine) (int32, *Slot, error) {
	if len(code.instrs) == 0 {
		return 0, nil, nil
	}
	smap := m.map_
	smap.highpassed = false

	ret, map_, err := m.run(code, int(smap.preContext)+1)
	if err != nil {
		smap.highwater = nil
		return 0, nil, err
	}

	return ret, m.map_.slots[map_], nil
}

func (pass *pass) adjustSlot(delta int32, slot *Slot, smap *slotMap) *Slot {
	if slot == nil {
		if smap.highpassed || slot == smap.highwater {
			slot = smap.segment.last
			delta++
			if smap.highwater == nil || smap.highwater == slot {
				smap.highpassed = false
			}
		} else {
			slot = smap.segment.First
			delta--
		}
	}
	if delta < 0 {
		for delta += 1; delta <= 0 && slot != nil; delta++ {
			slot = slot.prev
			if smap.highpassed && smap.highwater == slot {
				smap.highpassed = false
			}
		}
	} else if delta > 0 {
		for delta--; delta >= 0 && slot != nil; delta-- {
			if slot == smap.highwater && slot != nil {
				smap.highpassed = true
			}
			slot = slot.Next
		}
	}

	return slot
}

// Can slot s be kerned, or is it attached to something that can be kerned?
func inKernCluster(seg *Segment, s *Slot) bool {
	c := seg.getCollisionInfo(s)
	if c.flags&collKERN != 0 /** && c.flags & collFIX **/ {
		return true
	}
	for s.parent != nil {
		s = s.parent
		c = seg.getCollisionInfo(s)
		if c.flags&collKERN != 0 /** && c.flags & collFIX **/ {
			return true
		}
	}
	return false
}

// Fix collisions for the given slot.
// Return true if everything was fixed, false if there are still collisions remaining.
// isRev means be we are processing backwards.
func (pass *pass) resolveCollisions(seg *Segment, slotFix, start *Slot,
	coll *shiftCollider, isRev, isRTL bool, moved, hasCol *bool,
) (fixed bool) {
	var nbor *Slot // neighboring slot
	cFix := seg.getCollisionInfo(slotFix)
	if !coll.initSlot(seg, slotFix, cFix.limit, float32(cFix.margin), float32(cFix.marginWt),
		cFix.shift, cFix.offset, isRTL) {
		return false
	}
	collides := false
	// When we're processing forward, ignore kernable glyphs that preceed the target glyph.
	// When processing backward, don't ignore these until we pass slotFix.
	ignoreForKern := !isRev
	base := slotFix.findRoot()

	// Look for collisions with the neighboring glyphs.
	for nbor = start; nbor != nil; {
		cNbor := seg.getCollisionInfo(nbor)
		sameCluster := nbor.isChildOf(base)
		if nbor != slotFix && // don't process if this is the slot of interest
			!(cNbor.ignore()) && // don't process if ignoring
			(nbor == base || sameCluster || // process if in the same cluster as slotFix
				!inKernCluster(seg, nbor)) && // or this cluster is not to be kerned || (isRTL ^ ignoreForKern))       // or it comes before(ltr) or after(isRTL)
			(!isRev || // if processing forwards then good to merge otherwise only:
				!(cNbor.flags&collFIX != 0) || // merge in immovable stuff
				((cNbor.flags&collKERN != 0) && !sameCluster) || // ignore other kernable clusters
				(cNbor.flags&collISCOL != 0)) && // test against other collided glyphs
			!coll.mergeSlot(seg, nbor, cNbor, cNbor.shift, !ignoreForKern, sameCluster, false, &collides) {
			return false
		} else if nbor == slotFix {
			// Switching sides of this glyph - if we were ignoring kernable stuff before, don't anymore.
			ignoreForKern = !ignoreForKern
		}

		collConst := collEND
		if isRev {
			collConst = collSTART
		}
		if nbor != start && (cNbor.flags&collConst != 0) {
			break
		}

		if isRev {
			nbor = nbor.prev
		} else {
			nbor = nbor.Next
		}
	}
	isCol := false
	if collides || cFix.shift.X != 0. || cFix.shift.Y != 0. {
		var shift Position
		shift, isCol = coll.resolve(seg)
		// isCol has been set to true if a collision remains.
		if abs(shift.X) < 1e38 && abs(shift.Y) < 1e38 {
			if sqr(shift.X-cFix.shift.X)+sqr(shift.Y-cFix.shift.Y) >= sqr(pass.collisionThreshold) {
				*moved = true
			}
			cFix.shift = shift
			if slotFix.child != nil {
				var bbox rect
				here := slotFix.Position.add(shift)
				clusterMin := here.X
				slotFix.child.finalise(seg, nil, here, &bbox, 0, &clusterMin, isRTL, false, 0)
			}
		}
	} // else, This glyph is not colliding with anything.

	// Set the is-collision flag bit.
	if isCol {
		cFix.flags = cFix.flags | collISCOL | collKNOWN
	} else {
		cFix.flags = (cFix.flags & ^collISCOL) | collKNOWN
	}
	*hasCol = *hasCol || isCol
	return true
}

func (pass *pass) collisionShift(seg *Segment, isRTL bool) bool {
	var shiftcoll shiftCollider
	// bool isfirst = true;
	hasCollisions := false
	start := seg.First // turn on collision fixing for the first slot
	var end *Slot
	moved := false

	if debugMode >= 2 {
		tr.startDumpCollisions(pass.collisionLoops)
	}

	for start != nil {

		if debugMode >= 2 {
			tr.startDumpCollisionPhase("1", -1)
		}

		hasCollisions = false
		end = nil
		// phase 1 : position shiftable glyphs, ignoring kernable glyphs
		for s := start; s != nil; s = s.Next {
			c := seg.getCollisionInfo(s)
			if start != nil && (c.flags&(collFIX|collKERN)) == collFIX && !pass.resolveCollisions(seg, s, start, &shiftcoll, false, isRTL, &moved, &hasCollisions) {
				return false
			}
			if s != start && (c.flags&collEND) != 0 {
				end = s.Next
				break
			}
		}

		// #if !defined GRAPHITE2_NTRACING
		//         if (dbgout)
		//             *dbgout << json::close << json::close; // phase-1
		// #endif

		// phase 2 : loop until happy.
		for i := 0; i < int(pass.collisionLoops)-1; i++ {
			if hasCollisions || moved {

				if debugMode >= 2 {
					tr.startDumpCollisionPhase("2a", i)
				}

				// phase 2a : if any shiftable glyphs are in collision, iterate backwards,
				// fixing them and ignoring other non-collided glyphs. Note that this handles ONLY
				// glyphs that are actually in collision from phases 1 or 2b, and working backwards
				// has the intended effect of breaking logjams.
				if hasCollisions {
					hasCollisions = false
					// #if 0
					// moved = true;
					// for (Slot *s = start; s != end; s = s.Next)
					// {
					//     SlotCollision * c = seg.collisionInfo(s);
					//     c.setShift(Position(0, 0));
					// }
					// #endif
					lend := seg.last
					if end != nil {
						lend = end.prev
					}
					lstart := start.prev
					for s := lend; s != lstart; s = s.prev {
						c := seg.getCollisionInfo(s)
						if start != nil && (c.flags&(collFIX|collKERN|collISCOL)) == (collFIX|collISCOL) { // ONLY if this glyph is still colliding
							if !pass.resolveCollisions(seg, s, lend, &shiftcoll, true, isRTL, &moved, &hasCollisions) {
								return false
							}
							c.flags = c.flags | collTEMPLOCK
						}
					}
				}

				if debugMode >= 2 {
					tr.startDumpCollisionPhase("2b", i)
				}

				// phase 2b : redo basic diacritic positioning pass for ALL glyphs. Each successive loop adjusts
				// glyphs from their current adjusted position, which has the effect of gradually minimizing the
				// resulting adjustment; ie, the final result will be gradually closer to the original location.
				// Also it allows more flexibility in the final adjustment, since it is moving along the
				// possible 8 vectors from successively different starting locations.
				if moved {
					moved = false
					for s := start; s != end; s = s.Next {
						c := seg.getCollisionInfo(s)
						if start != nil && (c.flags&(collFIX|collTEMPLOCK|collKERN)) == collFIX &&
							!pass.resolveCollisions(seg, s, start, &shiftcoll, false, isRTL, &moved, &hasCollisions) {
							return false
						} else if c.flags&collTEMPLOCK != 0 {
							c.flags = c.flags & ^collTEMPLOCK
						}
					}
				}
				//      if (!hasCollisions) // no, don't leave yet because phase 2b will continue to improve things
				//          break;
				// #if !defined GRAPHITE2_NTRACING
				//                 if (dbgout)
				//                     *dbgout << json::close << json::close; // phase 2
				// #endif
			}
		}
		if end == nil {
			break
		}
		start = nil
		for s := end.prev; s != nil; s = s.Next {
			if seg.getCollisionInfo(s).flags&collSTART != 0 {
				start = s
				break
			}
		}
	}
	return true
}

func (pass *pass) collisionKern(seg *Segment, isRTL bool) bool {
	start := seg.First
	var (
		ymin float32 = 1e38
		ymax float32 = -1e38
	)

	// phase 3 : handle kerning of clusters
	if debugMode >= 2 {
		tr.startDumpCollisionPhase("3", -1)
	}

	for s := seg.First; s != nil; s = s.Next {
		if int(s.glyphID) >= len(seg.face.glyphs) {
			return false
		}
		c := seg.getCollisionInfo(s)
		bbox := seg.face.getGlyph(s.glyphID).bbox
		y := s.Position.Y + c.shift.Y
		if c.flags&collISSPACE == 0 {
			ymax = max(y+bbox.tr.Y, ymax)
			ymin = min(y+bbox.bl.Y, ymin)
		}
		if start != nil && (c.flags&(collKERN|collFIX)) == (collKERN|collFIX) {
			pass.resolveKern(seg, s, start, isRTL, &ymin, &ymax)
		}
		if c.flags&collEND != 0 {
			start = nil
		}
		if c.flags&collSTART != 0 {
			start = s
		}
	}

	return true
}

const (
	kernNone = iota
	kernCrossSpace
	kernInWord
	// Kernreserved
)

func (pass *pass) resolveKern(seg *Segment, slotFix, start *Slot, isRTL bool, ymin, ymax *float32) float32 {
	var currSpace float32
	collides := false
	spaceCount := 0
	base := slotFix.findRoot()
	cFix := seg.getCollisionInfo(base)
	// const GlyphCache &gc = seg.getFace().glyphs();
	bbb := seg.face.getGlyph(slotFix.glyphID).bbox
	by := slotFix.Position.Y + cFix.shift.Y

	if base != slotFix {
		cFix.flags = cFix.flags | collKERN | collFIX
		return 0
	}
	seenEnd := (cFix.flags & collEND) != 0
	isInit := false
	coll := newKernCollider()

	*ymax = max(by+bbb.tr.Y, *ymax)
	*ymin = min(by+bbb.bl.Y, *ymin)
	for nbor := slotFix.Next; nbor != nil; nbor = nbor.Next {
		if int(nbor.glyphID) >= len(seg.face.glyphs) {
			return 0.
		}
		bb := seg.face.getGlyph(nbor.glyphID).bbox
		cNbor := seg.getCollisionInfo(nbor)
		nby := nbor.Position.Y + cNbor.shift.Y
		if nbor.isChildOf(base) {
			*ymax = max(nby+bb.tr.Y, *ymax)
			*ymin = min(nby+bb.bl.Y, *ymin)
			continue
		}
		if (bb.bl.Y == 0. && bb.tr.Y == 0.) || (cNbor.flags&collISSPACE) != 0 {
			if pass.kerningColls == kernInWord {
				break
			}
			// Add space for a space glyph.
			currSpace += nbor.Advance.X
			spaceCount++
		} else {
			spaceCount = 0
			if nbor != slotFix && !cNbor.ignore() {
				seenEnd = true
				if !isInit {
					if !coll.initSlot(seg, slotFix, cFix.limit, float32(cFix.margin),
						cFix.shift, cFix.offset, isRTL, *ymin, *ymax) {
						return 0.
					}
					isInit = true
				}
				maybeCollide := coll.mergeSlot(seg, nbor, cNbor.shift, currSpace, isRTL)
				collides = collides || maybeCollide
			}
		}
		if cNbor.flags&collEND != 0 {
			if seenEnd && spaceCount < 2 {
				break
			} else {
				seenEnd = true
			}
		}
	}
	if collides {
		mv := coll.resolve(isRTL)
		coll.shift(mv, isRTL)
		delta := slotFix.Advance.add(mv).sub(cFix.shift)
		slotFix.Advance = delta
		cFix.shift = mv
		return mv.X
	}
	return 0.
}

func (pass *pass) collisionFinish(seg *Segment) {
	for s := seg.First; s != nil; s = s.Next {
		c := seg.getCollisionInfo(s)
		if c.shift.X != 0 || c.shift.Y != 0 {
			newOffset := c.shift
			var nullPosition Position
			c.offset = newOffset.add(c.offset)
			c.shift = nullPosition
		}
	}
	//    seg.positionSlots();

	// #if !defined GRAPHITE2_NTRACING
	//         if (dbgout)
	//             *dbgout << json::close;
	// #endif
}

func (pass *pass) runGraphite(m *machine, fsm *finiteStateMachine, reverse bool) (bool, error) {
	s := m.map_.segment.First
	if s == nil {
		return true, nil
	}

	if ok, err := pass.testPassConstraint(m); !ok {
		return true, err
	}
	if reverse {
		m.map_.segment.reverseSlots()
		s = m.map_.segment.First
	}
	if len(pass.ruleTable) != 0 {
		currHigh := s.Next

		m.map_.highwater = currHigh
		lc := pass.maxRuleLoop

		var err error
		for do := true; do; do = s != nil {
			s, err = pass.findAndDoRule(s, m, fsm)
			if err != nil {
				return false, err
			}
			if s != nil && (s == m.map_.highwater || m.map_.highpassed || decrease(&lc) == 0) {
				if lc == 0 {
					s = m.map_.highwater
				}
				lc = pass.maxRuleLoop
				if s != nil {
					m.map_.highwater = s.Next
				}
			}
		}
	}

	collisions := pass.collisionLoops != 0 || pass.kerningColls != 0

	if !collisions || !m.map_.segment.hasCollisionInfo() {
		return true, nil
	}

	if pass.collisionLoops != 0 {
		if (m.map_.segment.flags & initCollisions) == 0 {
			m.map_.segment.positionSlots(nil, nil, nil, m.map_.isRTL, true)
		}
		if !pass.collisionShift(m.map_.segment, m.map_.isRTL) {
			return false, nil
		}
	}
	if (pass.kerningColls != 0) && !pass.collisionKern(m.map_.segment, m.map_.isRTL) {
		return false, nil
	}

	if collisions {
		pass.collisionFinish(m.map_.segment)
	}

	return true, nil
}

// higher level version of a silf subtable
type passes struct {
	passes              []pass
	pseudoMaps          []pseudoMap
	justificationLevels []justificationLevel
	classMap            classMap

	userAttibutes      uint8 // Number of user-defined slot attributes
	attrPseudo         byte  // Glyph attribute number that is used for actual glyph ID for a pseudo glyph
	attrBreakWeight    byte  // Glyph attribute number of breakweight attribute
	attrDirectionality byte  // Glyph attribute number for directionality attribute
	attrMirroring      byte  // Glyph attribute number for mirror.glyph (mirror.isEncoded comes directly after)
	attrSkipPasses     byte  // Glyph attribute of bitmap indicating key glyphs for pass optimization
	attrCollision      byte  // Glyph attribute number for collision.flags attribute (several more collision attrs come after it...)

	indexBidiPass byte // (0xFF) means no bidi pass
	indexPosPass  byte // index of the first positionning pass
	hasCollision  bool
	isRTL         bool
}

// interprets and sanitizes the subtable
func newPasses(silf *silfSubtable, numAttributes, numFeatures uint16) (out passes, err error) {
	out.passes = make([]pass, len(silf.passes))

	context := codeContext{
		NumAttributes:     numAttributes,
		NumFeatures:       numFeatures,
		NumClasses:        silf.classMap.numClasses(),
		NumUserAttributes: silf.NumUserDefn,
	}
	for i := range silf.passes {
		pass := &silf.passes[i]

		// resolve the pass type
		context.Pt = ptUNKNOWN
		switch {
		case i >= int(silf.IJust):
			context.Pt = ptJUSTIFICATION
		case i >= int(silf.IPos):
			context.Pt = ptPOSITIONING
		case i >= int(silf.ISubst):
			context.Pt = ptSUBSTITUTE
		default:
			context.Pt = ptLINEBREAK
		}

		out.passes[i], err = newPass(pass, context)
		if err != nil {
			return out, fmt.Errorf("invalid silf pass %d: %s", i, err)
		}
	}

	out.pseudoMaps = silf.pseudoMap
	out.justificationLevels = silf.justificationLevels
	out.classMap = silf.classMap

	out.userAttibutes = silf.NumUserDefn
	out.attrPseudo = silf.AttrPseudo
	out.attrBreakWeight = silf.AttrBreakWeight
	out.attrDirectionality = silf.AttrDirectionality
	out.attrMirroring = silf.AttrMirroring
	out.attrSkipPasses = silf.AttrSkipPasses
	out.attrCollision = silf.AttrCollisions

	out.indexBidiPass = silf.IBidi
	out.indexPosPass = silf.IPos
	out.hasCollision = silf.Flags&0x20 != 0
	// see the reference implementation for this switch
	out.isRTL = (silf.Direction-1)&1 != 0
	return out, nil
}

func (s *passes) findPdseudoGlyph(r rune) GID {
	if s == nil {
		return 0
	}
	for _, rec := range s.pseudoMaps {
		if rec.Unicode == r {
			return GID(rec.NPseudo)
		}
	}
	return 0
}

func (s *passes) runGraphite(seg *Segment, firstPass, lastPass uint8, doBidi bool) bool {
	maxSize := len(seg.charinfo) * maxSegGrowthFactor

	fsm := &finiteStateMachine{slots: newSlotMap(seg, s.isRTL, maxSize)}
	m := newMachine(&fsm.slots) // sharing slots

	lbidi := s.indexBidiPass

	if lastPass == 0 {
		if firstPass == lastPass && lbidi == 0xFF {
			return true
		}
		lastPass = uint8(len(s.passes))
	}
	if (firstPass < lbidi || (doBidi && firstPass == lbidi)) && (lastPass >= lbidi || (doBidi && lastPass+1 == lbidi)) {
		lastPass++
	} else {
		lbidi = 0xFF
	}

	for i := firstPass; i < lastPass; i++ {
		if debugMode >= 1 {
			fmt.Printf("Pass %d, segment direction %v", i, seg.currdir())
		}

		// bidi and mirroring
		if i == lbidi {

			if seg.currdir() != s.isRTL {
				seg.reverseSlots()
			}
			if mirror := s.attrMirroring; mirror != 0 && (seg.dir&3) == 3 {
				seg.doMirror(mirror)
			}
			i--
			lbidi = lastPass
			lastPass--
			continue
		}

		if debugMode >= 2 {
			seg.positionSlots(nil, nil, nil, seg.currdir(), true)
			tr.appendPass(s, seg, i)
		}

		// test whether to reorder, prepare for positioning
		reverse := (lbidi == 0xFF) && (seg.currdir() != (s.isRTL != s.passes[i].isReverseDirection))
		var err error
		if i >= 32 || (seg.passBits&(1<<i)) == 0 || s.passes[i].collisionLoops != 0 {
			var ok bool
			ok, err = s.passes[i].runGraphite(m, fsm, reverse)
			if !ok {
				return false
			}
		}
		// only subsitution passes can change segment length, cached subsegments are short for their text
		if err != nil || (len(seg.charinfo) != 0 && len(seg.charinfo) > maxSize) {
			return false
		}
	}
	return true
}
