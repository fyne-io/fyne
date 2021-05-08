package gl

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
)

func (p *glPainter) drawTextureWithDetails(o fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32, pad float32) {

	texture := getTexture(o, creator)
	if texture == NoTexture {
		return
	}

	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = painter.GetAspect(img)
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
	}
	points := p.rectCoords(size, pos, frame, fill, aspect, pad)
	vbo := p.glCreateBuffer(points)

	p.glDrawTexture(texture, alpha)
	p.glFreeBuffer(vbo)
}

func (p *glPainter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(circle, p.newGlCircleTexture, pos, circle.Size(), frame, canvas.ImageFillStretch,
		1.0, painter.VectorPad(circle))
}

func (p *glPainter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	points, halfWidth, feather := p.lineCoords(pos, line.Position1, line.Position2, line.StrokeWidth, 0.5, frame)
	vbo := p.glCreateLineBuffer(points)
	p.glDrawLine(halfWidth, line.StrokeColor, feather)
	p.glFreeBuffer(vbo)
}

func (p *glPainter) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
}

func (p *glPainter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
}

func (p *glPainter) drawGradient(o fyne.CanvasObject, texCreator func(fyne.CanvasObject) Texture, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(o, texCreator, pos, o.Size(), frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *glPainter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(rect, p.newGlRectTexture, pos, rect.Size(), frame, canvas.ImageFillStretch,
		1.0, painter.VectorPad(rect))
}

func (p *glPainter) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if text.Text == "" {
		return
	}

	size := text.MinSize()
	containerSize := text.Size()
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		pos = fyne.NewPos(pos.X+containerSize.Width-size.Width, pos.Y)
	case fyne.TextAlignCenter:
		pos = fyne.NewPos(pos.X+(containerSize.Width-size.Width)/2, pos.Y)
	}

	if containerSize.Height > size.Height {
		pos = fyne.NewPos(pos.X, pos.Y+(containerSize.Height-size.Height)/2)
	}

	p.drawTextureWithDetails(text, p.newGlTextTexture, pos, size, frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *glPainter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if !o.Visible() {
		return
	}
	switch obj := o.(type) {
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		p.drawImage(obj, pos, frame)
	case *canvas.Raster:
		p.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame)
	case *canvas.LinearGradient:
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	}
}

func (p *glPainter) lineCoords(pos, pos1, pos2 fyne.Position, lineWidth, feather float32, frame fyne.Size) ([]float32, float32, float32) {
	// Shift line coordinates so that they match the target position.
	xPosDiff := pos.X - fyne.Min(pos1.X, pos2.X)
	yPosDiff := pos.Y - fyne.Min(pos1.Y, pos2.Y)
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	pos2.X = roundToPixel(pos2.X+xPosDiff, p.pixScale)
	pos2.Y = roundToPixel(pos2.Y+yPosDiff, p.pixScale)

	x1Pos := pos1.X / frame.Width
	x1 := -1 + x1Pos*2
	y1Pos := pos1.Y / frame.Height
	y1 := 1 - y1Pos*2
	x2Pos := pos2.X / frame.Width
	x2 := -1 + x2Pos*2
	y2Pos := pos2.Y / frame.Height
	y2 := 1 - y2Pos*2

	normalX := (pos2.Y - pos1.Y) / frame.Width
	normalY := (pos2.X - pos1.X) / frame.Height
	dirLength := float32(math.Sqrt(float64(normalX*normalX + normalY*normalY)))
	normalX /= dirLength
	normalY /= dirLength

	normalObjX := normalX * 0.5 * frame.Width
	normalObjY := normalY * 0.5 * frame.Height
	widthMultiplier := float32(math.Sqrt(float64(normalObjX*normalObjX + normalObjY*normalObjY)))
	halfWidth := (lineWidth*0.5 + feather) / widthMultiplier
	featherWidth := feather / widthMultiplier

	return []float32{
		// coord x, y normal x, y
		x1, y1, normalX, normalY,
		x2, y2, normalX, normalY,
		x2, y2, -normalX, -normalY,
		x2, y2, -normalX, -normalY,
		x1, y1, normalX, normalY,
		x1, y1, -normalX, -normalY,
	}, halfWidth, featherWidth
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (p *glPainter) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad float32) []float32 {
	size, pos = rectInnerCoords(size, pos, fill, aspect)
	size, pos = roundToPixelCoords(size, pos, p.pixScale)

	xPos := (pos.X - pad) / frame.Width
	x1 := -1 + xPos*2
	x2Pos := (pos.X + size.Width + pad) / frame.Width
	x2 := -1 + x2Pos*2

	yPos := (pos.Y - pad) / frame.Height
	y1 := 1 - yPos*2
	y2Pos := (pos.Y + size.Height + pad) / frame.Height
	y2 := 1 - y2Pos*2

	return []float32{
		// coord x, y, z texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}
}

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill == canvas.ImageFillContain || fill == canvas.ImageFillOriginal {
		// change pos and size accordingly

		viewAspect := size.Width / size.Height

		newWidth, newHeight := size.Width, size.Height
		widthPad, heightPad := float32(0), float32(0)
		if viewAspect > aspect {
			newWidth = size.Height * aspect
			widthPad = (size.Width - newWidth) / 2
		} else if viewAspect < aspect {
			newHeight = size.Width / aspect
			heightPad = (size.Height - newHeight) / 2
		}

		return fyne.NewSize(newWidth, newHeight), fyne.NewPos(pos.X+widthPad, pos.Y+heightPad)
	}

	return size, pos
}

func roundToPixel(v float32, pixScale float32) float32 {
	if pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}

	return float32(math.Round(float64(v*pixScale))) / pixScale
}

func roundToPixelCoords(size fyne.Size, pos fyne.Position, pixScale float32) (fyne.Size, fyne.Position) {
	size.Width = roundToPixel(size.Width, pixScale)
	size.Height = roundToPixel(size.Height, pixScale)
	pos.X = roundToPixel(pos.X, pixScale)
	pos.Y = roundToPixel(pos.Y, pixScale)

	return size, pos
}
