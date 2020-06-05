package canvas_test

import (
	"image"
	"testing"

	"fyne.io/fyne/canvas"

	"github.com/stretchr/testify/assert"
)

func TestRasterFromImage(t *testing.T) {
	source := image.Rect(2, 2, 4, 4)
	dest := canvas.NewRasterFromImage(source)
	img := dest.Generator(6, 6)

	// image.Rect is a 16 bit color model
	_, _, _, a := img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0x0000), a)
	_, _, _, a = img.At(2, 2).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = img.At(4, 4).RGBA()
	assert.Equal(t, uint32(0x0000), a)
}
