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
	assert.True(t, start.Height == 0)
	assert.True(t, start.Width == 0)

	circle.Resize(fyne.NewSize(targetWidth, targetHeight))
	target := circle.Size()
	assert.True(t, target.Height == targetHeight)
	assert.True(t, target.Width == targetWidth)
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

func TestCircle_shadow(t *testing.T) {
	circle := &canvas.Circle{
		FillColor:   color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor: color.NRGBA{R: 255, G: 120, B: 0, A: 255},
		StrokeWidth: 2.0,
		Shadow: canvas.Shadow{
			Color:      color.White,
			Offset:     fyne.NewPos(8, 5),
			BlurRadius: 3,
			Variant:    canvas.DropShadow,
		},
	}

	circle.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "circle_shadow.xml", circle)

	c := software.NewCanvas()
	c.SetContent(circle)
	c.Resize(fyne.NewSize(170, 170))
	circle.Resize(fyne.NewSize(150, 150))
	circle.Move(fyne.NewPos(6, 6))
	test.AssertRendersToImage(t, "circle_stroke_shadow.png", c)

	circle.StrokeWidth = 0
	circle.Shadow.Variant = canvas.DropShadow
	test.AssertRendersToImage(t, "circle_shadow.png", c)
}
