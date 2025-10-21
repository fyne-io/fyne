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

func TestBezierCurve_MinSizeAndStrokeColor(t *testing.T) {
	curve := canvas.NewQuadraticBezierCurve(fyne.NewPos(0, 0), fyne.NewPos(20, -5), fyne.NewPos(30, 0), color.Black)
	min := curve.MinSize()

	assert.Positive(t, min.Width)
	assert.Positive(t, min.Height)

	assert.Equal(t, color.Black, curve.StrokeColor)
}

func TestBezierCurve_Quadratic(t *testing.T) {
	curve := canvas.NewQuadraticBezierCurve(fyne.NewPos(0, 0), fyne.NewPos(30, 40), fyne.NewPos(50, 10), color.Black)
	curve.StrokeWidth = 2
	curve.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "bezier_curve_quadratic.xml", curve)

	c := software.NewCanvas()
	c.SetContent(curve)
	c.Resize(fyne.NewSize(70, 60))
	test.AssertRendersToImage(t, "bezier_curve_quadratic.png", c)

	curve.StartPoint = fyne.NewPos(-10, -10)
	test.AssertRendersToImage(t, "bezier_curve_quadratic.png", c)
}

func TestBezierCurve_Cubic(t *testing.T) {
	curve := canvas.NewCubicBezierCurve(fyne.NewPos(0, 60), fyne.NewPos(20, 30), fyne.NewPos(40, 60), fyne.NewPos(50, 10), color.RGBA{R: 0, G: 0, B: 255, A: 255})
	curve.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "bezier_curve_cubic.xml", curve)

	c := software.NewCanvas()
	c.SetContent(curve)
	c.Resize(fyne.NewSize(70, 60))
	test.AssertRendersToImage(t, "bezier_curve_cubic.png", c)

	curve.ControlPoints = []fyne.Position{fyne.NewPos(20, 30), fyne.NewPos(40, 500), fyne.NewPos(60, 60)}
	test.AssertRendersToImage(t, "bezier_curve_cubic.png", c)
}
