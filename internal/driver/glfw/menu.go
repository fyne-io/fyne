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
