//go:build wasm

package glfw

import (
	"regexp"
	"strings"
	"syscall/js"

	"fyne.io/fyne/v2"
)

var (
	isMobile bool
	isMacOS  bool

	setCursor func(name string)

	blurDummyEntry  func(...any) js.Value
	focusDummyEntry func(...any) js.Value
)

func init() {
	navigator := js.Global().Get("navigator")
	isMobile = regexp.MustCompile("Android|iPhone|iPad|iPod").
		MatchString(navigator.Get("userAgent").String())
	isMacOS = strings.Contains(navigator.Get("platform").String(), "Mac")

	document := js.Global().Get("document")
	style := document.Get("body").Get("style")
	setStyleProperty := style.Get("setProperty").Call("bind", style)
	setCursor = func(name string) {
		setStyleProperty.Invoke("cursor", name)
	}

	dummyEntry := document.Call("getElementById", "dummyEntry")
	if dummyEntry.IsNull() {
		return
	}

	blurDummyEntry = dummyEntry.Get("blur").Call("bind", dummyEntry).Invoke
	focusDummyEntry = dummyEntry.Get("focus").Call("bind", dummyEntry).Invoke
}

func (*glDevice) IsMobile() bool {
	return isMobile
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}

func (*glDevice) hideVirtualKeyboard() {
	if blurDummyEntry == nil {
		return
	}
	blurDummyEntry()
}

func (*glDevice) showVirtualKeyboard() {
	if focusDummyEntry == nil {
		return
	}
	focusDummyEntry()
}

func connectKeyboard(c *glCanvas) {
	c.OnFocus = handleKeyboard
	c.OnUnfocus = hideVirtualKeyboard
}

func isMacOSRuntime() bool {
	return isMacOS // Value depends on which OS the browser is running on.
}
