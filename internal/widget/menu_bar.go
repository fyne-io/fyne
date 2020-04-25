package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*MenuBar)(nil)

// MenuBar is a widget for displaying a fyne.MainMenu in a bar.
type MenuBar struct {
	base
	menuBase
	Items []fyne.CanvasObject

	active bool
	canvas fyne.Canvas
}

// NewMenuBar creates a menu bar populated with items from the passed main menu structure.
func NewMenuBar(mainMenu *fyne.MainMenu, canvas fyne.Canvas) *MenuBar {
	items := make([]fyne.CanvasObject, len(mainMenu.Items))
	b := &MenuBar{Items: items, canvas: canvas}
	for i, menu := range mainMenu.Items {
		items[i] = &MenuBarItem{Menu: menu, Parent: b}
	}
	return b
}

// CreateRenderer satisfies the fyne.Widget interface.
func (b *MenuBar) CreateRenderer() fyne.WidgetRenderer {
	cont := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), b.Items...)
	bg := &menuBarBackground{action: b.deactivate}
	return &menuBarRenderer{
		NewShadowingRenderer([]fyne.CanvasObject{bg, cont}, MenuLevel),
		b,
		bg,
		cont,
	}
}

// Hide satisfies the fyne.Widget interface.
func (b *MenuBar) Hide() {
	b.hide(b)
}

// MinSize satisfies the fyne.Widget interface.
func (b *MenuBar) MinSize() fyne.Size {
	return b.minSize(b)
}

// Refresh satisfies the fyne.Widget interface.
func (b *MenuBar) Refresh() {
	b.refresh(b)
}

// Resize satisfies the fyne.Widget interface.
func (b *MenuBar) Resize(size fyne.Size) {
	b.resize(size, b)
}

// Show satisfies the fyne.Widget interface.
func (b *MenuBar) Show() {
	b.show(b)
}

func (b *MenuBar) activate() {
	if b.active {
		return
	}

	b.active = true
	b.Refresh()
}

func (b *MenuBar) deactivate() {
	if !b.active {
		return
	}

	b.active = false
	b.dismiss()
	b.Refresh()
}

type menuBarRenderer struct {
	*ShadowingRenderer
	b    *MenuBar
	bg   *menuBarBackground
	cont *fyne.Container
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *menuBarRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuBarRenderer) Layout(size fyne.Size) {
	r.LayoutShadow(size, fyne.NewPos(0, 0))
	padding := r.padding()
	if r.b.active {
		r.bg.Resize(r.b.canvas.Size())
	} else {
		r.bg.Resize(fyne.NewSize(0, 0))
	}
	r.b.Resize(r.b.MinSize().Max(fyne.NewSize(size.Width, 0)))
	r.cont.Resize(size.Subtract(padding))
	r.cont.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(r.padding())
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	canvas.Refresh(r.b)
}

func (r *menuBarRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, 0)
}

// Transparent overlay shown as soon as menu is active.
// It catches mouse events outside the menu's objects.
type menuBarBackground struct {
	base
	action func()
}

var _ fyne.Widget = (*menuBarBackground)(nil)
var _ fyne.Tappable = (*menuBarBackground)(nil)     // unfocus menu on click outside
var _ desktop.Hoverable = (*menuBarBackground)(nil) // block hover events on main content

// CreateRenderer satisfies the fyne.Widget interface.
func (bg *menuBarBackground) CreateRenderer() fyne.WidgetRenderer {
	return &menuBackgroundRenderer{}
}

// Hide satisfies the fyne.Widget interface.
func (bg *menuBarBackground) Hide() {
	bg.hide(bg)
}

// MinSize satisfies the fyne.Widget interface.
func (bg *menuBarBackground) MinSize() fyne.Size {
	return bg.minSize(bg)
}

// MouseIn satisfies the desktop.Hoverable interface.
func (bg *menuBarBackground) MouseIn(*desktop.MouseEvent) {
}

// MouseOut satisfies the desktop.Hoverable interface.
func (bg *menuBarBackground) MouseOut() {
}

// MouseMoved satisfies the desktop.Hoverable interface.
func (bg *menuBarBackground) MouseMoved(*desktop.MouseEvent) {
}

// Refresh satisfies the fyne.Widget interface.
func (bg *menuBarBackground) Refresh() {
	bg.refresh(bg)
}

// Resize satisfies the fyne.Widget interface.
func (bg *menuBarBackground) Resize(size fyne.Size) {
	bg.resize(size, bg)
}

// Show satisfies the fyne.Widget interface.
func (bg *menuBarBackground) Show() {
	bg.show(bg)
}

// Tapped satisfies the fyne.Tappable interface.
func (bg *menuBarBackground) Tapped(*fyne.PointEvent) {
	bg.action()
}

type menuBackgroundRenderer struct {
	BaseRenderer
}

var _ fyne.WidgetRenderer = (*menuBackgroundRenderer)(nil)

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *menuBackgroundRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuBackgroundRenderer) Layout(fyne.Size) {
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuBackgroundRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuBackgroundRenderer) Refresh() {
}
