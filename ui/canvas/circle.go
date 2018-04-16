package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// Circle describes a coloured circle primitive in a Fyne canvas
type Circle struct {
	Position1 ui.Position // The current top-left position of the Circle
	Position2 ui.Position // The current bottomright position of the Circle

	FillColor   color.RGBA // The circle fill colour
	StrokeColor color.RGBA // The circle stroke colour
	StrokeWidth float32    // The stroke width of the circle
}

// CurrentSize returns the current size of bounding box for this circle object
func (l *Circle) CurrentSize() ui.Size {
	return ui.NewSize(l.Position2.X-l.Position1.X, l.Position2.Y-l.Position1.Y)
}

// Resize sets a new bottom-right position for the circle object
func (l *Circle) Resize(size ui.Size) {
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// CurrentPosition gets the current top-left position of this circle object, relative to it's parent / canvas
func (l *Circle) CurrentPosition() ui.Position {
	return l.Position1
}

// Move the circle object to a new position, relative to it's parent / canvas
func (l *Circle) Move(pos ui.Position) {
	size := l.CurrentSize()
	l.Position1 = pos
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a Circle simply returns Size{1, 1} as there is no
// explicit content
func (l *Circle) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewCircle returns a new Circle instance
func NewCircle(color color.RGBA) *Circle {
	return &Circle{
		FillColor: color,
	}
}
