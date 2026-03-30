package gl

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
)

type captureImage struct {
	pix           []uint8
	width, height int

	leftPad, rightPad, topPad, bottomPad int
}

func (c *captureImage) SubImage(inner image.Rectangle) *captureImage {
	return &captureImage{
		pix:    c.pix,
		width:  c.width,
		height: c.height,

		leftPad:   inner.Min.X,
		topPad:    inner.Min.Y,
		rightPad:  c.width - inner.Max.X,
		bottomPad: c.height - inner.Max.Y,
	}
}

func (c *captureImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *captureImage) Bounds() image.Rectangle {
	return image.Rect(c.leftPad, c.topPad, c.width-c.rightPad, c.height-c.bottomPad)
}

func (c *captureImage) At(x, y int) color.Color {
	start := ((c.height-y-1)*c.width + x) * 4
	return color.RGBA{R: c.pix[start], G: c.pix[start+1], B: c.pix[start+2], A: c.pix[start+3]}
}

func (p *painter) Capture(c fyne.Canvas) image.Image {
	pos := fyne.NewPos(c.Size().Width, c.Size().Height)
	width, height := c.PixelCoordinateForPosition(pos)
	pixels := make([]uint8, width*height*4)

	p.contextProvider.RunWithContext(func() {
		p.ctx.ReadBuffer(front)
		p.logError()
		p.ctx.ReadPixels(0, 0, width, height, colorFormatRGBA, unsignedByte, pixels)
		p.logError()
	})

	return &captureImage{
		pix:    pixels,
		width:  width,
		height: height,
	}
}
