package gl

import (
	"image"
	"image/color"
	"runtime"

	"fyne.io/fyne"
)

type captureImage struct {
	pix           []uint8
	width, height int
}

func (c *captureImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *captureImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height)
}

func (c *captureImage) At(x, y int) color.Color {
	start := ((c.height-y-1)*c.width + x) * 4
	return color.RGBA{R: c.pix[start], G: c.pix[start+1], B: c.pix[start+2], A: c.pix[start+3]}
}

type glCanvas interface {
	fyne.Canvas
	TextureScale() float32
}

func (p *glPainter) Capture(c fyne.Canvas) image.Image {
	scale := c.Scale()
	if gc, ok := c.(glCanvas); ok && runtime.GOOS == "darwin" { // macOS scaling is done at the texture level
		scale = gc.TextureScale()
	}
	width := int(float32(c.Size().Width) * scale)
	height := int(float32(c.Size().Height) * scale)
	pixels := make([]uint8, width*height*4)

	p.context.RunWithContext(func() {
		p.glCapture(int32(width), int32(height), &pixels)
	})

	return &captureImage{
		pix:    pixels,
		width:  width,
		height: height,
	}
}
