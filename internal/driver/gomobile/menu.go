package gomobile

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type menuLabel struct {
	*widget.Label

	menu   *fyne.Menu
	canvas fyne.Canvas
}

func (m *menuLabel) Tapped(*fyne.PointEvent) {
	p := widget.NewPopUpMenu(m.menu, m.canvas)
	p.Move(fyne.NewPos(m.Label.Size().Width, m.Label.Position().Y))
}

func (m *menuLabel) TappedSecondary(*fyne.PointEvent) {
}

func (m *menuLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(m.Label)
}

func newMenuLabel(item *fyne.Menu, c *mobileCanvas) *menuLabel {
	return &menuLabel{widget.NewLabel(item.Label), item, c}
}

func (c *mobileCanvas) showMenu(menu *fyne.MainMenu) {
	var panel *widget.Box
	top := widget.NewHBox(widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		panel.Hide()
		c.menu = nil
	}))
	panel = widget.NewVBox(top)
	for _, item := range menu.Items {
		panel.Append(newMenuLabel(item, c))
	}
	shadow := canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
	c.menu = fyne.NewContainer(panel, shadow)
	panel.Resize(fyne.NewSize(panel.MinSize().Width+internal.ScaleInt(c, theme.Padding()), c.size.Height))
	shadow.Resize(fyne.NewSize(internal.ScaleInt(c, theme.Padding())/2, c.size.Height))
	shadow.Move(fyne.NewPos(panel.Size().Width, 0))
}
