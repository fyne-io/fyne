package cache

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageCacheGet(t *testing.T) {
	ResetThemeCaches()
	img := image.NewNRGBA(image.Rect(0, 0, 25, 25))
	SetImage("shared.png", nil, img)
	assert.Equal(t, 1, images.Len())

	newImg := GetImage("shared.png", nil)
	assert.Same(t, img, newImg)

	miss := GetImage("missing.png", nil)
	assert.Nil(t, miss)
}

func TestImageCacheReset(t *testing.T) {
	ResetThemeCaches()
	SetImage("shared.png", nil, image.NewNRGBA(image.Rect(0, 0, 25, 25)))
	assert.Equal(t, 1, images.Len())

	ResetThemeCaches()
	assert.Equal(t, 0, images.Len())
}
