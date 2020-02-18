package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/theme"
)

const defaultPlaceHolder string = "(Select one)"

type selectRenderer struct {
	icon   *Icon
	label  *canvas.Text
	shadow fyne.CanvasObject

	objects []fyne.CanvasObject
	combo   *Select
}

// MinSize calculates the minimum size of a select button.
// This is based on the selected text, the drop icon and a standard amount of padding added.
func (s *selectRenderer) MinSize() fyne.Size {
	min := textMinSize(s.combo.PlaceHolder, s.label.TextSize, s.label.TextStyle)

	for _, option := range s.combo.Options {
		optionMin := textMinSize(option, s.label.TextSize, s.label.TextStyle)
		min = min.Union(optionMin)
	}

	min = min.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	return min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
}

// Layout the components of the button widget
func (s *selectRenderer) Layout(size fyne.Size) {
	s.shadow.Resize(size)
	inner := size.Subtract(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))

	offset := fyne.NewSize(theme.IconInlineSize(), 0)
	labelSize := inner.Subtract(offset)
	s.label.Resize(labelSize)
	s.label.Move(fyne.NewPos(theme.Padding()*2, theme.Padding()))

	s.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	s.icon.Move(fyne.NewPos(
		size.Width-theme.IconInlineSize()-theme.Padding()*2,
		(size.Height-theme.IconInlineSize())/2))
}

func (s *selectRenderer) BackgroundColor() color.Color {
	if s.combo.hovered {
		return theme.HoverColor()
	}
	return theme.ButtonColor()
}

func (s *selectRenderer) Refresh() {
	s.label.Color = theme.TextColor()
	s.label.TextSize = theme.TextSize()

	if s.combo.PlaceHolder == "" {
		s.combo.PlaceHolder = defaultPlaceHolder
	}

	if s.combo.Selected == "" {
		s.label.Text = s.combo.PlaceHolder
	} else {
		s.label.Text = s.combo.Selected
	}

	if false { // s.combo.down {
		s.icon.Resource = theme.MenuDropUpIcon()
	} else {
		s.icon.Resource = theme.MenuDropDownIcon()
	}

	s.Layout(s.combo.Size())
	canvas.Refresh(s.combo.super())
}

func (s *selectRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *selectRenderer) Destroy() {
	if s.combo.popUp != nil {
		c := fyne.CurrentApp().Driver().CanvasForObject(s.combo)
		c.SetOverlay(nil)
		cache.Renderer(s.combo.popUp).Destroy()
		s.combo.popUp = nil
	}
}

// Select widget has a list of options, with the current one shown, and triggers an event func when clicked
type Select struct {
	BaseWidget

	Selected    string
	Options     []string
	PlaceHolder string
	OnChanged   func(string) `json:"-"`

	hovered bool
	popUp   *PopUp
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *Select) Resize(size fyne.Size) {
	s.BaseWidget.Resize(size)

	if s.popUp != nil {
		s.popUp.Content.Resize(fyne.NewSize(size.Width, s.popUp.MinSize().Height))
	}
}

func (s *Select) optionTapped(text string) {
	s.SetSelected(text)
	s.popUp = nil
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (s *Select) Tapped(*fyne.PointEvent) {
	c := fyne.CurrentApp().Driver().CanvasForObject(s.super())

	var items []*fyne.MenuItem
	for _, option := range s.Options {
		text := option // capture
		item := fyne.NewMenuItem(option, func() {
			s.optionTapped(text)
		})
		items = append(items, item)
	}

	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.super())
	popUpPos := buttonPos.Add(fyne.NewPos(0, s.Size().Height))

	s.popUp = NewPopUpMenuAtPosition(fyne.NewMenu("", items...), c, popUpPos)
	s.popUp.Resize(fyne.NewSize(s.Size().Width, s.popUp.Content.MinSize().Height))
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (s *Select) TappedSecondary(*fyne.PointEvent) {
}

// MouseIn is called when a desktop pointer enters the widget
func (s *Select) MouseIn(*desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (s *Select) MouseOut() {
	s.hovered = false
	s.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (s *Select) MouseMoved(*desktop.MouseEvent) {
}

// MinSize returns the size that this widget should not shrink below
func (s *Select) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Select) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	icon := NewIcon(theme.MenuDropDownIcon())
	text := canvas.NewText(s.Selected, theme.TextColor())

	if s.PlaceHolder == "" {
		s.PlaceHolder = defaultPlaceHolder
	}
	if s.Selected == "" {
		text.Text = s.PlaceHolder
	}
	text.Alignment = fyne.TextAlignLeading

	shadow := newShadow(shadowAround, theme.Padding()/2)
	objects := []fyne.CanvasObject{
		text, icon, shadow,
	}

	return &selectRenderer{icon, text, shadow, objects, s}
}

// SetSelected sets the current option of the select widget
func (s *Select) SetSelected(text string) {
	for _, option := range s.Options {
		if text == option {
			s.Selected = text
		}
	}

	if s.OnChanged != nil {
		s.OnChanged(text)
	}

	s.Refresh()
}

// NewSelect creates a new select widget with the set list of options and changes handler
func NewSelect(options []string, changed func(string)) *Select {
	return &Select{BaseWidget{}, "", options, defaultPlaceHolder, changed, false, nil}
}
