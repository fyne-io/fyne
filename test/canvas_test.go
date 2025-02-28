package test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func Test_canvas_Capture(t *testing.T) {
	c := NewCanvas()
	c.Resize(fyne.NewSquareSize(100))

	img := c.Capture()
	assert.Positive(t, img.Bounds().Size().X)
	assert.Positive(t, img.Bounds().Size().Y)

	r1, g1, b1, a1 := theme.Color(theme.ColorNameBackground).RGBA()
	r2, g2, b2, a2 := img.At(1, 1).RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)
}

func Test_canvas_InteractiveArea(t *testing.T) {
	c := NewCanvas()
	c.Resize(fyne.NewSize(600, 400))
	pos, size := c.InteractiveArea()
	assert.Equal(t, fyne.NewPos(2, 3), pos)
	assert.Equal(t, fyne.NewSize(596, 395), size)
}

func Test_canvas_PixelCoordinateAtPosition(t *testing.T) {
	c := NewCanvas().(*canvas)

	pos := fyne.NewPos(4, 4)
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 10, x)
	assert.Equal(t, 10, y)
}

func Test_canvas_TransparentCapture(t *testing.T) {
	c := NewTransparentCanvasWithPainter(nil)
	c.Resize(fyne.NewSquareSize(10))

	img := c.Capture()
	assert.Positive(t, img.Bounds().Size().X)
	assert.Positive(t, img.Bounds().Size().Y)

	r1, g1, b1, a1 := color.Transparent.RGBA()
	r2, g2, b2, a2 := img.At(1, 1).RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)
}

func Test_canvas_Resize(t *testing.T) {
	c := NewCanvas()
	assert.Equal(t, float32(100), c.Size().Width) // backwards compatible

	smallSize := fyne.NewSize(10, 5)
	c.Resize(smallSize)
	c.SetContent(widget.NewLabel("Bigger"))
	assert.NotEqual(t, float32(100), c.Size().Width) // backwards compatible
	assert.Greater(t, c.Size().Width, smallSize.Width)
	assert.Greater(t, c.Size().Height, smallSize.Height)

	c = NewCanvas()
	largeSize := fyne.NewSize(100, 75)
	c.Resize(largeSize)
	c.SetContent(widget.NewLabel("Smaller"))
	assert.Equal(t, largeSize, c.Size())
}
