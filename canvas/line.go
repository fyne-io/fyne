package canvas

import "math"
import "image/color"

import "github.com/fyne-io/fyne"

// Line describes a coloured line primitive in a Fyne canvas.
// Lines are special as they can have a negative width or height to indicate
// an inverse slope (i.e. slope up vs down).
type Line struct {
	Position1 fyne.Position // The current top-left position of the Line
	Position2 fyne.Position // The current bottomright position of the Line
	Options   Options       // Options to pass to the renderer

	StrokeColor color.RGBA // The line stroke colour
	StrokeWidth float32    // The stroke width of the line
}

// CurrentSize returns the current size of bounding box for this line object
func (l *Line) CurrentSize() fyne.Size {
	return fyne.NewSize(int(math.Abs(float64(l.Position2.X))-math.Abs(float64(l.Position1.X))),
		int(math.Abs(float64(l.Position2.Y))-math.Abs(float64(l.Position1.Y))))
}

// Resize sets a new bottom-right position for the line object
func (l *Line) Resize(size fyne.Size) {
	l.Position2 = fyne.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// CurrentPosition gets the current top-left position of this line object, relative to it's parent / canvas
func (l *Line) CurrentPosition() fyne.Position {
	return fyne.NewPos(fyne.Min(l.Position1.X, l.Position2.X), fyne.Min(l.Position1.Y, l.Position2.Y))
}

// Move the line object to a new position, relative to it's parent / canvas
func (l *Line) Move(pos fyne.Position) {
	size := l.CurrentSize()
	l.Position1 = pos
	l.Position2 = fyne.NewPos(l.Position1.X+size.Width, l.Position1.Y+size.Height)
}

// MinSize for a Line simply returns Size{1, 1} as there is no
// explicit content
func (l *Line) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// NewLine returns a new Line instance
func NewLine(color color.RGBA) *Line {
	return &Line{
		StrokeColor: color,
		StrokeWidth: 1,
	}
}
