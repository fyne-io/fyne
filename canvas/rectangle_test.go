package canvas_test

import (
	"image/color"
	"testing"

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

func TestRectangle_Radius(t *testing.T) {
	rect := &canvas.Rectangle{
		FillColor:   color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor: color.NRGBA{R: 255, G: 120, B: 0, A: 255},
		StrokeWidth: 2.0,
		Radius:      25}

	assert.Equal(t, float32(25.0), rect.Radius)
}

func TestRectangle_StrokeColor(t *testing.T) {
	orange := color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	rect := &canvas.Rectangle{
		FillColor:   color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor: orange,
		StrokeWidth: 5,
		Radius:      25.0}

	assert.Equal(t, orange, rect.StrokeColor)

}
