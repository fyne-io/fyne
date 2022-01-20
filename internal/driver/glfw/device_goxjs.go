//go:build js || wasm
// +build js wasm

package glfw

import (
	"regexp"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var AndroidId = regexp.MustCompile(`Android`)
var BlackBerryID = regexp.MustCompile(`BlackBerry`)
var iOSID = regexp.MustCompile(`iPhone|iPad|iPod`)
var OperaID = regexp.MustCompile(`Opera Mini`)
var WindowsMobileID = regexp.MustCompile(`IEMobile`)

var navigator = js.Global().Get("navigator")
var userAgent = navigator.Get("userAgent").String()
var mobileCheck = AndroidId.MatchString(userAgent) ||
	BlackBerryID.MatchString(userAgent) ||
	iOSID.MatchString(userAgent) ||
	OperaID.MatchString(userAgent) ||
	WindowsMobileID.MatchString(userAgent)

func (*glDevice) IsMobile() bool {
	return mobileCheck
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}
