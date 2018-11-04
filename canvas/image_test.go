package canvas

import "image"
import "testing"

import "github.com/stretchr/testify/assert"

func TestImageFromImage(t *testing.T) {
	source := image.Rect(2, 2, 4, 4)
	dest := NewImageFromImage(source)

	// image.Rect is a 16 bit colour model
	_, _, _, a := dest.PixelColor(0, 0, 6, 6).RGBA()
	assert.Equal(t, uint32(0x0000), a)
	_, _, _, a = dest.PixelColor(2, 2, 6, 6).RGBA()
	assert.Equal(t, uint32(0xffff), a)
	_, _, _, a = dest.PixelColor(4, 4, 6, 6).RGBA()
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
