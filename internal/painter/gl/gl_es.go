//go:build (gles || arm || arm64) && !android && !ios && !mobile && !darwin && !js && !wasm && !test_web_driver
// +build gles arm arm64
// +build !android
// +build !ios
// +build !mobile
// +build !darwin
// +build !js
// +build !wasm
// +build !test_web_driver

package gl

import (
	"fmt"
	"image/color"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	clampToEdge      = gl.CLAMP_TO_EDGE
	colorFormatRGBA  = gl.RGBA
	texture0         = gl.TEXTURE0
	texture2D        = gl.TEXTURE_2D
	textureMinFilter = gl.TEXTURE_MIN_FILTER
	textureMagFilter = gl.TEXTURE_MAG_FILTER
	textureWrapS     = gl.TEXTURE_WRAP_S
	textureWrapT     = gl.TEXTURE_WRAP_T
	unsignedByte     = gl.UNSIGNED_BYTE
)

// Buffer represents a GL buffer
type Buffer uint32

// Program represents a compiled GL program
type Program uint32

var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST}

func (p *painter) glInit() {
	err := gl.Init()
	if err != nil {
		fyne.LogError("failed to initialise OpenGL", err)
		return
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	p.logError()
}

func (p *painter) compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	p.logError()
	free()
	gl.CompileShader(shader)
	p.logError()

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

var vertexShaderSource = string(shaderSimpleesVert.StaticContent) + "\x00"
var fragmentShaderSource = string(shaderSimpleesFrag.StaticContent) + "\x00"
var vertexLineShaderSource = string(shaderLineesVert.StaticContent) + "\x00"
var fragmentLineShaderSource = string(shaderLineesFrag.StaticContent) + "\x00"

func (p *painter) Init() {
	p.ctx = &esContext{}
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
	p.logError()

	p.program = Program(prog)

	vertexLineShader, err := p.compileShader(vertexLineShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentLineShader, err := p.compileShader(fragmentLineShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	lineProg := gl.CreateProgram()
	gl.AttachShader(lineProg, vertexLineShader)
	gl.AttachShader(lineProg, fragmentLineShader)
	gl.LinkProgram(lineProg)
	p.logError()

	p.lineProgram = Program(lineProg)
}

func (p *painter) glClearBuffer() {
	p.ctx.UseProgram(p.program)
	p.logError()

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	gl.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	p.logError()
}

func (p *painter) glScissorOpen(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
	gl.Enable(gl.SCISSOR_TEST)
	p.logError()
}

func (p *painter) glScissorClose() {
	gl.Disable(gl.SCISSOR_TEST)
	p.logError()
}

func (p *painter) glCreateBuffer(points []float32) Buffer {
	p.ctx.UseProgram(p.program)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	p.logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	p.logError()
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	p.logError()

	vertAttrib := uint32(gl.GetAttribLocation(uint32(p.program), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)
	p.logError()

	texCoordAttrib := uint32(gl.GetAttribLocation(uint32(p.program), gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)
	p.logError()

	return Buffer(vbo)
}

func (p *painter) glCreateLineBuffer(points []float32) Buffer {
	p.ctx.UseProgram(p.lineProgram)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	p.logError()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	p.logError()
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	p.logError()

	vertAttrib := uint32(gl.GetAttribLocation(uint32(p.lineProgram), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 2, gl.FLOAT, false, 4*4, 0)
	p.logError()

	normalAttrib := uint32(gl.GetAttribLocation(uint32(p.lineProgram), gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointerWithOffset(normalAttrib, 2, gl.FLOAT, false, 4*4, 2*4)
	p.logError()

	return Buffer(vbo)
}

func (p *painter) glFreeBuffer(vbo Buffer) {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	p.logError()
	buf := uint32(vbo)
	gl.DeleteBuffers(1, &buf)
	p.logError()
}

func (p *painter) glDrawTexture(texture Texture, alpha float32) {
	p.ctx.UseProgram(p.program)

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		gl.BlendColor(0, 0, 0, alpha)
		gl.BlendFunc(gl.CONSTANT_ALPHA, gl.ONE_MINUS_CONSTANT_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	p.logError()
}

func (p *painter) glDrawLine(width float32, col color.Color, feather float32) {
	p.ctx.UseProgram(p.lineProgram)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	p.logError()

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
	p.logError()
}

func (p *painter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadBuffer(gl.FRONT)
	p.logError()
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*pixels))
	p.logError()
}

type esContext struct{}

var _ context = (*esContext)(nil)

func (c *esContext) ActiveTexture(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
}

func (c *esContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(target, uint32(texture))
}

func (c *esContext) CreateTexture() (texture Texture) {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (c *esContext) DeleteTexture(texture Texture) {
	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
}

func (c *esContext) GetError() uint32 {
	return gl.GetError()
}

func (c *esContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	gl.TexImage2D(
		target,
		int32(level),
		int32(colorFormat),
		int32(width),
		int32(height),
		0,
		colorFormat,
		typ,
		gl.Ptr(data),
	)
}

func (c *esContext) TexParameteri(target, param uint32, value int32) {
	gl.TexParameteri(target, param, value)
}

func (c *esContext) UseProgram(program Program) {
	gl.UseProgram(uint32(program))
}

func (c *esContext) Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}
