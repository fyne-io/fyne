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

func TestPaintImage_SVG_StretchX(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillStretch
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(640, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-stretchx.png", c.Capture())
}

func TestPaintImage_SVG_StretchY(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillStretch
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(480, 480))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-stretchy.png", c.Capture())
}

func TestPaintImage_SVG_ContainX(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillContain
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(640, 240))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-containx.png", c.Capture())
}

func TestPaintImage_SVG_ContainY(t *testing.T) {
	img := canvas.NewImageFromFile("testdata/stroke.svg")
	img.FillMode = canvas.ImageFillContain
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetContent(img)
	c.Resize(fyne.NewSize(480, 480))
	img.Refresh()

	test.AssertImageMatches(t, "svg-stroke-containy.png", c.Capture())
}
