// +build !ci

package gl

import (
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

var d = NewGLDriver()

func init() {
	runtime.LockOSThread()
}

// TestMain makes sure that our driver is running on the main thread.
// This must be done for some of our tests to function correctly.
func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	d.Run()
}

func TestWindow_SetTitle(t *testing.T) {
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_PixelSize(t *testing.T) {
	w := d.CreateWindow("Test")

	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(fyne.NewContainer(rect))
	w.Canvas().Refresh(w.Content())

	scale := w.Canvas().Scale()
	winW, winH := w.(*window).sizeOnScreen()
	assert.Equal(t, int(100*scale), winW)
	assert.Equal(t, int(100*scale), winH)
}
