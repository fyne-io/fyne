//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package gl

import (
	"encoding/binary"
	"fmt"
	"image/color"

	"github.com/fyne-io/gl-js"
	"golang.org/x/mobile/exp/f32"
)

const (
	arrayBuffer      = gl.ARRAY_BUFFER
	bitColorBuffer   = gl.COLOR_BUFFER_BIT
	bitDepthBuffer   = gl.DEPTH_BUFFER_BIT
	clampToEdge      = gl.CLAMP_TO_EDGE
	colorFormatRGBA  = gl.RGBA
	float            = gl.FLOAT
	scissorTest      = gl.SCISSOR_TEST
	staticDraw       = gl.STATIC_DRAW
	texture0         = gl.TEXTURE0
	texture2D        = gl.TEXTURE_2D
	textureMinFilter = gl.TEXTURE_MIN_FILTER
	textureMagFilter = gl.TEXTURE_MAG_FILTER
	textureWrapS     = gl.TEXTURE_WRAP_S
	textureWrapT     = gl.TEXTURE_WRAP_T
	unsignedByte     = gl.UNSIGNED_BYTE
)

// Attribute represents a GL attribute.
type Attribute gl.Attrib

// Buffer represents a GL buffer
type Buffer gl.Buffer

// Program represents a compiled GL program
type Program gl.Program

var noBuffer = Buffer(gl.NoBuffer)
var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST}

func (p *painter) glInit() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	p.logError()
}

func (p *painter) compileShader(source string, shaderType gl.Enum) (gl.Shader, error) {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	p.logError()
	gl.CompileShader(shader)
	p.logError()

	status := gl.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		info := gl.GetShaderInfoLog(shader)

		return gl.NoShader, fmt.Errorf("failed to compile %v: %v", source, info)
	}

	return shader, nil
}

var vertexShaderSource = string(shaderSimpleesVert.StaticContent)
var fragmentShaderSource = string(shaderSimpleesFrag.StaticContent)
var vertexLineShaderSource = string(shaderLineesVert.StaticContent)
var fragmentLineShaderSource = string(shaderLineesFrag.StaticContent)

func (p *painter) Init() {
	p.ctx = &xjsContext{}
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

	colorUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "color")
	r, g, b, a := col.RGBA()
	if a == 0 {
		gl.Uniform4f(colorUniform, 0, 0, 0, 0)
	} else {
		alpha := float32(a)
		col := []float32{float32(r) / alpha, float32(g) / alpha, float32(b) / alpha, alpha / 0xffff}
		gl.Uniform4fv(colorUniform, col)
	}
	lineWidthUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "lineWidth")
	gl.Uniform1f(lineWidthUniform, width)

	featherUniform := gl.GetUniformLocation(gl.Program(p.lineProgram), "feather")
	gl.Uniform1f(featherUniform, feather)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	p.logError()
}

func (p *painter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadPixels(*pixels, 0, 0, int(width), int(height), gl.RGBA, gl.UNSIGNED_BYTE)
	p.logError()
}

type xjsContext struct{}

var _ context = (*xjsContext)(nil)

func (c *xjsContext) ActiveTexture(textureUnit uint32) {
	gl.ActiveTexture(gl.Enum(textureUnit))
}

func (c *xjsContext) BindBuffer(target uint32, buf Buffer) {
	gl.BindBuffer(gl.Enum(target), gl.Buffer(buf))
}

func (c *xjsContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(gl.Enum(target), gl.Texture(texture))
}

func (c *xjsContext) BlendColor(r, g, b, a float32) {
	gl.BlendColor(r, g, b, a)
}

func (c *xjsContext) BufferData(target uint32, points []float32, usage uint32) {
	gl.BufferData(gl.Enum(target), f32.Bytes(binary.LittleEndian, points...), gl.Enum(usage))
}

func (c *xjsContext) Clear(mask uint32) {
	gl.Clear(gl.Enum(mask))
}

func (c *xjsContext) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c *xjsContext) CreateBuffer() Buffer {
	return Buffer(gl.CreateBuffer())
}

func (c *xjsContext) CreateTexture() (texture Texture) {
	return Texture(gl.CreateTexture())
}

func (c *xjsContext) DeleteBuffer(buffer Buffer) {
	gl.DeleteBuffer(gl.Buffer(buffer))
}

func (c *xjsContext) DeleteTexture(texture Texture) {
	gl.DeleteTexture(gl.Texture(texture))
}

func (c *xjsContext) Disable(capability uint32) {
	gl.Disable(gl.Enum(capability))
}

func (c *xjsContext) Enable(capability uint32) {
	gl.Enable(gl.Enum(capability))
}

func (c *xjsContext) EnableVertexAttribArray(attribute Attribute) {
	gl.EnableVertexAttribArray(gl.Attrib(attribute))
}

func (c *xjsContext) GetAttribLocation(program Program, name string) Attribute {
	return Attribute(gl.GetAttribLocation(gl.Program(program), name))
}

func (c *xjsContext) GetError() uint32 {
	return uint32(gl.GetError())
}

func (c *xjsContext) Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
}

func (c *xjsContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	gl.TexImage2D(
		gl.Enum(target),
		level,
		width,
		height,
		gl.Enum(colorFormat),
		gl.Enum(typ),
		data,
	)
}

func (c *xjsContext) TexParameteri(target, param uint32, value int32) {
	gl.TexParameteri(gl.Enum(target), gl.Enum(param), int(value))
}

func (c *xjsContext) UseProgram(program Program) {
	gl.UseProgram(gl.Program(program))
}

func (c *xjsContext) VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(gl.Attrib(attribute), size, gl.Enum(typ), normalized, stride, offset)
}

func (c *xjsContext) Viewport(x, y, width, height int) {
	gl.Viewport(x, y, width, height)
}
