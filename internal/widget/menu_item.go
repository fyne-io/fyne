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
	return &menuItemRenderer{NewBaseRenderer([]fyne.CanvasObject{text}), i, text}
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

type menuItemRenderer struct {
	BaseRenderer
	i    *MenuItem
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
func (r *menuItemRenderer) Layout(fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuItemRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
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
