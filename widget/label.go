package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
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

	// The theme size name for the text size of the label
	//
	// Since: 2.6
	SizeName fyne.ThemeSizeName

	// If set to true, Selectable indicates that this label should support select interaction
	// to allow the text to be copied.
	//
	//Since: 2.6
	Selectable bool

	provider  *RichText
	binder    basicBinder
	selection *focusSelectable
}

// NewLabel creates a new label widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewLabelWithData returns a Label widget connected to the specified data source.
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

	l.selection = &focusSelectable{}
	l.selection.ExtendBaseWidget(l.selection)
	l.selection.focus = l.selection
	l.selection.style = l.TextStyle
	l.selection.theme = l.Theme()
	l.selection.provider = l.provider

	return &labelRenderer{l}
}

// MinSize returns the size that this label should not shrink below.
//
// Implements: fyne.Widget
func (l *Label) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)
	return l.BaseWidget.MinSize()
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

// SelectedText returns the text currently selected in this Label.
// If the label is not Selectable it will return an empty string.
// If there is no selection it will return the empty string.
//
// Since: 2.6
func (l *Label) SelectedText() string {
	if !l.Selectable || l.selection == nil {
		return ""
	}

	return l.selection.SelectedText()
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

	sizeName := l.SizeName
	if sizeName == "" {
		sizeName = theme.SizeNameText
	}
	l.provider.Wrapping = l.Wrapping
	l.provider.Truncation = l.Truncation
	l.provider.Segments[0].(*TextSegment).Style = RichTextStyle{
		Alignment: l.Alignment,
		ColorName: color,
		Inline:    true,
		TextStyle: l.TextStyle,
		SizeName:  sizeName,
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

type labelRenderer struct {
	l *Label
}

func (r *labelRenderer) Destroy() {
}

func (r *labelRenderer) Layout(s fyne.Size) {
	r.l.selection.Resize(s)
	r.l.provider.Resize(s)
}

func (r *labelRenderer) MinSize() fyne.Size {
	return r.l.provider.MinSize()
}

func (r *labelRenderer) Objects() []fyne.CanvasObject {
	if !r.l.Selectable {
		return []fyne.CanvasObject{r.l.provider}
	}

	return []fyne.CanvasObject{r.l.selection, r.l.provider}
}

func (r *labelRenderer) Refresh() {
	r.l.provider.Refresh()

	sel := r.l.selection
	if !r.l.Selectable || sel == nil {
		return
	}

	sel.sizeName = r.l.SizeName
	sel.style = r.l.TextStyle
	sel.theme = r.l.Theme()
	sel.Refresh()
}

type focusSelectable struct {
	selectable
}

func (f *focusSelectable) FocusGained() {
	f.focussed = true
	f.Refresh()
}

func (f *focusSelectable) FocusLost() {
	f.focussed = false
	f.Refresh()
}

func (f *focusSelectable) TypedKey(*fyne.KeyEvent) {
}

func (f *focusSelectable) TypedRune(rune) {
}
