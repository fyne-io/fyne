package painter_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestPaintImage_MinSize(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/svg-stroke-default.png")
	img.FillMode = canvas.ImageFillOriginal
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetScale(1.0)
	c.SetContent(img)

	// check fallback min
	assert.Equal(t, float32(1), img.MinSize().Width)
	assert.Equal(t, float32(1), img.MinSize().Height)

	// render it with original will set min
	painter.PaintImage(img, c, 100, 100)
	assert.Equal(t, float32(480), img.MinSize().Width)
	assert.Equal(t, float32(240), img.MinSize().Height)
}

func TestPaintImageWithBadSVGElement(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke-bad-element-data.svg")
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(480, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-default.png", c.Capture())
}

func TestPaintImage_SVG(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		width     float32
		height    float32
		fillMode  canvas.ImageFill
		wantImage string
	}{
		"default": {
			width:  480,
			height: 240,
		},
		"stretchx": {
			width:    640,
			height:   240,
			fillMode: canvas.ImageFillStretch,
		},
		"stretchy": {
			width:    480,
			height:   480,
			fillMode: canvas.ImageFillStretch,
		},
		"containx": {
			width:    640,
			height:   240,
			fillMode: canvas.ImageFillContain,
		},
		"containy": {
			width:    480,
			height:   480,
			fillMode: canvas.ImageFillContain,
		},
	} {
		t.Run(name, func(t *testing.T) {
			img := canvas.NewImageFromFile("testdata/stroke.svg")
			c := test.NewCanvasWithPainter(software.NewPainter())
			c.SetContent(img)
			c.Resize(fyne.NewSize(tt.width, tt.height))
			img.Refresh()
			img.FillMode = tt.fillMode

			test.AssertImageMatches(t, "svg-stroke-"+name+".png", c.Capture())
		})
	}
}
