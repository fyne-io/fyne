//go:build (!gles && !arm && !arm64 && !android && !ios && !mobile && !test_web_driver && !wasm) || (darwin && !mobile && !ios && !wasm && !test_web_driver)

package gl

import (
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"

	"fyne.io/fyne/v2"
)

const (
	arrayBuffer           = gl.ARRAY_BUFFER
	bitColorBuffer        = gl.COLOR_BUFFER_BIT
	bitDepthBuffer        = gl.DEPTH_BUFFER_BIT
	clampToEdge           = gl.CLAMP_TO_EDGE
	colorFormatRGBA       = gl.RGBA
	compileStatus         = gl.COMPILE_STATUS
	constantAlpha         = gl.CONSTANT_ALPHA
	float                 = gl.FLOAT
	fragmentShader        = gl.FRAGMENT_SHADER
	front                 = gl.FRONT
	back                  = gl.BACK
	glFalse               = gl.FALSE
	linkStatus            = gl.LINK_STATUS
	maxTextureSizeParam   = gl.MAX_TEXTURE_SIZE
	one                   = gl.ONE
	zero                  = gl.ZERO
	oneMinusConstantAlpha = gl.ONE_MINUS_CONSTANT_ALPHA
	oneMinusSrcAlpha      = gl.ONE_MINUS_SRC_ALPHA
	scissorTest           = gl.SCISSOR_TEST
	srcAlpha              = gl.SRC_ALPHA
	staticDraw            = gl.STATIC_DRAW
	texture0              = gl.TEXTURE0
	texture1              = gl.TEXTURE1
	texture2D             = gl.TEXTURE_2D
	textureNearest        = gl.NEAREST
	textureMinFilter      = gl.TEXTURE_MIN_FILTER
	textureMagFilter      = gl.TEXTURE_MAG_FILTER
	textureWrapS          = gl.TEXTURE_WRAP_S
	textureWrapT          = gl.TEXTURE_WRAP_T
	triangles             = gl.TRIANGLES
	triangleStrip         = gl.TRIANGLE_STRIP
	unsignedByte          = gl.UNSIGNED_BYTE
	vertexShader          = gl.VERTEX_SHADER
)

const (
	noBuffer  = Buffer(0)
	noProgram = Program(0)
	noShader  = Shader(0)
)

type (
	// Attribute represents a GL attribute
	Attribute uint32
	// Buffer represents a GL buffer
	Buffer uint32
	// Program represents a compiled GL program
	Program uint32
	// Shader represents a GL shader
	Shader uint32
	// Uniform represents a GL uniform
	Uniform int32
)

var textureFilterToGL = [...]int32{gl.LINEAR, gl.NEAREST, gl.LINEAR}

func (p *painter) Init() {
	p.ctx = &coreContext{}
	err := gl.Init()
	if err != nil {
		fyne.LogError("failed to initialise OpenGL", err)
		return
	}
	p.maxTextureSize = p.ctx.GetInteger(maxTextureSizeParam)

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	p.logError()
	p.programs = p.compilePrograms()
}

type coreContext struct{}

var _ context = (*coreContext)(nil)

func (c *coreContext) ActiveTexture(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
}

func (c *coreContext) AttachShader(program Program, shader Shader) {
	gl.AttachShader(uint32(program), uint32(shader))
}

func (c *coreContext) BindBuffer(target uint32, buf Buffer) {
	gl.BindBuffer(target, uint32(buf))
}

func (c *coreContext) BindTexture(target uint32, texture Texture) {
	gl.BindTexture(target, uint32(texture))
}

func (c *coreContext) BlendColor(r, g, b, a float32) {
	gl.BlendColor(r, g, b, a)
}

func (c *coreContext) BlendFunc(srcFactor, destFactor uint32) {
	gl.BlendFunc(srcFactor, destFactor)
}

func (c *coreContext) BufferData(target uint32, points []float32, usage uint32) {
	gl.BufferData(target, 4*len(points), gl.Ptr(points), usage)
}

func (c *coreContext) BufferSubData(target uint32, points []float32) {
	gl.BufferSubData(target, 0, 4*len(points), gl.Ptr(points))
}

func (c *coreContext) Clear(mask uint32) {
	gl.Clear(mask)
}

func (c *coreContext) ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c *coreContext) CompileShader(shader Shader) {
	gl.CompileShader(uint32(shader))
}

func (c *coreContext) CreateBuffer() Buffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	return Buffer(vbo)
}

func (c *coreContext) CreateProgram() Program {
	return Program(gl.CreateProgram())
}

func (c *coreContext) CreateShader(typ uint32) Shader {
	return Shader(gl.CreateShader(typ))
}

func (c *coreContext) CreateTexture() (texture Texture) {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (c *coreContext) DeleteBuffer(buffer Buffer) {
	gl.DeleteBuffers(1, (*uint32)(&buffer))
}

func (c *coreContext) DeleteProgram(program Program) {
	gl.DeleteProgram(uint32(program))
}

func (c *coreContext) DeleteTexture(texture Texture) {
	tex := uint32(texture)
	gl.DeleteTextures(1, &tex)
}

func (c *coreContext) Disable(capability uint32) {
	gl.Disable(capability)
}

func (c *coreContext) DrawArrays(mode uint32, first, count int) {
	gl.DrawArrays(
		mode,
		int32(first), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(count), //gosec:disable G115 -- we are definitely fine with limiting the value range here
	)
}

func (c *coreContext) Enable(capability uint32) {
	gl.Enable(capability)
}

func (c *coreContext) EnableVertexAttribArray(attribute Attribute) {
	gl.EnableVertexAttribArray(uint32(attribute))
}

func (c *coreContext) GetAttribLocation(program Program, name string) Attribute {
	return Attribute(gl.GetAttribLocation(uint32(program), gl.Str(name+"\x00"))) //gosec:disable G115 -- the attribute location is a pointer, so unsigned is fine
}

func (c *coreContext) GetError() uint32 {
	return gl.GetError()
}

func (c *coreContext) GetInteger(pname uint32) int {
	var value int32
	gl.GetIntegerv(pname, &value)
	return int(value)
}

func (c *coreContext) GetProgrami(program Program, param uint32) int {
	var value int32
	gl.GetProgramiv(uint32(program), param, &value)
	return int(value)
}

func (c *coreContext) GetProgramInfoLog(program Program) string {
	var logLength int32
	gl.GetProgramiv(uint32(program), gl.INFO_LOG_LENGTH, &logLength)
	info := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(uint32(program), logLength, nil, gl.Str(info))
	return info
}

func (c *coreContext) GetShaderi(shader Shader, param uint32) int {
	var value int32
	gl.GetShaderiv(uint32(shader), param, &value)
	return int(value)
}

func (c *coreContext) GetShaderInfoLog(shader Shader) string {
	var logLength int32
	gl.GetShaderiv(uint32(shader), gl.INFO_LOG_LENGTH, &logLength)
	info := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(uint32(shader), logLength, nil, gl.Str(info))
	return info
}

func (c *coreContext) GetUniformLocation(program Program, name string) Uniform {
	return Uniform(gl.GetUniformLocation(uint32(program), gl.Str(name+"\x00")))
}

func (c *coreContext) LinkProgram(program Program) {
	gl.LinkProgram(uint32(program))
}

func (c *coreContext) CopyTexSubImage2D(target uint32, level, xoffset, yoffset, x, y, width, height int) {
	gl.CopyTexSubImage2D(
		target,
		int32(level),   //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(xoffset), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(yoffset), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(x),       //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(y),       //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(width),   //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(height),  //gosec:disable G115 -- we are definitely fine with limiting the value range here
	)
}

func (c *coreContext) ReadBuffer(src uint32) {
	gl.ReadBuffer(src)
}

func (c *coreContext) ReadPixels(x, y, width, height int, colorFormat, typ uint32, pixels []uint8) {
	gl.ReadPixels(
		int32(x),      //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(y),      //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(width),  //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(height), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		colorFormat,
		typ,
		gl.Ptr(pixels),
	)
}

func (c *coreContext) Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
}

func (c *coreContext) ShaderSource(shader Shader, source string) {
	csources, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(uint32(shader), 1, csources, nil)
}

func (c *coreContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	var ptr unsafe.Pointer
	if len(data) > 0 {
		ptr = gl.Ptr(data)
	}
	gl.TexImage2D(
		target,
		int32(level),       //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(colorFormat), //gosec:disable G115 -- colorFormat is an enum behind the scenes while the internal format is an int
		int32(width),       //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(height),      //gosec:disable G115 -- we are definitely fine with limiting the value range here
		0,
		colorFormat,
		typ,
		ptr,
	)
}

func (c *coreContext) TexParameteri(target, param uint32, value int32) {
	gl.TexParameteri(target, param, value)
}

func (c *coreContext) Uniform1f(uniform Uniform, v float32) {
	gl.Uniform1f(int32(uniform), v)
}

func (c *coreContext) Uniform1fv(uniform Uniform, v []float32) {
	gl.Uniform1fv(
		int32(uniform),
		int32(len(v)), //gosec:disable G115 -- v should definitely not contain more than two billion values
		&v[0],
	)
}

func (c *coreContext) Uniform1i(uniform Uniform, v int32) {
	gl.Uniform1i(int32(uniform), v)
}

func (c *coreContext) Uniform2f(uniform Uniform, v0, v1 float32) {
	gl.Uniform2f(int32(uniform), v0, v1)
}

func (c *coreContext) Uniform2fv(uniform Uniform, v []float32) {
	gl.Uniform2fv(
		int32(uniform),
		int32(len(v)/2), //gosec:disable G115 -- v should definitely not contain more than two billion values
		&v[0],
	)
}

func (c *coreContext) Uniform4f(uniform Uniform, v0, v1, v2, v3 float32) {
	gl.Uniform4f(int32(uniform), v0, v1, v2, v3)
}

func (c *coreContext) UseProgram(program Program) {
	gl.UseProgram(uint32(program))
}

func (c *coreContext) VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int) {
	gl.VertexAttribPointerWithOffset(
		uint32(attribute),
		int32(size), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		typ,
		normalized,
		int32(stride), //gosec:disable G115 -- we are definitely fine with limiting the value range here
		uintptr(offset),
	)
}

func (c *coreContext) Viewport(x, y, width, height int) {
	gl.Viewport(
		int32(x),      //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(y),      //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(width),  //gosec:disable G115 -- we are definitely fine with limiting the value range here
		int32(height), //gosec:disable G115 -- we are definitely fine with limiting the value range here
	)
}
