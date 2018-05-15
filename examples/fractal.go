package examples

import "math"
import "image/color"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

var currIterations uint = 100

var currScale = 1.0
var currX = -0.75
var currY = 0.0

var fractalWindow ui.Window
var fractalCanvas ui.CanvasObject

type fractalLayout struct {
	canvas ui.CanvasObject
}

func (c *fractalLayout) Layout(objects []ui.CanvasObject, size ui.Size) {
	c.canvas.Resize(size)
}

func (c *fractalLayout) MinSize(objects []ui.CanvasObject) ui.Size {
	return ui.NewSize(320, 240)
}

func refresh() {
	currIterations = uint(100 * (1 + math.Pow((math.Log10(1/currScale)), 1.25)))

	fractalWindow.Canvas().Refresh(fractalCanvas)
}

func scaleColor(c float64, start, end uint8) uint8 {
	if end >= start {
		return (uint8)(c*float64(end-start)) + start
	} else {
		return (uint8)((1-c)*float64(start-end)) + end
	}
}

func mandelbrot(px, py, w, h int) color.RGBA {
	drawScale := 3.5 * currScale
	aspect := (float64(h) / float64(w))
	c_re := ((float64(px)/float64(w))-0.5)*drawScale + currX
	c_im := ((float64(py)/float64(w))-(0.5*aspect))*drawScale - currY

	var i uint
	var x, y, xsq, ysq float64

	for i = 0; i < currIterations && (xsq+ysq <= 4); i++ {
		x_new := float64(xsq-ysq) + c_re
		y = 2*x*y + c_im
		x = x_new

		xsq = x * x
		ysq = y * y
	}

	if i == currIterations {
		return theme.BackgroundColor()
	} else {
		mu := (float64(i) / float64(currIterations))
		c := math.Sin((mu / 2) * math.Pi)

		return color.RGBA{scaleColor(c, theme.PrimaryColor().R, theme.TextColor().R),
			scaleColor(c, theme.PrimaryColor().G, theme.TextColor().G),
			scaleColor(c, theme.PrimaryColor().B, theme.TextColor().B), 0xff}
	}
}

func fractalKeyDown(ev *ui.KeyEvent) {
	delta := currScale * 0.2
	if ev.Name == "Up" {
		currY -= delta
	} else if ev.Name == "Down" {
		currY += delta
	} else if ev.Name == "Left" {
		currX += delta
	} else if ev.Name == "Right" {
		currX -= delta
	} else if ev.String == "+" {
		currScale /= 1.1
	} else if ev.String == "-" {
		currScale *= 1.1
	}

	refresh()
}
func Fractal(app app.App) {
	fractalWindow = app.NewWindow("Fractal")
	fractalCanvas = canvas.NewRaster(mandelbrot)

	container := ui.NewContainer(fractalCanvas)
	container.Layout = &fractalLayout{canvas: fractalCanvas}
	fractalWindow.Canvas().SetContent(container)
	fractalWindow.Canvas().SetOnKeyDown(fractalKeyDown)
	fractalWindow.Show()
}
