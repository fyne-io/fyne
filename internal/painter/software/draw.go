package software

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"

	"github.com/anthonynsimon/bild/blur"

	"golang.org/x/image/draw"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func drawCircle(c fyne.Canvas, circle *canvas.Circle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(circle)
	scaledWidth := scale.ToScreenCoordinate(c, circle.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, circle.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawCircle(circle, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	if circle.ShadowColor != color.Transparent && circle.ShadowColor != nil {
		// circle.ShadowType has no effect, always BoxShadow is drawn
		shadowCircle := &canvas.Circle{FillColor: circle.ShadowColor}
		shadowCircle.Resize(circle.Size())
		shadow := painter.DrawCircle(shadowCircle, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})

		pads := circle.ShadowPaddings()
		shadowPadLeft := scale.ToScreenCoordinate(c, pads[0])
		shadowPadRight := scale.ToScreenCoordinate(c, pads[2])
		shadowPadTop := scale.ToScreenCoordinate(c, pads[1])
		shadowPadBottom := scale.ToScreenCoordinate(c, pads[3])
		shadowRect := image.Rect(
			scaledX+shadowPadLeft,
			scaledY+shadowPadTop,
			scaledX+scaledWidth+shadowPadRight+shadowPadLeft,
			scaledY+scaledHeight+shadowPadBottom+shadowPadTop,
		)
		bounds = clip.Intersect(shadowRect)

		offset := image.Point{
			X: scale.ToScreenCoordinate(c, float32(circle.ShadowOffset.X)),
			Y: scale.ToScreenCoordinate(c, float32(-circle.ShadowOffset.Y)),
		}
		shadowBounds := clip.Intersect(
			image.Rect(
				shadowRect.Min.X-offset.X, shadowRect.Min.Y-offset.Y,
				shadowRect.Max.X, shadowRect.Max.Y,
			),
		)

		blurred := blur.Gaussian(shadow, float64(scale.ToScreenCoordinate(c, circle.ShadowSoftness)))
		draw.Draw(base, shadowBounds, blurred, image.Point{}, draw.Over)
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
}

func drawGradient(c fyne.Canvas, g gradient, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := g.Size()
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	tex := g.Generate(width, height)
	drawTex(scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y), width, height, base, tex, clip, 1.0)
}

func drawImage(c fyne.Canvas, img *canvas.Image, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := img.Size()
	if bounds.IsZero() {
		return
	}
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)

	origImg := painter.PaintImage(img, c, width, height)

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

	drawPixels(scaledX, scaledY, width, height, img.ScaleMode, base, origImg, clip, img.Alpha())
}

func drawPixels(x, y, width, height int, mode canvas.ImageScale, base *image.NRGBA, origImg image.Image, clip image.Rectangle, alpha float64) {
	if origImg.Bounds().Dx() == width && origImg.Bounds().Dy() == height {
		// do not scale or duplicate image since not needed, draw directly
		drawTex(x, y, width, height, base, origImg, clip, alpha)
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

	drawTex(x, y, width, height, base, scaledImg, clip, alpha)
}

func drawLine(c fyne.Canvas, line *canvas.Line, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(line)
	scaledWidth := scale.ToScreenCoordinate(c, line.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, line.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawLine(line, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
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

func drawTex(x, y, width, height int, base *image.NRGBA, tex image.Image, clip image.Rectangle, alpha float64) {
	outBounds := image.Rect(x, y, x+width, y+height)
	clippedBounds := clip.Intersect(outBounds)
	srcPt := image.Point{X: clippedBounds.Min.X - outBounds.Min.X, Y: clippedBounds.Min.Y - outBounds.Min.Y}
	if alpha == 1.0 {
		draw.Draw(base, clippedBounds, tex, srcPt, draw.Over)
	} else {
		mask := &image.Uniform{C: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: uint8(float64(0xff) * alpha)}}
		draw.DrawMask(base, clippedBounds, tex, srcPt, mask, srcPt, draw.Over)
	}
}

func drawText(c fyne.Canvas, text *canvas.Text, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := text.MinSize()
	width := scale.ToScreenCoordinate(c, bounds.Width+painter.VectorPad(text))
	height := scale.ToScreenCoordinate(c, bounds.Height)
	txtImg := image.NewRGBA(image.Rect(0, 0, width, height))

	color := text.Color
	if color == nil {
		color = theme.Color(theme.ColorNameForeground)
	}

	face := painter.CachedFontFace(text.TextStyle, text.FontSource, text)
	painter.DrawString(txtImg, text.Text, color, face.Fonts, text.TextSize, c.Scale(), text.TextStyle)

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
}

func drawRaster(c fyne.Canvas, rast *canvas.Raster, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	bounds := rast.Size()
	if bounds.IsZero() {
		return
	}
	width := scale.ToScreenCoordinate(c, bounds.Width)
	height := scale.ToScreenCoordinate(c, bounds.Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)

	pix := rast.Generator(width, height)
	if pix.Bounds().Bounds().Dx() != width || pix.Bounds().Dy() != height {
		drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip, 1.0)
	} else {
		drawTex(scaledX, scaledY, width, height, base, pix, clip, 1.0)
	}
}

func drawOblongStroke(c fyne.Canvas, obj fyne.CanvasObject, width, height float32, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(obj)
	scaledWidth := scale.ToScreenCoordinate(c, width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	var raw *image.RGBA
	switch o := obj.(type) {
	case *canvas.Square:
		raw = painter.DrawSquare(o, width, height, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	default:
		raw = painter.DrawRectangle(obj.(*canvas.Rectangle), width, height, pad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
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
}

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	drawOblong(c, rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, rect.CornerRadius, rect.Aspect, rect.ShadowSoftness, rect.ShadowOffset, rect.ShadowColor, pos, base, clip)
}

func drawSquare(c fyne.Canvas, sq *canvas.Square, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	drawOblong(c, sq, sq.FillColor, sq.StrokeColor, sq.StrokeWidth, sq.CornerRadius, 1.0, sq.ShadowSoftness, sq.ShadowOffset, sq.ShadowColor, pos, base, clip)
}

func drawOblong(c fyne.Canvas, obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth, radius, aspect float32, shadowSoftness float32, shadowOffset fyne.Position, shadowColor color.Color, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	width, height := obj.Size().Components()
	if aspect != 0 {
		frameAspect := width / height

		xPad, yPad := float32(0), float32(0)
		if frameAspect > aspect {
			newWidth := height * aspect
			xPad = (width - newWidth) / 2
			width = newWidth
		} else if frameAspect < aspect {
			newHeight := width / aspect
			yPad = (height - newHeight) / 2
			height = newHeight
		}

		pos = pos.AddXY(xPad, yPad)
	}

	if (stroke != nil && strokeWidth > 0) || radius != 0 { // use a rasterizer if there is a stroke or radius
		drawOblongStroke(c, obj, width, height, pos, base, clip)
		return
	}

	scaledWidth := scale.ToScreenCoordinate(c, width)
	scaledHeight := scale.ToScreenCoordinate(c, height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	if shadowColor != color.Transparent && shadowColor != nil {
		// shadowType has no effect, always BoxShadow is drawn
		shadowRectangle := &canvas.Rectangle{FillColor: shadowColor}
		shadowRectangle.Resize(obj.Size())
		shadow := painter.DrawRectangle(shadowRectangle, width, height, shadowSoftness, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})

		var shadowPadLeft, shadowPadRight, shadowPadTop, shadowPadBottom int
		if s, ok := obj.(canvas.Shadowable); ok {
			pads := s.ShadowPaddings()
			shadowPadLeft = scale.ToScreenCoordinate(c, pads[0])
			shadowPadRight = scale.ToScreenCoordinate(c, pads[2])
			shadowPadTop = scale.ToScreenCoordinate(c, pads[1])
			shadowPadBottom = scale.ToScreenCoordinate(c, pads[3])
		}

		shadowRect := image.Rect(
			scaledX+shadowPadLeft,
			scaledY+shadowPadTop,
			scaledX+scaledWidth+shadowPadRight+shadowPadLeft,
			scaledY+scaledHeight+shadowPadBottom+shadowPadTop,
		)
		bounds = clip.Intersect(shadowRect)

		// shadowSoftness is used as a vector pad so the position is affected by this value 
		// adding shadow softness to the offset restore initial position
		offset := image.Point{
			X: scale.ToScreenCoordinate(c, float32(shadowOffset.X+shadowSoftness)),
			Y: scale.ToScreenCoordinate(c, float32(-shadowOffset.Y+shadowSoftness)),
		}
		shadowBounds := clip.Intersect(
			image.Rect(
				shadowRect.Min.X-offset.X, shadowRect.Min.Y-offset.Y,
				shadowRect.Max.X, shadowRect.Max.Y,
			),
		)

		blurred := blur.Gaussian(shadow, float64(scale.ToScreenCoordinate(c, shadowSoftness)))
		draw.Draw(base, shadowBounds, blurred, image.Point{}, draw.Over)

		// due to shadow draw rectangle with a certain width and height
		raw := painter.DrawRectangle(canvas.NewRectangle(fill), width, height, 0, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
		draw.Draw(base, bounds, raw, image.Point{}, draw.Over)
	} else {
		draw.Draw(base, bounds, image.NewUniform(fill), image.Point{}, draw.Over)
	}
}
