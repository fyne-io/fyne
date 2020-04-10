package glfw

import (
	"runtime"

	"fyne.io/fyne"
)

type glDevice struct {
}

// Declare conformity with Device
var _ fyne.Device = (*glDevice)(nil)

func (*glDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationHorizontalLeft // TODO should we consider the monitor orientation or topmost window?
}

func (*glDevice) HasKeyboard() bool {
	return true // TODO actually check - we could be in tablet mode
}

// Deprecated: SystemScaleForWindow provides an accurate answer, this method may return incorrect results!
func (d *glDevice) SystemScale() float32 {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return 1.0
	}

	return fyne.SettingsScaleAuto
}
