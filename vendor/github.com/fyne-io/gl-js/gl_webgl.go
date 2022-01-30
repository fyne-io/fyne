// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js,!wasm

package gl

import (
	"encoding/binary"
	"math"

	"github.com/gopherjs/gopherjs/js"
)

var ContextWatcher contextWatcher

type contextWatcher struct{}

func (contextWatcher) OnMakeCurrent(context interface{}) {
	// context must be a WebGLRenderingContext *js.Object.
	c = context.(*js.Object)
}
func (contextWatcher) OnDetach() {
	c = nil
}

// c is the current WebGL context, or nil if there is no current context.
var c *js.Object

func ActiveTexture(texture Enum) {
	c.Call("activeTexture", texture)
}

func AttachShader(p Program, s Shader) {
	c.Call("attachShader", p.Object, s.Object)
}

func BindAttribLocation(p Program, a Attrib, name string) {
	c.Call("bindAttribLocation", p.Object, a.Value, name)
}

func BindBuffer(target Enum, b Buffer) {
	c.Call("bindBuffer", target, b.Object)
}

func BindFramebuffer(target Enum, fb Framebuffer) {
	c.Call("bindFramebuffer", target, fb.Object)
}

func BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.Call("bindRenderbuffer", target, rb.Object)
}

func BindTexture(target Enum, t Texture) {
	c.Call("bindTexture", target, t.Object)
}

func BlendColor(red, green, blue, alpha float32) {
	c.Call("blendColor", red, green, blue, alpha)
}

func BlendEquation(mode Enum) {
	c.Call("blendEquation", mode)
}

func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func BlendFunc(sfactor, dfactor Enum) {
	c.Call("blendFunc", sfactor, dfactor)
}

func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	c.Call("blendFuncSeparate", sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha)
}

func BufferData(target Enum, data interface{}, usage Enum) {
	c.Call("bufferData", target, data, usage)
}

func BufferInit(target Enum, size int, usage Enum) {
	c.Call("bufferData", target, size, usage)
}

func BufferSubData(target Enum, offset int, data []byte) {
	c.Call("bufferSubData", target, offset, data)
}

func CheckFramebufferStatus(target Enum) Enum {
	return Enum(c.Call("checkFramebufferStatus", target).Int())
}

func Clear(mask Enum) {
	c.Call("clear", mask)
}

func ClearColor(red, green, blue, alpha float32) {
	c.Call("clearColor", red, green, blue, alpha)
}

func ClearDepthf(d float32) {
	c.Call("clearDepth", d)
}

func ClearStencil(s int) {
	c.Call("clearStencil", s)
}

func ColorMask(red, green, blue, alpha bool) {
	c.Call("colorMask", red, green, blue, alpha)
}

func CompileShader(s Shader) {
	c.Call("compileShader", s.Object)
}

func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	c.Call("compressedTexImage2D", target, level, internalformat, width, height, border, data)
}

func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	c.Call("compressedTexSubImage2D", target, level, xoffset, yoffset, width, height, format, data)
}

func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	c.Call("copyTexImage2D", target, level, internalformat, x, y, width, height, border)
}

func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	c.Call("copyTexSubImage2D", target, level, xoffset, yoffset, x, y, width, height)
}

func CreateBuffer() Buffer {
	return Buffer{Object: c.Call("createBuffer")}
}

func CreateFramebuffer() Framebuffer {
	return Framebuffer{Object: c.Call("createFramebuffer")}
}

func CreateProgram() Program {
	return Program{Object: c.Call("createProgram")}
}

func CreateRenderbuffer() Renderbuffer {
	return Renderbuffer{Object: c.Call("createRenderbuffer")}
}

func CreateShader(ty Enum) Shader {
	return Shader{Object: c.Call("createShader", ty)}
}

func CreateTexture() Texture {
	return Texture{Object: c.Call("createTexture")}
}

func CullFace(mode Enum) {
	c.Call("cullFace", mode)
}

func DeleteBuffer(v Buffer) {
	c.Call("deleteBuffer", v.Object)
}

func DeleteFramebuffer(v Framebuffer) {
	c.Call("deleteFramebuffer", v.Object)
}

func DeleteProgram(p Program) {
	c.Call("deleteProgram", p.Object)
}

func DeleteRenderbuffer(v Renderbuffer) {
	c.Call("deleteRenderbuffer", v.Object)
}

func DeleteShader(s Shader) {
	c.Call("deleteShader", s.Object)
}

func DeleteTexture(v Texture) {
	c.Call("deleteTexture", v.Object)
}

func DepthFunc(fn Enum) {
	c.Call("depthFunc", fn)
}

func DepthMask(flag bool) {
	c.Call("depthMask", flag)
}

func DepthRangef(n, f float32) {
	c.Call("depthRange", n, f)
}

func DetachShader(p Program, s Shader) {
	c.Call("detachShader", p.Object, s.Object)
}

func Disable(cap Enum) {
	c.Call("disable", cap)
}

func DisableVertexAttribArray(a Attrib) {
	c.Call("disableVertexAttribArray", a.Value)
}

func DrawArrays(mode Enum, first, count int) {
	c.Call("drawArrays", mode, first, count)
}

func DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.Call("drawElements", mode, count, ty, offset)
}

func Enable(cap Enum) {
	c.Call("enable", cap)
}

func EnableVertexAttribArray(a Attrib) {
	c.Call("enableVertexAttribArray", a.Value)
}

func Finish() {
	c.Call("finish")
}

func Flush() {
	c.Call("flush")
}

func FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	c.Call("framebufferRenderbuffer", target, attachment, rbTarget, rb.Object)
}

func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	c.Call("framebufferTexture2D", target, attachment, texTarget, t.Object, level)
}

func FrontFace(mode Enum) {
	c.Call("frontFace", mode)
}

func GenerateMipmap(target Enum) {
	c.Call("generateMipmap", target)
}

type activeInfo struct {
	*js.Object
	Size int    `js:"size"`
	Type int    `js:"type"`
	Name string `js:"name"`
}

func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	ai := activeInfo{Object: c.Call("getActiveAttrib", p.Object, index)}
	return ai.Name, ai.Size, Enum(ai.Type)
}

func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	ai := activeInfo{Object: c.Call("getActiveUniform", p.Object, index)}
	return ai.Name, ai.Size, Enum(ai.Type)
}

func GetAttachedShaders(p Program) []Shader {
	objs := c.Call("getAttachedShaders", p.Object)
	shaders := make([]Shader, objs.Length())
	for i := 0; i < objs.Length(); i++ {
		shaders[i] = Shader{Object: objs.Index(i)}
	}
	return shaders
}

func GetAttribLocation(p Program, name string) Attrib {
	return Attrib{Value: c.Call("getAttribLocation", p.Object, name).Int()}
}

func GetBooleanv(dst []bool, pname Enum) {
	println("GetBooleanv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getParameter", pname)
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = result.Index(i).Bool()
	}
}

func GetFloatv(dst []float32, pname Enum) {
	println("GetFloatv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getParameter", pname)
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = float32(result.Index(i).Float())
	}
}

func GetIntegerv(pname Enum, data []int32) {
	result := c.Call("getParameter", pname)
	length := result.Length()
	for i := 0; i < length; i++ {
		data[i] = int32(result.Index(i).Int())
	}
}

func GetInteger(pname Enum) int {
	return c.Call("getParameter", pname).Int()
}

func GetBufferParameteri(target, pname Enum) int {
	return c.Call("getBufferParameter", target, pname).Int()
}

func GetError() Enum {
	return Enum(c.Call("getError").Int())
}

func GetBoundFramebuffer() Framebuffer {
	return Framebuffer{Object: c.Call("getParameter", FRAMEBUFFER_BINDING)}
}

func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return c.Call("getFramebufferAttachmentParameter", target, attachment, pname).Int()
}

func GetProgrami(p Program, pname Enum) int {
	switch pname {
	case DELETE_STATUS, LINK_STATUS, VALIDATE_STATUS:
		if c.Call("getProgramParameter", p.Object, pname).Bool() {
			return TRUE
		}
		return FALSE
	default:
		return c.Call("getProgramParameter", p.Object, pname).Int()
	}
}

func GetProgramInfoLog(p Program) string {
	return c.Call("getProgramInfoLog", p.Object).String()
}

func GetRenderbufferParameteri(target, pname Enum) int {
	return c.Call("getRenderbufferParameter", target, pname).Int()
}

func GetShaderi(s Shader, pname Enum) int {
	switch pname {
	case DELETE_STATUS, COMPILE_STATUS:
		if c.Call("getShaderParameter", s.Object, pname).Bool() {
			return TRUE
		}
		return FALSE
	default:
		return c.Call("getShaderParameter", s.Object, pname).Int()
	}
}

func GetShaderInfoLog(s Shader) string {
	return c.Call("getShaderInfoLog", s.Object).String()
}

func GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeMin, rangeMax, precision int) {
	println("GetShaderPrecisionFormat: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	format := c.Call("getShaderPrecisionFormat", shadertype, precisiontype)
	rangeMin = format.Get("rangeMin").Int()
	rangeMax = format.Get("rangeMax").Int()
	precision = format.Get("precision").Int()
	return
}

func GetShaderSource(s Shader) string {
	return c.Call("getShaderSource", s.Object).String()
}

func GetString(pname Enum) string {
	return c.Call("getParameter", pname).String()
}

func GetTexParameterfv(dst []float32, target, pname Enum) {
	dst[0] = float32(c.Call("getTexParameter", pname).Float())
}

func GetTexParameteriv(dst []int32, target, pname Enum) {
	dst[0] = int32(c.Call("getTexParameter", pname).Int())
}

func GetUniformfv(dst []float32, src Uniform, p Program) {
	println("GetUniformfv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getUniform")
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = float32(result.Index(i).Float())
	}
}

func GetUniformiv(dst []int32, src Uniform, p Program) {
	println("GetUniformiv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getUniform")
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = int32(result.Index(i).Int())
	}
}

func GetUniformLocation(p Program, name string) Uniform {
	return Uniform{Object: c.Call("getUniformLocation", p.Object, name)}
}

func GetVertexAttribf(src Attrib, pname Enum) float32 {
	return float32(c.Call("getVertexAttrib", src.Value, pname).Float())
}

func GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	println("GetVertexAttribfv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getVertexAttrib")
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = float32(result.Index(i).Float())
	}
}

func GetVertexAttribi(src Attrib, pname Enum) int32 {
	return int32(c.Call("getVertexAttrib", src.Value, pname).Int())
}

func GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	println("GetVertexAttribiv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getVertexAttrib")
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = int32(result.Index(i).Int())
	}
}

func Hint(target, mode Enum) {
	c.Call("hint", target, mode)
}

func IsBuffer(b Buffer) bool {
	return c.Call("isBuffer", b.Object).Bool()
}

func IsEnabled(cap Enum) bool {
	return c.Call("isEnabled", cap).Bool()
}

func IsFramebuffer(fb Framebuffer) bool {
	return c.Call("isFramebuffer", fb.Object).Bool()
}

func IsProgram(p Program) bool {
	return c.Call("isProgram", p.Object).Bool()
}

func IsRenderbuffer(rb Renderbuffer) bool {
	return c.Call("isRenderbuffer", rb.Object).Bool()
}

func IsShader(s Shader) bool {
	return c.Call("isShader", s.Object).Bool()
}

func IsTexture(t Texture) bool {
	return c.Call("isTexture", t.Object).Bool()
}

func LineWidth(width float32) {
	c.Call("lineWidth", width)
}

func LinkProgram(p Program) {
	c.Call("linkProgram", p.Object)
}

func PixelStorei(pname Enum, param int32) {
	c.Call("pixelStorei", pname, param)
}

func PolygonOffset(factor, units float32) {
	c.Call("polygonOffset", factor, units)
}

func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	println("ReadPixels: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	if ty == Enum(UNSIGNED_BYTE) {
		c.Call("readPixels", x, y, width, height, format, ty, dst)
	} else {
		tmpDst := make([]float32, len(dst)/4)
		c.Call("readPixels", x, y, width, height, format, ty, tmpDst)
		for i, f := range tmpDst {
			binary.LittleEndian.PutUint32(dst[i*4:], math.Float32bits(f))
		}
	}
}

func ReleaseShaderCompiler() {
	// do nothing
}

func RenderbufferStorage(target, internalFormat Enum, width, height int) {
	c.Call("renderbufferStorage", target, internalFormat, width, height)
}

func SampleCoverage(value float32, invert bool) {
	c.Call("sampleCoverage", value, invert)
}

func Scissor(x, y, width, height int32) {
	c.Call("scissor", x, y, width, height)
}

func ShaderSource(s Shader, src string) {
	c.Call("shaderSource", s.Object, src)
}

func StencilFunc(fn Enum, ref int, mask uint32) {
	c.Call("stencilFunc", fn, ref, mask)
}

func StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	c.Call("stencilFuncSeparate", face, fn, ref, mask)
}

func StencilMask(mask uint32) {
	c.Call("stencilMask", mask)
}

func StencilMaskSeparate(face Enum, mask uint32) {
	c.Call("stencilMaskSeparate", face, mask)
}

func StencilOp(fail, zfail, zpass Enum) {
	c.Call("stencilOp", fail, zfail, zpass)
}

func StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	c.Call("stencilOpSeparate", face, sfail, dpfail, dppass)
}

func TexImage2D(target Enum, level int, width, height int, format Enum, ty Enum, data []byte) {
	var p interface{}
	if data != nil {
		p = data
	}
	c.Call("texImage2D", target, level, format, width, height, 0, format, ty, p)
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	c.Call("texSubImage2D", target, level, x, y, width, height, format, ty, data)
}

func TexParameterf(target, pname Enum, param float32) {
	c.Call("texParameterf", target, pname, param)
}

func TexParameterfv(target, pname Enum, params []float32) {
	println("TexParameterfv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	for _, param := range params {
		c.Call("texParameterf", target, pname, param)
	}
}

func TexParameteri(target, pname Enum, param int) {
	c.Call("texParameteri", target, pname, param)
}

func TexParameteriv(target, pname Enum, params []int32) {
	println("TexParameteriv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	for _, param := range params {
		c.Call("texParameteri", target, pname, param)
	}
}

func Uniform1f(dst Uniform, v float32) {
	c.Call("uniform1f", dst.Object, v)
}

func Uniform1fv(dst Uniform, src []float32) {
	c.Call("uniform1fv", dst.Object, src)
}

func Uniform1i(dst Uniform, v int) {
	c.Call("uniform1i", dst.Object, v)
}

func Uniform1iv(dst Uniform, src []int32) {
	c.Call("uniform1iv", dst.Object, src)
}

func Uniform2f(dst Uniform, v0, v1 float32) {
	c.Call("uniform2f", dst.Object, v0, v1)
}

func Uniform2fv(dst Uniform, src []float32) {
	c.Call("uniform2fv", dst.Object, src)
}

func Uniform2i(dst Uniform, v0, v1 int) {
	c.Call("uniform2i", dst.Object, v0, v1)
}

func Uniform2iv(dst Uniform, src []int32) {
	c.Call("uniform2iv", dst.Object, src)
}

func Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.Call("uniform3f", dst.Object, v0, v1, v2)
}

func Uniform3fv(dst Uniform, src []float32) {
	c.Call("uniform3fv", dst.Object, src)
}

func Uniform3i(dst Uniform, v0, v1, v2 int32) {
	c.Call("uniform3i", dst.Object, v0, v1, v2)
}

func Uniform3iv(dst Uniform, src []int32) {
	c.Call("uniform3iv", dst.Object, src)
}

func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.Call("uniform4f", dst.Object, v0, v1, v2, v3)
}

func Uniform4fv(dst Uniform, src []float32) {
	c.Call("uniform4fv", dst.Object, src)
}

func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	c.Call("uniform4i", dst.Object, v0, v1, v2, v3)
}

func Uniform4iv(dst Uniform, src []int32) {
	c.Call("uniform4iv", dst.Object, src)
}

func UniformMatrix2fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix2fv", dst.Object, false, src)
}

func UniformMatrix3fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix3fv", dst.Object, false, src)
}

func UniformMatrix4fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix4fv", dst.Object, false, src)
}

func UseProgram(p Program) {
	c.Call("useProgram", p.Object)
}

func ValidateProgram(p Program) {
	c.Call("validateProgram", p.Object)
}

func VertexAttrib1f(dst Attrib, x float32) {
	c.Call("vertexAttrib1f", dst.Value, x)
}

func VertexAttrib1fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib1fv", dst.Value, src)
}

func VertexAttrib2f(dst Attrib, x, y float32) {
	c.Call("vertexAttrib2f", dst.Value, x, y)
}

func VertexAttrib2fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib2fv", dst.Value, src)
}

func VertexAttrib3f(dst Attrib, x, y, z float32) {
	c.Call("vertexAttrib3f", dst.Value, x, y, z)
}

func VertexAttrib3fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib3fv", dst.Value, src)
}

func VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	c.Call("vertexAttrib4f", dst.Value, x, y, z, w)
}

func VertexAttrib4fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib4fv", dst.Value, src)
}

func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.Call("vertexAttribPointer", dst.Value, size, ty, normalized, stride, offset)
}

func Viewport(x, y, width, height int) {
	c.Call("viewport", x, y, width, height)
}
