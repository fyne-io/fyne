package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

// LineObject describes a coloured line primitive in a Fyne canvas
type LineObject struct {
	Position1 ui.Position // The current top-left position of the LineObject
	Position2 ui.Position // The current bottomright position of the LineObject

	Color color.RGBA // The line stroke colour
	Width float32    // The stroke width of the line
}

// CurrentSize returns the current size of bounding box for this line object
func (l *LineObject) CurrentSize() ui.Size {
	return ui.NewSize(l.Position2.X-l.Position1.X, l.Position2.Y-l.Position1.Y)
}

// Resize sets a new bottom-right position for the line object
func (l *LineObject) Resize(size ui.Size) {
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// CurrentPosition gets the current top-left position of this line object, relative to it's parent / canvas
func (l *LineObject) CurrentPosition() ui.Position {
	return l.Position1
}

// Move the line object to a new position, relative to it's parent / canvas
func (l *LineObject) Move(pos ui.Position) {
	size := l.CurrentSize()
	l.Position1 = pos
	l.Position2 = ui.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a LineObject simply returns Size{1, 1} as there is no
// explicit content
func (l *LineObject) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

// NewLine returns a new LineObject instance
func NewLine(color color.RGBA) *LineObject {
	return &LineObject{
		Color: color,
		Width: 1,
	}
}
