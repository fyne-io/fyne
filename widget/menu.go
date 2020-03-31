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
	return &hoverLabelRenderer{t.Label.CreateRenderer().(*textRenderer), t}
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

type hoverLabelRenderer struct {
	*textRenderer
	label *menuItemWidget
}

func (h *hoverLabelRenderer) BackgroundColor() color.Color {
	if h.label.hovered {
		return theme.HoverColor()
	}

	return theme.BackgroundColor()
}

func newSeparator() fyne.CanvasObject {
	return canvas.NewLine(theme.DisabledTextColor())
}
