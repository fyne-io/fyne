//go:build darwin

package glfw

import (
	"fyne.io/fyne/v2/driver"
)

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	runOnMain(func() {
		var nsWindow uintptr
		if v := w.view(); v != nil {
			nsWindow = uintptr(v.GetCocoaWindow())
		}
		f(driver.MacWindowContext{
			NSWindow: nsWindow,
		})
	})
}
