//go:build linux || freebsd || openbsd || netbsd
// +build linux freebsd openbsd netbsd

package glfw

import "fyne.io/fyne/v2"

func (w *window) platformResize(canvasSize fyne.Size) {
	w.canvas.Resize(canvasSize)
}
