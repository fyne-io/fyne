package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Ellipse)(nil)

// Ellipse describes a colored ellipse primitive in a Fyne canvas
// with radii being half the width and height of its size.
//
// Since: 2.8
type Ellipse struct {
	baseObject

	FillColor   color.Color // The ellipse fill color
	StrokeColor color.Color // The ellipse stroke color
	StrokeWidth float32     // The stroke width of the ellipse
	Angle       float32     // Angle of ellipse, in degrees (positive means clockwise, negative means counter-clockwise direction).
}

// Hide will set this ellipse to not be visible
func (e *Ellipse) Hide() {
	e.baseObject.Hide()

	repaint(e)
}

// Move the ellipse to a new position, relative to its parent / canvas
func (e *Ellipse) Move(pos fyne.Position) {
	if e.Position() == pos {
		return
	}

	e.baseObject.Move(pos)

	repaint(e)
}

// Refresh causes this ellipse to be redrawn with its configured state.
func (e *Ellipse) Refresh() {
	Refresh(e)
}

// Resize on a ellipse updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
// If Angle is non-zero, the ellipse may be rotated so that its boundaries
// extend beyond the shorter dimension of the requested size.
func (e *Ellipse) Resize(s fyne.Size) {
	if s == e.Size() {
		return
	}

	e.baseObject.Resize(s)
	if e.StrokeWidth == 0 {
		return
	}

	Refresh(e)
}

// NewEllipse returns a new Ellipse instance
func NewEllipse(color color.Color) *Ellipse {
	return &Ellipse{FillColor: color}
}
