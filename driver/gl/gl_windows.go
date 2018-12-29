package gl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := scaleInt(canvas, size.Width)
	winHeight := scaleInt(canvas, size.Height)

	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
}

// This forces a redraw of the window as we resize
func updateWinSize(w *window) {
	w.viewport.MakeContextCurrent()
	updateGLContext(w)

	size := w.canvas.Size()
	w.canvas.paint(size)

	w.viewport.SwapBuffers()
	glfw.DetachCurrentContext()
}
