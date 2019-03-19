package canvas

import (
	"image/color"
	"testing"

	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestCircle_MinSize(t *testing.T) {
	circle := NewCircle(color.Black)
	min := circle.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestCircle_FillColor(t *testing.T) {
	c := color.White
	circle := NewCircle(c)

	assert.Equal(t, c, circle.FillColor)
}
