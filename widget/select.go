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

	// Alignment sets the text alignment of the select and its list of options.
	//
	// Since: 2.1
	Alignment   fyne.TextAlign
	Selected    string
	Options     []string
	PlaceHolder string
	OnChanged   func(string) `json:"-"`

	focused bool
	hovered bool
	popUp   *PopUpMenu
	tapAnim *fyne.Animation
}

var _ fyne.Widget = (*Select)(nil)
var _ desktop.Hoverable = (*Select)(nil)
var _ fyne.Tappable = (*Select)(nil)
var _ fyne.Focusable = (*Select)(nil)
var _ fyne.Disableable = (*Select)(nil)

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
	txtProv := NewRichTextWithText(s.Selected)
	txtProv.inset = fyne.NewSize(theme.Padding(), theme.Padding())
	txtProv.ExtendBaseWidget(txtProv)
	txtProv.Truncation = fyne.TextTruncateEllipsis
	if s.disabled {
		txtProv.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameDisabled
	}

	background := &canvas.Rectangle{}
	tapBG := canvas.NewRectangle(color.Transparent)
	s.tapAnim = newButtonTapAnimation(tapBG, s)
	s.tapAnim.Curve = fyne.AnimationEaseOut
	objects := []fyne.CanvasObject{background, tapBG, txtProv, icon}
	r := &selectRenderer{icon, txtProv, background, objects, s}
	background.FillColor = r.bgColor()
	background.CornerRadius = theme.InputRadiusSize()
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

// SetOptions updates the list of options available and refreshes the widget
//
// Since: 2.4
func (s *Select) SetOptions(options []string) {
	s.Options = options
	s.Refresh()
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

func (s *Select) popUpPos() fyne.Position {
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.super())
	return buttonPos.Add(fyne.NewPos(0, s.Size().Height-theme.InputBorderSize()))
}

func (s *Select) showPopUp() {
	items := make([]*fyne.MenuItem, len(s.Options))
	for i := range s.Options {
		text := s.Options[i] // capture
		items[i] = fyne.NewMenuItem(text, func() {
			s.updateSelected(text)
			s.popUp = nil
		})
	}

	c := fyne.CurrentApp().Driver().CanvasForObject(s.super())
	s.popUp = NewPopUpMenu(fyne.NewMenu("", items...), c)
	s.popUp.alignment = s.Alignment
	s.popUp.ShowAtPosition(s.popUpPos())
	s.popUp.Resize(fyne.NewSize(s.Size().Width, s.popUp.MinSize().Height))
	s.popUp.OnDismiss = func() {
		s.popUp.Hide()
		s.popUp = nil
	}
}

func (s *Select) tapAnimation() {
	if s.tapAnim == nil {
		return
	}
	s.tapAnim.Stop()

	if fyne.CurrentApp().Settings().ShowAnimations() {
		s.tapAnim.Start()
	}
}

func (s *Select) updateSelected(text string) {
	s.Selected = text

	if s.OnChanged != nil {
		s.OnChanged(s.Selected)
	}

	s.Refresh()
}

type selectRenderer struct {
	icon       *Icon
	label      *RichText
	background *canvas.Rectangle

	objects []fyne.CanvasObject
	combo   *Select
}

func (s *selectRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *selectRenderer) Destroy() {}

// Layout the components of the button widget
func (s *selectRenderer) Layout(size fyne.Size) {
	s.background.Resize(fyne.NewSize(size.Width, size.Height))
	s.label.inset = fyne.NewSize(theme.Padding(), theme.Padding())

	iconPos := fyne.NewPos(size.Width-theme.IconInlineSize()-theme.InnerPadding(), (size.Height-theme.IconInlineSize())/2)
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

	minPlaceholderWidth := fyne.MeasureText(s.combo.PlaceHolder, theme.TextSize(), fyne.TextStyle{}).Width
	min := s.label.MinSize()
	min.Width = minPlaceholderWidth
	min = min.Add(fyne.NewSize(theme.InnerPadding()*3, theme.InnerPadding()))
	return min.Add(fyne.NewSize(theme.IconInlineSize()+theme.InnerPadding(), 0))
}

func (s *selectRenderer) Refresh() {
	s.combo.propertyLock.RLock()
	s.updateLabel()
	s.updateIcon()
	s.background.FillColor = s.bgColor()
	s.background.CornerRadius = theme.InputRadiusSize()
	s.combo.propertyLock.RUnlock()

	s.Layout(s.combo.Size())
	if s.combo.popUp != nil {
		s.combo.popUp.alignment = s.combo.Alignment
		s.combo.popUp.Move(s.combo.popUpPos())
		s.combo.popUp.Resize(fyne.NewSize(s.combo.size.Width, s.combo.popUp.MinSize().Height))
		s.combo.popUp.Refresh()
	}
	s.background.Refresh()
	canvas.Refresh(s.combo.super())
}

func (s *selectRenderer) bgColor() color.Color {
	if s.combo.disabled {
		return theme.DisabledButtonColor()
	}
	if s.combo.focused {
		return theme.FocusColor()
	}
	if s.combo.hovered {
		return theme.HoverColor()
	}
	return theme.InputBackgroundColor()
}

func (s *selectRenderer) updateIcon() {
	if s.combo.disabled {
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

	s.label.Segments[0].(*TextSegment).Style.Alignment = s.combo.Alignment
	if s.combo.disabled {
		s.label.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameDisabled
	} else {
		s.label.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameForeground
	}
	if s.combo.Selected == "" {
		s.label.Segments[0].(*TextSegment).Text = s.combo.PlaceHolder
	} else {
		s.label.Segments[0].(*TextSegment).Text = s.combo.Selected
	}
	s.label.Refresh()
}
