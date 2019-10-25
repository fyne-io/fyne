// +build !ci

package glfw

func repaintWindow(w *window) {
	runOnMain(func() {
		d.(*gLDriver).repaintWindow(w)
	})
}
