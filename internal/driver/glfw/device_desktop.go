//go:build !wasm

package glfw

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/goos"
)

func (*glDevice) IsMobile() bool {
	return false
}

func (*glDevice) SystemScaleForWindow(w fyne.Window) float32 {
	if runtime.GOOS == goos.Darwin {
		return 1.0 // macOS scaling is done at the texture level
	}
	if runtime.GOOS == goos.Windows {
		xScale, _ := w.(*window).viewport.GetContentScale()
		return xScale
	}

	return scaleAuto
}

func connectKeyboard(*glCanvas) {
	// no-op, mobile web compatibility
}

func isMacOSRuntime() bool {
	return runtime.GOOS == goos.Darwin
}
