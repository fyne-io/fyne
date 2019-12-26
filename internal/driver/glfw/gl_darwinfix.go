// +build darwin,!arm,!arm64

package glfw

import "C"

func updateGLContext(w *window) {
}

// This forces a redraw of the window as we resize
func forceWindowRefresh(w *window) {
	w.RunWithContext(func() {
		w.canvas.paint(w.canvas.Size())

		w.viewport.SwapBuffers()
	})
}
