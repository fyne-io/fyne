package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestCanvas_Capture(t *testing.T) {
	c := NewCanvas()
	c.Size()

	img := c.Capture()
	assert.True(t, img.Bounds().Size().X > 0)
	assert.True(t, img.Bounds().Size().Y > 0)

	r1, g1, b1, a1 := testBackground.RGBA()
	r2, g2, b2, a2 := img.At(1, 1).RGBA()
	assert.Equal(t, r1, r2)
	assert.Equal(t, g1, g2)
	assert.Equal(t, b1, b2)
	assert.Equal(t, a1, a2)
}
