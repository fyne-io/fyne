package canvas

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Circle)(nil)

// Circle describes a colored circle primitive in a Fyne canvas
type Circle struct {
	Position1 fyne.Position // The current top-left position of the Circle
	Position2 fyne.Position // The current bottomright position of the Circle
	Hidden    bool          // Is this circle currently hidden

	FillColor   color.Color // The circle fill color
	StrokeColor color.Color // The circle stroke color
	StrokeWidth float32     // The stroke width of the circle
}

// NewCircle returns a new Circle instance
func NewCircle(color color.Color) *Circle {
	return &Circle{
		FillColor: color,
	}
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
func (c *Circle) Move(pos fyne.Position) {
	size := c.Size()
	c.Position1 = pos
	c.Position2 = fyne.NewPos(c.Position1.X+size.Width, c.Position1.Y+size.Height)
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
func (c *Circle) Resize(size fyne.Size) {
	if size == c.Size() {
		return
	}

	c.Position2 = fyne.NewPos(c.Position1.X+size.Width, c.Position1.Y+size.Height)

	Refresh(c)
}

// Show will set this circle to be visible
func (c *Circle) Show() {
	c.Hidden = false

	c.Refresh()
}

// Size returns the current size of bounding box for this circle object
func (c *Circle) Size() fyne.Size {
	return fyne.NewSize(float32(math.Abs(float64(c.Position2.X)-float64(c.Position1.X))),
		float32(math.Abs(float64(c.Position2.Y)-float64(c.Position1.Y))))
}

// Visible returns true if this circle is visible, false otherwise
func (c *Circle) Visible() bool {
	return !c.Hidden
}
