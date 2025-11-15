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

func TestEllipse_MinSize(t *testing.T) {
	ellipse := canvas.NewEllipse(color.Black)
	min := ellipse.MinSize()

	assert.Positive(t, min.Width)
	assert.Positive(t, min.Height)
}

func TestEllipse_FillColor(t *testing.T) {
	c := color.White
	ellipse := canvas.NewEllipse(c)

	assert.Equal(t, c, ellipse.FillColor)
}

func TestEllipse_Resize(t *testing.T) {
	targetWidth := float32(50)
	targetHeight := float32(80)
	ellipse := canvas.NewEllipse(color.White)
	start := ellipse.Size()
	assert.True(t, start.Height == 0)
	assert.True(t, start.Width == 0)

	ellipse.Resize(fyne.NewSize(targetWidth, targetHeight))
	target := ellipse.Size()
	assert.True(t, target.Height == targetHeight)
	assert.True(t, target.Width == targetWidth)
}

func TestEllipse_Move(t *testing.T) {
	ellipse := canvas.NewEllipse(color.White)
	ellipse.Resize(fyne.NewSize(80, 50))

	start := fyne.Position{X: 0, Y: 0}
	assert.Equal(t, ellipse.Position(), start)

	target := fyne.Position{X: 10, Y: 75}
	ellipse.Move(target)
	assert.Equal(t, ellipse.Position(), target)
}

func TestEllipse_Rotation(t *testing.T) {
	ellipse := &canvas.Ellipse{
		FillColor:   color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		StrokeColor: color.Black,
		StrokeWidth: 2,
		Angle:       30,
	}

	ellipse.Resize(fyne.NewSize(50, 90))
	test.AssertObjectRendersToMarkup(t, "ellipse.xml", ellipse)

	c := software.NewCanvas()
	c.SetContent(ellipse)
	c.Resize(fyne.NewSize(60, 60))

	ellipse.Resize(fyne.NewSize(50, 30))
	ellipse.Angle = 0
	test.AssertRendersToImage(t, "ellipse_stroke.png", c)

	ellipse.StrokeWidth = 0
	test.AssertRendersToImage(t, "ellipse.png", c)

	ellipse.Angle = -20
	test.AssertRendersToImage(t, "ellipse_-20.png", c)

	ellipse.Angle = -55
	test.AssertRendersToImage(t, "ellipse_-55.png", c)

	ellipse.StrokeWidth = 3
	ellipse.Angle = 45
	test.AssertRendersToImage(t, "ellipse_45_stroke.png", c)

	ellipse.Angle = 90
	test.AssertRendersToImage(t, "ellipse_90_stroke.png", c)

	ellipse.Angle = 135
	test.AssertRendersToImage(t, "ellipse_135_stroke.png", c)
}
