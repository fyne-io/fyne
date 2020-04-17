package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/widget"
)

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	m := widget.NewMenu(menu)
	pop := newPopUp(m, c)
	pop.NotPadded = true
	focused := c.Focused()
	m.DismissAction = func() {
		if c.Focused() == nil {
			c.Focus(focused)
		}
		pop.Hide()
	}
	pop.ShowAtPosition(pos)
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}
