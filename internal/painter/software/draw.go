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

	"golang.org/x/image/draw"
)

type gradient interface {
	Generate(int, int) image.Image
	Size() fyne.Size
}

func drawArc(c fyne.Canvas, arc *canvas.Arc, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(arc)
	scaledWidth := scale.ToScreenCoordinate(c, arc.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, arc.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawArc(arc, pad, func(in float32) float32 {
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

func drawCircle(c fyne.Canvas, circle *canvas.Circle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(circle)
	scaledWidth := scale.ToScreenCoordinate(c, circle.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, circle.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawCircle(circle, pad, func(in float32) float32 {
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

	origImg := img.Image
	if img.FillMode != canvas.ImageFillCover {
		origImg = painter.PaintImage(img, c, width, height)
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
	} else if img.FillMode == canvas.ImageFillCover {
		inner := origImg.Bounds()
		imgAspect := img.Aspect()
		objAspect := float32(width) / float32(height)

		if objAspect > imgAspect {
			newHeight := float32(width) / imgAspect
			heightPad := (newHeight - float32(height)) / 2
			pixPad := int((heightPad / newHeight) * float32(inner.Dy()))

			inner = image.Rect(inner.Min.X, inner.Min.Y+pixPad, inner.Max.X, inner.Max.Y-pixPad)
		} else if objAspect < imgAspect {
			newWidth := float32(height) * imgAspect
			widthPad := (newWidth - float32(width)) / 2
			pixPad := int((widthPad / newWidth) * float32(inner.Dx()))

			inner = image.Rect(inner.Min.X+pixPad, inner.Min.Y, inner.Max.X-pixPad, inner.Max.Y)
		}

		subImg := image.NewRGBA(inner.Bounds())
		draw.Copy(subImg, inner.Min, origImg, inner, draw.Over, nil)
		origImg = subImg
	}

	cornerRadius := fyne.Min(painter.GetMaximumRadius(bounds), img.CornerRadius)
	drawPixels(scaledX, scaledY, width, height, img.ScaleMode, base, origImg, clip, img.Alpha(), cornerRadius*c.Scale())
}

func drawPixels(x, y, width, height int, mode canvas.ImageScale, base *image.NRGBA, origImg image.Image, clip image.Rectangle, alpha float64, radius float32) {
	if origImg.Bounds().Dx() == width && origImg.Bounds().Dy() == height && radius < 0.5 {
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

	if radius > 0.5 {
		applyRoundedCorners(scaledImg, width, height, radius)
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
		drawPixels(scaledX, scaledY, width, height, rast.ScaleMode, base, pix, clip, 1.0, 0.0)
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

	raw := painter.DrawRectangle(obj.(*canvas.Rectangle), width, height, pad, func(in float32) float32 {
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

func drawPolygon(c fyne.Canvas, polygon *canvas.Polygon, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(polygon)
	scaledWidth := scale.ToScreenCoordinate(c, polygon.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, polygon.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawPolygon(polygon, pad, func(in float32) float32 {
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

func drawRectangle(c fyne.Canvas, rect *canvas.Rectangle, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	topRightRadius := painter.GetCornerRadius(rect.TopRightCornerRadius, rect.CornerRadius)
	topLeftRadius := painter.GetCornerRadius(rect.TopLeftCornerRadius, rect.CornerRadius)
	bottomRightRadius := painter.GetCornerRadius(rect.BottomRightCornerRadius, rect.CornerRadius)
	bottomLeftRadius := painter.GetCornerRadius(rect.BottomLeftCornerRadius, rect.CornerRadius)
	drawOblong(c, rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rect.Aspect, pos, base, clip)
}

func drawOblong(c fyne.Canvas, obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, aspect float32, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
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

	if (stroke != nil && strokeWidth > 0) || topRightRadius != 0 || topLeftRadius != 0 || bottomRightRadius != 0 || bottomLeftRadius != 0 { // use a rasterizer if there is a stroke or radius
		drawOblongStroke(c, obj, width, height, pos, base, clip)
		return
	}

	scaledWidth := scale.ToScreenCoordinate(c, width)
	scaledHeight := scale.ToScreenCoordinate(c, height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))
	draw.Draw(base, bounds, image.NewUniform(fill), image.Point{}, draw.Over)
}

// applyRoundedCorners rounds the corners of the image in-place
func applyRoundedCorners(img *image.NRGBA, w, h int, radius float32) {
	rInt := int(math.Ceil(float64(radius)))

	aaWidth := float32(0.5)
	outerR2 := (radius + aaWidth) * (radius + aaWidth)
	innerR2 := (radius - aaWidth) * (radius - aaWidth)

	applyCorner := func(startX, endX, startY, endY int, cx, cy float32) {
		for y := startY; y < endY; y++ {
			for x := startX; x < endX; x++ {
				dx := float32(x) - cx
				dy := float32(y) - cy
				dist2 := dx*dx + dy*dy

				i := img.PixOffset(x, y)
				alpha := img.Pix[i+3]

				switch {
				case dist2 >= outerR2:
					img.Pix[i+3] = 0 // Fully transparent
				case dist2 > innerR2:
					// Linear falloff based on squared distance
					t := (outerR2 - dist2) / (outerR2 - innerR2) // t ranges from 0 to 1
					newAlpha := uint8(float32(alpha) * t)
					img.Pix[i+3] = newAlpha
				}
			}
		}
	}

	// Top-left
	r := minInt(rInt, minInt(w, h))
	applyCorner(0, r, 0, r, radius, radius)

	// Top-right
	applyCorner(w-r, w, 0, r, float32(w)-radius, radius)

	// Bottom-left
	applyCorner(0, r, h-r, h, radius, float32(h)-radius)

	// Bottom-right
	applyCorner(w-r, w, h-r, h, float32(w)-radius, float32(h)-radius)
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
