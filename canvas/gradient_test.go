package canvas

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHorizontalGradient(t *testing.T) {
	horizontal := NewHorizontalGradient(color.Black, color.Transparent)

	img := horizontal.Generate(50, 5)
	assert.Equal(t, img.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, img.At(25, i), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, img.At(50, 0), color.RGBA{0, 0, 0, 0x00})
}

func TestNewVerticalGradient(t *testing.T) {
	vertical := NewVerticalGradient(color.Black, color.Transparent)
	imgVert := vertical.Generate(5, 50)
	assert.Equal(t, imgVert.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, imgVert.At(i, 25), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, imgVert.At(50, 0), color.RGBA{0, 0, 0, 0x00})
}

func TestNewRadialGradient(t *testing.T) {
	circle := NewRadialGradient(color.Black, color.Transparent)
	imgCircle := circle.Generate(10, 10)
	assert.Equal(t, imgCircle.At(5, 5), color.RGBA{0, 0, 0, 0xff})
	assert.Equal(t, imgCircle.At(4, 5), color.RGBA{0, 0, 0, 0xcc})
	assert.Equal(t, imgCircle.At(3, 5), color.RGBA{0, 0, 0, 0x99})
	assert.Equal(t, imgCircle.At(2, 5), color.RGBA{0, 0, 0, 0x66})
	assert.Equal(t, imgCircle.At(1, 5), color.RGBA{0, 0, 0, 0x33})

	circle.CenterOffsetX = 0.1
	circle.CenterOffsetY = 0.1
	imgCircleOffset := circle.Generate(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0x79})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0x50})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x26})

	circle.CenterOffsetX = -0.1
	circle.CenterOffsetY = -0.1
	imgCircleOffset = circle.Generate(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xd5})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x79})
}
