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
	o := NewOverlayContainer(p.Menu, c, p.dismiss)
	o.Resize(o.MinSize())
	p.overlay = o
	return p
}

// CreateRenderer satisfies the fyne.Widget interface.
func (p *PopUpMenu) CreateRenderer() fyne.WidgetRenderer {
	return p.overlay.CreateRenderer()
}

// Hide satisfies the fyne.Widget interface.
func (p *PopUpMenu) Hide() {
	p.overlay.Hide()
	p.Menu.Hide()
}

// Move satisfies the fyne.Widget interface.
func (p *PopUpMenu) Move(pos fyne.Position) {
	p.Menu.Move(pos)
	p.adjustPosition()
}

// Resize satisfies the fyne.Widget interface.
func (p *PopUpMenu) Resize(size fyne.Size) {
	p.Menu.Resize(size)
	p.adjustPosition()
}

// Show satisfies the fyne.Widget interface.
func (p *PopUpMenu) Show() {
	p.overlay.Show()
	p.Menu.Show()
}

// ShowAtPosition shows the menu at the specified position.
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
