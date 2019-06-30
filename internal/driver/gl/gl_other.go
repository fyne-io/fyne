// +build !darwin arm arm64
// +build !windows

package gl

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := scaleInt(canvas, size.Width)
	winHeight := scaleInt(canvas, size.Height)

	setViewport(winWidth, winHeight)
}

func forceWindowRefresh(_ *window) {
}
