package software

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
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

func TestDrawImage(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage())
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestDrawImage_StretchX(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	img := canvas.NewImageFromImage(makeTestImage())
	img.Resize(fyne.NewSize(100, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(100, 50), target)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestDrawImage_StretchY(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 100))
	img := canvas.NewImageFromImage(makeTestImage())
	img.Resize(fyne.NewSize(50, 100))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 100), target)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestDrawImage_Contain(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestDrawImage_ContainX(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(100, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(100, 50), target)
	test.AssertImageMatches(t, "draw_image_containx.png", target)
}

func TestDrawImage_ContainY(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 100))
	img := canvas.NewImageFromImage(makeTestImage())
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(50, 100))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 100), target)
	test.AssertImageMatches(t, "draw_image_containy.png", target)
}
