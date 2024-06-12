//go:build !wayland && (linux || freebsd || openbsd || netbsd) && !wasm && !test_web_driver

package glfw

import (
	"strconv"

	"fyne.io/fyne/v2/driver"
)

// GetWindowHandle returns the window handle. Only implemented for X11 currently.
func (w *window) GetWindowHandle() string {
	xid := uint(w.viewport.GetX11Window())
	return "x11:" + strconv.FormatUint(uint64(xid), 16)
}

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	var handle uintptr
	if v := w.view(); v != nil {
		handle = v.GetX11Window()
	}
	runOnMain(func() {
		f(driver.X11WindowContext{
			WindowHandle: handle,
		})
	})
}
