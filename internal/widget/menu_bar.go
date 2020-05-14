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

// CreateRenderer returns a new renderer for the menu bar.
// Implements: fyne.Widget
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

// Hide hides the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) Hide() {
	b.hide(b)
}

// MinSize returns the minimal size of the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) MinSize() fyne.Size {
	return b.minSize(b)
}

// Refresh triggers a redraw of the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) Refresh() {
	b.refresh(b)
}

// Resize resizes the menu bar.
// It only affects the width because menu bars are always displayed with their minimal height.
// Implements: fyne.Widget
func (b *MenuBar) Resize(size fyne.Size) {
	b.resize(size, b)
}

// Show makes the menu bar visible.
// Implements: fyne.Widget
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

func (r *menuBarRenderer) BackgroundColor() color.Color {
	return theme.ButtonColor()
}

func (r *menuBarRenderer) Layout(size fyne.Size) {
	r.LayoutShadow(size, fyne.NewPos(0, 0))
	minSize := r.MinSize()
	if size.Height != minSize.Height || size.Width < minSize.Width {
		r.b.Resize(fyne.NewSize(fyne.Max(size.Width, minSize.Width), minSize.Height))
		return
	}

	padding := r.padding()
	if r.b.active {
		r.bg.Resize(r.b.canvas.Size())
	} else {
		r.bg.Resize(fyne.NewSize(0, 0))
	}
	r.cont.Resize(size.Subtract(padding))
	r.cont.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(r.padding())
}

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
var _ fyne.Tappable = (*menuBarBackground)(nil)     // deactivate menu on click outside
var _ desktop.Hoverable = (*menuBarBackground)(nil) // block hover events on main content

func (bg *menuBarBackground) CreateRenderer() fyne.WidgetRenderer {
	return &menuBackgroundRenderer{}
}

func (bg *menuBarBackground) Hide() {
	bg.hide(bg)
}

func (bg *menuBarBackground) MinSize() fyne.Size {
	return bg.minSize(bg)
}

func (bg *menuBarBackground) MouseIn(*desktop.MouseEvent) {
}

func (bg *menuBarBackground) MouseOut() {
}

func (bg *menuBarBackground) MouseMoved(*desktop.MouseEvent) {
}

func (bg *menuBarBackground) Refresh() {
	bg.refresh(bg)
}

func (bg *menuBarBackground) Resize(size fyne.Size) {
	bg.resize(size, bg)
}

func (bg *menuBarBackground) Show() {
	bg.show(bg)
}

func (bg *menuBarBackground) Tapped(*fyne.PointEvent) {
	bg.action()
}

type menuBackgroundRenderer struct {
	BaseRenderer
}

var _ fyne.WidgetRenderer = (*menuBackgroundRenderer)(nil)

func (r *menuBackgroundRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *menuBackgroundRenderer) Layout(fyne.Size) {
}

func (r *menuBackgroundRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *menuBackgroundRenderer) Refresh() {
}
