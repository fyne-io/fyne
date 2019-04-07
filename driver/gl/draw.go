package gl

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func (c *glCanvas) walkObjects(obj fyne.CanvasObject, pos fyne.Position,
	f func(object fyne.CanvasObject, pos fyne.Position)) {

	switch co := obj.(type) {
	case *fyne.Container:
		offset := co.Position().Add(pos)
		f(obj, offset)

		for _, child := range co.Objects {
			c.walkObjects(child, offset, f)
		}
	case *widget.ScrollContainer: // TODO should this be somehow not scroll container specific?
		offset := co.Position().Add(pos)

		scrollX := textureScaleInt(c, offset.X)
		scrollY := textureScaleInt(c, offset.Y)
		scrollWidth := textureScaleInt(c, co.Size().Width)
		scrollHeight := textureScaleInt(c, co.Size().Height)
		_, pixHeight := c.window.viewport.GetFramebufferSize()
		gl.Scissor(int32(scrollX), int32(pixHeight-scrollY-scrollHeight), int32(scrollWidth), int32(scrollHeight))
		gl.Enable(gl.SCISSOR_TEST)

		f(obj, offset)

		for _, child := range widget.Renderer(co).Objects() {
			c.walkObjects(child, offset, f)
		}

		gl.Disable(gl.SCISSOR_TEST)
	case fyne.Widget:
		offset := co.Position().Add(pos)
		f(obj, offset)

		for _, child := range widget.Renderer(co).Objects() {
			c.walkObjects(child, offset, f)
		}
	default:
		f(obj, pos)
	}
}

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill == canvas.ImageFillContain || fill == canvas.ImageFillOriginal {
		// change pos and size accordingly

		viewAspect := float32(size.Width) / float32(size.Height)

		newWidth, newHeight := size.Width, size.Height
		widthPad, heightPad := 0, 0
		if viewAspect > aspect {
			newWidth = int(float32(size.Height) * aspect)
			widthPad = (size.Width - newWidth) / 2
		} else if viewAspect < aspect {
			newHeight = int(float32(size.Width) / aspect)
			heightPad = (size.Height - newHeight) / 2
		}

		return fyne.NewSize(newWidth, newHeight), fyne.NewPos(pos.X+widthPad, pos.Y+heightPad)
	}

	return size, pos
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (c *glCanvas) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad int) ([]float32, uint32, uint32) {
	size, pos = rectInnerCoords(size, pos, fill, aspect)

	xPos := float32(pos.X-pad) / float32(frame.Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width+pad) / float32(frame.Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y-pad) / float32(frame.Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height+pad) / float32(frame.Height)
	y2 := 1 - y2Pos*2

	points := []float32{
		// coord x, y, x texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	textureUniform := gl.GetUniformLocation(c.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	vertAttrib := uint32(gl.GetAttribLocation(c.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(c.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return points, vao, vbo
}

func (c *glCanvas) freeCoords(vao, vbo uint32) {
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vbo)
}

func (c *glCanvas) drawTexture(texture uint32, points []float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
}

func (c *glCanvas) drawWidget(box fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	backCol := widget.Renderer(box.(fyne.Widget)).BackgroundColor()
	if !box.Visible() || backCol == color.Transparent {
		return
	}

	points, vao, vbo := c.rectCoords(box.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(box, c.newGlRectTexture)

	gl.Enable(gl.BLEND)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	if !circle.Visible() {
		return
	}

	points, vao, vbo := c.rectCoords(circle.Size(), pos, frame, canvas.ImageFillStretch, 0.0, vectorPad)
	texture := getTexture(circle, c.newGlCircleTexture)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	if !line.Visible() {
		return
	}

	points, vao, vbo := c.rectCoords(line.Size(), pos, frame, canvas.ImageFillStretch, 0.0, vectorPad)
	texture := getTexture(line, c.newGlLineTexture)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	if !img.Visible() {
		return
	}

	texture := getTexture(img, c.newGlImageTexture)
	if texture == 0 {
		return
	}

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	gl.Enable(gl.BLEND)
	if img.Alpha() != 1 {
		gl.BlendColor(0, 0, 0, float32(img.Alpha()))
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}

	aspect := aspects[img.Resource]
	if aspect == 0 {
		aspect = aspects[img]
	}
	points, vao, vbo := c.rectCoords(img.Size(), pos, frame, img.FillMode, aspect, 0)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	if !img.Visible() {
		return
	}

	texture := getTexture(img, c.newGlRasterTexture)
	if texture == 0 {
		return
	}

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	gl.Enable(gl.BLEND)
	if img.Alpha() != 1 {
		gl.BlendColor(0, 0, 0, float32(img.Alpha()))
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	points, vao, vbo := c.rectCoords(img.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if !rect.Visible() {
		return
	}

	points, vao, vbo := c.rectCoords(rect.Size(), pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(rect, c.newGlRectTexture)

	gl.Enable(gl.BLEND) // enable translucency
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if !text.Visible() || text.Text == "" {
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

	if text.Size().Height > text.MinSize().Height {
		pos = fyne.NewPos(pos.X, pos.Y+(text.Size().Height-text.MinSize().Height)/2)
	}

	points, vao, vbo := c.rectCoords(size, pos, frame, canvas.ImageFillStretch, 0.0, 0)
	texture := getTexture(text, c.newGlTextTexture)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	c.drawTexture(texture, points)
	c.freeCoords(vao, vbo)
}

func (c *glCanvas) drawObject(o fyne.CanvasObject, offset fyne.Position, frame fyne.Size) {
	canvasMutex.Lock()
	canvases[o] = c
	canvasMutex.Unlock()
	pos := o.Position().Add(offset)
	switch obj := o.(type) {
	case *canvas.Circle:
		c.drawCircle(obj, pos, frame)
	case *canvas.Line:
		c.drawLine(obj, pos, frame)
	case *canvas.Image:
		c.drawImage(obj, pos, frame)
	case *canvas.Raster:
		c.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		c.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		c.drawText(obj, pos, frame)
	case fyne.Widget:
		c.drawWidget(obj, offset, frame)
	}
}
