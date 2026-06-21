package software

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
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

func drawBlur(c fyne.Canvas, blurObj *canvas.Blur, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	if blurObj.Radius == 0 {
		return
	}

	scaledWidth := scale.ToScreenCoordinate(c, blurObj.Size().Width)
	scaledHeight := scale.ToScreenCoordinate(c, blurObj.Size().Height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	crop := base.SubImage(bounds)
	blurred := blur.Gaussian(crop, float64(blurObj.Radius*c.Scale()))
	draw.Draw(base, base.Bounds(), blurred, image.Point{}, draw.Over)
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

	if painter.IsShadowVisible(circle.Shadow) {
		drawShadow(c, circle, circle.Size(), circle.Shadow, pad, base, clip, pos)
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

func drawBezierCurve(c fyne.Canvas, bezierCurve *canvas.BezierCurve, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(bezierCurve)
	scaledWidth := scale.ToScreenCoordinate(c, bezierCurve.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, bezierCurve.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawBezierCurve(bezierCurve, pad, func(in float32) float32 {
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
	width := scale.ToScreenCoordinate(c, bounds.Width+painter.VectorPad(text)) // potentially italic overspill
	height := scale.ToScreenCoordinate(c, bounds.Height+painter.TextVectorPad) // space below for descenders / underline
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

	if text.TextStyle.Underline || text.TextStyle.Strikethrough {
		_, baseline := cache.GetFontMetrics(text.Text, text.TextSize, text.TextStyle, text.FontSource)
		line := canvas.NewLine(color)
		line.Resize(fyne.NewSize(bounds.Width, 0))
		if text.TextStyle.Underline {
			underlinePos := fyne.NewPos(pos.X, pos.Y+baseline+painter.UnderlineOffsetFromBaseline)
			drawLine(c, line, underlinePos, base, clip)
		}
		if text.TextStyle.Strikethrough {
			strikePos := fyne.NewPos(pos.X, pos.Y+baseline*painter.StrikethroughToBaselineFactor)
			drawLine(c, line, strikePos, base, clip)
		}
	}
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

func drawOblongStroke(c fyne.Canvas, obj fyne.CanvasObject, width, height float32, shadow canvas.Shadow, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(obj)
	scaledWidth := scale.ToScreenCoordinate(c, width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawRectangle(obj.(*canvas.Rectangle), width, height, pad, func(in float32) float32 {
		return float32(math.Round(float64(in) * float64(c.Scale())))
	})

	if painter.IsShadowVisible(shadow) {
		drawShadow(c, obj, fyne.NewSize(width, height), shadow, pad, base, clip, pos)
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

func drawPolygon(c fyne.Canvas, polygon *canvas.RegularPolygon, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
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

func drawArbitraryPolygon(c fyne.Canvas, polygon *canvas.ArbitraryPolygon, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	pad := painter.VectorPad(polygon)
	scaledWidth := scale.ToScreenCoordinate(c, polygon.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, polygon.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawArbitraryPolygon(polygon, pad, func(in float32) float32 {
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
	drawOblong(c, rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rect.Aspect, rect.Shadow, pos, base, clip)
}

func drawOblong(c fyne.Canvas, obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, aspect float32, shadow canvas.Shadow, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
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
		drawOblongStroke(c, obj, width, height, shadow, pos, base, clip)
		return
	}

	scaledWidth := scale.ToScreenCoordinate(c, width)
	scaledHeight := scale.ToScreenCoordinate(c, height)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X), scale.ToScreenCoordinate(c, pos.Y)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	if painter.IsShadowVisible(shadow) {
		drawShadow(c, obj, fyne.NewSize(width, height), shadow, 0, base, clip, pos)
	}

	draw.Draw(base, bounds, image.NewUniform(fill), image.Point{}, draw.Over)
}

func drawEllipse(c fyne.Canvas, ellipse *canvas.Ellipse, pos fyne.Position, base *image.NRGBA, clip image.Rectangle) {
	// when rotated, the ellipse needs more space
	// add half the difference between width and height as padding
	// with padding the final size is a square
	width, height := ellipse.Size().Components()
	rotPad := float32(math.Abs(float64(width)-float64(height)) / 2)
	xPad, yPad := float32(0), float32(0)

	if width > height {
		yPad = rotPad
	} else {
		xPad = rotPad
	}

	if painter.IsShadowVisible(ellipse.Shadow) {
		drawShadow(c, ellipse, ellipse.Size(), ellipse.Shadow, painter.VectorPad(ellipse), base, clip, pos)
	}

	pos = pos.AddXY(xPad, yPad)
	pad := painter.VectorPad(ellipse) + rotPad

	scaledWidth := scale.ToScreenCoordinate(c, ellipse.Size().Width+pad*2)
	scaledHeight := scale.ToScreenCoordinate(c, ellipse.Size().Height+pad*2)
	scaledX, scaledY := scale.ToScreenCoordinate(c, pos.X-pad), scale.ToScreenCoordinate(c, pos.Y-pad)
	bounds := clip.Intersect(image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight))

	raw := painter.DrawEllipse(ellipse, pad, func(in float32) float32 {
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

func drawShadow(c fyne.Canvas, obj fyne.CanvasObject, objSize fyne.Size, shadow canvas.Shadow, pad float32, base *image.NRGBA, clip image.Rectangle, pos fyne.Position) {
	shadowOffset := shadow.Offset
	shadowBlurRadius := shadow.BlurRadius
	shadowSpread := shadow.Spread
	shadowVariant := shadow.Variant
	shadowColor := shadow.Color

	var shadowRaw *image.RGBA
	var maskRaw *image.RGBA

	vPad := pad + shadowBlurRadius
	if shadowSpread < 0 {
		vPad -= shadowSpread
	}

	switch o := obj.(type) {
	case *canvas.Rectangle:
		shadowRaw = painter.DrawRectangle(&canvas.Rectangle{
			FillColor:               shadowColor,
			CornerRadius:            o.CornerRadius,
			TopRightCornerRadius:    o.TopRightCornerRadius,
			TopLeftCornerRadius:     o.TopLeftCornerRadius,
			BottomRightCornerRadius: o.BottomRightCornerRadius,
			BottomLeftCornerRadius:  o.BottomLeftCornerRadius,
		}, fyne.Max(objSize.Width+2*shadowSpread, 0), fyne.Max(objSize.Height+2*shadowSpread, 0), vPad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
		maskRaw = painter.DrawRectangle(&canvas.Rectangle{
			FillColor:               color.Opaque,
			CornerRadius:            o.CornerRadius,
			TopRightCornerRadius:    o.TopRightCornerRadius,
			TopLeftCornerRadius:     o.TopLeftCornerRadius,
			BottomRightCornerRadius: o.BottomRightCornerRadius,
			BottomLeftCornerRadius:  o.BottomLeftCornerRadius,
		}, objSize.Width, objSize.Height, vPad+shadowSpread, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	case *canvas.Circle:
		shadowCircle := &canvas.Circle{FillColor: shadowColor}
		shadowCircle.Resize(objSize.AddWidthHeight(2*shadowSpread, 2*shadowSpread))
		shadowRaw = painter.DrawCircle(shadowCircle, vPad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
		maskCircle := &canvas.Circle{FillColor: color.Opaque}
		maskCircle.Resize(objSize)
		maskRaw = painter.DrawCircle(maskCircle, vPad+shadowSpread, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	case *canvas.Ellipse:
		shadowEllipse := &canvas.Ellipse{FillColor: shadowColor}
		shadowEllipse.Resize(objSize.AddWidthHeight(2*shadowSpread, 2*shadowSpread))
		shadowRaw = painter.DrawEllipse(shadowEllipse, vPad, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
		maskEllipse := &canvas.Ellipse{FillColor: color.Opaque}
		maskEllipse.Resize(objSize)
		maskRaw = painter.DrawEllipse(maskEllipse, vPad+shadowSpread, func(in float32) float32 {
			return float32(math.Round(float64(in) * float64(c.Scale())))
		})
	}

	startX := pos.X + float32(shadowOffset.X) - shadowSpread - vPad
	startY := pos.Y + float32(shadowOffset.Y) - shadowSpread - vPad

	screenStartX := scale.ToScreenCoordinate(c, startX)
	screenStartY := scale.ToScreenCoordinate(c, startY)

	blurred := blur.Gaussian(shadowRaw, float64(shadowBlurRadius*c.Scale()))
	destRect := image.Rect(screenStartX, screenStartY, screenStartX+blurred.Bounds().Dx(), screenStartY+blurred.Bounds().Dy())
	shadowBounds := clip.Intersect(destRect)

	if shadowBounds.Empty() {
		return
	}

	// If DropShadow, subtract object from shadow
	if shadowVariant == canvas.DropShadow {
		dx := screenStartX - scale.ToScreenCoordinate(c, pos.X-shadowSpread-vPad)
		dy := screenStartY - scale.ToScreenCoordinate(c, pos.Y-shadowSpread-vPad)

		var fill, strokeCol color.Color
		var strokeWidth float32
		switch o := obj.(type) {
		case *canvas.Rectangle:
			fill, strokeCol, strokeWidth = o.FillColor, o.StrokeColor, o.StrokeWidth
		case *canvas.Circle:
			fill, strokeCol, strokeWidth = o.FillColor, o.StrokeColor, o.StrokeWidth
		case *canvas.Ellipse:
			fill, strokeCol, strokeWidth = o.FillColor, o.StrokeColor, o.StrokeWidth
		}

		var objAlpha float32
		if fill != nil {
			_, _, _, a := fill.RGBA()
			objAlpha = float32(a) / 65535.0
		}
		if strokeCol != nil && strokeWidth > 0 {
			_, _, _, a := strokeCol.RGBA()
			sa := float32(a) / 65535.0
			if sa > objAlpha {
				objAlpha = sa
			}
		}

		for y := 0; y < blurred.Bounds().Dy(); y++ {
			for x := 0; x < blurred.Bounds().Dx(); x++ {
				mx := x + dx
				my := y + dy

				_, _, _, maskA := maskRaw.At(mx, my).RGBA()
				if maskA > 0 {
					pixel := blurred.RGBAAt(x, y)
					cVal := float32(maskA) / 65535.0
					den := 1.0 - cVal*objAlpha
					var invMA float32
					if den <= 0 {
						invMA = 0
					} else {
						invMA = (1.0 - cVal) / den
					}
					pixel.R = uint8(float32(pixel.R) * invMA)
					pixel.G = uint8(float32(pixel.G) * invMA)
					pixel.B = uint8(float32(pixel.B) * invMA)
					pixel.A = uint8(float32(pixel.A) * invMA)
					blurred.SetRGBA(x, y, pixel)
				}
			}
		}
	}

	srcPt := image.Point{X: shadowBounds.Min.X - screenStartX, Y: shadowBounds.Min.Y - screenStartY}
	draw.Draw(base, shadowBounds, blurred, srcPt, draw.Over)
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
