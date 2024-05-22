package software

import (
	"fmt"
	"image"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"

	"golang.org/x/image/draw"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func (p *Painter) drawCircle(c fyne.Canvas, circle *canvas.Circle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	pad, bounds, clipBounds := circleCords(c, circle, pos, clip)

	raw := p.drawCircleWrapper(circle, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	mask := genMask(c, circle, bounds, base.Bounds())

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if bounds.Min.X < 0 {
		offX = -bounds.Min.X
	}
	if bounds.Min.Y < 0 {
		offY = -bounds.Min.Y
	}
	draw.DrawMask(base, clipBounds, raw, image.Point{offX, offY}, mask, clipBounds.Min, draw.Over)
	fmt.Println("drawCircle:", time.Since(t))
}

func (p *Painter) drawGradient(c fyne.Canvas, g gradient, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	bounds := gradientCords(c, g, pos)
	tex := g.Generate(bounds.Dx(), bounds.Dy())

	p.drawTex(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy(), base, tex, clip)
	fmt.Println("drawGradient:", time.Since(t))
}

func (p *Painter) drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	if img.Size().IsZero() {
		return
	}

	w, h, rect := imageCords(c, img, pos)

	origImg := p.paintImageWrapper(img, c, w, h)

	p.drawPixels(rect.Min.X, rect.Min.Y, w, h, img.ScaleMode, base, origImg, clip)
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
	pad, bounds, clipBounds := lineCords(c, line, pos, clip)

	raw := p.drawLineWrapper(line, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	mask := genMask(c, line, bounds, base.Bounds())

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if bounds.Min.X < 0 {
		offX = -bounds.Min.X
	}
	if bounds.Min.Y < 0 {
		offY = -bounds.Min.Y
	}
	draw.DrawMask(base, clipBounds, raw, image.Point{offX, offY}, mask, clipBounds.Min, draw.Over)
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

	color := text.Color
	if color == nil {
		color = theme.ForegroundColor()
	}

	bounds, clippedBounds := textCords(c, text, pos, clip)

	txtImg := p.drawStringWrapper(c, text, bounds.Dx(), bounds.Dy(), color)
	mask := genMask(c, text, bounds, base.Bounds())

	draw.DrawMask(base, clippedBounds, txtImg, image.Point{X: clippedBounds.Min.X - bounds.Min.X, Y: clippedBounds.Min.Y - bounds.Min.Y}, mask, clippedBounds.Min, draw.Over)
	fmt.Println("drawText:", time.Since(t))
}

func (p *Painter) drawRaster(c fyne.Canvas, rast *canvas.Raster, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	if rast.Size().IsZero() {
		return
	}
	bounds := rasterCords(c, rast, pos)
	width, height := bounds.Dx(), bounds.Dy()
	scaledX, scaledY := bounds.Min.X, bounds.Min.Y

	pix := p.drawRasterWrapper(rast, width, height)

	if pix.Bounds().Bounds().Dx() != width || pix.Bounds().Dy() != height {
		p.drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip)
	} else {
		p.drawTex(scaledX, scaledY, width, height, base, pix, clip)
	}
	fmt.Println("drawRaster:", time.Since(t))
}

func (p *Painter) drawRectangleStroke(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	t := time.Now()
	pad, bounds, clipBounds := rectangleStrokeCords(c, rect, pos, clip)

	raw := p.drawRectangleStrokeWrapper(rect, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	mask := genMask(c, rect, bounds, base.Bounds())

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if bounds.Min.X < 0 {
		offX = -bounds.Min.X
	}
	if bounds.Min.Y < 0 {
		offY = -bounds.Min.Y
	}
	draw.DrawMask(base, clipBounds, raw, image.Point{offX, offY}, mask, clipBounds.Min, draw.Over)
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
		return Texture(cache.NoTexture)
	})

	bounds, clipBounds := rectangleCords(c, rect, pos, clip)
	mask := genMask(c, rect, bounds, base.Bounds())

	draw.DrawMask(base, clipBounds, image.NewUniform(rect.FillColor), image.Point{}, mask, bounds.Min, draw.Over)
	fmt.Println("drawRectangle:", time.Since(t))
}
