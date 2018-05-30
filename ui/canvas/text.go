package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

// Text describes a text primitive in a Fyne canvas
type Text struct {
	Size      ui.Size      // The current size of the Text - the font will not scale to this Size
	Position  ui.Position  // The current position of the Text
	Color     color.RGBA   // The main text draw colour
	Alignment ui.TextAlign // The alignment of the text content within the current size

	Text     string // The string content of this Text
	FontSize int    // Size of the font - if the Canvas scale is 1.0 this will be equivalent to point size
	Bold     bool   // Should the text be bold
	Italic   bool   // Should the text be italic

	minSize ui.Size
}

// CurrentSize gets the current size of this text object
func (t *Text) CurrentSize() ui.Size {
	return t.Size
}

// Resize sets a new size for the text object
func (t *Text) Resize(size ui.Size) {
	t.Size = size
}

// CurrentPosition gets the current position of this text object, relative to it's parent / canvas
func (t *Text) CurrentPosition() ui.Position {
	return t.Position
}

// Move the text object to a new position, relative to it's parent / canvas
func (t *Text) Move(pos ui.Position) {
	t.Position = pos
}

// SetMinSize is used by the render implementation to specify the space the text requires.
// This should not be called by application code.
func (t *Text) SetMinSize(size ui.Size) {
	t.minSize = size
}

// MinSize returns the minimum size of this text objet based on it's font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() ui.Size {
	return t.minSize
}

// NewText returns a new Text implementation
func NewText(text string, color color.RGBA) *Text {
	return &Text{
		Color:    color,
		Text:     text,
		FontSize: theme.TextSize(),
	}
}
