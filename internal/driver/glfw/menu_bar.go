package glfw

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*MenuBar)(nil)

// MenuBar is a widget for displaying a fyne.MainMenu in a bar.
type MenuBar struct {
	widget.Base
	Items []fyne.CanvasObject

	active     bool
	activeItem *menuBarItem
	canvas     fyne.Canvas
}

// NewMenuBar creates a menu bar populated with items from the passed main menu structure.
func NewMenuBar(mainMenu *fyne.MainMenu, canvas fyne.Canvas) *MenuBar {
	items := make([]fyne.CanvasObject, len(mainMenu.Items))
	b := &MenuBar{Items: items, canvas: canvas}
	for i, menu := range mainMenu.Items {
		items[i] = &menuBarItem{Menu: menu, Parent: b}
	}
	return b
}

// CreateRenderer returns a new renderer for the menu bar.
//
// Implements: fyne.Widget
func (b *MenuBar) CreateRenderer() fyne.WidgetRenderer {
	cont := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), b.Items...)
	bg := &menuBarBackground{action: b.deactivate}
	objects := []fyne.CanvasObject{bg, cont}
	for _, item := range b.Items {
		objects = append(objects, item.(*menuBarItem).Child())
	}
	return &menuBarRenderer{
		widget.NewShadowingRenderer(objects, widget.MenuLevel),
		b,
		bg,
		cont,
	}
}

// IsActive returns whether the menu bar is active or not.
// An active menu bar shows the current selected menu and should have the focus.
func (b *MenuBar) IsActive() bool {
	return b.active
}

// Hide hides the menu bar.
//
// Implements: fyne.Widget
func (b *MenuBar) Hide() {
	widget.HideWidget(&b.Base, b)
}

// MinSize returns the minimal size of the menu bar.
//
// Implements: fyne.Widget
func (b *MenuBar) MinSize() fyne.Size {
	return widget.MinSizeOf(b)
}

// Move sets the position of the widget relative to its parent.
//
// Implements: fyne.Widget
func (b *MenuBar) Move(pos fyne.Position) {
	widget.MoveWidget(&b.Base, b, pos)
}

// Refresh triggers a redraw of the menu bar.
//
// Implements: fyne.Widget
func (b *MenuBar) Refresh() {
	widget.RefreshWidget(b)
}

// Resize resizes the menu bar.
// It only affects the width because menu bars are always displayed with their minimal height.
//
// Implements: fyne.Widget
func (b *MenuBar) Resize(size fyne.Size) {
	widget.ResizeWidget(&b.Base, b, size)
}

// Show makes the menu bar visible.
//
// Implements: fyne.Widget
func (b *MenuBar) Show() {
	widget.ShowWidget(&b.Base, b)
}

func (b *MenuBar) activateChild(item *menuBarItem) {
	if !b.active {
		b.active = true
	}
	if item.Child() != nil {
		item.Child().DeactivateChild()
	}
	if b.activeItem == item {
		return
	}

	if b.activeItem != nil {
		b.activeItem.Child().Hide()
	}
	b.activeItem = item
	if item == nil {
		return
	}

	item.Child().Show()
	b.Refresh()
}

func (b *MenuBar) deactivate() {
	if !b.active {
		return
	}

	b.active = false
	if b.activeItem != nil {
		defer b.activeItem.Child().Dismiss()
		b.activeItem.Child().Hide()
		b.activeItem = nil
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

	if r.b.active {
		r.bg.Resize(r.b.canvas.Size())
	} else {
		r.bg.Resize(fyne.NewSize(0, 0))
	}
	r.cont.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height))
	r.cont.Move(fyne.NewPos(theme.Padding(), 0))
	if item := r.b.activeItem; item != nil {
		if item.Child().Size().IsZero() {
			item.Child().Resize(item.Child().MinSize())
		}
		item.Child().Move(fyne.NewPos(item.Position().X+theme.Padding(), item.Size().Height))
	}
}

func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(fyne.NewSize(theme.Padding()*2, 0))
}

func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	canvas.Refresh(r.b)
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

func (bg *menuBarBackground) Move(pos fyne.Position) {
	widget.MoveWidget(&bg.Base, bg, pos)
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
