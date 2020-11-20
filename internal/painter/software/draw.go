package software

import (
	"fmt"
	"image"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/painter"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func drawCircle(c fyne.Canvas, circle *canvas.Circle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(circle)
	scaledWidth := internal.ScaleInt(c, circle.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, circle.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawCircle(circle, pad, func(in float32) int {
		return int(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawGradient(c fyne.Canvas, g gradient, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := g.Size()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	tex := g.Generate(width, height)
	drawTex(internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y), width, height, base, tex, clip)
}

func drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := img.Size()
	if bounds.IsZero() {
		return
	}
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)

	origImg := painter.PaintImage(img, c, width, height)

	if img.FillMode == canvas.ImageFillContain {
		imgAspect := painter.GetAspect(img)
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

	drawPixels(scaledX, scaledY, width, height, img.ScaleMode, base, origImg, clip)
}

func drawPixels(x, y, width, height int, mode canvas.ImageScale, base *image.NRGBA, origImg image.Image, clip image.Rectangle) {
	scaledBounds := image.Rect(0, 0, width, height)
	scaledImg := image.NewNRGBA(scaledBounds)
	switch mode {
	case canvas.ImageScalePixels:
		draw.NearestNeighbor.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	default:
		if mode != canvas.ImageScaleSmooth {
			fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", mode), nil)
		}
		draw.CatmullRom.Scale(scaledImg, scaledBounds, origImg, origImg.Bounds(), draw.Over, nil)
	}

	drawTex(x, y, width, height, base, scaledImg, clip)
}

func drawLine(c fyne.Canvas, line *canvas.Line, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(line)
	scaledWidth := internal.ScaleInt(c, line.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, line.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawLine(line, pad, func(in float32) int {
		return int(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawTex(x, y, width, height int, base *image.NRGBA, tex image.Image, clip image.Rectangle) {
	outBounds := image.Rect(x, y, x+width, y+height)
	clippedBounds := clip.Intersect(outBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - outBounds.Min.X, Y: clippedBounds.Min.Y - outBounds.Min.Y}
	draw.Draw(base, clippedBounds, tex, srcPt, draw.Over)
}

func drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := text.MinSize()
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	txtImg := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	fontSize := float64(text.TextSize) * float64(c.Scale())
	opts.Size = fontSize
	opts.DPI = painter.TextDPI
	face := painter.CachedFontFace(text.TextStyle, &opts)

	d := font.Drawer{}
	d.Dst = txtImg
	d.Src = &image.Uniform{C: text.Color}
	d.Face = face
	d.Dot = freetype.Pt(0, height-face.Metrics().Descent.Ceil())
	d.DrawString(text.Text)

	size := text.Size()
	offsetX := 0
	offsetY := 0
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		offsetX = size.Width - bounds.Width
	case fyne.TextAlignCenter:
		offsetX = (size.Width - bounds.Width) / 2
	}
	if size.Height > bounds.Height {
		offsetY = (size.Height - bounds.Height) / 2
	}
	scaledX := internal.ScaleInt(c, pos.X+offsetX)
	scaledY := internal.ScaleInt(c, pos.Y+offsetY)
	imgBounds := image.Rect(scaledX, scaledY, scaledX+width, scaledY+height)
	clippedBounds := clip.Intersect(imgBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - imgBounds.Min.X, Y: clippedBounds.Min.Y - imgBounds.Min.Y}
	draw.Draw(base, clippedBounds, txtImg, srcPt, draw.Over)
}

func drawRaster(c fyne.Canvas, rast *canvas.Raster, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := rast.Size()
	if bounds.IsZero() {
		return
	}
	width := internal.ScaleInt(c, bounds.Width)
	height := internal.ScaleInt(c, bounds.Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)

	pix := rast.Generator(width, height)
	if pix.Bounds().Bounds().Dx() != width || pix.Bounds().Dy() != height {
		drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip)
	} else {
		drawTex(scaledX, scaledY, width, height, base, pix, clip)
	}
}

func drawRectangleStroke(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(rect)
	scaledWidth := internal.ScaleInt(c, rect.Size().Width+pad*2)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height+pad*2)
	scaledX, scaledY := internal.ScaleInt(c, pos.X-pad), internal.ScaleInt(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawRectangle(rect, pad, func(in float32) int {
		return int(math.Round(float64(in) * float64(c.Scale())))
	})

	// the clip intersect above cannot be negative, so we may need to compensate
	offX, offY := 0, 0
	if scaledX < 0 {
		offX = -scaledX
	}
	if scaledY < 0 {
		offY = -scaledY
	}
	draw.Draw(base, bounds, raw, image.Point{offX, offY}, draw.Over)
}

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	if rect.StrokeColor != nil && rect.StrokeWidth > 0 { // use a rasterizer if there is a stroke
		drawRectangleStroke(c, rect, pos, base, clip)
		return
	}

	scaledWidth := internal.ScaleInt(c, rect.Size().Width)
	scaledHeight := internal.ScaleInt(c, rect.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))
	draw.Draw(base, bounds, image.NewUniform(rect.FillColor), image.Point{}, draw.Over)
}

func drawWidget(c fyne.Canvas, wid fyne.Widget, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	scaledWidth := internal.ScaleInt(c, wid.Size().Width)
	scaledHeight := internal.ScaleInt(c, wid.Size().Height)
	scaledX, scaledY := internal.ScaleInt(c, pos.X), internal.ScaleInt(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))
	draw.Draw(base, bounds, image.NewUniform(cache.Renderer(wid).BackgroundColor()), image.Point{}, draw.Over)
}
