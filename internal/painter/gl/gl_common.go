package gl

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/theme"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var textures = make(map[fyne.CanvasObject]Texture, 1024)

const vectorPad = 10

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
	radius := fyne.Min(circle.Size().Width, circle.Size().Height) / 2

	width := p.textureScaleInt(circle.Size().Width + vectorPad*2)
	height := p.textureScaleInt(circle.Size().Height + vectorPad*2)
	stroke := circle.StrokeWidth * p.canvas.Scale() * p.texScale

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(circle.Size().Width, circle.Size().Height, raw, raw.Bounds())

	if circle.FillColor != nil {
		filler := rasterx.NewFiller(width, height, scanner)
		filler.SetColor(circle.FillColor)
		rasterx.AddCircle(float64(width/2), float64(height/2), float64(p.textureScaleInt(radius)), filler)
		filler.Draw()
	}

	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(circle.StrokeColor)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	rasterx.AddCircle(float64(width/2), float64(height/2), float64(p.textureScaleInt(radius)), dasher)
	dasher.Draw()

	return p.imgToTexture(raw)
}

func (p *glPainter) newGlLineTexture(obj fyne.CanvasObject) Texture {
	line := obj.(*canvas.Line)

	col := line.StrokeColor
	width := p.textureScaleInt(line.Size().Width + vectorPad*2)
	height := p.textureScaleInt(line.Size().Height + vectorPad*2)
	stroke := line.StrokeWidth * p.canvas.Scale() * p.texScale

	raw := image.NewRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(line.Size().Width, line.Size().Height, raw, raw.Bounds())
	dasher := rasterx.NewDasher(width, height, scanner)
	dasher.SetColor(col)
	dasher.SetStroke(fixed.Int26_6(float64(stroke)*64), 0, nil, nil, nil, 0, nil, 0)
	p1x, p1y := p.textureScaleInt(line.Position1.X-line.Position().X+vectorPad), p.textureScaleInt(line.Position1.Y-line.Position().Y+vectorPad)
	p2x, p2y := p.textureScaleInt(line.Position2.X-line.Position().X+vectorPad), p.textureScaleInt(line.Position2.Y-line.Position().Y+vectorPad)

	dasher.Start(rasterx.ToFixedP(float64(p1x), float64(p1y)))
	dasher.Line(rasterx.ToFixedP(float64(p2x), float64(p2y)))
	dasher.Stop(true)
	dasher.Draw()

	return p.imgToTexture(raw)
}

func (p *glPainter) newGlRectTexture(rect fyne.CanvasObject) Texture {
	col := theme.BackgroundColor()
	if wid, ok := rect.(fyne.Widget); ok {
		widCol := cache.Renderer(wid).BackgroundColor()
		if widCol != nil {
			col = widCol
		}
	} else if rect, ok := rect.(*canvas.Rectangle); ok {
		if rect.FillColor != nil {
			col = rect.FillColor
		}
	}

	return p.imgToTexture(image.NewUniform(col))
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

	return p.imgToTexture(img)
}

func (p *glPainter) newGlImageTexture(obj fyne.CanvasObject) Texture {
	img := obj.(*canvas.Image)

	width := p.textureScaleInt(img.Size().Width)
	height := p.textureScaleInt(img.Size().Height)

	tex := painter.PaintImage(img, p.canvas, width, height)
	if tex == nil {
		return NoTexture
	}

	return p.imgToTexture(tex)
}

func (p *glPainter) newGlRasterTexture(obj fyne.CanvasObject) Texture {
	rast := obj.(*canvas.Raster)

	width := p.textureScaleInt(rast.Size().Width)
	height := p.textureScaleInt(rast.Size().Height)

	return p.imgToTexture(rast.Generator(width, height))
}

func (p *glPainter) newGlLinearGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.LinearGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height))
}

func (p *glPainter) newGlRadialGradientTexture(obj fyne.CanvasObject) Texture {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height))
}
