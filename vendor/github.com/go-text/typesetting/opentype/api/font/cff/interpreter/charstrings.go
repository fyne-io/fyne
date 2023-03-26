// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package psinterpreter

import (
	"errors"
	"fmt"

	"github.com/go-text/typesetting/opentype/api"
)

// PathBounds represents a control bounds for
// a glyph outline (in font units).
type PathBounds struct {
	Min, Max Point
}

// Enlarge enlarges the bounds to include pt
func (b *PathBounds) Enlarge(pt Point) {
	if pt.X < b.Min.X {
		b.Min.X = pt.X
	}
	if pt.X > b.Max.X {
		b.Max.X = pt.X
	}
	if pt.Y < b.Min.Y {
		b.Min.Y = pt.Y
	}
	if pt.Y > b.Max.Y {
		b.Max.Y = pt.Y
	}
}

// ToExtents converts a path bounds to the corresponding glyph extents.
func (b *PathBounds) ToExtents() api.GlyphExtents {
	return api.GlyphExtents{
		XBearing: float32(b.Min.X),
		YBearing: float32(b.Max.Y),
		Width:    float32(b.Max.X - b.Min.X),
		Height:   float32(b.Min.Y - b.Max.Y),
	}
}

// Point is a 2D Point in font units.
type Point struct{ X, Y int32 }

// Move translates the Point.
func (p *Point) Move(dx, dy int32) {
	p.X += dx
	p.Y += dy
}

func (p Point) toSP() api.SegmentPoint {
	return api.SegmentPoint{X: float32(p.X), Y: float32(p.Y)}
}

// CharstringReader provides implementation
// of the operators found in a font charstring.
type CharstringReader struct {
	// Acumulated segments for the glyph outlines
	Segments []api.Segment
	// Acumulated bounds for the glyph outlines
	Bounds PathBounds

	vstemCount   int32
	hstemCount   int32
	hintmaskSize int32

	CurrentPoint Point
	firstPoint   Point // first point in path, required to check if a path is closed
	isPathOpen   bool

	seenHintmask bool

	// bounds for an empty path is {0,0,0,0}
	// however, for the first point in the path,
	// we must not compare the coordinates with {0,0,0,0}
	seenPoint bool
}

// enlarges the current bounds to include the Point (x,y).
func (out *CharstringReader) updateBounds(pt Point) {
	if !out.seenPoint {
		out.Bounds.Min, out.Bounds.Max = pt, pt
		out.seenPoint = true
		return
	}
	out.Bounds.Enlarge(pt)
}

func (out *CharstringReader) Hstem(state *Machine) {
	out.hstemCount += state.ArgStack.Top / 2
}

func (out *CharstringReader) Vstem(state *Machine) {
	out.vstemCount += state.ArgStack.Top / 2
}

func (out *CharstringReader) determineHintmaskSize(state *Machine) {
	if !out.seenHintmask {
		out.vstemCount += state.ArgStack.Top / 2
		out.hintmaskSize = (out.hstemCount + out.vstemCount + 7) >> 3
		out.seenHintmask = true
	}
}

func (out *CharstringReader) Hintmask(state *Machine) {
	out.determineHintmaskSize(state)
	state.SkipBytes(out.hintmaskSize)
}

func (out *CharstringReader) move(pt Point) {
	out.ensureClosePath()

	out.CurrentPoint.Move(pt.X, pt.Y)
	out.isPathOpen = false
	out.firstPoint = out.CurrentPoint
	out.Segments = append(out.Segments, api.Segment{
		Op:   api.SegmentOpMoveTo,
		Args: [3]api.SegmentPoint{out.CurrentPoint.toSP()},
	})
}

// pt is in absolute coordinates
func (out *CharstringReader) line(pt Point) {
	if !out.isPathOpen {
		out.isPathOpen = true
		out.updateBounds(out.CurrentPoint)
	}
	out.CurrentPoint = pt
	out.updateBounds(pt)
	out.Segments = append(out.Segments, api.Segment{
		Op:   api.SegmentOpLineTo,
		Args: [3]api.SegmentPoint{pt.toSP()},
	})
}

func (out *CharstringReader) curve(pt1, pt2, pt3 Point) {
	if !out.isPathOpen {
		out.isPathOpen = true
		out.updateBounds(out.CurrentPoint)
	}
	/* include control Points */
	out.updateBounds(pt1)
	out.updateBounds(pt2)
	out.CurrentPoint = pt3
	out.updateBounds(pt3)
	out.Segments = append(out.Segments, api.Segment{
		Op:   api.SegmentOpCubeTo,
		Args: [3]api.SegmentPoint{pt1.toSP(), pt2.toSP(), pt3.toSP()},
	})
}

func (out *CharstringReader) doubleCurve(pt1, pt2, pt3, pt4, pt5, pt6 Point) {
	out.curve(pt1, pt2, pt3)
	out.curve(pt4, pt5, pt6)
}

func (out *CharstringReader) ensureClosePath() {
	if out.firstPoint != out.CurrentPoint {
		out.Segments = append(out.Segments, api.Segment{
			Op:   api.SegmentOpLineTo,
			Args: [3]api.SegmentPoint{out.firstPoint.toSP()},
		})
	}
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// ------------------------------------------------------------

// LocalSubr pops the subroutine index and call it
func LocalSubr(state *Machine) error {
	if state.ArgStack.Top < 1 {
		return errors.New("invalid callsubr operator (empty stack)")
	}
	index := state.ArgStack.Pop()
	return state.CallSubroutine(index, true)
}

// GlobalSubr pops the subroutine index and call it
func GlobalSubr(state *Machine) error {
	if state.ArgStack.Top < 1 {
		return errors.New("invalid callgsubr operator (empty stack)")
	}
	index := state.ArgStack.Pop()
	return state.CallSubroutine(index, false)
}

// ClosePath closes the current contour, adding
// a segment to the first point if needed.
func (out *CharstringReader) ClosePath() {
	out.ensureClosePath()
	out.isPathOpen = false
}

func (out *CharstringReader) Rmoveto(state *Machine) error {
	if state.ArgStack.Top < 2 {
		return errors.New("invalid rmoveto operator")
	}
	y := state.ArgStack.Pop()
	x := state.ArgStack.Pop()
	out.move(Point{x, y})
	return nil
}

func (out *CharstringReader) Vmoveto(state *Machine) error {
	if state.ArgStack.Top < 1 {
		return errors.New("invalid vmoveto operator")
	}
	y := state.ArgStack.Pop()
	out.move(Point{0, y})
	return nil
}

func (out *CharstringReader) Hmoveto(state *Machine) error {
	if state.ArgStack.Top < 1 {
		return errors.New("invalid hmoveto operator")
	}
	x := state.ArgStack.Pop()
	out.move(Point{x, 0})
	return nil
}

func (out *CharstringReader) Rlineto(state *Machine) {
	for i := int32(0); i+2 <= state.ArgStack.Top; i += 2 {
		newPoint := out.CurrentPoint
		newPoint.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
		out.line(newPoint)
	}
	state.ArgStack.Clear()
}

func (out *CharstringReader) Hlineto(state *Machine) {
	var i int32
	for ; i+2 <= state.ArgStack.Top; i += 2 {
		newPoint := out.CurrentPoint
		newPoint.X += state.ArgStack.Vals[i]
		out.line(newPoint)
		newPoint.Y += state.ArgStack.Vals[i+1]
		out.line(newPoint)
	}
	if i < state.ArgStack.Top {
		newPoint := out.CurrentPoint
		newPoint.X += state.ArgStack.Vals[i]
		out.line(newPoint)
	}
}

func (out *CharstringReader) Vlineto(state *Machine) {
	var i int32
	for ; i+2 <= state.ArgStack.Top; i += 2 {
		newPoint := out.CurrentPoint
		newPoint.Y += state.ArgStack.Vals[i]
		out.line(newPoint)
		newPoint.X += state.ArgStack.Vals[i+1]
		out.line(newPoint)
	}
	if i < state.ArgStack.Top {
		newPoint := out.CurrentPoint
		newPoint.Y += state.ArgStack.Vals[i]
		out.line(newPoint)
	}
}

// RelativeCurveTo draws a curve with controls points computed from
// the current point and `arg1`, `arg2`, `arg3`
func (out *CharstringReader) RelativeCurveTo(arg1, arg2, arg3 Point) {
	pt1 := out.CurrentPoint
	pt1.Move(arg1.X, arg1.Y)
	pt2 := pt1
	pt2.Move(arg2.X, arg2.Y)
	pt3 := pt2
	pt3.Move(arg3.X, arg3.Y)
	out.curve(pt1, pt2, pt3)
}

func (out *CharstringReader) Rrcurveto(state *Machine) {
	for i := int32(0); i+6 <= state.ArgStack.Top; i += 6 {
		out.RelativeCurveTo(
			Point{state.ArgStack.Vals[i], state.ArgStack.Vals[i+1]},
			Point{state.ArgStack.Vals[i+2], state.ArgStack.Vals[i+3]},
			Point{state.ArgStack.Vals[i+4], state.ArgStack.Vals[i+5]},
		)
	}
}

func (out *CharstringReader) Hhcurveto(state *Machine) {
	var (
		i   int32
		pt1 = out.CurrentPoint
	)
	if (state.ArgStack.Top & 1) != 0 {
		pt1.Y += (state.ArgStack.Vals[i])
		i++
	}
	for ; i+4 <= state.ArgStack.Top; i += 4 {
		pt1.X += state.ArgStack.Vals[i]
		pt2 := pt1
		pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
		pt3 := pt2
		pt3.X += state.ArgStack.Vals[i+3]
		out.curve(pt1, pt2, pt3)
		pt1 = out.CurrentPoint
	}
}

func (out *CharstringReader) Vhcurveto(state *Machine) {
	var i int32
	if (state.ArgStack.Top % 8) >= 4 {
		pt1 := out.CurrentPoint
		pt1.Y += state.ArgStack.Vals[i]
		pt2 := pt1
		pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
		pt3 := pt2
		pt3.X += state.ArgStack.Vals[i+3]
		i += 4

		for ; i+8 <= state.ArgStack.Top; i += 8 {
			out.curve(pt1, pt2, pt3)
			pt1 = out.CurrentPoint
			pt1.X += (state.ArgStack.Vals[i])
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
			pt3 = pt2
			pt3.Y += (state.ArgStack.Vals[i+3])
			out.curve(pt1, pt2, pt3)

			pt1 = pt3
			pt1.Y += (state.ArgStack.Vals[i+4])
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+5], state.ArgStack.Vals[i+6])
			pt3 = pt2
			pt3.X += (state.ArgStack.Vals[i+7])
		}
		if i < state.ArgStack.Top {
			pt3.Y += (state.ArgStack.Vals[i])
		}
		out.curve(pt1, pt2, pt3)
	} else {
		for ; i+8 <= state.ArgStack.Top; i += 8 {
			pt1 := out.CurrentPoint
			pt1.Y += (state.ArgStack.Vals[i])
			pt2 := pt1
			pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
			pt3 := pt2
			pt3.X += (state.ArgStack.Vals[i+3])
			out.curve(pt1, pt2, pt3)

			pt1 = pt3
			pt1.X += (state.ArgStack.Vals[i+4])
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+5], state.ArgStack.Vals[i+6])
			pt3 = pt2
			pt3.Y += (state.ArgStack.Vals[i+7])
			if (state.ArgStack.Top-i < 16) && ((state.ArgStack.Top & 1) != 0) {
				pt3.X += (state.ArgStack.Vals[i+8])
			}
			out.curve(pt1, pt2, pt3)
		}
	}
}

func (out *CharstringReader) Hvcurveto(state *Machine) {
	var i int32
	if (state.ArgStack.Top % 8) >= 4 {
		pt1 := out.CurrentPoint
		pt1.X += (state.ArgStack.Vals[i])
		pt2 := pt1
		pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
		pt3 := pt2
		pt3.Y += (state.ArgStack.Vals[i+3])
		i += 4

		for ; i+8 <= state.ArgStack.Top; i += 8 {
			out.curve(pt1, pt2, pt3)
			pt1 = out.CurrentPoint
			pt1.Y += (state.ArgStack.Vals[i])
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
			pt3 = pt2
			pt3.X += (state.ArgStack.Vals[i+3])
			out.curve(pt1, pt2, pt3)

			pt1 = pt3
			pt1.X += state.ArgStack.Vals[i+4]
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+5], state.ArgStack.Vals[i+6])
			pt3 = pt2
			pt3.Y += state.ArgStack.Vals[i+7]
		}
		if i < state.ArgStack.Top {
			pt3.X += (state.ArgStack.Vals[i])
		}
		out.curve(pt1, pt2, pt3)
	} else {
		for ; i+8 <= state.ArgStack.Top; i += 8 {
			pt1 := out.CurrentPoint
			pt1.X += (state.ArgStack.Vals[i])
			pt2 := pt1
			pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
			pt3 := pt2
			pt3.Y += (state.ArgStack.Vals[i+3])
			out.curve(pt1, pt2, pt3)

			pt1 = pt3
			pt1.Y += (state.ArgStack.Vals[i+4])
			pt2 = pt1
			pt2.Move(state.ArgStack.Vals[i+5], state.ArgStack.Vals[i+6])
			pt3 = pt2
			pt3.X += (state.ArgStack.Vals[i+7])
			if (state.ArgStack.Top-i < 16) && ((state.ArgStack.Top & 1) != 0) {
				pt3.Y += state.ArgStack.Vals[i+8]
			}
			out.curve(pt1, pt2, pt3)
		}
	}
}

func (out *CharstringReader) Rcurveline(state *Machine) error {
	argCount := state.ArgStack.Top
	if argCount < 8 {
		return fmt.Errorf("expected at least 8 operands for <rcurveline>, got %d", argCount)
	}

	var i int32
	curveLimit := argCount - 2
	for ; i+6 <= curveLimit; i += 6 {
		pt1 := out.CurrentPoint
		pt1.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
		pt2 := pt1
		pt2.Move(state.ArgStack.Vals[i+2], state.ArgStack.Vals[i+3])
		pt3 := pt2
		pt3.Move(state.ArgStack.Vals[i+4], state.ArgStack.Vals[i+5])
		out.curve(pt1, pt2, pt3)
	}

	pt1 := out.CurrentPoint
	pt1.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
	out.line(pt1)

	return nil
}

func (out *CharstringReader) Rlinecurve(state *Machine) error {
	argCount := state.ArgStack.Top
	if argCount < 8 {
		return fmt.Errorf("expected at least 8 operands for <rlinecurve>, got %d", argCount)
	}
	var i int32
	lineLimit := argCount - 6
	for ; i+2 <= lineLimit; i += 2 {
		pt1 := out.CurrentPoint
		pt1.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
		out.line(pt1)
	}

	pt1 := out.CurrentPoint
	pt1.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
	pt2 := pt1
	pt2.Move(state.ArgStack.Vals[i+2], state.ArgStack.Vals[i+3])
	pt3 := pt2
	pt3.Move(state.ArgStack.Vals[i+4], state.ArgStack.Vals[i+5])
	out.curve(pt1, pt2, pt3)

	return nil
}

func (out *CharstringReader) Vvcurveto(state *Machine) {
	var i int32
	pt1 := out.CurrentPoint
	if (state.ArgStack.Top & 1) != 0 {
		pt1.X += state.ArgStack.Vals[i]
		i++
	}
	for ; i+4 <= state.ArgStack.Top; i += 4 {
		pt1.Y += state.ArgStack.Vals[i]
		pt2 := pt1
		pt2.Move(state.ArgStack.Vals[i+1], state.ArgStack.Vals[i+2])
		pt3 := pt2
		pt3.Y += state.ArgStack.Vals[i+3]
		out.curve(pt1, pt2, pt3)
		pt1 = out.CurrentPoint
	}
}

func (out *CharstringReader) Hflex(state *Machine) error {
	if state.ArgStack.Top != 7 {
		return fmt.Errorf("expected 7 operands for <hflex>, got %d", state.ArgStack.Top)
	}

	pt1 := out.CurrentPoint
	pt1.X += state.ArgStack.Vals[0]
	pt2 := pt1
	pt2.Move(state.ArgStack.Vals[1], state.ArgStack.Vals[2])
	pt3 := pt2
	pt3.X += state.ArgStack.Vals[3]
	pt4 := pt3
	pt4.X += state.ArgStack.Vals[4]
	pt5 := pt4
	pt5.X += state.ArgStack.Vals[5]
	pt5.Y = pt1.Y
	pt6 := pt5
	pt6.X += state.ArgStack.Vals[6]

	out.doubleCurve(pt1, pt2, pt3, pt4, pt5, pt6)
	return nil
}

func (out *CharstringReader) Flex(state *Machine) error {
	if state.ArgStack.Top != 13 {
		return fmt.Errorf("expected 13 operands for <flex>, got %d", state.ArgStack.Top)
	}

	pt1 := out.CurrentPoint
	pt1.Move(state.ArgStack.Vals[0], state.ArgStack.Vals[1])
	pt2 := pt1
	pt2.Move(state.ArgStack.Vals[2], state.ArgStack.Vals[3])
	pt3 := pt2
	pt3.Move(state.ArgStack.Vals[4], state.ArgStack.Vals[5])
	pt4 := pt3
	pt4.Move(state.ArgStack.Vals[6], state.ArgStack.Vals[7])
	pt5 := pt4
	pt5.Move(state.ArgStack.Vals[8], state.ArgStack.Vals[9])
	pt6 := pt5
	pt6.Move(state.ArgStack.Vals[10], state.ArgStack.Vals[11])

	out.doubleCurve(pt1, pt2, pt3, pt4, pt5, pt6)
	return nil
}

func (out *CharstringReader) Hflex1(state *Machine) error {
	if state.ArgStack.Top != 9 {
		return fmt.Errorf("expected 9 operands for <hflex1>, got %d", state.ArgStack.Top)
	}
	pt1 := out.CurrentPoint
	pt1.Move(state.ArgStack.Vals[0], state.ArgStack.Vals[1])
	pt2 := pt1
	pt2.Move(state.ArgStack.Vals[2], state.ArgStack.Vals[3])
	pt3 := pt2
	pt3.X += state.ArgStack.Vals[4]
	pt4 := pt3
	pt4.X += state.ArgStack.Vals[5]
	pt5 := pt4
	pt5.Move(state.ArgStack.Vals[6], state.ArgStack.Vals[7])
	pt6 := pt5
	pt6.X += state.ArgStack.Vals[8]
	pt6.Y = out.CurrentPoint.Y

	out.doubleCurve(pt1, pt2, pt3, pt4, pt5, pt6)
	return nil
}

func (out *CharstringReader) Flex1(state *Machine) error {
	if state.ArgStack.Top != 11 {
		return fmt.Errorf("expected 11 operands for <flex1>, got %d", state.ArgStack.Top)
	}

	var d Point
	for i := 0; i < 10; i += 2 {
		d.Move(state.ArgStack.Vals[i], state.ArgStack.Vals[i+1])
	}

	pt1 := out.CurrentPoint
	pt1.Move(state.ArgStack.Vals[0], state.ArgStack.Vals[1])
	pt2 := pt1
	pt2.Move(state.ArgStack.Vals[2], state.ArgStack.Vals[3])
	pt3 := pt2
	pt3.Move(state.ArgStack.Vals[4], state.ArgStack.Vals[5])
	pt4 := pt3
	pt4.Move(state.ArgStack.Vals[6], state.ArgStack.Vals[7])
	pt5 := pt4
	pt5.Move(state.ArgStack.Vals[8], state.ArgStack.Vals[9])
	pt6 := pt5

	if abs(d.X) > abs(d.Y) {
		pt6.X += state.ArgStack.Vals[10]
		pt6.Y = out.CurrentPoint.Y
	} else {
		pt6.X = out.CurrentPoint.X
		pt6.Y += state.ArgStack.Vals[10]
	}

	out.doubleCurve(pt1, pt2, pt3, pt4, pt5, pt6)
	return nil
}
