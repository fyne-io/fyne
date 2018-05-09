package examples

import "math"
import "image/color"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

const maxIterations = 100
const currScale = 3.2
const currX = -0.75
const currY = 0.0

type fractalLayout struct {
	canvas ui.CanvasObject
}

func (c *fractalLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	c.canvas.Resize(size)
}

func (c *fractalLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	return ui.NewSize(320, 240)
}

func scaleColor(c float64, start, end uint8) uint8 {
	if end >= start {
		return (uint8)(c*float64(end-start)) + start
	} else {
		return (uint8)((1-c)*float64(start-end)) + end
	}
}

func mandelbrot(px, py, w, h int) color.RGBA {
	aspect := (float64(h) / float64(w))
	c_re := ((float64(px)/float64(w))-0.5)*currScale + currX
	c_im := ((float64(py)/float64(w))-(0.5*aspect))*currScale - currY

	var i int
	var x, y, xsq, ysq float64

	for i = 0; i < maxIterations && (xsq+ysq <= 4); i++ {
		x_new := float64(xsq-ysq) + c_re
		y = 2*x*y + c_im
		x = x_new

		xsq = x * x
		ysq = y * y
	}

	if i == maxIterations {
		return theme.BackgroundColor()
	} else {
		mu := (float64(i) / float64(maxIterations))
		c := math.Sin((mu / 2) * math.Pi)

		return color.RGBA{scaleColor(c, theme.PrimaryColor().R, theme.TextColor().R),
			scaleColor(c, theme.PrimaryColor().G, theme.TextColor().G),
			scaleColor(c, theme.PrimaryColor().B, theme.TextColor().B), 0xff}
	}
}

func Fractal(app app.App) {
	fractalWindow := app.NewWindow("Fractal")
	fractal := &fractalLayout{canvas: canvas.NewRaster(mandelbrot)}

	container := ui.NewContainer(fractal.canvas)
	container.Layout = fractal
	fractalWindow.Canvas().SetContent(container)
	fractalWindow.Show()
}
