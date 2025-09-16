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

func TestArc_FillColorStartAngleEndAngle(t *testing.T) {
	c := color.White
	arc := canvas.NewArc(35.5, 192.7, 0.5, c)

	assert.Equal(t, c, arc.FillColor)
	assert.Equal(t, float32(0.5), arc.CutoutRatio)
	assert.Equal(t, float32(35.5), arc.StartAngle)
	assert.Equal(t, float32(192.7), arc.EndAngle)
	assert.Equal(t, float32(0.0), arc.CornerRadius)
}

func TestArc_Pie(t *testing.T) {
	c := color.White
	arc := canvas.NewPieArc(0.0, 260.0, c)

	assert.Equal(t, c, arc.FillColor)
	assert.Equal(t, float32(0.0), arc.CutoutRatio)
	assert.Equal(t, float32(0.0), arc.StartAngle)
	assert.Equal(t, float32(260.0), arc.EndAngle)
	assert.Equal(t, float32(0.0), arc.CornerRadius)
}

func TestArc_Doughnut(t *testing.T) {
	c := color.White
	arc := canvas.NewDoughnutArc(15.0, 192.0, c)

	assert.Equal(t, c, arc.FillColor)
	assert.Equal(t, float32(0.5), arc.CutoutRatio)
	assert.Equal(t, float32(15.0), arc.StartAngle)
	assert.Equal(t, float32(192.0), arc.EndAngle)
	assert.Equal(t, float32(0.0), arc.CornerRadius)
}

func TestArc_InnerRadiusStartEndAngles(t *testing.T) {
	arc := &canvas.Arc{
		FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		CutoutRatio:  0.5,
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

	arc.CutoutRatio = 0
	test.AssertRendersToImage(t, "arc_full_no_inner_radius.png", c)

	arc.StartAngle = 120
	arc.EndAngle = 220
	test.AssertRendersToImage(t, "arc_120_220_rounded.png", c)

	arc.CutoutRatio = 0.38
	arc.StartAngle = 150
	arc.EndAngle = 350
	test.AssertRendersToImage(t, "arc_150_350_rounded.png", c)

	arc.CutoutRatio = 0.57
	arc.StartAngle = 55
	arc.EndAngle = 350
	arc.StrokeWidth = 1
	test.AssertRendersToImage(t, "arc_55_350_rounded_stroke.png", c)

	arc.CornerRadius = 0
	arc.CutoutRatio = 0.77
	arc.StartAngle = 80
	arc.EndAngle = 110
	arc.StrokeWidth = 0
	test.AssertRendersToImage(t, "arc_80_100.png", c)

	arc.CutoutRatio = 0
	arc.StartAngle = -30
	arc.EndAngle = -50
	arc.StrokeWidth = 3
	test.AssertRendersToImage(t, "arc_negative_30_50_stroke.png", c)
}

func TestArc_RadiusMaximum(t *testing.T) {
	arc := &canvas.Arc{
		FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
		CutoutRatio:  0.5,
		StartAngle:   80,
		EndAngle:     220,
		CornerRadius: canvas.RadiusMaximum,
		StrokeWidth:  2,
		StrokeColor:  color.Black,
	}

	arc.Resize(fyne.NewSize(50, 50))
	test.AssertObjectRendersToMarkup(t, "maximum_rounded_arc.xml", arc)

	c := software.NewCanvas()
	c.SetContent(arc)
	c.Resize(fyne.NewSize(60, 60))

	test.AssertRendersToImage(t, "maximum_rounded_arc_80_220_stroke.png", c)

	arc.StrokeWidth = 0
	test.AssertRendersToImage(t, "maximum_rounded_arc_80_220.png", c)

	arc.CutoutRatio = 0
	test.AssertRendersToImage(t, "maximum_rounded_arc_80_220_no_inner_radius.png", c)

	arc.EndAngle = 110
	arc.CutoutRatio = 0.38
	test.AssertRendersToImage(t, "maximum_rounded_arc_80_110.png", c)

	arc.CutoutRatio = 0
	test.AssertRendersToImage(t, "maximum_rounded_arc_80_110_no_inner_radius.png", c)

	arc.StartAngle = 0
	arc.EndAngle = 350
	arc.CutoutRatio = 0.77
	test.AssertRendersToImage(t, "maximum_rounded_arc_0_350.png", c)

	arc.CutoutRatio = 0
	test.AssertRendersToImage(t, "maximum_rounded_arc_0_350_no_inner_radius.png", c)
}
