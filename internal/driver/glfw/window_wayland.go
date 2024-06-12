//go:build wayland && (linux || freebsd || openbsd || netbsd)

package glfw

import (
	"unsafe"

	"fyne.io/fyne/v2/driver"
)

// GetWindowHandle returns the window handle. Only implemented for X11 currently.
func (w *window) GetWindowHandle() string {
	return "" // TODO: Find a way to get the Wayland handle for xdg_foreign protocol. Return "wayland:{id}".
}

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	runOnMain(func() {
		var waylandSurface, eglSurface uintptr
		if v := w.view(); v != nil {
			waylandSurface = uintptr(unsafe.Pointer(v.GetWaylandWindow()))
			eglSurface = uintptr(unsafe.Pointer(v.GetEGLSurface()))
		}
		f(driver.WaylandWindowContext{
			WaylandSurface: waylandSurface,
			EGLSurface:     eglSurface,
		})
	})
}
