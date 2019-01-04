package canvas

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Text describes a text primitive in a Fyne canvas.
// A text object can have a style set which will apply to the whole string.
// No formatting or text parsing will be performed
type Text struct {
	baseObject
	Alignment fyne.TextAlign // The alignment of the text content

	Color     color.Color    // The main text draw colour
	Text      string         // The string content of this Text
	TextSize  int            // Size of the text - if the Canvas scale is 1.0 this will be equivalent to point size
	TextStyle fyne.TextStyle // The style of the text content
}

// MinSize returns the minimum size of this text objet based on it's font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() fyne.Size {
	return fyne.CurrentApp().Driver().RenderedTextSize(t.Text, t.TextSize, t.TextStyle)
}

// NewText returns a new Text implementation
func NewText(text string, color color.Color) *Text {
	return &Text{
		Color:    color,
		Text:     text,
		TextSize: theme.TextSize(),
	}
}
