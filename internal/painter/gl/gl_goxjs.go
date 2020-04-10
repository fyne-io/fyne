// +build js wasm web

package gl

import (
	"fmt"
	"image"
	"image/draw"

	"encoding/binary"
	"golang.org/x/mobile/exp/f32"
	"github.com/goxjs/gl"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// Buffer represents a GL buffer
type Buffer gl.Buffer

// Program represents a compiled GL program
type Program gl.Program

// Texture represents an uploaded GL texture
type Texture gl.Texture

// NoTexture is the zero value for a Texture
var NoTexture = Texture(gl.NoTexture)

func (p *glPainter) newTexture() Texture {
	var texture = gl.CreateTexture()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	return Texture(texture)
}

func getTexture(object fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture) (Texture, error) {
	texture, ok := textures[object]

	if !ok {
		texture = creator(object)
		textures[object] = texture
	}
	if !gl.Texture(texture).Valid() {
		return NoTexture, fmt.Errorf("No texture available.")
	}
	return texture, nil
}

func (p *glPainter) imgToTexture(img image.Image) Texture {
	switch i := img.(type) {
	case *image.Uniform:
		texture := p.newTexture()
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, data)
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return NoTexture
		}
		texture := p.newTexture()
		gl.TexImage2D(gl.TEXTURE_2D, 0, i.Rect.Size().X, i.Rect.Size().Y,
			gl.RGBA, gl.UNSIGNED_BYTE, i.Pix)
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.ZP, draw.Over)
		return p.imgToTexture(rgba)
	}
}

func (p *glPainter) SetOutputSize(width, height int) {
	gl.Viewport(0, 0, width, height)
}

func (p *glPainter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := textures[obj]
	if ok {
		gl.DeleteTexture(gl.Texture(texture))
		delete(textures, obj)
	}
}

func (p *glPainter) compileShader(source string, shaderType gl.Enum) (gl.Shader, error) {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)

	status := gl.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		info := gl.GetShaderInfoLog(shader)

		return shader, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

const (
	vertexShaderSource = `
    #version 100
    attribute vec3 vert;
    attribute vec2 vertTexCoord;
    varying highp vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }`

	fragmentShaderSource = `
    #version 100
    uniform sampler2D tex;

    varying highp vec2 fragTexCoord;

    void main() {
        gl_FragColor = texture2D(tex, fragTexCoord);
    }`
)

func (p *glPainter) Init() {
	gl.Disable(gl.DEPTH_TEST)

	vertexShader, err := p.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := p.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	p.program = Program(prog)
}

func (p *glPainter) glClearBuffer() {
	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (p *glPainter) glScissorOpen(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
	gl.Enable(gl.SCISSOR_TEST)
}

func (p *glPainter) glScissorClose() {
	gl.Disable(gl.SCISSOR_TEST)
}

func (p *glPainter) glCreateBuffer(points []float32) Buffer {
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, points...), gl.STATIC_DRAW)

	vertAttrib := gl.GetAttribLocation(gl.Program(p.program), "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	texCoordAttrib := gl.GetAttribLocation(gl.Program(p.program), "vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	return Buffer(vbo)
}

func (p *glPainter) glFreeBuffer(vbo Buffer) {
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer(vbo))
	gl.DeleteBuffer(gl.Buffer(vbo))
}

func (p *glPainter) glDrawTexture(texture Texture, alpha float32) {
	gl.UseProgram(gl.Program(p.program))
	gl.Enable(gl.BLEND)

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		gl.BlendColor(0, 0, 0, alpha)
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture(texture))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (p *glPainter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadPixels(*pixels, 0, 0, int(width), int(height), gl.RGBA, gl.UNSIGNED_BYTE)
}

func glInit() {
	// no-op, gomobile does this
}
