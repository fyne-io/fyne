package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*MenuItem)(nil)

// MenuItem is a widget for displaying a fyne.MenuItem.
type MenuItem struct {
	base
	menuItemBase
	Item   *fyne.MenuItem
	Parent *Menu

	hovered bool
}

// NewMenuItem creates a new MenuItem.
func NewMenuItem(item *fyne.MenuItem, parent *Menu) *MenuItem {
	return &MenuItem{Item: item, Parent: parent}
}

// NewMenuItemSeparator creates a separator meant to separate MenuItems.
func NewMenuItemSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}

// CreateRenderer returns a new renderer for the menu item.
// Implements: fyne.Widget
func (i *MenuItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Item.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	var icon *canvas.Image
	if i.Item.ChildMenu != nil {
		icon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, icon)
		i.initChildWidget(i.Item.ChildMenu, i.Parent.dismiss)
		objects = append(objects, i.Child)
	}
	return &menuItemRenderer{
		BaseRenderer: NewBaseRenderer(objects),
		i:            i,
		icon:         icon,
		text:         text,
	}
}

// Hide hides the menu item.
// Implements: fyne.Widget
func (i *MenuItem) Hide() {
	i.hide(i)
}

// MinSize returns the minimal size of the menu item.
// Implements: fyne.Widget
func (i *MenuItem) MinSize() fyne.Size {
	return i.minSize(i)
}

// MouseIn changes the item to be hovered and shows the submenu if the item has one.
// The submenu of any sibling of the item will be hidden.
// Implements: desktop.Hoverable
func (i *MenuItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.activateChild(&i.Parent.menuBase, i.updateChildPosition)
	i.Refresh()
}

// MouseMoved does nothing.
// Implements: desktop.Hoverable
func (i *MenuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut changes the item to not be hovered but has no effect on the visibility of the item's submenu.
// Implements: desktop.Hoverable
func (i *MenuItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Refresh triggers a redraw of the menu item.
// Implements: fyne.Widget
func (i *MenuItem) Refresh() {
	i.refresh(i)
}

// Resize changes the size of the menu item.
// Implements: fyne.Widget
func (i *MenuItem) Resize(size fyne.Size) {
	i.resize(size, i)
	if i.Child != nil {
		i.updateChildPosition()
	}
}

// Show makes the menu item visible.
// Implements: fyne.Widget
func (i *MenuItem) Show() {
	i.show(i)
}

// Tapped performs the action of the item and dismisses the menu.
// It does nothing if the item doesn’t have an action.
// Implements: fyne.Tappable
func (i *MenuItem) Tapped(*fyne.PointEvent) {
	if i.Item.Action == nil {
		if fyne.CurrentDevice().IsMobile() {
			i.activateChild(&i.Parent.menuBase, i.updateChildPosition)
			i.Refresh()
		}
		return
	}

	i.Item.Action()
	i.Parent.dismiss()
}

func (i *MenuItem) updateChildPosition() {
	itemSize := i.Size()
	cp := fyne.NewPos(itemSize.Width, -theme.Padding())
	d := fyne.CurrentApp().Driver()
	c := d.CanvasForObject(i)
	if c != nil {
		absPos := d.AbsolutePositionForObject(i)
		childSize := i.Child.Size()
		if absPos.X+itemSize.Width+childSize.Width > c.Size().Width {
			if absPos.X-childSize.Width >= 0 {
				cp.X = -childSize.Width
			} else {
				cp.X = c.Size().Width - absPos.X - childSize.Width
			}
		}
		if absPos.Y+childSize.Height-theme.Padding() > c.Size().Height {
			cp.Y = c.Size().Height - absPos.Y - childSize.Height
		}
	}
	i.Child.Move(cp)
}

type menuItemRenderer struct {
	BaseRenderer
	i                *MenuItem
	icon             *canvas.Image
	lastThemePadding int
	minSize          fyne.Size
	text             *canvas.Text
}

func (r *menuItemRenderer) BackgroundColor() color.Color {
	if !fyne.CurrentDevice().IsMobile() && (r.i.hovered || (r.i.Child != nil && r.i.Child.Visible())) {
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
