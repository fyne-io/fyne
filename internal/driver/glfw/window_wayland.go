//go:build wayland && (linux || freebsd || openbsd || netbsd)

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

// GetWindowHandle returns the window handle. Only implemented for X11 currently.
func (w *window) GetWindowHandle() string {
	return "" // TODO: Find a way to get the Wayland handle for xdg_foreign protocol. Return "wayland:{id}".
}

// assert we are implementing fyne.NativeWindow
var _ fyne.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any) error) error {
	var err error
	done := make(chan struct{})
	runOnMain(func() {
		// TODO: define driver.WaylandWindowContext and pass window handle
		err = f(driver.UnknownContext{})
		close(done)
	})
	<-done
	return err
}
