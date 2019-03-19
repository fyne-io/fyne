package canvas

import (
	"image/color"
	"testing"

	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestRectangle_MinSize(t *testing.T) {
	rect := NewRectangle(color.Black)
	min := rect.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestRectangle_FillColor(t *testing.T) {
	c := color.White
	rect := NewRectangle(c)

	assert.Equal(t, c, rect.FillColor)
}
