package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Rectangle)(nil)

// Rectangle describes a colored rectangle primitive in a Fyne canvas
type Rectangle struct {
	baseObject

	FillColor   color.Color // The rectangle fill color
	StrokeColor color.Color // The rectangle stroke color
	StrokeWidth float32     // The stroke width of the rectangle
	// The radius of the rectangle corners
	//
	// Since: 2.4
	CornerRadius float32

	// Enforce an aspect ratio for the rectangle, the content will be made shorter or narrower
	// to meet the requested aspect, if set.
	//
	// Since: 2.7
	Aspect float32

	// The radius of the rectangle top-right corner only.
	//
	// Since: 2.7
	TopRightCornerRadius float32

	// The radius of the rectangle top-left corner only.
	//
	// Since: 2.7
	TopLeftCornerRadius float32

	// The radius of the rectangle bottom-right corner only.
	//
	// Since: 2.7
	BottomRightCornerRadius float32

	// The radius of the rectangle bottom-left corner only.
	//
	// Since: 2.7
	BottomLeftCornerRadius float32
}

// Hide will set this rectangle to not be visible
func (r *Rectangle) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the rectangle to a new position, relative to its parent / canvas
func (r *Rectangle) Move(pos fyne.Position) {
	if r.Position() == pos {
		return
	}

	r.baseObject.Move(pos)

	repaint(r)
}

// Refresh causes this rectangle to be redrawn with its configured state.
func (r *Rectangle) Refresh() {
	Refresh(r)
}

// Resize on a rectangle updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
// If Aspect is non-zero it may cause the rectangle to be smaller than the requested size.
func (r *Rectangle) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)
	if r.StrokeWidth == 0 {
		return
	}

	Refresh(r)
}

// NewRectangle returns a new Rectangle instance
func NewRectangle(color color.Color) *Rectangle {
	return &Rectangle{
		FillColor: color,
	}
}
