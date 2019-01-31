// +build !ci

package gl

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestWindow_SetTitle(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_PixelSize(t *testing.T) {
	d := NewGLDriver()
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
