package widget

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Focusable = (*Hyperlink)(nil)
var _ fyne.Widget = (*Hyperlink)(nil)

// Hyperlink widget is a text component with appropriate padding and layout.
// When clicked, the default web browser should open with a URL
type Hyperlink struct {
	BaseWidget
	Text      string
	URL       *url.URL
	Alignment fyne.TextAlign // The alignment of the Text
	Wrapping  fyne.TextWrap  // The wrapping of the Text
	TextStyle fyne.TextStyle // The style of the hyperlink text

	focused  bool
	provider *textProvider
}

// NewHyperlink creates a new hyperlink widget with the set text content
func NewHyperlink(text string, url *url.URL) *Hyperlink {
	return NewHyperlinkWithStyle(text, url, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewHyperlinkWithStyle creates a new hyperlink widget with the set text content
func NewHyperlinkWithStyle(text string, url *url.URL, alignment fyne.TextAlign, style fyne.TextStyle) *Hyperlink {
	hl := &Hyperlink{
		Text:      text,
		URL:       url,
		Alignment: alignment,
		TextStyle: style,
	}

	return hl
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (hl *Hyperlink) CreateRenderer() fyne.WidgetRenderer {
	hl.ExtendBaseWidget(hl)
	hl.provider = newTextProvider(hl.Text, hl)

	focus := canvas.NewRectangle(color.Transparent)
	focus.StrokeColor = theme.FocusColor()
	focus.StrokeWidth = 2
	focus.Hide()
	return &hyperlinkRenderer{hl: hl, objects: []fyne.CanvasObject{hl.provider, focus}, focus: focus}
}

// Cursor returns the cursor type of this widget
func (hl *Hyperlink) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (hl *Hyperlink) FocusGained() {
	hl.focused = true
	hl.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (hl *Hyperlink) FocusLost() {
	hl.focused = false
	hl.Refresh()
}

// MinSize returns the smallest size this widget can shrink to
func (hl *Hyperlink) MinSize() fyne.Size {
	hl.ExtendBaseWidget(hl)
	if p := hl.provider; p != nil && hl.Text != string(p.buffer) {
		p.setText(hl.Text)
	}

	label := canvas.NewText(hl.Text, theme.ForegroundColor())
	label.TextStyle = hl.textStyle()
	return label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*4))
}

// Resize sets a new size for the hyperlink.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (hl *Hyperlink) Resize(size fyne.Size) {
	hl.BaseWidget.Resize(size)
	if hl.provider == nil { // not created until visible
		return
	}
	hl.provider.Resize(size)
}

// SetText sets the text of the hyperlink
func (hl *Hyperlink) SetText(text string) {
	hl.Text = text
	if hl.provider == nil { // not created until visible
		return
	}
	hl.provider.setText(text) // calls refresh
}

// SetURL sets the URL of the hyperlink, taking in a URL type
func (hl *Hyperlink) SetURL(url *url.URL) {
	hl.URL = url
}

// SetURLFromString sets the URL of the hyperlink, taking in a string type
func (hl *Hyperlink) SetURLFromString(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	hl.URL = u
	return nil
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (hl *Hyperlink) Tapped(*fyne.PointEvent) {
	hl.openURL()
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (hl *Hyperlink) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (hl *Hyperlink) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeySpace {
		hl.openURL()
	}
}

// textAlign tells the rendering textProvider our alignment
func (hl *Hyperlink) textAlign() fyne.TextAlign {
	return hl.Alignment
}

// textWrap tells the rendering textProvider our wrapping
func (hl *Hyperlink) textWrap() fyne.TextWrap {
	return hl.Wrapping
}

// textStyle tells the rendering textProvider our style
func (hl *Hyperlink) textStyle() fyne.TextStyle {
	return hl.TextStyle
}

// textColor tells the rendering textProvider our color
func (hl *Hyperlink) textColor() color.Color {
	return theme.PrimaryColor()
}

// concealed tells the rendering textProvider if we are a concealed field
func (hl *Hyperlink) concealed() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (hl *Hyperlink) object() fyne.Widget {
	return hl.provider
}

func (hl *Hyperlink) openURL() {
	if hl.URL != nil {
		err := fyne.CurrentApp().OpenURL(hl.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}
}

var _ fyne.WidgetRenderer = (*hyperlinkRenderer)(nil)

type hyperlinkRenderer struct {
	hl    *Hyperlink
	focus *canvas.Rectangle

	objects []fyne.CanvasObject
}

func (r *hyperlinkRenderer) Destroy() {
}

func (r *hyperlinkRenderer) Layout(s fyne.Size) {
	r.hl.provider.Resize(s)
	r.focus.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	r.focus.Resize(fyne.NewSize(s.Width-theme.Padding()*2, s.Height-theme.Padding()*2))
}

func (r *hyperlinkRenderer) MinSize() fyne.Size {
	return r.hl.provider.MinSize()
}

func (r *hyperlinkRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *hyperlinkRenderer) Refresh() {
	r.hl.provider.Refresh()
	r.focus.StrokeColor = theme.FocusColor()
	r.focus.Hidden = !r.hl.focused
}
