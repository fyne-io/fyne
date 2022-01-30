//go:build !js && !wasm
// +build !js,!wasm

package glfw

import (
	"runtime"

	"fyne.io/fyne/v2"
)

func (*glDevice) IsMobile() bool {
	return false
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	if runtime.GOOS == "darwin" {
		return 1.0 // macOS scaling is done at the texture level
	}
	if runtime.GOOS == "windows" {
		xScale, _ := w.(*window).viewport.GetContentScale()
		return xScale
	}

	return scaleAuto
}
