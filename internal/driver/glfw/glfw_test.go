//go:build !no_glfw && !mobile

package glfw

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func repaintWindow(w *safeWindow) {
	runOnMain(func() {
		w.RunWithContext(func() {
			d.repaintWindow(w.window)
		})
	})

	time.Sleep(time.Millisecond * 150) // wait for the frames to be rendered... o
}

func waitForCanvasSize(t *testing.T, w *safeWindow, size fyne.Size, resizeIfNecessary bool) {
	attempts := 0
	for {
		if w.Canvas().Size() == size {
			break
		}
		attempts++
		if !assert.Less(t, attempts, 100, "canvas did not get correct size in time") {
			break
		}
		if resizeIfNecessary && attempts%20 == 0 {
			// sometimes the resize does not seem to reach the actual window at all
			w.Resize(size)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
