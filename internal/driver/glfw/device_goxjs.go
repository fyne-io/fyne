//go:build js || wasm
// +build js wasm

package glfw

import (
	"regexp"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var (
	AndroidId       = regexp.MustCompile(`Android`)
	BlackBerryID    = regexp.MustCompile(`BlackBerry`)
	iOSID           = regexp.MustCompile(`iPhone|iPad|iPod`)
	OperaID         = regexp.MustCompile(`Opera Mini`)
	WindowsMobileID = regexp.MustCompile(`IEMobile`)
)

var (
	navigator   = js.Global().Get("navigator")
	userAgent   = navigator.Get("userAgent").String()
	mobileCheck = AndroidId.MatchString(userAgent) ||
		BlackBerryID.MatchString(userAgent) ||
		iOSID.MatchString(userAgent) ||
		OperaID.MatchString(userAgent) ||
		WindowsMobileID.MatchString(userAgent)
)

func (*glDevice) IsMobile() bool {
	return mobileCheck
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}
