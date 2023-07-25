package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Text)(nil)

// Text describes a text primitive in a Fyne canvas.
// A text object can have a style set which will apply to the whole string.
// No formatting or text parsing will be performed
type Text struct {
	baseObject
	Alignment fyne.TextAlign // The alignment of the text content

	Color     color.Color    // The main text draw color
	Text      string         // The string content of this Text
	TextSize  float32        // Size of the text - if the Canvas scale is 1.0 this will be equivalent to point size
	TextStyle fyne.TextStyle // The style of the text content
}

// Hide will set this text to not be visible
func (t *Text) Hide() {
	t.baseObject.Hide()

	repaint(t)
}

// MinSize returns the minimum size of this text object based on its font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() fyne.Size {
	return fyne.MeasureText(t.Text, t.TextSize, t.TextStyle)
}

// Move the text to a new position, relative to its parent / canvas
func (t *Text) Move(pos fyne.Position) {
	t.baseObject.Move(pos)

	repaint(t)
}

// Resize on a text updates the new size of this object, which may not result in a visual change, depending on alignment.
func (t *Text) Resize(s fyne.Size) {
	if s == t.Size() {
		return
	}

	t.baseObject.Resize(s)
	Refresh(t)
}

// SetMinSize has no effect as the smallest size this canvas object can be is based on its font size and content.
func (t *Text) SetMinSize(fyne.Size) {
	// no-op
}

// Refresh causes this text to be redrawn with its configured state.
func (t *Text) Refresh() {
	Refresh(t)
}

// NewText returns a new Text implementation
func NewText(text string, color color.Color) *Text {
	size := float32(0)
	if fyne.CurrentApp() != nil { // nil app possible if app not started
		size = fyne.CurrentApp().Settings().Theme().Size("text") // manually name the size to avoid import loop
	}
	return &Text{
		Color:    color,
		Text:     text,
		TextSize: size,
	}
}
