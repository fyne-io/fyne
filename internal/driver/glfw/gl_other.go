// +build !darwin arm arm64
// +build !windows

package glfw

import util "fyne.io/fyne/internal"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := util.ScaleInt(canvas, size.Width)
	winHeight := util.ScaleInt(canvas, size.Height)

	w.canvas.painter.SetOutputSize(winWidth, winHeight)
}

func forceWindowRefresh(_ *window) {
}
