package glfw

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	publicWidget "fyne.io/fyne/widget"
)

var _ fyne.Widget = (*MenuBar)(nil)

// MenuBar is a widget for displaying a fyne.MainMenu in a bar.
type MenuBar struct {
	widget.Base
	Items []fyne.CanvasObject

	active      bool
	activeChild *publicWidget.Menu
	canvas      fyne.Canvas
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
		widget.NewShadowingRenderer([]fyne.CanvasObject{bg, cont}, widget.MenuLevel),
		b,
		bg,
		cont,
	}
}

// Hide hides the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) Hide() {
	widget.HideWidget(&b.Base, b)
}

// MinSize returns the minimal size of the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) MinSize() fyne.Size {
	return widget.MinSizeOf(b)
}

// Refresh triggers a redraw of the menu bar.
// Implements: fyne.Widget
func (b *MenuBar) Refresh() {
	widget.RefreshWidget(b)
}

// Resize resizes the menu bar.
// It only affects the width because menu bars are always displayed with their minimal height.
// Implements: fyne.Widget
func (b *MenuBar) Resize(size fyne.Size) {
	widget.ResizeWidget(&b.Base, b, size)
}

// Show makes the menu bar visible.
// Implements: fyne.Widget
func (b *MenuBar) Show() {
	widget.ShowWidget(&b.Base, b)
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
	if b.activeChild != nil {
		defer b.activeChild.Dismiss()
		b.activeChild.Hide()
		b.activeChild = nil
	}
	b.Refresh()
}

type menuBarRenderer struct {
	*widget.ShadowingRenderer
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
	widget.Base
	action func()
}

var _ fyne.Widget = (*menuBarBackground)(nil)
var _ fyne.Tappable = (*menuBarBackground)(nil)     // deactivate menu on click outside
var _ desktop.Hoverable = (*menuBarBackground)(nil) // block hover events on main content

func (bg *menuBarBackground) CreateRenderer() fyne.WidgetRenderer {
	return &menuBackgroundRenderer{}
}

func (bg *menuBarBackground) Hide() {
	widget.HideWidget(&bg.Base, bg)
}

func (bg *menuBarBackground) MinSize() fyne.Size {
	return widget.MinSizeOf(bg)
}

func (bg *menuBarBackground) MouseIn(*desktop.MouseEvent) {
}

func (bg *menuBarBackground) MouseOut() {
}

func (bg *menuBarBackground) MouseMoved(*desktop.MouseEvent) {
}

func (bg *menuBarBackground) Refresh() {
	widget.RefreshWidget(bg)
}

func (bg *menuBarBackground) Resize(size fyne.Size) {
	widget.ResizeWidget(&bg.Base, bg, size)
}

func (bg *menuBarBackground) Show() {
	widget.ShowWidget(&bg.Base, bg)
}

func (bg *menuBarBackground) Tapped(*fyne.PointEvent) {
	bg.action()
}

type menuBackgroundRenderer struct {
	widget.BaseRenderer
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
