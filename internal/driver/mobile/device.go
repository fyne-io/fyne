package mobile

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
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
var (
	_ fyne.Device   = (*device)(nil)
	_ mobile.Device = (*device)(nil)
)

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
	d.ShowVirtualKeyboardType(mobile.DefaultKeyboard)
}

func (d *device) ShowVirtualKeyboardType(keyboard mobile.KeyboardType) {
	if drv, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if drv.app == nil { // not yet running
			fyne.LogError("Cannot show keyboard before app is running", nil)
			return
		}

		d.keyboardShown = true
		drv.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}

func (d *device) HideVirtualKeyboard() {
	if drv, ok := fyne.CurrentApp().Driver().(*driver); ok {
		if drv.app == nil { // not yet running
			return
		}

		drv.app.HideVirtualKeyboard()
		d.keyboardShown = false
	}
}

// keyboardText returns the text and rune offset of a focused text input.
func keyboardText(obj fyne.Focusable) (string, int, bool) {
	entry, ok := obj.(*widget.Entry)
	if !ok {
		return "", 0, false
	}
	return entry.Text, entry.CursorTextOffset(), true
}

// syncKeyboardText mirrors the focused entry into the native keyboard proxy.
// While a soft backspace repeats, Java is ahead of Go by a few events, so this
// must not be called for those events - a setText there would end the repeat.
func (d *device) syncKeyboardText(obj fyne.Focusable) {
	if !d.keyboardShown {
		return
	}

	text, cursor, ok := keyboardText(obj)
	if !ok {
		return
	}

	drv, ok := fyne.CurrentApp().Driver().(*driver)
	if !ok || drv.app == nil {
		return
	}
	drv.app.SetKeyboardText(text, cursor)
}

func (d *device) handleKeyboard(obj fyne.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(fyne.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		if keyb, ok := obj.(mobile.Keyboardable); ok {
			d.ShowVirtualKeyboardType(keyb.Keyboard())
		} else {
			d.ShowVirtualKeyboardType(mobile.DefaultKeyboard)
		}
		d.syncKeyboardText(obj)
	} else {
		d.HideVirtualKeyboard()
	}
}
