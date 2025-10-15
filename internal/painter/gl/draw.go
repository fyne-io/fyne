package gl

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	paint "fyne.io/fyne/v2/internal/painter"
)

const edgeSoftness = 1.0

func (p *painter) createBuffer(size int) Buffer {
	vbo := p.ctx.CreateBuffer()
	p.logError()
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferData(arrayBuffer, make([]float32, size), staticDraw)
	p.logError()
	return vbo
}

func (p *painter) updateBuffer(vbo Buffer, points []float32) {
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferSubData(arrayBuffer, points)
	p.logError()
}

func (p *painter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	radius := paint.GetMaximumRadius(circle.Size())
	program := p.roundRectangleProgram

	// Vertex: BEG
	bounds, points := p.vecSquareCoords(pos, circle, frame)
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(circle.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width_half", strokeWidthScaled*0.5)

	rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
	rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
	p.SetUniform2f(program, "rect_size_half", rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

	radiusScaled := roundToPixel(radius*p.pixScale, 1.0)
	p.SetUniform4f(program, "radius", radiusScaled, radiusScaled, radiusScaled, radiusScaled)

	r, g, b, a := getFragmentColor(circle.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := circle.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	edgeSoftnessScaled := roundToPixel(edgeSoftness*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
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
	p.ctx.UseProgram(p.lineProgram.ref)
	p.updateBuffer(p.lineProgram.buff, points)
	p.UpdateVertexArray(p.lineProgram, "vert", 2, 4, 0)
	p.UpdateVertexArray(p.lineProgram, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	r, g, b, a := getFragmentColor(line.StrokeColor)
	p.SetUniform4f(p.lineProgram, "color", r, g, b, a)

	p.SetUniform1f(p.lineProgram, "lineWidth", halfWidth)

	p.SetUniform1f(p.lineProgram, "feather", feather)

	p.ctx.DrawArrays(triangles, 0, 6)
	p.logError()
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
	case *canvas.Polygon:
		p.drawPolygon(obj, pos, frame)
	case *canvas.Arc:
		p.drawArc(obj, pos, frame)
	}
}

func (p *painter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
}

func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	topRightRadius := paint.GetCornerRadius(rect.TopRightCornerRadius, rect.CornerRadius)
	topLeftRadius := paint.GetCornerRadius(rect.TopLeftCornerRadius, rect.CornerRadius)
	bottomRightRadius := paint.GetCornerRadius(rect.BottomRightCornerRadius, rect.CornerRadius)
	bottomLeftRadius := paint.GetCornerRadius(rect.BottomLeftCornerRadius, rect.CornerRadius)
	p.drawOblong(rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rect.Aspect, pos, frame)
}

func (p *painter) drawOblong(obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, aspect float32, pos fyne.Position, frame fyne.Size) {
	if (fill == color.Transparent || fill == nil) && (stroke == color.Transparent || stroke == nil || strokeWidth == 0) {
		return
	}

	roundedCorners := topRightRadius != 0 || topLeftRadius != 0 || bottomRightRadius != 0 || bottomLeftRadius != 0
	var program ProgramState
	if roundedCorners {
		program = p.roundRectangleProgram
	} else {
		program = p.rectangleProgram
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, obj, frame, aspect)
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(strokeWidth*p.pixScale, 1.0)
	if roundedCorners {
		p.SetUniform1f(program, "stroke_width_half", strokeWidthScaled*0.5)

		rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
		rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
		p.SetUniform2f(program, "rect_size_half", rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

		// the maximum possible corner radii for a circular shape, calculated taking into account the rect coords with aspect ratio
		size := fyne.NewSize(bounds[2]-bounds[0], bounds[3]-bounds[1])
		topRightRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(topRightRadius, topLeftRadius, bottomRightRadius, size)*p.pixScale,
			1.0,
		)
		topLeftRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(topLeftRadius, topRightRadius, bottomLeftRadius, size)*p.pixScale,
			1.0,
		)
		bottomRightRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(bottomRightRadius, bottomLeftRadius, topRightRadius, size)*p.pixScale,
			1.0,
		)
		bottomLeftRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(bottomLeftRadius, bottomRightRadius, topLeftRadius, size)*p.pixScale,
			1.0,
		)
		p.SetUniform4f(program, "radius", topRightRadiusScaled, bottomRightRadiusScaled, topLeftRadiusScaled, bottomLeftRadiusScaled)

		edgeSoftnessScaled := roundToPixel(edgeSoftness*p.pixScale, 1.0)
		p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)
	} else {
		p.SetUniform1f(program, "stroke_width", strokeWidthScaled)
	}

	r, g, b, a := getFragmentColor(fill)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := stroke
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawPolygon(polygon *canvas.Polygon, pos fyne.Position, frame fyne.Size) {
	if ((polygon.FillColor == color.Transparent || polygon.FillColor == nil) && (polygon.StrokeColor == color.Transparent || polygon.StrokeColor == nil || polygon.StrokeWidth == 0)) || polygon.Sides < 3 {
		return
	}
	size := polygon.Size()

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, polygon, frame, 0.0)
	program := p.polygonProgram
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	edgeSoftnessScaled := roundToPixel(edgeSoftness*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)

	outerRadius := fyne.Min(size.Width, size.Height) / 2
	outerRadiusScaled := roundToPixel(outerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "outer_radius", outerRadiusScaled)

	p.SetUniform1f(program, "angle", polygon.Angle)
	p.SetUniform1f(program, "sides", float32(polygon.Sides))

	cornerRadius := fyne.Min(paint.GetMaximumRadius(size), polygon.CornerRadius)
	cornerRadiusScaled := roundToPixel(cornerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "corner_radius", cornerRadiusScaled)

	strokeWidthScaled := roundToPixel(polygon.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width", strokeWidthScaled)

	r, g, b, a := getFragmentColor(polygon.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := polygon.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawArc(arc *canvas.Arc, pos fyne.Position, frame fyne.Size) {
	if ((arc.FillColor == color.Transparent || arc.FillColor == nil) && (arc.StrokeColor == color.Transparent || arc.StrokeColor == nil || arc.StrokeWidth == 0)) || arc.StartAngle == arc.EndAngle {
		return
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, arc, frame, 0.0)
	program := p.arcProgram
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	edgeSoftnessScaled := roundToPixel(edgeSoftness*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)

	outerRadius := fyne.Min(arc.Size().Width, arc.Size().Height) / 2
	outerRadiusScaled := roundToPixel(outerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "outer_radius", outerRadiusScaled)

	innerRadius := outerRadius * float32(math.Min(1.0, math.Max(0.0, float64(arc.CutoutRatio))))
	innerRadiusScaled := roundToPixel(innerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "inner_radius", innerRadiusScaled)

	startAngle, endAngle := paint.NormalizeArcAngles(arc.StartAngle, arc.EndAngle)
	p.SetUniform1f(program, "start_angle", startAngle)
	p.SetUniform1f(program, "end_angle", endAngle)

	cornerRadius := fyne.Min(paint.GetMaximumRadiusArc(outerRadius, innerRadius, arc.EndAngle-arc.StartAngle), arc.CornerRadius)
	cornerRadiusScaled := roundToPixel(cornerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "corner_radius", cornerRadiusScaled)

	strokeWidthScaled := roundToPixel(arc.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width", strokeWidthScaled)

	r, g, b, a := getFragmentColor(arc.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := arc.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
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
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32, pad float32,
) {
	texture, err := p.getTexture(o, creator)
	if err != nil {
		return
	}

	cornerRadius := float32(0)
	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = img.Aspect()
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
		if img.CornerRadius > 0 {
			cornerRadius = img.CornerRadius
		}
	}
	points := p.rectCoords(size, pos, frame, fill, aspect, pad)
	inner, _ := rectInnerCoords(size, pos, fill, aspect)

	p.ctx.UseProgram(p.program.ref)
	p.updateBuffer(p.program.buff, points)
	p.UpdateVertexArray(p.program, "vert", 3, 5, 0)
	p.UpdateVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	// Set corner radius and texture size in pixels
	cornerRadius = fyne.Min(paint.GetMaximumRadius(size), cornerRadius)
	p.SetUniform1f(p.program, "cornerRadius", cornerRadius*p.pixScale)
	p.SetUniform2f(p.program, "size", inner.Width*p.pixScale, inner.Height*p.pixScale)

	p.SetUniform1f(p.program, "alpha", alpha)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
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
	fill canvas.ImageFill, aspect float32, pad float32,
) []float32 {
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

	xInset := float32(0.0)
	yInset := float32(0.0)

	if fill == canvas.ImageFillCover {
		viewAspect := size.Width / size.Height

		if viewAspect > aspect {
			newHeight := size.Width / aspect
			heightPad := (newHeight - size.Height) / 2
			yInset = heightPad / newHeight
		} else if viewAspect < aspect {
			newWidth := size.Height * aspect
			widthPad := (newWidth - size.Width) / 2
			xInset = widthPad / newWidth
		}
	}

	return []float32{
		// coord x, y, z texture x, y
		x1, y2, 0, xInset, 1.0 - yInset, // top left
		x1, y1, 0, xInset, yInset, // bottom left
		x2, y2, 0, 1.0 - xInset, 1.0 - yInset, // top right
		x2, y1, 0, 1.0 - xInset, yInset, // bottom right
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

func (p *painter) vecRectCoords(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size, aspect float32) ([4]float32, []float32) {
	xPad, yPad := float32(0), float32(0)

	if aspect != 0 {
		inner := rect.Size()
		frameAspect := inner.Width / inner.Height

		if frameAspect > aspect {
			newWidth := inner.Height * aspect
			xPad = (inner.Width - newWidth) / 2
		} else if frameAspect < aspect {
			newHeight := inner.Width / aspect
			yPad = (inner.Height - newHeight) / 2
		}
	}

	return p.vecRectCoordsWithPad(pos, rect, frame, xPad, yPad)
}

func (p *painter) vecRectCoordsWithPad(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size, xPad, yPad float32) ([4]float32, []float32) {
	size := rect.Size()
	pos1 := rect.Position()

	xPosDiff := pos.X - pos1.X + xPad
	yPosDiff := pos.Y - pos1.Y + yPad
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width-2*xPad, p.pixScale)
	size.Height = roundToPixel(size.Height-2*yPad, p.pixScale)

	// without edge softness adjustment the rectangle has cropped edges
	edgeSoftnessScaled := roundToPixel(edgeSoftness*p.pixScale, 1.0)
	x1Pos := pos1.X
	x1Norm := -1 + (x1Pos-edgeSoftnessScaled)*2/frame.Width
	x2Pos := pos1.X + size.Width
	x2Norm := -1 + (x2Pos+edgeSoftnessScaled)*2/frame.Width
	y1Pos := pos1.Y
	y1Norm := 1 - (y1Pos-edgeSoftnessScaled)*2/frame.Height
	y2Pos := pos1.Y + size.Height
	y2Norm := 1 - (y2Pos+edgeSoftnessScaled)*2/frame.Height

	// output a norm for the fill and the vert is unused, but we pass 0 to avoid optimisation issues
	coords := []float32{
		0, 0, x1Norm, y1Norm, // first triangle
		0, 0, x2Norm, y1Norm, // second triangle
		0, 0, x1Norm, y2Norm,
		0, 0, x2Norm, y2Norm,
	}

	return [4]float32{x1Pos, y1Pos, x2Pos, y2Pos}, coords
}

func (p *painter) vecSquareCoords(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size) ([4]float32, []float32) {
	return p.vecRectCoordsWithPad(pos, rect, frame, 0, 0)
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
