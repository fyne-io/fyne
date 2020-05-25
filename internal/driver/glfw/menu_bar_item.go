package glfw

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
	publicWidget "fyne.io/fyne/widget"
)

var _ fyne.Widget = (*menuBarItem)(nil)
var _ desktop.Hoverable = (*menuBarItem)(nil)

// menuBarItem is a widget for displaying an item for a fyne.Menu in a MenuBar.
type menuBarItem struct {
	widget.Base
	Menu   *fyne.Menu
	Parent *MenuBar

	child   *publicWidget.Menu
	hovered bool
}

func (i *menuBarItem) Child() *publicWidget.Menu {
	if i.child == nil {
		child := publicWidget.NewMenu(i.Menu)
		child.Hide()
		child.OnDismiss = i.Parent.deactivate
		i.child = child
	}
	return i.child
}

// CreateRenderer returns a new renderer for the menu bar item.
// Implements: fyne.Widget
func (i *menuBarItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Menu.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}

	return &menuBarItemRenderer{
		widget.NewBaseRenderer(objects),
		i,
		text,
	}
}

// Hide hides the menu bar item.
// Implements: fyne.Widget
func (i *menuBarItem) Hide() {
	widget.HideWidget(&i.Base, i)
}

// MinSize returns the minimal size of the menu bar item.
// Implements: fyne.Widget
func (i *menuBarItem) MinSize() fyne.Size {
	return widget.MinSizeOf(i)
}

// MouseIn changes the item to be hovered and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.hovered = true
		i.Parent.activateChild(i)
		i.Refresh()
	} else {
		i.hovered = true
		i.Refresh()
	}
}

// MouseMoved does nothing.
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut changes the item to not be hovered but has no effect on the visibility of the menu.
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Move sets the position of the widget relative to its parent.
// Implements: fyne.Widget
func (i *menuBarItem) Move(pos fyne.Position) {
	widget.MoveWidget(&i.Base, i, pos)
}

// Refresh triggers a redraw of the menu bar item.
// Implements: fyne.Widget
func (i *menuBarItem) Refresh() {
	widget.RefreshWidget(i)
}

// Resize changes the size of the menu bar item.
// Implements: fyne.Widget
func (i *menuBarItem) Resize(size fyne.Size) {
	widget.ResizeWidget(&i.Base, i, size)
}

// Show makes the menu bar item visible.
// Implements: fyne.Widget
func (i *menuBarItem) Show() {
	widget.ShowWidget(&i.Base, i)
}

// Tapped toggles the activation state of the menu bar.
// It shows the item’s menu if the bar is activated and hides it if the bar is deactivated.
// Implements: fyne.Tappable
func (i *menuBarItem) Tapped(*fyne.PointEvent) {
	if i.Parent.active {
		i.Parent.deactivate()
	} else {
		i.Parent.activateChild(i)
	}
}

type menuBarItemRenderer struct {
	widget.BaseRenderer
	i    *menuBarItem
	text *canvas.Text
}

func (r *menuBarItemRenderer) BackgroundColor() color.Color {
	if r.i.hovered || (r.i.child != nil && r.i.child.Visible()) {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuBarItemRenderer) Layout(_ fyne.Size) {
	padding := r.padding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuBarItemRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuBarItemRenderer) Refresh() {
	canvas.Refresh(r.i)
}

func (r *menuBarItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
