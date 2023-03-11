// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ios android

package gl

/*
#include <stdlib.h>

#ifdef os_ios
#include <OpenGLES/ES2/glext.h>
#endif
#ifdef os_android
#include <GLES2/gl2.h>
#endif
*/
import "C"

import "unsafe"

var ContextWatcher contextWatcher

type contextWatcher struct{}

func (contextWatcher) OnMakeCurrent(context interface{}) {}
func (contextWatcher) OnDetach()                         {}

func ActiveTexture(texture Enum) {
	C.glActiveTexture(texture.c())
}

func AttachShader(p Program, s Shader) {
	C.glAttachShader(p.c(), s.c())
}

func BindAttribLocation(p Program, a Attrib, name string) {
	str := unsafe.Pointer(C.CString(name))
	defer C.free(str)
	C.glBindAttribLocation(p.c(), a.c(), (*C.GLchar)(str))
}

func BindBuffer(target Enum, b Buffer) {
	C.glBindBuffer(target.c(), b.c())
}

func BindFramebuffer(target Enum, fb Framebuffer) {
	C.glBindFramebuffer(target.c(), fb.c())
}

func BindRenderbuffer(target Enum, rb Renderbuffer) {
	C.glBindRenderbuffer(target.c(), rb.c())
}

func BindTexture(target Enum, t Texture) {
	C.glBindTexture(target.c(), t.c())
}

func BlendColor(red, green, blue, alpha float32) {
	blendColor(red, green, blue, alpha)
}

func BlendEquation(mode Enum) {
	C.glBlendEquation(mode.c())
}

func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	C.glBlendEquationSeparate(modeRGB.c(), modeAlpha.c())
}

func BlendFunc(sfactor, dfactor Enum) {
	C.glBlendFunc(sfactor.c(), dfactor.c())
}

func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	C.glBlendFuncSeparate(sfactorRGB.c(), dfactorRGB.c(), sfactorAlpha.c(), dfactorAlpha.c())
}

func BufferData(target Enum, src []byte, usage Enum) {
	C.glBufferData(target.c(), C.GLsizeiptr(len(src)), unsafe.Pointer(&src[0]), usage.c())
}

func BufferInit(target Enum, size int, usage Enum) {
	C.glBufferData(target.c(), C.GLsizeiptr(size), nil, usage.c())
}

func BufferSubData(target Enum, offset int, data []byte) {
	C.glBufferSubData(target.c(), C.GLintptr(offset), C.GLsizeiptr(len(data)), unsafe.Pointer(&data[0]))
}

func CheckFramebufferStatus(target Enum) Enum {
	return Enum(C.glCheckFramebufferStatus(target.c()))
}

func Clear(mask Enum) {
	C.glClear(C.GLbitfield(mask))
}

func ClearColor(red, green, blue, alpha float32) {
	clearColor(red, green, blue, alpha)
}

func ClearDepthf(d float32) {
	clearDepthf(d)
}

func ClearStencil(s int) {
	C.glClearStencil(C.GLint(s))
}

func ColorMask(red, green, blue, alpha bool) {
	C.glColorMask(glBoolean(red), glBoolean(green), glBoolean(blue), glBoolean(alpha))
}

func CompileShader(s Shader) {
	C.glCompileShader(s.c())
}

func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	C.glCompressedTexImage2D(target.c(), C.GLint(level), internalformat.c(), C.GLsizei(width), C.GLsizei(height), C.GLint(border), C.GLsizei(len(data)), unsafe.Pointer(&data[0]))
}

func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	C.glCompressedTexSubImage2D(target.c(), C.GLint(level), C.GLint(xoffset), C.GLint(yoffset), C.GLsizei(width), C.GLsizei(height), format.c(), C.GLsizei(len(data)), unsafe.Pointer(&data[0]))
}

func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	C.glCopyTexImage2D(target.c(), C.GLint(level), internalformat.c(), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), C.GLint(border))
}

func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	C.glCopyTexSubImage2D(target.c(), C.GLint(level), C.GLint(xoffset), C.GLint(yoffset), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func CreateBuffer() Buffer {
	var b Buffer
	C.glGenBuffers(1, (*C.GLuint)(&b.Value))
	return b
}

func CreateFramebuffer() Framebuffer {
	var b Framebuffer
	C.glGenFramebuffers(1, (*C.GLuint)(&b.Value))
	return b
}

func CreateProgram() Program {
	return Program{Value: uint32(C.glCreateProgram())}
}

func CreateRenderbuffer() Renderbuffer {
	var b Renderbuffer
	C.glGenRenderbuffers(1, (*C.GLuint)(&b.Value))
	return b
}

func CreateShader(ty Enum) Shader {
	return Shader{Value: uint32(C.glCreateShader(ty.c()))}
}

func CreateTexture() Texture {
	var t Texture
	C.glGenTextures(1, (*C.GLuint)(&t.Value))
	return t
}

func CullFace(mode Enum) {
	C.glCullFace(mode.c())
}

func DeleteBuffer(v Buffer) {
	C.glDeleteBuffers(1, (*C.GLuint)(&v.Value))
}

func DeleteFramebuffer(v Framebuffer) {
	C.glDeleteFramebuffers(1, (*C.GLuint)(&v.Value))
}

func DeleteProgram(p Program) {
	C.glDeleteProgram(p.c())
}

func DeleteRenderbuffer(v Renderbuffer) {
	C.glDeleteRenderbuffers(1, (*C.GLuint)(&v.Value))
}

func DeleteShader(s Shader) {
	C.glDeleteShader(s.c())
}

func DeleteTexture(v Texture) {
	C.glDeleteTextures(1, (*C.GLuint)(&v.Value))
}

func DepthFunc(fn Enum) {
	C.glDepthFunc(fn.c())
}

func DepthMask(flag bool) {
	C.glDepthMask(glBoolean(flag))
}

func DepthRangef(n, f float32) {
	depthRangef(n, f)
}

func DetachShader(p Program, s Shader) {
	C.glDetachShader(p.c(), s.c())
}

func Disable(cap Enum) {
	C.glDisable(cap.c())
}

func DisableVertexAttribArray(a Attrib) {
	C.glDisableVertexAttribArray(a.c())
}

func DrawArrays(mode Enum, first, count int) {
	C.glDrawArrays(mode.c(), C.GLint(first), C.GLsizei(count))
}

func DrawElements(mode Enum, count int, ty Enum, offset int) {
	C.glDrawElements(mode.c(), C.GLsizei(count), ty.c(), unsafe.Pointer(uintptr(offset)))
}

func Enable(cap Enum) {
	C.glEnable(cap.c())
}

func EnableVertexAttribArray(a Attrib) {
	C.glEnableVertexAttribArray(a.c())
}

func Finish() {
	C.glFinish()
}

func Flush() {
	C.glFlush()
}

func FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	C.glFramebufferRenderbuffer(target.c(), attachment.c(), rbTarget.c(), rb.c())
}

func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	C.glFramebufferTexture2D(target.c(), attachment.c(), texTarget.c(), t.c(), C.GLint(level))
}

func FrontFace(mode Enum) {
	C.glFrontFace(mode.c())
}

func GenerateMipmap(target Enum) {
	C.glGenerateMipmap(target.c())
}

func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	bufSize := GetProgrami(p, ACTIVE_ATTRIBUTE_MAX_LENGTH)
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	var cSize C.GLint
	var cType C.GLenum
	C.glGetActiveAttrib(p.c(), C.GLuint(index), C.GLsizei(bufSize), nil, &cSize, &cType, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf)), int(cSize), Enum(cType)
}

func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	bufSize := GetProgrami(p, ACTIVE_UNIFORM_MAX_LENGTH)
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)

	var cSize C.GLint
	var cType C.GLenum

	C.glGetActiveUniform(p.c(), C.GLuint(index), C.GLsizei(bufSize), nil, &cSize, &cType, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf)), int(cSize), Enum(cType)
}

func GetAttachedShaders(p Program) []Shader {
	shadersLen := GetProgrami(p, ATTACHED_SHADERS)
	var n C.GLsizei
	buf := make([]C.GLuint, shadersLen)
	C.glGetAttachedShaders(p.c(), C.GLsizei(shadersLen), &n, &buf[0])
	buf = buf[:int(n)]
	shaders := make([]Shader, len(buf))
	for i, s := range buf {
		shaders[i] = Shader{Value: uint32(s)}
	}
	return shaders
}

func GetAttribLocation(p Program, name string) Attrib {
	str := unsafe.Pointer(C.CString(name))
	defer C.free(str)
	return Attrib{Value: uint(C.glGetAttribLocation(p.c(), (*C.GLchar)(str)))}
}

func GetBooleanv(dst []bool, pname Enum) {
	buf := make([]C.GLboolean, len(dst))
	C.glGetBooleanv(pname.c(), &buf[0])
	for i, v := range buf {
		dst[i] = v != 0
	}
}

func GetFloatv(dst []float32, pname Enum) {
	C.glGetFloatv(pname.c(), (*C.GLfloat)(&dst[0]))
}

func GetIntegerv(pname Enum, data []int32) {
	buf := make([]C.GLint, len(data))
	C.glGetIntegerv(pname.c(), &buf[0])
	for i, v := range buf {
		data[i] = int32(v)
	}
}

func GetInteger(pname Enum) int {
	var v C.GLint
	C.glGetIntegerv(pname.c(), &v)
	return int(v)
}

func GetBufferParameteri(target, pname Enum) int {
	var params C.GLint
	C.glGetBufferParameteriv(target.c(), pname.c(), &params)
	return int(params)
}

func GetError() Enum {
	return Enum(C.glGetError())
}

func GetBoundFramebuffer() Framebuffer {
	println("GetBoundFramebuffer: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	var b C.GLint
	C.glGetIntegerv(FRAMEBUFFER_BINDING, &b)
	return Framebuffer{Value: uint32(b)}
}

func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	var params C.GLint
	C.glGetFramebufferAttachmentParameteriv(target.c(), attachment.c(), pname.c(), &params)
	return int(params)
}

func GetProgrami(p Program, pname Enum) int {
	var params C.GLint
	C.glGetProgramiv(p.c(), pname.c(), &params)
	return int(params)
}

func GetProgramInfoLog(p Program) string {
	infoLen := GetProgrami(p, INFO_LOG_LENGTH)
	buf := C.malloc(C.size_t(infoLen))
	C.free(buf)
	C.glGetProgramInfoLog(p.c(), C.GLsizei(infoLen), nil, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf))
}

func GetRenderbufferParameteri(target, pname Enum) int {
	var params C.GLint
	C.glGetRenderbufferParameteriv(target.c(), pname.c(), &params)
	return int(params)
}

func GetShaderi(s Shader, pname Enum) int {
	var params C.GLint
	C.glGetShaderiv(s.c(), pname.c(), &params)
	return int(params)
}

func GetShaderInfoLog(s Shader) string {
	infoLen := GetShaderi(s, INFO_LOG_LENGTH)
	buf := C.malloc(C.size_t(infoLen))
	defer C.free(buf)
	C.glGetShaderInfoLog(s.c(), C.GLsizei(infoLen), nil, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf))
}

func GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	const glintSize = 4
	var cRange [2]C.GLint
	var cPrecision C.GLint

	C.glGetShaderPrecisionFormat(shadertype.c(), precisiontype.c(), &cRange[0], &cPrecision)
	return int(cRange[0]), int(cRange[1]), int(cPrecision)
}

func GetShaderSource(s Shader) string {
	sourceLen := GetShaderi(s, SHADER_SOURCE_LENGTH)
	if sourceLen == 0 {
		return ""
	}
	buf := C.malloc(C.size_t(sourceLen))
	defer C.free(buf)
	C.glGetShaderSource(s.c(), C.GLsizei(sourceLen), nil, (*C.GLchar)(buf))
	return C.GoString((*C.char)(buf))
}

func GetString(pname Enum) string {
	// Bounce through unsafe.Pointer, because on some platforms
	// GetString returns an *unsigned char which doesn't convert.
	return C.GoString((*C.char)((unsafe.Pointer)(C.glGetString(pname.c()))))
}

func GetTexParameterfv(dst []float32, target, pname Enum) {
	C.glGetTexParameterfv(target.c(), pname.c(), (*C.GLfloat)(&dst[0]))
}

func GetTexParameteriv(dst []int32, target, pname Enum) {
	C.glGetTexParameteriv(target.c(), pname.c(), (*C.GLint)(&dst[0]))
}

func GetUniformfv(dst []float32, src Uniform, p Program) {
	C.glGetUniformfv(p.c(), src.c(), (*C.GLfloat)(&dst[0]))
}

func GetUniformiv(dst []int32, src Uniform, p Program) {
	C.glGetUniformiv(p.c(), src.c(), (*C.GLint)(&dst[0]))
}

func GetUniformLocation(p Program, name string) Uniform {
	str := unsafe.Pointer(C.CString(name))
	defer C.free(str)
	return Uniform{Value: int32(C.glGetUniformLocation(p.c(), (*C.GLchar)(str)))}
}

func GetVertexAttribf(src Attrib, pname Enum) float32 {
	var params C.GLfloat
	C.glGetVertexAttribfv(src.c(), pname.c(), &params)
	return float32(params)
}

func GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	C.glGetVertexAttribfv(src.c(), pname.c(), (*C.GLfloat)(&dst[0]))
}

func GetVertexAttribi(src Attrib, pname Enum) int32 {
	var params C.GLint
	C.glGetVertexAttribiv(src.c(), pname.c(), &params)
	return int32(params)
}

func GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	C.glGetVertexAttribiv(src.c(), pname.c(), (*C.GLint)(&dst[0]))
}

func Hint(target, mode Enum) {
	C.glHint(target.c(), mode.c())
}

func IsBuffer(b Buffer) bool {
	return C.glIsBuffer(b.c()) != 0
}

func IsEnabled(cap Enum) bool {
	return C.glIsEnabled(cap.c()) != 0
}

func IsFramebuffer(fb Framebuffer) bool {
	return C.glIsFramebuffer(fb.c()) != 0
}

func IsProgram(p Program) bool {
	return C.glIsProgram(p.c()) != 0
}

func IsRenderbuffer(rb Renderbuffer) bool {
	return C.glIsRenderbuffer(rb.c()) != 0
}

func IsShader(s Shader) bool {
	return C.glIsShader(s.c()) != 0
}

func IsTexture(t Texture) bool {
	return C.glIsTexture(t.c()) != 0
}

func LineWidth(width float32) {
	C.glLineWidth(C.GLfloat(width))
}

func LinkProgram(p Program) {
	C.glLinkProgram(p.c())
}

func PixelStorei(pname Enum, param int32) {
	C.glPixelStorei(pname.c(), C.GLint(param))
}

func PolygonOffset(factor, units float32) {
	C.glPolygonOffset(C.GLfloat(factor), C.GLfloat(units))
}

func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	C.glReadPixels(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), format.c(), ty.c(), unsafe.Pointer(&dst[0]))
}

func ReleaseShaderCompiler() {
	C.glReleaseShaderCompiler()
}

func RenderbufferStorage(target, internalFormat Enum, width, height int) {
	C.glRenderbufferStorage(target.c(), internalFormat.c(), C.GLsizei(width), C.GLsizei(height))
}

func SampleCoverage(value float32, invert bool) {
	sampleCoverage(value, invert)
}

func Scissor(x, y, width, height int32) {
	C.glScissor(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func ShaderSource(s Shader, src string) {
	str := (*C.GLchar)(C.CString(src))
	defer C.free(unsafe.Pointer(str))
	C.glShaderSource(s.c(), 1, &str, nil)
}

func StencilFunc(fn Enum, ref int, mask uint32) {
	C.glStencilFunc(fn.c(), C.GLint(ref), C.GLuint(mask))
}

func StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	C.glStencilFuncSeparate(face.c(), fn.c(), C.GLint(ref), C.GLuint(mask))
}

func StencilMask(mask uint32) {
	C.glStencilMask(C.GLuint(mask))
}

func StencilMaskSeparate(face Enum, mask uint32) {
	C.glStencilMaskSeparate(face.c(), C.GLuint(mask))
}

func StencilOp(fail, zfail, zpass Enum) {
	C.glStencilOp(fail.c(), zfail.c(), zpass.c())
}

func StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	C.glStencilOpSeparate(face.c(), sfail.c(), dpfail.c(), dppass.c())
}

func TexImage2D(target Enum, level int, width, height int, format Enum, ty Enum, data []byte) {
	p := unsafe.Pointer(nil)
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.glTexImage2D(target.c(), C.GLint(level), C.GLint(format), C.GLsizei(width), C.GLsizei(height), 0, format.c(), ty.c(), p)
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	C.glTexSubImage2D(target.c(), C.GLint(level), C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height), format.c(), ty.c(), unsafe.Pointer(&data[0]))
}

func TexParameterf(target, pname Enum, param float32) {
	C.glTexParameterf(target.c(), pname.c(), C.GLfloat(param))
}

func TexParameterfv(target, pname Enum, params []float32) {
	C.glTexParameterfv(target.c(), pname.c(), (*C.GLfloat)(&params[0]))
}

func TexParameteri(target, pname Enum, param int) {
	C.glTexParameteri(target.c(), pname.c(), C.GLint(param))
}

func TexParameteriv(target, pname Enum, params []int32) {
	C.glTexParameteriv(target.c(), pname.c(), (*C.GLint)(&params[0]))
}

func Uniform1f(dst Uniform, v float32) {
	C.glUniform1f(dst.c(), C.GLfloat(v))
}

func Uniform1fv(dst Uniform, src []float32) {
	C.glUniform1fv(dst.c(), C.GLsizei(len(src)), (*C.GLfloat)(&src[0]))
}

func Uniform1i(dst Uniform, v int) {
	C.glUniform1i(dst.c(), C.GLint(v))
}

func Uniform1iv(dst Uniform, src []int32) {
	C.glUniform1iv(dst.c(), C.GLsizei(len(src)), (*C.GLint)(&src[0]))
}

func Uniform2f(dst Uniform, v0, v1 float32) {
	C.glUniform2f(dst.c(), C.GLfloat(v0), C.GLfloat(v1))
}

func Uniform2fv(dst Uniform, src []float32) {
	C.glUniform2fv(dst.c(), C.GLsizei(len(src)/2), (*C.GLfloat)(&src[0]))
}

func Uniform2i(dst Uniform, v0, v1 int) {
	C.glUniform2i(dst.c(), C.GLint(v0), C.GLint(v1))
}

func Uniform2iv(dst Uniform, src []int32) {
	C.glUniform2iv(dst.c(), C.GLsizei(len(src)/2), (*C.GLint)(&src[0]))
}

func Uniform3f(dst Uniform, v0, v1, v2 float32) {
	C.glUniform3f(dst.c(), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2))
}

func Uniform3fv(dst Uniform, src []float32) {
	C.glUniform3fv(dst.c(), C.GLsizei(len(src)/3), (*C.GLfloat)(&src[0]))
}

func Uniform3i(dst Uniform, v0, v1, v2 int32) {
	C.glUniform3i(dst.c(), C.GLint(v0), C.GLint(v1), C.GLint(v2))
}

func Uniform3iv(dst Uniform, src []int32) {
	C.glUniform3iv(dst.c(), C.GLsizei(len(src)/3), (*C.GLint)(&src[0]))
}

func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	C.glUniform4f(dst.c(), C.GLfloat(v0), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
}

func Uniform4fv(dst Uniform, src []float32) {
	C.glUniform4fv(dst.c(), C.GLsizei(len(src)/4), (*C.GLfloat)(&src[0]))
}

func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	C.glUniform4i(dst.c(), C.GLint(v0), C.GLint(v1), C.GLint(v2), C.GLint(v3))
}

func Uniform4iv(dst Uniform, src []int32) {
	C.glUniform4iv(dst.c(), C.GLsizei(len(src)/4), (*C.GLint)(&src[0]))
}

func UniformMatrix2fv(dst Uniform, src []float32) {
	// OpenGL ES 2 does not support transpose.
	C.glUniformMatrix2fv(dst.c(), C.GLsizei(len(src)/4), 0, (*C.GLfloat)(&src[0]))
}

func UniformMatrix3fv(dst Uniform, src []float32) {
	C.glUniformMatrix3fv(dst.c(), C.GLsizei(len(src)/9), 0, (*C.GLfloat)(&src[0]))
}

func UniformMatrix4fv(dst Uniform, src []float32) {
	C.glUniformMatrix4fv(dst.c(), C.GLsizei(len(src)/16), 0, (*C.GLfloat)(&src[0]))
}

func UseProgram(p Program) {
	C.glUseProgram(p.c())
}

func ValidateProgram(p Program) {
	C.glValidateProgram(p.c())
}

func VertexAttrib1f(dst Attrib, x float32) {
	C.glVertexAttrib1f(dst.c(), C.GLfloat(x))
}

func VertexAttrib1fv(dst Attrib, src []float32) {
	C.glVertexAttrib1fv(dst.c(), (*C.GLfloat)(&src[0]))
}

func VertexAttrib2f(dst Attrib, x, y float32) {
	C.glVertexAttrib2f(dst.c(), C.GLfloat(x), C.GLfloat(y))
}

func VertexAttrib2fv(dst Attrib, src []float32) {
	C.glVertexAttrib2fv(dst.c(), (*C.GLfloat)(&src[0]))
}

func VertexAttrib3f(dst Attrib, x, y, z float32) {
	C.glVertexAttrib3f(dst.c(), C.GLfloat(x), C.GLfloat(y), C.GLfloat(z))
}

func VertexAttrib3fv(dst Attrib, src []float32) {
	C.glVertexAttrib3fv(dst.c(), (*C.GLfloat)(&src[0]))
}

func VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	C.glVertexAttrib4f(dst.c(), C.GLfloat(x), C.GLfloat(y), C.GLfloat(z), C.GLfloat(w))
}

func VertexAttrib4fv(dst Attrib, src []float32) {
	C.glVertexAttrib4fv(dst.c(), (*C.GLfloat)(&src[0]))
}

func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	n := glBoolean(normalized)
	s := C.GLsizei(stride)
	C.glVertexAttribPointer(dst.c(), C.GLint(size), ty.c(), n, s, unsafe.Pointer(uintptr(offset)))
}

func Viewport(x, y, width, height int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}
