package canvas

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Line)(nil)

// Line describes a colored line primitive in a Fyne canvas.
// Lines are special as they can have a negative width or height to indicate
// an inverse slope (i.e. slope up vs down).
type Line struct {
	Position1 fyne.Position // The current top-left position of the Line
	Position2 fyne.Position // The current bottom-right position of the Line
	Hidden    bool          // Is this Line currently hidden

	StrokeColor color.Color // The line stroke color
	StrokeWidth float32     // The stroke width of the line
}

// Size returns the current size of bounding box for this line object
func (l *Line) Size() fyne.Size {
	return fyne.NewSize(float32(math.Abs(float64(l.Position2.X)-float64(l.Position1.X))),
		float32(math.Abs(float64(l.Position2.Y)-float64(l.Position1.Y))))
}

// Resize sets a new bottom-right position for the line object, then it will then be refreshed.
func (l *Line) Resize(size fyne.Size) {
	if size == l.Size() {
		return
	}

	if l.Position1.X <= l.Position2.X {
		l.Position2.X = l.Position1.X + size.Width
	} else {
		l.Position1.X = l.Position2.X + size.Width
	}
	if l.Position1.Y <= l.Position2.Y {
		l.Position2.Y = l.Position1.Y + size.Height
	} else {
		l.Position1.Y = l.Position2.Y + size.Height
	}
	Refresh(l)
}

// Position gets the current top-left position of this line object, relative to its parent / canvas
func (l *Line) Position() fyne.Position {
	return fyne.NewPos(fyne.Min(l.Position1.X, l.Position2.X), fyne.Min(l.Position1.Y, l.Position2.Y))
}

// Move the line object to a new position, relative to its parent / canvas
func (l *Line) Move(pos fyne.Position) {
	oldPos := l.Position()
	deltaX := pos.X - oldPos.X
	deltaY := pos.Y - oldPos.Y

	l.Position1 = l.Position1.Add(fyne.NewPos(deltaX, deltaY))
	l.Position2 = l.Position2.Add(fyne.NewPos(deltaX, deltaY))
	repaint(l)
}

// MinSize for a Line simply returns Size{1, 1} as there is no
// explicit content
func (l *Line) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

// Visible returns true if this line// Show will set this circle to be visible is visible, false otherwise
func (l *Line) Visible() bool {
	return !l.Hidden
}

// Show will set this line to be visible
func (l *Line) Show() {
	l.Hidden = false

	l.Refresh()
}

// Hide will set this line to not be visible
func (l *Line) Hide() {
	l.Hidden = true

	repaint(l)
}

// Refresh causes this line to be redrawn with its configured state.
func (l *Line) Refresh() {
	Refresh(l)
}

// NewLine returns a new Line instance
func NewLine(color color.Color) *Line {
	return &Line{
		StrokeColor: color,
		StrokeWidth: 1,
	}
}
