// +build !darwin arm arm64
// +build !windows

package glfw

import "fyne.io/fyne/internal"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := internal.ScaleInt(canvas, size.Width)
	winHeight := internal.ScaleInt(canvas, size.Height)

	w.canvas.painter.SetOutputSize(winWidth, winHeight)
}

func forceWindowRefresh(_ *window) {
}
