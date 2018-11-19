// +build !ci,gl

package gl

import (
	"github.com/fyne-io/fyne/canvas"
	"image"
	"image/color"
)

type pixelImage struct {
	source *canvas.Image
}

func (i *pixelImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *pixelImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.source.Size.Width, i.source.Size.Height)
}

func (i *pixelImage) At(x, y int) color.Color {
	if i.source.PixelColor == nil {
		return color.Transparent
	}

	return i.source.PixelColor(x, y, i.source.Size.Width, i.source.Size.Height)
}

func NewPixelImage(source *canvas.Image) image.Image {
	return &pixelImage{source: source}
}
