package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/widget"
)

// PopUpMenu is a Menu which displays itself in an OverlayContainer.
type PopUpMenu struct {
	*widget.Menu
	canvas  fyne.Canvas
	overlay *widget.OverlayContainer
}

// ShowPopUpMenuAtPosition creates a PopUp menu populated with items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func ShowPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) {
	m := NewPopUpMenu2(menu, c)
	m.ShowAtPosition(pos)
}

func NewPopUpMenu2(menu *fyne.Menu, c fyne.Canvas) *PopUpMenu {
	p := &PopUpMenu{Menu: widget.NewMenu(menu), canvas: c}
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := widget.NewOverlayContainer(p.Menu, c, p.Dismiss)
	o.Resize(o.MinSize())
	p.overlay = o

	focused := c.Focused()
	p.DismissAction = func() {
		if c.Focused() == nil {
			c.Focus(focused)
		}
		p.Hide()
	}
	return p
}

// CreateRenderer returns a new renderer for the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) CreateRenderer() fyne.WidgetRenderer {
	return p.overlay.CreateRenderer()
}

// Hide hides the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) Hide() {
	p.overlay.Hide()
	p.Menu.Hide()
}

// Move moves the pop-up menu.
// The position is absolute because pop-up menus are shown in an overlay which covers the whole canvas.
// Implements: fyne.Widget
func (p *PopUpMenu) Move(pos fyne.Position) {
	p.Menu.Move(pos)
	p.adjustPosition()
}

// Resize changes the size of the pop-up menu.
// Implements: fyne.Widget
func (p *PopUpMenu) Resize(size fyne.Size) {
	p.Menu.Resize(size)
	p.adjustPosition()
}

// Show makes the pop-up menu visible.
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

func (p *PopUpMenu) adjustPosition() {
	newX := p.Position().X
	newY := p.Position().Y
	if p.Position().X+p.Size().Width > p.canvas.Size().Width {
		newX = p.canvas.Size().Width - p.Size().Width
		if newX < 0 {
			newX = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if p.Position().Y+p.Size().Height > p.canvas.Size().Height {
		newY = p.canvas.Size().Height - p.Size().Height
		if newY < 0 {
			newY = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
	if newX != p.Position().X || newY != p.Position().Y {
		p.Base.Move(fyne.NewPos(newX, newY))
	}
}
