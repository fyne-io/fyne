package apps

import "math"
import "image/color"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type fractal struct {
	currIterations          uint
	currScale, currX, currY float64

	window fyne.Window
	canvas fyne.CanvasObject
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

func (f *fractal) scaleColor(c float64, start, end uint8) uint8 {
	if end >= start {
		return (uint8)(c*float64(end-start)) + start
	}

	return (uint8)((1-c)*float64(start-end)) + end
}

func (f *fractal) mandelbrot(px, py, w, h int) color.RGBA {
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

	return color.RGBA{f.scaleColor(c, theme.PrimaryColor().R, theme.TextColor().R),
		f.scaleColor(c, theme.PrimaryColor().G, theme.TextColor().G),
		f.scaleColor(c, theme.PrimaryColor().B, theme.TextColor().B), 0xff}
}

func (f *fractal) fractalKeyDown(ev *fyne.KeyEvent) {
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

// Fractal loads a Mandelbrot fractal example window for the specified app context
func Fractal(app fyne.App) {
	window := app.NewWindow("Fractal")
	fractal := &fractal{window: window}
	fractal.canvas = canvas.NewRaster(fractal.mandelbrot)

	fractal.currIterations = 100
	fractal.currScale = 1.0
	fractal.currX = -0.75
	fractal.currY = 0.0

	window.SetContent(fyne.NewContainerWithLayout(fractal, fractal.canvas))
	window.Canvas().SetOnKeyDown(fractal.fractalKeyDown)
	window.Show()
}
