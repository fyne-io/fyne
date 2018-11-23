// +build !ci,gl

package gl

import (
	_ "image/png"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func (c *glCanvas) drawContainer(cont *fyne.Container, offset fyne.Position) {
	pos := cont.Position.Add(offset)
	for _, child := range cont.Objects {
		switch co := child.(type) {
		case *fyne.Container:
			c.drawContainer(co, pos)
		case fyne.Widget:
			c.drawWidget(co, pos)
		default:
			c.drawObject(co, pos)
		}
	}
}

func (c *glCanvas) drawWidget(w fyne.Widget, offset fyne.Position) {
	pos := w.CurrentPosition().Add(offset)
	for _, child := range w.Renderer().Objects() {
		switch co := child.(type) {
		case *fyne.Container:
			c.drawContainer(co, pos)
		case fyne.Widget:
			c.drawWidget(co, pos)
		default:
			c.drawObject(co, pos)
		}
	}
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (c *glCanvas) rectCoords(size fyne.Size, pos fyne.Position) []float32 {
	xPos := float32(pos.X) / float32(c.Size().Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width) / float32(c.Size().Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y) / float32(c.Size().Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height) / float32(c.Size().Height)
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

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position) {
	points := c.rectCoords(rect.Size, pos)
	texture := getTexture(rect, c.newGlRectTexture)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawImage(img *canvas.Image, pos fyne.Position) {
	points := c.rectCoords(img.Size, pos)
	texture := c.newGlImageTexture(img)
	if texture == 0 {
		return
	}

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawText(text *canvas.Text, pos fyne.Position) {
	if text.Text == "" {
		return
	}

	points := c.rectCoords(text.MinSize(), pos)
	texture := c.newGlTextTexture(text)

	c.drawTexture(texture, points)
}

func (c *glCanvas) drawObject(o fyne.CanvasObject, offset fyne.Position) {
	pos := o.CurrentPosition().Add(offset)
	switch obj := o.(type) {
	case *canvas.Rectangle:
		c.drawRectangle(obj, pos)
	case *canvas.Image:
		c.drawImage(obj, pos)
	case *canvas.Text:
		c.drawText(obj, pos)
	}
}
