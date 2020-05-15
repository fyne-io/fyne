package painter_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
)

func TestPaintImage_SVG(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(480, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-default.png", c.Capture())
}

func TestPaintImage_SVG_Stretch(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillStretch
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(640, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-stretch.png", c.Capture())
}

func TestPaintImage_SVG_Contain(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillContain
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(640, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-contain.png", c.Capture())
}
