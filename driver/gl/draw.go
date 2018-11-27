// +build !ci,gl

package gl

import (
	_ "image/png"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func walkObjects(obj fyne.CanvasObject, pos fyne.Position,
	f func(object fyne.CanvasObject, pos fyne.Position)) {

	switch co := obj.(type) {
	case *fyne.Container:
		offset := co.Position.Add(pos)
		f(obj, offset)

		for _, child := range co.Objects {
			walkObjects(child, offset, f)
		}
	case fyne.Widget:
		offset := co.CurrentPosition().Add(pos)
		f(obj, offset)

		for _, child := range co.Renderer().Objects() {
			walkObjects(child, offset, f)
		}
	default:
		f(obj, pos)
	}
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (c *glCanvas) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size) []float32 {
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

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if !rect.IsVisible() {
		return
	}

	points := c.rectCoords(rect.Size, pos, frame)
	texture := getTexture(rect, c.newGlRectTexture)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	if !img.IsVisible() {
		return
	}

	points := c.rectCoords(img.Size, pos, frame)
	texture := c.newGlImageTexture(img)
	if texture == 0 {
		return
	}

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if !text.IsVisible() || text.Text == "" {
		return
	}

	size := text.MinSize()
	containerSize := text.CurrentSize()
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		pos = fyne.NewPos(pos.X + containerSize.Width - size.Width, pos.Y)
	case fyne.TextAlignCenter:
		pos = fyne.NewPos(pos.X + (containerSize.Width - size.Width)/2, pos.Y)
	}

	if text.CurrentSize().Height > text.MinSize().Height {
		pos = fyne.NewPos(pos.X, pos.Y + (text.CurrentSize().Height-text.MinSize().Height)/2)
	}

	points := c.rectCoords(size, pos, frame)
	texture := c.newGlTextTexture(text)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawObject(o fyne.CanvasObject, offset fyne.Position, frame fyne.Size) {
	pos := o.CurrentPosition().Add(offset)
	switch obj := o.(type) {
	case *canvas.Rectangle:
		c.drawRectangle(obj, pos, frame)
	case *canvas.Image:
		c.drawImage(obj, pos, frame)
	case *canvas.Text:
		c.drawText(obj, pos, frame)
	}
}
