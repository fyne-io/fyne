//go:build !windows || !ci

package gl

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureImage_At(t *testing.T) {
	// a 2x1 buffer, the second pixel having the reduced alpha that blending a
	// shadow or an anti-aliased edge leaves in the framebuffer
	img := &captureImage{
		pix:    []uint8{0xff, 0x00, 0x00, 0xff, 0x00, 0x20, 0x40, 0x80},
		width:  2,
		height: 1,
	}

	assert.Equal(t, color.RGBA{R: 0xff, A: 0xff}, img.At(0, 0))
	assert.Equal(t, color.RGBA{G: 0x20, B: 0x40, A: 0xff}, img.At(1, 0))
}

func TestCaptureImage_SubImage(t *testing.T) {
	img := &captureImage{
		pix:    make([]uint8, 4*4*4),
		width:  4,
		height: 4,
	}

	sub := img.SubImage(image.Rect(1, 1, 3, 4))
	assert.Equal(t, image.Rect(1, 1, 3, 4), sub.Bounds())
}
