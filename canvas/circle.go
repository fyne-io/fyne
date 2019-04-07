package canvas

import (
	"image/color"

	"fyne.io/fyne"
)

// Circle describes a coloured circle primitive in a Fyne canvas
type Circle struct {
	Position1 fyne.Position // The current top-left position of the Circle
	Position2 fyne.Position // The current bottomright position of the Circle
	Hidden    bool          // Is this circle currently hidden

	FillColor   color.Color // The circle fill colour
	StrokeColor color.Color // The circle stroke colour
	StrokeWidth float32     // The stroke width of the circle
}

// Size returns the current size of bounding box for this circle object
func (l *Circle) Size() fyne.Size {
	return fyne.NewSize(l.Position2.X-l.Position1.X, l.Position2.Y-l.Position1.Y)
}

// Resize sets a new bottom-right position for the circle object
func (l *Circle) Resize(size fyne.Size) {
	l.Position2 = fyne.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// Position gets the current top-left position of this circle object, relative to it's parent / canvas
func (l *Circle) Position() fyne.Position {
	return l.Position1
}

// Move the circle object to a new position, relative to it's parent / canvas
func (l *Circle) Move(pos fyne.Position) {
	size := l.Size()
	l.Position1 = pos
	l.Position2 = fyne.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a Circle simply returns Size{1, 1} as there is no
// explicit content
func (l *Circle) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// Visible returns true if this circle is visible, false otherwise
func (l *Circle) Visible() bool {
	return !l.Hidden
}

// Show will set this circle to be visible
func (l *Circle) Show() {
	l.Hidden = false

	Refresh(l)
}

// Hide will set this circle to not be visible
func (l *Circle) Hide() {
	l.Hidden = true

	Refresh(l)
}

// NewCircle returns a new Circle instance
func NewCircle(color color.Color) *Circle {
	return &Circle{
		FillColor: color,
	}
}
