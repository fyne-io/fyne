package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
	publicWidget "fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = (*menuBarItem)(nil)
var _ fyne.Focusable = (*menuBarItem)(nil)
var _ fyne.Widget = (*menuBarItem)(nil)

// menuBarItem is a widget for displaying an item for a fyne.Menu in a MenuBar.
type menuBarItem struct {
	widget.Base
	Menu   *fyne.Menu
	Parent *MenuBar

	active  bool
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
//
// Implements: fyne.Widget
func (i *menuBarItem) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.HoverColor())
	background.Hide()
	text := canvas.NewText(i.Menu.Label, theme.ForegroundColor())
	objects := []fyne.CanvasObject{background, text}

	return &menuBarItemRenderer{
		widget.NewBaseRenderer(objects),
		i,
		text,
		background,
	}
}

func (i *menuBarItem) FocusGained() {
	i.active = true
	if i.Parent.active {
		i.Parent.activateChild(i)
	}
	i.Refresh()
}

func (i *menuBarItem) FocusLost() {
	i.active = false
	i.Refresh()
}

func (i *menuBarItem) Focused() bool {
	return i.active
}

// MouseIn activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
//
// If the bar is not active, the item will be hovered.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseIn(_ *desktop.MouseEvent) {
	i.hovered = true
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
	i.Refresh()
}

// MouseMoved activates the item and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
// This might have an effect when mouse and keyboard control are mixed.
// Changing the active menu with the keyboard will make the hovered menu bar item inactive.
// On the next mouse move the hovered item is activated again.
//
// If the bar is not active, this will do nothing.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseMoved(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.Parent.canvas.Focus(i)
	}
}

// MouseOut does nothing if the bar is active.
//
// IF the bar is not active, it changes the item to not be hovered.
//
// Implements: desktop.Hoverable
func (i *menuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Tapped toggles the activation state of the menu bar.
// It shows the item’s menu if the bar is activated and hides it if the bar is deactivated.
//
// Implements: fyne.Tappable
func (i *menuBarItem) Tapped(*fyne.PointEvent) {
	i.Parent.toggle(i)
}

func (i *menuBarItem) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeyLeft:
		if !i.Child().DeactivateLastSubmenu() {
			i.Parent.canvas.FocusPrevious()
		}
	case fyne.KeyRight:
		if !i.Child().ActivateLastSubmenu() {
			i.Parent.canvas.FocusNext()
		}
	case fyne.KeyDown:
		i.Child().ActivateNext()
	case fyne.KeyUp:
		i.Child().ActivatePrevious()
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		i.Child().TriggerLast()
	}
}

func (i *menuBarItem) TypedRune(_ rune) {
}

type menuBarItemRenderer struct {
	widget.BaseRenderer
	i          *menuBarItem
	text       *canvas.Text
	background *canvas.Rectangle
}

func (r *menuBarItemRenderer) Layout(size fyne.Size) {
	padding := r.padding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.ForegroundColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))

	r.background.Resize(size)
}

func (r *menuBarItemRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuBarItemRenderer) Refresh() {
	if r.i.active && r.i.Parent.active {
		r.background.FillColor = theme.FocusColor()
		r.background.Show()
	} else if r.i.hovered && !r.i.Parent.active {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	canvas.Refresh(r.i)
}

func (r *menuBarItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
