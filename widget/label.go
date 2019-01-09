package widget

import (
	"fyne.io/fyne"
)

// Label widget is a label component with appropriate padding and layout.
type Label struct {
	textWidget
	Text      string
	Alignment fyne.TextAlign // The alignment of the Text
	TextStyle fyne.TextStyle // The style of the label text
}

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewLabelWithStyle creates a new layout widget with the set text content
func NewLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle) *Label {
	l := &Label{
		Text:      text,
		Alignment: alignment,
		TextStyle: style,
	}

	Renderer(l).Refresh()
	return l
}

// SetText sets the text of the label
func (l *Label) SetText(text string) {
	l.Text = text
	l.textWidget.SetText(text)
	Renderer(l).Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	l.textWidget = textWidget{
		Alignment: l.Alignment,
		TextStyle: l.TextStyle,
	}
	l.textWidget.SetText(l.Text)
	r := l.textWidget.CreateRenderer()
	r.Refresh()
	return r
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Label) Resize(size fyne.Size) {
	l.resize(size, l)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Label) Move(pos fyne.Position) {
	l.move(pos, l)
}

// MinSize returns the smallest size this widget can shrink to
func (l *Label) MinSize() fyne.Size {
	return l.minSize(l)
}

// Show this widget, if it was previously hidden
func (l *Label) Show() {
	l.show(l)
}

// Hide this widget, if it was previously visible
func (l *Label) Hide() {
	l.hide(l)
}
