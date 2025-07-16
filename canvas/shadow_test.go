package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestShadow_ShadowPaddings(t *testing.T) {
	b := &shadow{
		ShadowColor:    color.NRGBA{0, 0, 0, 128},
		ShadowSoftness: 4,
		ShadowOffset:   fyne.NewPos(2, 3),
		ShadowType:     DropShadow,
	}

	pads := b.ShadowPaddings()
	assert.Equal(t, float32(6), pads[0])
	assert.Equal(t, float32(1), pads[1])
	assert.Equal(t, float32(2), pads[2])
	assert.Equal(t, float32(7), pads[3])
}
