package gomobile

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/mobile"
	"github.com/fyne-io/mobile/app"
)

func showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if driver, ok := fyne.CurrentApp().Driver().(*mobileDriver); ok {
		driver.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}

func hideVirtualKeyboard() {
	if driver, ok := fyne.CurrentApp().Driver().(*mobileDriver); ok {
		driver.app.HideVirtualKeyboard()
	}
}
