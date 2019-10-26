package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	options := NewVBox()
	pop := NewPopUp(options, c)
	focused := c.Focused()
	for _, option := range menu.Items {
		opt := option // capture value
		options.Append(newTappableLabel(opt.Label, func() {
			if c.Focused() == nil {
				c.Focus(focused)
			}
			c.SetOverlay(nil)
			Renderer(pop).Destroy()

			opt.Action()
		}))
	}

	options.Resize(options.MinSize()) // make sure we have updated after appending
	return pop
}

type menuItemWidget struct {
	*Label
	OnTapped func()
	hovered  bool
}

func (t *menuItemWidget) Tapped(*fyne.PointEvent) {
	t.OnTapped()
}

func (t *menuItemWidget) TappedSecondary(*fyne.PointEvent) {
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

func newTappableLabel(label string, tapped func()) *menuItemWidget {
	ret := &menuItemWidget{NewLabel(label), tapped, false}
	Renderer(ret).Refresh() // trigger the textProvider to refresh metrics for our MinSize
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
