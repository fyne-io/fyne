// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux openbsd windows
// +build !gldebug

package gl

// TODO(crawshaw): should functions on specific types become methods? E.g.
//                 func (t Texture) Bind(target Enum)
//                 this seems natural in Go, but moves us slightly
//                 further away from the underlying OpenGL spec.

import (
	"math"
	"unsafe"
)

func (ctx *context) ActiveTexture(texture Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnActiveTexture,
			a0: texture.c(),
		},
	})
}

func (ctx *context) AttachShader(p Program, s Shader) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnAttachShader,
			a0: p.c(),
			a1: s.c(),
		},
	})
}

func (ctx *context) BindAttribLocation(p Program, a Attrib, name string) {
	s, free := ctx.cString(name)
	defer free()
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindAttribLocation,
			a0: p.c(),
			a1: a.c(),
			a2: s,
		},
		blocking: true,
	})
}

func (ctx *context) BindBuffer(target Enum, b Buffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindBuffer,
			a0: target.c(),
			a1: b.c(),
		},
	})
}

func (ctx *context) BindFramebuffer(target Enum, fb Framebuffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindFramebuffer,
			a0: target.c(),
			a1: fb.c(),
		},
	})
}

func (ctx *context) BindRenderbuffer(target Enum, rb Renderbuffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindRenderbuffer,
			a0: target.c(),
			a1: rb.c(),
		},
	})
}

func (ctx *context) BindTexture(target Enum, t Texture) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindTexture,
			a0: target.c(),
			a1: t.c(),
		},
	})
}

func (ctx *context) BindVertexArray(va VertexArray) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBindVertexArray,
			a0: va.c(),
		},
	})
}

func (ctx *context) BlendColor(red, green, blue, alpha float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlendColor,
			a0: uintptr(math.Float32bits(red)),
			a1: uintptr(math.Float32bits(green)),
			a2: uintptr(math.Float32bits(blue)),
			a3: uintptr(math.Float32bits(alpha)),
		},
	})
}

func (ctx *context) BlendEquation(mode Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlendEquation,
			a0: mode.c(),
		},
	})
}

func (ctx *context) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlendEquationSeparate,
			a0: modeRGB.c(),
			a1: modeAlpha.c(),
		},
	})
}

func (ctx *context) BlendFunc(sfactor, dfactor Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlendFunc,
			a0: sfactor.c(),
			a1: dfactor.c(),
		},
	})
}

func (ctx *context) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlendFuncSeparate,
			a0: sfactorRGB.c(),
			a1: dfactorRGB.c(),
			a2: sfactorAlpha.c(),
			a3: dfactorAlpha.c(),
		},
	})
}

func (ctx *context) BufferData(target Enum, src []byte, usage Enum) {
	parg := unsafe.Pointer(nil)
	if len(src) > 0 {
		parg = unsafe.Pointer(&src[0])
	}
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBufferData,
			a0: target.c(),
			a1: uintptr(len(src)),
			a2: usage.c(),
		},
		parg:     parg,
		blocking: true,
	})
}

func (ctx *context) BufferInit(target Enum, size int, usage Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBufferData,
			a0: target.c(),
			a1: uintptr(size),
			a2: usage.c(),
		},
		parg: unsafe.Pointer(nil),
	})
}

func (ctx *context) BufferSubData(target Enum, offset int, data []byte) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBufferSubData,
			a0: target.c(),
			a1: uintptr(offset),
			a2: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CheckFramebufferStatus(target Enum) Enum {
	return Enum(ctx.enqueue(call{
		args: fnargs{
			fn: glfnCheckFramebufferStatus,
			a0: target.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) Clear(mask Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnClear,
			a0: uintptr(mask),
		},
	})
}

func (ctx *context) ClearColor(red, green, blue, alpha float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnClearColor,
			a0: uintptr(math.Float32bits(red)),
			a1: uintptr(math.Float32bits(green)),
			a2: uintptr(math.Float32bits(blue)),
			a3: uintptr(math.Float32bits(alpha)),
		},
	})
}

func (ctx *context) ClearDepthf(d float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnClearDepthf,
			a0: uintptr(math.Float32bits(d)),
		},
	})
}

func (ctx *context) ClearStencil(s int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnClearStencil,
			a0: uintptr(s),
		},
	})
}

func (ctx *context) ColorMask(red, green, blue, alpha bool) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnColorMask,
			a0: glBoolean(red),
			a1: glBoolean(green),
			a2: glBoolean(blue),
			a3: glBoolean(alpha),
		},
	})
}

func (ctx *context) CompileShader(s Shader) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCompileShader,
			a0: s.c(),
		},
	})
}

func (ctx *context) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCompressedTexImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: internalformat.c(),
			a3: uintptr(width),
			a4: uintptr(height),
			a5: uintptr(border),
			a6: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCompressedTexSubImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(xoffset),
			a3: uintptr(yoffset),
			a4: uintptr(width),
			a5: uintptr(height),
			a6: format.c(),
			a7: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCopyTexImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: internalformat.c(),
			a3: uintptr(x),
			a4: uintptr(y),
			a5: uintptr(width),
			a6: uintptr(height),
			a7: uintptr(border),
		},
	})
}

func (ctx *context) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCopyTexSubImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(xoffset),
			a3: uintptr(yoffset),
			a4: uintptr(x),
			a5: uintptr(y),
			a6: uintptr(width),
			a7: uintptr(height),
		},
	})
}

func (ctx *context) CreateBuffer() Buffer {
	return Buffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenBuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateFramebuffer() Framebuffer {
	return Framebuffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenFramebuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateProgram() Program {
	return Program{
		Init: true,
		Value: uint32(ctx.enqueue(call{
			args: fnargs{
				fn: glfnCreateProgram,
			},
			blocking: true,
		},
		))}
}

func (ctx *context) CreateRenderbuffer() Renderbuffer {
	return Renderbuffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenRenderbuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateShader(ty Enum) Shader {
	return Shader{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnCreateShader,
			a0: uintptr(ty),
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateTexture() Texture {
	return Texture{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenTexture,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateVertexArray() VertexArray {
	return VertexArray{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenVertexArray,
		},
		blocking: true,
	}))}
}

func (ctx *context) CullFace(mode Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnCullFace,
			a0: mode.c(),
		},
	})
}

func (ctx *context) DeleteBuffer(v Buffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteBuffer,
			a0: v.c(),
		},
	})
}

func (ctx *context) DeleteFramebuffer(v Framebuffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteFramebuffer,
			a0: v.c(),
		},
	})
}

func (ctx *context) DeleteProgram(p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteProgram,
			a0: p.c(),
		},
	})
}

func (ctx *context) DeleteRenderbuffer(v Renderbuffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteRenderbuffer,
			a0: v.c(),
		},
	})
}

func (ctx *context) DeleteShader(s Shader) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteShader,
			a0: s.c(),
		},
	})
}

func (ctx *context) DeleteTexture(v Texture) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteTexture,
			a0: v.c(),
		},
	})
}

func (ctx *context) DeleteVertexArray(v VertexArray) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDeleteVertexArray,
			a0: v.c(),
		},
	})
}

func (ctx *context) DepthFunc(fn Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDepthFunc,
			a0: fn.c(),
		},
	})
}

func (ctx *context) DepthMask(flag bool) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDepthMask,
			a0: glBoolean(flag),
		},
	})
}

func (ctx *context) DepthRangef(n, f float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDepthRangef,
			a0: uintptr(math.Float32bits(n)),
			a1: uintptr(math.Float32bits(f)),
		},
	})
}

func (ctx *context) DetachShader(p Program, s Shader) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDetachShader,
			a0: p.c(),
			a1: s.c(),
		},
	})
}

func (ctx *context) Disable(cap Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDisable,
			a0: cap.c(),
		},
	})
}

func (ctx *context) DisableVertexAttribArray(a Attrib) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDisableVertexAttribArray,
			a0: a.c(),
		},
	})
}

func (ctx *context) DrawArrays(mode Enum, first, count int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDrawArrays,
			a0: mode.c(),
			a1: uintptr(first),
			a2: uintptr(count),
		},
	})
}

func (ctx *context) DrawElements(mode Enum, count int, ty Enum, offset int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnDrawElements,
			a0: mode.c(),
			a1: uintptr(count),
			a2: ty.c(),
			a3: uintptr(offset),
		},
	})
}

func (ctx *context) Enable(cap Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnEnable,
			a0: cap.c(),
		},
	})
}

func (ctx *context) EnableVertexAttribArray(a Attrib) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnEnableVertexAttribArray,
			a0: a.c(),
		},
	})
}

func (ctx *context) Finish() {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnFinish,
		},
		blocking: true,
	})
}

func (ctx *context) Flush() {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnFlush,
		},
		blocking: true,
	})
}

func (ctx *context) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnFramebufferRenderbuffer,
			a0: target.c(),
			a1: attachment.c(),
			a2: rbTarget.c(),
			a3: rb.c(),
		},
	})
}

func (ctx *context) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnFramebufferTexture2D,
			a0: target.c(),
			a1: attachment.c(),
			a2: texTarget.c(),
			a3: t.c(),
			a4: uintptr(level),
		},
	})
}

func (ctx *context) FrontFace(mode Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnFrontFace,
			a0: mode.c(),
		},
	})
}

func (ctx *context) GenerateMipmap(target Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenerateMipmap,
			a0: target.c(),
		},
	})
}

func (ctx *context) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	bufSize := ctx.GetProgrami(p, ACTIVE_ATTRIBUTE_MAX_LENGTH)
	buf := make([]byte, bufSize)
	var cType int

	cSize := ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetActiveAttrib,
			a0: p.c(),
			a1: uintptr(index),
			a2: uintptr(bufSize),
			a3: uintptr(unsafe.Pointer(&cType)), // TODO(crawshaw): not safe for a moving collector
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	return goString(buf), int(cSize), Enum(cType)
}

func (ctx *context) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	bufSize := ctx.GetProgrami(p, ACTIVE_UNIFORM_MAX_LENGTH)
	buf := make([]byte, bufSize+8) // extra space for cType
	var cType int

	cSize := ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetActiveUniform,
			a0: p.c(),
			a1: uintptr(index),
			a2: uintptr(bufSize),
			a3: uintptr(unsafe.Pointer(&cType)), // TODO(crawshaw): not safe for a moving collector
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	return goString(buf), int(cSize), Enum(cType)
}

func (ctx *context) GetAttachedShaders(p Program) []Shader {
	shadersLen := ctx.GetProgrami(p, ATTACHED_SHADERS)
	if shadersLen == 0 {
		return nil
	}
	buf := make([]uint32, shadersLen)

	n := int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetAttachedShaders,
			a0: p.c(),
			a1: uintptr(shadersLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	}))

	buf = buf[:int(n)]
	shaders := make([]Shader, len(buf))
	for i, s := range buf {
		shaders[i] = Shader{Value: uint32(s)}
	}
	return shaders
}

func (ctx *context) GetAttribLocation(p Program, name string) Attrib {
	s, free := ctx.cString(name)
	defer free()
	return Attrib{Value: uint(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetAttribLocation,
			a0: p.c(),
			a1: s,
		},
		blocking: true,
	}))}
}

func (ctx *context) GetBooleanv(dst []bool, pname Enum) {
	buf := make([]int32, len(dst))

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetBooleanv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	for i, v := range buf {
		dst[i] = v != 0
	}
}

func (ctx *context) GetFloatv(dst []float32, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetFloatv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetIntegerv(dst []int32, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetIntegerv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetInteger(pname Enum) int {
	var v [1]int32
	ctx.GetIntegerv(v[:], pname)
	return int(v[0])
}

func (ctx *context) GetBufferParameteri(target, value Enum) int {
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetBufferParameteri,
			a0: target.c(),
			a1: value.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetError() Enum {
	return Enum(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetError,
		},
		blocking: true,
	}))
}

func (ctx *context) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) int {
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetFramebufferAttachmentParameteriv,
			a0: target.c(),
			a1: attachment.c(),
			a2: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetProgrami(p Program, pname Enum) int {
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetProgramiv,
			a0: p.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetProgramInfoLog(p Program) string {
	infoLen := ctx.GetProgrami(p, INFO_LOG_LENGTH)
	if infoLen == 0 {
		return ""
	}
	buf := make([]byte, infoLen)

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetProgramInfoLog,
			a0: p.c(),
			a1: uintptr(infoLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	return goString(buf)
}

func (ctx *context) GetRenderbufferParameteri(target, pname Enum) int {
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetRenderbufferParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetShaderi(s Shader, pname Enum) int {
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetShaderiv,
			a0: s.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetShaderInfoLog(s Shader) string {
	infoLen := ctx.GetShaderi(s, INFO_LOG_LENGTH)
	if infoLen == 0 {
		return ""
	}
	buf := make([]byte, infoLen)

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetShaderInfoLog,
			a0: s.c(),
			a1: uintptr(infoLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	return goString(buf)
}

func (ctx *context) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	var rangeAndPrec [3]int32

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetShaderPrecisionFormat,
			a0: shadertype.c(),
			a1: precisiontype.c(),
		},
		parg:     unsafe.Pointer(&rangeAndPrec[0]),
		blocking: true,
	})

	return int(rangeAndPrec[0]), int(rangeAndPrec[1]), int(rangeAndPrec[2])
}

func (ctx *context) GetShaderSource(s Shader) string {
	sourceLen := ctx.GetShaderi(s, SHADER_SOURCE_LENGTH)
	if sourceLen == 0 {
		return ""
	}
	buf := make([]byte, sourceLen)

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetShaderSource,
			a0: s.c(),
			a1: uintptr(sourceLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})

	return goString(buf)
}

func (ctx *context) GetString(pname Enum) string {
	ret := ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetString,
			a0: pname.c(),
		},
		blocking: true,
	})
	retp := unsafe.Pointer(ret)
	buf := (*[1 << 24]byte)(retp)[:]
	return goString(buf)
}

func (ctx *context) GetTexParameterfv(dst []float32, target, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetTexParameterfv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetTexParameteriv(dst []int32, target, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetTexParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		blocking: true,
	})
}

func (ctx *context) GetUniformfv(dst []float32, src Uniform, p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetUniformfv,
			a0: p.c(),
			a1: src.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetUniformiv(dst []int32, src Uniform, p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetUniformiv,
			a0: p.c(),
			a1: src.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetUniformLocation(p Program, name string) Uniform {
	s, free := ctx.cString(name)
	defer free()
	return Uniform{Value: int32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetUniformLocation,
			a0: p.c(),
			a1: s,
		},
		blocking: true,
	}))}
}

func (ctx *context) GetVertexAttribf(src Attrib, pname Enum) float32 {
	var params [1]float32
	ctx.GetVertexAttribfv(params[:], src, pname)
	return params[0]
}

func (ctx *context) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetVertexAttribfv,
			a0: src.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetVertexAttribi(src Attrib, pname Enum) int32 {
	var params [1]int32
	ctx.GetVertexAttribiv(params[:], src, pname)
	return params[0]
}

func (ctx *context) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetVertexAttribiv,
			a0: src.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) Hint(target, mode Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnHint,
			a0: target.c(),
			a1: mode.c(),
		},
	})
}

func (ctx *context) IsBuffer(b Buffer) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsBuffer,
			a0: b.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsEnabled(cap Enum) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsEnabled,
			a0: cap.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsFramebuffer(fb Framebuffer) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsFramebuffer,
			a0: fb.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsProgram(p Program) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsProgram,
			a0: p.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsRenderbuffer(rb Renderbuffer) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsRenderbuffer,
			a0: rb.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsShader(s Shader) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsShader,
			a0: s.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsTexture(t Texture) bool {
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsTexture,
			a0: t.c(),
		},
		blocking: true,
	})
}

func (ctx *context) LineWidth(width float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnLineWidth,
			a0: uintptr(math.Float32bits(width)),
		},
	})
}

func (ctx *context) LinkProgram(p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnLinkProgram,
			a0: p.c(),
		},
	})
}

func (ctx *context) PixelStorei(pname Enum, param int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnPixelStorei,
			a0: pname.c(),
			a1: uintptr(param),
		},
	})
}

func (ctx *context) PolygonOffset(factor, units float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnPolygonOffset,
			a0: uintptr(math.Float32bits(factor)),
			a1: uintptr(math.Float32bits(units)),
		},
	})
}

func (ctx *context) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnReadPixels,
			// TODO(crawshaw): support PIXEL_PACK_BUFFER in GLES3, uses offset.
			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
			a4: format.c(),
			a5: ty.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) ReleaseShaderCompiler() {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnReleaseShaderCompiler,
		},
	})
}

func (ctx *context) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnRenderbufferStorage,
			a0: target.c(),
			a1: internalFormat.c(),
			a2: uintptr(width),
			a3: uintptr(height),
		},
	})
}

func (ctx *context) SampleCoverage(value float32, invert bool) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnSampleCoverage,
			a0: uintptr(math.Float32bits(value)),
			a1: glBoolean(invert),
		},
	})
}

func (ctx *context) Scissor(x, y, width, height int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnScissor,
			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
		},
	})
}

func (ctx *context) ShaderSource(s Shader, src string) {
	strp, free := ctx.cStringPtr(src)
	defer free()
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnShaderSource,
			a0: s.c(),
			a1: 1,
			a2: strp,
		},
		blocking: true,
	})
}

func (ctx *context) StencilFunc(fn Enum, ref int, mask uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilFunc,
			a0: fn.c(),
			a1: uintptr(ref),
			a2: uintptr(mask),
		},
	})
}

func (ctx *context) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilFuncSeparate,
			a0: face.c(),
			a1: fn.c(),
			a2: uintptr(ref),
			a3: uintptr(mask),
		},
	})
}

func (ctx *context) StencilMask(mask uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilMask,
			a0: uintptr(mask),
		},
	})
}

func (ctx *context) StencilMaskSeparate(face Enum, mask uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilMaskSeparate,
			a0: face.c(),
			a1: uintptr(mask),
		},
	})
}

func (ctx *context) StencilOp(fail, zfail, zpass Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilOp,
			a0: fail.c(),
			a1: zfail.c(),
			a2: zpass.c(),
		},
	})
}

func (ctx *context) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnStencilOpSeparate,
			a0: face.c(),
			a1: sfail.c(),
			a2: dpfail.c(),
			a3: dppass.c(),
		},
	})
}

func (ctx *context) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	// It is common to pass TexImage2D a nil data, indicating that a
	// bound GL buffer is being used as the source. In that case, it
	// is not necessary to block.
	parg := unsafe.Pointer(nil)
	if len(data) > 0 {
		parg = unsafe.Pointer(&data[0])
	}

	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexImage2D,
			// TODO(crawshaw): GLES3 offset for PIXEL_UNPACK_BUFFER and PIXEL_PACK_BUFFER.
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(internalFormat),
			a3: uintptr(width),
			a4: uintptr(height),
			a5: format.c(),
			a6: ty.c(),
		},
		parg:     parg,
		blocking: parg != nil,
	})
}

func (ctx *context) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexSubImage2D,
			// TODO(crawshaw): GLES3 offset for PIXEL_UNPACK_BUFFER and PIXEL_PACK_BUFFER.
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(x),
			a3: uintptr(y),
			a4: uintptr(width),
			a5: uintptr(height),
			a6: format.c(),
			a7: ty.c(),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) TexParameterf(target, pname Enum, param float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexParameterf,
			a0: target.c(),
			a1: pname.c(),
			a2: uintptr(math.Float32bits(param)),
		},
	})
}

func (ctx *context) TexParameterfv(target, pname Enum, params []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexParameterfv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&params[0]),
		blocking: true,
	})
}

func (ctx *context) TexParameteri(target, pname Enum, param int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexParameteri,
			a0: target.c(),
			a1: pname.c(),
			a2: uintptr(param),
		},
	})
}

func (ctx *context) TexParameteriv(target, pname Enum, params []int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnTexParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&params[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform1f(dst Uniform, v float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform1f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v)),
		},
	})
}

func (ctx *context) Uniform1fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform1fv,
			a0: dst.c(),
			a1: uintptr(len(src)),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform1i(dst Uniform, v int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform1i,
			a0: dst.c(),
			a1: uintptr(v),
		},
	})
}

func (ctx *context) Uniform1iv(dst Uniform, src []int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform1iv,
			a0: dst.c(),
			a1: uintptr(len(src)),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform2f(dst Uniform, v0, v1 float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform2f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
		},
	})
}

func (ctx *context) Uniform2fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 2),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform2i(dst Uniform, v0, v1 int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform2i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
		},
	})
}

func (ctx *context) Uniform2iv(dst Uniform, src []int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform2iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 2),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform3f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
			a3: uintptr(math.Float32bits(v2)),
		},
	})
}

func (ctx *context) Uniform3fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 3),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform3i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
		},
	})
}

func (ctx *context) Uniform3iv(dst Uniform, src []int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform3iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 3),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform4f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
			a3: uintptr(math.Float32bits(v2)),
			a4: uintptr(math.Float32bits(v3)),
		},
	})
}

func (ctx *context) Uniform4fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform4i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
			a4: uintptr(v3),
		},
	})
}

func (ctx *context) Uniform4iv(dst Uniform, src []int32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform4iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix2fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix2fv,
			// OpenGL ES 2 does not support transpose.
			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix3fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 9),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix4fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 16),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UseProgram(p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUseProgram,
			a0: p.c(),
		},
	})
}

func (ctx *context) ValidateProgram(p Program) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnValidateProgram,
			a0: p.c(),
		},
	})
}

func (ctx *context) VertexAttrib1f(dst Attrib, x float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib1f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
		},
	})
}

func (ctx *context) VertexAttrib1fv(dst Attrib, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib1fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib2f(dst Attrib, x, y float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib2f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
		},
	})
}

func (ctx *context) VertexAttrib2fv(dst Attrib, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib2fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib3f(dst Attrib, x, y, z float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib3f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
			a3: uintptr(math.Float32bits(z)),
		},
	})
}

func (ctx *context) VertexAttrib3fv(dst Attrib, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib3fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib4f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
			a3: uintptr(math.Float32bits(z)),
			a4: uintptr(math.Float32bits(w)),
		},
	})
}

func (ctx *context) VertexAttrib4fv(dst Attrib, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttrib4fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnVertexAttribPointer,
			a0: dst.c(),
			a1: uintptr(size),
			a2: ty.c(),
			a3: glBoolean(normalized),
			a4: uintptr(stride),
			a5: uintptr(offset),
		},
	})
}

func (ctx *context) Viewport(x, y, width, height int) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnViewport,
			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
		},
	})
}

func (ctx context3) UniformMatrix2x3fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix2x3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 6),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix3x2fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix3x2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 6),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix2x4fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix2x4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 8),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix4x2fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix4x2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 8),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix3x4fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix3x4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 12),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix4x3fv(dst Uniform, src []float32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniformMatrix4x3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 12),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int, mask uint, filter Enum) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnBlitFramebuffer,
			a0: uintptr(srcX0),
			a1: uintptr(srcY0),
			a2: uintptr(srcX1),
			a3: uintptr(srcY1),
			a4: uintptr(dstX0),
			a5: uintptr(dstY0),
			a6: uintptr(dstX1),
			a7: uintptr(dstY1),
			a8: uintptr(mask),
			a9: filter.c(),
		},
	})
}

func (ctx context3) Uniform1ui(dst Uniform, v uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform1ui,
			a0: dst.c(),
			a1: uintptr(v),
		},
	})
}

func (ctx context3) Uniform2ui(dst Uniform, v0, v1 uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform2ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
		},
	})
}

func (ctx context3) Uniform3ui(dst Uniform, v0, v1, v2 uint) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform3ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
		},
	})
}

func (ctx context3) Uniform4ui(dst Uniform, v0, v1, v2, v3 uint32) {
	ctx.enqueue(call{
		args: fnargs{
			fn: glfnUniform4ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
			a4: uintptr(v3),
		},
	})
}
