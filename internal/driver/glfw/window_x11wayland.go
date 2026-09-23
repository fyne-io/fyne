//go:build ((x11 && wayland) || (!x11 && !wayland)) && (linux || freebsd || openbsd || netbsd) && !wasm && !test_web_driver

package glfw

import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/internal/build"
)

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	v := w.view()

	if build.IsWayland {
		context := driver.WaylandWindowContext{}
		if v != nil {
			context.WaylandSurface = uintptr(unsafe.Pointer(v.GetWaylandWindow()))
		}

		f(context)
		return
	}

	context := driver.X11WindowContext{}
	if v != nil {
		context.WindowHandle = uintptr(v.GetX11Window())
	}

	f(context)
}
