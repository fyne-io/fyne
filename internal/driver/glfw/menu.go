package glfw

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"
)

// “overlay” shown as soon as menu is active.
// It catches mouse events outside the menu's objects.
type menuBackground struct {
	widget.BaseWidget
	b *menuBar
}

var _ fyne.Widget = (*menuBackground)(nil)
var _ fyne.Tappable = (*menuBackground)(nil)     // unfocus menu on click outside
var _ desktop.Hoverable = (*menuBackground)(nil) // block hover events on main content

func (b *menuBackground) CreateRenderer() fyne.WidgetRenderer {
	return &menuBackgroundRenderer{}
}

func (b *menuBackground) MouseIn(*desktop.MouseEvent) {
}

func (b *menuBackground) MouseOut() {
}

func (b *menuBackground) MouseMoved(*desktop.MouseEvent) {
}

func (b *menuBackground) Tapped(*fyne.PointEvent) {
	b.b.bar.Deactivate()
}

type menuBackgroundRenderer struct {
}

var _ fyne.WidgetRenderer = (*menuBackgroundRenderer)(nil)

func (r *menuBackgroundRenderer) Layout(fyne.Size) {
}

func (r *menuBackgroundRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *menuBackgroundRenderer) Refresh() {
}

func (r *menuBackgroundRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *menuBackgroundRenderer) Objects() []fyne.CanvasObject {
	return nil
}

func (r *menuBackgroundRenderer) Destroy() {
}

type menuBar struct {
	widget.BaseWidget
	c      fyne.Canvas
	bar    *widget.MenuBarWidget
	active bool
}

var _ fyne.Widget = (*menuBar)(nil)

func (b *menuBar) CreateRenderer() fyne.WidgetRenderer {
	return &menuBarRenderer{b: b, bg: &menuBackground{b: b}}
}

type menuBarRenderer struct {
	b  *menuBar
	bg *menuBackground
}

var _ fyne.WidgetRenderer = (*menuBarRenderer)(nil)

func (r *menuBarRenderer) Layout(size fyne.Size) {
	if r.b.active {
		r.bg.Resize(r.b.c.Size())
	} else {
		r.bg.Resize(fyne.NewSize(0, 0))
	}
	r.b.bar.Resize(r.b.bar.MinSize().Max(fyne.NewSize(size.Width, 0)))
}

func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.b.bar.MinSize().Max(fyne.NewSize(r.b.c.Size().Width, 0))
}

func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	r.b.bar.Refresh()
}

func (r *menuBarRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *menuBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.b.bar}
}

func (r *menuBarRenderer) Destroy() {
}

func buildMenuBar(menus *fyne.MainMenu, c fyne.Canvas) *menuBar {
	b := &menuBar{c: c}
	b.ExtendBaseWidget(b)

	if menus.Items[0].Items[len(menus.Items[0].Items)-1].Label != "Quit" { // make sure the first menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", func() {
			fyne.CurrentApp().Quit()
		})
		menus.Items[0].Items = append(menus.Items[0].Items, fyne.NewMenuItemSeparator(), quitItem)
	}

	menuWidget := widget.NewMenuBarWidget(menus)
	menuWidget.ActivateAction = func() {
		b.active = true
		b.Refresh()
	}
	menuWidget.DismissAction = func() {
		b.active = false
		b.Refresh()
	}
	b.bar = menuWidget

	return b
}
