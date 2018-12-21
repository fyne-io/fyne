package gl

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/widget"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func walkObjects(obj fyne.CanvasObject, pos fyne.Position,
	f func(object fyne.CanvasObject, pos fyne.Position)) {

	switch co := obj.(type) {
	case *fyne.Container:
		offset := co.Position().Add(pos)
		f(obj, offset)

		for _, child := range co.Objects {
			walkObjects(child, offset, f)
		}
	case fyne.Widget:
		offset := co.Position().Add(pos)
		f(obj, offset)

		for _, child := range widget.Renderer(co).Objects() {
			walkObjects(child, offset, f)
		}
	default:
		f(obj, pos)
	}
}

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill == canvas.ImageFillContain {
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
func (c *glCanvas) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size, fill canvas.ImageFill, aspect float32) []float32 {
	size, pos = rectInnerCoords(size, pos, fill, aspect)

	xPos := float32(pos.X) / float32(frame.Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width) / float32(frame.Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y) / float32(frame.Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height) / float32(frame.Height)
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

	return points
}

func (c *glCanvas) drawTexture(texture uint32, points []float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
}

func (c *glCanvas) drawWidget(box fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if !box.Visible() {
		return
	}

	points := c.rectCoords(box.Size(), pos, frame, canvas.ImageFillStretch, 0.0)
	texture := getTexture(box, c.newGlRectTexture)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if !rect.Visible() {
		return
	}

	points := c.rectCoords(rect.Size(), pos, frame, canvas.ImageFillStretch, 0.0)
	texture := getTexture(rect, c.newGlRectTexture)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	if !img.Visible() {
		return
	}

	texture := getTexture(img, c.newGlImageTexture)
	if texture == 0 {
		return
	}

	points := c.rectCoords(img.Size(), pos, frame, img.FillMode, img.PixelAspect)
	c.drawTexture(texture, points)
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

	points := c.rectCoords(size, pos, frame, canvas.ImageFillStretch, 0.0)
	texture := c.newGlTextTexture(text)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawObject(o fyne.CanvasObject, offset fyne.Position, frame fyne.Size) {
	canvasMutex.Lock()
	canvases[o] = c
	canvasMutex.Unlock()
	pos := o.Position().Add(offset)
	switch obj := o.(type) {
	case *canvas.Rectangle:
		c.drawRectangle(obj, pos, frame)
	case *canvas.Image:
		c.drawImage(obj, pos, frame)
	case *canvas.Text:
		c.drawText(obj, pos, frame)
	case fyne.Widget:
		c.drawWidget(obj, offset, frame)
	}
}
