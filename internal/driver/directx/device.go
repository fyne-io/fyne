//go:build windows

package directx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

// Declare conformity with Device.
var _ fyne.Device = (*dxDevice)(nil)

type dxDevice struct{}

func (*dxDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationHorizontalLeft
}

func (*dxDevice) IsMobile() bool  { return false }
func (*dxDevice) IsBrowser() bool { return false }

func (*dxDevice) HasKeyboard() bool { return true }

func (*dxDevice) SystemScaleForWindow(w fyne.Window) float32 {
	if win, ok := w.(*window); ok && win.hwnd != 0 {
		return win.calculatedScale()
	}
	return 1.0
}

// Locale reports the user's language and region from the OS settings.
func (*dxDevice) Locale() fyne.Locale {
	return lang.SystemLocale()
}
