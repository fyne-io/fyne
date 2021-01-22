package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const defaultPlaceHolder string = "(Select one)"

// Select widget has a list of options, with the current one shown, and triggers an event func when clicked
type Select struct {
	DisableableWidget

	Selected    string
	Options     []string
	PlaceHolder string
	OnChanged   func(string) `json:"-"`

	focused bool
	hovered bool
	popUp   *PopUpMenu
	tapAnim *fyne.Animation
	tapBG   *canvas.Rectangle
}

var _ fyne.Widget = (*Select)(nil)
var _ desktop.Hoverable = (*Select)(nil)
var _ fyne.Tappable = (*Select)(nil)
var _ fyne.Focusable = (*Select)(nil)
var _ fyne.Disableable = (*Select)(nil)

var _ textPresenter = (*Select)(nil)

// NewSelect creates a new select widget with the set list of options and changes handler
func NewSelect(options []string, changed func(string)) *Select {
	s := &Select{
		OnChanged:   changed,
		Options:     options,
		PlaceHolder: defaultPlaceHolder,
	}
	s.ExtendBaseWidget(s)
	return s
}

// ClearSelected clears the current option of the select widget.  After
// clearing the current option, the Select widget's PlaceHolder will
// be displayed.
func (s *Select) ClearSelected() {
	s.updateSelected("")
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Select) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	s.propertyLock.RLock()
	icon := NewIcon(theme.MenuDropDownIcon())
	if s.PlaceHolder == "" {
		s.PlaceHolder = defaultPlaceHolder
	}
	txtProv := newTextProvider(s.Selected, s)
	txtProv.ExtendBaseWidget(txtProv)

	background := &canvas.Rectangle{}
	line := canvas.NewRectangle(theme.ShadowColor())
	s.tapBG = canvas.NewRectangle(color.Transparent)
	objects := []fyne.CanvasObject{background, line, s.tapBG, txtProv, icon}
	r := &selectRenderer{icon, txtProv, background, line, objects, s}
	background.FillColor, line.FillColor = r.bgLineColor()
	r.updateIcon()
	s.propertyLock.RUnlock() // updateLabel and some text handling isn't quite right, resolve in text refactor for 2.0
	r.updateLabel()
	return r
}

// FocusGained is called after this Select has gained focus.
//
// Implements: fyne.Focusable
func (s *Select) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called after this Select has lost focus.
//
// Implements: fyne.Focusable
func (s *Select) FocusLost() {
	s.focused = false
	s.Refresh()
}

// Hide hides the select.
//
// Implements: fyne.Widget
func (s *Select) Hide() {
	if s.popUp != nil {
		s.popUp.Hide()
		s.popUp = nil
	}
	s.BaseWidget.Hide()
}

// MinSize returns the size that this widget should not shrink below
func (s *Select) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (s *Select) MouseIn(*desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (s *Select) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (s *Select) MouseOut() {
	s.hovered = false
	s.Refresh()
}

// Move changes the relative position of the select.
//
// Implements: fyne.Widget
func (s *Select) Move(pos fyne.Position) {
	s.BaseWidget.Move(pos)

	if s.popUp != nil {
		s.popUp.Move(s.popUpPos())
	}
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *Select) Resize(size fyne.Size) {
	s.BaseWidget.Resize(size)

	if s.popUp != nil {
		s.popUp.Resize(fyne.NewSize(size.Width, s.popUp.MinSize().Height))
	}
}

// SelectedIndex returns the index value of the currently selected item in Options list.
// It will return -1 if there is no selection.
func (s *Select) SelectedIndex() int {
	for i, option := range s.Options {
		if s.Selected == option {
			return i
		}
	}
	return -1 // not selected/found
}

// SetSelected sets the current option of the select widget
func (s *Select) SetSelected(text string) {
	for _, option := range s.Options {
		if text == option {
			s.updateSelected(text)
		}
	}
}

// SetSelectedIndex will set the Selected option from the value in Options list at index position.
func (s *Select) SetSelectedIndex(index int) {
	if index < 0 || index >= len(s.Options) {
		return
	}

	s.updateSelected(s.Options[index])
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (s *Select) Tapped(*fyne.PointEvent) {
	if s.Disabled() {
		return
	}

	s.tapAnimation()
	s.Refresh()

	s.showPopUp()
}

// TypedKey is called if a key event happens while this Select is focused.
//
// Implements: fyne.Focusable
func (s *Select) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace, fyne.KeyUp, fyne.KeyDown:
		s.showPopUp()
	case fyne.KeyRight:
		i := s.SelectedIndex() + 1
		if i >= len(s.Options) {
			i = 0
		}
		s.SetSelectedIndex(i)
	case fyne.KeyLeft:
		i := s.SelectedIndex() - 1
		if i < 0 {
			i = len(s.Options) - 1
		}
		s.SetSelectedIndex(i)
	}
}

// TypedRune is called if a text event happens while this Select is focused.
//
// Implements: fyne.Focusable
func (s *Select) TypedRune(_ rune) {
	// intentionally left blank
}

func (s *Select) concealed() bool {
	return false
}

func (s *Select) object() fyne.Widget {
	return nil
}

func (s *Select) optionTapped(text string) {
	s.SetSelected(text)
	s.popUp = nil
}

func (s *Select) popUpPos() fyne.Position {
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.super())
	return buttonPos.Add(fyne.NewPos(0, s.Size().Height-theme.InputBorderSize()))
}

func (s *Select) showPopUp() {
	var items []*fyne.MenuItem
	for _, option := range s.Options {
		text := option // capture
		item := fyne.NewMenuItem(option, func() {
			s.optionTapped(text)
		})
		items = append(items, item)
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(s.super())
	s.popUp = NewPopUpMenu(fyne.NewMenu("", items...), c)
	s.popUp.ShowAtPosition(s.popUpPos())
	s.popUp.Resize(fyne.NewSize(s.Size().Width, s.popUp.MinSize().Height))
}

func (s *Select) tapAnimation() {
	if s.tapBG == nil { // not rendered yet? (tests)
		return
	}

	if s.tapAnim == nil {
		s.tapAnim = newButtonTapAnimation(s.tapBG, s)
		s.tapAnim.Curve = fyne.AnimationEaseOut
	} else {
		s.tapAnim.Stop()
	}

	s.tapAnim.Start()
}

func (s *Select) textAlign() fyne.TextAlign {
	return fyne.TextAlignLeading
}

func (s *Select) textColor() color.Color {
	if s.Disabled() {
		return theme.DisabledColor()
	}
	return theme.ForegroundColor()
}

func (s *Select) textStyle() fyne.TextStyle {
	return fyne.TextStyle{Bold: false}
}

func (s *Select) textWrap() fyne.TextWrap {
	return fyne.TextTruncate
}

func (s *Select) updateSelected(text string) {
	s.Selected = text

	if s.OnChanged != nil {
		s.OnChanged(s.Selected)
	}

	s.Refresh()
}

type selectRenderer struct {
	icon             *Icon
	label            *textProvider
	background, line *canvas.Rectangle

	objects []fyne.CanvasObject
	combo   *Select
}

func (s *selectRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *selectRenderer) Destroy() {}

// Layout the components of the button widget
func (s *selectRenderer) Layout(size fyne.Size) {
	s.line.Resize(fyne.NewSize(size.Width, theme.InputBorderSize()))
	s.line.Move(fyne.NewPos(0, size.Height-theme.InputBorderSize()))
	s.background.Resize(fyne.NewSize(size.Width, size.Height-theme.InputBorderSize()*2))
	s.background.Move(fyne.NewPos(0, theme.InputBorderSize()))

	iconPos := fyne.NewPos(size.Width-theme.IconInlineSize()-theme.Padding()*2, (size.Height-theme.IconInlineSize())/2)
	labelSize := fyne.NewSize(iconPos.X-theme.Padding(), s.label.MinSize().Height)

	s.label.Resize(labelSize)
	s.label.Move(fyne.NewPos(theme.Padding(), (size.Height-labelSize.Height)/2))

	s.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	s.icon.Move(iconPos)
}

// MinSize calculates the minimum size of a select button.
// This is based on the selected text, the drop icon and a standard amount of padding added.
func (s *selectRenderer) MinSize() fyne.Size {
	s.combo.propertyLock.RLock()
	defer s.combo.propertyLock.RUnlock()

	minPlaceholderWidth := fyne.MeasureText(s.combo.PlaceHolder, theme.TextSize(), s.combo.textStyle()).Width
	min := s.label.MinSize()
	min.Width = minPlaceholderWidth
	min = min.Add(fyne.NewSize(theme.Padding()*6, theme.Padding()*2))
	return min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding()*2, 0))
}

func (s *selectRenderer) Refresh() {
	s.combo.propertyLock.RLock()
	s.updateLabel()
	s.updateIcon()
	s.background.FillColor, s.line.FillColor = s.bgLineColor()
	s.combo.propertyLock.RUnlock()

	s.Layout(s.combo.Size())
	if s.combo.popUp != nil {
		s.combo.Move(s.combo.position)
		s.combo.Resize(s.combo.size)
	}
	s.background.Refresh()
	canvas.Refresh(s.combo.super())
}

func (s *selectRenderer) bgLineColor() (bg color.Color, line color.Color) {
	if s.combo.Disabled() {
		return theme.InputBackgroundColor(), theme.DisabledTextColor()
	}
	if s.combo.focused {
		return theme.FocusColor(), theme.PrimaryColor()
	}
	if s.combo.hovered {
		return theme.HoverColor(), theme.ShadowColor()
	}
	return theme.InputBackgroundColor(), theme.ShadowColor()
}

func (s *selectRenderer) updateIcon() {
	if s.combo.Disabled() {
		s.icon.Resource = theme.NewDisabledResource(theme.MenuDropDownIcon())
	} else {
		s.icon.Resource = theme.MenuDropDownIcon()
	}
	s.icon.Refresh()
}

func (s *selectRenderer) updateLabel() {
	if s.combo.PlaceHolder == "" {
		s.combo.PlaceHolder = defaultPlaceHolder
	}

	if s.combo.Selected == "" {
		s.label.setText(s.combo.PlaceHolder)
	} else {
		s.label.setText(s.combo.Selected)
	}
}
