package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
	background := canvas.NewRectangle(theme.ButtonColor())
	underlay := &menuBarUnderlay{action: b.deactivate}
	objects := []fyne.CanvasObject{underlay, background, cont}
	for _, item := range b.Items {
		objects = append(objects, item.(*menuBarItem).Child())
	}
	return &menuBarRenderer{
		widget.NewShadowingRenderer(objects, widget.MenuLevel),
		b,
		background,
		underlay,
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

// Toggle changes the activation state of the menu bar.
// On activation, the first item will become active.
func (b *MenuBar) Toggle() {
	b.toggle(b.Items[0].(*menuBarItem))
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
		if c := b.activeItem.Child(); c != nil {
			c.Hide()
		}
		b.activeItem.Refresh()
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
		if c := b.activeItem.Child(); c != nil {
			defer c.Dismiss()
			c.Hide()
		}
		b.activeItem.Refresh()
		b.activeItem = nil
	}
	b.Refresh()
}

func (b *MenuBar) toggle(item *menuBarItem) {
	if b.active {
		b.canvas.Unfocus()
		b.deactivate()
	} else {
		b.activateChild(item)
		b.canvas.Focus(item)
	}
}

type menuBarRenderer struct {
	*widget.ShadowingRenderer
	b          *MenuBar
	background *canvas.Rectangle
	underlay   *menuBarUnderlay
	cont       *fyne.Container
}

func (r *menuBarRenderer) Layout(size fyne.Size) {
	r.LayoutShadow(size, fyne.NewPos(0, 0))
	minSize := r.MinSize()
	if size.Height != minSize.Height || size.Width < minSize.Width {
		r.b.Resize(fyne.NewSize(fyne.Max(size.Width, minSize.Width), minSize.Height))
		return
	}

	if r.b.active {
		r.underlay.Resize(r.b.canvas.Size())
	} else {
		r.underlay.Resize(fyne.NewSize(0, 0))
	}
	r.cont.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height))
	r.cont.Move(fyne.NewPos(theme.Padding(), 0))
	if item := r.b.activeItem; item != nil {
		if item.Child().Size().IsZero() {
			item.Child().Resize(item.Child().MinSize())
		}
		item.Child().Move(fyne.NewPos(item.Position().X+theme.Padding(), item.Size().Height))
	}
	r.background.Resize(size)
}

func (r *menuBarRenderer) MinSize() fyne.Size {
	return r.cont.MinSize().Add(fyne.NewSize(theme.Padding()*2, 0))
}

func (r *menuBarRenderer) Refresh() {
	r.Layout(r.b.Size())
	r.background.FillColor = theme.ButtonColor()
	r.background.Refresh()
	canvas.Refresh(r.b)
}

// Transparent underlay shown as soon as menu is active.
// It catches mouse events outside the menu's objects.
type menuBarUnderlay struct {
	widget.Base
	action func()
}

var _ fyne.Widget = (*menuBarUnderlay)(nil)
var _ fyne.Tappable = (*menuBarUnderlay)(nil)     // deactivate menu on click outside
var _ desktop.Hoverable = (*menuBarUnderlay)(nil) // block hover events on main content

func (u *menuBarUnderlay) CreateRenderer() fyne.WidgetRenderer {
	return &menuUnderlayRenderer{}
}

func (u *menuBarUnderlay) Hide() {
	widget.HideWidget(&u.Base, u)
}

func (u *menuBarUnderlay) MinSize() fyne.Size {
	return widget.MinSizeOf(u)
}

func (u *menuBarUnderlay) MouseIn(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) MouseOut() {
}

func (u *menuBarUnderlay) MouseMoved(*desktop.MouseEvent) {
}

func (u *menuBarUnderlay) Move(pos fyne.Position) {
	widget.MoveWidget(&u.Base, u, pos)
}

func (u *menuBarUnderlay) Refresh() {
	widget.RefreshWidget(u)
}

func (u *menuBarUnderlay) Resize(size fyne.Size) {
	widget.ResizeWidget(&u.Base, u, size)
}

func (u *menuBarUnderlay) Show() {
	widget.ShowWidget(&u.Base, u)
}

func (u *menuBarUnderlay) Tapped(*fyne.PointEvent) {
	u.action()
}

type menuUnderlayRenderer struct {
	widget.BaseRenderer
}

var _ fyne.WidgetRenderer = (*menuUnderlayRenderer)(nil)

func (r *menuUnderlayRenderer) Layout(fyne.Size) {
}

func (r *menuUnderlayRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *menuUnderlayRenderer) Refresh() {
}
