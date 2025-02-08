//go:build wasm

package glfw

import (
	"regexp"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var isMobile = regexp.MustCompile("Android|BlackBerry|iPhone|iPad|iPod|Opera Mini|IEMobile").
	MatchString(js.Global().Get("navigator").Get("userAgent").String())

func (*glDevice) IsMobile() bool {
	return isMobile
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}
