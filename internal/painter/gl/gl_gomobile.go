//go:build (android || ios || mobile) && (!wasm || !test_web_driver)

package gl

import (
	"math"

	"fyne.io/fyne/v2/internal/driver/mobile/gl"
)

const (
	arrayBuffer           = gl.ArrayBuffer
	bitColorBuffer        = gl.ColorBufferBit
	bitDepthBuffer        = gl.DepthBufferBit
	clampToEdge           = gl.ClampToEdge
	colorFormatRGBA       = gl.RGBA
	compileStatus         = gl.CompileStatus
	constantAlpha         = gl.ConstantAlpha
	float                 = gl.Float
	fragmentShader        = gl.FragmentShader
	front                 = gl.Front
	back                  = gl.Back
	glFalse               = gl.False
	linkStatus            = gl.LinkStatus
	maxTextureSizeParam   = gl.MaxTextureSize
	one                   = gl.One
	zero                  = gl.Zero
	oneMinusConstantAlpha = gl.OneMinusConstantAlpha
	oneMinusSrcAlpha      = gl.OneMinusSrcAlpha
	scissorTest           = gl.ScissorTest
	srcAlpha              = gl.SrcAlpha
	staticDraw            = gl.StaticDraw
	texture0              = gl.Texture0
	texture1              = gl.Texture1
	texture2D             = gl.Texture2D
	textureNearest        = gl.Nearest
	textureMinFilter      = gl.TextureMinFilter
	textureMagFilter      = gl.TextureMagFilter
	textureWrapS          = gl.TextureWrapS
	textureWrapT          = gl.TextureWrapT
	triangles             = gl.Triangles
	triangleStrip         = gl.TriangleStrip
	unsignedByte          = gl.UnsignedByte
	vertexShader          = gl.VertexShader
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

var (
	compiled          *programs // avoid multiple compilations with the re-used mobile GUI context
	noBuffer          = Buffer{}
	noProgram         = Program{}
	noShader          = Shader{}
	textureFilterToGL = [...]int32{gl.Linear, gl.Nearest, gl.Linear}
)

func (p *painter) Init() {
	glctx := p.contextProvider.Context().(gl.Context)
	glctx.Disable(gl.DepthTest)
	glctx.Enable(gl.Blend)
	p.ctx = &mobileContext{glContext: glctx}
	p.maxTextureSize = p.ctx.GetInteger(maxTextureSizeParam)
	p.blurSnap.texValid = false   // reset on context recreation; old texture IDs are no longer valid
	p.blurKernel.texValid = false // kernel texture must also be re-created
	if compiled == nil {
		compiled = p.compilePrograms()
	}

	p.programs = compiled
}

type mobileContext struct {
	glContext gl.Context
}

var _ context = (*mobileContext)(nil)

func (c *mobileContext) ActiveTexture(textureUnit uint32) {
	c.glContext.ActiveTexture(gl.Enum(textureUnit))
}

func (c *mobileContext) AttachShader(program Program, shader Shader) {
	c.glContext.AttachShader(gl.Program(program), gl.Shader(shader))
}

func (c *mobileContext) BindBuffer(target uint32, buf Buffer) {
	c.glContext.BindBuffer(gl.Enum(target), gl.Buffer(buf))
}

func (c *mobileContext) BindTexture(target uint32, texture Texture) {
	c.glContext.BindTexture(gl.Enum(target), gl.Texture(texture))
}

func (c *mobileContext) BlendColor(r, g, b, a float32) {
	c.glContext.BlendColor(r, g, b, a)
}

func (c *mobileContext) BlendFunc(srcFactor, destFactor uint32) {
	c.glContext.BlendFunc(gl.Enum(srcFactor), gl.Enum(destFactor))
}

func (c *mobileContext) BufferData(target uint32, points []float32, usage uint32) {
	data := toLEByteOrder(points...)
	c.glContext.BufferData(gl.Enum(target), data, gl.Enum(usage))
}

func (c *mobileContext) BufferSubData(target uint32, points []float32) {
	data := toLEByteOrder(points...)
	c.glContext.BufferSubData(gl.Enum(target), data)
}

func (c *mobileContext) Clear(mask uint32) {
	c.glContext.Clear(gl.Enum(mask))
}

func (c *mobileContext) ClearColor(r, g, b, a float32) {
	c.glContext.ClearColor(r, g, b, a)
}

func (c *mobileContext) CompileShader(shader Shader) {
	c.glContext.CompileShader(gl.Shader(shader))
}

func (c *mobileContext) CreateBuffer() Buffer {
	return Buffer(c.glContext.CreateBuffer())
}

func (c *mobileContext) CreateProgram() Program {
	return Program(c.glContext.CreateProgram())
}

func (c *mobileContext) CreateShader(typ uint32) Shader {
	return Shader(c.glContext.CreateShader(gl.Enum(typ)))
}

func (c *mobileContext) CreateTexture() (texture Texture) {
	return Texture(c.glContext.CreateTexture())
}

func (c *mobileContext) DeleteBuffer(buffer Buffer) {
	c.glContext.DeleteBuffer(gl.Buffer(buffer))
}

func (c *mobileContext) DeleteProgram(program Program) {
	c.glContext.DeleteProgram(gl.Program(program))
}

func (c *mobileContext) DeleteTexture(texture Texture) {
	c.glContext.DeleteTexture(gl.Texture(texture))
}

func (c *mobileContext) Disable(capability uint32) {
	c.glContext.Disable(gl.Enum(capability))
}

func (c *mobileContext) DrawArrays(mode uint32, first, count int) {
	c.glContext.DrawArrays(gl.Enum(mode), first, count)
}

func (c *mobileContext) Enable(capability uint32) {
	c.glContext.Enable(gl.Enum(capability))
}

func (c *mobileContext) EnableVertexAttribArray(attribute Attribute) {
	c.glContext.EnableVertexAttribArray(gl.Attrib(attribute))
}

func (c *mobileContext) GetAttribLocation(program Program, name string) Attribute {
	return Attribute(c.glContext.GetAttribLocation(gl.Program(program), name))
}

func (c *mobileContext) GetError() uint32 {
	return uint32(c.glContext.GetError())
}

func (c *mobileContext) GetInteger(pname uint32) int {
	return c.glContext.GetInteger(gl.Enum(pname))
}

func (c *mobileContext) GetProgrami(program Program, param uint32) int {
	return c.glContext.GetProgrami(gl.Program(program), gl.Enum(param))
}

func (c *mobileContext) GetProgramInfoLog(program Program) string {
	return c.glContext.GetProgramInfoLog(gl.Program(program))
}

func (c *mobileContext) GetShaderi(shader Shader, param uint32) int {
	return c.glContext.GetShaderi(gl.Shader(shader), gl.Enum(param))
}

func (c *mobileContext) GetShaderInfoLog(shader Shader) string {
	return c.glContext.GetShaderInfoLog(gl.Shader(shader))
}

func (c *mobileContext) GetUniformLocation(program Program, name string) Uniform {
	return Uniform(c.glContext.GetUniformLocation(gl.Program(program), name))
}

func (c *mobileContext) LinkProgram(program Program) {
	c.glContext.LinkProgram(gl.Program(program))
}

func (c *mobileContext) CopyTexSubImage2D(target uint32, level, xoffset, yoffset, x, y, width, height int) {
	c.glContext.CopyTexSubImage2D(gl.Enum(target), level, xoffset, yoffset, x, y, width, height)
}

func (c *mobileContext) ReadBuffer(_ uint32) {
}

func (c *mobileContext) ReadPixels(x, y, width, height int, colorFormat, typ uint32, pixels []uint8) {
	c.glContext.ReadPixels(pixels, x, y, width, height, gl.Enum(colorFormat), gl.Enum(typ))
}

func (c *mobileContext) Scissor(x, y, w, h int32) {
	c.glContext.Scissor(x, y, w, h)
}

func (c *mobileContext) ShaderSource(shader Shader, source string) {
	c.glContext.ShaderSource(gl.Shader(shader), source)
}

func (c *mobileContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	c.glContext.TexImage2D(
		gl.Enum(target),
		level,
		int(colorFormat),
		width,
		height,
		gl.Enum(colorFormat),
		gl.Enum(typ),
		data,
	)
}

func (c *mobileContext) TexParameteri(target, param uint32, value int32) {
	c.glContext.TexParameteri(gl.Enum(target), gl.Enum(param), int(value))
}

func (c *mobileContext) Uniform1f(uniform Uniform, v float32) {
	c.glContext.Uniform1f(gl.Uniform(uniform), v)
}

func (c *mobileContext) Uniform1fv(uniform Uniform, v []float32) {
	c.glContext.Uniform1fv(gl.Uniform(uniform), v)
}

func (c *mobileContext) Uniform1i(uniform Uniform, v int32) {
	c.glContext.Uniform1i(gl.Uniform(uniform), int(v))
}

func (c *mobileContext) Uniform2f(uniform Uniform, v0, v1 float32) {
	c.glContext.Uniform2f(gl.Uniform(uniform), v0, v1)
}

func (c *mobileContext) Uniform2fv(uniform Uniform, v []float32) {
	c.glContext.Uniform2fv(gl.Uniform(uniform), v)
}

func (c *mobileContext) Uniform4f(uniform Uniform, v0, v1, v2, v3 float32) {
	c.glContext.Uniform4f(gl.Uniform(uniform), v0, v1, v2, v3)
}

func (c *mobileContext) UseProgram(program Program) {
	c.glContext.UseProgram(gl.Program(program))
}

func (c *mobileContext) VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int) {
	c.glContext.VertexAttribPointer(gl.Attrib(attribute), size, gl.Enum(typ), normalized, stride, offset)
}

func (c *mobileContext) Viewport(x, y, width, height int) {
	c.glContext.Viewport(x, y, width, height)
}

// toLEByteOrder returns the byte representation of float32 values in little endian byte order.
func toLEByteOrder(values ...float32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		b[4*i+0] = byte(u >> 0)
		b[4*i+1] = byte(u >> 8)
		b[4*i+2] = byte(u >> 16)
		b[4*i+3] = byte(u >> 24)
	}
	return b
}
