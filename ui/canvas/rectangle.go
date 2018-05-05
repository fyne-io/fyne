package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// Rectangle describes a coloured rectangle primitive in a Fyne canvas
type Rectangle struct {
	Size     ui.Size     // The current size of the Rectangle
	Position ui.Position // The current position of the Rectangle

	FillColor   color.RGBA // The rectangle fill colour
	StrokeColor color.RGBA // The rectangle stroke colour
	StrokeWidth float32    // The stroke width of the rectangle
}

// CurrentSize returns the current size of this rectangle object
func (r *Rectangle) CurrentSize() ui.Size {
	return r.Size
}

// Resize sets a new size for the rectangle object
func (r *Rectangle) Resize(size ui.Size) {
	r.Size = size
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *Rectangle) CurrentPosition() ui.Position {
	return r.Position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *Rectangle) Move(pos ui.Position) {
	r.Position = pos
}

// MinSize for a Rectangle simply returns Size{1, 1} as there is no
// explicit content
func (r *Rectangle) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewRectangle returns a new Rectangle instance
func NewRectangle(color color.RGBA) *Rectangle {
	return &Rectangle{
		FillColor: color,
	}
}
