// +build !js

package gl

import (
	"log"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

// ContextWatcher is this library's context watcher, satisfying glfw.ContextWatcher interface.
// It must be notified when context is made current or detached.
var ContextWatcher = new(contextWatcher)

type contextWatcher struct {
	initGL bool
}

func (cw *contextWatcher) OnMakeCurrent(context interface{}) {
	if !cw.initGL {
		// Initialise gl bindings using the current context.
		err := gl.Init()
		if err != nil {
			log.Fatalln("gl.Init:", err)
		}
		cw.initGL = true
	}
}
func (contextWatcher) OnDetach() {}

// ActiveTexture sets the active texture unit.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
func ActiveTexture(texture Enum) {
	gl.ActiveTexture(uint32(texture))
}

// AttachShader attaches a shader to a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
func AttachShader(p Program, s Shader) {
	gl.AttachShader(p.Value, s.Value)
}

// BindAttribLocation binds a vertex attribute index with a named
// variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
func BindAttribLocation(p Program, a Attrib, name string) {
	gl.BindAttribLocation(p.Value, uint32(a.Value), gl.Str(name+"\x00"))
}

// BindBuffer binds a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
func BindBuffer(target Enum, b Buffer) {
	gl.BindBuffer(uint32(target), b.Value)
}

// BindFramebuffer binds a framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
func BindFramebuffer(target Enum, fb Framebuffer) {
	gl.BindFramebuffer(uint32(target), fb.Value)
}

// BindRenderbuffer binds a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
func BindRenderbuffer(target Enum, rb Renderbuffer) {
	gl.BindRenderbuffer(uint32(target), rb.Value)
}

// BindTexture binds a texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
func BindTexture(target Enum, t Texture) {
	gl.BindTexture(uint32(target), t.Value)
}

// BlendColor sets the blend color.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
func BlendColor(red, green, blue, alpha float32) {
	gl.BlendColor(red, green, blue, alpha)
}

// BlendEquation sets both RGB and alpha blend equations.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
func BlendEquation(mode Enum) {
	gl.BlendEquation(uint32(mode))
}

// BlendEquationSeparate sets RGB and alpha blend equations separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	gl.BlendEquationSeparate(uint32(modeRGB), uint32(modeAlpha))
}

// BlendFunc sets the pixel blending factors.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
func BlendFunc(sfactor, dfactor Enum) {
	gl.BlendFunc(uint32(sfactor), uint32(dfactor))
}

// BlendFunc sets the pixel RGB and alpha blending factors separately.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFuncSeparate.xhtml
func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	gl.BlendFuncSeparate(uint32(sfactorRGB), uint32(dfactorRGB), uint32(sfactorAlpha), uint32(dfactorAlpha))
}

// BufferData creates a new data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func BufferData(target Enum, src []byte, usage Enum) {
	gl.BufferData(uint32(target), int(len(src)), gl.Ptr(&src[0]), uint32(usage))
}

// BufferInit creates a new unitialized data store for the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
func BufferInit(target Enum, size int, usage Enum) {
	gl.BufferData(uint32(target), size, nil, uint32(usage))
}

// BufferSubData sets some of data in the bound buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
func BufferSubData(target Enum, offset int, data []byte) {
	gl.BufferSubData(uint32(target), offset, int(len(data)), gl.Ptr(&data[0]))
}

// CheckFramebufferStatus reports the completeness status of the
// active framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
func CheckFramebufferStatus(target Enum) Enum {
	return Enum(gl.CheckFramebufferStatus(uint32(target)))
}

// Clear clears the window.
//
// The behavior of Clear is influenced by the pixel ownership test,
// the scissor test, dithering, and the buffer writemasks.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
func Clear(mask Enum) {
	gl.Clear(uint32(mask))
}

// ClearColor specifies the RGBA values used to clear color buffers.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
func ClearColor(red, green, blue, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

// ClearDepthf sets the depth value used to clear the depth buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
func ClearDepthf(d float32) {
	gl.ClearDepthf(d)
}

// ClearStencil sets the index used to clear the stencil buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearStencil.xhtml
func ClearStencil(s int) {
	gl.ClearStencil(int32(s))
}

// ColorMask specifies whether color components in the framebuffer
// can be written.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glColorMask.xhtml
func ColorMask(red, green, blue, alpha bool) {
	gl.ColorMask(red, green, blue, alpha)
}

// CompileShader compiles the source code of s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
func CompileShader(s Shader) {
	gl.CompileShader(s.Value)
}

// CompressedTexImage2D writes a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexImage2D.xhtml
func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	gl.CompressedTexImage2D(uint32(target), int32(level), uint32(internalformat), int32(width), int32(height), int32(border), int32(len(data)), gl.Ptr(data))
}

// CompressedTexSubImage2D writes a subregion of a compressed 2D texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexSubImage2D.xhtml
func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	gl.CompressedTexSubImage2D(uint32(target), int32(level), int32(xoffset), int32(yoffset), int32(width), int32(height), uint32(format), int32(len(data)), gl.Ptr(data))
}

// CopyTexImage2D writes a 2D texture from the current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexImage2D.xhtml
func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	gl.CopyTexImage2D(uint32(target), int32(level), uint32(internalformat), int32(x), int32(y), int32(width), int32(height), int32(border))
}

// CopyTexSubImage2D writes a 2D texture subregion from the
// current framebuffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexSubImage2D.xhtml
func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	gl.CopyTexSubImage2D(uint32(target), int32(level), int32(xoffset), int32(yoffset), int32(x), int32(y), int32(width), int32(height))
}

// CreateBuffer creates a buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
func CreateBuffer() Buffer {
	var b Buffer
	gl.GenBuffers(1, &b.Value)
	return b
}

// CreateFramebuffer creates a framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
func CreateFramebuffer() Framebuffer {
	var b Framebuffer
	gl.GenFramebuffers(1, &b.Value)
	return b
}

// CreateProgram creates a new empty program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
func CreateProgram() Program {
	return Program{Value: uint32(gl.CreateProgram())}
}

// CreateRenderbuffer create a renderbuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenRenderbuffers.xhtml
func CreateRenderbuffer() Renderbuffer {
	var b Renderbuffer
	gl.GenRenderbuffers(1, &b.Value)
	return b
}

// CreateShader creates a new empty shader object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
func CreateShader(ty Enum) Shader {
	return Shader{Value: uint32(gl.CreateShader(uint32(ty)))}
}

// CreateTexture creates a texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
func CreateTexture() Texture {
	var t Texture
	gl.GenTextures(1, &t.Value)
	return t
}

// CullFace specifies which polygons are candidates for culling.
//
// Valid modes: FRONT, BACK, FRONT_AND_BACK.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
func CullFace(mode Enum) {
	gl.CullFace(uint32(mode))
}

// DeleteBuffer deletes the given buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
func DeleteBuffer(v Buffer) {
	gl.DeleteBuffers(1, &v.Value)
}

// DeleteFramebuffer deletes the given framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
func DeleteFramebuffer(v Framebuffer) {
	gl.DeleteFramebuffers(1, &v.Value)
}

// DeleteProgram deletes the given program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
func DeleteProgram(p Program) {
	gl.DeleteProgram(p.Value)
}

// DeleteRenderbuffer deletes the given render buffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteRenderbuffers.xhtml
func DeleteRenderbuffer(v Renderbuffer) {
	gl.DeleteRenderbuffers(1, &v.Value)
}

// DeleteShader deletes shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteShader.xhtml
func DeleteShader(s Shader) {
	gl.DeleteShader(s.Value)
}

// DeleteTexture deletes the given texture object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
func DeleteTexture(v Texture) {
	gl.DeleteTextures(1, &v.Value)
}

// DepthFunc sets the function used for depth buffer comparisons.
//
// Valid fn values:
//	NEVER
//	LESS
//	EQUAL
//	LEQUAL
//	GREATER
//	NOTEQUAL
//	GEQUAL
//	ALWAYS
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthFunc.xhtml
func DepthFunc(fn Enum) {
	gl.DepthFunc(uint32(fn))
}

// DepthMask sets the depth buffer enabled for writing.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
func DepthMask(flag bool) {
	gl.DepthMask(flag)
}

// DepthRangef sets the mapping from normalized device coordinates to
// window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthRangef.xhtml
func DepthRangef(n, f float32) {
	gl.DepthRangef(n, f)
}

// DetachShader detaches the shader s from the program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDetachShader.xhtml
func DetachShader(p Program, s Shader) {
	gl.DetachShader(p.Value, s.Value)
}

// Disable disables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
func Disable(cap Enum) {
	gl.Disable(uint32(cap))
}

// DisableVertexAttribArray disables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
func DisableVertexAttribArray(a Attrib) {
	gl.DisableVertexAttribArray(uint32(a.Value))
}

// DrawArrays renders geometric primitives from the bound data.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
func DrawArrays(mode Enum, first, count int) {
	gl.DrawArrays(uint32(mode), int32(first), int32(count))
}

// DrawElements renders primitives from a bound buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
func DrawElements(mode Enum, count int, ty Enum, offset int) {
	gl.DrawElements(uint32(mode), int32(count), uint32(ty), gl.PtrOffset(offset))
}

// Enable enables various GL capabilities.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
func Enable(cap Enum) {
	gl.Enable(uint32(cap))
}

// EnableVertexAttribArray enables a vertex attribute array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
func EnableVertexAttribArray(a Attrib) {
	gl.EnableVertexAttribArray(uint32(a.Value))
}

// Finish blocks until the effects of all previously called GL
// commands are complete.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
func Finish() {
	gl.Finish()
}

// Flush empties all buffers. It does not block.
//
// An OpenGL implementation may buffer network communication,
// the command stream, or data inside the graphics accelerator.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
func Flush() {
	gl.Flush()
}

// FramebufferRenderbuffer attaches rb to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferRenderbuffer.xhtml
func FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	gl.FramebufferRenderbuffer(uint32(target), uint32(attachment), uint32(rbTarget), rb.Value)
}

// FramebufferTexture2D attaches the t to the current frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	gl.FramebufferTexture2D(uint32(target), uint32(attachment), uint32(texTarget), t.Value, int32(level))
}

// FrontFace defines which polygons are front-facing.
//
// Valid modes: CW, CCW.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
func FrontFace(mode Enum) {
	gl.FrontFace(uint32(mode))
}

// GenerateMipmap generates mipmaps for the current texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
func GenerateMipmap(target Enum) {
	gl.GenerateMipmap(uint32(target))
}

// GetActiveAttrib returns details about an active attribute variable.
// A value of 0 for index selects the first active attribute variable.
// Permissible values for index range from 0 to the number of active
// attribute variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	var length, si int32
	var typ uint32
	name = strings.Repeat("\x00", 256)
	cname := gl.Str(name)
	gl.GetActiveAttrib(p.Value, uint32(index), int32(len(name)-1), &length, &si, &typ, cname)
	name = name[:strings.IndexRune(name, 0)]
	return name, int(si), Enum(typ)
}

// GetActiveUniform returns details about an active uniform variable.
// A value of 0 for index selects the first active uniform variable.
// Permissible values for index range from 0 to the number of active
// uniform variables minus 1.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	var length, si int32
	var typ uint32
	name = strings.Repeat("\x00", 256)
	cname := gl.Str(name)
	gl.GetActiveUniform(p.Value, uint32(index), int32(len(name)-1), &length, &si, &typ, cname)
	name = name[:strings.IndexRune(name, 0)]
	return name, int(si), Enum(typ)
}

// GetAttachedShaders returns the shader objects attached to program p.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttachedShaders.xhtml
func GetAttachedShaders(p Program) []Shader {
	log.Println("GetAttachedShaders: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	shadersLen := GetProgrami(p, ATTACHED_SHADERS)
	var n int32
	buf := make([]uint32, shadersLen)
	gl.GetAttachedShaders(uint32(p.Value), int32(shadersLen), &n, &buf[0])
	buf = buf[:int(n)]
	shaders := make([]Shader, int(n))
	for i, s := range buf {
		shaders[i] = Shader{Value: uint32(s)}
	}
	return shaders
}

// GetAttribLocation returns the location of an attribute variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
func GetAttribLocation(p Program, name string) Attrib {
	return Attrib{Value: uint(gl.GetAttribLocation(p.Value, gl.Str(name+"\x00")))}
}

// GetBooleanv returns the boolean values of parameter pname.
//
// Many boolean parameters can be queried more easily using IsEnabled.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func GetBooleanv(dst []bool, pname Enum) {
	gl.GetBooleanv(uint32(pname), &dst[0])
}

// GetFloatv returns the float values of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func GetFloatv(dst []float32, pname Enum) {
	gl.GetFloatv(uint32(pname), &dst[0])
}

// GetIntegerv returns the int values of parameter pname.
//
// Single values may be queried more easily using GetInteger.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func GetIntegerv(pname Enum, data []int32) {
	gl.GetIntegerv(uint32(pname), &data[0])
}

// GetInteger returns the int value of parameter pname.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
func GetInteger(pname Enum) int {
	var data int32
	gl.GetIntegerv(uint32(pname), &data)
	return int(data)
}

// GetBufferParameteri returns a parameter for the active buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetBufferParameteriv.xhtml
func GetBufferParameteri(target, pname Enum) int {
	var params int32
	gl.GetBufferParameteriv(uint32(target), uint32(pname), &params)
	return int(params)
}

// GetError returns the next error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
func GetError() Enum {
	return Enum(gl.GetError())
}

// GetBoundFramebuffer returns the currently bound framebuffer.
// Use this method instead of gl.GetInteger(gl.FRAMEBUFFER_BINDING) to
// enable support on all platforms
func GetBoundFramebuffer() Framebuffer {
	var b int32
	gl.GetIntegerv(FRAMEBUFFER_BINDING, &b)
	return Framebuffer{Value: uint32(b)}
}

// GetFramebufferAttachmentParameteri returns attachment parameters
// for the active framebuffer object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetFramebufferAttachmentParameteriv.xhtml
func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	log.Println("GetFramebufferAttachmentParameteri: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	var param int32
	gl.GetFramebufferAttachmentParameteriv(uint32(target), uint32(attachment), uint32(pname), &param)
	return int(param)
}

// GetProgrami returns a parameter value for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
func GetProgrami(p Program, pname Enum) int {
	var result int32
	gl.GetProgramiv(p.Value, uint32(pname), &result)
	return int(result)
}

// GetProgramInfoLog returns the information log for a program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
func GetProgramInfoLog(p Program) string {
	var logLength int32
	gl.GetProgramiv(p.Value, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p.Value, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}

// GetRenderbufferParameteri returns a parameter value for a render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetRenderbufferParameteriv.xhtml
func GetRenderbufferParameteri(target, pname Enum) int {
	log.Println("GetRenderbufferParameteri: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	var result int32
	gl.GetRenderbufferParameteriv(uint32(target), uint32(pname), &result)
	return int(result)
}

// GetRenderbufferParameteri returns a parameter value for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
func GetShaderi(s Shader, pname Enum) int {
	var result int32
	gl.GetShaderiv(s.Value, uint32(pname), &result)
	return int(result)
}

// GetShaderInfoLog returns the information log for a shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
func GetShaderInfoLog(s Shader) string {
	var logLength int32
	gl.GetShaderiv(s.Value, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s.Value, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}

// GetShaderPrecisionFormat returns range and precision limits for
// shader types.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderPrecisionFormat.xhtml
func GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	log.Println("GetShaderPrecisionFormat: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	var cRange [2]int32
	var cPrecision int32

	gl.GetShaderPrecisionFormat(uint32(shadertype), uint32(precisiontype), &cRange[0], &cPrecision)
	return int(cRange[0]), int(cRange[1]), int(cPrecision)
}

// GetShaderSource returns source code of shader s.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderSource.xhtml
func GetShaderSource(s Shader) string {
	log.Println("GetShaderSource: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	sourceLen := GetShaderi(s, gl.SHADER_SOURCE_LENGTH)
	if sourceLen == 0 {
		return ""
	}
	buf := make([]byte, sourceLen)
	gl.GetShaderSource(s.Value, int32(sourceLen), nil, &buf[0])
	return string(buf)
}

// GetString reports current GL state.
//
// Valid name values:
//	EXTENSIONS
//	RENDERER
//	SHADING_LANGUAGE_VERSION
//	VENDOR
//	VERSION
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetString.xhtml
func GetString(pname Enum) string {
	return gl.GoStr(gl.GetString(uint32(pname)))
}

// GetTexParameterfv returns the float values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func GetTexParameterfv(dst []float32, target, pname Enum) {
	gl.GetTexParameterfv(uint32(target), uint32(pname), &dst[0])
}

// GetTexParameteriv returns the int values of a texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
func GetTexParameteriv(dst []int32, target, pname Enum) {
	gl.GetTexParameteriv(uint32(target), uint32(pname), &dst[0])
}

// GetUniformfv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func GetUniformfv(dst []float32, src Uniform, p Program) {
	gl.GetUniformfv(p.Value, src.Value, &dst[0])
}

// GetUniformiv returns the float values of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
func GetUniformiv(dst []int32, src Uniform, p Program) {
	gl.GetUniformiv(p.Value, src.Value, &dst[0])
}

// GetUniformLocation returns the location of a uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
func GetUniformLocation(p Program, name string) Uniform {
	return Uniform{Value: gl.GetUniformLocation(p.Value, gl.Str(name+"\x00"))}
}

// GetVertexAttribf reads the float value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func GetVertexAttribf(src Attrib, pname Enum) float32 {
	log.Println("GetVertexAttribf: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	var result float32
	gl.GetVertexAttribfv(uint32(src.Value), uint32(pname), &result)
	return result
}

// GetVertexAttribfv reads float values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	gl.GetVertexAttribfv(uint32(src.Value), uint32(pname), &dst[0])
}

// GetVertexAttribi reads the int value of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func GetVertexAttribi(src Attrib, pname Enum) int32 {
	var result int32
	gl.GetVertexAttribiv(uint32(src.Value), uint32(pname), &result)
	return result
}

// GetVertexAttribiv reads int values of a vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
func GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	gl.GetVertexAttribiv(uint32(src.Value), uint32(pname), &dst[0])
}

// Hint sets implementation-specific modes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
func Hint(target, mode Enum) {
	gl.Hint(uint32(target), uint32(mode))
}

// IsBuffer reports if b is a valid buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsBuffer.xhtml
func IsBuffer(b Buffer) bool {
	return gl.IsBuffer(b.Value)
}

// IsEnabled reports if cap is an enabled capability.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsEnabled.xhtml
func IsEnabled(cap Enum) bool {
	return gl.IsEnabled(uint32(cap))
}

// IsFramebuffer reports if fb is a valid frame buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsFramebuffer.xhtml
func IsFramebuffer(fb Framebuffer) bool {
	return gl.IsFramebuffer(fb.Value)
}

// IsProgram reports if p is a valid program object.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsProgram.xhtml
func IsProgram(p Program) bool {
	return gl.IsProgram(p.Value)
}

// IsRenderbuffer reports if rb is a valid render buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsRenderbuffer.xhtml
func IsRenderbuffer(rb Renderbuffer) bool {
	return gl.IsRenderbuffer(rb.Value)
}

// IsShader reports if s is valid shader.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsShader.xhtml
func IsShader(s Shader) bool {
	return gl.IsShader(s.Value)
}

// IsTexture reports if t is a valid texture.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsTexture.xhtml
func IsTexture(t Texture) bool {
	return gl.IsTexture(t.Value)
}

// LineWidth specifies the width of lines.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
func LineWidth(width float32) {
	gl.LineWidth(width)
}

// LinkProgram links the specified program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
func LinkProgram(p Program) {
	gl.LinkProgram(p.Value)
}

// PixelStorei sets pixel storage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPixelStorei.xhtml
func PixelStorei(pname Enum, param int32) {
	gl.PixelStorei(uint32(pname), param)
}

// PolygonOffset sets the scaling factors for depth offsets.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glPolygonOffset.xhtml
func PolygonOffset(factor, units float32) {
	gl.PolygonOffset(factor, units)
}

// ReadPixels returns pixel data from a buffer.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	gl.ReadPixels(int32(x), int32(y), int32(width), int32(height), uint32(format), uint32(ty), gl.Ptr(&dst[0]))
}

// ReleaseShaderCompiler frees resources allocated by the shader compiler.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glReleaseShaderCompiler.xhtml
func ReleaseShaderCompiler() {
	gl.ReleaseShaderCompiler()
}

// RenderbufferStorage establishes the data storage, format, and
// dimensions of a renderbuffer object's image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glRenderbufferStorage.xhtml
func RenderbufferStorage(target, internalFormat Enum, width, height int) {
	gl.RenderbufferStorage(uint32(target), uint32(internalFormat), int32(width), int32(height))
}

// SampleCoverage sets multisample coverage parameters.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glSampleCoverage.xhtml
func SampleCoverage(value float32, invert bool) {
	gl.SampleCoverage(value, invert)
}

// Scissor defines the scissor box rectangle, in window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
func Scissor(x, y, width, height int32) {
	gl.Scissor(x, y, width, height)
}

// ShaderSource sets the source code of s to the given source code.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
func ShaderSource(s Shader, src string) {
	glsource, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s.Value, 1, glsource, nil)
	free()
}

// StencilFunc sets the front and back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFunc.xhtml
func StencilFunc(fn Enum, ref int, mask uint32) {
	gl.StencilFunc(uint32(fn), int32(ref), mask)
}

// StencilFunc sets the front or back stencil test reference value.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFuncSeparate.xhtml
func StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	gl.StencilFuncSeparate(uint32(face), uint32(fn), int32(ref), mask)
}

// StencilMask controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMask.xhtml
func StencilMask(mask uint32) {
	gl.StencilMask(mask)
}

// StencilMaskSeparate controls the writing of bits in the stencil planes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMaskSeparate.xhtml
func StencilMaskSeparate(face Enum, mask uint32) {
	gl.StencilMaskSeparate(uint32(face), mask)
}

// StencilOp sets front and back stencil test actions.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOp.xhtml
func StencilOp(fail, zfail, zpass Enum) {
	gl.StencilOp(uint32(fail), uint32(zfail), uint32(zpass))
}

// StencilOpSeparate sets front or back stencil tests.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOpSeparate.xhtml
func StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	gl.StencilOpSeparate(uint32(face), uint32(sfail), uint32(dpfail), uint32(dppass))
}

// TexImage2D writes a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func TexImage2D(target Enum, level int, width, height int, format Enum, ty Enum, data []byte) {
	p := unsafe.Pointer(nil)
	if len(data) > 0 {
		p = gl.Ptr(&data[0])
	}
	gl.TexImage2D(uint32(target), int32(level), int32(format), int32(width), int32(height), 0, uint32(format), uint32(ty), p)
}

// TexSubImage2D writes a subregion of a 2D texture image.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexSubImage2D.xhtml
func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	gl.TexSubImage2D(uint32(target), int32(level), int32(x), int32(y), int32(width), int32(height), uint32(format), uint32(ty), gl.Ptr(&data[0]))
}

// TexParameterf sets a float texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func TexParameterf(target, pname Enum, param float32) {
	gl.TexParameterf(uint32(target), uint32(pname), param)
}

// TexParameterfv sets a float texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func TexParameterfv(target, pname Enum, params []float32) {
	gl.TexParameterfv(uint32(target), uint32(pname), &params[0])
}

// TexParameteri sets an integer texture parameter.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func TexParameteri(target, pname Enum, param int) {
	gl.TexParameteri(uint32(target), uint32(pname), int32(param))
}

// TexParameteriv sets an integer texture parameter array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
func TexParameteriv(target, pname Enum, params []int32) {
	gl.TexParameteriv(uint32(target), uint32(pname), &params[0])
}

// Uniform1f writes a float uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform1f(dst Uniform, v float32) {
	gl.Uniform1f(dst.Value, v)
}

// Uniform1fv writes a [len(src)]float uniform array.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform1fv(dst Uniform, src []float32) {
	gl.Uniform1fv(dst.Value, int32(len(src)), &src[0])
}

// Uniform1i writes an int uniform variable.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform1i(dst Uniform, v int) {
	gl.Uniform1i(dst.Value, int32(v))
}

// Uniform1iv writes a int uniform array of len(src) elements.
//
// Uniform1i and Uniform1iv are the only two functions that may be used
// to load uniform variables defined as sampler types. Loading samplers
// with any other function will result in a INVALID_OPERATION error.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform1iv(dst Uniform, src []int32) {
	gl.Uniform1iv(dst.Value, int32(len(src)), &src[0])
}

// Uniform2f writes a vec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform2f(dst Uniform, v0, v1 float32) {
	gl.Uniform2f(dst.Value, v0, v1)
}

// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform2fv(dst Uniform, src []float32) {
	gl.Uniform2fv(dst.Value, int32(len(src)/2), &src[0])
}

// Uniform2i writes an ivec2 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform2i(dst Uniform, v0, v1 int) {
	gl.Uniform2i(dst.Value, int32(v0), int32(v1))
}

// Uniform2iv writes an ivec2 uniform array of len(src)/2 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform2iv(dst Uniform, src []int32) {
	gl.Uniform2iv(dst.Value, int32(len(src)/2), &src[0])
}

// Uniform3f writes a vec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform3f(dst Uniform, v0, v1, v2 float32) {
	gl.Uniform3f(dst.Value, v0, v1, v2)
}

// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform3fv(dst Uniform, src []float32) {
	gl.Uniform3fv(dst.Value, int32(len(src)/3), &src[0])
}

// Uniform3i writes an ivec3 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform3i(dst Uniform, v0, v1, v2 int32) {
	gl.Uniform3i(dst.Value, v0, v1, v2)
}

// Uniform3iv writes an ivec3 uniform array of len(src)/3 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform3iv(dst Uniform, src []int32) {
	gl.Uniform3iv(dst.Value, int32(len(src)/3), &src[0])
}

// Uniform4f writes a vec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	gl.Uniform4f(dst.Value, v0, v1, v2, v3)
}

// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform4fv(dst Uniform, src []float32) {
	gl.Uniform4fv(dst.Value, int32(len(src)/4), &src[0])
}

// Uniform4i writes an ivec4 uniform variable.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	gl.Uniform4i(dst.Value, v0, v1, v2, v3)
}

// Uniform4i writes an ivec4 uniform array of len(src)/4 elements.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func Uniform4iv(dst Uniform, src []int32) {
	gl.Uniform4iv(dst.Value, int32(len(src)/4), &src[0])
}

// UniformMatrix2fv writes 2x2 matrices. Each matrix uses four
// float32 values, so the number of matrices written is len(src)/4.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func UniformMatrix2fv(dst Uniform, src []float32) {
	gl.UniformMatrix2fv(dst.Value, int32(len(src)/(2*2)), false, &src[0])
}

// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
// float32 values, so the number of matrices written is len(src)/9.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func UniformMatrix3fv(dst Uniform, src []float32) {
	gl.UniformMatrix3fv(dst.Value, int32(len(src)/(3*3)), false, &src[0])
}

// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
// float32 values, so the number of matrices written is len(src)/16.
//
// Each matrix must be supplied in column major order.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
func UniformMatrix4fv(dst Uniform, src []float32) {
	gl.UniformMatrix4fv(dst.Value, int32(len(src)/(4*4)), false, &src[0])
}

// UseProgram sets the active program.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
func UseProgram(p Program) {
	gl.UseProgram(p.Value)
}

// ValidateProgram checks to see whether the executables contained in
// program can execute given the current OpenGL state.
//
// Typically only used for debugging.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glValidateProgram.xhtml
func ValidateProgram(p Program) {
	gl.ValidateProgram(uint32(p.Value))
}

// VertexAttrib1f writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib1f(dst Attrib, x float32) {
	gl.VertexAttrib1f(uint32(dst.Value), x)
}

// VertexAttrib1fv writes a float vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib1fv(dst Attrib, src []float32) {
	gl.VertexAttrib1fv(uint32(dst.Value), &src[0])
}

// VertexAttrib2f writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib2f(dst Attrib, x, y float32) {
	gl.VertexAttrib2f(uint32(dst.Value), x, y)
}

// VertexAttrib2fv writes a vec2 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib2fv(dst Attrib, src []float32) {
	gl.VertexAttrib2fv(uint32(dst.Value), &src[0])
}

// VertexAttrib3f writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib3f(dst Attrib, x, y, z float32) {
	gl.VertexAttrib3f(uint32(dst.Value), x, y, z)
}

// VertexAttrib3fv writes a vec3 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib3fv(dst Attrib, src []float32) {
	gl.VertexAttrib3fv(uint32(dst.Value), &src[0])
}

// VertexAttrib4f writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	gl.VertexAttrib4f(uint32(dst.Value), x, y, z, w)
}

// VertexAttrib4fv writes a vec4 vertex attribute.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
func VertexAttrib4fv(dst Attrib, src []float32) {
	gl.VertexAttrib4fv(uint32(dst.Value), &src[0])
}

// VertexAttribPointer uses a bound buffer to define vertex attribute data.
//
// Direct use of VertexAttribPointer to load data into OpenGL is not
// supported via the Go bindings. Instead, use BindBuffer with an
// ARRAY_BUFFER and then fill it using BufferData.
//
// The size argument specifies the number of components per attribute,
// between 1-4. The stride argument specifies the byte offset between
// consecutive vertex attributes.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttribPointer.xhtml
func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	gl.VertexAttribPointer(uint32(dst.Value), int32(size), uint32(ty), normalized, int32(stride), gl.PtrOffset(offset))
}

// Viewport sets the viewport, an affine transformation that
// normalizes device coordinates to window coordinates.
//
// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
func Viewport(x, y, width, height int) {
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}
