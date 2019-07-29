// +build !darwin arm arm64
// +build !windows

package glfw

import "fyne.io/fyne/internal/driver"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := driver.ScaleInt(canvas, size.Width)
	winHeight := driver.ScaleInt(canvas, size.Height)

	w.canvas.painter.SetOutputSize(winWidth, winHeight)
}

func forceWindowRefresh(_ *window) {
}
