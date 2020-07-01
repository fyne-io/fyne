// +build !gles,!arm,!arm64,!android,!ios,!mobile

package gl

import (
	"fmt"
	"image"
	"image/draw"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// Buffer represents a GL buffer
type Buffer uint32

// Program represents a compiled GL program
type Program uint32

// Texture represents an uploaded GL texture
type Texture uint32

// NoTexture is the zero value for a Texture
var NoTexture = Texture(0)

var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST}

func newTexture(textureFilter canvas.ImageScale) Texture {
	var texture uint32

	if int(textureFilter) >= len(textureFilterToGL) {
		fyne.LogError(fmt.Sprintf("Invalid canvas.ImageScale value (%d), using canvas.ImageScaleSmooth as default value", textureFilter), nil)
		textureFilter = canvas.ImageScaleSmooth
	}

	gl.GenTextures(1, &texture)
	logError()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	logError()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, textureFilterToGL[textureFilter])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, textureFilterToGL[textureFilter])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	logError()

	return Texture(texture)
}

func (p *glPainter) imgToTexture(img image.Image, textureFilter canvas.ImageScale) Texture {
	switch i := img.(type) {
	case *image.Uniform:
		texture := newTexture(textureFilter)
		r, g, b, a := i.RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		data := []uint8{r8, g8, b8, a8}
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA,
			gl.UNSIGNED_BYTE, gl.Ptr(data))
		logError()
		return texture
	case *image.RGBA:
		if len(i.Pix) == 0 { // image is empty
			return 0
		}

		texture := newTexture(textureFilter)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(i.Rect.Size().X), int32(i.Rect.Size().Y),
			0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(i.Pix))
		logError()
		return texture
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Over)
		return p.imgToTexture(rgba, textureFilter)
	}
}

func (p *glPainter) SetOutputSize(width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	logError()
}

func (p *glPainter) freeTexture(obj fyne.CanvasObject) {
	texture, ok := textures[obj]
	if !ok {
		return
	}

	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
	logError()
	delete(textures, obj)
}

func glInit() {
	err := gl.Init()
	if err != nil {
		fyne.LogError("failed to initialise OpenGL", err)
		return
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	logError()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	logError()
	free()
	gl.CompileShader(shader)
	logError()

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		info := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(info))

		return 0, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

const (
	vertexShaderSource = `
    #version 110
    attribute vec3 vert;
    attribute vec2 vertTexCoord;
    varying vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 110
    uniform sampler2D tex;

    varying vec2 fragTexCoord;

    void main() {
        gl_FragColor = texture2D(tex, fragTexCoord);
    }
` + "\x00"
)

func (p *glPainter) Init() {
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
	logError()

	p.program = Program(prog)
}

func (p *glPainter) glClearBuffer() {
	gl.UseProgram(uint32(p.program))
	logError()

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	logError()
}

func (p *glPainter) glScissorOpen(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
	gl.Enable(gl.SCISSOR_TEST)
	logError()
}

func (p *glPainter) glScissorClose() {
	gl.Disable(gl.SCISSOR_TEST)
	logError()
}

func (p *glPainter) glCreateBuffer(points []float32) Buffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	logError()
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	logError()

	vertAttrib := uint32(gl.GetAttribLocation(uint32(p.program), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	logError()

	texCoordAttrib := uint32(gl.GetAttribLocation(uint32(p.program), gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(12))
	logError()

	return Buffer(vbo)
}

func (p *glPainter) glFreeBuffer(vbo Buffer) {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	logError()
	buf := uint32(vbo)
	gl.DeleteBuffers(1, &buf)
	logError()
}

func (p *glPainter) glDrawTexture(texture Texture, alpha float32) {
	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		gl.BlendColor(0, 0, 0, alpha)
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	logError()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(texture))
	logError()

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	logError()
}

func (p *glPainter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadBuffer(gl.FRONT)
	logError()
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*pixels))
	logError()
}

func logError() {
	logGLError(gl.GetError())
}
