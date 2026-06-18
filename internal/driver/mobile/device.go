package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/lang"
)

type device struct {
	safeTop, safeLeft, safeBottom, safeRight int

	keyboardShown bool
}

//lint:file-ignore U1000 Var currentDPI is used in other files, but not here
var (
	currentOrientation size.Orientation
	currentDPI         float32
)

// Declare conformity with Device
var _ fyne.Device = (*device)(nil)

func (*device) Locale() fyne.Locale {
	return lang.SystemLocale()
}

func (*device) Orientation() fyne.DeviceOrientation {
	switch currentOrientation {
	case size.OrientationLandscape:
		return fyne.OrientationHorizontalLeft
	default:
		return fyne.OrientationVertical
	}
}

func (*device) IsMobile() bool {
	return true
}

func (*device) IsBrowser() bool {
	return false
}

func (*device) HasKeyboard() bool {
	return false
}

func (d *device) ShowVirtualKeyboard() {
	d.showVirtualKeyboard(mobile.DefaultKeyboard)
}

func (d *device) ShowVirtualKeyboardType(keyboard mobile.KeyboardType) {
	d.showVirtualKeyboard(keyboard)
}

func (d *device) HideVirtualKeyboard() {
	d.hideVirtualKeyboard()
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

func (d *device) hideVirtualKeyboard() {
	if drv, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if drv.app == nil { // not yet running
			return
		}

		drv.app.HideVirtualKeyboard()
		d.keyboardShown = false
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
