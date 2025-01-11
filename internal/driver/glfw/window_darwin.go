//go:build darwin

package glfw

import (
	"fyne.io/fyne/v2/driver"
)

// assert we are implementing driver.NativeWindow
var _ driver.NativeWindow = (*window)(nil)

func (w *window) RunNative(f func(any)) {
	context := driver.MacWindowContext{}
	if v := w.view(); v != nil {
		context.NSWindow = uintptr(v.GetCocoaWindow())
	}

	f(context)
}
