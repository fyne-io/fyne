package software

import (
	"fmt"
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"

	"golang.org/x/image/draw"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func (p *Painter) drawCircle(c fyne.Canvas, circle *canvas.Circle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	pad := painter.VectorPad(circle)
	scaledWidth := scale.ToScreenCoordinate(c, circle.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, circle.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw, err := p.getTexture(circle, func(object fyne.CanvasObject) Texture {
		return painter.DrawCircle(circle, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	})
	if err != nil {
		return
	}

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
	fmt.Println("drawCircle:", time.Since(t))
}

func (p *Painter) drawGradient(c fyne.Canvas, g gradient, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	bounds := g.Size()
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	tex := g.Generate(width, height)
	p.drawTex(scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y), width, height, base, tex, clip)
	fmt.Println("drawGradient:", time.Since(t))
}

func (p *Painter) drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	bounds := img.Size()
	if bounds.IsZero() {
		return
	}
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)

	origImg, err := p.getTexture(img, func(object fyne.CanvasObject) Texture {
		return painter.PaintImage(img, c, width, height)
	})
	if err != nil {
		return
	}

	if img.FillMode == canvas.ImageFillContain {
		imgAspect := img.Aspect()
		objAspect := float32(width) / float32(height)

		if objAspect > imgAspect {
			newWidth := int(float32(height) * imgAspect)
			scaledX += (width - newWidth) / 2
			width = newWidth
		} else if objAspect < imgAspect {
			newHeight := int(float32(width) / imgAspect)
			scaledY += (height - newHeight) / 2
			height = newHeight
		}
	}

	p.drawPixels(scaledX, scaledY, width, height, img.ScaleMode, base, origImg, clip)
	fmt.Println("drawImage:", time.Since(t))
}

func (p *Painter) drawPixels(x, y, width, height int, mode canvas.ImageScale, base *image.NRGBA, origImg image.Image, clip image.Rectangle) {
	t := time.Now()
	if origImg.Bounds().Dx() == width && origImg.Bounds().Dy() == height {
		// do not scale or duplicate image since not needed, draw directly
		p.drawTex(x, y, width, height, base, origImg, clip)
		return
	}

	scaledBounds := image.Rect(0, 0, width, height)
	scaledImg := image.NewNRGBA(scaledBounds)
	switch mode {
	case canvas.ImageScalePixels:
		draw.NearestNeighbor.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	case canvas.ImageScaleFastest:
		draw.ApproxBiLinear.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	default:
		if mode != canvas.ImageScaleSmooth {
			fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", mode), nil)
		}
		draw.CatmullRom.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	}

	p.drawTex(x, y, width, height, base, scaledImg, clip)
	fmt.Println("drawPixels:", time.Since(t))
}

func (p *Painter) drawLine(c fyne.Canvas, line *canvas.Line, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	pad := painter.VectorPad(line)
	scaledWidth := scale.ToScreenCoordinate(c, line.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, line.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw, err := p.getTexture(line, func(object fyne.CanvasObject) Texture {
		return painter.DrawLine(line, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	})
	if err != nil {
		return
	}

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
	fmt.Println("drawLine:", time.Since(t))
}

func (p *Painter) drawTex(x, y, width, height int, base *image.NRGBA, tex image.Image, clip image.Rectangle) {
	t := time.Now()
	outBounds := image.Rect(x, y, x+width, y+height)
	clippedBounds := clip.Intersect(outBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - outBounds.Min.X, Y: clippedBounds.Min.Y - outBounds.Min.Y}
	draw.Draw(base, clippedBounds, tex, srcPt, draw.Over)
	fmt.Println("drawTex:", time.Since(t))
}

func (p *Painter) drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	bounds := text.MinSize()
	width := scale.ToScreenCoordinate(c, bounds.Width+painter.VectorPad(text))
	height := scale.ToScreenCoordinate(c, bounds.Height)

	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	txtImg, err := p.getTexture(text, func(object fyne.CanvasObject) Texture {
		txtImg := image.NewRGBA(image.Rect(0, 0, width, height))
		face := painter.CachedFontFace(text.TextStyle, text.TextSize*c.Scale(), 1)
		painter.DrawString(txtImg, text.Text, color, face.Fonts, text.TextSize, c.Scale(), text.TextStyle.TabWidth)
		return txtImg
	})
	if err != nil {
		fyne.LogError("Failed to render text", err)
		return
	}

	size := text.Size()
	offsetX := float32(0)
	offsetY := float32(0)
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		offsetX = size.Width - bounds.Width
	case fyne.TextAlignCenter:
		offsetX = (size.Width - bounds.Width) / 2
	}
	if size.Height > bounds.Height {
		offsetY = (size.Height - bounds.Height) / 2
	}
	scaledX := scale.ToScreenCoordinate(c, pos.X+offsetX)
	scaledY := scale.ToScreenCoordinate(c, pos.Y+offsetY)
	imgBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	clippedBounds := clip.Intersect(imgBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - imgBounds.Min.X, Y: clippedBounds.Min.Y - imgBounds.Min.Y}
	draw.Draw(base, clippedBounds, txtImg, srcPt, draw.Over)
	fmt.Println("drawText:", time.Since(t))
}

func (p *Painter) drawRaster(c fyne.Canvas, rast *canvas.Raster, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	bounds := rast.Size()
	if bounds.IsZero() {
		return
	}
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)

	pix, err := p.getTexture(rast, func(object fyne.CanvasObject) Texture {
		return rast.Generator(width, height)
	})
	if err != nil {
		return
	}

	if pix.Bounds().Bounds().Dx() != width || pix.Bounds().Dy() != height {
		p.drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip)
	} else {
		p.drawTex(scaledX, scaledY, width, height, base, pix, clip)
	}
	fmt.Println("drawRaster:", time.Since(t))
}

func (p *Painter) drawRectangleStroke(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	pad := painter.VectorPad(rect)
	scaledWidth := scale.ToScreenCoordinate(c, rect.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, rect.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw, err := p.getTexture(rect, func(object fyne.CanvasObject) Texture {
		return painter.DrawRectangle(rect, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	})
	if err != nil {
		return
	}

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
	fmt.Println("drawRectangleStroke:", time.Since(t))
}

func (p *Painter) drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	if (rect.StrokeColor != nil && rect.StrokeWidth > 0) || rect.CornerRadius != 0 { // use a rasterizer if there is a stroke or radius
		p.drawRectangleStroke(c, rect, pos, base, clip)
		return
	}

	// allows us to keep track of if it's been drawn before
	p.getTexture(rect, func(object fyne.CanvasObject) Texture {
		return cache.NoTexture
	})

	scaledWidth := scale.ToScreenCoordinate(c, rect.Size().Width)
	scaledHeight := scale.ToScreenCoordinate(c, rect.Size().Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.Point{}, draw.Over)
	fmt.Println("drawRectangle:", time.Since(t))
}
