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
	Child         *Menu
	DismissAction func()
	Item          *fyne.MenuItem
	Parent        *Menu

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

// CreateRenderer satisfies the fyne.Widget interface.
func (i *MenuItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Item.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	var icon *canvas.Image
	if i.Item.ChildMenu != nil {
		icon = canvas.NewImageFromResource(theme.MenuExpandIcon())
		objects = append(objects, icon)
		i.initChildWidget()
		objects = append(objects, i.Child)
	}
	return &menuItemRenderer{NewBaseRenderer(objects), i, icon, text}
}

// Hide satisfies the fyne.Widget interface.
func (i *MenuItem) Hide() {
	i.hide(i)
}

// MinSize satisfies the fyne.Widget interface.
func (i *MenuItem) MinSize() fyne.Size {
	return i.minSize(i)
}

// MouseIn satisfies the desktop.Hoverable interface.
func (i *MenuItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.activateChild()
	i.Refresh()
}

// MouseMoved satisfies the desktop.Hoverable interface.
func (i *MenuItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut satisfies the desktop.Hoverable interface.
func (i *MenuItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Refresh satisfies the fyne.Widget interface.
func (i *MenuItem) Refresh() {
	i.refresh(i)
}

// Resize satisfies the fyne.Widget interface.
func (i *MenuItem) Resize(size fyne.Size) {
	i.resize(size, i)
}

// Show satisfies the fyne.Widget interface.
func (i *MenuItem) Show() {
	i.show(i)
}

// Tapped satisfies the fyne.Tappable interface.
func (i *MenuItem) Tapped(*fyne.PointEvent) {
	if i.Item.Action == nil {
		return
	}

	i.Item.Action()
	i.Parent.dismiss()
}

func (i *MenuItem) activateChild() {
	if i.Child != nil && i.Child.activeChild != nil {
		i.Child.activeChild.Hide()
		i.Child.activeChild = nil
	}

	if i.Parent.activeChild == i.Child {
		return
	}

	if i.Parent.activeChild != nil {
		i.Parent.activeChild.Hide()
	}
	i.Parent.activeChild = i.Child
	if i.Child != nil {
		i.Child.Show()
	}
}

func (i *MenuItem) initChildWidget() {
	if i.Child != nil {
		return
	}
	i.Child = NewMenu(i.Item.ChildMenu)
	i.Child.Hide()
	i.Child.DismissAction = func() { i.Parent.dismiss() }
}

type menuItemRenderer struct {
	BaseRenderer
	i    *MenuItem
	icon *canvas.Image
	text *canvas.Text
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) BackgroundColor() color.Color {
	if r.i.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}

	if r.i.Child != nil {
		r.i.Child.Resize(r.i.Child.MinSize())
		var cp fyne.Position
		itemSize := r.i.Size()
		cp.X = itemSize.Width

		d := fyne.CurrentApp().Driver()
		c := d.CanvasForObject(r.i)
		if c != nil {
			absPos := d.AbsolutePositionForObject(r.i)
			childSize := r.i.Child.Size()
			if absPos.X+itemSize.Width+childSize.Width > c.Size().Width {
				if absPos.X-childSize.Width >= 0 {
					cp.X = -childSize.Width
				} else {
					cp.X = c.Size().Width - absPos.X - childSize.Width
				}
			}
			if absPos.Y+childSize.Height > c.Size().Height {
				cp.Y = c.Size().Height - absPos.Y - childSize.Height
			}
		}
		r.i.Child.Move(cp)
	}
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) MinSize() fyne.Size {
	minSize := r.text.MinSize().Add(r.padding())
	if r.icon == nil {
		return minSize
	}
	return minSize.Add(fyne.NewSize(theme.IconInlineSize(), 0))
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) Refresh() {
	if r.text.TextSize != theme.TextSize() {
		defer r.Layout(r.i.Size())
	}
	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	canvas.Refresh(r.text)
}

func (r *menuItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
