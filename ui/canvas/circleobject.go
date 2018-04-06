package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// CircleObject describes a coloured circle primitive in a Fyne canvas
type CircleObject struct {
	Position1 ui.Position // The current top-left position of the CircleObject
	Position2 ui.Position // The current bottomright position of the CircleObject

	FillColor   color.RGBA // The circle fill colour
	StrokeColor color.RGBA // The circle stroke colour
	StrokeWidth float32    // The stroke width of the circle
}

// CurrentSize returns the current size of bounding box for this circle object
func (l *CircleObject) CurrentSize() ui.Size {
	return ui.NewSize(l.Position2.X-l.Position1.X, l.Position2.Y-l.Position1.Y)
}

// Resize sets a new bottom-right position for the circle object
func (l *CircleObject) Resize(size ui.Size) {
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// CurrentPosition gets the current top-left position of this circle object, relative to it's parent / canvas
func (l *CircleObject) CurrentPosition() ui.Position {
	return l.Position1
}

// Move the circle object to a new position, relative to it's parent / canvas
func (l *CircleObject) Move(pos ui.Position) {
	size := l.CurrentSize()
	l.Position1 = pos
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a CircleObject simply returns Size{1, 1} as there is no
// explicit content
func (l *CircleObject) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewLine returns a new CircleObject instance
func NewCircle(color color.RGBA) *CircleObject {
	return &CircleObject{
		FillColor: color,
	}
}
