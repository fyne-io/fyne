package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// RectangleObject describes a coloured rectangle primitive in a Fyne canvas
type RectangleObject struct {
	Size     ui.Size     // The current size of the RectangleObject - the font will not scale to this Size
	Position ui.Position // The current position of the RectangleObject

	Color color.RGBA // The rectangle fill colour
}

// CurrentSize returns the current size of this rectangle object
func (r *RectangleObject) CurrentSize() ui.Size {
	return r.Size
}

// Resize sets a new size for the rectangle object
func (r *RectangleObject) Resize(size ui.Size) {
	r.Size = size
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *RectangleObject) CurrentPosition() ui.Position {
	return r.Position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *RectangleObject) Move(pos ui.Position) {
	r.Position = pos
}

// MinSize for a RectangleObject simply returns Size{1, 1} as there is no
// explicit content
func (r *RectangleObject) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewRectangle returns a new RectangleObject instance
func NewRectangle(color color.RGBA) *RectangleObject {
	return &RectangleObject{
		Color: color,
	}
}
