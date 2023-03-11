package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/widget"
)

var _ fyne.Widget = (*PopUpMenu)(nil)
var _ fyne.Focusable = (*PopUpMenu)(nil)

// PopUpMenu is a Menu which displays itself in an OverlayContainer.
type PopUpMenu struct {
	*Menu
	canvas  fyne.Canvas
	overlay *widget.OverlayContainer
}

// NewPopUpMenu creates a new, reusable popup menu. You can show it using ShowAtPosition.
//
// Since: 2.0
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUpMenu {
	m := &Menu{}
	m.setMenu(menu)
	p := &PopUpMenu{Menu: m, canvas: c}
	p.ExtendBaseWidget(p)
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := widget.NewOverlayContainer(p, c, p.Dismiss)
	o.Resize(o.MinSize())
	p.overlay = o
	p.OnDismiss = func() {
		p.Hide()
	}
	return p
}

// ShowPopUpMenuAtPosition creates a PopUp menu populated with items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func ShowPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) {
	m := NewPopUpMenu(menu, c)
	m.ShowAtPosition(pos)
}

// FocusGained is triggered when the object gained focus. For the pop-up menu it does nothing.
//
// Implements: fyne.Focusable
func (p *PopUpMenu) FocusGained() {}

// FocusLost is triggered when the object lost focus. For the pop-up menu it does nothing.
//
// Implements: fyne.Focusable
func (p *PopUpMenu) FocusLost() {}

// Hide hides the pop-up menu.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Hide() {
	p.overlay.Hide()
	p.Menu.Hide()
}

// Move moves the pop-up menu.
// The position is absolute because pop-up menus are shown in an overlay which covers the whole canvas.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Move(pos fyne.Position) {
	p.BaseWidget.Move(p.adjustedPosition(pos, p.Size()))
}

// Resize changes the size of the pop-up menu.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Resize(size fyne.Size) {
	p.BaseWidget.Move(p.adjustedPosition(p.Position(), size))
	p.Menu.Resize(size)
}

// Show makes the pop-up menu visible.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Show() {
	p.Menu.alignment = p.alignment
	p.Menu.Refresh()

	p.overlay.Show()
	p.Menu.Show()
	if !fyne.CurrentDevice().IsMobile() {
		p.canvas.Focus(p)
	}
}

// ShowAtPosition shows the pop-up menu at the specified position.
func (p *PopUpMenu) ShowAtPosition(pos fyne.Position) {
	p.Move(pos)
	p.Show()
}

// TypedKey handles key events. It allows keyboard control of the pop-up menu.
//
// Implements: fyne.Focusable
func (p *PopUpMenu) TypedKey(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyDown:
		p.ActivateNext()
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		p.TriggerLast()
	case fyne.KeyEscape:
		p.Dismiss()
	case fyne.KeyLeft:
		p.DeactivateLastSubmenu()
	case fyne.KeyRight:
		p.ActivateLastSubmenu()
	case fyne.KeyUp:
		p.ActivatePrevious()
	}
}

// TypedRune handles text events. For pop-up menus this does nothing.
//
// Implements: fyne.Focusable
func (p *PopUpMenu) TypedRune(rune) {}

func (p *PopUpMenu) adjustedPosition(pos fyne.Position, size fyne.Size) fyne.Position {
	x := pos.X
	y := pos.Y
	if x+size.Width > p.canvas.Size().Width {
		x = p.canvas.Size().Width - size.Width
		if x < 0 {
			x = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if y+size.Height > p.canvas.Size().Height {
		y = p.canvas.Size().Height - size.Height
		if y < 0 {
			y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	return fyne.NewPos(x, y)
}
