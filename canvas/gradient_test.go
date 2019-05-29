package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestNewLinearGradient(t *testing.T) {
	// Horizontal
	horizontal, err := NewLinearGradient(color.Black, color.Transparent, GradientDirectionHorizontal)
	if err != nil {
		t.Errorf("Expected no error but received: %v", err)
	}
	assert.Equal(t, horizontal.Center(), fyne.NewPos(0, 0))

	img := horizontal.Generator(50, 5)
	assert.Equal(t, img.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, img.At(25, i), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, img.At(50, 0), color.RGBA{0, 0, 0, 0x00})
	horizontal.SetCenter(fyne.NewPos(3, 3))
	assert.Equal(t, horizontal.Center(), fyne.NewPos(3, 3))

	// Vertical
	vertical, err := NewLinearGradient(color.Black, color.Transparent, GradientDirectionVertical)
	if err != nil {
		t.Errorf("Expected no error but received: %v", err)
	}
	imgVert := vertical.Generator(5, 50)
	assert.Equal(t, imgVert.At(0, 0), color.RGBA{0, 0, 0, 0xff})
	for i := 0; i < 5; i++ {
		assert.Equal(t, imgVert.At(i, 25), color.RGBA{0, 0, 0, 0x7f})
	}
	assert.Equal(t, imgVert.At(50, 0), color.RGBA{0, 0, 0, 0x00})

	// Radial and offsets
	circle, err := NewLinearGradient(color.Black, color.Transparent, GradientDirectionCircular)
	if err != nil {
		t.Errorf("Expected no error but received: %v", err)
	}
	imgCircle := circle.Generator(10, 10)
	assert.Equal(t, imgCircle.At(5, 5), color.RGBA{0, 0, 0, 0xff})
	assert.Equal(t, imgCircle.At(4, 5), color.RGBA{0, 0, 0, 0xcc})
	assert.Equal(t, imgCircle.At(3, 5), color.RGBA{0, 0, 0, 0x99})
	assert.Equal(t, imgCircle.At(2, 5), color.RGBA{0, 0, 0, 0x66})
	assert.Equal(t, imgCircle.At(1, 5), color.RGBA{0, 0, 0, 0x33})

	circle.SetCenter(fyne.NewPos(1, 1))
	imgCircleOffset := circle.Generator(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0x79})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0x50})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x26})

	circle.SetCenter(fyne.NewPos(-1, -1))
	imgCircleOffset = circle.Generator(10, 10)
	assert.Equal(t, imgCircleOffset.At(5, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(4, 5), color.RGBA{0, 0, 0, 0xd5})
	assert.Equal(t, imgCircleOffset.At(3, 5), color.RGBA{0, 0, 0, 0xc3})
	assert.Equal(t, imgCircleOffset.At(2, 5), color.RGBA{0, 0, 0, 0xa0})
	assert.Equal(t, imgCircleOffset.At(1, 5), color.RGBA{0, 0, 0, 0x79})

	// Test error response
	_, err = NewLinearGradient(nil, nil, GradientDirectionCircular)
	if err == nil {
		t.Error("Expected an error but received none")
	}
}
