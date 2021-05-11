package color

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnmultiplyAlpha(t *testing.T) {
	c := color.RGBA{R: 100, G: 100, B: 100, A: 100}
	r, g, b, a := unmultiplyAlpha(c)

	assert.Equal(t, 255, r)
	assert.Equal(t, 255, g)
	assert.Equal(t, 255, b)
	assert.Equal(t, 100, a)

	c = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	r, g, b, a = unmultiplyAlpha(c)

	assert.Equal(t, 100, r)
	assert.Equal(t, 100, g)
	assert.Equal(t, 100, b)
	assert.Equal(t, 255, a)
}
