package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
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

	// FontSource defines a resource that can be used instead of the theme for looking up the font.
	// When a font source is set the `TextStyle` may not be effective, as it will be limited to the styles
	// present in the data provided.
	//
	// Since: 2.5
	FontSource fyne.Resource

	// measured is what MinSize last measured, kept while the fields it was
	// measured from and the fonts stay the same.
	measured textMeasure
}

// textMeasure is a measured size with what it was measured for.
type textMeasure struct {
	valid  bool
	gen    uint64 // of the font metrics, see cache.FontMetricsGeneration
	text   string
	size   float32
	style  fyne.TextStyle
	result fyne.Size
}

// Hide will set this text to not be visible
func (t *Text) Hide() {
	t.baseObject.Hide()

	repaint(t)
}

// MinSize returns the minimum size of this text object based on its font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() fyne.Size {
	// Every repaint asks, from the layouts above and from the painter, and
	// measuring is a cache lookup that hashes the string, while a Text mostly
	// keeps its string. A font source is left to the cache: not every
	// resource is comparable.
	m := &t.measured
	gen := cache.FontMetricsGeneration()
	if t.FontSource == nil && m.valid && m.gen == gen && m.text == t.Text && m.size == t.TextSize && m.style == t.TextStyle {
		return m.result
	}
	s, _ := fyne.CurrentApp().Driver().RenderedTextSize(t.Text, t.TextSize, t.TextStyle, t.FontSource)
	if t.FontSource == nil {
		*m = textMeasure{valid: true, gen: gen, text: t.Text, size: t.TextSize, style: t.TextStyle, result: s}
	}
	return s
}

// Move the text to a new position, relative to its parent / canvas
func (t *Text) Move(pos fyne.Position) {
	if t.Position() == pos {
		return
	}

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
func (*Text) SetMinSize(fyne.Size) {
	// no-op
}

// Refresh causes this text to be redrawn with its configured state.
func (t *Text) Refresh() {
	Refresh(t)
}

// NewText returns a new Text implementation
func NewText(text string, c color.Color) *Text {
	return &Text{
		Color:    c,
		Text:     text,
		TextSize: theme.TextSize(),
	}
}
