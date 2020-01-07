package test

import "fyne.io/fyne"

type device struct {
}

// Declare conformity with Device
var _ fyne.Device = (*device)(nil)

func (d *device) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (d *device) IsMobile() bool {
	return false
}

func (d *device) HasKeyboard() bool {
	return false
}

func (d *device) SystemScale() float32 {
	return 1
}
