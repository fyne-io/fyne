package test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func Test_canvas_Capture(t *testing.T) {
	c := NewCanvas()
	c.Size()

	img := c.Capture()
	assert.True(t, img.Bounds().Size().X > 0)
	assert.True(t, img.Bounds().Size().Y > 0)

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
	c.Size()

	img := c.Capture()
	assert.True(t, img.Bounds().Size().X > 0)
	assert.True(t, img.Bounds().Size().Y > 0)

	r1, g1, b1, a1 := color.Transparent.RGBA()
	r2, g2, b2, a2 := img.At(1, 1).RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)
}
