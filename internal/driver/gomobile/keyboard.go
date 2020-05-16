package gomobile

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/mobile"
	"github.com/fyne-io/mobile/app"
)

func showVirtualKeyboard(keyboard mobile.Keyboard) {
	fyne.CurrentApp().Driver().(*mobileDriver).app.ShowVirtualKeyboard(app.Keyboard(keyboard))
}

func hideVirtualKeyboard() {
	fyne.CurrentApp().Driver().(*mobileDriver).app.HideVirtualKeyboard()
}
