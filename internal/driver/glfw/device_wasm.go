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

var dummyEntry = js.Global().Get("document").Call("getElementById", "dummyEntry")

func (*glDevice) IsMobile() bool {
	return isMobile
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}

func (*glDevice) hideVirtualKeyboard() {
	if dummyEntry.IsNull() {
		return
	}
	dummyEntry.Call("blur")
}

func (*glDevice) showVirtualKeyboard() {
	if dummyEntry.IsNull() {
		return
	}
	dummyEntry.Call("focus")
}

func connectKeyboard(c *glCanvas) {
	c.OnFocus = handleKeyboard
	c.OnUnfocus = hideVirtualKeyboard
}

func isMacOSRuntime() bool {
	return isMacOS // Value depends on which OS the browser is running on.
}
