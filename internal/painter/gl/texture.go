package gl

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/goki/freetype/truetype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/theme"
)

var noTexture = Texture(cache.NoTexture)

// Texture represents an uploaded GL texture
type Texture cache.TextureType

func (p *painter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := cache.GetTexture(obj)
	if !ok {
		return
	}

	p.ctx.DeleteTexture(Texture(texture))
	p.logError()
	cache.DeleteTexture(obj)
}

func (p *painter) getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
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

func (p *painter) imgToTexture(img image.Image, textureFilter canvas.ImageScale) Texture {
	switch i := img.(type) {
	case *image.Uniform:
		texture := p.newTexture(textureFilter)
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		p.ctx.TexImage2D(
			texture2D,
			0,
			1,
			1,
			colorFormatRGBA,
			unsignedByte,
			data,
		)
		p.logError()
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return noTexture
		}

		texture := p.newTexture(textureFilter)
		p.ctx.TexImage2D(
			texture2D,
			0,
			i.Rect.Size().X,
			i.Rect.Size().Y,
			colorFormatRGBA,
			unsignedByte,
			i.Pix,
		)
		p.logError()
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Over)
		return p.imgToTexture(rgba, textureFilter)
	}
}

func (p *painter) newGlCircleTexture(obj fyne.CanvasObject) Texture {
	circle := obj.(*canvas.Circle)
	raw := paint.DrawCircle(circle, paint.VectorPad(circle), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *painter) newGlImageTexture(obj fyne.CanvasObject) Texture {
	img := obj.(*canvas.Image)

	width := p.textureScale(img.Size().Width)
	height := p.textureScale(img.Size().Height)

	tex := paint.PaintImage(img, p.canvas, int(width), int(height))
	if tex == nil {
		return noTexture
	}

	return p.imgToTexture(tex, img.ScaleMode)
}

func (p *painter) newGlLinearGradientTexture(obj fyne.CanvasObject) Texture {
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

func (p *painter) newGlRadialGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScale(gradient.Size().Width)
	height := p.textureScale(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(int(width), int(height)), canvas.ImageScaleSmooth)
}

func (p *painter) newGlRasterTexture(obj fyne.CanvasObject) Texture {
	rast := obj.(*canvas.Raster)

	width := p.textureScale(rast.Size().Width)
	height := p.textureScale(rast.Size().Height)

	return p.imgToTexture(rast.Generator(int(width), int(height)), rast.ScaleMode)
}

func (p *painter) newGlRectTexture(obj fyne.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	if rect.StrokeColor != nil && rect.StrokeWidth > 0 {
		return p.newGlStrokedRectTexture(rect)
	}
	if rect.FillColor == nil {
		return noTexture
	}
	return p.imgToTexture(image.NewUniform(rect.FillColor), canvas.ImageScaleSmooth)
}

func (p *painter) newGlStrokedRectTexture(obj fyne.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	raw := paint.DrawRectangle(rect, paint.VectorPad(rect), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *painter) newGlTextTexture(obj fyne.CanvasObject) Texture {
	text := obj.(*canvas.Text)
	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	bounds := text.MinSize()
	width := int(p.textureScale(bounds.Width + paint.VectorPad(text))) // potentially italic overspill
	height := int(p.textureScale(bounds.Height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(p.canvas.Scale())
	opts.Size = fontSize
	opts.DPI = float64(paint.TextDPI * p.texScale)
	face := paint.CachedFontFace(text.TextStyle, &opts)

	paint.DrawString(img, text.Text, color, face, height, text.TextStyle.TabWidth)
	return p.imgToTexture(img, canvas.ImageScaleSmooth)
}

func (p *painter) newTexture(textureFilter canvas.ImageScale) Texture {
	if int(textureFilter) >= len(textureFilterToGL) {
		fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", textureFilter), nil)
		textureFilter = canvas.ImageScaleSmooth
	}

	texture := p.ctx.CreateTexture()
	p.logError()
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()
	p.ctx.TexParameteri(texture2D, textureMinFilter, textureFilterToGL[textureFilter])
	p.ctx.TexParameteri(texture2D, textureMagFilter, textureFilterToGL[textureFilter])
	p.ctx.TexParameteri(texture2D, textureWrapS, clampToEdge)
	p.ctx.TexParameteri(texture2D, textureWrapT, clampToEdge)
	p.logError()

	return texture
}

func (p *painter) textureScale(v float32) float32 {
	if p.pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}

	return float32(math.Round(float64(v * p.pixScale)))
}
