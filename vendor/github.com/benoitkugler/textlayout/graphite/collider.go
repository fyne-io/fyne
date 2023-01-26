package graphite

import (
	"math"
)

const (
	collFIX      uint16 = 1 << iota // fix collisions involving this glyph
	collIGNORE                      // ignore this glyph altogether
	collSTART                       // start of range of possible collisions
	collEND                         // end of range of possible collisions
	collKERN                        // collisions with this glyph are fixed by adding kerning space after it
	collISCOL                       // this glyph has a collision
	collKNOWN                       // we've figured out what's happening with this glyph
	collISSPACE                     // treat this glyph as a space with regard to kerning
	collTEMPLOCK                    // Lock glyphs that have been given priority positioning
)

// Behavior for the collision.order attribute. To GDL this is an enum, to us it's a bitfield, with only 1 bit set
// Allows for easier inversion.
const (
	seqOrderLEFTDOWN = 1 << iota
	seqOrderRIGHTUP
	seqOrderNOABOVE
	seqOrderNOBELOW
	seqOrderNOLEFT
	seqOrderNORIGHT
)

// slot attributes related to collision-fixing
type slotCollision struct {
	limit        rect
	shift        Position // adjustment within the given pass
	offset       Position // total adjustment for collisions
	exclOffset   Position
	exclGlyph    GID
	margin       uint16
	marginWt     uint16
	flags        uint16
	seqClass     uint16
	seqProxClass uint16
	seqOrder     uint16
	seqAboveXoff int16
	seqAboveWt   uint16
	seqBelowXlim int16
	seqBelowWt   uint16
	seqValignHt  uint16
	seqValignWt  uint16
}

// Initialize the collision attributes for the given slot.
func (sc *slotCollision) init(seg *Segment, slot *Slot) {
	// Initialize slot attributes from glyph attributes.
	// The order here must match the order in the grcompiler code,
	// GrcSymbolTable::AssignInternalGlyphAttrIDs.
	gid := slot.glyphID
	aCol := uint16(seg.silf.attrCollision) // flags attr ID
	glyphFace := seg.face.getGlyph(gid)
	if glyphFace == nil {
		return
	}
	p := glyphFace.attrs
	sc.flags = uint16(p.get(aCol))
	sc.limit = rect{
		Position{float32(int16(p.get(aCol + 1))), float32(int16(p.get(aCol + 2)))},
		Position{float32(int16(p.get(aCol + 3))), float32(int16(p.get(aCol + 4)))},
	}
	sc.margin = uint16(p.get(aCol + 5))
	sc.marginWt = uint16(p.get(aCol + 6))

	sc.seqClass = uint16(p.get(aCol + 7))
	sc.seqProxClass = uint16(p.get(aCol + 8))
	sc.seqOrder = uint16(p.get(aCol + 9))
	sc.seqAboveXoff = p.get(aCol + 10)
	sc.seqAboveWt = uint16(p.get(aCol + 11))
	sc.seqBelowXlim = p.get(aCol + 12)
	sc.seqBelowWt = uint16(p.get(aCol + 13))
	sc.seqValignHt = uint16(p.get(aCol + 14))
	sc.seqValignWt = uint16(p.get(aCol + 15))

	// These attributes do not have corresponding glyph attribute:
	sc.exclGlyph = 0
	sc.exclOffset = Position{}
}

func (sc *slotCollision) ignore() bool {
	return (sc.flags&collIGNORE) != 0 || (sc.flags&collISSPACE) != 0
}

type exclusion struct {
	x    float32 // x position
	xm   float32 // xmax position
	c    float32 // constant + sum(MiXi^2)
	sm   float32 // sum(Mi)
	smx  float32 // sum(MiXi)
	open bool
}

func newExclusionWeightedXY(xmin, xmax, f, a0, m, xi, c float32) exclusion {
	return exclusion{
		x: xmin, xm: xmax,
		sm:  m + f,
		smx: m * xi,
		c:   m*xi*xi + f*a0*a0 + c,
	}
}

func newExclusionWeightedSD(xmin, xmax, f, a0,
	m, xi, ai, c float32, nega bool) exclusion {
	xia := xi + ai
	if nega {
		xia = xi - ai
	}
	return exclusion{
		x: xmin, xm: xmax,
		sm:  0.25 * (m + 2.*f),
		smx: 0.25 * m * xia,
		c:   0.25*(m*xia*xia+2.*f*a0*a0) + c,
	}
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (e exclusion) outcode(val float32) uint8 {
	// float d = std::numeric_limits<float>::epsilon();
	var zero float32
	return (boolToUint8((val-e.xm) >= zero) << 1) | boolToUint8(e.x-val > zero)
}

// add other to e
func (e *exclusion) add(other exclusion) {
	e.c += other.c
	e.sm += other.sm
	e.smx += other.smx
	e.open = false
}

func (e *exclusion) splitAt(p float32) exclusion {
	r := *e
	r.xm = p
	e.x = p
	return r
}

func (e *exclusion) leftTrim(p float32) { e.x = p }

// Cost and test position functions

func (e *exclusion) trackCost(bestCost, bestPos *float32, origin float32) bool {
	p := e.testPosition(origin)
	localc := e.cost(p - origin)
	if e.open && localc > *bestCost {
		return true
	}

	if localc < *bestCost {
		*bestCost = localc
		*bestPos = p
	}
	return false
}

func (e exclusion) cost(p float32) float32 {
	return (e.sm*p-2*e.smx)*p + e.c
}

func (e exclusion) testPosition(origin float32) float32 {
	if e.sm < 0 {
		// sigh, test both ends and perhaps the middle too!
		res := e.x
		cl := e.cost(e.x)
		if e.x < origin && e.xm > origin {
			co := e.cost(origin)
			if co < cl {
				cl = co
				res = origin
			}
		}
		cr := e.cost(e.xm)
		return pick(cl > cr, e.xm, res)
	} else {
		zerox := e.smx/e.sm + origin
		if zerox < e.x {
			return e.x
		} else if zerox > e.xm {
			return e.xm
		} else {
			return zerox
		}
	}
}

func separated(v1, v2 float32) bool { return v1 != v2 }
func sqr(v float32) float32         { return v * v }

// represents the possible movement of a given glyph in a given direction
// (horizontally, vertically, or diagonally).
// A vector is needed to represent disjoint ranges, eg, -300..-150, 20..200, 500..750.
// Each pair represents the min/max of a sub-range.
type zones struct {
	exclusions []exclusion

	debugs []zoneDebug // always empty when debug is disabled

	marginLen, marginWeight, pos, posm float32
}

func (zo *zones) initialise(xmin, xmax, marginLen, marginWeight, a0 float32, isXY bool) {
	zo.marginLen = marginLen
	zo.marginWeight = marginWeight
	zo.pos = xmin
	zo.posm = xmax
	zo.exclusions = zo.exclusions[:0]
	var ex exclusion
	if isXY {
		ex = newExclusionWeightedXY(xmin, xmax, 1, a0, 0, 0, 0)
	} else {
		ex = newExclusionWeightedSD(xmin, xmax, 1, a0, 0, 0, 0, 0, false)
	}
	zo.exclusions = append(zo.exclusions, ex)
	zo.exclusions[0].open = true

	if debugMode >= 2 {
		zo.debugs = zo.debugs[:0]
	}
}

func (zo *zones) weightedAxis(axis int, xmin, xmax, f, a0,
	m, xi, ai, c float32, nega bool) {
	if axis < 2 {
		zo.weighted_XY(xmin, xmax, f, a0, m, xi, ai, c, nega)
	} else {
		zo.weighted_SD(xmin, xmax, f, a0, m, xi, ai, c, nega)
	}
}

func (zo *zones) weighted_XY(xmin, xmax, f, a0,
	m, xi, _, c float32, _ bool) {
	zo.insert(newExclusionWeightedXY(xmin, xmax, f, a0, m, xi, c))
}

func (zo *zones) weighted_SD(xmin, xmax, f, a0,
	m, xi, ai, c float32, nega bool) {
	zo.insert(newExclusionWeightedSD(xmin, xmax, f, a0, m, xi, ai, c, nega))
}

func insertExclusion(s []exclusion, i int, x exclusion) []exclusion {
	s = append(s, exclusion{})
	copy(s[i+1:], s[i:])
	s[i] = x
	return s
}

func (zo *zones) insert(e exclusion) {
	if debugMode >= 2 {
		zo.debugs = append(zo.debugs, zoneDebug{excl: e, isDel: false, env: tr.colliderEnv})
	}

	e.x = max(e.x, zo.pos)
	e.xm = min(e.xm, zo.posm)
	if e.x >= e.xm {
		return
	}

	for i := 0; i < len(zo.exclusions) && e.x < e.xm; i++ {
		iter := &zo.exclusions[i]
		oca := e.outcode(iter.x)
		ocb := e.outcode(iter.xm)
		if (oca & ocb) != 0 {
			continue
		}

		switch oca ^ ocb { // What kind of overlap?
		case 0: // e completely covers i
			// split e at iter.x into e1,e2
			// split e2 at iter.mx into e2,e3
			// drop e1 ,i+e2, e=e3
			iter.add(e)
			e.leftTrim(iter.xm)
		case 1: // e overlaps on the rhs of i
			// split i at e.x into i1,i2
			// split e at iter.mx into e1,e2
			// trim i1, insert i2+e1, e=e2
			if !separated(iter.xm, e.x) {
				break
			}
			if separated(iter.x, e.x) {
				zo.exclusions = insertExclusion(zo.exclusions, i, iter.splitAt(e.x))
				i++
				iter = &zo.exclusions[i]
			}
			iter.add(e)
			e.leftTrim(iter.xm)
		case 2: // e overlaps on the lhs of i
			// split e at iter.x into e1,e2
			// split i at e.mx into i1,i2
			// drop e1, insert e2+i1, trim i2
			if !separated(e.xm, iter.x) {
				return
			}
			if separated(e.xm, iter.xm) {
				zo.exclusions = insertExclusion(zo.exclusions, i, iter.splitAt(e.xm))
				iter = &zo.exclusions[i]
			}
			iter.add(e)
			return
		case 3: // i completely covers e
			// split i at e.x into i1,i2
			// split i2 at e.mx into i2,i3
			// insert i1, insert e+i2
			if separated(e.xm, iter.xm) {
				zo.exclusions = insertExclusion(zo.exclusions, i, iter.splitAt(e.xm))
				iter = &zo.exclusions[i]
			}
			zo.exclusions = insertExclusion(zo.exclusions, i, iter.splitAt(e.x))
			i++
			iter = &zo.exclusions[i]
			iter.add(e)
			return
		}

	}
}

func (zo *zones) remove(x, xm float32) {
	if debugMode >= 2 {
		e := exclusion{x: x, xm: xm}
		zo.debugs = append(zo.debugs, zoneDebug{excl: e, isDel: true, env: tr.colliderEnv})
	}

	x = max(x, zo.pos)
	xm = min(xm, zo.posm)
	if x >= xm {
		return
	}

	for i := 0; i < len(zo.exclusions); i++ {
		iter := &zo.exclusions[i]
		oca := iter.outcode(x)
		ocb := iter.outcode(xm)
		if (oca & ocb) != 0 {
			continue
		}

		switch oca ^ ocb { // What kind of overlap?
		case 0: // i completely covers e
			if separated(iter.x, x) {
				zo.exclusions = insertExclusion(zo.exclusions, i, iter.splitAt(x))
				i++
				iter = &zo.exclusions[i]
			}
			fallthrough
			// no break
		case 1: // i overlaps on the rhs of e
			iter.leftTrim(xm)
			return
		case 2: // i overlaps on the lhs of e
			iter.xm = x
			if separated(iter.x, iter.xm) {
				break
			}
			fallthrough
			// no break
		case 3: // e completely covers i
			zo.exclusions = append(zo.exclusions[:i], zo.exclusions[i+1:]...) // erase i
			i--
		}
	}
}

func (zo *zones) excludeWithMargins(xmin, xmax float32, axis int) {
	zo.remove(xmin, xmax)
	zo.weightedAxis(axis, xmin-zo.marginLen, xmin, 0, 0, zo.marginWeight, xmin-zo.marginLen, 0, 0, false)
	zo.weightedAxis(axis, xmax, xmax+zo.marginLen, 0, 0, zo.marginWeight, xmax+zo.marginLen, 0, 0, false)
}

func (zo *zones) findExclusionUnder(x float32) int {
	l, h := 0, len(zo.exclusions)

	for l < h {
		p := (l + h) >> 1
		switch zo.exclusions[p].outcode(x) {
		case 0:
			return p
		case 1:
			h = p
		case 2, 3:
			l = p + 1
		}
	}

	return l
}

func (zo *zones) closest(origin float32) (best, cost float32) {
	var (
		bestC float32 = math.MaxFloat32
		bestX float32
	)

	start := zo.findExclusionUnder(origin)

	// Forward scan looking for lowest cost
	for i := start; i < len(zo.exclusions); i++ {
		if zo.exclusions[i].trackCost(&bestC, &bestX, origin) {
			break
		}
	}

	// Backward scan looking for lowest cost
	//  We start from the exclusion to the immediate left of start since we've
	//  already tested start with the right most scan above.
	for i := start - 1; i >= 0; i-- {
		if zo.exclusions[i].trackCost(&bestC, &bestX, origin) {
			break
		}
	}

	cost = pick(bestC == math.MaxFloat32, -1, bestC)
	return bestX, cost
}

type shiftCollider struct {
	target *Slot // the glyph to fix

	ranges       [4]zones // possible movements in 4 directions (horizontally, vertically, diagonally)
	len          [4]float32
	limit        rect
	currShift    Position
	currOffset   Position
	origin       Position // Base for all relative calculations
	margin       float32
	marginWt     float32
	seqClass     uint16
	seqProxClass uint16
	seqOrder     uint16
}

// initialize the Collider to hold the basic movement limits for the
// target slot, the one we are focusing on fixing.
func (sc *shiftCollider) initSlot(seg *Segment, aSlot *Slot, limit rect, margin, marginWeight float32,
	currShift, currOffset Position, isRTL bool) bool {

	gid := aSlot.glyphID
	glyph := seg.face.getGlyph(gid)
	if glyph == nil {
		return false
	}
	bb := glyph.bbox
	sb := glyph.boxes.slant
	// float sx = aSlot.Position.x + currShift.x;
	// float sy = aSlot.Position.y + currShift.y;
	if currOffset.X != 0. || currOffset.Y != 0. {
		sc.limit = rect{limit.bl.sub(currOffset), limit.tr.sub(currOffset)}
	} else {
		sc.limit = limit
	}

	// For a ShiftCollider, these indices indicate which vector we are moving by:
	// each sc.ranges represents absolute space with respect to the origin of the slot.
	// Thus take into account true origins but subtract the vmin for the slot
	// case 0: // x direction
	mn := sc.limit.bl.X + currOffset.X
	mx := sc.limit.tr.X + currOffset.X
	sc.len[0] = bb.tr.X - bb.bl.X
	a := currOffset.Y + currShift.Y
	sc.ranges[0].initialise(mn, mx, margin, marginWeight, a, true)
	// case 1: // y direction
	mn = sc.limit.bl.Y + currOffset.Y
	mx = sc.limit.tr.Y + currOffset.Y
	sc.len[1] = bb.tr.Y - bb.bl.Y
	a = currOffset.X + currShift.X
	sc.ranges[1].initialise(mn, mx, margin, marginWeight, a, true)
	// case 2: // sum (negatively sloped diagonal boundaries)
	// pick closest x,y limit boundaries in s direction
	shift := currOffset.X + currOffset.Y + currShift.X + currShift.Y
	mn = -2*min(currShift.X-sc.limit.bl.X, currShift.Y-sc.limit.bl.Y) + shift
	mx = 2*min(sc.limit.tr.X-currShift.X, sc.limit.tr.Y-currShift.Y) + shift
	sc.len[2] = sb.tr.X - sb.bl.X
	a = currOffset.X - currOffset.Y + currShift.X - currShift.Y
	sc.ranges[2].initialise(mn, mx, margin/iSQRT2, marginWeight, a, false)
	// case 3: // diff (positively sloped diagonal boundaries)
	// pick closest x,y limit boundaries in d direction
	shift = currOffset.X - currOffset.Y + currShift.X - currShift.Y
	mn = -2*min(currShift.X-sc.limit.bl.X, sc.limit.tr.Y-currShift.Y) + shift
	mx = 2*min(sc.limit.tr.X-currShift.X, currShift.Y-sc.limit.bl.Y) + shift
	sc.len[3] = sb.tr.Y - sb.bl.Y
	a = currOffset.X + currOffset.Y + currShift.X + currShift.Y
	sc.ranges[3].initialise(mn, mx, margin/iSQRT2, marginWeight, a, false)

	sc.target = aSlot
	if !isRTL {
		// For LTR, switch and negate x limits.
		sc.limit.bl.X = -1 * limit.tr.X
		// sc.limit.tr.x = -1 * limit.bl.x;
	}
	sc.currOffset = currOffset
	sc.currShift = currShift
	sc.origin = aSlot.Position.sub(currOffset) // the original anchor position of the glyph

	sc.margin = margin
	sc.marginWt = marginWeight

	c := seg.getCollisionInfo(aSlot)
	sc.seqClass = c.seqClass
	sc.seqProxClass = c.seqProxClass
	sc.seqOrder = c.seqOrder
	return true
}

func sdm(vi, va, mx, my float32, op func(a, b float32) bool) float32 {
	res := 2*mx - vi
	if op(res, vi+2*my) {
		res = va + 2*my
		if op(res, 2*mx-va) {
			res = mx + my
		}
	}
	return res
}

// return b ? v1 : v2
func pick(b bool, v1, v2 float32) float32 {
	if b {
		return v1
	}
	return v2
}

// Mark an area with a cost that can vary along the x or y axis. The region is expressed in terms of the centre of the target glyph in each axis
func (sc *shiftCollider) addBoxSlope(isx bool, box, bb, sb rect, org Position, weight, m float32, minright bool, axis int) {
	switch axis {
	case 0:
		if box.bl.Y < org.Y+bb.tr.Y && box.tr.Y > org.Y+bb.bl.Y && box.width() > 0 {
			a := org.Y + 0.5*(bb.bl.Y+bb.tr.Y)
			c := 0.5 * (bb.bl.X + bb.tr.X)
			if isx {
				sc.ranges[axis].weighted_XY(box.bl.X-c, box.tr.X-c, weight, a, m,
					pick(minright, box.tr.X, box.bl.X)-c, a, 0, false)
			} else {
				sc.ranges[axis].weighted_XY(box.bl.X-c, box.tr.X-c, weight, a, 0, 0, org.Y,
					m*(a*a+sqr(pick(minright, box.tr.Y, box.bl.Y)-0.5*(bb.bl.Y+bb.tr.Y))), false)
			}
		}
	case 1:
		if box.bl.X < org.X+bb.tr.X && box.tr.X > org.X+bb.bl.X && box.height() > 0 {
			a := org.X + 0.5*(bb.bl.X+bb.tr.X)
			c := 0.5 * (bb.bl.Y + bb.tr.Y)
			if isx {
				sc.ranges[axis].weighted_XY(box.bl.Y-c, box.tr.Y-c, weight, a, 0, 0, org.X,
					m*(a*a+sqr(pick(minright, box.tr.X, box.bl.X)-0.5*(bb.bl.X+bb.tr.X))), false)
			} else {
				sc.ranges[axis].weighted_XY(box.bl.Y-c, box.tr.Y-c, weight, a, m,
					pick(minright, box.tr.Y, box.bl.Y)-c, a, 0, false)
			}
		}
	case 2:
		if box.bl.X-box.tr.Y < org.X-org.Y+sb.tr.Y && box.tr.X-box.bl.Y > org.X-org.Y+sb.bl.Y {
			d := org.X - org.Y + 0.5*(sb.bl.Y+sb.tr.Y)
			c := 0.5 * (sb.bl.X + sb.tr.X)
			smax := min(2*box.tr.X-d, 2*box.tr.Y+d)
			smin := max(2*box.bl.X-d, 2*box.bl.Y+d)
			if smin > smax {
				return
			}
			var si float32
			a := d
			if isx {
				si = 2*pick(minright, box.tr.X, box.bl.X) - a
			} else {
				si = 2*pick(minright, box.tr.Y, box.bl.Y) + a
			}
			sc.ranges[axis].weighted_SD(smin-c, smax-c, weight/2, a, m/2, si, 0, 0, isx)
		}
	case 3:
		if box.bl.X+box.bl.Y < org.X+org.Y+sb.tr.X && box.tr.X+box.tr.Y > org.X+org.Y+sb.bl.X {
			s := org.X + org.Y + 0.5*(sb.bl.X+sb.tr.X)
			c := 0.5 * (sb.bl.Y + sb.tr.Y)
			dmax := min(2*box.tr.X-s, s-2*box.bl.Y)
			dmin := max(2*box.bl.X-s, s-2*box.tr.Y)
			if dmin > dmax {
				return
			}
			var di float32
			a := s
			if isx {
				di = 2*pick(minright, box.tr.X, box.bl.X) - a
			} else {
				di = 2*pick(minright, box.tr.Y, box.bl.Y) + a
			}
			sc.ranges[axis].weighted_SD(dmin-c, dmax-c, weight/2, a, m/2, di, 0, 0, !isx)
		}
	}
}

// Mark an area with an absolute cost, making it completely inaccessible.
func (sc *shiftCollider) removeBox(box, bb, sb rect, org Position, axis int) {
	switch axis {
	case 0:
		if box.bl.Y < org.Y+bb.tr.Y && box.tr.Y > org.Y+bb.bl.Y && box.width() > 0 {
			c := 0.5 * (bb.bl.X + bb.tr.X)
			sc.ranges[axis].remove(box.bl.X-c, box.tr.X-c)
		}
	case 1:
		if box.bl.X < org.X+bb.tr.X && box.tr.X > org.X+bb.bl.X && box.height() > 0 {
			c := 0.5 * (bb.bl.Y + bb.tr.Y)
			sc.ranges[axis].remove(box.bl.Y-c, box.tr.Y-c)
		}
	case 2:
		if box.bl.X-box.tr.Y < org.X-org.Y+sb.tr.Y && box.tr.X-box.bl.Y > org.X-org.Y+sb.bl.Y && box.width() > 0 && box.height() > 0 {
			di := org.X - org.Y + sb.bl.Y
			da := org.X - org.Y + sb.tr.Y
			smax := sdm(di, da, box.tr.X, box.tr.Y, func(a, b float32) bool { return a > b })
			smin := sdm(da, di, box.bl.X, box.bl.Y, func(a, b float32) bool { return a < b })
			c := 0.5 * (sb.bl.X + sb.tr.X)
			sc.ranges[axis].remove(smin-c, smax-c)
		}
	case 3:
		if box.bl.X+box.bl.Y < org.X+org.Y+sb.tr.X && box.tr.X+box.tr.Y > org.X+org.Y+sb.bl.X && box.width() > 0 && box.height() > 0 {
			si := org.X + org.Y + sb.bl.X
			sa := org.X + org.Y + sb.tr.X
			dmax := sdm(si, sa, box.tr.X, -box.bl.Y, func(a, b float32) bool { return a > b })
			dmin := sdm(sa, si, box.bl.X, -box.tr.Y, func(a, b float32) bool { return a < b })
			c := 0.5 * (sb.bl.Y + sb.tr.Y)
			sc.ranges[axis].remove(dmin-c, dmax-c)
		}
	}
}

// Adjust the movement limits for the target to avoid having it collide
// with the given neighbor slot. Also determine if there is in fact a collision
// between the target and the given slot.
func (sc *shiftCollider) mergeSlot(seg *Segment, slot *Slot, cslot *slotCollision, currShift Position,
	isAfter, // slot is logically after _target
	sameCluster, isExclusion bool, collides *bool) bool {
	sx := slot.Position.X - sc.origin.X + currShift.X
	sy := slot.Position.Y - sc.origin.Y + currShift.Y
	sd := sx - sy
	ss := sx + sy
	var (
		vmin, vmax, omin, omax, otmin, otmax float32
		cmin, cmax                           float32 // target limits
		torg                                 float32
	)
	glyph := seg.face.getGlyph(slot.glyphID)
	if glyph == nil {
		return false
	}
	bb := glyph.bbox

	// SlotCollision * cslot = seg.collisionInfo(slot);
	var orderFlags uint16
	sameClass := sc.seqProxClass == 0 && cslot.seqClass == sc.seqClass
	if sameCluster && sc.seqClass != 0 && (sameClass || (sc.seqProxClass != 0 && cslot.seqClass == sc.seqProxClass)) {
		// Force the target glyph to be in the specified direction from the slot we're testing.
		orderFlags = sc.seqOrder
	}

	// short circuit if only interested in direct collision and we are out of range
	if orderFlags != 0 || (sx+bb.tr.X+sc.margin >= sc.limit.bl.X && sx+bb.bl.X-sc.margin <= sc.limit.tr.X) ||
		(sy+bb.tr.Y+sc.margin >= sc.limit.bl.Y && sy+bb.bl.Y-sc.margin <= sc.limit.tr.Y) {
		tx := sc.currOffset.X + sc.currShift.X
		ty := sc.currOffset.Y + sc.currShift.Y
		td := tx - ty
		ts := tx + ty
		sb := glyph.boxes.slant
		var tbb, tsb rect
		if tglyph := seg.face.getGlyph(sc.target.glyphID); tglyph != nil {
			tbb, tsb = tglyph.bbox, tglyph.boxes.slant
		}
		seqAboveWt := float32(cslot.seqAboveWt)
		seqBelowWt := float32(cslot.seqBelowWt)
		seqValignWt := float32(cslot.seqValignWt)
		seqValignHt := float32(cslot.seqValignHt)

		var lmargin float32
		// if isAfter, invert orderFlags for diagonal orders.
		if isAfter {
			// invert appropriate bits
			if sameClass {
				orderFlags ^= 0x3F
			} else {
				orderFlags ^= 0x3
			}
			// consider 2 bits at a time, non overlapping. If both bits set, clear them
			orderFlags = orderFlags ^ ((((orderFlags >> 1) & orderFlags) & 0x15) * 3)
		}

		if debugMode >= 2 {
			tr.colliderEnv.sl = slot
		}

		// Process main bounding octabox.
		for i := range sc.ranges {

			switch i {
			case 0: // x direction
				vmin = max(max(bb.bl.X-tbb.tr.X+sx, sb.bl.Y-tsb.tr.Y+ty+sd), sb.bl.X-tsb.tr.X-ty+ss)
				vmax = min(min(bb.tr.X-tbb.bl.X+sx, sb.tr.Y-tsb.bl.Y+ty+sd), sb.tr.X-tsb.bl.X-ty+ss)
				otmin = tbb.bl.Y + ty
				otmax = tbb.tr.Y + ty
				omin = bb.bl.Y + sy
				omax = bb.tr.Y + sy
				torg = sc.currOffset.X
				cmin = sc.limit.bl.X + torg
				cmax = sc.limit.tr.X - tbb.bl.X + tbb.tr.X + torg
				lmargin = sc.margin

			case 1: // y direction
				vmin = max(max(bb.bl.Y-tbb.tr.Y+sy, tsb.bl.Y-sb.tr.Y+tx-sd), sb.bl.X-tsb.tr.X-tx+ss)
				vmax = min(min(bb.tr.Y-tbb.bl.Y+sy, tsb.tr.Y-sb.bl.Y+tx-sd), sb.tr.X-tsb.bl.X-tx+ss)
				otmin = tbb.bl.X + tx
				otmax = tbb.tr.X + tx
				omin = bb.bl.X + sx
				omax = bb.tr.X + sx
				torg = sc.currOffset.Y
				cmin = sc.limit.bl.Y + torg
				cmax = sc.limit.tr.Y - tbb.bl.Y + tbb.tr.Y + torg
				lmargin = sc.margin

			case 2: // sum - moving along the positively-sloped vector, so the boundaries are the
				// negatively-sloped boundaries.
				vmin = max(max(sb.bl.X-tsb.tr.X+ss, 2*(bb.bl.Y-tbb.tr.Y+sy)+td), 2*(bb.bl.X-tbb.tr.X+sx)-td)
				vmax = min(min(sb.tr.X-tsb.bl.X+ss, 2*(bb.tr.Y-tbb.bl.Y+sy)+td), 2*(bb.tr.X-tbb.bl.X+sx)-td)
				otmin = tsb.bl.Y + td
				otmax = tsb.tr.Y + td
				omin = sb.bl.Y + sd
				omax = sb.tr.Y + sd
				torg = sc.currOffset.X + sc.currOffset.Y
				cmin = sc.limit.bl.X + sc.limit.bl.Y + torg
				cmax = sc.limit.tr.X + sc.limit.tr.Y - tsb.bl.X + tsb.tr.X + torg
				lmargin = sc.margin / iSQRT2

			case 3: // diff - moving along the negatively-sloped vector, so the boundaries are the
				// positively-sloped boundaries.
				vmin = max(max(sb.bl.Y-tsb.tr.Y+sd, 2*(bb.bl.X-tbb.tr.X+sx)-ts), -2*(bb.tr.Y-tbb.bl.Y+sy)+ts)
				vmax = min(min(sb.tr.Y-tsb.bl.Y+sd, 2*(bb.tr.X-tbb.bl.X+sx)-ts), -2*(bb.bl.Y-tbb.tr.Y+sy)+ts)
				otmin = tsb.bl.X + ts
				otmax = tsb.tr.X + ts
				omin = sb.bl.X + ss
				omax = sb.tr.X + ss
				torg = sc.currOffset.X - sc.currOffset.Y
				cmin = sc.limit.bl.X - sc.limit.tr.Y + torg
				cmax = sc.limit.tr.X - sc.limit.bl.Y - tsb.bl.Y + tsb.tr.Y + torg
				lmargin = sc.margin / iSQRT2
			}

			if debugMode >= 2 {
				tr.colliderEnv.val = -1
			}

			if orderFlags != 0 {
				org := Position{tx, ty}
				xminf := sc.limit.bl.X + sc.currOffset.X + tbb.bl.X
				xpinf := sc.limit.tr.X + sc.currOffset.X + tbb.tr.X
				ypinf := sc.limit.tr.Y + sc.currOffset.Y + tbb.tr.Y
				yminf := sc.limit.bl.Y + sc.currOffset.Y + tbb.bl.Y
				switch orderFlags {
				case seqOrderRIGHTUP:
					r1Xedge := float32(cslot.seqAboveXoff) + 0.5*(bb.bl.X+bb.tr.X) + sx
					r3Xedge := float32(cslot.seqBelowXlim) + bb.tr.X + sx + 0.5*(tbb.tr.X-tbb.bl.X)
					r2Yedge := 0.5*(bb.bl.Y+bb.tr.Y) + sy

					// region 1
					// DBGTAG(1x) means the regions are up and right
					if debugMode >= 2 {
						tr.colliderEnv.val = -11
					}
					sc.addBoxSlope(true, rect{Position{xminf, r2Yedge}, Position{r1Xedge, ypinf}},
						tbb, tsb, org, 0, seqAboveWt, true, i)
					// region 2
					if debugMode >= 2 {
						tr.colliderEnv.val = -12
					}
					sc.removeBox(rect{Position{xminf, yminf}, Position{r3Xedge, r2Yedge}}, tbb, tsb, org, i)
					// region 3, which end is zero is irrelevant since m weight is 0
					if debugMode >= 2 {
						tr.colliderEnv.val = -13
					}
					sc.addBoxSlope(true, rect{Position{r3Xedge, yminf}, Position{xpinf, r2Yedge - seqValignHt}},
						tbb, tsb, org, seqBelowWt, 0, true, i)
					// region 4
					if debugMode >= 2 {
						tr.colliderEnv.val = -14
					}
					sc.addBoxSlope(false, rect{Position{sx + bb.bl.X, r2Yedge}, Position{xpinf, r2Yedge + seqValignHt}},
						tbb, tsb, org, 0, seqValignWt, true, i)
					// region 5
					if debugMode >= 2 {
						tr.colliderEnv.val = -15
					}
					sc.addBoxSlope(false, rect{Position{sx + bb.bl.X, r2Yedge - seqValignHt}, Position{xpinf, r2Yedge}},
						tbb, tsb, org, seqBelowWt, seqValignWt, false, i)
				case seqOrderLEFTDOWN:
					r1Xedge := 0.5*(bb.bl.X+bb.tr.X) + float32(cslot.seqAboveXoff) + sx
					r3Xedge := bb.bl.X - float32(cslot.seqBelowXlim) + sx - 0.5*(tbb.tr.X-tbb.bl.X)
					r2Yedge := 0.5*(bb.bl.Y+bb.tr.Y) + sy
					// DBGTAG(2x) means the regions are up and right
					// region 1
					if debugMode >= 2 {
						tr.colliderEnv.val = -21
					}
					sc.addBoxSlope(true, rect{Position{r1Xedge, yminf}, Position{xpinf, r2Yedge}},
						tbb, tsb, org, 0, seqAboveWt, false, i)
					// region 2
					if debugMode >= 2 {
						tr.colliderEnv.val = -22
					}
					sc.removeBox(rect{Position{r3Xedge, r2Yedge}, Position{xpinf, ypinf}}, tbb, tsb, org, i)
					// region 3
					if debugMode >= 2 {
						tr.colliderEnv.val = -23
					}
					sc.addBoxSlope(true, rect{Position{xminf, r2Yedge - seqValignHt}, Position{r3Xedge, ypinf}},
						tbb, tsb, org, seqBelowWt, 0, false, i)
					// region 4
					if debugMode >= 2 {
						tr.colliderEnv.val = -24
					}
					sc.addBoxSlope(false, rect{Position{xminf, r2Yedge}, Position{sx + bb.tr.X, r2Yedge + seqValignHt}},
						tbb, tsb, org, 0, seqValignWt, true, i)
					// region 5
					if debugMode >= 2 {
						tr.colliderEnv.val = -25
					}
					sc.addBoxSlope(false, rect{
						Position{xminf, r2Yedge - seqValignHt},
						Position{sx + bb.tr.X, r2Yedge},
					}, tbb, tsb, org, seqBelowWt, seqValignWt, false, i)
				case seqOrderNOABOVE: // enforce neighboring glyph being above
					if debugMode >= 2 {
						tr.colliderEnv.val = -31
					}
					sc.removeBox(rect{
						Position{bb.bl.X - tbb.tr.X + sx, sy + bb.tr.Y},
						Position{bb.tr.X - tbb.bl.X + sx, ypinf},
					}, tbb, tsb, org, i)
				case seqOrderNOBELOW: // enforce neighboring glyph being below
					if debugMode >= 2 {
						tr.colliderEnv.val = -32
					}
					sc.removeBox(rect{
						Position{bb.bl.X - tbb.tr.X + sx, yminf},
						Position{bb.tr.X - tbb.bl.X + sx, sy + bb.bl.Y},
					}, tbb, tsb, org, i)
				case seqOrderNOLEFT: // enforce neighboring glyph being to the left
					if debugMode >= 2 {
						tr.colliderEnv.val = -33
					}
					sc.removeBox(rect{
						Position{xminf, bb.bl.Y - tbb.tr.Y + sy},
						Position{bb.bl.X - tbb.tr.X + sx, bb.tr.Y - tbb.bl.Y + sy},
					}, tbb, tsb, org, i)
				case seqOrderNORIGHT: // enforce neighboring glyph being to the right
					if debugMode >= 2 {
						tr.colliderEnv.val = -34
					}
					sc.removeBox(rect{
						Position{bb.tr.X - tbb.bl.X + sx, bb.bl.Y - tbb.tr.Y + sy},
						Position{xpinf, bb.tr.Y - tbb.bl.Y + sy},
					}, tbb, tsb, org, i)
				}
			}

			if vmax < cmin-lmargin || vmin > cmax+lmargin || omax < otmin-lmargin || omin > otmax+lmargin {
				continue
			}
			// Process sub-boxes that are defined for this glyph.
			// We only need to do this if there was in fact a collision with the main octabox.
			numsub := len(glyph.boxes.subBboxes)
			if numsub > 0 {
				anyhits := false
				for j := range glyph.boxes.slantSubBboxes {
					sbb := glyph.boxes.subBboxes[j]
					ssb := glyph.boxes.slantSubBboxes[j]
					switch i {
					case 0: // x
						vmin = max(max(sbb.bl.X-tbb.tr.X+sx, ssb.bl.Y-tsb.tr.Y+sd+ty), ssb.bl.X-tsb.tr.X+ss-ty)
						vmax = min(min(sbb.tr.X-tbb.bl.X+sx, ssb.tr.Y-tsb.bl.Y+sd+ty), ssb.tr.X-tsb.bl.X+ss-ty)
						omin = sbb.bl.Y + sy
						omax = sbb.tr.Y + sy
					case 1: // y
						vmin = max(max(sbb.bl.Y-tbb.tr.Y+sy, tsb.bl.Y-ssb.tr.Y-sd+tx), ssb.bl.X-tsb.tr.X+ss-tx)
						vmax = min(min(sbb.tr.Y-tbb.bl.Y+sy, tsb.tr.Y-ssb.bl.Y-sd+tx), ssb.tr.X-tsb.bl.X+ss-tx)
						omin = sbb.bl.X + sx
						omax = sbb.tr.X + sx
					case 2: // sum
						vmin = max(max(ssb.bl.X-tsb.tr.X+ss, 2*(sbb.bl.Y-tbb.tr.Y+sy)+td), 2*(sbb.bl.X-tbb.tr.X+sx)-td)
						vmax = min(min(ssb.tr.X-tsb.bl.X+ss, 2*(sbb.tr.Y-tbb.bl.Y+sy)+td), 2*(sbb.tr.X-tbb.bl.X+sx)-td)
						omin = ssb.bl.Y + sd
						omax = ssb.tr.Y + sd
					case 3: // diff
						vmin = max(max(ssb.bl.Y-tsb.tr.Y+sd, 2*(sbb.bl.X-tbb.tr.X+sx)-ts), -2*(sbb.tr.Y-tbb.bl.Y+sy)+ts)
						vmax = min(min(ssb.tr.Y-tsb.bl.Y+sd, 2*(sbb.tr.X-tbb.bl.X+sx)-ts), -2*(sbb.bl.Y-tbb.tr.Y+sy)+ts)
						omin = ssb.bl.X + ss
						omax = ssb.tr.X + ss
					}
					if vmax < cmin-lmargin || vmin > cmax+lmargin || omax < otmin-lmargin || omin > otmax+lmargin {
						continue
					}

					if debugMode >= 2 {
						tr.colliderEnv.val = j
					}

					if omin > otmax {
						sc.ranges[i].weightedAxis(i, vmin-lmargin, vmax+lmargin, 0, 0, 0, 0, 0,
							sqr(lmargin-omin+otmax)*sc.marginWt, false)
					} else if omax < otmin {
						sc.ranges[i].weightedAxis(i, vmin-lmargin, vmax+lmargin, 0, 0, 0, 0, 0,
							sqr(lmargin-otmin+omax)*sc.marginWt, false)
					} else {
						sc.ranges[i].excludeWithMargins(vmin, vmax, i)
					}
					anyhits = true
				}
				if anyhits {
					*collides = true
				}
			} else { // no sub-boxes

				if debugMode >= 2 {
					tr.colliderEnv.val = -1
				}

				*collides = true
				if omin > otmax {
					sc.ranges[i].weightedAxis(i, vmin-lmargin, vmax+lmargin, 0, 0, 0, 0, 0,
						sqr(lmargin-omin+otmax)*sc.marginWt, false)
				} else if omax < otmin {
					sc.ranges[i].weightedAxis(i, vmin-lmargin, vmax+lmargin, 0, 0, 0, 0, 0,
						sqr(lmargin-otmin+omax)*sc.marginWt, false)
				} else {
					sc.ranges[i].excludeWithMargins(vmin, vmax, i)
				}

			}
		}
	}
	res := true
	if cslot.exclGlyph > 0 && int(cslot.exclGlyph) < len(seg.face.glyphs) && !isExclusion {
		// Set up the bogus slot representing the exclusion glyph.
		exclSlot := seg.newSlot()
		if exclSlot == nil {
			return res
		}
		exclSlot.setGlyph(seg, cslot.exclGlyph)
		exclOrigin := slot.Position.add(cslot.exclOffset)
		exclSlot.setPosition(exclOrigin)
		var exclInfo slotCollision
		exclInfo.init(seg, exclSlot)
		resExl := sc.mergeSlot(seg, exclSlot, &exclInfo, currShift, isAfter, sameCluster, true, collides)
		res = res && resExl
		seg.freeSlot(exclSlot)
	}
	return res
}

// Figure out where to move the target glyph to, and return the amount to shift by.
func (sc *shiftCollider) resolve(seg *Segment) (Position, bool) {
	totalCost := float32(math.MaxFloat32) / 2
	var resultPos Position
	bestAxis := -1

	if debugMode >= 2 {
		tr.addCollisionMove(sc, seg)
	}

	isCol := true
	for i := range sc.ranges {
		var bestPos, tbase float32
		// Calculate the margin depending on whether we are moving diagonally or not:
		switch i {
		case 0: // x direction
			tbase = sc.currOffset.X
		case 1: // y direction
			tbase = sc.currOffset.Y
		case 2: // sum (negatively-sloped diagonals)
			tbase = sc.currOffset.X + sc.currOffset.Y
		case 3: // diff (positively-sloped diagonals)
			tbase = sc.currOffset.X - sc.currOffset.Y
		}
		var testp Position
		tmp, bestCost := sc.ranges[i].closest(0)
		bestPos = tmp - tbase // Get the best relative position

		if debugMode >= 2 {
			tr.addCollisionVector(sc, seg, i, tbase, bestCost, bestPos)
		}

		if bestCost >= 0.0 {
			isCol = false
			switch i {
			case 0:
				testp = Position{bestPos, sc.currShift.Y}
			case 1:
				testp = Position{sc.currShift.X, bestPos}
			case 2:
				testp = Position{0.5 * (sc.currShift.X - sc.currShift.Y + bestPos), 0.5 * (sc.currShift.Y - sc.currShift.X + bestPos)}
			case 3:
				testp = Position{0.5 * (sc.currShift.X + sc.currShift.Y + bestPos), 0.5 * (sc.currShift.X + sc.currShift.Y - bestPos)}
			}
			if bestCost < totalCost-0.01 {
				totalCost = bestCost
				resultPos = testp
				bestAxis = i
			}
		}
	} // end of loop over 4 directions

	if debugMode >= 2 {
		tr.endCollisionMove(resultPos, bestAxis, isCol)
	}

	return resultPos, isCol
}

type kernCollider struct {
	target *Slot     // the glyph to fix
	edges  []float32 // edges of horizontal slices

	// always empty outside of debug mode
	seg       *Segment
	nearEdges []float32
	slotNear  []*Slot

	limit      rect
	offsetPrev Position // kern from a previous pass
	margin     float32

	sliceWidth float32 // width of each slice
	mingap     float32
	xbound     float32 // max or min edge
	miny, maxy float32 // y-coordinates offset by global slot position
	hit        bool
}

func newKernCollider() *kernCollider {
	return &kernCollider{miny: -1e38, maxy: 1e38}
}

func localmax(al, au, bl, bu, x float32) float32 {
	if al < bl {
		if au < bu {
			return min(au, x)
		}
	} else if au > bu {
		return min(bl, x)
	}
	return x
}

func localmin(al, au, bl, bu, x float32) float32 {
	if bl > al {
		if bu > au {
			return max(bl, x)
		}
	} else if au > bu {
		return max(al, x)
	}
	return x
}

// Return the given edge of the glyph at height y, taking any slant box into account.
func getEdge(seg *Segment, s *Slot, shift Position, y, width, margin float32, isRight bool) float32 {
	sx := s.Position.X + shift.X
	sy := s.Position.Y + shift.Y
	res := pick(isRight, -1e38, 1e38)

	glyph := seg.face.getGlyph(s.glyphID)
	if len(glyph.boxes.subBboxes) != 0 {
		for i := range glyph.boxes.subBboxes {
			sbb := glyph.boxes.subBboxes[i]
			ssb := glyph.boxes.slantSubBboxes[i]
			if sy+sbb.bl.Y-margin > y+width/2 || sy+sbb.tr.Y+margin < y-width/2 {
				continue
			}
			if isRight {
				x := sx + sbb.tr.X + margin
				if x > res {
					td := sx - sy + ssb.tr.Y + margin + y
					ts := sx + sy + ssb.tr.X + margin - y
					x = localmax(td-width/2, td+width/2, ts-width/2, ts+width/2, x)
					if x > res {
						res = x
					}
				}
			} else {
				x := sx + sbb.bl.X - margin
				if x < res {
					td := sx - sy + ssb.bl.Y - margin + y
					ts := sx + sy + ssb.bl.X - margin - y
					x = localmin(td-width/2, td+width/2, ts-width/2, ts+width/2, x)
					if x < res {
						res = x
					}
				}
			}
		}
	} else {
		bb := glyph.bbox
		sb := glyph.boxes.slant
		if sy+bb.bl.Y-margin > y+width/2 || sy+bb.tr.Y+margin < y-width/2 {
			return res
		}
		td := sx - sy + y
		ts := sx + sy - y
		if isRight {
			res = localmax(td+sb.tr.Y-width/2, td+sb.tr.Y+width/2, ts+sb.tr.X-width/2, ts+sb.tr.X+width/2, sx+bb.tr.X) + margin
		} else {
			res = localmin(td+sb.bl.Y-width/2, td+sb.bl.Y+width/2, ts+sb.bl.X-width/2, ts+sb.bl.X+width/2, sx+bb.bl.X) - margin
		}
	}
	return res
}

// append n copies of val to dst
func insertN(dst []float32, n int, val float32) []float32 {
	L := len(dst)
	dst = append(dst, make([]float32, n)...)
	subSlice := dst[L : L+n]
	for i := range subSlice {
		subSlice[i] = val
	}
	return dst
}

func (kc *kernCollider) initSlot(seg *Segment, aSlot *Slot, limit rect, margin float32,
	currShift, offsetPrev Position, isRTL bool, ymin, ymax float32) bool {
	base := aSlot.findRoot()
	if margin < 10 {
		margin = 10
	}

	kc.limit = limit
	kc.offsetPrev = offsetPrev // kern from a previous pass

	var numSlices int
	// Calculate the height of the glyph and how many horizontal slices to use.
	if kc.maxy >= 1e37 {
		kc.sliceWidth = margin / 1.5
		kc.maxy = ymax + margin
		kc.miny = ymin - margin
		numSlices = int((kc.maxy-kc.miny+2)/(kc.sliceWidth/1.5) + 1.) // +2 helps with rounding errors
		kc.edges = kc.edges[:0]
		kc.edges = insertN(kc.edges, numSlices, pick(isRTL, 1e38, -1e38))
		kc.xbound = pick(isRTL, 1e38, -1e38)
	} else if kc.maxy != ymax || kc.miny != ymin {
		if kc.miny != ymin {
			numSlices = int((ymin-margin-kc.miny)/kc.sliceWidth - 1)
			kc.miny += float32(numSlices) * kc.sliceWidth
			if numSlices < 0 {
				kc.edges = append(insertN(nil, -numSlices, pick(isRTL, 1e38, -1e38)), kc.edges...)
			} else if numSlices < len(kc.edges) { // this shouldn't fire since we always grow the range
				kc.edges = kc.edges[numSlices:]
			}
		}
		if kc.maxy != ymax {
			numSlices = int((ymax+margin-kc.miny)/kc.sliceWidth + 1)
			kc.maxy = float32(numSlices)*kc.sliceWidth + kc.miny
			if numSlices > len(kc.edges) {
				kc.edges = insertN(kc.edges, numSlices-len(kc.edges), pick(isRTL, 1e38, -1e38))
			} else if numSlices < len(kc.edges) { // this shouldn't fire since we always grow the range
				kc.edges = kc.edges[:numSlices]
			}
		}
		goto done
	}
	numSlices = len(kc.edges)

	if debugMode >= 2 {
		kc.seg = seg
		kc.slotNear = make([]*Slot, numSlices)
		kc.nearEdges = make([]float32, numSlices)
		for i := range kc.nearEdges {
			kc.nearEdges[i] = pick(isRTL, -1e38, +1e38)
		}
	}

	// Determine the trailing edge of each slice (ie, left edge for a RTL glyph).
	for s := base; s != nil; s = s.nextInCluster(s) {
		c := seg.getCollisionInfo(s)
		if int(s.glyphID) >= len(seg.face.glyphs) {
			return false
		}
		bs := seg.face.getGlyph(s.glyphID).bbox
		x := s.Position.X + c.shift.X + (pick(isRTL, bs.bl.X, bs.tr.X))
		// Loop over slices.
		// Note smin might not be zero if glyph s is not at the bottom of the cluster; similarly for smax.
		toffset := c.shift.Y - kc.miny + 1 + s.Position.Y
		smin := int(max(0, (bs.bl.Y+toffset)/kc.sliceWidth))
		smax := int(min(float32(numSlices-1), (bs.tr.Y+toffset)/kc.sliceWidth+1))
		for i := smin; i <= smax; i++ {
			y := kc.miny - 1 + (float32(i)+0.5)*kc.sliceWidth // vertical center of slice
			if isRTL && x < kc.edges[i] {
				t := getEdge(seg, s, c.shift, y, kc.sliceWidth, margin, false)
				if t < kc.edges[i] {
					kc.edges[i] = t
					if t < kc.xbound {
						kc.xbound = t
					}
				}
			} else if !isRTL && x > kc.edges[i] {
				t := getEdge(seg, s, c.shift, y, kc.sliceWidth, margin, true)
				if t > kc.edges[i] {
					kc.edges[i] = t
					if t > kc.xbound {
						kc.xbound = t
					}
				}
			}
		}
	}
done:
	kc.mingap = 1e37 // less than 1e38 s.t. 1e38-_mingap is really big
	kc.target = aSlot
	kc.margin = margin
	return true
}

// Determine how much the target slot needs to kern away from the given slot.
// In other words, merge information from given slot's position with what the target slot knows
// about how it can kern.
// Return false if we know there is no collision, true if we think there might be one.
func (kc *kernCollider) mergeSlot(seg *Segment, slot *Slot, currShift Position, currSpace float32, isRTL bool) bool {
	rtl := pick(isRTL, 1, -1)
	glyph := seg.face.getGlyph(slot.glyphID)
	if glyph == nil {
		return false
	}
	bb := glyph.bbox
	sx := slot.Position.X + currShift.X
	x := (sx + pick(rtl > 0, bb.tr.X, bb.bl.X)) * rtl
	// this isn't going to reduce _mingap so skip
	if kc.hit && x < rtl*(kc.xbound-kc.mingap-currSpace) {
		return false
	}

	sy := slot.Position.Y + currShift.Y
	smin := int(max(1, (bb.bl.Y+(1-kc.miny+sy))/kc.sliceWidth+1) - 1)
	smax := int(min(float32(len(kc.edges)-2), (bb.tr.Y+(1-kc.miny+sy))/kc.sliceWidth+1) + 1)
	if smin > smax {
		return false
	}
	collides := false
	nooverlap := true

	for i := smin; i <= smax; i++ {
		here := kc.edges[i] * rtl
		if here > 9e37 {
			continue
		}
		if !kc.hit || x > here-kc.mingap-currSpace {
			y := (kc.miny - 1 + (float32(i)+0.5)*kc.sliceWidth) // vertical center of slice
			// 2 * currSpace to account for the space that is already separating them and the space we want to add
			m := getEdge(seg, slot, currShift, y, kc.sliceWidth, 0., rtl > 0)*rtl + 2*currSpace
			if m < -8e37 { // only true if the glyph has a gap in it
				continue
			}
			nooverlap = false
			t := here - m
			// kc.mingap is positive to shrink
			if t < kc.mingap || (!kc.hit && !collides) {
				kc.mingap = t
				collides = true
			}

			if debugMode >= 2 {
				// Debugging - remember the closest neighboring edge for this slice.
				if m > rtl*kc.nearEdges[i] {
					kc.slotNear[i] = slot
					kc.nearEdges[i] = m * rtl
				}
			}

		} else {
			nooverlap = false
		}
	}
	if nooverlap {
		kc.mingap = max(kc.mingap, kc.xbound-rtl*(currSpace+kc.margin+x))
	}
	if collides && !nooverlap {
		kc.hit = true
	}
	return collides || nooverlap // note that true is not a necessarily reliable value
}

// Return the amount to kern by.
func (kc *kernCollider) resolve(isRTL bool) Position {
	resultNeeded := pick(isRTL, -1, 1) * kc.mingap
	result := min(kc.limit.tr.X-kc.offsetPrev.X, max(resultNeeded, kc.limit.bl.X-kc.offsetPrev.X))

	if debugMode >= 2 {
		tr.addKern(kc, kc.seg, result, resultNeeded)
	}

	return Position{result, 0.}
}

func (kc *kernCollider) shift(mv Position, isRTL bool) {
	for i := range kc.edges {
		kc.edges[i] += mv.X
	}
	kc.xbound += pick(isRTL, -1, 1) * mv.X
}
