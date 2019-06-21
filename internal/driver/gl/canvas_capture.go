package gl

import (
	"image"
	"image/color"

	"github.com/go-gl/gl/v3.2-core/gl"
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

func (c *glCanvas) Capture() image.Image {
	width := int(float32(c.Size().Width) * c.scale)
	height := int(float32(c.Size().Height) * c.scale)
	pixels := make([]uint8, width*height*4)

	runOnMain(func() {
		c.context.runWithContext(func() {
			gl.ReadBuffer(gl.FRONT)
			gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

		})
	})

	return &captureImage{
		pix:    pixels,
		width:  width,
		height: height,
	}
}
