package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const defaultPlaceHolder string = "(Select one)"

var (
	_ fyne.Widget       = (*Select)(nil)
	_ desktop.Hoverable = (*Select)(nil)
	_ fyne.Tappable     = (*Select)(nil)
	_ fyne.Focusable    = (*Select)(nil)
	_ fyne.Disableable  = (*Select)(nil)
)

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

	binder basicBinder

	focused bool
	hovered bool
	popUp   *PopUpMenu
	tapAnim *fyne.Animation
}

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

// NewSelectWithData returns a new select widget connected to the specified data source.
//
// Since: 2.6
func NewSelectWithData(options []string, data binding.String) *Select {
	sel := NewSelect(options, nil)
	sel.Bind(data)

	return sel
}

// Bind connects the specified data source to this select.
// The current value will be displayed and any changes in the data will cause the widget
// to update.
//
// Since: 2.6
func (s *Select) Bind(data binding.String) {
	s.binder.SetCallback(s.updateFromData)
	s.binder.Bind(data)

	s.OnChanged = func(_ string) {
		s.binder.CallWithData(s.writeData)
	}
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
	th := s.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	icon := NewIcon(th.Icon(theme.IconNameArrowDropDown))
	if s.PlaceHolder == "" {
		s.PlaceHolder = defaultPlaceHolder
	}
	txtProv := NewRichTextWithText(s.Selected)
	txtProv.inset = fyne.NewSquareSize(th.Size(theme.SizeNamePadding))
	txtProv.ExtendBaseWidget(txtProv)
	txtProv.Truncation = fyne.TextTruncateEllipsis
	if s.Disabled() {
		txtProv.Segments[0].(*TextSegment).Style.ColorName = theme.ColorNameDisabled
	}

	background := &canvas.Rectangle{}
	tapBG := canvas.NewRectangle(color.Transparent)
	s.tapAnim = newButtonTapAnimation(tapBG, s, th)
	s.tapAnim.Curve = fyne.AnimationEaseOut
	objects := []fyne.CanvasObject{background, tapBG, txtProv, icon}
	r := &selectRenderer{icon, txtProv, background, objects, s}
	background.FillColor = r.bgColor(th, v)
	background.CornerRadius = th.Size(theme.SizeNameInputRadius)
	r.updateIcon(th)
	r.updateLabel()
	return r
}

// FocusGained is called after this Select has gained focus.
func (s *Select) FocusGained() {
	s.focused = true
	s.Refresh()
}

// FocusLost is called after this Select has lost focus.
func (s *Select) FocusLost() {
	s.focused = false
	s.Refresh()
}

// Hide hides the select.
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

	if !s.focused {
		focusIfNotMobile(s.super())
	}

	s.tapAnimation()
	s.Refresh()

	s.showPopUp()
}

// TypedKey is called if a key event happens while this Select is focused.
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
func (s *Select) TypedRune(_ rune) {
	// intentionally left blank
}

// Unbind disconnects any configured data source from this Select.
// The current value will remain at the last value of the data source.
//
// Since: 2.6
func (s *Select) Unbind() {
	s.OnChanged = nil
	s.binder.Unbind()
}

func (s *Select) popUpPos() fyne.Position {
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.super())
	return buttonPos.Add(fyne.NewPos(0, s.Size().Height-s.Theme().Size(theme.SizeNameInputBorder)))
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
	pop := NewPopUpMenu(fyne.NewMenu("", items...), c)
	pop.alignment = s.Alignment
	pop.ShowAtPosition(s.popUpPos())
	pop.Resize(fyne.NewSize(s.Size().Width, pop.MinSize().Height))
	pop.OnDismiss = func() {
		pop.Hide()
		if s.popUp == pop {
			s.popUp = nil
		}
	}
	s.popUp = pop
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

func (s *Select) updateFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	stringSource, ok := data.(binding.String)
	if !ok {
		return
	}

	val, err := stringSource.Get()
	if err != nil {
		return
	}
	s.SetSelected(val)
}

func (s *Select) updateSelected(text string) {
	s.Selected = text

	if s.OnChanged != nil {
		s.OnChanged(s.Selected)
	}

	s.Refresh()
}

func (s *Select) writeData(data binding.DataItem) {
	if data == nil {
		return
	}
	stringTarget, ok := data.(binding.String)
	if !ok {
		return
	}
	currentValue, err := stringTarget.Get()
	if err != nil {
		return
	}
	if currentValue != s.Selected {
		err := stringTarget.Set(s.Selected)
		if err != nil {
			fyne.LogError(fmt.Sprintf("Failed to set binding value to %s", s.Selected), err)
		}
	}
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
	th := s.combo.Theme()
	pad := th.Size(theme.SizeNamePadding)
	iconSize := th.Size(theme.SizeNameInlineIcon)
	innerPad := th.Size(theme.SizeNameInnerPadding)
	s.background.Resize(fyne.NewSize(size.Width, size.Height))
	s.label.inset = fyne.NewSquareSize(pad)

	iconPos := fyne.NewPos(size.Width-iconSize-innerPad, (size.Height-iconSize)/2)
	labelSize := fyne.NewSize(iconPos.X-pad, s.label.MinSize().Height)

	s.label.Resize(labelSize)
	s.label.Move(fyne.NewPos(pad, (size.Height-labelSize.Height)/2))

	s.icon.Resize(fyne.NewSquareSize(iconSize))
	s.icon.Move(iconPos)
}

// MinSize calculates the minimum size of a select button.
// This is based on the selected text, the drop icon and a standard amount of padding added.
func (s *selectRenderer) MinSize() fyne.Size {
	th := s.combo.Theme()
	innerPad := th.Size(theme.SizeNameInnerPadding)

	minPlaceholderWidth := fyne.MeasureText(s.combo.PlaceHolder, th.Size(theme.SizeNameText), fyne.TextStyle{}).Width
	min := s.label.MinSize()
	min.Width = minPlaceholderWidth
	min = min.Add(fyne.NewSize(innerPad*3, innerPad))
	return min.Add(fyne.NewSize(th.Size(theme.SizeNameInlineIcon)+innerPad, 0))
}

func (s *selectRenderer) Refresh() {
	th := s.combo.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	s.updateLabel()
	s.updateIcon(th)
	s.background.FillColor = s.bgColor(th, v)
	s.background.CornerRadius = s.combo.Theme().Size(theme.SizeNameInputRadius)

	s.Layout(s.combo.Size())
	if s.combo.popUp != nil {
		s.combo.popUp.alignment = s.combo.Alignment
		s.combo.popUp.Move(s.combo.popUpPos())
		s.combo.popUp.Resize(fyne.NewSize(s.combo.Size().Width, s.combo.popUp.MinSize().Height))
		s.combo.popUp.Refresh()
	}
	s.background.Refresh()
	canvas.Refresh(s.combo.super())
}

func (s *selectRenderer) bgColor(th fyne.Theme, v fyne.ThemeVariant) color.Color {
	if s.combo.Disabled() {
		return th.Color(theme.ColorNameDisabledButton, v)
	}
	if s.combo.focused {
		return th.Color(theme.ColorNameFocus, v)
	}
	if s.combo.hovered {
		return th.Color(theme.ColorNameHover, v)
	}
	return th.Color(theme.ColorNameInputBackground, v)
}

func (s *selectRenderer) updateIcon(th fyne.Theme) {
	icon := th.Icon(theme.IconNameArrowDropDown)
	if s.combo.Disabled() {
		s.icon.Resource = theme.NewDisabledResource(icon)
	} else {
		s.icon.Resource = icon
	}
	s.icon.Refresh()
}

func (s *selectRenderer) updateLabel() {
	if s.combo.PlaceHolder == "" {
		s.combo.PlaceHolder = defaultPlaceHolder
	}

	segment := s.label.Segments[0].(*TextSegment)
	segment.Style.Alignment = s.combo.Alignment
	if s.combo.Disabled() {
		segment.Style.ColorName = theme.ColorNameDisabled
	} else {
		segment.Style.ColorName = theme.ColorNameForeground
	}
	if s.combo.Selected == "" {
		segment.Text = s.combo.PlaceHolder
	} else {
		segment.Text = s.combo.Selected
	}
	s.label.Refresh()
}
