package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/widget"
)

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	options := NewVBox()
	for _, item := range menu.Items {
		if item.IsSeparator {
			options.Append(widget.NewMenuItemSeparator())
		} else {
			options.Append(widget.NewMenuItem(item))
		}
	}
	pop := newPopUp(options, c)
	pop.NotPadded = true
	focused := c.Focused()
	for _, o := range options.Children {
		if item, ok := o.(*widget.MenuItem); ok {
			item.DismissAction = func() {
				if c.Focused() == nil {
					c.Focus(focused)
				}
				pop.Hide()
			}
		}
	}
	pop.ShowAtPosition(pos)
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}
