//go:build darwin

package glfw

import (
	"fyne.io/fyne/v2/driver"
)

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any) error) error {
	var err error
	done := make(chan struct{})
	runOnMain(func() {
		err = f(driver.MacWindowContext{
			NSWindow: uintptr(w.view().GetCocoaWindow()),
		})
		close(done)
	})
	<-done
	return err
}
