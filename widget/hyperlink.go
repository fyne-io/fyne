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

	// OnTapped overrides the default `fyne.OpenURL` call when the link is tapped
	//
	// Since: 2.2
	OnTapped func() `json:"-"`

	textSize         fyne.Size // updated in syncSegments
	focused, hovered bool
	provider         *RichText
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
	hl.provider = NewRichTextWithText(hl.Text)
	hl.provider.ExtendBaseWidget(hl.provider)
	hl.syncSegments()

	focus := canvas.NewRectangle(color.Transparent)
	focus.StrokeColor = theme.FocusColor()
	focus.StrokeWidth = 2
	focus.Hide()
	under := canvas.NewRectangle(theme.HyperlinkColor())
	under.Hide()
	return &hyperlinkRenderer{hl: hl, objects: []fyne.CanvasObject{hl.provider, focus, under}, focus: focus, under: under}
}

// Cursor returns the cursor type of this widget
func (hl *Hyperlink) Cursor() desktop.Cursor {
	if hl.hovered {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (hl *Hyperlink) FocusGained() {
	hl.focused = true
	hl.BaseWidget.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (hl *Hyperlink) FocusLost() {
	hl.focused = false
	hl.BaseWidget.Refresh()
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (hl *Hyperlink) MouseIn(e *desktop.MouseEvent) {
	hl.MouseMoved(e)
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (hl *Hyperlink) MouseMoved(e *desktop.MouseEvent) {
	oldHovered := hl.hovered
	hl.hovered = hl.isPosOverText(e.Position)
	if hl.hovered != oldHovered {
		hl.BaseWidget.Refresh()
	}
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (hl *Hyperlink) MouseOut() {
	changed := hl.hovered
	hl.hovered = false
	if changed {
		hl.BaseWidget.Refresh()
	}
}

func (hl *Hyperlink) focusWidth() float32 {
	innerPad := theme.InnerPadding()
	return fyne.Min(hl.size.Width, hl.textSize.Width+innerPad+theme.Padding()*2) - innerPad
}

func (hl *Hyperlink) focusXPos() float32 {
	switch hl.Alignment {
	case fyne.TextAlignLeading:
		return theme.InnerPadding() / 2
	case fyne.TextAlignCenter:
		return (hl.size.Width - hl.focusWidth()) / 2
	case fyne.TextAlignTrailing:
		return (hl.size.Width - hl.focusWidth()) - theme.InnerPadding()/2
	default:
		return 0 // unreached
	}
}

func (hl *Hyperlink) isPosOverText(pos fyne.Position) bool {
	innerPad := theme.InnerPadding()
	pad := theme.Padding()
	// If not rendered yet provider will be nil
	lineCount := float32(1)
	if hl.provider != nil {
		lineCount = fyne.Max(lineCount, float32(len(hl.provider.rowBounds)))
	}

	xpos := hl.focusXPos()
	return pos.X >= xpos && pos.X <= xpos+hl.focusWidth() &&
		pos.Y >= innerPad/2 && pos.Y <= hl.textSize.Height*lineCount+pad*2+innerPad/2
}

// Refresh triggers a redraw of the hyperlink.
//
// Implements: fyne.Widget
func (hl *Hyperlink) Refresh() {
	if hl.provider == nil { // not created until visible
		return
	}
	hl.syncSegments()

	hl.provider.Refresh()
	hl.BaseWidget.Refresh()
}

// MinSize returns the smallest size this widget can shrink to
func (hl *Hyperlink) MinSize() fyne.Size {
	if hl.provider == nil {
		hl.CreateRenderer()
	}

	return hl.provider.MinSize()
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
	hl.syncSegments()
	hl.provider.Refresh()
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
func (hl *Hyperlink) Tapped(e *fyne.PointEvent) {
	// If not rendered yet (hl.provider == nil), register all taps
	// in practice this probably only happens in our unit tests
	if hl.provider != nil && !hl.isPosOverText(e.Position) {
		return
	}
	hl.invokeAction()
}

func (hl *Hyperlink) invokeAction() {
	if hl.OnTapped != nil {
		hl.OnTapped()
		return
	}
	hl.openURL()
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (hl *Hyperlink) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (hl *Hyperlink) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeySpace {
		hl.invokeAction()
	}
}

func (hl *Hyperlink) openURL() {
	if hl.URL != nil {
		err := fyne.CurrentApp().OpenURL(hl.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}
}

func (hl *Hyperlink) syncSegments() {
	hl.provider.Wrapping = hl.Wrapping
	hl.provider.Segments = []RichTextSegment{&TextSegment{
		Style: RichTextStyle{
			Alignment: hl.Alignment,
			ColorName: theme.ColorNameHyperlink,
			Inline:    true,
			TextStyle: hl.TextStyle,
		},
		Text: hl.Text,
	}}
	hl.textSize = fyne.MeasureText(hl.Text, theme.TextSize(), hl.TextStyle)
}

var _ fyne.WidgetRenderer = (*hyperlinkRenderer)(nil)

type hyperlinkRenderer struct {
	hl    *Hyperlink
	focus *canvas.Rectangle
	under *canvas.Rectangle

	objects []fyne.CanvasObject
}

func (r *hyperlinkRenderer) Destroy() {
}

func (r *hyperlinkRenderer) Layout(s fyne.Size) {
	innerPad := theme.InnerPadding()
	w := r.hl.focusWidth()
	xposFocus := r.hl.focusXPos()
	xposUnderline := xposFocus + innerPad/2

	r.hl.provider.Resize(s)
	lineCount := float32(len(r.hl.provider.rowBounds))
	r.focus.Move(fyne.NewPos(xposFocus, innerPad/2))
	r.focus.Resize(fyne.NewSize(w, r.hl.textSize.Height*lineCount+innerPad))
	r.under.Move(fyne.NewPos(xposUnderline, r.hl.textSize.Height*lineCount+theme.Padding()*2))
	r.under.Resize(fyne.NewSize(w-innerPad, 1))
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
	r.under.StrokeColor = theme.HyperlinkColor()
	r.under.Hidden = !r.hl.hovered
}
