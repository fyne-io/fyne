package gl

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type menuBarAction struct {
	Label string

	menu   *fyne.Menu
	canvas fyne.Canvas
}

// ToolbarObject gets a button to render this ToolbarAction
func (m *menuBarAction) ToolbarObject() fyne.CanvasObject {
	button := widget.NewButton(m.Label, nil)

	button.OnTapped = func() {
		pos := button.Position().Add(fyne.NewPos(0, button.Size().Height))
		showMenu(m.menu, pos, m.canvas)
	}

	return button
}

func newMenuBarAction(menu *fyne.Menu, w fyne.Window) widget.ToolbarItem {
	return &menuBarAction{menu.Label, menu, w.Canvas()}
}

func newMenuBarActionWithQuit(menu *fyne.Menu, w fyne.Window) widget.ToolbarItem {
	if menu.Items[len(menu.Items)-1].Label != "Quit" { // make sure the first menu always has a quit option
		// TODO append a separator as well
		menu.Items = append(menu.Items, fyne.NewMenuItem("Quit", func() {
			fyne.CurrentApp().Quit()
		}))
	}
	return &menuBarAction{menu.Label, menu, w.Canvas()}
}

func buildMenuBar(menus *fyne.MainMenu, w fyne.Window) *widget.Toolbar {
	var items []widget.ToolbarItem

	for i, menu := range menus.Items {
		if i == 0 {
			items = append(items, newMenuBarActionWithQuit(menu, w))
		} else {
			items = append(items, newMenuBarAction(menu, w))
		}
	}
	return widget.NewToolbar(items...)
}

func showMenu(menu *fyne.Menu, pos fyne.Position, c fyne.Canvas) {
	pop := widget.NewPopUpMenu(menu, c)
	pop.Move(pos)
}
