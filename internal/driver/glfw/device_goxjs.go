//go:build js || wasm
// +build js wasm

package glfw

import (
	"regexp"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var isMobile = regexp.MustCompile("Android|BlackBerry|iPhone|iPad|iPod|Opera Mini|IEMobile")

var navigator = js.Global().Get("navigator")
var userAgent = navigator.Get("userAgent").String()
var mobileCheck = isMobile.MatchString(userAgent)

func (*glDevice) IsMobile() bool {
	return mobileCheck
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}
