package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

func TestShadow_New(t *testing.T) {
	s := canvas.Shadow{Color: color.Black, BlurRadius: 3, Spread: 2, Offset: fyne.NewPos(4, 5), Variant: canvas.DropShadow}

	assert.Equal(t, color.Black, s.Color)
	assert.Equal(t, float32(3), s.BlurRadius)
	assert.Equal(t, float32(2), s.Spread)
	assert.Equal(t, fyne.NewPos(4, 5), s.Offset)
	assert.Equal(t, canvas.DropShadow, s.Variant)
}
