// +build !darwin arm arm64
// +build !windows

package gl

import "github.com/go-gl/gl/v3.2-core/gl"

func updateGLContext(w *window) {
	canvas := w.Canvas()
	size := canvas.Size()

	winWidth := scaleInt(canvas, size.Width)
	winHeight := scaleInt(canvas, size.Height)

	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
}

func updateWinSize(_ *window) {
}
