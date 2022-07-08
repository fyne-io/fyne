package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"github.com/stretchr/testify/assert"
)

func TestRectangle_MinSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	min := rect.MinSize()

	assert.True(t, min.Width > 0)
	assert.True(t, min.Height > 0)
}

func TestRectangle_FillColor(t *testing.T) {
	c := color.White
	rect := canvas.NewRectangle(c)

	assert.Equal(t, c, rect.FillColor)
}

func TestRectangle_StrokeColor(t *testing.T) {
	orange := color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	rect := canvas.Rectangle{
		FillColor:   orange,
		StrokeColor: red,
		Radius: fyne.RectangleRadius{
			Left:  30.0,
			Right: 30.0,
		},
	}
	rect.Resize((fyne.NewSize(300, 150)))

	assert.Equal(t, red, rect.StrokeColor)
	assert.Equal(t, float32(30.0), rect.Radius.Left)
}
