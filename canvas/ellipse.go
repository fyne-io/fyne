package canvas

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Ellipse)(nil)

// Ellipse describes a colored ellipse primitive in a Fyne canvas
//
// Since: 2.8
type Ellipse struct {
	Position1 fyne.Position // The current top-left position of the Ellipse
	Position2 fyne.Position // The current bottom-right position of the Ellipse
	Hidden    bool          // Is this ellipse currently hidden

	FillColor   color.Color // The ellipse fill color
	StrokeColor color.Color // The ellipse stroke color
	StrokeWidth float32     // The stroke width of the ellipse
	Angle       float32     // Angle of ellipse, in degrees (positive means clockwise, negative means counter-clockwise direction).
}

// NewEllipse returns a new Ellipse instance
func NewEllipse(color color.Color) *Ellipse {
	return &Ellipse{FillColor: color}
}

// Hide will set this ellipse to not be visible
func (c *Ellipse) Hide() {
	c.Hidden = true

	repaint(c)
}

// MinSize for a Ellipse simply returns Size{1, 1} as there is no
// explicit content
func (c *Ellipse) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// Move the ellipse object to a new position, relative to its parent / canvas
func (c *Ellipse) Move(pos fyne.Position) {
	if c.Position1 == pos {
		return
	}

	size := c.Size()
	c.Position1 = pos
	c.Position2 = c.Position1.Add(size)

	repaint(c)
}

// Position gets the current top-left position of this ellipse object, relative to its parent / canvas
func (c *Ellipse) Position() fyne.Position {
	return c.Position1
}

// Refresh causes this object to be redrawn with its configured state.
func (c *Ellipse) Refresh() {
	Refresh(c)
}

// Resize sets a new bottom-right position for the ellipse object
// If it has a stroke width this will cause it to Refresh.
func (c *Ellipse) Resize(size fyne.Size) {
	if size == c.Size() {
		return
	}

	c.Position2 = c.Position1.Add(size)

	Refresh(c)
}

// Show will set this ellipse to be visible
func (c *Ellipse) Show() {
	c.Hidden = false

	c.Refresh()
}

// Size returns the current size of bounding box for this ellipse object
func (c *Ellipse) Size() fyne.Size {
	return fyne.NewSize(
		float32(math.Abs(float64(c.Position2.X)-float64(c.Position1.X))),
		float32(math.Abs(float64(c.Position2.Y)-float64(c.Position1.Y))),
	)
}

// Visible returns true if this ellipse is visible, false otherwise
func (c *Ellipse) Visible() bool {
	return !c.Hidden
}
