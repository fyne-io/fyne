package glfw

import (
	"fyne.io/fyne/v2"
)

func buildMenuOverlay(menus *fyne.MainMenu, w *window) fyne.CanvasObject {
	if len(menus.Items) == 0 {
		fyne.LogError("Main menu must have at least one child menu", nil)
		return nil
	}

	menus = addMissingQuitForMainMenu(menus, w)
	return NewMenuBar(menus, w.canvas)
}

func addMissingQuitForMainMenu(menus *fyne.MainMenu, w *window) *fyne.MainMenu {
	var lastItem *fyne.MenuItem
	if len(menus.Items[0].Items) > 0 {
		lastItem = menus.Items[0].Items[len(menus.Items[0].Items)-1]
		if lastItem.Label == "Quit" {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the first menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", nil)
		quitItem.IsQuit = true
		menus.Items[0].Items = append(menus.Items[0].Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menus.Items[0].Items {
		if item.IsQuit && item.Action == nil {
			item.Action = func() {
				for _, win := range w.driver.AllWindows() {
					if glWin, ok := win.(*window); ok {
						glWin.closed(glWin.view())
					} else {
						win.Close() // for test windows
					}
				}
			}
		}
	}
	return menus
}
