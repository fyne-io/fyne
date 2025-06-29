package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Square)(nil)

// Square describes a colored rectangle primitive in a Fyne canvas that always has the same
// width and height.
// The rendered portion will be in the center of the available space, as large as it can be.
//
// Since: 2.7
type Square struct {
	baseObject

	FillColor    color.Color // The square fill color
	StrokeColor  color.Color // The square stroke color
	StrokeWidth  float32     // The stroke width of the square
	CornerRadius float32     // How large the corner radius should be
}

// Hide will set this square to not be visible
func (s *Square) Hide() {
	s.baseObject.Hide()

	repaint(s)
}

// Move the square to a new position, relative to its parent / canvas
func (s *Square) Move(pos fyne.Position) {
	if s.Position() == pos {
		return
	}
	s.baseObject.Move(pos)

	repaint(s)
}

// Refresh causes this square to be redrawn with its configured state.
func (s *Square) Refresh() {
	Refresh(s)
}

// Resize on a square updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
// The square may have transparent padding on the top/bottom or left/right
// so that it is still a square shape in the specified size.
func (s *Square) Resize(size fyne.Size) {
	if size == s.Size() {
		return
	}

	s.baseObject.Resize(size)
	if s.StrokeWidth == 0 {
		return
	}

	Refresh(s)
}

// NewSquare returns a new Square instance.
//
// Since: 2.7
func NewSquare(color color.Color) *Square {
	return &Square{
		FillColor: color,
	}
}
