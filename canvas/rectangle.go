package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Rectangle)(nil)

// Rectangle describes a colored rectangle primitive in a Fyne canvas
// Shadow support since 2.7
type Rectangle struct {
	baseShadow

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

	// The `CornerRadius` field sets the default radius for all corners. If a specific corner
	// radius (such as `TopRightCornerRadius`, `TopLeftCornerRadius`, `BottomRightCornerRadius`,
	// or `BottomLeftCornerRadius`) is set to a value greater than `CornerRadius`, that value
	// will override `CornerRadius` for the respective corner. Otherwise, `CornerRadius` is used.
	//
	// Since: 2.7
	TopRightCornerRadius    float32
	TopLeftCornerRadius     float32
	BottomRightCornerRadius float32
	BottomLeftCornerRadius  float32

	// When true, the rectangle keeps the requested size for its content, and the total object size expands to include the shadow.
	// This means the rectangle will be properly scaled and moved to fit the requested size and position.
	// The total size including the shadow will be larger than the requested size.
	//
	// Since: 2.7
	ExpandForShadow bool
}

// Hide will set this rectangle to not be visible
func (r *Rectangle) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the rectangle to a new position, relative to its parent / canvas
// If ExpandForShadow is true, the position is adjusted to account for the shadow paddings.
// The rectangle position is then updated accordingly to reflect the new position.
func (r *Rectangle) Move(pos fyne.Position) {
	if r.Position() == pos {
		return
	}

	if r.ExpandForShadow {
		_, p := r.SizeAndPositionWithShadow(r.Size())
		pos = pos.Add(p)
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
// If ExpandForShadow is true, the size is adjusted to accommodate the shadow,
// ensuring the rectangle size matches the requested size.
// The total size including the shadow will be larger than the requested size.
func (r *Rectangle) Resize(s fyne.Size) {
	if r.ExpandForShadow {
		if s == r.ContentSize() {
			return
		}
		size, _ := r.SizeAndPositionWithShadow(s)
		s = size
	} else if s == r.Size() {
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
