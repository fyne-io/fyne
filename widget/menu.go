package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// ShowPopUpMenuAtPosition creates a PopUp menu populated with items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func ShowPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) {
	m := newPopUpMenu(menu, c)
	m.ShowAtPosition(pos)
}

func newPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *widget.PopUpMenu {
	m := widget.NewPopUpMenu(menu, c)
	focused := c.Focused()
	m.DismissAction = func() {
		if c.Focused() == nil {
			c.Focus(focused)
		}
		m.Hide()
	}
	return m
}

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
// Deprecated: Use ShowPopUpMenuAtPosition instead.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	options := NewVBox()
	for _, option := range menu.Items {
		opt := option // capture value
		if opt.IsSeparator {
			options.Append(newSeparator())
		} else {
			options.Append(newMenuItemWidget(opt.Label))
		}
	}
	pop := NewPopUpAtPosition(options, c, pos)
	focused := c.Focused()
	for i, o := range options.Children {
		if label, ok := o.(*menuItemWidget); ok {
			item := menu.Items[i]
			label.OnTapped = func() {
				if c.Focused() == nil {
					c.Focus(focused)
				}
				pop.Hide()
				item.Action()
			}
		}
	}
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
// Deprecated: Use ShowPopUpMenuAtPosition instead.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}

type menuItemWidget struct {
	*Label
	OnTapped func()
	hovered  bool
}

func (t *menuItemWidget) Tapped(*fyne.PointEvent) {
	t.OnTapped()
}

func (t *menuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	return &menuItemWidgetRenderer{t.Label.CreateRenderer().(*textRenderer), t}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *menuItemWidget) MouseIn(*desktop.MouseEvent) {
	t.hovered = true

	canvas.Refresh(t)
}

// MouseOut is called when a desktop pointer exits the widget
func (t *menuItemWidget) MouseOut() {
	t.hovered = false

	canvas.Refresh(t)
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *menuItemWidget) MouseMoved(*desktop.MouseEvent) {
}

func newMenuItemWidget(label string) *menuItemWidget {
	ret := &menuItemWidget{Label: NewLabel(label)}
	ret.ExtendBaseWidget(ret)
	return ret
}

type menuItemWidgetRenderer struct {
	*textRenderer
	label *menuItemWidget
}

func (h *menuItemWidgetRenderer) BackgroundColor() color.Color {
	if h.label.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

func newSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}
