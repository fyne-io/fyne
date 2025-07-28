package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
)

func (dev *device) hideVirtualKeyboard() {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			return
		}

		d.app.HideVirtualKeyboard()
		dev.keyboardShown = false
	}
}

func (dev *device) handleKeyboard(obj fyne.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(fyne.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		if keyb, ok := obj.(mobile.Keyboardable); ok {
			dev.showVirtualKeyboard(keyb.Keyboard())
		} else {
			dev.showVirtualKeyboard(mobile.DefaultKeyboard)
		}
	} else {
		dev.hideVirtualKeyboard()
	}
}

func (dev *device) showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if d, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if d.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		dev.keyboardShown = true
		d.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}
