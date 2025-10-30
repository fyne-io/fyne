package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

func TestShadow_New(t *testing.T) {
	s := canvas.NewShadow(color.Black, 3, fyne.NewPos(4, 5), canvas.DropShadow)

	assert.Equal(t, color.Black, s.ShadowColor)
	assert.Equal(t, float32(3), s.ShadowSoftness)
	assert.Equal(t, fyne.NewPos(4, 5), s.ShadowOffset)
	assert.Equal(t, canvas.DropShadow, s.ShadowType)
}
