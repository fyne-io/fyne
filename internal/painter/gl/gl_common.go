package gl

import (
	"image"
	"log"
	"runtime"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/theme"
)

var textures = make(map[fyne.CanvasObject]Texture, 1024)

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

func getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) Texture {
	texture, ok := textures[object]

	if !ok {
		texture = creator(object)
		textures[object] = texture
	}
	return texture
}

func (p *glPainter) newGlCircleTexture(obj fyne.CanvasObject) Texture {
	circle := obj.(*canvas.Circle)
	raw := painter.DrawCircle(circle, painter.VectorPad(circle), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlLineTexture(obj fyne.CanvasObject) Texture {
	line := obj.(*canvas.Line)
	raw := painter.DrawLine(line, painter.VectorPad(line), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRectTexture(rect fyne.CanvasObject) Texture {
	col := theme.BackgroundColor()
	if wid, ok := rect.(fyne.Widget); ok {
		widCol := cache.Renderer(wid).BackgroundColor()
		if widCol != nil {
			col = widCol
		}
	} else if rect, ok := rect.(*canvas.Rectangle); ok {
		if rect.StrokeColor != nil && rect.StrokeWidth > 0 {
			return p.newGlStrokedRectTexture(rect)
		}
		if rect.FillColor != nil {
			col = rect.FillColor
		}
	}

	return p.imgToTexture(image.NewUniform(col), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlStrokedRectTexture(obj fyne.CanvasObject) Texture {
	rect := obj.(*canvas.Rectangle)
	raw := painter.DrawRectangle(rect, painter.VectorPad(rect), p.textureScale)

	return p.imgToTexture(raw, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlTextTexture(obj fyne.CanvasObject) Texture {
	text := obj.(*canvas.Text)

	bounds := text.MinSize()
	width := p.textureScaleInt(bounds.Width)
	height := p.textureScaleInt(bounds.Height)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(p.canvas.Scale())
	opts.Size = fontSize
	opts.DPI = float64(painter.TextDPI * p.texScale)
	face := painter.CachedFontFace(text.TextStyle, &opts)

	d := font.Drawer{}
	d.Dst = img
	d.Src = &image.Uniform{C: text.Color}
	d.Face = face
	d.Dot = freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	d.DrawString(text.Text)

	return p.imgToTexture(img, canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlImageTexture(obj fyne.CanvasObject) Texture {
	img := obj.(*canvas.Image)

	width := p.textureScaleInt(img.Size().Width)
	height := p.textureScaleInt(img.Size().Height)

	tex := painter.PaintImage(img, p.canvas, width, height)
	if tex == nil {
		return NoTexture
	}

	return p.imgToTexture(tex, img.ScaleMode)
}

func (p *glPainter) newGlRasterTexture(obj fyne.CanvasObject) Texture {
	rast := obj.(*canvas.Raster)

	width := p.textureScaleInt(rast.Size().Width)
	height := p.textureScaleInt(rast.Size().Height)

	return p.imgToTexture(rast.Generator(width, height), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlLinearGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.LinearGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height), canvas.ImageScaleSmooth)
}

func (p *glPainter) newGlRadialGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height), canvas.ImageScaleSmooth)
}
