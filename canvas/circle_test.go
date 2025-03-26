package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

func TestCircle_MinSize(t *testing.T) {
	circle := canvas.NewCircle(color.Black)
	min := circle.MinSize()

	assert.Positive(t, min.Width)
	assert.Positive(t, min.Height)
}

func TestCircle_FillColor(t *testing.T) {
	c := color.White
	circle := canvas.NewCircle(c)

	assert.Equal(t, c, circle.FillColor)
}

func TestCircle_Resize(t *testing.T) {
	targetWidth := float32(50)
	targetHeight := float32(50)
	circle := canvas.NewCircle(color.White)
	start := circle.Size()
	assert.Equal(t, float32(0), start.Height)
	assert.Equal(t, float32(0), start.Width)

	circle.Resize(fyne.NewSize(targetWidth, targetHeight))
	target := circle.Size()
	assert.Equal(t, targetWidth, target.Width)
	assert.Equal(t, targetHeight, target.Height)

}

func TestCircle_Move(t *testing.T) {
	circle := canvas.NewCircle(color.White)
	circle.Resize(fyne.NewSize(50, 50))

	start := fyne.Position{X: 0, Y: 0}
	assert.Equal(t, circle.Position(), start)

	target := fyne.Position{X: 10, Y: 75}
	circle.Move(target)
	assert.Equal(t, circle.Position(), target)
}
