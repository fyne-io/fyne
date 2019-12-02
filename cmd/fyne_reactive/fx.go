package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dataapi"
	"fyne.io/fyne/theme"
	"image/color"
	"math"
)

func newFx(app fyne.App, data *dataModel) *fractal {
	window := app.NewWindow("Fractal")
	window.SetIcon(theme.DocumentCreateIcon())
	window.SetPadded(false)
	fractal := &fractal{data: data, window: window}
	fractal.canvas = canvas.NewRasterWithPixels(fractal.mandelbrot)

	fractal.currIterations = 100
	fractal.currScale = 1.0
	fractal.currX = -0.75
	fractal.currY = 0.0

	window.SetContent(fyne.NewContainerWithLayout(fractal, fractal.canvas))
	window.Canvas().SetOnTypedRune(fractal.fractalRune)
	window.Canvas().SetOnTypedKey(fractal.fractalKey)
	window.Show()

	fractal.dataid = data.DeliveryTime.AddListener(fractal.ChangeTransparency)
	return fractal
}

type fractal struct {
	data                    *dataModel
	currIterations          uint
	currScale, currX, currY float64

	window fyne.Window
	canvas fyne.CanvasObject

	transparency float64
	dataid       int
}

func (f *fractal) ChangeTransparency(value dataapi.DataItem) {
	if vv, ok := value.(*dataapi.Float); ok {
		f.transparency = vv.Value()
		f.currScale = 10 / (f.transparency + 1)
		println("transparency is now", f.transparency)
		f.refresh()
	}
}

func (f *fractal) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	f.canvas.Resize(size)
}

func (f *fractal) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(320, 240)
}

func (f *fractal) refresh() {
	if f.currScale >= 1.0 {
		f.currIterations = 100
	} else {
		f.currIterations = uint(100 * (1 + math.Pow((math.Log10(1/f.currScale)), 1.25)))
	}

	f.window.Canvas().Refresh(f.canvas)
}

func (f *fractal) scaleChannel(c float64, start, end uint32) uint8 {
	if end >= start {
		return (uint8)(c*float64(uint8(end-start))) + uint8(start)
	}

	return (uint8)((1-c)*float64(uint8(start-end))) + uint8(end)
}

func (f *fractal) scaleColor(c float64, start, end color.Color) color.Color {
	r1, g1, b1, _ := start.RGBA()
	r2, g2, b2, _ := end.RGBA()
	return color.RGBA{
		f.scaleChannel(c, r1, r2),
		f.scaleChannel(c, g1, g2),
		f.scaleChannel(c, b1, b2),
		uint8(255 * f.transparency / 100),
	}
}

func (f *fractal) mandelbrot(px, py, w, h int) color.Color {
	drawScale := 3.5 * f.currScale
	aspect := (float64(h) / float64(w))
	cRe := ((float64(px)/float64(w))-0.5)*drawScale + f.currX
	cIm := ((float64(py)/float64(w))-(0.5*aspect))*drawScale - f.currY

	var i uint
	var x, y, xsq, ysq float64

	for i = 0; i < f.currIterations && (xsq+ysq <= 4); i++ {
		xNew := float64(xsq-ysq) + cRe
		y = 2*x*y + cIm
		x = xNew

		xsq = x * x
		ysq = y * y
	}

	if i == f.currIterations {
		return theme.BackgroundColor()
	}

	mu := (float64(i) / float64(f.currIterations))
	c := math.Sin((mu / 2) * math.Pi)

	return f.scaleColor(c, theme.PrimaryColor(), theme.TextColor())
}

func (f *fractal) fractalRune(r rune) {
	if r == '+' {
		f.currScale /= 1.1
	} else if r == '-' {
		f.currScale *= 1.1
	}

	f.transparency = (10 / f.currScale) - 1
	f.data.DeliveryTime.SetFloat(f.transparency, f.dataid)
	f.refresh()
}

func (f *fractal) fractalKey(ev *fyne.KeyEvent) {
	delta := f.currScale * 0.2
	if ev.Name == fyne.KeyUp {
		f.currY -= delta
	} else if ev.Name == fyne.KeyDown {
		f.currY += delta
	} else if ev.Name == fyne.KeyLeft {
		f.currX += delta
	} else if ev.Name == fyne.KeyRight {
		f.currX -= delta
	} else if ev.Name == fyne.KeyEscape {
		f.data.DeliveryTime.DeleteListener(f.dataid)
		println("Deregistering listener")
	}

	f.refresh()
}
