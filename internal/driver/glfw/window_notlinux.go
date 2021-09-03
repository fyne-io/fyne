// +build !linux

package glfw

import "fyne.io/fyne/v2"

func (w *window) platformResize(canvasSize fyne.Size) {
	d, ok := fyne.CurrentApp().Driver().(*gLDriver)
	if !ok { // don't wait to redraw in this way if we are running on test

		// A resize event is triggered on the main thread, whereas the
		// related rendering operations (draw calls) are signaled on a
		// different thread, the draw thread.
		//
		// Theoretically, when resizing a window, the underlying hardware
		// frame buffer resizes (call from the event thread) concurrently
		// with other relevant draw calls (especially glClear and glClearColor),
		// which are issued from the draw thread. Intrinsically, there is
		// a data race on the frame buffer. (e.g. a race on the glClear and
		// the framebuffer resize) Hence, the following Resize must run on
		// the draw thread too. Otherwise, it executes graphics calls
		// concurrently with the draw thread, and the data race happens
		// and will result in a crash.
		//
		// See https://github.com/fyne-io/fyne/issues/2188
		runOnDraw(w, func() { w.canvas.Resize(canvasSize) })
		return
	}

	runOnDraw(w, func() {
		w.canvas.Resize(canvasSize)
		d.repaintWindow(w)
	})
}
