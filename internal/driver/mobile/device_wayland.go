// +build wayland

package gomobile

import "fyne.io/fyne/v2"

func (*device) SystemScaleForWindow(_ fyne.Window) float32 {
	return 1 // PinePhone simplification, our only wayland mobile currently
}
