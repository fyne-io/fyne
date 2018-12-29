// +build darwin,!arm,!arm64

package gl

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include <Cocoa/Cocoa.h>
void cocoa_update_nsgl_context(void* id) {
    NSOpenGLContext *ctx = id;
    [ctx update];
}
*/
import "C"

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// This fixes the initially blank window on macOS Mojave
// We can remove this (and the gl_other.go) once GLFW 3.3 arrives in go_gl
func updateGLContext(w *window) {
	if w.painted < 2 {
		ctx := w.viewport.GetNSGLContext()
		C.cocoa_update_nsgl_context(unsafe.Pointer(ctx))
		w.painted++
	}
}

// This forces a redraw of the window as we resize
func updateWinSize(w *window) {
	w.viewport.MakeContextCurrent()

	size := w.canvas.Size()
	w.canvas.paint(size)

	w.viewport.SwapBuffers()
	glfw.DetachCurrentContext()

}
