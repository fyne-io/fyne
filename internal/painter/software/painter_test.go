package software_test

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func makeTestImage(w, h int) image.Image {
	src := image.NewNRGBA(image.Rect(0, 0, w, h))
	colors := []color.Color{color.White, color.Black}
	c := 0
	for y := 0; y < h; y++ {
		c = y % 2
		for x := 0; x < w; x++ {
			src.Set(x, y, colors[c])
			c = 1 - c
		}
	}
	return src
}

func TestPainter_PaintImage(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestPainter_PaintImageScalePixels(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScalePixels.png", target)
}

func TestPainter_PaintImageScaleSmooth(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScaleSmooth

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScaleSmooth.png", target)
}

func TestPainter_PaintImage_StretchX(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(fyne.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestPainter_PaintImage_StretchY(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(fyne.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestPainter_PaintImage_Contain(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_contain.png", target)
}

func TestPainter_PaintImage_ContainX(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	img := canvas.NewImageFromImage(makeTestImage(3, 4))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_containx.png", target)
}

func TestPainter_PaintImage_ContainY(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	img := canvas.NewImageFromImage(makeTestImage(4, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_containy.png", target)
}
