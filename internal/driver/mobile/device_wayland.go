//go:build wayland
// +build wayland

package mobile

import "fyne.io/fyne/v2"

const tapYOffset = -4.0 // to compensate for how we hold our fingers on the device

func (*device) SystemScaleForWindow(_ fyne.Window) float32 {
	return 1 // PinePhone simplification, our only wayland mobile currently
}
