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

/* deactivate temporary
func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
		return
	}
	p.drawTextureWithDetails(rect, p.newGlRectTexture, pos, rect.Size(), frame, canvas.ImageFillStretch,
		1.0, paint.VectorPad(rect))
}
*/

func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	var points []float32
	points = p.flexRectCoords(pos, rect, 0.5, frame)
	p.ctx.UseProgram(p.rectangleProgram)
	vbo := p.createBuffer(points)
	p.defineVertexArray(p.rectangleProgram, "vert", 2, 7, 0)
	p.defineVertexArray(p.rectangleProgram, "normal", 2, 7, 2)
	p.defineVertexArray(p.rectangleProgram, "colorSwitch", 1, 7, 4)
	p.defineVertexArray(p.rectangleProgram, "lineWidth", 1, 7, 5)
	p.defineVertexArray(p.rectangleProgram, "feather", 1, 7, 6)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	//println("drawRect: ", rect)
	triangleXYPoints := len(points) / 4

	var col color.Color
	if rect.StrokeColor == col {
		rect.StrokeColor = color.NRGBA{0.0, 0.0, 0.0, 0.0}
	}
	fillColorUniform := p.ctx.GetUniformLocation(p.rectangleProgram, "fill_color")
	rF, gF, bF, aF := rect.FillColor.RGBA()
	if aF == 0 {
		p.ctx.Uniform4f(fillColorUniform, 0, 0, 0, 0)
	} else {
		alphaF := float32(aF)
		colF := []float32{float32(rF) / alphaF, float32(gF) / alphaF, float32(bF) / alphaF, alphaF / 0xffff}
		p.ctx.Uniform4f(fillColorUniform, colF[0], colF[1], colF[2], colF[3])
	}
	strokeColorUniform := p.ctx.GetUniformLocation(p.rectangleProgram, "stroke_color")
	rS, gS, bS, aS := rect.StrokeColor.RGBA()
	if aS == 0 {
		p.ctx.Uniform4f(strokeColorUniform, 0, 0, 0, 0)
	} else {
		alphaS := float32(aS)
		colF := []float32{float32(rS) / alphaS, float32(gS) / alphaS, float32(bS) / alphaS, alphaS / 0xffff}
		p.ctx.Uniform4f(strokeColorUniform, colF[0], colF[1], colF[2], colF[3])
	}
	p.logError()

	p.ctx.DrawArrays(triangles, 0, triangleXYPoints)
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

func (p *painter) flexLineCoords(pos, pos1, pos2 fyne.Position, lineWidth, feather float32,
	frame fyne.Size, lineOut bool, strokeWidth float32) []float32 {
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
	//println(halfWidth, featherWidth)
	var colorType float32
	if strokeWidth == 0.0 {
		colorType = 1.0 // fillColor
	} else {
		colorType = 2.0 // strokeColor
	}

	if lineOut == true {
		return []float32{
			// coord x, y normal x, y, (fillColor or strokeColor)
			x1, y1, normalX, normalY, colorType, halfWidth, featherWidth,
			x2, y2, normalX, normalY, colorType, halfWidth, featherWidth,
			x2, y2, 0.0, 0.0, colorType, halfWidth, featherWidth,
			x2, y2, 0.0, 0.0, colorType, halfWidth, featherWidth,
			x1, y1, normalX, normalY, colorType, halfWidth, featherWidth,
			x1, y1, 0.0, 0.0, colorType, halfWidth, featherWidth,
		}
	} else {
		return []float32{
			// coord x, y normal x, y, strokeColor
			x1, y1, 0.0, 0.0, colorType, halfWidth, featherWidth,
			x2, y2, 0.0, 0.0, colorType, halfWidth, featherWidth,
			x2, y2, -normalX, -normalY, colorType, halfWidth, featherWidth,
			x2, y2, -normalX, -normalY, colorType, halfWidth, featherWidth,
			x1, y1, 0.0, 0.0, colorType, halfWidth, featherWidth,
			x1, y1, -normalX, -normalY, colorType, halfWidth, featherWidth,
		}
	}
}

func (p *painter) flexLineCoordsNew(pos, pos1, pos2 fyne.Position, lineWidth, feather float32, frame fyne.Size) ([]float32, float32, float32) {
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

	// Line calculation: y = k * x + d
	// Opposite slope of line (-k ... k_minus)
	y_lenght := pos1.Y + pos2.Y*(-1)
	x_lenght := pos1.X + pos2.X*(-1)
	k := y_lenght / x_lenght
	k_minus := k * (-1)

	// d = P (0/y)
	y_xnull := pos1.Y - (k_minus * pos1.X)
	// P (x/0)
	x_ynull := ((-1) * y_xnull) / k_minus
	// h_relation = Pythagoras of y_xnull and x_ynull
	h_rel := math.Sqrt(float64(y_xnull*y_xnull) + (float64(x_ynull * x_ynull)))
	// calculate x_dif and y_dif on ralation
	// x_ynull : h_rel = x_dif : lineWidth
	// y_xnull : h_rel = y_dif : lineWidth
	x_dif := x_ynull / float32(h_rel) * lineWidth
	y_dif := y_xnull / float32(h_rel) * lineWidth

	normalX := -1 + x_dif/frame.Width
	normalY := 1 - y_dif/frame.Height

	return []float32{
		// coord x, y normal x, y
		x1, y1, normalX, normalY,
		x2, y2, normalX, normalY,
		x2, y2, -normalX, -normalY,
		x2, y2, -normalX, -normalY,
		x1, y1, normalX, normalY,
		x1, y1, -normalX, -normalY,
	}, 0.0, 0.0
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

func (p *painter) flexRectCoords(pos fyne.Position, rect *canvas.Rectangle, feather float32, frame fyne.Size) []float32 {
	radius := rect.Radius
	size := rect.Size()
	pos1 := rect.Position()
	strokeWidth := rect.StrokeWidth

	var coords []float32
	var x1Pos, x1, y1Pos, y1, x2Pos, x2, y2Pos, y2 float32
	var leftRadius, rightRadius, leftRadiusInn, rightRadiusInn, theta, xxInn, yyInn, xxOut, yyOut float32
	// Preparations for LineCoords/Antializing
	var pos1LOut, pos2LOut fyne.Position
	var pos1LInn, pos2LInn fyne.Position
	var linePoints []float32
	//var rightRadiusInn float32

	leftRadius = radius.Left
	rightRadius = radius.Right
	if radius.Left > size.Height || radius.Left > size.Width {
		leftRadius = fyne.Min(size.Height, size.Width)
	}
	if radius.Right > size.Height || radius.Right > size.Width {
		rightRadius = fyne.Min(size.Height, size.Width)
	}

	leftRadius = roundToPixel(leftRadius, p.pixScale)
	rightRadius = roundToPixel(rightRadius, p.pixScale)
	if leftRadius < 5.0 {
		leftRadius = 0.0
	}
	if rightRadius < 5.0 {
		rightRadius = 0.0
	}

	if radius.LeftSegments == 0 {
		radius.LeftSegments = 8
	}
	if radius.RightSegments == 0 {
		radius.RightSegments = 8
	}

	//
	xPosDiff := pos.X - pos1.X
	yPosDiff := pos.Y - pos1.Y
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)
	strokeWidth = roundToPixel(strokeWidth, p.pixScale)
	var aLine, aLineRaw, aLineOneSeg, aLineOneSegRaw float32
	aLineRaw = 0.5
	//aLine = roundToPixel(aLineRaw, p.pixScale)
	aLineOneSegRaw = 0.5
	aLineOneSeg = roundToPixel(aLineOneSegRaw, p.pixScale)
	if strokeWidth < 1.0 {
		strokeWidth = 0.0
		aLine = 0.0
	}
	/*
		    ---------------------
		   / |                 | \
		  /#1|      #2         |#3\
		 /   |                 |   \
		|----|-----------------|----|
		| #4 |      #5         | #6 |
		|----|-----------------|----|
		 \   |                 |   /
		  \#7|      #8         |#9/
		   \ |                 | /
		    ---------------------
	*/
	// Slice #2 #5 #8 with stroke
	if leftRadius == 0.0 {
		x1Pos = (pos1.X + strokeWidth) / frame.Width
		x1 = -1 + x1Pos*2
	} else {
		x1Pos = (pos1.X + leftRadius) / frame.Width
		x1 = -1 + x1Pos*2
	}
	if rightRadius == 0.0 {
		x2Pos = (pos1.X + size.Width - strokeWidth) / frame.Width
		x2 = -1 + x2Pos*2
	} else {
		x2Pos = (pos1.X + size.Width - rightRadius) / frame.Width
		x2 = -1 + x2Pos*2
	}
	y1Pos = (pos1.Y + aLine + strokeWidth + aLine) / frame.Height
	y1 = 1 - y1Pos*2
	y2Pos = (pos1.Y + size.Height - aLine - strokeWidth - aLine) / frame.Height
	y2 = 1 - y2Pos*2
	coords = append(coords,
		x1, y1, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. triangle
		x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
		x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0,
		x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0, // 2. triangle
		x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
		x2, y2, 0.0, 0.0, 1.0, 0.0, 0.0)
	if strokeWidth >= 1.0 {
		// Stroke #2
		y1Pos = (pos1.Y + aLine) / frame.Height
		y1 = 1 - y1Pos*2
		y2Pos = (pos1.Y + aLine + strokeWidth + aLine) / frame.Height
		y2 = 1 - y2Pos*2
		coords = append(coords,
			x1, y1, 0.0, 0.0, 2.0, 0.0, 0.0, // 1. triangle
			x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
			x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0,
			x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0, // 2. triangle
			x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
			x2, y2, 0.0, 0.0, 2.0, 0.0, 0.0)

		// Stroke #8
		y1Pos = (pos1.Y + size.Height - aLine - strokeWidth - aLine) / frame.Height
		y1 = 1 - y1Pos*2
		y2Pos = (pos1.Y + size.Height - aLine) / frame.Height
		y2 = 1 - y2Pos*2
		coords = append(coords,
			x1, y1, 0.0, 0.0, 2.0, 0.0, 0.0, // 1. triangle
			x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
			x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0,
			x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0, // 2. triangle
			x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
			x2, y2, 0.0, 0.0, 2.0, 0.0, 0.0)
	}

	// Slice #4
	if leftRadius < size.Height*0.5 {
		x1Pos = (pos1.X + aLine + strokeWidth + aLine) / frame.Width
		x1 = -1 + x1Pos*2
		y1Pos = (pos1.Y + leftRadius) / frame.Height
		y1 = 1 - y1Pos*2
		x2Pos = (pos1.X + leftRadius) / frame.Width
		x2 = -1 + x2Pos*2
		y2Pos = (pos1.Y + size.Height - leftRadius) / frame.Height
		y2 = 1 - y2Pos*2
		if leftRadius != 0 {
			coords = append(coords,
				x1, y1, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. triangle
				x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0, // 2. triangle
				x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
				x2, y2, 0.0, 0.0, 1.0, 0.0, 0.0)
		}
		if strokeWidth >= 1.0 {
			x1Pos = (pos1.X + aLine) / frame.Width
			x1 = -1 + x1Pos*2
			x2Pos = (pos1.X + aLine + strokeWidth + aLine) / frame.Width
			x2 = -1 + x2Pos*2
			coords = append(coords,
				x1, y1, 0.0, 0.0, 2.0, 0.0, 0.0, // 1. triangle
				x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0, // 2. triangle
				x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
				x2, y2, 0.0, 0.0, 2.0, 0.0, 0.0)
		}
	}

	// Slice #6
	if rightRadius < size.Height*0.5 {
		x1Pos := (pos1.X + size.Width - rightRadius) / frame.Width
		x1 = -1 + x1Pos*2
		y1Pos = (pos1.Y + rightRadius) / frame.Height
		y1 = 1 - y1Pos*2
		x2Pos = (pos1.X + size.Width - aLine - strokeWidth - aLine) / frame.Width
		x2 = -1 + x2Pos*2
		y2Pos = (pos1.Y + size.Height - rightRadius) / frame.Height
		y2 = 1 - y2Pos*2
		if rightRadius != 0 {
			coords = append(coords,
				x1, y1, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. triangle
				x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 1.0, 0.0, 0.0, // 2. triangle
				x2, y1, 0.0, 0.0, 1.0, 0.0, 0.0,
				x2, y2, 0.0, 0.0, 1.0, 0.0, 0.0)
		}
		if strokeWidth >= 1.0 {
			x1Pos = (pos1.X + size.Width - aLine - strokeWidth - aLine) / frame.Width
			x1 = -1 + x1Pos*2
			x2Pos = (pos1.X + size.Width - aLine) / frame.Width
			x2 = -1 + x2Pos*2
			coords = append(coords,
				x1, y1, 0.0, 0.0, 2.0, 0.0, 0.0, // 1. triangle
				x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0,
				x1, y2, 0.0, 0.0, 2.0, 0.0, 0.0, // 2. triangle
				x2, y1, 0.0, 0.0, 2.0, 0.0, 0.0,
				x2, y2, 0.0, 0.0, 2.0, 0.0, 0.0)
		}
	}

	// Preparations for round corners
	leftCircleSeg := float32(radius.LeftSegments) * 4
	leftSeg4 := int32(radius.LeftSegments)
	rightCircleSeg := float32(radius.RightSegments) * 4
	rightSeg4 := int32(radius.RightSegments)
	rb_beg := int32(0)
	rb_end := int32(rightSeg4)
	lb_beg := int32(leftSeg4)
	lb_end := int32(lb_beg + leftSeg4)
	lt_beg := int32(lb_end)
	lt_end := int32(lt_beg + leftSeg4)
	rt_beg := int32(3 * rightSeg4)
	rt_end := int32(4 * rightSeg4)

	if strokeWidth >= 1.0 {
		leftRadiusInn = leftRadius - aLine - strokeWidth - aLine
		rightRadiusInn = rightRadius - aLine - strokeWidth - aLine
	} else {
		leftRadiusInn = leftRadius - aLine
		rightRadiusInn = rightRadius - aLine
	}

	var cx1Pos, cx1, cy1Pos, cy1 float32
	var x1PosInn, y1PosInn, x1Inn, y1Inn, x1PosOut, y1PosOut, x1Out, y1Out float32
	var cx2Inn, cy2Inn, cx3PosInn, cx3Inn, cy3PosInn, cy3Inn float32
	var cx2Out, cy2Out, cx3PosOut, cx3Out, cy3PosOut, cy3Out float32

	if leftRadius != 0.0 {
		// Slice #1
		// center x/y of circle
		cx1Pos = (pos1.X + leftRadius)
		cx1 = -1 + cx1Pos/frame.Width*2
		cy1Pos = (pos1.Y + leftRadius)
		cy1 = 1 - cy1Pos/frame.Height*2
		// first x/y on circle
		if strokeWidth >= 1.0 {
			// Innner
			x1PosInn = pos1.X + aLine + strokeWidth + aLine
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = pos1.Y + leftRadius
			y1Inn = 1 - y1PosInn/frame.Height*2
			// Outer
			x1PosOut = pos1.X + aLine
			x1Out = -1 + x1PosOut/frame.Width*2
			y1PosOut = pos1.Y + leftRadius
			y1Out = 1 - y1PosOut/frame.Height*2
		} else {
			// Innner
			x1PosInn = pos1.X + aLine
			if radius.LeftSegments == 1 {
				x1PosInn += aLineOneSegRaw
			}
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = pos1.Y + leftRadius
			y1Inn = 1 - y1PosInn/frame.Height*2
		}
		for i := lt_beg; i < lt_end+1; i++ {
			if i == lt_beg {
				cx2Inn = x1Inn
				cy2Inn = y1Inn
				if strokeWidth >= 1.0 {
					// Outer out
					cx2Out = x1Out
					cy2Out = y1Out
				}
				// BEG: Line Antializing 1. x/y
				if strokeWidth >= 1.0 {
					pos1LOut.X = x1PosOut
					pos1LOut.Y = y1PosOut
					pos1LInn.X = x1PosInn
					pos1LInn.Y = y1PosInn
				} else {
					pos1LOut.X = x1PosInn
					pos1LOut.Y = y1PosInn
				}
				// END: Line Antializing 1. x/y
			} else {
				theta = 2 * float32(math.Pi) * float32(i) / leftCircleSeg
				if radius.LeftSegments == 1 {
					xxInn = (leftRadiusInn - aLineOneSegRaw) * float32(math.Cos(float64(theta)))
					yyInn = (leftRadiusInn - aLineOneSegRaw) * float32(math.Sin(float64(theta)))
				} else {
					xxInn = (leftRadiusInn) * float32(math.Cos(float64(theta)))
					yyInn = (leftRadiusInn) * float32(math.Sin(float64(theta)))
				}
				cx3PosInn = xxInn + cx1Pos
				cx3Inn = -1 + cx3PosInn/frame.Width*2
				cy3PosInn = yyInn + cy1Pos
				cy3Inn = 1 - cy3PosInn/frame.Height*2
				//if i%2 == 0 {
				coords = append(coords,
					// segPx, segPy, lineNormX, lineNormY, color (1.0 = Inner/FillColor, 2.0 = Outer/StrokeColor)
					cx1, cy1, 0.0, 0.0, 1.0, 0.0, 0.0, // center x/y = const
					cx2Inn, cy2Inn, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. x/y on circle
					cx3Inn, cy3Inn, 0.0, 0.0, 1.0, 0.0, 0.0) // 2. x/y on circle
				//}
				if strokeWidth >= 1.0 {
					// Outer
					xxOut = (leftRadius - aLine) * float32(math.Cos(float64(theta)))
					yyOut = (leftRadius - aLine) * float32(math.Sin(float64(theta)))
					cx3PosOut = xxOut + cx1Pos
					cx3Out = -1 + cx3PosOut/frame.Width*2
					cy3PosOut = yyOut + cy1Pos
					cy3Out = 1 - cy3PosOut/frame.Height*2
					coords = append(coords,
						// 1. Triangle of stroke-segment
						cx2Out, cy2Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						// 2. Triangle of stroke-segment
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Inn, cy3Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
					)
					cx2Out = cx3Out
					cy2Out = cy3Out
				}
				cx2Inn = cx3Inn
				cy2Inn = cy3Inn

				// BEG: Line Antializing 2. x/y
				if strokeWidth >= 1.0 {
					pos2LOut.X = cx3PosOut
					pos2LOut.Y = cy3PosOut
					pos2LInn.X = cx3PosInn
					pos2LInn.Y = cy3PosInn
				} else {
					pos2LOut.X = cx3PosInn
					pos2LOut.Y = cy3PosInn
				}
				if radius.LeftSegments == 1 {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineOneSeg, aLineOneSeg, frame, true, strokeWidth)
				} else {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineRaw, feather, frame, true, strokeWidth)
				}
				coords = append(coords, linePoints...)
				pos1LOut = pos2LOut
				if strokeWidth >= 1.0 {
					linePoints = p.flexLineCoords(pos, pos1LInn, pos2LInn, aLineRaw, feather, frame, false, strokeWidth)
					coords = append(coords, linePoints...)
					pos1LInn = pos2LInn
				}
				// END:
			}
		}

		// Slice #7
		// center x/y of circle
		cx1Pos = (pos1.X + leftRadius)
		cx1 = -1 + cx1Pos/frame.Width*2
		cy1Pos = (pos1.Y + size.Height - leftRadius)
		cy1 = 1 - cy1Pos/frame.Height*2
		// first x/y on circle
		if strokeWidth >= 1.0 {
			// Innner
			x1PosInn = (pos1.X + leftRadius)
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + size.Height - aLine - strokeWidth - aLine)
			y1Inn = 1 - y1PosInn/frame.Height*2
			// Outer
			x1PosOut = (pos1.X + leftRadius)
			x1Out = -1 + x1PosOut/frame.Width*2
			y1PosOut = (pos1.Y + size.Height - aLine)
			y1Out = 1 - y1PosOut/frame.Height*2
		} else {
			// Innner
			x1PosInn = (pos1.X + leftRadius)
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + size.Height - aLine)
			if radius.LeftSegments == 1 {
				y1PosInn -= aLineOneSegRaw
			}
			y1Inn = 1 - y1PosInn/frame.Height*2
		}
		for i := lb_beg; i < lb_end+1; i++ {
			if i == lb_beg {
				cx2Inn = x1Inn
				cy2Inn = y1Inn
				if strokeWidth >= 1.0 {
					cx2Out = x1Out
					cy2Out = y1Out
				}
				// BEG: Line Antializing 1. x/y
				if strokeWidth >= 1.0 {
					pos1LOut.X = x1PosOut
					pos1LOut.Y = y1PosOut
					pos1LInn.X = x1PosInn
					pos1LInn.Y = y1PosInn
				} else {
					pos1LOut.X = x1PosInn
					pos1LOut.Y = y1PosInn
				}
				// END: Line Antializing 1. x/y
			} else {
				theta = 2 * float32(math.Pi) * float32(i) / leftCircleSeg
				if radius.LeftSegments == 1 {
					xxInn = (leftRadiusInn - aLineOneSegRaw) * float32(math.Cos(float64(theta)))
					yyInn = (leftRadiusInn - aLineOneSegRaw) * float32(math.Sin(float64(theta)))
				} else {
					xxInn = (leftRadiusInn) * float32(math.Cos(float64(theta)))
					yyInn = (leftRadiusInn) * float32(math.Sin(float64(theta)))
				}
				cx3PosInn = xxInn + cx1Pos
				cx3Inn = -1 + cx3PosInn/frame.Width*2
				cy3PosInn = yyInn + cy1Pos
				cy3Inn = 1 - cy3PosInn/frame.Height*2
				coords = append(coords,
					// segPx, segPy, lineNormX, lineNormY, color (1.0 = Inner/FillColor, 2.0 = Outer/StrokeColor)
					cx1, cy1, 0.0, 0.0, 1.0, 0.0, 0.0, // center x/y = const
					cx2Inn, cy2Inn, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. x/y on circle
					cx3Inn, cy3Inn, 0.0, 0.0, 1.0, 0.0, 0.0) // 2. x/y on circle
				if strokeWidth >= 1.0 {
					xxOut = (leftRadius - aLine) * float32(math.Cos(float64(theta)))
					yyOut = (leftRadius - aLine) * float32(math.Sin(float64(theta)))
					cx3PosOut = xxOut + cx1Pos
					cx3Out = -1 + cx3PosOut/frame.Width*2
					cy3PosOut = yyOut + cy1Pos
					cy3Out = 1 - cy3PosOut/frame.Height*2
					coords = append(coords,
						// 1. Triangle of stroke-segment
						cx2Out, cy2Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						// 2. Triangle of stroke-segment
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Inn, cy3Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
					)
					cx2Out = cx3Out
					cy2Out = cy3Out
				}
				cx2Inn = cx3Inn
				cy2Inn = cy3Inn

				// BEG: Line Antializing 2. x/y
				if strokeWidth >= 1.0 {
					pos2LOut.X = cx3PosOut
					pos2LOut.Y = cy3PosOut
					pos2LInn.X = cx3PosInn
					pos2LInn.Y = cy3PosInn
				} else {
					pos2LOut.X = cx3PosInn
					pos2LOut.Y = cy3PosInn
				}
				if radius.LeftSegments == 1 {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineOneSeg, aLineOneSeg, frame, true, strokeWidth)
				} else {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineRaw, feather, frame, true, strokeWidth)
				}
				coords = append(coords, linePoints...)
				pos1LOut = pos2LOut
				if strokeWidth >= 1.0 {
					linePoints = p.flexLineCoords(pos, pos1LInn, pos2LInn, aLineRaw, feather, frame, false, strokeWidth)
					coords = append(coords, linePoints...)
					pos1LInn = pos2LInn
				}
				// END:
			}
		}
	}

	if rightRadius != 0.0 {
		// Slice #3
		// center x/y of circle
		cx1Pos = (pos1.X + size.Width - rightRadius)
		cx1 = -1 + cx1Pos/frame.Width*2
		cy1Pos = (pos1.Y + rightRadius)
		cy1 = 1 - cy1Pos/frame.Height*2
		// first x/y on circle
		if strokeWidth >= 1.0 {
			// Innner
			x1PosInn = (pos1.X + size.Width - rightRadius)
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + aLine + strokeWidth + aLine)
			y1Inn = 1 - y1PosInn/frame.Height*2
			// Outer
			x1PosOut = (pos1.X + size.Width - rightRadius)
			x1Out = -1 + x1PosOut/frame.Width*2
			y1PosOut = (pos1.Y + aLine)
			y1Out = 1 - y1PosOut/frame.Height*2
		} else {
			// Innner
			x1PosInn = (pos1.X + size.Width - rightRadius)
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + aLine)
			if radius.RightSegments == 1 {
				y1PosInn += aLineOneSegRaw
			}
			y1Inn = 1 - y1PosInn/frame.Height*2
		}
		for i := rt_beg; i < rt_end+1; i++ {
			if i == rt_beg {
				cx2Inn = x1Inn
				cy2Inn = y1Inn
				if strokeWidth >= 1.0 {
					cx2Out = x1Out
					cy2Out = y1Out
				}
				// BEG: Line Antializing 1. x/y
				if strokeWidth >= 1.0 {
					pos1LOut.X = x1PosOut
					pos1LOut.Y = y1PosOut
					pos1LInn.X = x1PosInn
					pos1LInn.Y = y1PosInn
				} else {
					pos1LOut.X = x1PosInn
					pos1LOut.Y = y1PosInn
				}
				// END: Line Antializing 1. x/y
			} else {
				theta = 2 * float32(math.Pi) * float32(i) / rightCircleSeg
				if radius.RightSegments == 1 {
					xxInn = (rightRadiusInn - aLineOneSegRaw) * float32(math.Cos(float64(theta)))
					yyInn = (rightRadiusInn - aLineOneSegRaw) * float32(math.Sin(float64(theta)))
				} else {
					xxInn = (rightRadiusInn) * float32(math.Cos(float64(theta)))
					yyInn = (rightRadiusInn) * float32(math.Sin(float64(theta)))
				}
				cx3PosInn = xxInn + cx1Pos
				cx3Inn = -1 + cx3PosInn/frame.Width*2
				cy3PosInn = yyInn + cy1Pos
				cy3Inn = 1 - cy3PosInn/frame.Height*2
				coords = append(coords,
					// segPx, segPy, lineNormX, lineNormY, color (1.0 = Inner/FillColor, 2.0 = Outer/StrokeColor)
					cx1, cy1, 0.0, 0.0, 1.0, 0.0, 0.0, // center x/y = const
					cx2Inn, cy2Inn, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. x/y on circle
					cx3Inn, cy3Inn, 0.0, 0.0, 1.0, 0.0, 0.0) // 2. x/y on circle
				if strokeWidth >= 1.0 {
					xxOut = (rightRadius - aLine) * float32(math.Cos(float64(theta)))
					yyOut = (rightRadius - aLine) * float32(math.Sin(float64(theta)))
					cx3PosOut = xxOut + cx1Pos
					cx3Out = -1 + cx3PosOut/frame.Width*2
					cy3PosOut = yyOut + cy1Pos
					cy3Out = 1 - cy3PosOut/frame.Height*2
					coords = append(coords,
						// 1. Triangle of stroke-segment
						cx2Out, cy2Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						// 2. Triangle of stroke-segment
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Inn, cy3Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
					)
					cx2Out = cx3Out
					cy2Out = cy3Out
				}
				cx2Inn = cx3Inn
				cy2Inn = cy3Inn

				// BEG: Line Antializing 2. x/y
				if strokeWidth >= 1.0 {
					pos2LOut.X = cx3PosOut
					pos2LOut.Y = cy3PosOut
					pos2LInn.X = cx3PosInn
					pos2LInn.Y = cy3PosInn
				} else {
					pos2LOut.X = cx3PosInn
					pos2LOut.Y = cy3PosInn
				}
				if radius.RightSegments == 1 {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineOneSeg, aLineOneSeg, frame, true, strokeWidth)
				} else {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineRaw, feather, frame, true, strokeWidth)
				}
				coords = append(coords, linePoints...)
				pos1LOut = pos2LOut
				if strokeWidth >= 1.0 {
					linePoints = p.flexLineCoords(pos, pos1LInn, pos2LInn, aLineRaw, feather, frame, false, strokeWidth)
					coords = append(coords, linePoints...)
					pos1LInn = pos2LInn
				}
				// END:
			}
		}

		// Slice #9
		// center x/y of circle
		cx1Pos = (pos1.X + size.Width - rightRadius)
		cx1 = -1 + cx1Pos/frame.Width*2
		cy1Pos = (pos1.Y + size.Height - rightRadius)
		cy1 = 1 - cy1Pos/frame.Height*2
		// first x/y on circle
		if strokeWidth >= 1.0 {
			// Innner
			x1PosInn = (pos1.X + size.Width - aLine - strokeWidth - aLine)
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + size.Height - rightRadius)
			y1Inn = 1 - y1PosInn/frame.Height*2
			// Outer
			x1PosOut = (pos1.X + size.Width - aLine)
			x1Out = -1 + x1PosOut/frame.Width*2
			y1PosOut = (pos1.Y + size.Height - rightRadius)
			y1Out = 1 - y1PosOut/frame.Height*2
		} else {
			// Innner
			x1PosInn = (pos1.X + size.Width - aLine)
			if radius.RightSegments == 1 {
				x1PosInn -= aLineOneSegRaw
			}
			x1Inn = -1 + x1PosInn/frame.Width*2
			y1PosInn = (pos1.Y + size.Height - rightRadius)
			y1Inn = 1 - y1PosInn/frame.Height*2
		}
		for i := rb_beg; i < rb_end+1; i++ {
			if i == rb_beg {
				cx2Inn = x1Inn
				cy2Inn = y1Inn
				if strokeWidth >= 1.0 {
					cx2Out = x1Out
					cy2Out = y1Out
				}
				// BEG: Line Antializing 1. x/y
				if strokeWidth >= 1.0 {
					pos1LOut.X = x1PosOut
					pos1LOut.Y = y1PosOut
					pos1LInn.X = x1PosInn
					pos1LInn.Y = y1PosInn
				} else {
					pos1LOut.X = x1PosInn
					pos1LOut.Y = y1PosInn
				}
				// END: Line Antializing 1. x/y
			} else {
				theta = 2 * float32(math.Pi) * float32(i) / rightCircleSeg
				if radius.RightSegments == 1 {
					xxInn = (rightRadiusInn - aLineOneSegRaw) * float32(math.Cos(float64(theta)))
					yyInn = (rightRadiusInn - aLineOneSegRaw) * float32(math.Sin(float64(theta)))
				} else {
					xxInn = (rightRadiusInn) * float32(math.Cos(float64(theta)))
					yyInn = (rightRadiusInn) * float32(math.Sin(float64(theta)))
				}
				cx3PosInn = xxInn + cx1Pos
				cx3Inn = -1 + cx3PosInn/frame.Width*2
				cy3PosInn = yyInn + cy1Pos
				cy3Inn = 1 - cy3PosInn/frame.Height*2
				coords = append(coords,
					// segPx, segPy, lineNormX, lineNormY, color (1.0 = Inner/FillColor, 2.0 = Outer/StrokeColor)
					cx1, cy1, 0.0, 0.0, 1.0, 0.0, 0.0, // center x/y = const
					cx2Inn, cy2Inn, 0.0, 0.0, 1.0, 0.0, 0.0, // 1. x/y on circle
					cx3Inn, cy3Inn, 0.0, 0.0, 1.0, 0.0, 0.0) // 2. x/y on circle
				if strokeWidth >= 1.0 {
					xxOut = (rightRadius - aLine) * float32(math.Cos(float64(theta)))
					yyOut = (rightRadius - aLine) * float32(math.Sin(float64(theta)))
					cx3PosOut = xxOut + cx1Pos
					cx3Out = -1 + cx3PosOut/frame.Width*2
					cy3PosOut = yyOut + cy1Pos
					cy3Out = 1 - cy3PosOut/frame.Height*2
					coords = append(coords,
						// 1. Triangle of stroke-segment
						cx2Out, cy2Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						// 2. Triangle of stroke-segment
						cx3Out, cy3Out, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx2Inn, cy2Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
						cx3Inn, cy3Inn, 0.0, 0.0, 2.0, 0.0, 0.0,
					)
					cx2Out = cx3Out
					cy2Out = cy3Out
				}
				cx2Inn = cx3Inn
				cy2Inn = cy3Inn

				// BEG: Line Antializing 2. x/y
				if strokeWidth >= 1.0 {
					pos2LOut.X = cx3PosOut
					pos2LOut.Y = cy3PosOut
					pos2LInn.X = cx3PosInn
					pos2LInn.Y = cy3PosInn
				} else {
					pos2LOut.X = cx3PosInn
					pos2LOut.Y = cy3PosInn
				}
				if radius.RightSegments == 1 {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineOneSeg, aLineOneSeg, frame, true, strokeWidth)
				} else {
					linePoints = p.flexLineCoords(pos, pos1LOut, pos2LOut, aLineRaw, feather, frame, true, strokeWidth)
				}
				coords = append(coords, linePoints...)
				pos1LOut = pos2LOut
				if strokeWidth >= 1.0 {
					linePoints = p.flexLineCoords(pos, pos1LInn, pos2LInn, aLineRaw, feather, frame, false, strokeWidth)
					coords = append(coords, linePoints...)
					pos1LInn = pos2LInn
				}
				// END:
			}
		}
	}

	return coords
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
