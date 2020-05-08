package software

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/test"
)

type painter struct {
}

// NewPainter creates a new software painter that can paint a canvas in memory
func NewPainter() test.SoftwarePainter {
	return &painter{}
}

// Paint is the main entry point for a simple software painter.
// The canvas to be drawn is passed in as a parameter and the return is an
// image containing the result of rendering.
func (*painter) Paint(c fyne.Canvas) image.Image {
	theme := fyne.CurrentApp().Settings().Theme()

	size := c.Size().Max(c.Content().MinSize())
	bounds := image.Rect(0, 0, int(float32(size.Width)*c.Scale()), int(float32(size.Height)*c.Scale()))
	base := image.NewNRGBA(bounds)
	draw.Draw(base, bounds, image.NewUniform(theme.BackgroundColor()), image.ZP, draw.Src)

	paint := func(obj fyne.CanvasObject, pos, _ fyne.Position, _ fyne.Size) bool {
		switch o := obj.(type) {
		case *canvas.Image:
			drawImage(c, o, pos, base)
		case *canvas.Text:
			drawText(c, o, pos, base)
		case gradient:
			drawGradient(c, o, pos, base)
		case *canvas.Rectangle:
			drawRectangle(c, o, pos, base)
		case fyne.Widget:
			drawWidget(c, o, pos, base)
		}

		return false
	}

	driver.WalkVisibleObjectTree(c.Content(), paint, nil)
	for _, o := range c.Overlays().List() {
		driver.WalkVisibleObjectTree(o, paint, nil)
	}

	return base
}
