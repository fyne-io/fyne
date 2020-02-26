package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Label widget is a label component with appropriate padding and layout.
type Label struct {
	textProvider
	Text      string
	Alignment fyne.TextAlign // The alignment of the Text
	TextStyle fyne.TextStyle // The style of the label text
}

// NewLabel creates a new label widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewLabelWithStyle creates a new label widget with the set text content
func NewLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle) *Label {
	l := &Label{
		Text:      text,
		Alignment: alignment,
		TextStyle: style,
	}

	return l
}

// Refresh checks if the text content should be updated then refreshes the graphical context
func (l *Label) Refresh() {
	if l.Text != string(l.buffer) {
		l.textProvider.SetText(l.Text)
	}

	l.BaseWidget.Refresh()
}

// SetText sets the text of the label
func (l *Label) SetText(text string) {
	l.Text = text
	l.textProvider.SetText(text) // calls refresh
}

// textAlign tells the rendering textProvider our alignment
func (l *Label) textAlign() fyne.TextAlign {
	return l.Alignment
}

// textStyle tells the rendering textProvider our style
func (l *Label) textStyle() fyne.TextStyle {
	return l.TextStyle
}

// textColor tells the rendering textProvider our color
func (l *Label) textColor() color.Color {
	return theme.TextColor()
}

// password tells the rendering textProvider if we are a password field
func (l *Label) password() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (l *Label) object() fyne.Widget {
	return l
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	hidden := l.Hidden
	l.textProvider = newTextProvider(l.Text, l)
	l.textProvider.Hidden = hidden
	return l.textProvider.CreateRenderer()
}

// MinSize returns the size that this widget should not shrink below
func (l *Label) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)
	return l.BaseWidget.MinSize()
}
