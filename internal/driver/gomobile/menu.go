package gomobile

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalWidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type menuLabel struct {
	*widget.Box
	label *widget.Label

	menu   *fyne.Menu
	bar    *widget.Box
	canvas *mobileCanvas
}

func (m *menuLabel) Tapped(*fyne.PointEvent) {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(m)
	widget.ShowPopUpMenuAtPosition(m.menu, m.canvas, fyne.NewPos(pos.X+m.Size().Width, pos.Y))

	// TODO use NewPopUpMenu in 2.0 once the Deprecated
	menu := m.canvas.Overlays().Top().(*internalWidget.OverlayContainer).Content.(*widget.Menu)
	menuDismiss := menu.OnDismiss // this dismisses the menu stack
	menu.OnDismiss = func() {
		menuDismiss()
		m.bar.Hide() // dismiss the overlay menu bar
		m.canvas.menu = nil
	}
}

func (m *menuLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.Renderer(m.Box)
}

func newMenuLabel(item *fyne.Menu, parent *widget.Box, c *mobileCanvas) *menuLabel {
	label := widget.NewLabel(item.Label)
	box := widget.NewHBox(layout.NewSpacer(), label, layout.NewSpacer(), widget.NewIcon(theme.MenuExpandIcon()))

	m := &menuLabel{box, label, item, parent, c}
	return m
}

func (c *mobileCanvas) showMenu(menu *fyne.MainMenu) {
	var panel *widget.Box
	top := widget.NewHBox(widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		panel.Hide()
		c.menu = nil
	}))
	panel = widget.NewVBox(top)
	for _, item := range menu.Items {
		panel.Append(newMenuLabel(item, panel, c))
	}
	shadow := canvas.NewHorizontalGradient(theme.ShadowColor(), color.Transparent)
	c.menu = fyne.NewContainer(panel, shadow)

	devicePadTopLeft, devicePadBottomRight := c.edgePadding()
	padY := devicePadTopLeft.Height + devicePadBottomRight.Height
	panel.Move(fyne.NewPos(devicePadTopLeft.Width, devicePadTopLeft.Height))
	panel.Resize(fyne.NewSize(panel.MinSize().Width+theme.Padding(), c.size.Height-padY))
	shadow.Resize(fyne.NewSize(theme.Padding()/2, c.size.Height-padY))
	shadow.Move(fyne.NewPos(panel.Size().Width+devicePadTopLeft.Width, devicePadTopLeft.Height))
}

func (d *mobileDriver) findMenu(win *window) *fyne.MainMenu {
	if win.menu != nil {
		return win.menu
	}

	matched := false
	for x := len(d.windows) - 1; x >= 0; x-- {
		w := d.windows[x]
		if !matched {
			if w == win {
				matched = true
			}
			continue
		}

		if w.(*window).menu != nil {
			return w.(*window).menu
		}
	}

	return nil
}
