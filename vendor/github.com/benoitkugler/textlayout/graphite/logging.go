package graphite

import (
	"encoding/json"
	"fmt"
)

// this file implements tracing helpers, which are only used
// in debug mode

var tr = &traceOutput{}

type traceOutput struct {
	colliderEnv colliderEnv
	Outputdir   string     `json:"outputdir"`
	Id          string     `json:"id"`
	Output      []slotJSON `json:"output"`
	Chars       []charInfo `json:"chars"`
	Passes      []passJSON `json:"passes"`
	Advance     Position   `json:"advance"`
}

type colliderEnv struct {
	sl  *Slot
	val int
}

func (tr *traceOutput) appendPass(s *passes, seg *Segment, i uint8) {
	sd, pd := "ltr", "ltr"
	if seg.currdir() {
		sd = "rtl"
	}
	if s.isRTL != s.passes[i].isReverseDirection {
		pd = "rtl"
	}
	debug := passJSON{
		ID:       i + 1,
		Slotsdir: sd,
		Passdir:  pd,
		Slots:    seg.slotsJSON(),
		Rules:    make([]ruleDump, 0),
	}
	tr.Passes = append(tr.Passes, debug)
}

func (tr *traceOutput) finaliseOutput(seg *Segment) {
	tr.Outputdir = "ltr"
	if seg.currdir() {
		tr.Outputdir = "rtl"
	}
	tr.Output = seg.slotsJSON()
	tr.Advance = seg.Advance
	tr.Chars = seg.charinfo
}

func (ci charInfo) MarshalJSON() ([]byte, error) {
	type charInfoSlotJSON struct {
		Before int `json:"before"`
		After  int `json:"after"`
	}
	type charInfoJSON struct {
		Offset  int              `json:"offset"`
		Unicode rune             `json:"unicode"`
		Break   int16            `json:"break"`
		Flags   uint8            `json:"flags"`
		Slot    charInfoSlotJSON `json:"slot"`
	}
	out := charInfoJSON{
		Offset:  ci.base,
		Unicode: ci.char,
		Break:   ci.breakWeight,
		Flags:   ci.flags,
		Slot: charInfoSlotJSON{
			Before: ci.before,
			After:  ci.after,
		},
	}
	return json.Marshal(out)
}

type slotCharInfoJSON struct {
	Original int `json:"original"`
	Before   int `json:"before"`
	After    int `json:"after"`
}

func (pos Position) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%.4f, %.4f]", pos.X, pos.Y)), nil
}

func (r rect) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%.4f, %.4f, %.4f, %.4f]", r.bl.X, r.bl.Y, r.tr.X, r.tr.Y)), nil
}

func (s *Slot) objectID() string {
	if s == nil {
		return "0000-00-0000"
	}
	return fmt.Sprintf("%p", s)
}

type slotParentJSON struct {
	Id     string   `json:"id"`
	Level  int32    `json:"level"`
	Offset Position `json:"offset"`
}

type collisionSeq struct {
	Seqclass  Position `json:"seqclass"`
	Seqorder  uint16   `json:"seqorder"`
	Seqabove  Position `json:"seqabove"`
	Seqbelow  Position `json:"seqbelow"`
	Seqvalign Position `json:"seqvalign"`
}

type collisionJSON struct {
	Limit         rect     `json:"limit"`
	Offset        Position `json:"offset"`
	Margin        Position `json:"margin"`
	Excludeoffset Position `json:"excludeoffset"`
	Flags         uint16   `json:"flags"`
	Exclude       GID      `json:"exclude"`
}

type collisionJSONWithSeq struct {
	collisionJSON
	collisionSeq
}

func (col collisionJSONWithSeq) MarshalJSON() ([]byte, error) {
	if col.Seqorder != 0 {
		type noMethod collisionJSONWithSeq
		return json.Marshal(noMethod(col))
	}
	return json.Marshal(col.collisionJSON)
}

type slotJSON struct {
	Collision     *collisionJSONWithSeq `json:"collision,omitempty"`
	Parent        *slotParentJSON       `json:"parent,omitempty"`
	Id            string                `json:"id"`
	User          []int16               `json:"user"`
	Children      []string              `json:"children,omitempty"`
	Charinfo      slotCharInfoJSON      `json:"charinfo"`
	Origin        Position              `json:"origin"`
	Shift         Position              `json:"shift"`
	Advance       Position              `json:"advance"`
	Justification float32               `json:"justification,omitempty"`
	Break         int32                 `json:"break"`
	Gid           GID                   `json:"gid"`
	Bidi          uint8                 `json:"bidi,omitempty"`
	Insert        bool                  `json:"insert"`
}

// returns a JSON compatible representation of the slot
func (s *Slot) json(seg *Segment) slotJSON {
	out := slotJSON{
		Id:  s.objectID(),
		Gid: s.glyphID,
		Charinfo: slotCharInfoJSON{
			Original: s.original,
			Before:   s.Before,
			After:    s.After,
		},
		Origin: s.Position,
		Shift: Position{
			X: float32(s.getAttr(nil, acShiftX, 0)),
			Y: float32(s.getAttr(nil, acShiftY, 0)),
		},
		Advance:       s.Advance,
		Insert:        s.CanInsertBefore(),
		Break:         s.getAttr(seg, acBreak, 0),
		Justification: s.just,
		Bidi:          s.bidiLevel,
		User:          append([]int16(nil), s.userAttrs...),
	}
	if !s.isBase() {
		out.Parent = &slotParentJSON{
			Id:     s.parent.objectID(),
			Level:  s.getAttr(nil, acAttLevel, 0),
			Offset: s.attach.sub(s.with),
		}
	}
	if s.child != nil {
		for c := s.child; c != nil; c = c.sibling {
			out.Children = append(out.Children, c.objectID())
		}
	}
	if cslot := seg.getCollisionInfo(s); cslot != nil {
		// Note: the reason for using Positions to lump together related attributes is to make the
		// JSON output slightly more compact.
		out.Collision = &collisionJSONWithSeq{
			collisionJSON: collisionJSON{
				Offset:        cslot.offset,
				Limit:         cslot.limit,
				Flags:         cslot.flags,
				Margin:        Position{float32(cslot.margin), float32(cslot.marginWt)},
				Exclude:       cslot.exclGlyph,
				Excludeoffset: cslot.exclOffset,
			},
			collisionSeq: collisionSeq{
				Seqclass:  Position{float32(cslot.seqClass), float32(cslot.seqProxClass)},
				Seqorder:  cslot.seqOrder,
				Seqabove:  Position{float32(cslot.seqAboveXoff), float32(cslot.seqAboveWt)},
				Seqbelow:  Position{float32(cslot.seqBelowXlim), float32(cslot.seqBelowWt)},
				Seqvalign: Position{float32(cslot.seqValignHt), float32(cslot.seqValignWt)},
			},
		}
	}
	return out
}

func (seg *Segment) slotsJSON() (out []slotJSON) {
	for s := seg.First; s != nil; s = s.Next {
		out = append(out, s.json(seg))
	}
	return out
}

type passJSON struct {
	Collisions *passCollisions `json:"collisions,omitempty"`
	Constraint *bool           `json:"constraint,omitempty"`
	Slotsdir   string          `json:"slotsdir"`
	Passdir    string          `json:"passdir"`
	Slots      []slotJSON      `json:"slots"`
	Rules      []ruleDump      `json:"rules"`
	ID         uint8           `json:"id"`
}

type ruleDump struct {
	Output     *ruleOutput `json:"output"`
	Cursor     string      `json:"cursor"`
	Considered []ruleJSON  `json:"considered"`
}

type slotRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ruleOutput struct {
	Range     slotRange  `json:"range"`
	Slots     []slotJSON `json:"slots"`
	Postshift Position   `json:"postshift"`
}

func inputSlot(slots *slotMap, n int) *Slot {
	s := slots.get(int(slots.preContext) + n)
	if !s.isCopied() {
		return s
	}

	if s.prev != nil {
		return s.prev.Next
	}
	if s.Next != nil {
		return s.Next.prev
	}
	return slots.segment.last
}

func outputSlot(slots *slotMap, n int) *Slot {
	s := slots.get(int(slots.preContext) + n - 1)
	if s != nil {
		return s.Next
	}
	return slots.segment.First
}

type ruleJSON struct {
	Input struct {
		Start  string `json:"start"`
		Length uint16 `json:"length"`
	} `json:"input,omitempty"`
	ID     uint16 `json:"id"`
	Failed bool   `json:"failed"`
}

func (tr *traceOutput) startDumpRule(fsm *finiteStateMachine, length int) {
	pass := &tr.Passes[len(tr.Passes)-1]

	var considered []ruleJSON
	for _, ruleIndex := range fsm.rules[:length] {
		r := fsm.ruleTable[ruleIndex]
		if uint16(r.preContext) > fsm.slots.preContext {
			continue
		}
		rj := ruleJSON{
			ID:     ruleIndex,
			Failed: true,
			Input: struct {
				Start  string `json:"start"`
				Length uint16 `json:"length"`
			}{
				Start:  inputSlot(&fsm.slots, -int(r.preContext)).objectID(),
				Length: r.sortKey,
			},
		}
		considered = append(considered, rj)
	}

	pass.Rules = append(pass.Rules, ruleDump{Considered: considered})
}

func (tr *traceOutput) dumpRuleOutput(fsm *finiteStateMachine, ruleIndex uint16, lastSlot *Slot) {
	r := fsm.ruleTable[ruleIndex]
	rj := ruleJSON{
		ID:     ruleIndex,
		Failed: false,
		Input: struct {
			Start  string `json:"start"`
			Length uint16 `json:"length"`
		}{
			Start:  inputSlot(&fsm.slots, 0).objectID(),
			Length: r.sortKey - uint16(r.preContext),
		},
	}

	pass := &tr.Passes[len(tr.Passes)-1]
	rule := &pass.Rules[len(pass.Rules)-1]

	rule.Considered = append(rule.Considered, rj)

	oj := ruleOutput{
		Range: slotRange{inputSlot(&fsm.slots, 0).objectID(), lastSlot.objectID()},
	}
	rsbPrepos := fsm.slots.segment.Advance
	if lastSlot != nil {
		rsbPrepos = lastSlot.Position
	}
	fsm.slots.segment.positionSlots(nil, nil, nil, fsm.slots.segment.currdir(), true)
	for slot := outputSlot(&fsm.slots, 0); slot != lastSlot; slot = slot.Next {
		oj.Slots = append(oj.Slots, slot.json(fsm.slots.segment))
	}

	if lastSlot != nil {
		oj.Postshift = lastSlot.Position
	} else {
		oj.Postshift = fsm.slots.segment.Advance
	}
	oj.Postshift = oj.Postshift.sub(rsbPrepos)

	rule.Output = &oj
}

func (tr *traceOutput) dumpRuleCursor(s *Slot) {
	pass := &tr.Passes[len(tr.Passes)-1]
	rule := &pass.Rules[len(pass.Rules)-1]
	rule.Cursor = s.objectID()
}

func (tr *traceOutput) setCurrentPassConstraint(b bool) {
	tr.Passes[len(tr.Passes)-1].Constraint = &b
}

type passCollisions struct {
	phases   []collisionPhase
	numLoops uint8
}

func (pc passCollisions) MarshalJSON() ([]byte, error) {
	tmp := []interface{}{
		map[string]uint8{"num-loops": pc.numLoops},
	}
	for _, v := range pc.phases {
		tmp = append(tmp, v)
	}
	return json.Marshal(tmp)
}

type phaseMove interface {
	isPhaseMove()
}

func (*collisionMove) isPhaseMove() {}
func (kernMove) isPhaseMove()       {}

type collisionPhase struct {
	Phase string      `json:"phase"`
	Moves []phaseMove `json:"moves"`
	Loop  int         `json:"loop"`
}

func (cl collisionPhase) MarshalJSON() ([]byte, error) {
	if cl.Loop == -1 { // ignore Loop here
		type noLoop struct {
			Phase string      `json:"phase"`
			Moves []phaseMove `json:"moves"`
		}
		tmp := noLoop{Phase: cl.Phase, Moves: cl.Moves}
		return json.Marshal(tmp)
	}
	type noMethod collisionPhase
	return json.Marshal(noMethod(cl))
}

type collisionMove struct {
	Slot     string              `json:"slot"`
	Vectors  []collisionVector   `json:"vectors"`
	Target   collisionMoveTarget `json:"target"`
	BestAxis int                 `json:"bestAxis"`
	Limit    rect                `json:"limit"`
	Result   Position            `json:"result"`
	Gid      GID                 `json:"gid"`
	StillBad bool                `json:"stillBad"`
}

type collisionMoveTarget struct {
	Fix        string   `json:"fix"`
	Bbox       rect     `json:"bbox"`
	SlantBox   rect     `json:"slantBox"`
	Origin     Position `json:"origin"`
	CurrShift  Position `json:"currShift"`
	CurrOffset Position `json:"currOffset"`
}

type collisionVector struct {
	Direction string          `json:"direction"`
	Removals  [][]interface{} `json:"removals"`
	Ranges    []interface{}   `json:"ranges"`
	TargetMin float32         `json:"targetMin"`
	BestCost  float32         `json:"bestCost"`
	BestVal   float32         `json:"bestVal"`
}

func (tr *traceOutput) startDumpCollisions(numLoops uint8) {
	tr.Passes[len(tr.Passes)-1].Collisions = &passCollisions{numLoops: numLoops}
}

func (tr *traceOutput) startDumpCollisionPhase(phase string, loop int) {
	cl := tr.Passes[len(tr.Passes)-1].Collisions
	cl.phases = append(cl.phases, collisionPhase{Phase: phase, Loop: loop})
}

func (tr *traceOutput) addCollisionMove(sc *shiftCollider, seg *Segment) {
	cl := tr.Passes[len(tr.Passes)-1].Collisions
	phase := &cl.phases[len(cl.phases)-1]
	phase.Moves = append(phase.Moves, &collisionMove{
		Slot:  sc.target.objectID(),
		Gid:   sc.target.glyphID,
		Limit: sc.limit,
		Target: collisionMoveTarget{
			Origin:     sc.origin,
			CurrShift:  sc.currShift,
			CurrOffset: seg.getCollisionInfo(sc.target).offset,
			Bbox:       seg.face.getGlyph(sc.target.glyphID).bbox,
			SlantBox:   seg.face.getGlyph(sc.target.glyphID).boxes.slant,
			Fix:        "shift",
		},
	})
}

func (tr *traceOutput) currentPhaseMove() phaseMove {
	cl := tr.Passes[len(tr.Passes)-1].Collisions
	phase := &cl.phases[len(cl.phases)-1]
	return phase.Moves[len(phase.Moves)-1]
}

func (tr *traceOutput) endCollisionMove(resultPos Position, bestAxis int, isCol bool) {
	move := tr.currentPhaseMove().(*collisionMove)
	move.Result = resultPos
	//<< "scraping" << _scraping[bestAxis]
	move.BestAxis = bestAxis
	move.StillBad = isCol
}

func (tr *traceOutput) addCollisionVector(sc *shiftCollider, seg *Segment, axis int,
	tleft, bestCost, bestVal float32,
) {
	var out collisionVector
	switch axis {
	case 0:
		out.Direction = "x"
	case 1:
		out.Direction = "y"
	case 2:
		out.Direction = "sum (NE-SW)"
	case 3:
		out.Direction = "diff (NW-SE)"
	default:
		out.Direction = "???"
	}

	out.TargetMin = tleft

	out.Removals = sc.ranges[axis].formatDebugs(seg)
	out.Ranges = sc.debugAxis(seg, axis)

	out.BestCost = bestCost
	out.BestVal = bestVal + tleft

	move := tr.currentPhaseMove().(*collisionMove)
	move.Vectors = append(move.Vectors, out)
}

type zoneDebug struct {
	env   colliderEnv
	excl  exclusion
	isDel bool
}

type fl5 [5]float32

func (fs fl5) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%.4f, %.4f, %.4f, %.4f, %.4f]",
		fs[0], fs[1], fs[2], fs[3], fs[4])), nil
}

func (zo zones) formatDebugs(seg *Segment) [][]interface{} {
	var out [][]interface{}
	for _, s := range zo.debugs {
		l := []interface{}{
			s.env.sl.objectID(),
			s.env.val,
		}
		if s.isDel {
			l = append(l, "remove", Position{s.excl.x, s.excl.xm})
		} else {
			l = append(l, "exclude", fl5{
				s.excl.x, s.excl.xm,
				s.excl.sm, s.excl.smx, s.excl.c,
			})
		}
		out = append(out, l)
	}
	return out
}

func (sc *shiftCollider) debugAxis(seg *Segment, axis int) []interface{} {
	var out []interface{}
	out = append(out, Position{sc.ranges[axis].pos, sc.ranges[axis].posm})
	// *dbgout << json::flat << json::array << _ranges[axis].position();
	for _, s := range sc.ranges[axis].exclusions {
		l := []interface{}{
			Position{s.x, s.xm},
			s.sm,
			s.smx,
			s.c,
		}
		out = append(out, l)
	}
	return out
}

type kernTarget struct {
	Fix        string   `json:"fix"`
	Bbox       rect     `json:"bbox"`
	SlantBox   rect     `json:"slantBox"`
	Origin     Position `json:"origin"`
	OffsetPrev Position `json:"offsetPrev"`
}

type kernSlice struct {
	Neighbor   string  `json:"neighbor"`
	I          int     `json:"i"`
	TargetEdge float32 `json:"targetEdge"`
	NearEdge   float32 `json:"nearEdge"`
}

type kernMove struct {
	Slot       string      `json:"slot"`
	Slices     []kernSlice `json:"slices"`
	Target     kernTarget  `json:"target"`
	Limit      rect        `json:"limit"`
	Miny       float32     `json:"miny"`
	Maxy       float32     `json:"maxy"`
	Slicewidth float32     `json:"slicewidth"`
	Xbound     float32     `json:"xbound"`
	MinGap     float32     `json:"minGap"`
	Needed     float32     `json:"needed"`
	Result     float32     `json:"result"`
	Gid        GID         `json:"gid"`
	StillBad   bool        `json:"stillBad"`
}

func (tr *traceOutput) addKern(kc *kernCollider, seg *Segment, result, resultNeeded float32) {
	var out kernMove
	out.Slot = kc.target.objectID()
	out.Gid = kc.target.glyphID
	out.Limit = kc.limit
	out.Miny = kc.miny
	out.Maxy = kc.maxy
	out.Slicewidth = kc.sliceWidth
	out.Target = kernTarget{
		Origin: kc.target.Position,
		// CurrShift: kc.currShift,
		OffsetPrev: kc.offsetPrev,
		Bbox:       seg.face.getGlyph(kc.target.glyphID).bbox,
		SlantBox:   seg.face.getGlyph(kc.target.glyphID).boxes.slant,
		Fix:        "kern",
	}

	for is := range kc.edges {
		out.Slices = append(out.Slices, kernSlice{
			I:          is,
			TargetEdge: kc.edges[is],
			Neighbor:   kc.slotNear[is].objectID(),
			NearEdge:   kc.nearEdges[is],
		})
	}

	out.Xbound = kc.xbound
	out.MinGap = kc.mingap
	out.Needed = resultNeeded
	out.Result = result
	out.StillBad = (result != resultNeeded)

	cl := tr.Passes[len(tr.Passes)-1].Collisions
	phase := &cl.phases[len(cl.phases)-1]
	phase.Moves = append(phase.Moves, out)
}
