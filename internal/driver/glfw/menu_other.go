// +build !darwin no_native_menus

package glfw

import "fyne.io/fyne/v2"

func hasNativeMenu() bool {
	return false
}

func setupNativeMenu(_ *window, _ *fyne.MainMenu) {
	// no-op
}
