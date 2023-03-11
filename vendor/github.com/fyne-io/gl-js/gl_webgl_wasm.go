// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm
// +build js,wasm

package gl

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"syscall/js"
	"unsafe"
)

var ContextWatcher contextWatcher

type contextWatcher struct{}

func (contextWatcher) OnMakeCurrent(context interface{}) {
	// context must be a WebGLRenderingContext js.Value.
	c = context.(js.Value)
}
func (contextWatcher) OnDetach() {
	c = js.Null()
}

func sliceToByteSlice(s interface{}) []byte {
	switch s := s.(type) {
	case []int8:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		return *(*[]byte)(unsafe.Pointer(h))
	case []int16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []int32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []int64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint8:
		return s
	case []uint16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []float32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []float64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at sliceToBytesSlice: %T", s))
	}
}

func SliceToTypedArray(s interface{}) js.Value {
	if s == nil {
		return js.Null()
	}

	switch s := s.(type) {
	case []int8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int8Array").New(buf, a.Get("byteOffset"), a.Get("byteLength"))
	case []int16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []int32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []uint8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, s)
		runtime.KeepAlive(s)
		return a
	case []uint16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []uint32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float64:
		a := js.Global().Get("Uint8Array").New(len(s) * 8)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float64Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at SliceToTypedArray: %T", s))
	}
}

// c is the current WebGL context, or nil if there is no current context.
var c js.Value

func ActiveTexture(texture Enum) {
	c.Call("activeTexture", int(texture))
}

func AttachShader(p Program, s Shader) {
	c.Call("attachShader", p.Value, s.Value)
}

func BindAttribLocation(p Program, a Attrib, name string) {
	c.Call("bindAttribLocation", p.Value, a.Value, name)
}

func BindBuffer(target Enum, b Buffer) {
	c.Call("bindBuffer", int(target), b.Value)
}

func BindFramebuffer(target Enum, fb Framebuffer) {
	c.Call("bindFramebuffer", int(target), fb.Value)
}

func BindRenderbuffer(target Enum, rb Renderbuffer) {
	c.Call("bindRenderbuffer", int(target), rb.Value)
}

func BindTexture(target Enum, t Texture) {
	c.Call("bindTexture", int(target), t.Value)
}

func BlendColor(red, green, blue, alpha float32) {
	c.Call("blendColor", red, green, blue, alpha)
}

func BlendEquation(mode Enum) {
	c.Call("blendEquation", int(mode))
}

func BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func BlendFunc(sfactor, dfactor Enum) {
	c.Call("blendFunc", int(sfactor), int(dfactor))
}

func BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	c.Call("blendFuncSeparate", int(sfactorRGB), int(dfactorRGB), int(sfactorAlpha), int(dfactorAlpha))
}

func BufferData(target Enum, data interface{}, usage Enum) {
	c.Call("bufferData", int(target), SliceToTypedArray(data), int(usage))
}

func BufferInit(target Enum, size int, usage Enum) {
	c.Call("bufferData", int(target), size, int(usage))
}

func BufferSubData(target Enum, offset int, data interface{}) {
	c.Call("bufferSubData", int(target), offset, SliceToTypedArray(data))
}

func CheckFramebufferStatus(target Enum) Enum {
	return Enum(c.Call("checkFramebufferStatus", int(target)).Int())
}

func Clear(mask Enum) {
	c.Call("clear", int(mask))
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
	c.Call("compileShader", s.Value)
}

func CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data interface{}) {
	c.Call("compressedTexImage2D", int(target), level, internalformat, width, height, border, SliceToTypedArray(data))
}

func CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data interface{}) {
	c.Call("compressedTexSubImage2D", int(target), level, xoffset, yoffset, width, height, format, SliceToTypedArray(data))
}

func CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	c.Call("copyTexImage2D", int(target), level, internalformat, x, y, width, height, border)
}

func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	c.Call("copyTexSubImage2D", int(target), level, xoffset, yoffset, x, y, width, height)
}

func CreateBuffer() Buffer {
	return Buffer{Value: c.Call("createBuffer")}
}

func CreateFramebuffer() Framebuffer {
	return Framebuffer{Value: c.Call("createFramebuffer")}
}

func CreateProgram() Program {
	return Program{Value: c.Call("createProgram")}
}

func CreateRenderbuffer() Renderbuffer {
	return Renderbuffer{Value: c.Call("createRenderbuffer")}
}

func CreateShader(ty Enum) Shader {
	return Shader{Value: c.Call("createShader", int(ty))}
}

func CreateTexture() Texture {
	return Texture{Value: c.Call("createTexture")}
}

func CullFace(mode Enum) {
	c.Call("cullFace", int(mode))
}

func DeleteBuffer(v Buffer) {
	c.Call("deleteBuffer", v.Value)
}

func DeleteFramebuffer(v Framebuffer) {
	c.Call("deleteFramebuffer", v.Value)
}

func DeleteProgram(p Program) {
	c.Call("deleteProgram", p.Value)
}

func DeleteRenderbuffer(v Renderbuffer) {
	c.Call("deleteRenderbuffer", v.Value)
}

func DeleteShader(s Shader) {
	c.Call("deleteShader", s.Value)
}

func DeleteTexture(v Texture) {
	c.Call("deleteTexture", v.Value)
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
	c.Call("detachShader", p.Value, s.Value)
}

func Disable(cap Enum) {
	c.Call("disable", int(cap))
}

func DisableVertexAttribArray(a Attrib) {
	c.Call("disableVertexAttribArray", a.Value)
}

func DrawArrays(mode Enum, first, count int) {
	c.Call("drawArrays", int(mode), first, count)
}

func DrawElements(mode Enum, count int, ty Enum, offset int) {
	c.Call("drawElements", int(mode), count, int(ty), offset)
}

func Enable(cap Enum) {
	c.Call("enable", int(cap))
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
	c.Call("framebufferRenderbuffer", target, attachment, int(rbTarget), rb.Value)
}

func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	c.Call("framebufferTexture2D", target, attachment, int(texTarget), t.Value, level)
}

func FrontFace(mode Enum) {
	c.Call("frontFace", int(mode))
}

func GenerateMipmap(target Enum) {
	c.Call("generateMipmap", int(target))
}

func GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	ai := c.Call("getActiveAttrib", p.Value, index)
	return ai.Get("name").String(), ai.Get("size").Int(), Enum(ai.Get("type").Int())
}

func GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	ai := c.Call("getActiveUniform", p.Value, index)
	return ai.Get("name").String(), ai.Get("size").Int(), Enum(ai.Get("type").Int())
}

func GetAttachedShaders(p Program) []Shader {
	objs := c.Call("getAttachedShaders", p.Value)
	shaders := make([]Shader, objs.Length())
	for i := 0; i < objs.Length(); i++ {
		shaders[i] = Shader{Value: objs.Index(i)}
	}
	return shaders
}

func GetAttribLocation(p Program, name string) Attrib {
	return Attrib{Value: c.Call("getAttribLocation", p.Value, name).Int()}
}

func GetBooleanv(dst []bool, pname Enum) {
	println("GetBooleanv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getParameter", int(pname))
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = result.Index(i).Bool()
	}
}

func GetFloatv(dst []float32, pname Enum) {
	println("GetFloatv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	result := c.Call("getParameter", int(pname))
	length := result.Length()
	for i := 0; i < length; i++ {
		dst[i] = float32(result.Index(i).Float())
	}
}

func GetIntegerv(pname Enum, data []int32) {
	result := c.Call("getParameter", int(pname))
	length := result.Length()
	for i := 0; i < length; i++ {
		data[i] = int32(result.Index(i).Int())
	}
}

func GetInteger(pname Enum) int {
	return c.Call("getParameter", int(pname)).Int()
}

func GetBufferParameteri(target, pname Enum) int {
	return c.Call("getBufferParameter", int(target), int(pname)).Int()
}

func GetError() Enum {
	return Enum(c.Call("getError").Int())
}

func GetBoundFramebuffer() Framebuffer {
	return Framebuffer{Value: c.Call("getParameter", FRAMEBUFFER_BINDING)}
}

func GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return c.Call("getFramebufferAttachmentParameter", int(target), int(attachment), int(pname)).Int()
}

func GetProgrami(p Program, pname Enum) int {
	switch pname {
	case DELETE_STATUS, LINK_STATUS, VALIDATE_STATUS:
		if c.Call("getProgramParameter", p.Value, int(pname)).Bool() {
			return TRUE
		}
		return FALSE
	default:
		return c.Call("getProgramParameter", p.Value, int(pname)).Int()
	}
}

func GetProgramInfoLog(p Program) string {
	return c.Call("getProgramInfoLog", p.Value).String()
}

func GetRenderbufferParameteri(target, pname Enum) int {
	return c.Call("getRenderbufferParameter", int(target), int(pname)).Int()
}

func GetShaderi(s Shader, pname Enum) int {
	switch pname {
	case DELETE_STATUS, COMPILE_STATUS:
		if c.Call("getShaderParameter", s.Value, int(pname)).Bool() {
			return TRUE
		}
		return FALSE
	default:
		return c.Call("getShaderParameter", s.Value, int(pname)).Int()
	}
}

func GetShaderInfoLog(s Shader) string {
	return c.Call("getShaderInfoLog", s.Value).String()
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
	return c.Call("getShaderSource", s.Value).String()
}

func GetString(pname Enum) string {
	return c.Call("getParameter", int(pname)).String()
}

func GetTexParameterfv(dst []float32, target, pname Enum) {
	dst[0] = float32(c.Call("getTexParameter", int(pname)).Float())
}

func GetTexParameteriv(dst []int32, target, pname Enum) {
	dst[0] = int32(c.Call("getTexParameter", int(pname)).Int())
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
	return Uniform{Value: c.Call("getUniformLocation", p.Value, name)}
}

func GetVertexAttribf(src Attrib, pname Enum) float32 {
	return float32(c.Call("getVertexAttrib", src.Value, int(pname)).Float())
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
	return int32(c.Call("getVertexAttrib", src.Value, int(pname)).Int())
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
	c.Call("hint", int(target), int(mode))
}

func IsBuffer(b Buffer) bool {
	return c.Call("isBuffer", b.Value).Bool()
}

func IsEnabled(cap Enum) bool {
	return c.Call("isEnabled", int(cap)).Bool()
}

func IsFramebuffer(fb Framebuffer) bool {
	return c.Call("isFramebuffer", fb.Value).Bool()
}

func IsProgram(p Program) bool {
	return c.Call("isProgram", p.Value).Bool()
}

func IsRenderbuffer(rb Renderbuffer) bool {
	return c.Call("isRenderbuffer", rb.Value).Bool()
}

func IsShader(s Shader) bool {
	return c.Call("isShader", s.Value).Bool()
}

func IsTexture(t Texture) bool {
	return c.Call("isTexture", t.Value).Bool()
}

func LineWidth(width float32) {
	c.Call("lineWidth", width)
}

func LinkProgram(p Program) {
	c.Call("linkProgram", p.Value)
}

func PixelStorei(pname Enum, param int32) {
	c.Call("pixelStorei", int(pname), param)
}

func PolygonOffset(factor, units float32) {
	c.Call("polygonOffset", factor, units)
}

func ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	println("ReadPixels: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	if ty == Enum(UNSIGNED_BYTE) {
		c.Call("readPixels", x, y, width, height, format, int(ty), dst)
	} else {
		tmpDst := make([]float32, len(dst)/4)
		c.Call("readPixels", x, y, width, height, format, int(ty), tmpDst)
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
	c.Call("shaderSource", s.Value, src)
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

func TexImage2D(target Enum, level int, width, height int, format Enum, ty Enum, data interface{}) {
	c.Call("texImage2D", int(target), level, int(format), width, height, 0, int(format), int(ty), SliceToTypedArray(data))
}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data interface{}) {
	c.Call("texSubImage2D", int(target), level, x, y, width, height, format, int(ty), SliceToTypedArray(data))
}

func TexParameterf(target, pname Enum, param float32) {
	c.Call("texParameterf", int(target), int(pname), param)
}

func TexParameterfv(target, pname Enum, params []float32) {
	println("TexParameterfv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	for _, param := range params {
		c.Call("texParameterf", int(target), int(pname), SliceToTypedArray(param))
	}
}

func TexParameteri(target, pname Enum, param int) {
	c.Call("texParameteri", int(target), int(pname), param)
}

func TexParameteriv(target, pname Enum, params []int32) {
	println("TexParameteriv: not yet tested (TODO: remove this after it's confirmed to work. Your feedback is welcome.)")
	for _, param := range params {
		c.Call("texParameteri", int(target), int(pname), SliceToTypedArray(param))
	}
}

func Uniform1f(dst Uniform, v float32) {
	c.Call("uniform1f", dst.Value, v)
}

func Uniform1fv(dst Uniform, src []float32) {
	c.Call("uniform1fv", dst.Value, SliceToTypedArray(src))
}

func Uniform1i(dst Uniform, v int) {
	c.Call("uniform1i", dst.Value, v)
}

func Uniform1iv(dst Uniform, src []int32) {
	c.Call("uniform1iv", dst.Value, SliceToTypedArray(src))
}

func Uniform2f(dst Uniform, v0, v1 float32) {
	c.Call("uniform2f", dst.Value, v0, v1)
}

func Uniform2fv(dst Uniform, src []float32) {
	c.Call("uniform2fv", dst.Value, SliceToTypedArray(src))
}

func Uniform2i(dst Uniform, v0, v1 int) {
	c.Call("uniform2i", dst.Value, v0, v1)
}

func Uniform2iv(dst Uniform, src []int32) {
	c.Call("uniform2iv", dst.Value, SliceToTypedArray(src))
}

func Uniform3f(dst Uniform, v0, v1, v2 float32) {
	c.Call("uniform3f", dst.Value, v0, v1, v2)
}

func Uniform3fv(dst Uniform, src []float32) {
	c.Call("uniform3fv", dst.Value, SliceToTypedArray(src))
}

func Uniform3i(dst Uniform, v0, v1, v2 int32) {
	c.Call("uniform3i", dst.Value, v0, v1, v2)
}

func Uniform3iv(dst Uniform, src []int32) {
	c.Call("uniform3iv", dst.Value, SliceToTypedArray(src))
}

func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	c.Call("uniform4f", dst.Value, v0, v1, v2, v3)
}

func Uniform4fv(dst Uniform, src []float32) {
	c.Call("uniform4fv", dst.Value, SliceToTypedArray(src))
}

func Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	c.Call("uniform4i", dst.Value, v0, v1, v2, v3)
}

func Uniform4iv(dst Uniform, src []int32) {
	c.Call("uniform4iv", dst.Value, SliceToTypedArray(src))
}

func UniformMatrix2fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix2fv", dst.Value, false, SliceToTypedArray(src))
}

func UniformMatrix3fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix3fv", dst.Value, false, SliceToTypedArray(src))
}

func UniformMatrix4fv(dst Uniform, src []float32) {
	c.Call("uniformMatrix4fv", dst.Value, false, SliceToTypedArray(src))
}

func UseProgram(p Program) {
	// Workaround for js.Value zero value.
	if p.Value.Equal(js.Value{}) {
		p.Value = js.Null()
	}
	c.Call("useProgram", p.Value)
}

func ValidateProgram(p Program) {
	c.Call("validateProgram", p.Value)
}

func VertexAttrib1f(dst Attrib, x float32) {
	c.Call("vertexAttrib1f", dst.Value, x)
}

func VertexAttrib1fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib1fv", dst.Value, SliceToTypedArray(src))
}

func VertexAttrib2f(dst Attrib, x, y float32) {
	c.Call("vertexAttrib2f", dst.Value, x, y)
}

func VertexAttrib2fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib2fv", dst.Value, SliceToTypedArray(src))
}

func VertexAttrib3f(dst Attrib, x, y, z float32) {
	c.Call("vertexAttrib3f", dst.Value, x, y, z)
}

func VertexAttrib3fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib3fv", dst.Value, SliceToTypedArray(src))
}

func VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	c.Call("vertexAttrib4f", dst.Value, x, y, z, w)
}

func VertexAttrib4fv(dst Attrib, src []float32) {
	c.Call("vertexAttrib4fv", dst.Value, SliceToTypedArray(src))
}

func VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	c.Call("vertexAttribPointer", dst.Value, size, int(ty), normalized, stride, offset)
}

func Viewport(x, y, width, height int) {
	c.Call("viewport", x, y, width, height)
}
