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

func TestSquare_FillColor(t *testing.T) {
	c := color.White
	rect := canvas.NewSquare(c)

	assert.Equal(t, c, rect.FillColor)
}

func TestSquare_Radius(t *testing.T) {
	rect := &canvas.Square{
		FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor:  color.NRGBA{R: 255, G: 120, B: 0, A: 255},
		StrokeWidth:  2.0,
		CornerRadius: 12}

	rect.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "rounded_square.xml", rect)

	c := software.NewCanvas()
	c.SetContent(rect)
	c.Resize(fyne.NewSize(60, 60))
	test.AssertRendersToImage(t, "rounded_rect_stroke.png", c)

	rect.StrokeWidth = 0
	test.AssertRendersToImage(t, "rounded_rect.png", c)
}
