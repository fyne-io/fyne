package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

const noSelection = "(Select one)"

type selectRenderer struct {
	icon  *Icon
	label *canvas.Text

	objects []fyne.CanvasObject
	combo   *Select
}

// MinSize calculates the minimum size of a select button.
// This is based on the selected text, the drop icon and a standard amount of padding added.
func (s *selectRenderer) MinSize() fyne.Size {
	min := textMinSize(noSelection, s.label.TextSize, s.label.TextStyle)

	for _, option := range s.combo.Options {
		optionMin := textMinSize(option, s.label.TextSize, s.label.TextStyle)
		min = min.Union(optionMin)
	}

	min = min.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	return min.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))
}

// Layout the components of the button widget
func (s *selectRenderer) Layout(size fyne.Size) {
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

// ApplyTheme is called when the Button may need to update its look
func (s *selectRenderer) ApplyTheme() {
	s.label.Color = theme.TextColor()

	s.Refresh()
}

func (s *selectRenderer) BackgroundColor() color.Color {
	if s.combo.hovered {
		return theme.HoverColor()
	}
	return theme.ButtonColor()
}

func (s *selectRenderer) Refresh() {
	if s.combo.Selected == "" {
		s.label.Text = noSelection
	} else {
		s.label.Text = s.combo.Selected
	}

	if false { //s.combo.down {
		s.icon.Resource = theme.MenuDropUpIcon()
	} else {
		s.icon.Resource = theme.MenuDropDownIcon()
	}

	s.Layout(s.combo.Size())
	canvas.Refresh(s.combo)
}

func (s *selectRenderer) Objects() []fyne.CanvasObject {
	return s.objects
}

func (s *selectRenderer) Destroy() {
	if s.combo.popover != nil {
		c := fyne.CurrentApp().Driver().CanvasForObject(s.combo)
		c.SetOverlay(nil)
		Renderer(s.combo.popover).Destroy()
		s.combo.popover = nil
	}
}

// Select widget has a list of options, with the current one shown, and triggers an event func when clicked
type Select struct {
	baseWidget
	Selected string
	Options  []string

	OnChanged func(string) `json:"-"`
	hovered   bool
	popover   *PopOver
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *Select) Resize(size fyne.Size) {
	s.resize(size, s)

	if s.popover != nil {
		s.popover.Content.Resize(fyne.NewSize(size.Width, s.popover.MinSize().Height))
	}
}

// Move the widget to a new position, relative to its parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (s *Select) Move(pos fyne.Position) {
	s.move(pos, s)
}

// MinSize returns the smallest size this widget can shrink to
func (s *Select) MinSize() fyne.Size {
	return s.minSize(s)
}

// Show this widget, if it was previously hidden
func (s *Select) Show() {
	s.show(s)
}

// Hide this widget, if it was previously visible
func (s *Select) Hide() {
	s.hide(s)
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (s *Select) Tapped(*fyne.PointEvent) {
	c := fyne.CurrentApp().Driver().CanvasForObject(s)
	if s.popover != nil {
		c.SetOverlay(nil)
		Renderer(s.popover).Destroy()
		s.popover = nil
	}

	options := NewVBox()
	for _, option := range s.Options {
		text := option // capture value
		options.Append(newTappableLabel(text, func() {
			s.SetSelected(text)
			c.SetOverlay(nil)
			Renderer(options).Destroy()
			s.popover = nil
		}))
	}

	s.popover = NewPopOver(options, c)
	buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(s)
	popOverPos := buttonPos.Add(fyne.NewPos(0, s.Size().Height+theme.Padding()))

	if popOverPos.Y+s.popover.Content.Size().Height > c.Size().Height-theme.Padding()*3 {
		popOverPos.Y = c.Size().Height - s.popover.Content.Size().Height - theme.Padding()*3
		if popOverPos.Y < 0 {
			popOverPos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	s.popover.Content.Move(popOverPos)
	c.SetOverlay(s.popover)
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (s *Select) TappedSecondary(*fyne.PointEvent) {
}

// MouseIn is called when a desktop pointer enters the widget
func (s *Select) MouseIn(*desktop.MouseEvent) {
	s.hovered = true
	Refresh(s)
}

// MouseOut is called when a desktop pointer exits the widget
func (s *Select) MouseOut() {
	s.hovered = false
	Refresh(s)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (s *Select) MouseMoved(*desktop.MouseEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (s *Select) CreateRenderer() fyne.WidgetRenderer {
	icon := NewIcon(theme.MenuDropDownIcon())

	text := canvas.NewText(s.Selected, theme.TextColor())
	if s.Selected == "" {
		text.Text = noSelection
	}
	text.Alignment = fyne.TextAlignLeading

	objects := []fyne.CanvasObject{
		text, icon,
	}

	return &selectRenderer{icon, text, objects, s}
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

	Refresh(s)
}

// NewSelect creates a new select widget with the set list of options and changes handler
func NewSelect(options []string, changed func(string)) *Select {
	combo := &Select{baseWidget{}, "", options, changed, false, nil}

	Renderer(combo).Layout(combo.MinSize())
	return combo
}

type tappableLabel struct {
	*Label
	OnTapped func()
	hovered  bool
}

func (t *tappableLabel) Tapped(*fyne.PointEvent) {
	t.OnTapped()
}

func (t *tappableLabel) TappedSecondary(*fyne.PointEvent) {
}

func (t *tappableLabel) CreateRenderer() fyne.WidgetRenderer {
	return &hoverLabelRenderer{t.Label.CreateRenderer().(*textRenderer), t}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *tappableLabel) MouseIn(*desktop.MouseEvent) {
	t.hovered = true

	canvas.Refresh(t)
}

// MouseOut is called when a desktop pointer exits the widget
func (t *tappableLabel) MouseOut() {
	t.hovered = false

	canvas.Refresh(t)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *tappableLabel) MouseMoved(*desktop.MouseEvent) {
}

func newTappableLabel(label string, tapped func()) *tappableLabel {
	ret := &tappableLabel{NewLabel(label), tapped, false}
	Renderer(ret).Refresh() // trigger the textProvider to refresh metrics for our MinSize
	return ret
}

type hoverLabelRenderer struct {
	*textRenderer
	label *tappableLabel
}

func (h *hoverLabelRenderer) BackgroundColor() color.Color {
	if h.label.hovered {
		return theme.HoverColor()
	}

	return theme.BackgroundColor()
}
