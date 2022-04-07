package test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestTestCanvas_Capture(t *testing.T) {
	c := NewCanvas()
	c.Size()

	img := c.Capture()
	assert.True(t, img.Bounds().Size().X > 0)
	assert.True(t, img.Bounds().Size().Y > 0)

	r1, g1, b1, a1 := theme.BackgroundColor().RGBA()
	r2, g2, b2, a2 := img.At(1, 1).RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)
}

func TestTestCanvas_TransparentCapture(t *testing.T) {
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

func TestGlCanvas_PixelCoordinateAtPosition(t *testing.T) {
	c := NewCanvas().(*testCanvas)

	pos := fyne.NewPos(4, 4)
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 10, x)
	assert.Equal(t, 10, y)
}
