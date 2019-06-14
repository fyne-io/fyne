package gl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := scaleInt(canvas, size.Width)
	winHeight := scaleInt(canvas, size.Height)

	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
}

// This forces a redraw of the window as we resize or move
func forceWindowRefresh(w *window) {
	w.runWithContext(func() {
		updateGLContext(w)

		w.canvas.paint(w.canvas.Size())
		w.viewport.SwapBuffers()
	})
}
