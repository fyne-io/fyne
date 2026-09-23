package test

import (
	"runtime"

	"fyne.io/fyne/v2"
)

type device struct{}

// Declare conformity with Device
var _ fyne.Device = (*device)(nil)

func (*device) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (*device) HasKeyboard() bool {
	return false
}

func (d *device) SystemScale() float32 {
	return d.SystemScaleForWindow(nil)
}

func (*device) SystemScaleForWindow(fyne.Window) float32 {
	return 1
}

func (*device) Locale() fyne.Locale {
	return "en"
}

func (*device) IsBrowser() bool {
	return runtime.GOARCH == "js" || runtime.GOOS == "js"
}
