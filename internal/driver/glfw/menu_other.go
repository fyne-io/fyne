// +build !darwin no_native_menus

package glfw

import "fyne.io/fyne"

func hasNativeMenu() bool {
	return false
}

func setupNativeMenu(menu *fyne.MainMenu) {
	// no-op
}
