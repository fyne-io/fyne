package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Arc)(nil)

// Arc represents a filled arc or annular sector primitive that can be drawn on a Fyne canvas.
// It allows for the creation of circular, ring-shaped or pie-shaped segment, with configurable cutout ratio
// as well as customizable start and end angles to define the arc's length as the absolute difference between the two angles. The arc is always centered on its position.
// The arc is drawn from StartAngle to EndAngle (in degrees, positive is clockwise, negative is counter-clockwise).
// 0°/360 is top, 90° is right, 180° is bottom, 270° is left
// 0°/-360 is top, -90° is left, -180° is bottom, -270° is right
//
// Since: 2.7
type Arc struct {
	baseObject

	FillColor    color.Color // The arc fill colour
	StartAngle   float32     // Start angle in degrees
	EndAngle     float32     // End angle in degrees
	CornerRadius float32     // Radius used to round the corners
	StrokeColor  color.Color // The arc stroke color
	StrokeWidth  float32     // The stroke width of the arc
	CutoutRatio  float32     // Controls what portion of the inner should be cut out. A value of 0.0 results in a pie slice, while 1.0 results in a stroke.
}

// Hide will set this arc to not be visible.
func (a *Arc) Hide() {
	a.baseObject.Hide()

	repaint(a)
}

// Move the arc to a new position, relative to its parent / canvas.
// The position specifies the **center** of the arc.
func (a *Arc) Move(pos fyne.Position) {
	if a.Position() == pos {
		return
	}
	a.baseObject.Move(pos)

	repaint(a)
}

// Refresh causes this arc to be redrawn with its configured state.
func (a *Arc) Refresh() {
	Refresh(a)
}

// Resize updates the logical size of the arc.
// The arc is always drawn centered on its Position().
func (a *Arc) Resize(s fyne.Size) {
	if s == a.Size() {
		return
	}

	a.baseObject.Resize(s)

	repaint(a)
}

// NewArc returns a new Arc instance with the specified start and end angles (in degrees), fill color and cutout ratio.
func NewArc(startAngle, endAngle, cutoutRatio float32, color color.Color) *Arc {
	return &Arc{
		StartAngle:  startAngle,
		EndAngle:    endAngle,
		FillColor:   color,
		CutoutRatio: cutoutRatio,
	}
}

// NewPieArc returns a new pie-shaped Arc instance with the specified start and end angles (in degrees), fill color and cutout ratio set to 0.
func NewPieArc(startAngle, endAngle float32, color color.Color) *Arc {
	return &Arc{
		StartAngle:  startAngle,
		EndAngle:    endAngle,
		FillColor:   color,
		CutoutRatio: 0.0,
	}
}

// NewDoughnutArc returns a new doughnut-shaped Arc instance with the specified start and end angles (in degrees), fill color and cutout ratio set to 0.5.
func NewDoughnutArc(startAngle, endAngle float32, color color.Color) *Arc {
	return &Arc{
		StartAngle:  startAngle,
		EndAngle:    endAngle,
		FillColor:   color,
		CutoutRatio: 0.5,
	}
}
