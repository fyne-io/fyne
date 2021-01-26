package software_test

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/painter/software"
	internalTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeTestImage(w, h int) image.Image {
	return internalTest.NewCheckedImage(w, h, w, h)
}

func TestPainter_paintCircle(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewCircle(color.Black)

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(fyne.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_circle.png", p.Paint(c))
}

func TestPainter_paintCircleStroke(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewCircle(color.White)
	obj.StrokeColor = color.Black
	obj.StrokeWidth = 4

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(fyne.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_circle_stroke.png", p.Paint(c))
}

func TestPainter_paintGradient_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	g := canvas.NewRadialGradient(color.NRGBA{R: 200, A: 255}, color.NRGBA{B: 200, A: 255})
	g.SetMinSize(fyne.NewSize(100, 100))
	scroll := container.NewScroll(g)
	scroll.Move(fyne.NewPos(10, 10))
	scroll.Resize(fyne.NewSize(50, 50))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-30, -30)})
	cont := fyne.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(fyne.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_gradient_clipped.png", p.Paint(c))
}

func TestPainter_paintImage(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_default.png", target)
}

func TestPainter_paintImage_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	img := canvas.NewImageFromImage(makeTestImage(5, 5))
	img.ScaleMode = canvas.ImageScalePixels
	img.SetMinSize(fyne.NewSize(100, 100))
	scroll := container.NewScroll(img)
	scroll.Move(fyne.NewPos(10, 10))
	scroll.Resize(fyne.NewSize(50, 50))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-15, -15)})
	cont := fyne.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(fyne.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_image_clipped.png", p.Paint(c))
}

func TestPainter_paintImage_scalePixels(t *testing.T) {
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

func TestPainter_paintImage_scaleSmooth(t *testing.T) {
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

func TestPainter_paintImage_scaleFastest(t *testing.T) {
	img := canvas.NewImageFromImage(makeTestImage(3, 3))
	img.ScaleMode = canvas.ImageScaleFastest

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_ImageScaleFastest.png", target)
}

func TestPainter_paintImage_stretchX(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(fyne.NewSize(100, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchx.png", target)
}

func TestPainter_paintImage_stretchY(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(canvas.NewImageFromImage(makeTestImage(3, 3)))
	c.Resize(fyne.NewSize(50, 100))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_image_stretchy.png", target)
}

func TestPainter_paintImage_contain(t *testing.T) {
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

func TestPainter_paintImage_containX(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
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

func TestPainter_paintImage_containY(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
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

func TestPainter_paintLine(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewLine(color.Black)
	obj.StrokeWidth = 6

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(fyne.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_line.png", p.Paint(c))
}

func TestPainter_paintRaster(t *testing.T) {
	img := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		x = x / 5
		y = y / 5
		if x%2 == y%2 {
			return color.White
		} else {
			return color.Black
		}
	})

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSize(50, 50))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_raster.png", target)
}

func TestPainter_paintRaster_scaled(t *testing.T) {
	img := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		x = x / 5
		y = y / 5
		if x%2 == y%2 {
			return color.White
		} else {
			return color.Black
		}
	})

	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.SetScale(5.0)
	c.Resize(fyne.NewSize(5, 5))
	p := software.NewPainter()

	target := p.Paint(c)
	test.AssertImageMatches(t, "draw_raster_scale.png", target)
}

func TestPainter_paintRectangle_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	red1 := canvas.NewRectangle(color.NRGBA{R: 200, A: 255})
	red1.SetMinSize(fyne.NewSize(20, 20))
	red2 := canvas.NewRectangle(color.NRGBA{R: 150, A: 255})
	red2.SetMinSize(fyne.NewSize(20, 20))
	red3 := canvas.NewRectangle(color.NRGBA{R: 100, A: 255})
	red3.SetMinSize(fyne.NewSize(20, 20))
	reds := container.NewHBox(red1, red2, red3)
	green1 := canvas.NewRectangle(color.NRGBA{G: 200, A: 255})
	green1.SetMinSize(fyne.NewSize(20, 20))
	green2 := canvas.NewRectangle(color.NRGBA{G: 150, A: 255})
	green2.SetMinSize(fyne.NewSize(20, 20))
	green3 := canvas.NewRectangle(color.NRGBA{G: 100, A: 255})
	green3.SetMinSize(fyne.NewSize(20, 20))
	greens := container.NewHBox(green1, green2, green3)
	blue1 := canvas.NewRectangle(color.NRGBA{B: 200, A: 255})
	blue1.SetMinSize(fyne.NewSize(20, 20))
	blue2 := canvas.NewRectangle(color.NRGBA{B: 150, A: 255})
	blue2.SetMinSize(fyne.NewSize(20, 20))
	blue3 := canvas.NewRectangle(color.NRGBA{B: 100, A: 255})
	blue3.SetMinSize(fyne.NewSize(20, 20))
	blues := container.NewHBox(blue1, blue2, blue3)
	box := container.NewVBox(reds, greens, blues)
	scroll := container.NewScroll(box)
	scroll.Move(fyne.NewPos(10, 10))
	scroll.Resize(fyne.NewSize(50, 50))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-10, -10)})
	cont := fyne.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(fyne.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_rect_clipped.png", p.Paint(c))
}

func TestPainter_paintRectangle_stroke(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	obj := canvas.NewRectangle(color.Black)
	obj.StrokeWidth = 5
	obj.StrokeColor = &color.RGBA{R: 0xFF, G: 0x33, B: 0x33, A: 0xFF}

	c := test.NewCanvas()
	c.SetPadded(true)
	c.SetContent(obj)
	c.Resize(fyne.NewSize(70+2*theme.Padding(), 70+2*theme.Padding()))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_rectangle_stroke.png", p.Paint(c))
}

func TestPainter_paintText_clipped(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	scroll := container.NewScroll(widget.NewLabel("some text\nis here\nand here"))
	scroll.Move(fyne.NewPos(10, 10))
	scroll.Resize(fyne.NewSize(50, 50))
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(-10, -10)})
	cont := fyne.NewContainer(scroll)
	c := test.NewCanvas()
	c.SetPadded(false)
	c.SetContent(cont)
	c.Resize(fyne.NewSize(70, 70))
	p := software.NewPainter()

	test.AssertImageMatches(t, "draw_text_clipped.png", p.Paint(c))
}
