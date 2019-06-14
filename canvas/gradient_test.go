package canvas

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinearGradient(t *testing.T) {
	// Horizontal
	horizontal := NewLinearGradient(color.Black, color.Transparent, GradientDirectionHorizontal)
	assert.Equal(t, horizontal.CenterOffsetX, 0.0)
	assert.Equal(t, horizontal.CenterOffsetY, 0.0)

	img := horizontal.Generator(50, 5)
	assert.Equal(t, img.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, img.At(25, i), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, img.At(50, 0), color.RGBA{0, 0, 0, 0x00})
	horizontal.CenterOffsetX = -0.4
	horizontal.CenterOffsetY = 0.35
	assert.Equal(t, horizontal.CenterOffsetX, -0.4)
	assert.Equal(t, horizontal.CenterOffsetY, 0.35)

	// Vertical
	vertical := NewLinearGradient(color.Black, color.Transparent, GradientDirectionVertical)
	imgVert := vertical.Generator(5, 50)
	assert.Equal(t, imgVert.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, imgVert.At(i, 25), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, imgVert.At(50, 0), color.RGBA{0, 0, 0, 0x00})

	// Radial and offsets
	circle := NewLinearGradient(color.Black, color.Transparent, GradientDirectionCircular)
	imgCircle := circle.Generator(10, 10)
	assert.Equal(t, imgCircle.At(5, 5), color.RGBA{0, 0, 0, 0xff})
	assert.Equal(t, imgCircle.At(4, 5), color.RGBA{0, 0, 0, 0xcc})
	assert.Equal(t, imgCircle.At(3, 5), color.RGBA{0, 0, 0, 0x99})
	assert.Equal(t, imgCircle.At(2, 5), color.RGBA{0, 0, 0, 0x66})
	assert.Equal(t, imgCircle.At(1, 5), color.RGBA{0, 0, 0, 0x33})

	circle.CenterOffsetX = 0.1
	circle.CenterOffsetY = 0.1
	imgCircleOffset := circle.Generator(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0x79})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0x50})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x26})

	circle.CenterOffsetX = -0.1
	circle.CenterOffsetY = -0.1
	imgCircleOffset = circle.Generator(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xd5})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x79})
}
