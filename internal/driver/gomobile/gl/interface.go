// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// Context is an OpenGL ES context.
//
// A Context has a method for every GL function supported by ES 2 or later.
// In a program compiled with ES 3 support, a Context is also a Context3.
// For example, a program can:
//
//	func f(glctx gl.Context) {
//		glctx.(gl.Context3).BlitFramebuffer(...)
//	}
//
// Calls are not safe for concurrent use. However calls can be made from
// any goroutine, the gl package removes the notion of thread-local
// context.
//
// Contexts are independent. Two contexts can be used concurrently.
type Context interface {
	// ActiveTexture sets the active texture unit.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
	ActiveTexture(texture Enum)

	// AttachShader attaches a shader to a program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
	AttachShader(p Program, s Shader)

	// BindAttribLocation binds a vertex attribute index with a named
	// variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindAttribLocation.xhtml
	BindAttribLocation(p Program, a Attrib, name string)

	// BindBuffer binds a buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
	BindBuffer(target Enum, b Buffer)

	// BindFramebuffer binds a framebuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindFramebuffer.xhtml
	BindFramebuffer(target Enum, fb Framebuffer)

	// BindRenderbuffer binds a render buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindRenderbuffer.xhtml
	BindRenderbuffer(target Enum, rb Renderbuffer)

	// BindTexture binds a texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
	BindTexture(target Enum, t Texture)

	// BindVertexArray binds a vertex array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
	BindVertexArray(rb VertexArray)

	// BlendColor sets the blend color.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
	BlendColor(red, green, blue, alpha float32)

	// BlendEquation sets both RGB and alpha blend equations.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquation.xhtml
	BlendEquation(mode Enum)

	// BlendEquationSeparate sets RGB and alpha blend equations separately.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendEquationSeparate.xhtml
	BlendEquationSeparate(modeRGB, modeAlpha Enum)

	// BlendFunc sets the pixel blending factors.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
	BlendFunc(sfactor, dfactor Enum)

	// BlendFunc sets the pixel RGB and alpha blending factors separately.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFuncSeparate.xhtml
	BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum)

	// BufferData creates a new data store for the bound buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
	BufferData(target Enum, src []byte, usage Enum)

	// BufferInit creates a new uninitialized data store for the bound buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
	BufferInit(target Enum, size int, usage Enum)

	// BufferSubData sets some of data in the bound buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferSubData.xhtml
	BufferSubData(target Enum, offset int, data []byte)

	// CheckFramebufferStatus reports the completeness status of the
	// active framebuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCheckFramebufferStatus.xhtml
	CheckFramebufferStatus(target Enum) Enum

	// Clear clears the window.
	//
	// The behavior of Clear is influenced by the pixel ownership test,
	// the scissor test, dithering, and the buffer writemasks.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
	Clear(mask Enum)

	// ClearColor specifies the RGBA values used to clear color buffers.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
	ClearColor(red, green, blue, alpha float32)

	// ClearDepthf sets the depth value used to clear the depth buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearDepthf.xhtml
	ClearDepthf(d float32)

	// ClearStencil sets the index used to clear the stencil buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearStencil.xhtml
	ClearStencil(s int)

	// ColorMask specifies whether color components in the framebuffer
	// can be written.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glColorMask.xhtml
	ColorMask(red, green, blue, alpha bool)

	// CompileShader compiles the source code of s.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
	CompileShader(s Shader)

	// CompressedTexImage2D writes a compressed 2D texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexImage2D.xhtml
	CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte)

	// CompressedTexSubImage2D writes a subregion of a compressed 2D texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompressedTexSubImage2D.xhtml
	CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte)

	// CopyTexImage2D writes a 2D texture from the current framebuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexImage2D.xhtml
	CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int)

	// CopyTexSubImage2D writes a 2D texture subregion from the
	// current framebuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCopyTexSubImage2D.xhtml
	CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int)

	// CreateBuffer creates a buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
	CreateBuffer() Buffer

	// CreateFramebuffer creates a framebuffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenFramebuffers.xhtml
	CreateFramebuffer() Framebuffer

	// CreateProgram creates a new empty program object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
	CreateProgram() Program

	// CreateRenderbuffer create a renderbuffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenRenderbuffers.xhtml
	CreateRenderbuffer() Renderbuffer

	// CreateShader creates a new empty shader object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
	CreateShader(ty Enum) Shader

	// CreateTexture creates a texture object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
	CreateTexture() Texture

	// CreateTVertexArray creates a vertex array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
	CreateVertexArray() VertexArray

	// CullFace specifies which polygons are candidates for culling.
	//
	// Valid modes: FRONT, BACK, FRONT_AND_BACK.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCullFace.xhtml
	CullFace(mode Enum)

	// DeleteBuffer deletes the given buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
	DeleteBuffer(v Buffer)

	// DeleteFramebuffer deletes the given framebuffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteFramebuffers.xhtml
	DeleteFramebuffer(v Framebuffer)

	// DeleteProgram deletes the given program object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteProgram.xhtml
	DeleteProgram(p Program)

	// DeleteRenderbuffer deletes the given render buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteRenderbuffers.xhtml
	DeleteRenderbuffer(v Renderbuffer)

	// DeleteShader deletes shader s.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteShader.xhtml
	DeleteShader(s Shader)

	// DeleteTexture deletes the given texture object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
	DeleteTexture(v Texture)

	// DeleteVertexArray deletes the given render buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteVertexArrays.xhtml
	DeleteVertexArray(v VertexArray)

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
	DepthFunc(fn Enum)

	// DepthMask sets the depth buffer enabled for writing.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthMask.xhtml
	DepthMask(flag bool)

	// DepthRangef sets the mapping from normalized device coordinates to
	// window coordinates.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDepthRangef.xhtml
	DepthRangef(n, f float32)

	// DetachShader detaches the shader s from the program p.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDetachShader.xhtml
	DetachShader(p Program, s Shader)

	// Disable disables various GL capabilities.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
	Disable(cap Enum)

	// DisableVertexAttribArray disables a vertex attribute array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisableVertexAttribArray.xhtml
	DisableVertexAttribArray(a Attrib)

	// DrawArrays renders geometric primitives from the bound data.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
	DrawArrays(mode Enum, first, count int)

	// DrawElements renders primitives from a bound buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawElements.xhtml
	DrawElements(mode Enum, count int, ty Enum, offset int)

	// TODO(crawshaw): consider DrawElements8 / DrawElements16 / DrawElements32

	// Enable enables various GL capabilities.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
	Enable(cap Enum)

	// EnableVertexAttribArray enables a vertex attribute array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
	EnableVertexAttribArray(a Attrib)

	// Finish blocks until the effects of all previously called GL
	// commands are complete.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFinish.xhtml
	Finish()

	// Flush empties all buffers. It does not block.
	//
	// An OpenGL implementation may buffer network communication,
	// the command stream, or data inside the graphics accelerator.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
	Flush()

	// FramebufferRenderbuffer attaches rb to the current frame buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferRenderbuffer.xhtml
	FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer)

	// FramebufferTexture2D attaches the t to the current frame buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFramebufferTexture2D.xhtml
	FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int)

	// FrontFace defines which polygons are front-facing.
	//
	// Valid modes: CW, CCW.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFrontFace.xhtml
	FrontFace(mode Enum)

	// GenerateMipmap generates mipmaps for the current texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenerateMipmap.xhtml
	GenerateMipmap(target Enum)

	// GetActiveAttrib returns details about an active attribute variable.
	// A value of 0 for index selects the first active attribute variable.
	// Permissible values for index range from 0 to the number of active
	// attribute variables minus 1.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveAttrib.xhtml
	GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum)

	// GetActiveUniform returns details about an active uniform variable.
	// A value of 0 for index selects the first active uniform variable.
	// Permissible values for index range from 0 to the number of active
	// uniform variables minus 1.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetActiveUniform.xhtml
	GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum)

	// GetAttachedShaders returns the shader objects attached to program p.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttachedShaders.xhtml
	GetAttachedShaders(p Program) []Shader

	// GetAttribLocation returns the location of an attribute variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
	GetAttribLocation(p Program, name string) Attrib

	// GetBooleanv returns the boolean values of parameter pname.
	//
	// Many boolean parameters can be queried more easily using IsEnabled.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
	GetBooleanv(dst []bool, pname Enum)

	// GetFloatv returns the float values of parameter pname.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
	GetFloatv(dst []float32, pname Enum)

	// GetIntegerv returns the int values of parameter pname.
	//
	// Single values may be queried more easily using GetInteger.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
	GetIntegerv(dst []int32, pname Enum)

	// GetInteger returns the int value of parameter pname.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGet.xhtml
	GetInteger(pname Enum) int

	// GetBufferParameteri returns a parameter for the active buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetBufferParameter.xhtml
	GetBufferParameteri(target, value Enum) int

	// GetError returns the next error.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
	GetError() Enum

	// GetFramebufferAttachmentParameteri returns attachment parameters
	// for the active framebuffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetFramebufferAttachmentParameteriv.xhtml
	GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int

	// GetProgrami returns a parameter value for a program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
	GetProgrami(p Program, pname Enum) int

	// GetProgramInfoLog returns the information log for a program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
	GetProgramInfoLog(p Program) string

	// GetRenderbufferParameteri returns a parameter value for a render buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetRenderbufferParameteriv.xhtml
	GetRenderbufferParameteri(target, pname Enum) int

	// GetShaderi returns a parameter value for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
	GetShaderi(s Shader, pname Enum) int

	// GetShaderInfoLog returns the information log for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
	GetShaderInfoLog(s Shader) string

	// GetShaderPrecisionFormat returns range and precision limits for
	// shader types.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderPrecisionFormat.xhtml
	GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int)

	// GetShaderSource returns source code of shader s.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderSource.xhtml
	GetShaderSource(s Shader) string

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
	GetString(pname Enum) string

	// GetTexParameterfv returns the float values of a texture parameter.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
	GetTexParameterfv(dst []float32, target, pname Enum)

	// GetTexParameteriv returns the int values of a texture parameter.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetTexParameter.xhtml
	GetTexParameteriv(dst []int32, target, pname Enum)

	// GetUniformfv returns the float values of a uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
	GetUniformfv(dst []float32, src Uniform, p Program)

	// GetUniformiv returns the float values of a uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniform.xhtml
	GetUniformiv(dst []int32, src Uniform, p Program)

	// GetUniformLocation returns the location of a uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
	GetUniformLocation(p Program, name string) Uniform

	// GetVertexAttribf reads the float value of a vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
	GetVertexAttribf(src Attrib, pname Enum) float32

	// GetVertexAttribfv reads float values of a vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
	GetVertexAttribfv(dst []float32, src Attrib, pname Enum)

	// GetVertexAttribi reads the int value of a vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
	GetVertexAttribi(src Attrib, pname Enum) int32

	// GetVertexAttribiv reads int values of a vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetVertexAttrib.xhtml
	GetVertexAttribiv(dst []int32, src Attrib, pname Enum)

	// TODO(crawshaw): glGetVertexAttribPointerv

	// Hint sets implementation-specific modes.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glHint.xhtml
	Hint(target, mode Enum)

	// IsBuffer reports if b is a valid buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsBuffer.xhtml
	IsBuffer(b Buffer) bool

	// IsEnabled reports if cap is an enabled capability.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsEnabled.xhtml
	IsEnabled(cap Enum) bool

	// IsFramebuffer reports if fb is a valid frame buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsFramebuffer.xhtml
	IsFramebuffer(fb Framebuffer) bool

	// IsProgram reports if p is a valid program object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsProgram.xhtml
	IsProgram(p Program) bool

	// IsRenderbuffer reports if rb is a valid render buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsRenderbuffer.xhtml
	IsRenderbuffer(rb Renderbuffer) bool

	// IsShader reports if s is valid shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsShader.xhtml
	IsShader(s Shader) bool

	// IsTexture reports if t is a valid texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glIsTexture.xhtml
	IsTexture(t Texture) bool

	// LineWidth specifies the width of lines.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glLineWidth.xhtml
	LineWidth(width float32)

	// LinkProgram links the specified program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
	LinkProgram(p Program)

	// PixelStorei sets pixel storage parameters.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glPixelStorei.xhtml
	PixelStorei(pname Enum, param int32)

	// PolygonOffset sets the scaling factors for depth offsets.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glPolygonOffset.xhtml
	PolygonOffset(factor, units float32)

	// ReadPixels returns pixel data from a buffer.
	//
	// In GLES 3, the source buffer is controlled with ReadBuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
	ReadPixels(dst []byte, x, y, width, height int, format, ty Enum)

	// ReleaseShaderCompiler frees resources allocated by the shader compiler.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glReleaseShaderCompiler.xhtml
	ReleaseShaderCompiler()

	// RenderbufferStorage establishes the data storage, format, and
	// dimensions of a renderbuffer object's image.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glRenderbufferStorage.xhtml
	RenderbufferStorage(target, internalFormat Enum, width, height int)

	// SampleCoverage sets multisample coverage parameters.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glSampleCoverage.xhtml
	SampleCoverage(value float32, invert bool)

	// Scissor defines the scissor box rectangle, in window coordinates.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
	Scissor(x, y, width, height int32)

	// TODO(crawshaw): ShaderBinary

	// ShaderSource sets the source code of s to the given source code.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
	ShaderSource(s Shader, src string)

	// StencilFunc sets the front and back stencil test reference value.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFunc.xhtml
	StencilFunc(fn Enum, ref int, mask uint32)

	// StencilFunc sets the front or back stencil test reference value.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilFuncSeparate.xhtml
	StencilFuncSeparate(face, fn Enum, ref int, mask uint32)

	// StencilMask controls the writing of bits in the stencil planes.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMask.xhtml
	StencilMask(mask uint32)

	// StencilMaskSeparate controls the writing of bits in the stencil planes.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilMaskSeparate.xhtml
	StencilMaskSeparate(face Enum, mask uint32)

	// StencilOp sets front and back stencil test actions.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOp.xhtml
	StencilOp(fail, zfail, zpass Enum)

	// StencilOpSeparate sets front or back stencil tests.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glStencilOpSeparate.xhtml
	StencilOpSeparate(face, sfail, dpfail, dppass Enum)

	// TexImage2D writes a 2D texture image.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
	TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte)

	// TexSubImage2D writes a subregion of a 2D texture image.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexSubImage2D.xhtml
	TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte)

	// TexParameterf sets a float texture parameter.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
	TexParameterf(target, pname Enum, param float32)

	// TexParameterfv sets a float texture parameter array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
	TexParameterfv(target, pname Enum, params []float32)

	// TexParameteri sets an integer texture parameter.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
	TexParameteri(target, pname Enum, param int)

	// TexParameteriv sets an integer texture parameter array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
	TexParameteriv(target, pname Enum, params []int32)

	// Uniform1f writes a float uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform1f(dst Uniform, v float32)

	// Uniform1fv writes a [len(src)]float uniform array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform1fv(dst Uniform, src []float32)

	// Uniform1i writes an int uniform variable.
	//
	// Uniform1i and Uniform1iv are the only two functions that may be used
	// to load uniform variables defined as sampler types. Loading samplers
	// with any other function will result in a INVALID_OPERATION error.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform1i(dst Uniform, v int)

	// Uniform1iv writes a int uniform array of len(src) elements.
	//
	// Uniform1i and Uniform1iv are the only two functions that may be used
	// to load uniform variables defined as sampler types. Loading samplers
	// with any other function will result in a INVALID_OPERATION error.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform1iv(dst Uniform, src []int32)

	// Uniform2f writes a vec2 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform2f(dst Uniform, v0, v1 float32)

	// Uniform2fv writes a vec2 uniform array of len(src)/2 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform2fv(dst Uniform, src []float32)

	// Uniform2i writes an ivec2 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform2i(dst Uniform, v0, v1 int)

	// Uniform2iv writes an ivec2 uniform array of len(src)/2 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform2iv(dst Uniform, src []int32)

	// Uniform3f writes a vec3 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform3f(dst Uniform, v0, v1, v2 float32)

	// Uniform3fv writes a vec3 uniform array of len(src)/3 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform3fv(dst Uniform, src []float32)

	// Uniform3i writes an ivec3 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform3i(dst Uniform, v0, v1, v2 int32)

	// Uniform3iv writes an ivec3 uniform array of len(src)/3 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform3iv(dst Uniform, src []int32)

	// Uniform4f writes a vec4 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4f(dst Uniform, v0, v1, v2, v3 float32)

	// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4fv(dst Uniform, src []float32)

	// Uniform4i writes an ivec4 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4i(dst Uniform, v0, v1, v2, v3 int32)

	// Uniform4i writes an ivec4 uniform array of len(src)/4 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4iv(dst Uniform, src []int32)

	// UniformMatrix2fv writes 2x2 matrices. Each matrix uses four
	// float32 values, so the number of matrices written is len(src)/4.
	//
	// Each matrix must be supplied in column major order.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	UniformMatrix2fv(dst Uniform, src []float32)

	// UniformMatrix3fv writes 3x3 matrices. Each matrix uses nine
	// float32 values, so the number of matrices written is len(src)/9.
	//
	// Each matrix must be supplied in column major order.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	UniformMatrix3fv(dst Uniform, src []float32)

	// UniformMatrix4fv writes 4x4 matrices. Each matrix uses 16
	// float32 values, so the number of matrices written is len(src)/16.
	//
	// Each matrix must be supplied in column major order.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	UniformMatrix4fv(dst Uniform, src []float32)

	// UseProgram sets the active program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
	UseProgram(p Program)

	// ValidateProgram checks to see whether the executables contained in
	// program can execute given the current OpenGL state.
	//
	// Typically only used for debugging.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glValidateProgram.xhtml
	ValidateProgram(p Program)

	// VertexAttrib1f writes a float vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib1f(dst Attrib, x float32)

	// VertexAttrib1fv writes a float vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib1fv(dst Attrib, src []float32)

	// VertexAttrib2f writes a vec2 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib2f(dst Attrib, x, y float32)

	// VertexAttrib2fv writes a vec2 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib2fv(dst Attrib, src []float32)

	// VertexAttrib3f writes a vec3 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib3f(dst Attrib, x, y, z float32)

	// VertexAttrib3fv writes a vec3 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib3fv(dst Attrib, src []float32)

	// VertexAttrib4f writes a vec4 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib4f(dst Attrib, x, y, z, w float32)

	// VertexAttrib4fv writes a vec4 vertex attribute.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttrib.xhtml
	VertexAttrib4fv(dst Attrib, src []float32)

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
	VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int)

	// Viewport sets the viewport, an affine transformation that
	// normalizes device coordinates to window coordinates.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
	Viewport(x, y, width, height int)
}

// Context3 is an OpenGL ES 3 context.
//
// When the gl package is compiled with GL ES 3 support, the produced
// Context object also implements the Context3 interface.
type Context3 interface {
	Context

	// BlitFramebuffer copies a block of pixels between framebuffers.
	//
	// https://www.khronos.org/opengles/sdk/docs/man3/html/glBlitFramebuffer.xhtml
	BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int, mask uint, filter Enum)
}

// Worker is used by display driver code to execute OpenGL calls.
//
// Typically display driver code creates a gl.Context for an application,
// and along with it establishes a locked OS thread to execute the cgo
// calls:
//
//	go func() {
//		runtime.LockOSThread()
//		// ... platform-specific cgo call to bind a C OpenGL context
//		// into thread-local storage.
//
//		glctx, worker := gl.NewContext()
//		workAvailable := worker.WorkAvailable()
//		go userAppCode(glctx)
//		for {
//			select {
//			case <-workAvailable:
//				worker.DoWork()
//			case <-drawEvent:
//				// ... platform-specific cgo call to draw screen
//			}
//		}
//	}()
//
// This interface is an internal implementation detail and should only be used
// by the package responsible for managing the screen, such as
// golang.org/x/mobile/app.
type Worker interface {
	// WorkAvailable returns a channel that communicates when DoWork should be
	// called.
	WorkAvailable() <-chan struct{}

	// DoWork performs any pending OpenGL calls.
	DoWork()
}
