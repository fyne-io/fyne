package glfw

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := scaleInt(canvas, size.Width)
	winHeight := scaleInt(canvas, size.Height)

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
