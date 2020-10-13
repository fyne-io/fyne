package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*menuItem)(nil)

// menuItem is a widget for displaying a fyne.menuItem.
type menuItem struct {
	widget.Base
	Item   *fyne.MenuItem
	Parent *Menu

	child           *Menu
	hovered         bool
	onActivateChild func(*menuItem)
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *fyne.MenuItem, parent *Menu, onActivateChild func(*menuItem)) *menuItem {
	return &menuItem{Item: item, Parent: parent, onActivateChild: onActivateChild}
}

func (i *menuItem) Child() *Menu {
	if i.Item.ChildMenu != nil && i.child == nil {
		child := NewMenu(i.Item.ChildMenu)
		child.Hide()
		child.OnDismiss = i.Parent.Dismiss
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Item.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	var icon *canvas.Image
	if i.Item.ChildMenu != nil {
		icon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, icon)
	}
	return &menuItemRenderer{
		BaseRenderer: widget.NewBaseRenderer(objects),
		i:            i,
		icon:         icon,
		text:         text,
	}
}

// Hide hides the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) Hide() {
	widget.HideWidget(&i.Base, i)
}

// MinSize returns the minimal size of the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) MinSize() fyne.Size {
	return widget.MinSizeOf(i)
}

// MouseIn changes the item to be hovered and shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.onActivateChild(i)
	i.Refresh()
}

// MouseMoved does nothing.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut changes the item to not be hovered but has no effect on the visibility of the item's submenu.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Move sets the position of the widget relative to its parent.
//
// Implements: fyne.Widget
func (i *menuItem) Move(pos fyne.Position) {
	widget.MoveWidget(&i.Base, i, pos)
}

// Refresh triggers a redraw of the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) Refresh() {
	widget.RefreshWidget(i)
}

// Resize changes the size of the menu item.
//
// Implements: fyne.Widget
func (i *menuItem) Resize(size fyne.Size) {
	widget.ResizeWidget(&i.Base, i, size)
}

// Show makes the menu item visible.
//
// Implements: fyne.Widget
func (i *menuItem) Show() {
	widget.ShowWidget(&i.Base, i)
}

// Tapped performs the action of the item and dismisses the menu.
// It does nothing if the item doesn’t have an action.
//
// Implements: fyne.Tappable
func (i *menuItem) Tapped(*fyne.PointEvent) {
	if i.Item.Action == nil {
		if fyne.CurrentDevice().IsMobile() {
			i.onActivateChild(i)
		}
		return
	}

	i.Parent.Dismiss()
	i.Item.Action()
}

type menuItemRenderer struct {
	widget.BaseRenderer
	i                *menuItem
	icon             *canvas.Image
	lastThemePadding int
	minSize          fyne.Size
	text             *canvas.Text
}

func (r *menuItemRenderer) BackgroundColor() color.Color {
	if !fyne.CurrentDevice().IsMobile() && (r.i.hovered || (r.i.child != nil && r.i.child.Visible())) {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	padding := r.itemPadding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}
}

func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	minSize := r.text.MinSize().Add(r.itemPadding())
	if r.icon != nil {
		minSize = minSize.Add(fyne.NewSize(theme.IconInlineSize(), 0))
	}
	r.minSize = minSize
	return r.minSize
}

func (r *menuItemRenderer) Refresh() {
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	return !r.minSize.IsZero() &&
		r.text.TextSize == theme.TextSize() &&
		(r.icon == nil || r.icon.Size().Width == theme.IconInlineSize()) &&
		r.lastThemePadding == theme.Padding()
}

func (r *menuItemRenderer) itemPadding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
