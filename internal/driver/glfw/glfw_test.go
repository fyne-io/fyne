// +build !ci
// +build !mobile

package glfw

import "time"

func repaintWindow(w *window) {
	for !running() {
		// Wait for GLFW loop to be running.
		// If we try to paint windows before the context is created, we will end up on the wrong thread.
		time.Sleep(time.Millisecond * 10)
	}
	runOnDraw(w, func() {
		d.(*gLDriver).repaintWindow(w)
	})

	time.Sleep(time.Millisecond * 150) // wait for the frames to be rendered... o
}
