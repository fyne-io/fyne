package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	m := newMenuWidget(menu)
	pop := newPopUp(m, c)
	pop.WithoutPadding = true
	focused := c.Focused()
	for _, o := range m.Children {
		if item, ok := o.(*menuItemWidget); ok {
			item.DismissAction = func() {
				if c.Focused() == nil {
					c.Focus(focused)
				}
				pop.Hide()
			}
		}
	}
	pop.ShowAtPosition(pos)
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}

func newMenuItemWidget(label string, action func()) *menuItemWidget {
	ret := &menuItemWidget{Label: label, Action: action}
	ret.ExtendBaseWidget(ret)
	return ret
}

func newSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}

func newMenuWidget(menu *fyne.Menu) *Box {
	m := NewVBox()
	for _, item := range menu.Items {
		if item.IsSeparator {
			m.Append(newSeparator())
		} else {
			m.Append(newMenuItemWidget(item.Label, item.Action))
		}
	}
	return m
}

type menuItemWidget struct {
	BaseWidget
	Label         string
	Action        func()
	DismissAction func()
	hovered       bool
}

func (t *menuItemWidget) Tapped(*fyne.PointEvent) {
	t.Action()
	if t.DismissAction != nil {
		t.DismissAction()
	}
}

func (t *menuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(t.Label, theme.TextColor())
	return &menuItemWidgetRenderer{baseRenderer{[]fyne.CanvasObject{text}}, text, t}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *menuItemWidget) MouseIn(*desktop.MouseEvent) {
	t.hovered = true
	t.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (t *menuItemWidget) MouseOut() {
	t.hovered = false
	t.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *menuItemWidget) MouseMoved(*desktop.MouseEvent) {
}

type menuItemWidgetRenderer struct {
	baseRenderer
	text *canvas.Text
	w    *menuItemWidget
}

func (r *menuItemWidgetRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuItemWidgetRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuItemWidgetRenderer) Refresh() {
	if r.text.TextSize != theme.TextSize() {
		defer r.Layout(r.w.Size())
	}
	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	canvas.Refresh(r.text)
}

func (r *menuItemWidgetRenderer) BackgroundColor() color.Color {
	if r.w.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuItemWidgetRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
