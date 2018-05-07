package examples

import "image/color"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

const maxIterations = 20

type fractalLayout struct {
	canvas ui.CanvasObject
}

func (c *fractalLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	c.canvas.Resize(size)
}

func (c *fractalLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	return ui.NewSize(320, 240)
}

func mandelbrot(px, py, w, h int) color.RGBA {
	c_re := ((float64(px) / float64(w)) - .75) * 3
	c_im := ((float64(py) / float64(h)) - .5) * 2.5

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

		return color.RGBA{(uint8)(mu*199 + 56), 72, (uint8)((1 - mu) * 255), 0xff}
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
