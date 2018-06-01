package examples

import "math"
import "image/color"

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/theme"

type fractal struct {
	currIterations          uint
	currScale, currX, currY float64

	window ui.Window
	canvas ui.CanvasObject
}

func (f *fractal) Layout(objects []ui.CanvasObject, size ui.Size) {
	f.canvas.Resize(size)
}

func (f *fractal) MinSize(objects []ui.CanvasObject) ui.Size {
	return ui.NewSize(320, 240)
}

func (f *fractal) refresh() {
	if f.currScale >= 1.0 {
		f.currIterations = 100
	} else {
		f.currIterations = uint(100 * (1 + math.Pow((math.Log10(1/f.currScale)), 1.25)))
	}

	f.window.Canvas().Refresh(f.canvas)
}

func (f *fractal) scaleColor(c float64, start, end uint8) uint8 {
	if end >= start {
		return (uint8)(c*float64(end-start)) + start
	} else {
		return (uint8)((1-c)*float64(start-end)) + end
	}
}

func (f *fractal) mandelbrot(px, py, w, h int) color.RGBA {
	drawScale := 3.5 * f.currScale
	aspect := (float64(h) / float64(w))
	c_re := ((float64(px)/float64(w))-0.5)*drawScale + f.currX
	c_im := ((float64(py)/float64(w))-(0.5*aspect))*drawScale - f.currY

	var i uint
	var x, y, xsq, ysq float64

	for i = 0; i < f.currIterations && (xsq+ysq <= 4); i++ {
		x_new := float64(xsq-ysq) + c_re
		y = 2*x*y + c_im
		x = x_new

		xsq = x * x
		ysq = y * y
	}

	if i == f.currIterations {
		return theme.BackgroundColor()
	} else {
		mu := (float64(i) / float64(f.currIterations))
		c := math.Sin((mu / 2) * math.Pi)

		return color.RGBA{f.scaleColor(c, theme.PrimaryColor().R, theme.TextColor().R),
			f.scaleColor(c, theme.PrimaryColor().G, theme.TextColor().G),
			f.scaleColor(c, theme.PrimaryColor().B, theme.TextColor().B), 0xff}
	}
}

func (f *fractal) fractalKeyDown(ev *ui.KeyEvent) {
	delta := f.currScale * 0.2
	if ev.Name == "Up" {
		f.currY -= delta
	} else if ev.Name == "Down" {
		f.currY += delta
	} else if ev.Name == "Left" {
		f.currX += delta
	} else if ev.Name == "Right" {
		f.currX -= delta
	} else if ev.String == "+" {
		f.currScale /= 1.1
	} else if ev.String == "-" {
		f.currScale *= 1.1
	}

	f.refresh()
}

func Fractal(app app.App) {
	window := app.NewWindow("Fractal")
	fractal := &fractal{window: window}
	fractal.canvas = canvas.NewRaster(fractal.mandelbrot)

	fractal.currIterations = 100
	fractal.currScale = 1.0
	fractal.currX = -0.75
	fractal.currY = 0.0

	container := ui.NewContainer(fractal.canvas)
	container.Layout = fractal

	window.Canvas().SetContent(container)
	window.Canvas().SetOnKeyDown(fractal.fractalKeyDown)
	window.Show()
}
