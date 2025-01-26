//go:build darwin

package glfw

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit

#import <stdbool.h>

void setFullScreen(bool full, void *window);
*/
import "C"
import (
	"runtime"

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

func (w *window) doSetFullScreen(full bool) {
	if runtime.GOOS == "darwin" {
		win := w.view().GetCocoaWindow()
		C.setFullScreen(C.bool(full), win)
		return
	}
}
