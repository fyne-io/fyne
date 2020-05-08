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

func makeTestImage() image.Image {
	src := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if x == y {
				src.Set(x, y, color.Black)
			} else {
				src.Set(x, y, color.White)
			}
		}
	}

	return src
}

func TestPainter_PaintImage(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage()))
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestPainter_PaintImage_StretchX(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage()))
	c.Resize(fyne.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestPainter_PaintImage_StretchY(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage()))
	c.Resize(fyne.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestPainter_PaintImage_Contain(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestPainter_PaintImage_ContainX(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
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
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_containy.png", target)
}
