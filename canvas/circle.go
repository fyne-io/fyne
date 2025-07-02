package canvas

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Circle)(nil)

// Circle describes a colored circle primitive in a Fyne canvas
// Shadow support since 2.7
type Circle struct {
	baseShadow

	Position1 fyne.Position // The current top-left position of the Circle
	Position2 fyne.Position // The current bottomright position of the Circle

	FillColor   color.Color // The circle fill color
	StrokeColor color.Color // The circle stroke color
	StrokeWidth float32     // The stroke width of the circle

	// When true, the circle keeps the requested size for its content, and the total object size expands to include the shadow.
	// This means the object will be properly scaled and moved to fit the requested size and position.
	// The total size including the shadow will be larger than the requested size.
	//
	// Since: 2.7
	ExpandForShadow bool
}

// NewCircle returns a new Circle instance
func NewCircle(color color.Color) *Circle {
	return &Circle{FillColor: color}
}

// Hide will set this circle to not be visible
func (c *Circle) Hide() {
	c.Hidden = true

	repaint(c)
}

// MinSize for a Circle simply returns Size{1, 1} as there is no
// explicit content
func (c *Circle) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// Move the circle object to a new position, relative to its parent / canvas
// If ExpandForShadow is true, the position is adjusted to account for the shadow paddings.
// The circle position is then updated accordingly to reflect the new position.
func (c *Circle) Move(pos fyne.Position) {
	if c.Position1 == pos {
		return
	}

	size := c.Size()
	if c.ExpandForShadow {
		_, p := c.SizeAndPositionWithShadow(size)
		pos = pos.Add(p)
	}
	c.Position1 = pos
	c.Position2 = c.Position1.Add(size)

	repaint(c)
}

// Position gets the current top-left position of this circle object, relative to its parent / canvas
func (c *Circle) Position() fyne.Position {
	return c.Position1
}

// Refresh causes this object to be redrawn with its configured state.
func (c *Circle) Refresh() {
	Refresh(c)
}

// Resize sets a new bottom-right position for the circle object
// If it has a stroke width this will cause it to Refresh.
// If ExpandForShadow is true, the size is adjusted to accommodate the shadow,
// ensuring the circle size matches the requested size.
// The total size including the shadow will be larger than the requested size.
func (c *Circle) Resize(size fyne.Size) {
	if c.ExpandForShadow {
		if size == c.ContentSize() {
			return
		}
		sizeWithShadow, _ := c.SizeAndPositionWithShadow(size)
		size = sizeWithShadow
	} else if size == c.Size() {
		return
	}

	c.Position2 = c.Position1.Add(size)

	Refresh(c)
}

// Show will set this circle to be visible
func (c *Circle) Show() {
	c.Hidden = false

	c.Refresh()
}

// Size returns the current size of bounding box for this circle object
func (c *Circle) Size() fyne.Size {
	return fyne.NewSize(
		float32(math.Abs(float64(c.Position2.X)-float64(c.Position1.X))),
		float32(math.Abs(float64(c.Position2.Y)-float64(c.Position1.Y))),
	)
}

// Visible returns true if this circle is visible, false otherwise
func (c *Circle) Visible() bool {
	return !c.Hidden
}
