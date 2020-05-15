package widget

import (
	"fyne.io/fyne"
)

// PopUpMenu is a Menu which displays itself in an OverlayContainer.
type PopUpMenu struct {
	*Menu
	canvas  fyne.Canvas
	overlay *OverlayContainer
}

// NewPopUpMenu creates a new PopUpMenu.
func NewPopUpMenu(m *fyne.Menu, c fyne.Canvas) *PopUpMenu {
	p := &PopUpMenu{Menu: NewMenu(m), canvas: c}
	p.Menu.Resize(p.Menu.MinSize())
	p.Menu.customSized = true
	o := NewOverlayContainer(p.Menu, c, p.dismiss)
	o.Resize(o.MinSize())
	p.overlay = o
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
	if p.pos.X+p.size.Width > p.canvas.Size().Width {
		p.pos.X = p.canvas.Size().Width - p.size.Width
		if p.pos.X < 0 {
			p.pos.X = 0 // TODO here we may need a scroller as it's wider than our canvas
		}
	}
	if p.pos.Y+p.size.Height > p.canvas.Size().Height {
		p.pos.Y = p.canvas.Size().Height - p.size.Height
		if p.pos.Y < 0 {
			p.pos.Y = 0 // TODO here we may need a scroller as it's longer than our canvas
		}
	}
}
