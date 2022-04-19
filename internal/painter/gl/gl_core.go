//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !js && !test_web_driver && !wasm) || (darwin && !mobile && !ios)
// +build !gles,!arm,!arm64,!android,!ios,!mobile,!js,!test_web_driver,!wasm darwin,!mobile,!ios

package gl

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"

	"fyne.io/fyne/v2"
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

const noBuffer = Buffer(0)

// Attribute represents a GL attribute.
type Attribute uint32

// Buffer represents a GL buffer
type Buffer uint32

// Program represents a compiled GL program
type Program uint32

var textureFilterToGL = []int32{gl.LINEAR, gl.NEAREST, gl.LINEAR}

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

	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	info := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(info))

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		return 0, fmt.Errorf("failed to compile OpenGL shader:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE", info, source)
	}

	if logLength > 0 {
		fmt.Printf("OpenGL shader compilation output:\n%s\n>>> SHADER SOURCE\n%s\n<<< SHADER SOURCE\n", info, source)
	}

	return shader, nil
}

func (p *painter) Init() {
	p.ctx = &coreContext{}
	p.program = p.createProgram("simple")
	p.lineProgram = p.createProgram("line")
}

func (p *painter) createProgram(shaderFilename string) Program {
	var vertexSrc []byte
	var fragmentSrc []byte

	// Why a switch over a filename?
	// Because this allows for a minimal change, once we reach Go 1.16 and use go:embed instead of
	// fyne bundle.
	switch shaderFilename {
	case "line":
		vertexSrc = shaderLineVert.StaticContent
		fragmentSrc = shaderLineFrag.StaticContent
	case "simple":
		vertexSrc = shaderSimpleVert.StaticContent
		fragmentSrc = shaderSimpleFrag.StaticContent
	default:
		panic("shader not found: " + shaderFilename)
	}

	vertexShader, err := p.compileShader(string(vertexSrc)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := p.compileShader(string(fragmentSrc)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	var logLength int32
	gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
	info := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(info))

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		panic(fmt.Errorf("failed to link OpenGL program:\n%s", info))
	}

	if logLength > 0 {
		fmt.Printf("OpenGL program linking output:\n%s\n", info)
	}

	if glErr := p.ctx.GetError(); glErr != 0 {
		panic(fmt.Sprintf("failed to link OpenGL program; error code: %x", glErr))
	}

	return Program(prog)
}

func (p *painter) glFreeBuffer(vbo Buffer) {
	p.ctx.BindBuffer(arrayBuffer, noBuffer)
	p.logError()
	p.ctx.DeleteBuffer(vbo)
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
	p.logError()

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	p.logError()
}

func (p *painter) glCapture(width, height int32, pixels *[]uint8) {
	gl.ReadBuffer(gl.FRONT)
	p.logError()
	gl.ReadPixels(0, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(*pixels))
	p.logError()
}

type coreContext struct{}

var _ context = (*coreContext)(nil)

func (c *coreContext) ActiveTexture(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
}

func (c *coreContext) BindBuffer(target uint32, buf Buffer) {
	gl.BindBuffer(target, uint32(buf))
}

func (c *coreContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(target, uint32(texture))
}

func (c *coreContext) BufferData(target uint32, points []float32, usage uint32) {
	gl.BufferData(target, 4*len(points), gl.Ptr(points), usage)
}

func (c *coreContext) Clear(mask uint32) {
	gl.Clear(mask)
}

func (c *coreContext) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c *coreContext) CreateBuffer() Buffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	return Buffer(vbo)
}

func (c *coreContext) CreateTexture() (texture Texture) {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (c *coreContext) DeleteBuffer(buffer Buffer) {
	gl.DeleteBuffers(1, (*uint32)(&buffer))
}

func (c *coreContext) DeleteTexture(texture Texture) {
	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
}

func (c *coreContext) Disable(capability uint32) {
	gl.Disable(capability)
}

func (c *coreContext) Enable(capability uint32) {
	gl.Enable(capability)
}

func (c *coreContext) EnableVertexAttribArray(attribute Attribute) {
	gl.EnableVertexAttribArray(uint32(attribute))
}

func (c *coreContext) GetAttribLocation(program Program, name string) Attribute {
	return Attribute(gl.GetAttribLocation(uint32(program), gl.Str(name+"\x00")))
}

func (c *coreContext) GetError() uint32 {
	return gl.GetError()
}

func (c *coreContext) Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
}

func (c *coreContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
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

func (c *coreContext) TexParameteri(target, param uint32, value int32) {
	gl.TexParameteri(target, param, value)
}

func (c *coreContext) UseProgram(program Program) {
	gl.UseProgram(uint32(program))
}

func (c *coreContext) VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int) {
	gl.VertexAttribPointerWithOffset(uint32(attribute), int32(size), typ, normalized, int32(stride), uintptr(offset))
}

func (c *coreContext) Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}
