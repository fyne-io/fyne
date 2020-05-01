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

// CreateRenderer satisfies the fyne.Widget interface.
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
	i.activateChild(&i.Parent.menuBase, i.updateChildPosition)
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
	if i.Child != nil {
		i.updateChildPosition()
	}
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
	i       *MenuItem
	icon    *canvas.Image
	minSize fyne.Size
	pad     fyne.Size
	msPad   fyne.Size
	tPad    int
	text    *canvas.Text
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) BackgroundColor() color.Color {
	if r.i.hovered || (r.i.Child != nil && r.i.Child.Visible()) {
		return theme.HoverColor()
	}

	return color.Transparent
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) Layout(size fyne.Size) {
	padding := r.padding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	if r.icon != nil {
		r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.icon.Move(fyne.NewPos(size.Width-theme.IconInlineSize(), (size.Height-theme.IconInlineSize())/2))
	}
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) MinSize() fyne.Size {
	if r.minSizeUnchanged() {
		return r.minSize
	}

	minSize := r.text.MinSize().Add(r.padding())
	r.msPad = r.padding()
	if r.icon != nil {
		minSize = minSize.Add(fyne.NewSize(theme.IconInlineSize(), 0))
	}
	r.minSize = minSize
	return r.minSize
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) Refresh() {
	canvas.Refresh(r.i)
}

func (r *menuItemRenderer) minSizeUnchanged() bool {
	return !r.minSize.IsZero() &&
		r.text.TextSize == theme.TextSize() &&
		(r.icon == nil || r.icon.Size().Width == theme.IconInlineSize()) &&
		r.msPad == r.padding()
}

func (r *menuItemRenderer) padding() fyne.Size {
	if r.tPad != theme.Padding() {
		r.pad = fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
	}
	return r.pad
}
