package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

func hideVirtualKeyboard() {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			return
		}

		d.app.HideVirtualKeyboard()
	}
}

func showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		d.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}
