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
)

const (
	arrayBuffer           = gl.ARRAY_BUFFER
	bitColorBuffer        = gl.COLOR_BUFFER_BIT
	bitDepthBuffer        = gl.DEPTH_BUFFER_BIT
	clampToEdge           = gl.CLAMP_TO_EDGE
	colorFormatRGBA       = gl.RGBA
	constantAlpha         = gl.CONSTANT_ALPHA
	float                 = gl.FLOAT
	one                   = gl.ONE
	oneMinusConstantAlpha = gl.ONE_MINUS_CONSTANT_ALPHA
	oneMinusSrcAlpha      = gl.ONE_MINUS_SRC_ALPHA
	scissorTest           = gl.SCISSOR_TEST
	srcAlpha              = gl.SRC_ALPHA
	staticDraw            = gl.STATIC_DRAW
	texture0              = gl.TEXTURE0
	texture2D             = gl.TEXTURE_2D
	textureMinFilter      = gl.TEXTURE_MIN_FILTER
	textureMagFilter      = gl.TEXTURE_MAG_FILTER
	textureWrapS          = gl.TEXTURE_WRAP_S
	textureWrapT          = gl.TEXTURE_WRAP_T
	triangles             = gl.TRIANGLES
	triangleStrip         = gl.TRIANGLE_STRIP
	unsignedByte          = gl.UNSIGNED_BYTE
)

const noBuffer = Buffer(0)

// Attribute represents a GL attribute.
type Attribute uint32

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

func (p *painter) glDrawTexture(texture Texture, alpha float32) {
	p.ctx.UseProgram(p.program)

	// here we have to choose between blending the image alpha or fading it...
	// TODO find a way to support both
	if alpha != 1.0 {
		p.ctx.BlendColor(0, 0, 0, alpha)
		p.ctx.BlendFunc(constantAlpha, oneMinusConstantAlpha)
	} else {
		p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	}
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) glDrawLine(width float32, col color.Color, feather float32) {
	p.ctx.UseProgram(p.lineProgram)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
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

	p.ctx.DrawArrays(triangles, 0, 6)
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

func (c *esContext) BindBuffer(target uint32, buf Buffer) {
	gl.BindBuffer(target, uint32(buf))
}

func (c *esContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(target, uint32(texture))
}

func (c *esContext) BlendColor(r, g, b, a float32) {
	gl.BlendColor(r, g, b, a)
}

func (c *esContext) BlendFunc(srcFactor, destFactor uint32) {
	gl.BlendFunc(srcFactor, destFactor)
}

func (c *esContext) BufferData(target uint32, points []float32, usage uint32) {
	gl.BufferData(target, 4*len(points), gl.Ptr(points), usage)
}

func (c *esContext) Clear(mask uint32) {
	gl.Clear(mask)
}

func (c *esContext) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c *esContext) CreateBuffer() Buffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	return Buffer(vbo)
}

func (c *esContext) CreateTexture() (texture Texture) {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (c *esContext) DeleteBuffer(buffer Buffer) {
	gl.DeleteBuffers(1, (*uint32)(&buffer))
}

func (c *esContext) DeleteTexture(texture Texture) {
	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
}

func (c *esContext) Disable(capability uint32) {
	gl.Disable(capability)
}

func (c *esContext) DrawArrays(mode uint32, first, count int) {
	gl.DrawArrays(mode, int32(first), int32(count))
}

func (c *esContext) Enable(capability uint32) {
	gl.Enable(capability)
}

func (c *esContext) EnableVertexAttribArray(attribute Attribute) {
	gl.EnableVertexAttribArray(uint32(attribute))
}

func (c *esContext) GetAttribLocation(program Program, name string) Attribute {
	return Attribute(gl.GetAttribLocation(uint32(program), gl.Str(name+"\x00")))
}

func (c *esContext) GetError() uint32 {
	return gl.GetError()
}

func (c *esContext) Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
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

func (c *esContext) VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int) {
	gl.VertexAttribPointerWithOffset(uint32(attribute), int32(size), typ, normalized, int32(stride), uintptr(offset))
}

func (c *esContext) Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}
