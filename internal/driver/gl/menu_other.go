// +build !darwin no_native_menus

package gl

import "fyne.io/fyne"

func hasNativeMenu() bool {
	return false
}

func setupNativeMenu(menu *fyne.MainMenu) {
	// no-op
}
