// +build !ci,gl

package gl

import (
	"fmt"
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/go-gl/gl/v3.3-core/gl"
	"image/color"
	"strings"
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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

const (
	vertexShaderSource = `
    #version 130
    in vec3 vp;
    uniform vec4 obj_fill_color;
    out vec4 fill_color;
    void main() {
        fill_color = obj_fill_color;
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 130
    in vec4 fill_color;
    out vec4 frag_colour;
    void main() {
        frag_colour = fill_color;
    }
` + "\x00"
)

func square(size fyne.Size, pos fyne.Position, full fyne.Size) []float32 {
	xPos := float32(pos.X) / float32(full.Width)
	x1 := -1 + xPos*2
	x2Pos := float32(pos.X+size.Width) / float32(full.Width)
	x2 := -1 + x2Pos*2

	yPos := float32(pos.Y) / float32(full.Height)
	y1 := 1 - yPos*2
	y2Pos := float32(pos.Y+size.Height) / float32(full.Height)
	y2 := 1 - y2Pos*2

	return []float32{
		x1, y1, 0,
		x1, y2, 0,
		x2, y2, 0,

		x1, y1, 0,
		x2, y1, 0,
		x2, y2, 0,
	}
}

func (d *gLDriver) initOpenGL() {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	d.program = prog
}

// makeVao initializes and returns a vertex array from the points provided.
func (c *glCanvas) makeVao(points []float32, fill color.Color) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	r, g, b, a := fill.RGBA()
	loc := gl.GetUniformLocation(c.program, gl.Str("obj_fill_color\x00"))
	gl.Uniform4f(loc, float32(uint8(r))/255,
		float32(uint8(g))/255, float32(uint8(b))/255, float32(uint8(a))/255)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func draw(vao uint32, len int) {
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len/3))
}

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position) {
	// Add the padding so that our calculations fit in a smaller area
	points := square(rect.Size, pos, c.Size())
	square := c.makeVao(points, rect.FillColor)
	draw(square, len(points))
}

func (c *glCanvas) drawObject(o fyne.CanvasObject, offset fyne.Position) {
	pos := o.CurrentPosition().Add(offset)
	switch obj := o.(type) {
	case *canvas.Rectangle:
		c.drawRectangle(obj, pos)
	}
}
