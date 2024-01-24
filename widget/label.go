package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

// Label widget is a label component with appropriate padding and layout.
type Label struct {
	BaseWidget
	Text      string
	Alignment fyne.TextAlign // The alignment of the text
	Wrapping  fyne.TextWrap  // The wrapping of the text
	TextStyle fyne.TextStyle // The style of the label text

	// The truncation mode of the text
	//
	// Since: 2.4
	Truncation fyne.TextTruncation
	// Importance informs how the label should be styled, i.e. warning or disabled
	//
	// Since: 2.4
	Importance Importance

	provider *RichText
	binder   basicBinder
}

// NewLabel creates a new label widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewLabelWithData returns an Label widget connected to the specified data source.
//
// Since: 2.0
func NewLabelWithData(data binding.String) *Label {
	label := NewLabel("")
	label.Bind(data)

	return label
}

// NewLabelWithStyle creates a new label widget with the set text content
func NewLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle) *Label {
	l := &Label{
		Text:      text,
		Alignment: alignment,
		TextStyle: style,
	}

	l.ExtendBaseWidget(l)
	return l
}

// Bind connects the specified data source to this Label.
// The current value will be displayed and any changes in the data will cause the widget to update.
//
// Since: 2.0
func (l *Label) Bind(data binding.String) {
	l.binder.SetCallback(l.updateFromData) // This could only be done once, maybe in ExtendBaseWidget?
	l.binder.Bind(data)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	l.provider = NewRichTextWithText(l.Text)
	l.ExtendBaseWidget(l)
	l.syncSegments()

	return l.provider.CreateRenderer()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (l *Label) ExtendBaseWidget(w fyne.Widget) {
	if w == nil {
		w = l
	}
	l.BaseWidget.ExtendBaseWidget(w)
	if l.provider != nil {
		l.provider.ExtendBaseWidget(l.super())
	}
}

// MinSize returns the size that this label should not shrink below.
//
// Implements: fyne.Widget
func (l *Label) MinSize() fyne.Size {
	if l.provider == nil {
		l.ExtendBaseWidget(l)
		cache.Renderer(l.super())
	}

	return l.provider.MinSize()
}

// Refresh triggers a redraw of the label.
//
// Implements: fyne.Widget
func (l *Label) Refresh() {
	if l.provider == nil { // not created until visible
		return
	}
	l.syncSegments()
	l.provider.Refresh()
	l.BaseWidget.Refresh()
}

// Resize sets a new size for the label.
// This should only be called if it is not in a container with a layout manager.
//
// Implements: fyne.Widget
func (l *Label) Resize(s fyne.Size) {
	l.BaseWidget.Resize(s)
	if l.provider != nil {
		l.provider.Resize(s)
	}
}

// SetText sets the text of the label
func (l *Label) SetText(text string) {
	l.Text = text
	l.Refresh()
}

// Unbind disconnects any configured data source from this Label.
// The current value will remain at the last value of the data source.
//
// Since: 2.0
func (l *Label) Unbind() {
	l.binder.Unbind()
}

func (l *Label) syncSegments() {
	var color fyne.ThemeColorName
	switch l.Importance {
	case LowImportance:
		color = theme.ColorNameDisabled
	case MediumImportance:
		color = theme.ColorNameForeground
	case HighImportance:
		color = theme.ColorNamePrimary
	case DangerImportance:
		color = theme.ColorNameError
	case WarningImportance:
		color = theme.ColorNameWarning
	case SuccessImportance:
		color = theme.ColorNameSuccess
	default:
		color = theme.ColorNameForeground
	}

	l.provider.Wrapping = l.Wrapping
	l.provider.Truncation = l.Truncation
	l.provider.Segments[0].(*TextSegment).Style = RichTextStyle{
		Alignment: l.Alignment,
		ColorName: color,
		Inline:    true,
		TextStyle: l.TextStyle,
	}
	l.provider.Segments[0].(*TextSegment).Text = l.Text
}

func (l *Label) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		fyne.LogError("Error getting current data value", err)
		return
	}
	l.SetText(val)
}
