//go:build !wayland && (linux || freebsd || openbsd || netbsd) && !wasm && !test_web_driver

package glfw

import "fyne.io/fyne/v2/driver"

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	var handle uintptr
	if v := w.view(); v != nil {
		handle = uintptr(v.GetX11Window())
	}
	runOnMain(func() {
		f(driver.X11WindowContext{
			WindowHandle: handle,
		})
	})
}
