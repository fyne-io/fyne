package embedded

import (
	"fyne.io/fyne/v2"
)

type noosDevice struct{}

func (n *noosDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (n *noosDevice) IsMobile() bool {
	return false
}

func (n *noosDevice) IsBrowser() bool {
	return false
}

func (n *noosDevice) HasKeyboard() bool {
	return true
}

func (n *noosDevice) SystemScaleForWindow(fyne.Window) float32 {
	return 1.0
}

func (n *noosDevice) Locale() fyne.Locale {
	return "en"
}
