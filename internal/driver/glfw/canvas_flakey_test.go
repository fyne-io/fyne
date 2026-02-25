//go:build !no_glfw && !mobile && flakey

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_Resize(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	w.SetContent(content)
	ensureCanvasSize(t, w, fyne.NewSize(69, 36))

	size := fyne.NewSize(200, 100)
	runOnMain(func() {
		assert.NotEqual(t, size, content.Size())
	})

	w.Resize(size)
	ensureCanvasSize(t, w, size)

	runOnMain(func() {
		assert.Equal(t, size, content.Size())
	})
}
