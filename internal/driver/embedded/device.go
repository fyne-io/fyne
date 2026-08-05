package embedded

import (
	"fyne.io/fyne/v2"
)

type noosDevice struct{}

func (*noosDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (*noosDevice) IsMobile() bool {
	return false
}

func (*noosDevice) IsBrowser() bool {
	return false
}

func (*noosDevice) HasKeyboard() bool {
	return true
}

func (*noosDevice) SystemScaleForWindow(fyne.Window) float32 {
	return 1.0
}

func (*noosDevice) Locale() fyne.Locale {
	return "en"
}
