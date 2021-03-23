package glfw

import (
	"runtime"

	"fyne.io/fyne/v2"
)

type glDevice struct {
}

// Declare conformity with Device
var _ fyne.Device = (*glDevice)(nil)

func (*glDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationHorizontalLeft // TODO should we consider the monitor orientation or topmost window?
}

func (*glDevice) IsMobile() bool {
	return false
}

func (*glDevice) HasKeyboard() bool {
	return true // TODO actually check - we could be in tablet mode
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	if runtime.GOOS == "darwin" {
		return 1.0 // macOS scaling is done at the texture level
	}
	if runtime.GOOS == "windows" {
		xScale, _ := w.(*window).viewport.GetContentScale()
		return xScale
	}

	return scaleAuto
}
