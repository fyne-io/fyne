package gl

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg" // avoid users having to import when using image widget
	_ "image/png"  // avoid the same for PNG images
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var textures = make(map[fyne.CanvasObject]uint32)

const vectorPad = 10

func getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) uint32) uint32 {
	texture := textures[object]

	if texture == 0 {
		texture = creator(object)
		textures[object] = texture
	}
	return texture
}

func (p *glPainter) newGlCircleTexture(obj fyne.CanvasObject) uint32 {
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

func (p *glPainter) newGlLineTexture(obj fyne.CanvasObject) uint32 {
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

func (p *glPainter) newGlRectTexture(rect fyne.CanvasObject) uint32 {
	col := theme.BackgroundColor()
	if wid, ok := rect.(fyne.Widget); ok {
		widCol := widget.Renderer(wid).BackgroundColor()
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

func (p *glPainter) newGlTextTexture(obj fyne.CanvasObject) uint32 {
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

func (p *glPainter) newGlImageTexture(obj fyne.CanvasObject) uint32 {
	img := obj.(*canvas.Image)

	width := p.textureScaleInt(img.Size().Width)
	height := p.textureScaleInt(img.Size().Height)
	if width <= 0 || height <= 0 {
		return 0
	}

	switch {
	case img.File != "" || img.Resource != nil:
		var file io.Reader
		var name string
		if img.Resource != nil {
			name = img.Resource.Name()
			file = bytes.NewReader(img.Resource.Content())
		} else {
			name = img.File
			handle, _ := os.Open(img.File)
			defer handle.Close()
			file = handle
		}

		if strings.ToLower(filepath.Ext(name)) == ".svg" {
			tex := svgCacheGet(img.Resource, width, height)
			if tex == nil {
				// Not in cache, so load the item and add to cache

				icon, err := oksvg.ReadIconStream(file)
				if err != nil {
					fyne.LogError("SVG Load error:", err)

					return 0
				}
				icon.SetTarget(0, 0, float64(width), float64(height))

				w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
				// this is used by our render code, so let's set it to the file aspect
				aspects[img.Resource] = float32(w) / float32(h)
				// if the image specifies it should be original size we need at least that many pixels on screen
				if img.FillMode == canvas.ImageFillOriginal {
					p.checkImageMinSize(img, w, h)
				}

				tex = image.NewRGBA(image.Rect(0, 0, width, height))
				scanner := rasterx.NewScannerGV(w, h, tex, tex.Bounds())
				raster := rasterx.NewDasher(width, height, scanner)

				icon.Draw(raster, 1)
				svgCachePut(img.Resource, tex, width, height)
			}

			return p.imgToTexture(tex)
		}

		pixels, _, err := image.Decode(file)

		if err != nil {
			fyne.LogError("image err", err)

			return 0
		}
		origSize := pixels.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			p.checkImageMinSize(img, origSize.X, origSize.Y)
		}

		tex := image.NewRGBA(pixels.Bounds())
		draw.Draw(tex, pixels.Bounds(), pixels, image.ZP, draw.Src)

		return p.imgToTexture(tex)
	case img.Image != nil:
		origSize := img.Image.Bounds().Size()
		// this is used by our render code, so let's set it to the file aspect
		aspects[img] = float32(origSize.X) / float32(origSize.Y)
		// if the image specifies it should be original size we need at least that many pixels on screen
		if img.FillMode == canvas.ImageFillOriginal {
			p.checkImageMinSize(img, origSize.X, origSize.Y)
		}

		tex := image.NewRGBA(img.Image.Bounds())
		draw.Draw(tex, img.Image.Bounds(), img.Image, image.ZP, draw.Src)

		return p.imgToTexture(tex)
	default:
		return p.imgToTexture(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	}
}

func (p *glPainter) checkImageMinSize(img *canvas.Image, pixX, pixY int) {
	pixSize := fyne.NewSize(unscaleInt(p.canvas, pixX), unscaleInt(p.canvas, pixY))

	if img.MinSize() != pixSize {
		img.SetMinSize(pixSize)
		canvas.Refresh(img) // force the initial size to be respected
	}
}

func (p *glPainter) newGlRasterTexture(obj fyne.CanvasObject) uint32 {
	rast := obj.(*canvas.Raster)

	width := p.textureScaleInt(rast.Size().Width)
	height := p.textureScaleInt(rast.Size().Height)

	return p.imgToTexture(rast.Generator(width, height))
}

func (p *glPainter) newGlLinearGradientTexture(obj fyne.CanvasObject) uint32 {
	gradient := obj.(*canvas.LinearGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height))
}

func (p *glPainter) newGlRadialGradientTexture(obj fyne.CanvasObject) uint32 {
	gradient := obj.(*canvas.RadialGradient)

	width := p.textureScaleInt(gradient.Size().Width)
	height := p.textureScaleInt(gradient.Size().Height)

	return p.imgToTexture(gradient.Generate(width, height))
}
