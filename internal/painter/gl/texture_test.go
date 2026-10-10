//go:build !windows || !ci

package gl

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

// subImageContext records the pixels sent to replace part of a texture.
type subImageContext struct {
	context

	bound         Texture
	x, y          int
	width, height int
	data          []uint8
}

func (*subImageContext) ActiveTexture(uint32) {}

func (c *subImageContext) BindTexture(_ uint32, texture Texture) {
	c.bound = texture
}

func (*subImageContext) DeleteTexture(Texture) {}

func (*subImageContext) GetError() uint32 {
	return 0
}

func (c *subImageContext) TexSubImage2D(_ uint32, _, x, y, width, height int, _, _ uint32, data []uint8) {
	c.x, c.y, c.width, c.height = x, y, width, height
	c.data = append([]uint8{}, data...)
}

func TestPainter_updateImageTexture(t *testing.T) {
	pixels := image.NewRGBA(image.Rect(0, 0, 4, 3))
	for i := range pixels.Pix {
		pixels.Pix[i] = uint8(i)
	}
	img := canvas.NewImageFromImage(pixels)
	ctx := &subImageContext{}
	p := &painter{ctx: ctx, imageSizes: map[*canvas.Image]image.Point{img: image.Pt(4, 3)}}
	cache.SetTexture(img, cache.TextureType(testTexture()), nil)
	defer cache.DeleteTexture(img)

	t.Run("part of the rows", func(t *testing.T) {
		assert.True(t, p.updateImageTexture(img, image.Rect(1, 1, 3, 3)))
		assert.Equal(t, testTexture(), ctx.bound)
		assert.Equal(t, []int{1, 1, 2, 2}, []int{ctx.x, ctx.y, ctx.width, ctx.height})
		assert.Equal(t, []uint8{20, 21, 22, 23, 24, 25, 26, 27, 36, 37, 38, 39, 40, 41, 42, 43}, ctx.data)
	})

	t.Run("whole rows", func(t *testing.T) {
		assert.True(t, p.updateImageTexture(img, image.Rect(0, 2, 4, 3)))
		assert.Equal(t, []int{0, 2, 4, 1}, []int{ctx.x, ctx.y, ctx.width, ctx.height})
		assert.Equal(t, pixels.Pix[32:], ctx.data)
	})

	t.Run("beyond the image", func(t *testing.T) {
		assert.True(t, p.updateImageTexture(img, image.Rect(3, 2, 9, 9)))
		assert.Equal(t, []int{3, 2, 1, 1}, []int{ctx.x, ctx.y, ctx.width, ctx.height})
		assert.Equal(t, pixels.Pix[44:], ctx.data)
	})

	t.Run("other pixel formats", func(t *testing.T) {
		other := image.NewNRGBA(image.Rect(10, 10, 14, 13))
		other.SetNRGBA(12, 11, color.NRGBA{R: 0xff, A: 0x80})
		img.Image = other
		defer func() { img.Image = pixels }()

		assert.True(t, p.updateImageTexture(img, image.Rect(12, 11, 13, 12)))
		assert.Equal(t, []int{2, 1, 1, 1}, []int{ctx.x, ctx.y, ctx.width, ctx.height})
		assert.Equal(t, []uint8{0x80, 0, 0, 0x80}, ctx.data)
	})

	t.Run("image of another size", func(t *testing.T) {
		img.Image = image.NewRGBA(image.Rect(0, 0, 8, 3))
		defer func() { img.Image = pixels }()

		assert.False(t, p.updateImageTexture(img, image.Rect(1, 1, 3, 3)))
	})
}
