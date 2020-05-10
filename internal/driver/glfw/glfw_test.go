// +build !ci

package glfw

import "time"

func repaintWindow(w *window) {
	for !running() {
		// Wait for GLFW loop to be running.
		// If we try to create windows before the context is created, this will fail with an exception.
		time.Sleep(time.Millisecond * 10)
	}
	runOnDraw(w, func() {
		d.(*gLDriver).repaintWindow(w)
	})
}
