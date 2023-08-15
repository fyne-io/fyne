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
	r, g, b, a := getFragmentColor(line.StrokeColor)
	p.ctx.Uniform4f(colorUniform, r, g, b, a)

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

	roundedCorners := rect.CornerRadius != 0
	var program Program
	if roundedCorners {
		program = p.roundRectangleProgram
	} else {
		program = p.rectangleProgram
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, rect, frame)
	p.ctx.UseProgram(program)
	vbo := p.createBuffer(points)
	p.defineVertexArray(program, "vert", 2, 4, 0)
	p.defineVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameSizeUniform := p.ctx.GetUniformLocation(program, "frame_size")
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.ctx.Uniform2f(frameSizeUniform, frameWidthScaled, frameHeightScaled)

	rectCoordsUniform := p.ctx.GetUniformLocation(program, "rect_coords")
	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.ctx.Uniform4f(rectCoordsUniform, x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(rect.StrokeWidth*p.pixScale, 1.0)
	if roundedCorners {
		strokeUniform := p.ctx.GetUniformLocation(program, "stroke_width_half")
		p.ctx.Uniform1f(strokeUniform, strokeWidthScaled*0.5)

		rectSizeUniform := p.ctx.GetUniformLocation(program, "rect_size_half")
		rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
		rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
		p.ctx.Uniform2f(rectSizeUniform, rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

		radiusUniform := p.ctx.GetUniformLocation(program, "radius")
		radiusScaled := roundToPixel(rect.CornerRadius*p.pixScale, 1.0)
		p.ctx.Uniform1f(radiusUniform, radiusScaled)
	} else {
		strokeUniform := p.ctx.GetUniformLocation(program, "stroke_width")
		p.ctx.Uniform1f(strokeUniform, strokeWidthScaled)
	}

	var r, g, b, a float32
	fillColorUniform := p.ctx.GetUniformLocation(program, "fill_color")
	r, g, b, a = getFragmentColor(rect.FillColor)
	p.ctx.Uniform4f(fillColorUniform, r, g, b, a)

	strokeColorUniform := p.ctx.GetUniformLocation(program, "stroke_color")
	strokeColor := rect.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.ctx.Uniform4f(strokeColorUniform, r, g, b, a)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
	p.freeBuffer(vbo)
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

	// text size is sensitive to position on screen
	size, _ = roundToPixelCoords(size, text.Position(), p.pixScale)
	size.Width += roundToPixel(paint.VectorPad(text), p.pixScale)
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
		aspect = img.Aspect()
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

func (p *painter) vecRectCoords(pos fyne.Position, rect *canvas.Rectangle, frame fyne.Size) ([4]float32, []float32) {
	size := rect.Size()
	pos1 := rect.Position()

	xPosDiff := pos.X - pos1.X
	yPosDiff := pos.Y - pos1.Y
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)

	x1Pos := pos1.X
	x1Norm := -1 + x1Pos*2/frame.Width
	x2Pos := pos1.X + size.Width
	x2Norm := -1 + x2Pos*2/frame.Width
	y1Pos := pos1.Y
	y1Norm := 1 - y1Pos*2/frame.Height
	y2Pos := pos1.Y + size.Height
	y2Norm := 1 - y2Pos*2/frame.Height

	// output a norm for the fill and the vert is unused, but we pass 0 to avoid optimisation issues
	coords := []float32{
		0, 0, x1Norm, y1Norm, // first triangle
		0, 0, x2Norm, y1Norm, // second triangle
		0, 0, x1Norm, y2Norm,
		0, 0, x2Norm, y2Norm}

	return [4]float32{x1Pos, y1Pos, x2Pos, y2Pos}, coords
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

// Returns FragmentColor(red,green,blue,alpha) from fyne.Color
func getFragmentColor(col color.Color) (float32, float32, float32, float32) {
	if col == nil {
		return 0, 0, 0, 0
	}
	r, g, b, a := col.RGBA()
	if a == 0 {
		return 0, 0, 0, 0
	}
	alpha := float32(a)
	return float32(r) / alpha, float32(g) / alpha, float32(b) / alpha, alpha / 0xffff
}

func (p *painter) scaleFrameSize(frame fyne.Size) (float32, float32) {
	frameWidthScaled := roundToPixel(frame.Width*p.pixScale, 1.0)
	frameHeightScaled := roundToPixel(frame.Height*p.pixScale, 1.0)
	return frameWidthScaled, frameHeightScaled
}

// Returns scaled RectCoords(x1,x2,y1,y2) in same order
func (p *painter) scaleRectCoords(x1, x2, y1, y2 float32) (float32, float32, float32, float32) {
	x1Scaled := roundToPixel(x1*p.pixScale, 1.0)
	x2Scaled := roundToPixel(x2*p.pixScale, 1.0)
	y1Scaled := roundToPixel(y1*p.pixScale, 1.0)
	y2Scaled := roundToPixel(y2*p.pixScale, 1.0)
	return x1Scaled, x2Scaled, y1Scaled, y2Scaled
}
