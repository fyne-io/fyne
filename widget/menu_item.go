package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Widget = (*menuItem)(nil)

// menuItem is a widget for displaying a fyne.menuItem.
type menuItem struct {
	widget.Base
	Item   *fyne.MenuItem
	Parent *Menu

	child *Menu
}

// newMenuItem creates a new menuItem.
func newMenuItem(item *fyne.MenuItem, parent *Menu) *menuItem {
	return &menuItem{Item: item, Parent: parent}
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
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	text := canvas.NewText(i.Item.Label, theme.ForegroundColor())
	objects := []fyne.CanvasObject{background, text}
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
		background:   background,
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

// MouseIn activates the item which shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseIn(*desktop.MouseEvent) {
	i.activate()
}

// MouseMoved does nothing.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut deactivates the item unless it has an open submenu.
//
// Implements: desktop.Hoverable
func (i *menuItem) MouseOut() {
	if !i.isSubmenuOpen() {
		i.deactivate()
	}
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
			i.activate()
		}
		return
	}

	i.trigger()
}

func (i *menuItem) activate() {
	if i.Child() != nil {
		i.Child().Show()
	}
	i.Parent.activateItem(i)
}

func (i *menuItem) activateLastSubmenu() bool {
	if i.Child() == nil {
		return false
	}
	if i.isSubmenuOpen() {
		return i.Child().ActivateLastSubmenu()
	}
	i.Child().Show()
	i.Child().ActivateNext()
	return true
}

func (i *menuItem) deactivate() {
	if i.Child() != nil {
		i.Child().Hide()
	}
	i.Parent.DeactivateChild()
}

func (i *menuItem) deactivateLastSubmenu() bool {
	if !i.isSubmenuOpen() {
		return false
	}
	if !i.Child().DeactivateLastSubmenu() {
		i.Child().DeactivateChild()
		i.Child().Hide()
	}
	return true
}

func (i *menuItem) isActive() bool {
	return i.Parent.activeItem == i
}

func (i *menuItem) isSubmenuOpen() bool {
	return i.Child() != nil && i.Child().Visible()
}

func (i *menuItem) trigger() {
	i.Parent.Dismiss()
	if i.Item.Action != nil {
		i.Item.Action()
	}
}

func (i *menuItem) triggerLast() {
	if i.isSubmenuOpen() {
		i.Child().TriggerLast()
		return
	}
	i.trigger()
}

type menuItemRenderer struct {
	widget.BaseRenderer
	i                *menuItem
	icon             *canvas.Image
	lastThemePadding float32
	minSize          fyne.Size
	text             *canvas.Text
	background       *canvas.Rectangle
}

func (r *menuItemRenderer) Layout(size fyne.Size) {
	padding := r.itemPadding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.ForegroundColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}

	r.background.Resize(size)
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
	if fyne.CurrentDevice().IsMobile() {
		r.background.Hide()
	} else if r.i.isActive() {
		r.background.FillColor = theme.FocusColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
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
