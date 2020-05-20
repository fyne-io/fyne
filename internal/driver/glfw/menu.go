package glfw

import (
	"fyne.io/fyne"
)

func buildMenuOverlay(menus *fyne.MainMenu, c fyne.Canvas) fyne.CanvasObject {
	if menus.Items[0].Items[len(menus.Items[0].Items)-1].Label != "Quit" { // make sure the first menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", func() {
			fyne.CurrentApp().Quit()
		})
		menus.Items[0].Items = append(menus.Items[0].Items, fyne.NewMenuItemSeparator(), quitItem)
	}

	return NewMenuBar(menus, c)
}
