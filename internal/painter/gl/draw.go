package gl

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	paint "fyne.io/fyne/v2/internal/painter"
)

func (p *painter) createBuffer(points []float32) Buffer {
	vbo := p.ctx.CreateBuffer()
	p.logError()
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferData(arrayBuffer, points, staticDraw)
	p.logError()
	return vbo
}

func (p *painter) defineVertexArray(prog Program, name string, size, stride, offset int) {
	vertAttrib := p.ctx.GetAttribLocation(prog, name)
	p.ctx.EnableVertexAttribArray(vertAttrib)
	p.ctx.VertexAttribPointerWithOffset(vertAttrib, size, float, false, stride*floatSize, offset*floatSize)
	p.logError()
}

func (p *painter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(circle, p.newGlCircleTexture, pos, circle.Size(), frame, canvas.ImageFillStretch,
		1.0, paint.VectorPad(circle))
}

func (p *painter) drawGradient(o fyne.CanvasObject, texCreator func(fyne.CanvasObject) Texture, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(o, texCreator, pos, o.Size(), frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *painter) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
}

func (p *painter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	if line.StrokeColor == color.Transparent || line.StrokeColor == nil || line.StrokeWidth == 0 {
		return
	}

	points, halfWidth, feather := p.lineCoords(pos, line.Position1, line.Position2, line.StrokeWidth, 0.5, frame)
	p.ctx.UseProgram(p.lineProgram)
	vbo := p.createBuffer(points)
	p.defineVertexArray(p.lineProgram, "vert", 2, 4, 0)
	p.defineVertexArray(p.lineProgram, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	colorUniform := p.ctx.GetUniformLocation(p.lineProgram, "color")
	r, g, b, a := line.StrokeColor.RGBA()
	if a == 0 {
		p.ctx.Uniform4f(colorUniform, 0, 0, 0, 0)
	} else {
		alpha := float32(a)
		p.ctx.Uniform4f(colorUniform, float32(r)/alpha, float32(g)/alpha, float32(b)/alpha, alpha/0xffff)
	}
	lineWidthUniform := p.ctx.GetUniformLocation(p.lineProgram, "lineWidth")
	p.ctx.Uniform1f(lineWidthUniform, halfWidth)

	featherUniform := p.ctx.GetUniformLocation(p.lineProgram, "feather")
	p.ctx.Uniform1f(featherUniform, feather)

	p.ctx.DrawArrays(triangles, 0, 6)
	p.logError()
	p.freeBuffer(vbo)
}

func (p *painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
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

func (p *painter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
}

func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
		return
	}
	p.drawTextureWithDetails(rect, p.newGlRectTexture, pos, rect.Size(), frame, canvas.ImageFillStretch,
		1.0, paint.VectorPad(rect))
}

func (p *painter) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if text.Text == "" || text.Text == " " {
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

func (p *painter) drawTextureWithDetails(o fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32, pad float32) {

	texture, err := p.getTexture(o, creator)
	if err != nil {
		return
	}

	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = paint.GetAspect(img)
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
	}
	points := p.rectCoords(size, pos, frame, fill, aspect, pad)
	p.ctx.UseProgram(p.program)
	vbo := p.createBuffer(points)
	p.defineVertexArray(p.program, "vert", 3, 5, 0)
	p.defineVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		p.ctx.BlendColor(0, 0, 0, alpha)
		p.ctx.BlendFunc(constantAlpha, oneMinusConstantAlpha)
	} else {
		p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	}
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
	p.freeBuffer(vbo)
}

func (p *painter) freeBuffer(vbo Buffer) {
	p.ctx.BindBuffer(arrayBuffer, noBuffer)
	p.logError()
	p.ctx.DeleteBuffer(vbo)
	p.logError()
}

func (p *painter) lineCoords(pos, pos1, pos2 fyne.Position, lineWidth, feather float32, frame fyne.Size) ([]float32, float32, float32) {
	// Shift line coordinates so that they match the target position.
	xPosDiff := pos.X - fyne.Min(pos1.X, pos2.X)
	yPosDiff := pos.Y - fyne.Min(pos1.Y, pos2.Y)
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	pos2.X = roundToPixel(pos2.X+xPosDiff, p.pixScale)
	pos2.Y = roundToPixel(pos2.Y+yPosDiff, p.pixScale)

	if lineWidth <= 1 {
		offset := float32(0.5)                  // adjust location for lines < 1pt on regular display
		if lineWidth <= 0.5 && p.pixScale > 1 { // and for 1px drawing on HiDPI (width 0.5)
			offset = 0.25
		}
		if pos1.X == pos2.X {
			pos1.X -= offset
			pos2.X -= offset
		}
		if pos1.Y == pos2.Y {
			pos1.Y -= offset
			pos2.Y -= offset
		}
	}

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
	halfWidth := (roundToPixel(lineWidth+feather, p.pixScale) * 0.5) / widthMultiplier
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
func (p *painter) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
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
	end := pos.Add(size)
	end.X = roundToPixel(end.X, pixScale)
	end.Y = roundToPixel(end.Y, pixScale)
	pos.X = roundToPixel(pos.X, pixScale)
	pos.Y = roundToPixel(pos.Y, pixScale)
	size.Width = end.X - pos.X
	size.Height = end.Y - pos.Y

	return size, pos
}
