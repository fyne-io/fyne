// +build !ci,gl

package gl

import (
	"fmt"
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image/draw"
	"log"
	"os"
	"path/filepath"

	"image"
	"image/color"
	_ "image/png"
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
    in vec2 vertTexCoord;
    out vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 130
    in vec2 fragTexCoord;
    out vec4 frag_colour;
    uniform sampler2D tex;
    
    void main() {
        vec4 color = texture(tex, fragTexCoord);
        if(color.a < 0.01)
            discard;

        frag_colour = color;
    }
` + "\x00"
)

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

	return []float32{
		// coord x, y, x texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
	}
}

// textureForPoints initializes a vertex array and prepares a texture to draw on it
func (c *glCanvas) textureForPoints(points []float32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)

	textureUniform := gl.GetUniformLocation(c.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(c.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindVertexArray(vao)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
}

func (c *glCanvas) drawRectangle(rect *canvas.Rectangle, pos fyne.Position) {
	points := c.rectCoords(rect.Size, pos)
	c.textureForPoints(points)

	r, g, b, a := rect.FillColor.RGBA()
	data := []uint8{uint8(r), uint8(g), uint8(b), uint8(a)}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA,
		gl.UNSIGNED_BYTE, gl.Ptr(data))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
}

func (c *glCanvas) drawRawImage(img *image.RGBA, size fyne.Size, pos fyne.Position) {
	points := c.rectCoords(size, pos)
	c.textureForPoints(points)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X), int32(img.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(points)/5))
}

func (c *glCanvas) drawImage(img *canvas.Image, pos fyne.Position) {
	raw := image.NewRGBA(image.Rect(0, 0, img.Size.Width, img.Size.Height))

	if img.File != "" {
		if strings.ToLower(filepath.Ext(img.File)) == ".svg" {
			placeholderColor := &image.Uniform{color.RGBA{0, 0, 255, 255}}
			draw.Draw(raw, raw.Bounds(), placeholderColor, image.ZP, draw.Src)
		} else {
			file, _ := os.Open(img.File)
			pixels, _, err := image.Decode(file)

			if err != nil {
				log.Println("image err", err)

				errColor := &image.Uniform{color.RGBA{255, 0, 0, 255}}
				draw.Draw(raw, raw.Bounds(), errColor, image.ZP, draw.Src)
			} else {
				raw = image.NewRGBA(pixels.Bounds())

				draw.Draw(raw, raw.Bounds(), pixels, image.ZP, draw.Src)
			}
		}
	} else if img.PixelColor != nil {
		pixels := NewPixelImage(img)
		draw.Draw(raw, raw.Bounds(), pixels, image.ZP, draw.Src)
	}

	c.drawRawImage(raw, img.Size, pos)
}

func (c *glCanvas) drawText(text *canvas.Text, pos fyne.Position) {
	bounds := text.MinSize()
	width := scaleInt(c, bounds.Width)
	height := scaleInt(c, bounds.Height)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var opts truetype.Options
	font := fontCache()
	fontSize := float64(text.TextSize) * float64(c.Scale())
	opts.Size = fontSize
	face := truetype.NewFace(fontCache(), &opts)

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(font)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(&image.Uniform{theme.TextColor()})

	ctx.DrawString(text.Text, freetype.Pt(0, height+2-face.Metrics().Descent.Ceil()))
	c.drawRawImage(img, bounds, pos)
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
