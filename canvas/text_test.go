package canvas

import (
	"image/color"
	"testing"

	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestText_MinSize(t *testing.T) {
	text := NewText("Test", color.RGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)

	text = NewText("Test2", color.RGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
}

func TestText_MinSize_NoMultiline(t *testing.T) {
	text := NewText("Break", color.RGBA{0, 0, 0, 0xff})
	min := text.MinSize()

	text = NewText("Bre\nak", color.RGBA{0, 0, 0, 0xff})
	min2 := text.MinSize()
	assert.True(t, min2.Width > min.Width)
	assert.True(t, min2.Height == min.Height)
}
