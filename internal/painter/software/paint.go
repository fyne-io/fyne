package software

import (
	"image"
	"image/draw"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/test"
)

type softwarePainter struct {
}

// Paint is the main entry point for a simple software painter.
// The canvas to be drawn is passed in as a parameter and the return is an
// image containing the result of rendering.
func (*softwarePainter) Paint(c fyne.Canvas) image.Image {
	theme := fyne.CurrentApp().Settings().Theme()

	size := c.Size().Union(c.Content().MinSize())
	bounds := image.Rect(0, 0, int(float32(size.Width)*c.Scale()), int(float32(size.Height)*c.Scale()))
	base := image.NewRGBA(bounds)
	draw.Draw(base, bounds, image.NewUniform(theme.BackgroundColor()), image.ZP, draw.Src)

	paint := func(obj fyne.CanvasObject, pos, _ fyne.Position, _ fyne.Size) bool {
		if img, ok := obj.(*canvas.Image); ok {
			drawImage(c, img, pos, size, base)
		} else if text, ok := obj.(*canvas.Text); ok {
			drawText(c, text, pos, size, base)
		} else if rect, ok := obj.(*canvas.Rectangle); ok {
			drawRectangle(c, rect, pos, size, base)
		} else if wid, ok := obj.(fyne.Widget); ok {
			drawWidget(c, wid, pos, size, base)
		}

		return false
	}

	driver.WalkVisibleObjectTree(c.Content(), paint, nil)
	return base
}

// NewPainter creates a new software painter that can paint a canvas in memory
func NewPainter() test.SoftwarePainter {
	return &softwarePainter{}
}
