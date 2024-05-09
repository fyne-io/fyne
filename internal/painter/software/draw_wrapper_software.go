//go:build software

package software

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
)

func (p *Painter) drawCircleWrapper(circle *canvas.Circle, pad float32, scale func(float32) float32) *image.RGBA {
	raw, _ := p.getTexture(circle, func(object fyne.CanvasObject) Texture {
		return painter.DrawCircle(circle, pad, scale)
	})
	return raw
}

func (p *Painter) paintImageWrapper(img *canvas.Image, c fyne.Canvas, width, height int) image.Image {
	raw, _ := p.getTexture(img, func(object fyne.CanvasObject) Texture {
		return painter.PaintImage(img, c, width, height)
	})
	return raw
}

func (p *Painter) drawLineWrapper(line *canvas.Line, pad float32, scale func(float32) float32) *image.RGBA {
	raw, _ := p.getTexture(line, func(object fyne.CanvasObject) Texture {
		return painter.DrawLine(line, pad, scale)
	})
	return raw
}

func (p *Painter) drawStringWrapper(c fyne.Canvas, text *canvas.Text, width, height int, color color.Color) *image.RGBA {
	txtImg, _ := p.getTexture(text, func(object fyne.CanvasObject) Texture {
		txtImg := image.NewRGBA(image.Rect(0, 0, width, height))
		face := painter.CachedFontFace(text.TextStyle, text.TextSize*c.Scale(), 1)
		painter.DrawString(txtImg, text.Text, color, face.Fonts, text.TextSize, c.Scale(), text.TextStyle.TabWidth)
		return txtImg
	})
	return txtImg
}

func (p *Painter) drawRasterWrapper(raster *canvas.Raster, width, height int) image.Image {
	raw, _ := p.getTexture(raster, func(object fyne.CanvasObject) Texture {
		return raster.Generator(width, height)
	})
	return raw
}

func (p *Painter) drawRectangleStrokeWrapper(rect *canvas.Rectangle, pad float32, scale func(float32) float32) *image.RGBA {
	raw, _ := p.getTexture(rect, func(object fyne.CanvasObject) Texture {
		return painter.DrawRectangle(rect, pad, scale)
	})
	return raw
}
