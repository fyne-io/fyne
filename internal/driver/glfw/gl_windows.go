package glfw

import "fyne.io/fyne/internal"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := internal.ScaleInt(canvas, size.Width)
	winHeight := internal.ScaleInt(canvas, size.Height)

	w.canvas.painter.SetOutputSize(winWidth, winHeight)
}

// This forces a redraw of the window as we resize or move
func forceWindowRefresh(w *window) {
	w.RunWithContext(func() {
		updateGLContext(w)

		w.canvas.paint(w.canvas.Size())
		w.viewport.SwapBuffers()
	})
}
