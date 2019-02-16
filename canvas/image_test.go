package canvas

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRasterFromImage(t *testing.T) {
	source := image.Rect(2, 2, 4, 4)
	dest := NewRasterFromImage(source)
	img := dest.Raster(6, 6)

	// image.Rect is a 16 bit colour model
	_, _, _, a := img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0x0000), a)
	_, _, _, a = img.At(2, 2).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = img.At(4, 4).RGBA()
	assert.Equal(t, uint32(0x0000), a)
}

func TestImage_AlphaDefault(t *testing.T) {
	img := &Image{}

	assert.Equal(t, 1.0, img.Alpha())
}

func TestImage_TranslucencyDefault(t *testing.T) {
	img := &Image{}

	assert.Equal(t, 0.0, img.Translucency)
}
