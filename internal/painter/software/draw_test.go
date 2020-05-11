package software

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
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

func TestDrawImage(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestDrawImageNearestFilter(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_nearest_default.png", target)
}

func TestDrawImageLinearFilter(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScalingFilter = canvas.LinearFilter // default value
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_linear_default.png", target)
}

func TestDrawImage_StretchX(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(100, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(100, 50), target)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestDrawImage_StretchY(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 100))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(50, 100))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 100), target)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestDrawImage_Contain(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.FillMode = canvas.ImageFillContain
	img.Resize(fyne.NewSize(50, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 50), target)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestDrawImage_ContainHV(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	img := canvas.NewImageFromImage(makeTestImage(3, 5))
	img.FillMode = canvas.ImageFillContain
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(100, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(100, 50), target)
	test.AssertImageMatches(t, "draw_image_contain_HV.png", target)
}
func TestDrawImage_ContainHH(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 100, 50))
	img := canvas.NewImageFromImage(makeTestImage(9, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(100, 50))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(100, 50), target)
	test.AssertImageMatches(t, "draw_image_contain_HH.png", target)
}

func TestDrawImage_ContainVH(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 100))
	img := canvas.NewImageFromImage(makeTestImage(5, 3))
	img.FillMode = canvas.ImageFillContain
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(50, 100))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 100), target)
	test.AssertImageMatches(t, "draw_image_contain_VH.png", target)
}

func TestDrawImage_ContainVV(t *testing.T) {
	target := image.NewNRGBA(image.Rect(0, 0, 50, 100))
	img := canvas.NewImageFromImage(makeTestImage(3, 9))
	img.FillMode = canvas.ImageFillContain
	img.ScalingFilter = canvas.NearestFilter
	img.Resize(fyne.NewSize(50, 100))

	drawImage(test.Canvas(), img, fyne.NewPos(0, 0), fyne.NewSize(50, 100), target)
	test.AssertImageMatches(t, "draw_image_contain_VV.png", target)
}
