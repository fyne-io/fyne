//go:build !no_glfw && !mobile

package glfw

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
)

func assertCanvasSize(t *testing.T, w *window, size fyne.Size) {
	if runtime.GOOS == "linux" {
		// TODO: find the root cause for these problems and solve them without additional repaint
		// fixes issues where the window does not have the correct size
		waitForCanvasSize(t, w, size, false)
	}
	assert.Equal(t, size, w.canvas.Size())
}

func ensureCanvasSize(t *testing.T, w *window, size fyne.Size) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// TODO: find the root cause for these problems and solve them without additional repaint
		// fixes issues where the window does not have the correct size
		waitForCanvasSize(t, w, size, true)
	}
	require.Equal(t, size, w.canvas.Size())
}

func repaintWindow(w *window) {
	// Wait for GLFW loop to be running.
	// If we try to paint windows before the context is created, we will end up on the wrong thread.
	<-w.driver.waitForStart

	runOnDraw(w, func() {
		d.repaintWindow(w)
	})

	time.Sleep(time.Millisecond * 150) // wait for the frames to be rendered... o
}

func waitForCanvasSize(t *testing.T, w *window, size fyne.Size, resizeIfNecessary bool) {
	attempts := 0
	for {
		if w.canvas.Size() == size {
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
