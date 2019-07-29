package glfw

import "fyne.io/fyne/internal/driver"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := driver.ScaleInt(canvas, size.Width)
	winHeight := driver.ScaleInt(canvas, size.Height)

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
