package painter_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/internal/painter"
	"github.com/stretchr/testify/assert"
)

func TestBlendPhotometric_Opaque(t *testing.T) {
	dst := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	src := color.NRGBA{R: 100, G: 150, B: 200, A: 255}

	out := painter.BlendPhotometric(dst, src)
	assert.Equal(t, src, out)
}

func TestBlendPhotometric_Transparent(t *testing.T) {
	dst := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	src := color.NRGBA{R: 100, G: 150, B: 200, A: 0}

	out := painter.BlendPhotometric(dst, src)
	assert.Equal(t, dst, out)
}

func TestBlendPhotometric_WhiteOverBlack50Percent(t *testing.T) {
	// Blending 50% white over black in linear space produces ~188 sRGB value (not 128)
	dst := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	src := color.NRGBA{R: 255, G: 255, B: 255, A: 128}

	out := painter.BlendPhotometric(dst, src)
	assert.Equal(t, uint8(188), out.R)
	assert.Equal(t, uint8(188), out.G)
	assert.Equal(t, uint8(188), out.B)
	assert.Equal(t, uint8(255), out.A)
}
