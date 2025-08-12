package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

func (d *device) hideVirtualKeyboard() {
	if drv, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if drv.app == nil { // not yet running
			return
		}

		drv.app.HideVirtualKeyboard()
		d.keyboardShown = false
	}
}

func (d *device) handleKeyboard(obj fyne.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(fyne.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		if keyb, ok := obj.(mobile.Keyboardable); ok {
			d.showVirtualKeyboard(keyb.Keyboard())
		} else {
			d.showVirtualKeyboard(mobile.DefaultKeyboard)
		}
	} else {
		d.hideVirtualKeyboard()
	}
}

func (d *device) showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if drv, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if drv.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		d.keyboardShown = true
		drv.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}
