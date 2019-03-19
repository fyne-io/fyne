package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestLine_MinSize(t *testing.T) {
	line := NewLine(color.Black)
	min := line.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestLine_Resize(t *testing.T) {
	line := NewLine(color.Black)

	line.Resize(fyne.NewSize(10, 0))
	size := line.Size()

	assert.Equal(t, size.Width, 10)
	assert.Equal(t, size.Height, 0)
}

func TestLine_StrokeColor(t *testing.T) {
	c := color.White
	line := NewLine(c)

	assert.Equal(t, c, line.StrokeColor)
}
