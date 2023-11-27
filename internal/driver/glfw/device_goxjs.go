//go:build js || wasm
// +build js wasm

package glfw

import (
	"syscall/js"

	"fyne.io/fyne/v2"
)

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	// Get the scale information from the web browser directly
	return float32(js.Global().Get("devicePixelRatio").Float())
}
