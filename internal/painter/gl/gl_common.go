//go:generate fyne bundle -o shaders.go --prefix shader --package gl shaders/

package gl

import (
	"fmt"
	"image"
	"log"
	"runtime"

	"github.com/goki/freetype/truetype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/theme"
)

// Texture represents an uploaded GL texture
type Texture cache.TextureType

var noTexture = Texture(cache.NoTexture)

func logGLError(err uint32) {
	if err == 0 {
		return
	}

	log.Printf("Error %x in GL Renderer", err)
	_, file, line, ok := runtime.Caller(2)
	if ok {
		log.Printf("  At: %s:%d", file, line)
	}
}

func (p *glPainter) getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
	texture, ok := cache.GetTexture(object)

	if !ok {
		texture = cache.TextureType(creator(object))
		cache.SetTexture(object, texture, p.canvas)
	}
	if !cache.IsValid(texture) {
		return noTexture, fmt.Errorf("no texture available")
	}
	return Texture(texture), nil
}

func (p *glPainter) newGlCircleTexture(obj fyne.CanvasObject) Texture {
	circle := obj.(*canvas.Circle)
	raw := painter.DrawCircle(circle, painter.VectorPad(circle), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRectTexture(obj fyne.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	if rect.StrokeColor != nil && rect.StrokeWidth > 0 {
		return p.newGlStrokedRectTexture(rect)
	}
	if rect.FillColor == nil {
		return noTexture
	}
	return p.imgToTexture(image.NewUniform(rect.FillColor), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlStrokedRectTexture(obj fyne.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	raw := painter.DrawRectangle(rect, painter.VectorPad(rect), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlTextTexture(obj fyne.CanvasObject) Texture {
	text := obj.(*canvas.Text)
	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	bounds := text.MinSize()
	width := int(p.textureScale(bounds.Width))
	height := int(p.textureScale(bounds.Height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(p.canvas.Scale())
	opts.Size = fontSize
	opts.DPI = float64(painter.TextDPI * p.texScale)
	face := painter.CachedFontFace(text.TextStyle, &opts)

	painter.DrawString(img, text.Text, color, face, height, text.TextStyle.TabWidth)
	return p.imgToTexture(img, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlImageTexture(obj fyne.CanvasObject) Texture {
	img := obj.(*canvas.Image)

	width := p.textureScale(img.Size().Width)
	height := p.textureScale(img.Size().Height)

	tex := painter.PaintImage(img, p.canvas, int(width), int(height))
	if tex == nil {
		return noTexture
	}

	return p.imgToTexture(tex, img.ScaleMode)
}

func (p *glPainter) newGlRasterTexture(obj fyne.CanvasObject) Texture {
	rast := obj.(*canvas.Raster)

	width := p.textureScale(rast.Size().Width)
	height := p.textureScale(rast.Size().Height)

	return p.imgToTexture(rast.Generator(int(width), int(height)), rast.ScaleMode)
}

func (p *glPainter) newGlLinearGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.LinearGradient)

	w := gradient.Size().Width
	h := gradient.Size().Height
	switch gradient.Angle {
	case 90, 270:
		h = 1
	case 0, 180:
		w = 1
	}
	width := p.textureScale(w)
	height := p.textureScale(h)

	return p.imgToTexture(gradient.Generate(int(width), int(height)), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRadialGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScale(gradient.Size().Width)
	height := p.textureScale(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(int(width), int(height)), canvas.ImageScaleSmooth)
}
