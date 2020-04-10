package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/theme"
)

// Label widget is a label component with appropriate padding and layout.
type Label struct {
	BaseWidget
	Text      string
	Alignment fyne.TextAlign // The alignment of the Text
	Wrapping  fyne.TextWrap  // The wrapping of the Text
	TextStyle fyne.TextStyle // The style of the label text

	provider textProvider

	textBinding *binding.StringBinding
	textNotify  *binding.NotifyFunction
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
	if l.Text != string(l.provider.buffer) {
		l.provider.SetText(l.Text)
	}

	l.BaseWidget.Refresh()
}

// SetText sets the text of the label
func (l *Label) SetText(text string) {
	if l.Text == text {
		return
	}
	l.Text = text
	l.provider.SetText(text) // calls refresh
}

// textAlign tells the rendering textProvider our alignment
func (l *Label) textAlign() fyne.TextAlign {
	return l.Alignment
}

// textWrap tells the rendering textProvider our wrapping
func (l *Label) textWrap() fyne.TextWrap {
	return l.Wrapping
}

// textStyle tells the rendering textProvider our style
func (l *Label) textStyle() fyne.TextStyle {
	return l.TextStyle
}

// textColor tells the rendering textProvider our color
func (l *Label) textColor() color.Color {
	return theme.TextColor()
}

// concealed tells the rendering textProvider if we are a concealed field
func (l *Label) concealed() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (l *Label) object() fyne.Widget {
	return l
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	l.provider = newTextProvider(l.Text, l)
	l.provider.Hidden = l.Hidden
	return l.provider.CreateRenderer()
}

// MinSize returns the size that this widget should not shrink below
func (l *Label) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)
	return l.BaseWidget.MinSize()
}

// BindText binds the Label's Text to the given data binding.
// Returns the Label for chaining.
func (l *Label) BindText(data *binding.StringBinding) *Label {
	l.textBinding = data
	l.textNotify = data.AddStringListener(l.SetText)
	return l
}

// UnbindText unbinds the Label's Text from the data binding (if any).
// Returns the Label for chaining.
func (l *Label) UnbindText() *Label {
	l.textBinding.DeleteListener(l.textNotify)
	l.textBinding = nil
	l.textNotify = nil
	return l
}
