package canvas_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestRectangle_MinSize(t *testing.T) {
	rect := canvas.NewRectangle(color.Black)
	min := rect.MinSize()

	assert.Positive(t, min.Width)
	assert.Positive(t, min.Height)
}

func TestRectangle_FillColor(t *testing.T) {
	c := color.White
	rect := canvas.NewRectangle(c)

	assert.Equal(t, c, rect.FillColor)
}

func TestRectangle_Radius(t *testing.T) {
	rect := &canvas.Rectangle{
		FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor:  color.NRGBA{R: 255, G: 120, B: 0, A: 255},
		StrokeWidth:  2.0,
		CornerRadius: 12}

	rect.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "rounded_rect.xml", rect)

	c := software.NewCanvas()
	c.SetContent(rect)
	c.Resize(fyne.NewSize(60, 60))
	test.AssertRendersToImage(t, "rounded_rect_stroke.png", c)

	rect.StrokeWidth = 0
	test.AssertRendersToImage(t, "rounded_rect.png", c)

	rect.Aspect = 2.0
	test.AssertRendersToImage(t, "rounded_rect_aspect.png", c)
}
