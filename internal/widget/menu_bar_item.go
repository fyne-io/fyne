package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*MenuBarItem)(nil)
var _ desktop.Hoverable = (*MenuBarItem)(nil)

// MenuBarItem is a widget for displaying an item for a fyne.Menu in a MenuBar.
type MenuBarItem struct {
	base
	menuItemBase
	Menu   *fyne.Menu
	Parent *MenuBar

	hovered bool
}

// CreateRenderer satisfies the fyne.Widget interface.
func (i *MenuBarItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Menu.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	i.initChildWidget(i.Menu, i.Parent.deactivate)
	objects = append(objects, i.Child)

	return &menuBarItemRenderer{
		NewBaseRenderer(objects),
		i,
		text,
	}
}

// Hide satisfies the fyne.Widget interface.
func (i *MenuBarItem) Hide() {
	i.hide(i)
}

// MinSize satisfies the fyne.Widget interface.
func (i *MenuBarItem) MinSize() fyne.Size {
	return i.minSize(i)
}

// MouseIn satisfies the desktop.Hoverable interface.
func (i *MenuBarItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.hovered = true
		i.activateChild(&i.Parent.menuBase, i.updateChildPosition)
		i.Refresh()
	} else {
		i.hovered = true
		i.Refresh()
	}
}

// MouseMoved satisfies the desktop.Hoverable interface.
func (i *MenuBarItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut satisfies the desktop.Hoverable interface.
func (i *MenuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Refresh satisfies the fyne.Widget interface.
func (i *MenuBarItem) Refresh() {
	i.refresh(i)
}

// Resize satisfies the fyne.Widget interface.
func (i *MenuBarItem) Resize(size fyne.Size) {
	i.resize(size, i)
}

// Show satisfies the fyne.Widget interface.
func (i *MenuBarItem) Show() {
	i.show(i)
}

// Tapped satisfies the fyne.Tappable interface.
func (i *MenuBarItem) Tapped(*fyne.PointEvent) {
	if i.Parent.active {
		i.Parent.deactivate()
	} else {
		i.Parent.activate()
		i.activateChild(&i.Parent.menuBase, i.updateChildPosition)
	}
	i.Refresh()
}

func (i *MenuBarItem) updateChildPosition() {
	i.Child.Move(fyne.NewPos(0, i.Size().Height))
}

type menuBarItemRenderer struct {
	BaseRenderer
	i    *MenuBarItem
	text *canvas.Text
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *menuBarItemRenderer) BackgroundColor() color.Color {
	if r.i.hovered || (r.i.Child != nil && r.i.Child.Visible()) {
		return theme.HoverColor()
	}

	return color.Transparent
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *menuBarItemRenderer) Layout(_ fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *menuBarItemRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *menuBarItemRenderer) Refresh() {
	if r.text.TextSize != theme.TextSize() {
		defer r.Layout(r.i.Size())
	}
	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	canvas.Refresh(r.text)
}

func (r *menuBarItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
