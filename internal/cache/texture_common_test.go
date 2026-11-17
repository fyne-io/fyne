package cache

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextureDirty(t *testing.T) {
	img := &dummyWidget{}
	SetTextureDirty(img, image.Rect(1, 1, 3, 3))
	assert.True(t, TakeTextureDirty(img).Empty(), "nothing is uploaded to be out of date")

	SetTexture(img, NoTexture, nil)
	defer DeleteTexture(img)
	assert.True(t, TakeTextureDirty(img).Empty())

	SetTextureDirty(img, image.Rect(1, 1, 3, 3))
	SetTextureDirty(img, image.Rect(5, 2, 6, 8))
	assert.Equal(t, image.Rect(1, 1, 6, 8), TakeTextureDirty(img))
	assert.True(t, TakeTextureDirty(img).Empty())

	SetTextureDirty(img, image.Rect(1, 1, 3, 3))
	DeleteTexture(img)
	SetTexture(img, NoTexture, nil)
	assert.True(t, TakeTextureDirty(img).Empty(), "a new texture has every pixel")
}
