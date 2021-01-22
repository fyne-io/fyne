package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/widget"
)

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
	p := &PopUpMenu{Menu: NewMenu(menu), canvas: c}
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := widget.NewOverlayContainer(p.Menu, c, p.Dismiss)
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

// CreateRenderer returns a new renderer for the pop-up menu.
//
// Implements: fyne.Widget
func (p *PopUpMenu) CreateRenderer() fyne.WidgetRenderer {
	return p.overlay.CreateRenderer()
}

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
	widget.MoveWidget(&p.Base, p, p.adjustedPosition(pos, p.Size()))
}

// Resize changes the size of the pop-up menu.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Resize(size fyne.Size) {
	widget.MoveWidget(&p.Base, p, p.adjustedPosition(p.Position(), size))
	p.Menu.Resize(size)
}

// Show makes the pop-up menu visible.
//
// Implements: fyne.Widget
func (p *PopUpMenu) Show() {
	p.overlay.Show()
	p.Menu.Show()
}

// ShowAtPosition shows the pop-up menu at the specified position.
func (p *PopUpMenu) ShowAtPosition(pos fyne.Position) {
	p.Move(pos)
	p.Show()
}

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
