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

func TestArc_FillColorEndAngle(t *testing.T) {
	c := color.White
	arc := canvas.NewArc(0.0, 360.0, c)

	assert.Equal(t, c, arc.FillColor)
	assert.Equal(t, float32(0.0), arc.InnerRadius)
	assert.Equal(t, float32(0.0), arc.StartAngle)
	assert.Equal(t, float32(360.0), arc.EndAngle)
	assert.Equal(t, float32(0.0), arc.CornerRadius)
}

func TestArc_InnerRadiusStartEndAngles(t *testing.T) {
	arc := &canvas.Arc{
		FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		InnerRadius:  15,
		StartAngle:   90,
		EndAngle:     180,
		CornerRadius: 5,
		StrokeWidth:  2,
		StrokeColor:  color.Black,
	}

	arc.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "arc.xml", arc)

	c := software.NewCanvas()
	c.SetContent(arc)
	c.Resize(fyne.NewSize(60, 60))

	arc.StartAngle = 0
	arc.EndAngle = 360
	test.AssertRendersToImage(t, "arc_full_stroke.png", c)

	arc.StrokeWidth = 0
	test.AssertRendersToImage(t, "arc_full.png", c)

	arc.InnerRadius = 0
	test.AssertRendersToImage(t, "arc_full_no_inner_radius.png", c)

	arc.InnerRadius = 0
	arc.StartAngle = 120
	arc.EndAngle = 220
	test.AssertRendersToImage(t, "arc_120_220_rounded.png", c)

	arc.InnerRadius = 10
	arc.StartAngle = 150
	arc.EndAngle = 350
	test.AssertRendersToImage(t, "arc_150_350_rounded.png", c)

	arc.InnerRadius = 15
	arc.StartAngle = 55
	arc.EndAngle = 350
	arc.StrokeWidth = 1
	test.AssertRendersToImage(t, "arc_55_350_rounded_stroke.png", c)

	arc.CornerRadius = 0
	arc.InnerRadius = 20
	arc.StartAngle = 80
	arc.EndAngle = 110
	arc.StrokeWidth = 0
	test.AssertRendersToImage(t, "arc_80_100.png", c)

	arc.InnerRadius = 0
	arc.StartAngle = -30
	arc.EndAngle = -50
	arc.StrokeWidth = 3
	test.AssertRendersToImage(t, "arc_negative_30_50_stroke.png", c)
}
