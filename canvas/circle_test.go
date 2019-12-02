package canvas

import (
	"image/color"
	"testing"

	"fyne.io/fyne"

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

func TestCircle_Resize(t *testing.T) {
	targetWidth := 50
	targetHeight := 50
	circle := NewCircle(color.White)
	start := circle.Size()
	assert.True(t, start.Height == 0)
	assert.True(t, start.Width == 0)

	circle.Resize(fyne.NewSize(targetWidth, targetHeight))
	target := circle.Size()
	assert.True(t, target.Height == targetHeight)
	assert.True(t, target.Width == targetWidth)
}

func TestCircle_Move(t *testing.T) {
	circle := NewCircle(color.White)
	circle.Resize(fyne.NewSize(50, 50))

	start := fyne.Position{X: 0, Y: 0}
	assert.True(t, circle.Position() == start)

	target := fyne.Position{X: 10, Y: 75}
	circle.Move(target)
	assert.True(t, circle.Position() == target)
}

func TestCircle_RefreshDuringResize(t *testing.T) {
	circle := NewCircle(color.White)
	assert.True(t, circle.SkipRefreshDuringResize())
	circle.RefreshDuringResize = true
	assert.False(t, circle.SkipRefreshDuringResize())
}
