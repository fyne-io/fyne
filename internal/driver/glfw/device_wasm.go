//go:build wasm

package glfw

import (
	"regexp"
	"strings"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var (
	isMobile = regexp.MustCompile("Android|iPhone|iPad|iPod").
			MatchString(js.Global().Get("navigator").Get("userAgent").String())
	isMacOS = strings.Contains(js.Global().Get("window").Get("navigator").Get("platform").String(), "Mac")
)

func (*glDevice) IsMobile() bool {
	return isMobile
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}

func isMacOSRuntime() bool {
	return isMacOS // Value depends on which OS the browser is running on.
}
