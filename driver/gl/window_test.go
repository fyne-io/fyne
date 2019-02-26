// +build !ci

package gl

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var d = NewGLDriver()

func TestWindow_SetTitle(t *testing.T) {
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_PixelSize(t *testing.T) {
	w := d.CreateWindow("Test")
	w.SetPadded(false)

	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(fyne.NewContainer(rect))
	w.Canvas().Refresh(w.Content())

	winW, winH := w.(*window).sizeOnScreen()
	assert.Equal(t, scaleInt(w.Canvas(), 100), winW)
	assert.Equal(t, scaleInt(w.Canvas(), 100), winH)
}

func TestWindow_Padded(t *testing.T) {
	w := d.CreateWindow("Test")
	content := canvas.NewRectangle(color.White)
	w.Canvas().SetScale(1.0)
	w.SetContent(content)

	width, _ := w.(*window).sizeOnScreen()
	assert.Equal(t, theme.Padding()*2+content.MinSize().Width, width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}

func TestWindow_SetPadded(t *testing.T) {
	w := d.CreateWindow("Test")
	w.SetPadded(false)

	content := canvas.NewRectangle(color.White)
	w.Canvas().SetScale(1.0)
	w.SetContent(content)

	width, _ := w.(*window).sizeOnScreen()
	assert.Equal(t, content.MinSize().Width, width)
	assert.Equal(t, 0, content.Position().X)

	w.SetPadded(true)
	width, _ = w.(*window).sizeOnScreen()
	assert.Equal(t, theme.Padding()*2+content.MinSize().Width, width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}
