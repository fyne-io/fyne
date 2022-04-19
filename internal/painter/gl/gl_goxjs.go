//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package gl

import (
	"encoding/binary"
	"fmt"

	"github.com/fyne-io/gl-js"
	"golang.org/x/mobile/exp/f32"
)

const (
	arrayBuffer           = gl.ARRAY_BUFFER
	bitColorBuffer        = gl.COLOR_BUFFER_BIT
	bitDepthBuffer        = gl.DEPTH_BUFFER_BIT
	clampToEdge           = gl.CLAMP_TO_EDGE
	colorFormatRGBA       = gl.RGBA
	constantAlpha         = gl.CONSTANT_ALPHA
	float                 = gl.FLOAT
	front                 = gl.FRONT
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

type (
	// Attribute represents a GL attribute
	Attribute gl.Attrib
	// Buffer represents a GL buffer
	Buffer gl.Buffer
	// Program represents a compiled GL program
	Program gl.Program
	// Shader represents a GL shader
	Shader gl.Shader
	// Uniform represents a GL uniform
	Uniform gl.Uniform
)

var noBuffer = Buffer(gl.NoBuffer)
var noShader = Shader(gl.NoShader)
var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST}

func (p *painter) glInit() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	p.logError()
}

func (p *painter) compileShader(source string, shaderType gl.Enum) (Shader, error) {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	p.logError()
	gl.CompileShader(shader)
	p.logError()

	info := gl.GetShaderInfoLog(shader)
	status := gl.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		return noShader, fmt.Errorf("failed to compile OpenGL shader:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE", info, source)
	}

	if len(info) > 0 {
		fmt.Printf("OpenGL shader compilation output:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE\n", info, source)
	}

	return Shader(shader), nil
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
	gl.AttachShader(prog, gl.Shader(vertexShader))
	gl.AttachShader(prog, gl.Shader(fragmentShader))
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
	gl.AttachShader(lineProg, gl.Shader(vertexLineShader))
	gl.AttachShader(lineProg, gl.Shader(fragmentLineShader))
	gl.LinkProgram(lineProg)
	p.logError()

	p.lineProgram = Program(lineProg)
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

func (c *xjsContext) BlendFunc(srcFactor, destFactor uint32) {
	gl.BlendFunc(gl.Enum(srcFactor), gl.Enum(destFactor))
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

func (c *xjsContext) DrawArrays(mode uint32, first, count int) {
	gl.DrawArrays(gl.Enum(mode), first, count)
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

func (c *xjsContext) GetUniformLocation(program Program, name string) Uniform {
	return Uniform(gl.GetUniformLocation(gl.Program(program), name))
}

func (c *xjsContext) ReadBuffer(_ uint32) {
}

func (c *xjsContext) ReadPixels(x, y, width, height int, colorFormat, typ uint32, pixels []uint8) {
	gl.ReadPixels(pixels, x, y, width, height, gl.Enum(colorFormat), gl.Enum(typ))
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

func (c *xjsContext) Uniform1f(uniform Uniform, v float32) {
	gl.Uniform1f(gl.Uniform(uniform), v)
}

func (c *xjsContext) Uniform4f(uniform Uniform, v0, v1, v2, v3 float32) {
	gl.Uniform4f(gl.Uniform(uniform), v0, v1, v2, v3)
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
