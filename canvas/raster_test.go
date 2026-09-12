package canvas_test

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

func TestRasterFromImage(t *testing.T) {
	source := image.Rect(2, 2, 4, 4)
	dest := canvas.NewRasterFromImage(source)
	img := dest.Generator(6, 6)

	// the source is drawn from the top left, its bounds offset is not a position
	// image.Rect is a 16 bit color model
	_, _, _, a := img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = img.At(2, 2).RGBA()
	assert.Equal(t, uint32(0x0000), a)
}

func TestRasterFromImage_offsetBounds(t *testing.T) {
	red := color.NRGBA{R: 0xff, A: 0xff}
	blue := color.NRGBA{B: 0xff, A: 0xff}
	source := image.NewNRGBA(image.Rect(2, 3, 10, 10)) // 8x7 pixels, not at the origin
	source.SetNRGBA(2, 3, red)
	source.SetNRGBA(9, 9, blue)
	dest := canvas.NewRasterFromImage(source)

	// the top left source pixel lands in the corner of the raster, whether it is
	// generated at the source size, larger (padded) or smaller (truncated)
	for _, img := range []image.Image{dest.Generator(8, 7), dest.Generator(12, 12), dest.Generator(4, 4)} {
		assert.Equal(t, image.Point{}, img.Bounds().Min)
		assert.Equal(t, red, color.NRGBAModel.Convert(img.At(0, 0)))
	}

	assert.Equal(t, blue, color.NRGBAModel.Convert(dest.Generator(8, 7).At(7, 6)))

	// a raster with no area still has to be a finite image, the painter scales whatever it gets
	empty := dest.Generator(0, 5)
	assert.Equal(t, image.Rect(0, 0, 0, 5), empty.Bounds())
}
