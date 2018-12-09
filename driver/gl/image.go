// +build !ci,gl

package gl

import (
	"github.com/fyne-io/fyne/canvas"
	"image"
	"image/color"
)

type pixelImage struct {
	source *canvas.Image
	scale  float32
}

func (i *pixelImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *pixelImage) Bounds() image.Rectangle {
	width := int(float32(i.source.Size().Width) * i.scale)
	height := int(float32(i.source.Size().Height) * i.scale)

	return image.Rect(0, 0, width, height)
}

func (i *pixelImage) At(x, y int) color.Color {
	if i.source.PixelColor == nil {
		return color.Transparent
	}

	width := int(float32(i.source.Size().Width) * i.scale)
	height := int(float32(i.source.Size().Height) * i.scale)

	return i.source.PixelColor(x, y, width, height)
}

func NewPixelImage(source *canvas.Image, scale float32) image.Image {
	return &pixelImage{source: source, scale: scale}
}
