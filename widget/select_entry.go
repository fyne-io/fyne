package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// SelectEntry is an input field which supports selecting from a fixed set of options.
type SelectEntry struct {
	Entry
	dropDown *fyne.Menu
	popUp    *PopUp
}

// NewSelectEntry creates a SelectEntry.
func NewSelectEntry(options []string) *SelectEntry {
	e := &SelectEntry{}
	e.ExtendBaseWidget(e)
	e.SetOptions(options)
	return e
}

// MinSize satisfies the fyne.CanvasObject interface.
func (e *SelectEntry) MinSize() fyne.Size {
	min := e.Entry.MinSize()
	if e.dropDown != nil {
		for _, item := range e.dropDown.Items {
			itemMin := textMinSize(item.Label, theme.TextSize(), fyne.TextStyle{}).Add(fyne.NewSize(4*theme.Padding(), 0))
			min = min.Union(itemMin)
		}
	}
	return min
}

// SetOptions sets the options the user might select from.
func (e *SelectEntry) SetOptions(options []string) {
	if len(options) == 0 {
		e.ActionItem = nil
		return
	}

	var items []*fyne.MenuItem
	for _, option := range options {
		option := option // capture
		items = append(items, fyne.NewMenuItem(option, func() { e.SetText(option) }))
	}
	e.dropDown = fyne.NewMenu("", items...)
	e.ActionItem = newDropDownSwitch(func() {
		c := fyne.CurrentApp().Driver().CanvasForObject(e.super())

		entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(e.super())
		popUpPos := entryPos.Add(fyne.NewPos(0, e.Size().Height))

		e.popUp = NewPopUpMenuAtPosition(fyne.NewMenu("", items...), c, popUpPos)
		e.popUp.Resize(fyne.NewSize(e.Size().Width, e.popUp.MinSize().Height))
	})
}

type dropDownSwitch struct {
	BaseWidget
	icon    *canvas.Image
	onClick func()
}

var _ fyne.Tappable = (*dropDownSwitch)(nil)

func newDropDownSwitch(onClick func()) *dropDownSwitch {
	s := &dropDownSwitch{
		icon:    canvas.NewImageFromResource(theme.MenuDropDownIcon()),
		onClick: onClick,
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *dropDownSwitch) CreateRenderer() fyne.WidgetRenderer {
	return &dropDownSwitchRenderer{icon: s.icon}
}

func (s *dropDownSwitch) Tapped(*fyne.PointEvent) {
	s.onClick()
}

func (s *dropDownSwitch) TappedSecondary(*fyne.PointEvent) {
}

type dropDownSwitchRenderer struct {
	icon *canvas.Image
}

var _ fyne.WidgetRenderer = (*dropDownSwitchRenderer)(nil)

func (r *dropDownSwitchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *dropDownSwitchRenderer) Layout(size fyne.Size) {
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *dropDownSwitchRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *dropDownSwitchRenderer) Refresh() {
}

func (r *dropDownSwitchRenderer) Destroy() {
}

func (r *dropDownSwitchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.icon}
}
