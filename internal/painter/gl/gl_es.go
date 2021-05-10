// +build gles arm arm64
// +build !android,!ios,!mobile
// +build !darwin

package gl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
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
		draw.Draw(rgba, rgba.Rect, img, image.ZP, draw.Over)
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
    #version 100
    attribute vec3 vert;
    attribute vec2 vertTexCoord;
    varying highp vec2 fragTexCoord;

    void main() {
        fragTexCoord = vertTexCoord;

        gl_Position = vec4(vert, 1);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 100
    uniform sampler2D tex;

    varying highp vec2 fragTexCoord;

    void main() {
        gl_FragColor = texture2D(tex, fragTexCoord);
    }
` + "\x00"

	vertexLineShaderSource = `
    #version 100
    attribute vec2 vert;
    attribute vec2 normal;
    
    uniform float lineWidth;

    varying highp vec2 delta;

    void main() {
        delta = normal * lineWidth;

        gl_Position = vec4(vert + delta, 0, 1);
    }
` + "\x00"

	fragmentLineShaderSource = `
    #version 100
    uniform highp vec4 color;
    uniform highp float lineWidth;
    uniform highp float feather;

    varying highp vec2 delta;

    void main() {
        highp float alpha = color.a;
        highp float distance = length(delta);

        if (feather == 0.0 || distance <= lineWidth - feather) {
           gl_FragColor = color;
        } else {
           gl_FragColor = vec4(color.r, color.g, color.b, mix(color.a, 0.0, (distance - (lineWidth - feather)) / feather));
        }
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

	vertexLineShader, err := compileShader(vertexLineShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentLineShader, err := compileShader(fragmentLineShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	lineProg := gl.CreateProgram()
	gl.AttachShader(lineProg, vertexLineShader)
	gl.AttachShader(lineProg, fragmentLineShader)
	gl.LinkProgram(lineProg)
	logError()

	p.lineProgram = Program(lineProg)
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
	gl.UseProgram(uint32(p.program))

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
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	logError()

	return Buffer(vbo)
}

func (p *glPainter) glCreateLineBuffer(points []float32) Buffer {
	gl.UseProgram(uint32(p.lineProgram))

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	logError()
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	logError()

	vertAttrib := uint32(gl.GetAttribLocation(uint32(p.lineProgram), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	logError()

	normalAttrib := uint32(gl.GetAttribLocation(uint32(p.lineProgram), gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointer(normalAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
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
	gl.UseProgram(uint32(p.program))

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

func (p *glPainter) glDrawLine(width float32, col color.Color, feather float32) {
	gl.UseProgram(uint32(p.lineProgram))

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	logError()

	colorUniform := gl.GetUniformLocation(uint32(p.lineProgram), gl.Str("color\x00"))
	r, g, b, a := col.RGBA()
	if a == 0 {
		gl.Uniform4f(colorUniform, 0, 0, 0, 0)
	} else {
		alpha := float32(a)
		col := []float32{float32(r) / alpha, float32(g) / alpha, float32(b) / alpha, alpha / 0xffff}
		gl.Uniform4fv(colorUniform, 1, &col[0])
	}
	lineWidthUniform := gl.GetUniformLocation(uint32(p.lineProgram), gl.Str("lineWidth\x00"))
	gl.Uniform1f(lineWidthUniform, width)

	featherUniform := gl.GetUniformLocation(uint32(p.lineProgram), gl.Str("feather\x00"))
	gl.Uniform1f(featherUniform, feather)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	logError()
}

func (p *glPainter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadBuffer(gl.FRONT)
	logError()
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*pixels))
	logError()
}

func logError() {
	if fyne.CurrentApp().Settings().BuildType() != fyne.BuildDebug {
		return
	}
	logGLError(gl.GetError())
}
