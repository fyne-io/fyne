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
	baseShadow

	FillColor    color.Color // The square fill color
	StrokeColor  color.Color // The square stroke color
	StrokeWidth  float32     // The stroke width of the square
	CornerRadius float32     // How large the corner radius should be

	// The `CornerRadius` field sets the default radius for all corners. If a specific corner
	// radius (such as `TopRightCornerRadius`, `TopLeftCornerRadius`, `BottomRightCornerRadius`,
	// or `BottomLeftCornerRadius`) is set to a value greater than `CornerRadius`, that value
	// will override `CornerRadius` for the respective corner. Otherwise, `CornerRadius` is used.
	TopRightCornerRadius    float32
	TopLeftCornerRadius     float32
	BottomRightCornerRadius float32
	BottomLeftCornerRadius  float32

	// When true, the square keeps the requested size for its content, and the total object size expands to include the shadow.
	// This means the square will be properly scaled and moved to fit the requested size and position.
	// The total size including the shadow will be larger than the requested size.
	//
	// Since: 2.7
	ExpandForShadow bool
}

// Hide will set this square to not be visible
func (s *Square) Hide() {
	s.baseObject.Hide()

	repaint(s)
}

// Move the square to a new position, relative to its parent / canvas
// If ExpandForShadow is true, the position is adjusted to account for the shadow paddings.
// The square position is then updated accordingly to reflect the new position.
func (s *Square) Move(pos fyne.Position) {
	if s.ExpandForShadow {
		_, p := s.SizeAndPositionWithShadow(s.Size())
		pos = pos.Add(p)
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
// If ExpandForShadow is true, the size is adjusted to accommodate the shadow,
// ensuring the square size matches the requested size. 
// The total size including the shadow will be larger than the requested size.
func (s *Square) Resize(size fyne.Size) {
	if s.ExpandForShadow {
		if size == s.ContentSize() {
			return
		}
		sizeWithShadow, _ := s.SizeAndPositionWithShadow(size)
		size = sizeWithShadow
	} else if size == s.Size() {
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
