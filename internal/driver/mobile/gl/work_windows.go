// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"runtime"
	"syscall"
	"unsafe"
)

// context is described in work.go.
type context struct {
	debug         int32
	workAvailable chan struct{}
	work          chan call
	retvalue      chan uintptr

	// TODO(crawshaw): will not work with a moving collector
	cStringCounter int
	cStrings       map[int]unsafe.Pointer
}

func (ctx *context) WorkAvailable() <-chan struct{} { return ctx.workAvailable }

type context3 struct {
	*context
}

func NewContext() (Context, Worker) {
	if err := findDLLs(); err != nil {
		panic(err)
	}
	glctx := &context{
		workAvailable: make(chan struct{}, 1),
		work:          make(chan call, 3),
		retvalue:      make(chan uintptr),
		cStrings:      make(map[int]unsafe.Pointer),
	}
	return glctx, glctx
}

func (ctx *context) enqueue(c call) uintptr {
	ctx.work <- c

	select {
	case ctx.workAvailable <- struct{}{}:
	default:
	}

	if c.blocking {
		return <-ctx.retvalue
	}
	return 0
}

func (ctx *context) DoWork() {
	// TODO: add a work queue
	for {
		select {
		case w := <-ctx.work:
			ret := ctx.doWork(w)
			if w.blocking {
				ctx.retvalue <- ret
			}
		default:
			return
		}
	}
}

func (ctx *context) cString(s string) (uintptr, func()) {
	buf := make([]byte, len(s)+1)
	for i := 0; i < len(s); i++ {
		buf[i] = s[i]
	}
	ret := unsafe.Pointer(&buf[0])
	id := ctx.cStringCounter
	ctx.cStringCounter++
	ctx.cStrings[id] = ret
	return uintptr(ret), func() { delete(ctx.cStrings, id) }
}

func (ctx *context) cStringPtr(str string) (uintptr, func()) {
	s, sfree := ctx.cString(str)
	sptr := [2]uintptr{s, 0}
	ret := unsafe.Pointer(&sptr[0])
	id := ctx.cStringCounter
	ctx.cStringCounter++
	ctx.cStrings[id] = ret
	return uintptr(ret), func() { sfree(); delete(ctx.cStrings, id) }
}

// fixFloat copies the first four arguments into the XMM registers.
// This is for the windows/amd64 calling convention, that wants
// floating point arguments to be passed in XMM.
//
// Mercifully, type information is not required to implement
// this calling convention. In particular see the mixed int/float
// examples:
//
//	https://msdn.microsoft.com/en-us/library/zthk2dkh.aspx
//
// This means it could be fixed in syscall.Syscall. The relevant
// issue is
//
//	https://golang.org/issue/6510
func fixFloat(x0, x1, x2, x3 uintptr)

var glfnMap = map[glfn]func(c call) (ret uintptr){
	glfnActiveTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glActiveTexture.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnAttachShader: func(c call) (ret uintptr) {
		syscall.Syscall(glAttachShader.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBindAttribLocation: func(c call) (ret uintptr) {
		syscall.Syscall(glBindAttribLocation.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnBindBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glBindBuffer.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBindFramebuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glBindFramebuffer.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBindRenderbuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glBindRenderbuffer.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBindTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glBindTexture.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBindVertexArray: func(c call) (ret uintptr) {
		syscall.Syscall(glBindVertexArray.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnBlendColor: func(c call) (ret uintptr) {
		syscall.Syscall6(glBlendColor.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnBlendEquation: func(c call) (ret uintptr) {
		syscall.Syscall(glBlendEquation.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnBlendEquationSeparate: func(c call) (ret uintptr) {
		syscall.Syscall(glBlendEquationSeparate.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBlendFunc: func(c call) (ret uintptr) {
		syscall.Syscall(glBlendFunc.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBlendFuncSeparate: func(c call) (ret uintptr) {
		syscall.Syscall6(glBlendFuncSeparate.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnBufferData: func(c call) (ret uintptr) {
		syscall.Syscall6(glBufferData.Addr(), 4, c.args.a0, c.args.a1, uintptr(c.parg), c.args.a2, 0, 0)
		return
	},
	glfnBufferSubData: func(c call) (ret uintptr) {
		syscall.Syscall6(glBufferSubData.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, uintptr(c.parg), 0, 0)
		return
	},
	glfnCheckFramebufferStatus: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glCheckFramebufferStatus.Addr(), 1, c.args.a0, 0, 0)
		return ret
	},
	glfnClear: func(c call) (ret uintptr) {
		syscall.Syscall(glClear.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnClearColor: func(c call) (ret uintptr) {
		syscall.Syscall6(glClearColor.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnClearDepthf: func(c call) (ret uintptr) {
		syscall.Syscall6(glClearDepthf.Addr(), 1, c.args.a0, 0, 0, 0, 0, 0)
		return
	},
	glfnClearStencil: func(c call) (ret uintptr) {
		syscall.Syscall(glClearStencil.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnColorMask: func(c call) (ret uintptr) {
		syscall.Syscall6(glColorMask.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnCompileShader: func(c call) (ret uintptr) {
		syscall.Syscall(glCompileShader.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnCompressedTexImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glCompressedTexImage2D.Addr(), 8, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, c.args.a6, uintptr(c.parg), 0)
		return
	},
	glfnCompressedTexSubImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glCompressedTexSubImage2D.Addr(), 9, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, c.args.a6, c.args.a7, uintptr(c.parg))
		return
	},
	glfnCopyTexImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glCopyTexImage2D.Addr(), 8, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, c.args.a6, c.args.a7, 0)
		return
	},
	glfnCopyTexSubImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glCopyTexSubImage2D.Addr(), 8, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, c.args.a6, c.args.a7, 0)
		return
	},
	glfnCreateProgram: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glCreateProgram.Addr(), 0, 0, 0, 0)
		return ret
	},
	glfnCreateShader: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glCreateShader.Addr(), 1, c.args.a0, 0, 0)
		return ret
	},
	glfnCullFace: func(c call) (ret uintptr) {
		syscall.Syscall(glCullFace.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDeleteBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDeleteFramebuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteFramebuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDeleteProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteProgram.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDeleteRenderbuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteRenderbuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDeleteShader: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteShader.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDeleteVertexArray: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteVertexArrays.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDeleteTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteTextures.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDepthFunc: func(c call) (ret uintptr) {
		syscall.Syscall(glDepthFunc.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDepthRangef: func(c call) (ret uintptr) {
		syscall.Syscall6(glDepthRangef.Addr(), 2, c.args.a0, c.args.a1, 0, 0, 0, 0)
		return
	},
	glfnDepthMask: func(c call) (ret uintptr) {
		syscall.Syscall(glDepthMask.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDetachShader: func(c call) (ret uintptr) {
		syscall.Syscall(glDetachShader.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnDisable: func(c call) (ret uintptr) {
		syscall.Syscall(glDisable.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDisableVertexAttribArray: func(c call) (ret uintptr) {
		syscall.Syscall(glDisableVertexAttribArray.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDrawArrays: func(c call) (ret uintptr) {
		syscall.Syscall(glDrawArrays.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnDrawElements: func(c call) (ret uintptr) {
		syscall.Syscall6(glDrawElements.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnEnable: func(c call) (ret uintptr) {
		syscall.Syscall(glEnable.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnEnableVertexAttribArray: func(c call) (ret uintptr) {
		syscall.Syscall(glEnableVertexAttribArray.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnFinish: func(c call) (ret uintptr) {
		syscall.Syscall(glFinish.Addr(), 0, 0, 0, 0)
		return
	},
	glfnFlush: func(c call) (ret uintptr) {
		syscall.Syscall(glFlush.Addr(), 0, 0, 0, 0)
		return
	},
	glfnFramebufferRenderbuffer: func(c call) (ret uintptr) {
		syscall.Syscall6(glFramebufferRenderbuffer.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnFramebufferTexture2D: func(c call) (ret uintptr) {
		syscall.Syscall6(glFramebufferTexture2D.Addr(), 5, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, 0)
		return
	},
	glfnFrontFace: func(c call) (ret uintptr) {
		syscall.Syscall(glFrontFace.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnGenBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glGenBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
		return
	},
	glfnGenFramebuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glGenFramebuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
		return
	},
	glfnGenRenderbuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glGenRenderbuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
		return
	},
	glfnGenVertexArray: func(c call) (ret uintptr) {
		syscall.Syscall(glGenVertexArrays.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
		return
	},
	glfnGenTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glGenTextures.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
		return
	},
	glfnGenerateMipmap: func(c call) (ret uintptr) {
		syscall.Syscall(glGenerateMipmap.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnGetActiveAttrib: func(c call) (ret uintptr) {
		syscall.Syscall9(glGetActiveAttrib.Addr(), 7, c.args.a0, c.args.a1, c.args.a2, 0, uintptr(unsafe.Pointer(&ret)), c.args.a3, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetActiveUniform: func(c call) (ret uintptr) {
		syscall.Syscall9(glGetActiveUniform.Addr(), 7, c.args.a0, c.args.a1, c.args.a2, 0, uintptr(unsafe.Pointer(&ret)), c.args.a3, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetAttachedShaders: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetAttachedShaders.Addr(), 4, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)), uintptr(c.parg), 0, 0)
		return
	},
	glfnGetAttribLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetAttribLocation.Addr(), 2, c.args.a0, c.args.a1, 0)
		return ret
	},
	glfnGetBooleanv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetBooleanv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnGetBufferParameteri: func(c call) (ret uintptr) {
		syscall.Syscall(glGetBufferParameteri.Addr(), 3, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return
	},
	glfnGetError: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetError.Addr(), 0, 0, 0, 0)
		return ret
	},
	glfnGetFloatv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetFloatv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnGetFramebufferAttachmentParameteriv: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetFramebufferAttachmentParameteriv.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, uintptr(unsafe.Pointer(&ret)), 0, 0)
		return
	},
	glfnGetIntegerv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetIntegerv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnGetProgramInfoLog: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetProgramInfoLog.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetProgramiv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetProgramiv.Addr(), 3, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return
	},
	glfnGetRenderbufferParameteriv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetRenderbufferParameteriv.Addr(), 3, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return
	},
	glfnGetShaderInfoLog: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetShaderInfoLog.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetShaderPrecisionFormat: func(c call) (ret uintptr) {
		// c.parg is a [3]int32. The first [2]int32 of the array is one
		// parameter, the final *int32 is another parameter.
		syscall.Syscall6(glGetShaderPrecisionFormat.Addr(), 4, c.args.a0, c.args.a1, uintptr(c.parg), uintptr(c.parg)+2*unsafe.Sizeof(uintptr(0)), 0, 0)
		return
	},
	glfnGetShaderSource: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetShaderSource.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetShaderiv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetShaderiv.Addr(), 3, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return
	},
	glfnGetString: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetString.Addr(), 1, c.args.a0, 0, 0)
		return ret
	},
	glfnGetTexParameterfv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetTexParameterfv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetTexParameteriv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetTexParameteriv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetUniformLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetUniformLocation.Addr(), 2, c.args.a0, c.args.a1, 0)
		return ret
	},
	glfnGetUniformfv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetUniformfv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetUniformiv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetUniformiv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetVertexAttribfv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetVertexAttribfv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetVertexAttribiv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetVertexAttribiv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnHint: func(c call) (ret uintptr) {
		syscall.Syscall(glHint.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnIsBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glIsBuffer.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnIsEnabled: func(c call) (ret uintptr) {
		syscall.Syscall(glIsEnabled.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnIsFramebuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glIsFramebuffer.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnIsProgram: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glIsProgram.Addr(), 1, c.args.a0, 0, 0)
		return ret
	},
	glfnIsRenderbuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glIsRenderbuffer.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnIsShader: func(c call) (ret uintptr) {
		syscall.Syscall(glIsShader.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnIsTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glIsTexture.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnLineWidth: func(c call) (ret uintptr) {
		syscall.Syscall(glLineWidth.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnLinkProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glLinkProgram.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnPixelStorei: func(c call) (ret uintptr) {
		syscall.Syscall(glPixelStorei.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnPolygonOffset: func(c call) (ret uintptr) {
		syscall.Syscall6(glPolygonOffset.Addr(), 2, c.args.a0, c.args.a1, 0, 0, 0, 0)
		return
	},
	glfnReadPixels: func(c call) (ret uintptr) {
		syscall.Syscall9(glReadPixels.Addr(), 7, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, uintptr(c.parg), 0, 0)
		return
	},
	glfnReleaseShaderCompiler: func(c call) (ret uintptr) {
		syscall.Syscall(glReleaseShaderCompiler.Addr(), 0, 0, 0, 0)
		return
	},
	glfnRenderbufferStorage: func(c call) (ret uintptr) {
		syscall.Syscall6(glRenderbufferStorage.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnSampleCoverage: func(c call) (ret uintptr) {
		syscall.Syscall6(glSampleCoverage.Addr(), 1, c.args.a0, 0, 0, 0, 0, 0)
		return
	},
	glfnScissor: func(c call) (ret uintptr) {
		syscall.Syscall6(glScissor.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnShaderSource: func(c call) (ret uintptr) {
		syscall.Syscall6(glShaderSource.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, 0, 0, 0)
		return
	},
	glfnStencilFunc: func(c call) (ret uintptr) {
		syscall.Syscall(glStencilFunc.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnStencilFuncSeparate: func(c call) (ret uintptr) {
		syscall.Syscall6(glStencilFuncSeparate.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnStencilMask: func(c call) (ret uintptr) {
		syscall.Syscall(glStencilMask.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnStencilMaskSeparate: func(c call) (ret uintptr) {
		syscall.Syscall(glStencilMaskSeparate.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnStencilOp: func(c call) (ret uintptr) {
		syscall.Syscall(glStencilOp.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnStencilOpSeparate: func(c call) (ret uintptr) {
		syscall.Syscall6(glStencilOpSeparate.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnTexImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glTexImage2D.Addr(), 9, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, 0, c.args.a5, c.args.a6, uintptr(c.parg))
		return
	},
	glfnTexParameterf: func(c call) (ret uintptr) {
		syscall.Syscall6(glTexParameterf.Addr(), 3, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnTexParameterfv: func(c call) (ret uintptr) {
		syscall.Syscall(glTexParameterfv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnTexParameteri: func(c call) (ret uintptr) {
		syscall.Syscall(glTexParameteri.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnTexParameteriv: func(c call) (ret uintptr) {
		syscall.Syscall(glTexParameteriv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnTexSubImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glTexSubImage2D.Addr(), 9, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, c.args.a6, c.args.a7, uintptr(c.parg))
		return
	},
	glfnUniform1f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform1f.Addr(), 2, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnUniform1fv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform1fv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform1i: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform1i.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnUniform1iv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform1iv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform2f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform2f.Addr(), 3, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnUniform2fv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform2fv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform2i: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform2i.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnUniform2iv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform2iv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform3f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform3f.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnUniform3fv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform3fv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform3i: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform3i.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnUniform3iv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform3iv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform4f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform4f.Addr(), 5, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnUniform4fv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform4fv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniform4i: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform4i.Addr(), 5, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, 0)
		return
	},
	glfnUniform4iv: func(c call) (ret uintptr) {
		syscall.Syscall(glUniform4iv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnUniformMatrix2fv: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniformMatrix2fv.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnUniformMatrix3fv: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniformMatrix3fv.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnUniformMatrix4fv: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniformMatrix4fv.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnUseProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glUseProgram.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnValidateProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glValidateProgram.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnVertexAttrib1f: func(c call) (ret uintptr) {
		syscall.Syscall6(glVertexAttrib1f.Addr(), 2, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnVertexAttrib1fv: func(c call) (ret uintptr) {
		syscall.Syscall(glVertexAttrib1fv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnVertexAttrib2f: func(c call) (ret uintptr) {
		syscall.Syscall6(glVertexAttrib2f.Addr(), 3, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnVertexAttrib2fv: func(c call) (ret uintptr) {
		syscall.Syscall(glVertexAttrib2fv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnVertexAttrib3f: func(c call) (ret uintptr) {
		syscall.Syscall6(glVertexAttrib3f.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnVertexAttrib3fv: func(c call) (ret uintptr) {
		syscall.Syscall(glVertexAttrib3fv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnVertexAttrib4f: func(c call) (ret uintptr) {
		syscall.Syscall6(glVertexAttrib4f.Addr(), 5, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnVertexAttrib4fv: func(c call) (ret uintptr) {
		syscall.Syscall(glVertexAttrib4fv.Addr(), 2, c.args.a0, uintptr(c.parg), 0)
		return
	},
	glfnVertexAttribPointer: func(c call) (ret uintptr) {
		syscall.Syscall6(glVertexAttribPointer.Addr(), 6, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnViewport: func(c call) (ret uintptr) {
		syscall.Syscall6(glViewport.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
}

func (ctx *context) doWork(c call) (ret uintptr) {
	if runtime.GOARCH == "amd64" {
		fixFloat(c.args.a0, c.args.a1, c.args.a2, c.args.a3)
	}

	if f, ok := glfnMap[c.args.fn]; ok {
		return f(c)
	}
	panic("unknown GL function")
}

// Exported libraries for a Windows GUI driver.
//
// LibEGL is not used directly by the gl package, but is needed by any
// driver hoping to use OpenGL ES.
//
// LibD3DCompiler is needed by libglesv2.dll for compiling shaders.
var (
	LibGLESv2      = syscall.NewLazyDLL("libglesv2.dll")
	LibEGL         = syscall.NewLazyDLL("libegl.dll")
	LibD3DCompiler = syscall.NewLazyDLL("d3dcompiler_47.dll")
)

var (
	libGLESv2                             = LibGLESv2
	glActiveTexture                       = libGLESv2.NewProc("glActiveTexture")
	glAttachShader                        = libGLESv2.NewProc("glAttachShader")
	glBindAttribLocation                  = libGLESv2.NewProc("glBindAttribLocation")
	glBindBuffer                          = libGLESv2.NewProc("glBindBuffer")
	glBindFramebuffer                     = libGLESv2.NewProc("glBindFramebuffer")
	glBindRenderbuffer                    = libGLESv2.NewProc("glBindRenderbuffer")
	glBindTexture                         = libGLESv2.NewProc("glBindTexture")
	glBindVertexArray                     = libGLESv2.NewProc("glBindVertexArray")
	glBlendColor                          = libGLESv2.NewProc("glBlendColor")
	glBlendEquation                       = libGLESv2.NewProc("glBlendEquation")
	glBlendEquationSeparate               = libGLESv2.NewProc("glBlendEquationSeparate")
	glBlendFunc                           = libGLESv2.NewProc("glBlendFunc")
	glBlendFuncSeparate                   = libGLESv2.NewProc("glBlendFuncSeparate")
	glBufferData                          = libGLESv2.NewProc("glBufferData")
	glBufferSubData                       = libGLESv2.NewProc("glBufferSubData")
	glCheckFramebufferStatus              = libGLESv2.NewProc("glCheckFramebufferStatus")
	glClear                               = libGLESv2.NewProc("glClear")
	glClearColor                          = libGLESv2.NewProc("glClearColor")
	glClearDepthf                         = libGLESv2.NewProc("glClearDepthf")
	glClearStencil                        = libGLESv2.NewProc("glClearStencil")
	glColorMask                           = libGLESv2.NewProc("glColorMask")
	glCompileShader                       = libGLESv2.NewProc("glCompileShader")
	glCompressedTexImage2D                = libGLESv2.NewProc("glCompressedTexImage2D")
	glCompressedTexSubImage2D             = libGLESv2.NewProc("glCompressedTexSubImage2D")
	glCopyTexImage2D                      = libGLESv2.NewProc("glCopyTexImage2D")
	glCopyTexSubImage2D                   = libGLESv2.NewProc("glCopyTexSubImage2D")
	glCreateProgram                       = libGLESv2.NewProc("glCreateProgram")
	glCreateShader                        = libGLESv2.NewProc("glCreateShader")
	glCullFace                            = libGLESv2.NewProc("glCullFace")
	glDeleteBuffers                       = libGLESv2.NewProc("glDeleteBuffers")
	glDeleteFramebuffers                  = libGLESv2.NewProc("glDeleteFramebuffers")
	glDeleteProgram                       = libGLESv2.NewProc("glDeleteProgram")
	glDeleteRenderbuffers                 = libGLESv2.NewProc("glDeleteRenderbuffers")
	glDeleteShader                        = libGLESv2.NewProc("glDeleteShader")
	glDeleteTextures                      = libGLESv2.NewProc("glDeleteTextures")
	glDeleteVertexArrays                  = libGLESv2.NewProc("glDeleteVertexArrays")
	glDepthFunc                           = libGLESv2.NewProc("glDepthFunc")
	glDepthRangef                         = libGLESv2.NewProc("glDepthRangef")
	glDepthMask                           = libGLESv2.NewProc("glDepthMask")
	glDetachShader                        = libGLESv2.NewProc("glDetachShader")
	glDisable                             = libGLESv2.NewProc("glDisable")
	glDisableVertexAttribArray            = libGLESv2.NewProc("glDisableVertexAttribArray")
	glDrawArrays                          = libGLESv2.NewProc("glDrawArrays")
	glDrawElements                        = libGLESv2.NewProc("glDrawElements")
	glEnable                              = libGLESv2.NewProc("glEnable")
	glEnableVertexAttribArray             = libGLESv2.NewProc("glEnableVertexAttribArray")
	glFinish                              = libGLESv2.NewProc("glFinish")
	glFlush                               = libGLESv2.NewProc("glFlush")
	glFramebufferRenderbuffer             = libGLESv2.NewProc("glFramebufferRenderbuffer")
	glFramebufferTexture2D                = libGLESv2.NewProc("glFramebufferTexture2D")
	glFrontFace                           = libGLESv2.NewProc("glFrontFace")
	glGenBuffers                          = libGLESv2.NewProc("glGenBuffers")
	glGenFramebuffers                     = libGLESv2.NewProc("glGenFramebuffers")
	glGenRenderbuffers                    = libGLESv2.NewProc("glGenRenderbuffers")
	glGenTextures                         = libGLESv2.NewProc("glGenTextures")
	glGenVertexArrays                     = libGLESv2.NewProc("glGenVertexArrays")
	glGenerateMipmap                      = libGLESv2.NewProc("glGenerateMipmap")
	glGetActiveAttrib                     = libGLESv2.NewProc("glGetActiveAttrib")
	glGetActiveUniform                    = libGLESv2.NewProc("glGetActiveUniform")
	glGetAttachedShaders                  = libGLESv2.NewProc("glGetAttachedShaders")
	glGetAttribLocation                   = libGLESv2.NewProc("glGetAttribLocation")
	glGetBooleanv                         = libGLESv2.NewProc("glGetBooleanv")
	glGetBufferParameteri                 = libGLESv2.NewProc("glGetBufferParameteri")
	glGetError                            = libGLESv2.NewProc("glGetError")
	glGetFloatv                           = libGLESv2.NewProc("glGetFloatv")
	glGetFramebufferAttachmentParameteriv = libGLESv2.NewProc("glGetFramebufferAttachmentParameteriv")
	glGetIntegerv                         = libGLESv2.NewProc("glGetIntegerv")
	glGetProgramInfoLog                   = libGLESv2.NewProc("glGetProgramInfoLog")
	glGetProgramiv                        = libGLESv2.NewProc("glGetProgramiv")
	glGetRenderbufferParameteriv          = libGLESv2.NewProc("glGetRenderbufferParameteriv")
	glGetShaderInfoLog                    = libGLESv2.NewProc("glGetShaderInfoLog")
	glGetShaderPrecisionFormat            = libGLESv2.NewProc("glGetShaderPrecisionFormat")
	glGetShaderSource                     = libGLESv2.NewProc("glGetShaderSource")
	glGetShaderiv                         = libGLESv2.NewProc("glGetShaderiv")
	glGetString                           = libGLESv2.NewProc("glGetString")
	glGetTexParameterfv                   = libGLESv2.NewProc("glGetTexParameterfv")
	glGetTexParameteriv                   = libGLESv2.NewProc("glGetTexParameteriv")
	glGetUniformLocation                  = libGLESv2.NewProc("glGetUniformLocation")
	glGetUniformfv                        = libGLESv2.NewProc("glGetUniformfv")
	glGetUniformiv                        = libGLESv2.NewProc("glGetUniformiv")
	glGetVertexAttribfv                   = libGLESv2.NewProc("glGetVertexAttribfv")
	glGetVertexAttribiv                   = libGLESv2.NewProc("glGetVertexAttribiv")
	glHint                                = libGLESv2.NewProc("glHint")
	glIsBuffer                            = libGLESv2.NewProc("glIsBuffer")
	glIsEnabled                           = libGLESv2.NewProc("glIsEnabled")
	glIsFramebuffer                       = libGLESv2.NewProc("glIsFramebuffer")
	glIsProgram                           = libGLESv2.NewProc("glIsProgram")
	glIsRenderbuffer                      = libGLESv2.NewProc("glIsRenderbuffer")
	glIsShader                            = libGLESv2.NewProc("glIsShader")
	glIsTexture                           = libGLESv2.NewProc("glIsTexture")
	glLineWidth                           = libGLESv2.NewProc("glLineWidth")
	glLinkProgram                         = libGLESv2.NewProc("glLinkProgram")
	glPixelStorei                         = libGLESv2.NewProc("glPixelStorei")
	glPolygonOffset                       = libGLESv2.NewProc("glPolygonOffset")
	glReadPixels                          = libGLESv2.NewProc("glReadPixels")
	glReleaseShaderCompiler               = libGLESv2.NewProc("glReleaseShaderCompiler")
	glRenderbufferStorage                 = libGLESv2.NewProc("glRenderbufferStorage")
	glSampleCoverage                      = libGLESv2.NewProc("glSampleCoverage")
	glScissor                             = libGLESv2.NewProc("glScissor")
	glShaderSource                        = libGLESv2.NewProc("glShaderSource")
	glStencilFunc                         = libGLESv2.NewProc("glStencilFunc")
	glStencilFuncSeparate                 = libGLESv2.NewProc("glStencilFuncSeparate")
	glStencilMask                         = libGLESv2.NewProc("glStencilMask")
	glStencilMaskSeparate                 = libGLESv2.NewProc("glStencilMaskSeparate")
	glStencilOp                           = libGLESv2.NewProc("glStencilOp")
	glStencilOpSeparate                   = libGLESv2.NewProc("glStencilOpSeparate")
	glTexImage2D                          = libGLESv2.NewProc("glTexImage2D")
	glTexParameterf                       = libGLESv2.NewProc("glTexParameterf")
	glTexParameterfv                      = libGLESv2.NewProc("glTexParameterfv")
	glTexParameteri                       = libGLESv2.NewProc("glTexParameteri")
	glTexParameteriv                      = libGLESv2.NewProc("glTexParameteriv")
	glTexSubImage2D                       = libGLESv2.NewProc("glTexSubImage2D")
	glUniform1f                           = libGLESv2.NewProc("glUniform1f")
	glUniform1fv                          = libGLESv2.NewProc("glUniform1fv")
	glUniform1i                           = libGLESv2.NewProc("glUniform1i")
	glUniform1iv                          = libGLESv2.NewProc("glUniform1iv")
	glUniform2f                           = libGLESv2.NewProc("glUniform2f")
	glUniform2fv                          = libGLESv2.NewProc("glUniform2fv")
	glUniform2i                           = libGLESv2.NewProc("glUniform2i")
	glUniform2iv                          = libGLESv2.NewProc("glUniform2iv")
	glUniform3f                           = libGLESv2.NewProc("glUniform3f")
	glUniform3fv                          = libGLESv2.NewProc("glUniform3fv")
	glUniform3i                           = libGLESv2.NewProc("glUniform3i")
	glUniform3iv                          = libGLESv2.NewProc("glUniform3iv")
	glUniform4f                           = libGLESv2.NewProc("glUniform4f")
	glUniform4fv                          = libGLESv2.NewProc("glUniform4fv")
	glUniform4i                           = libGLESv2.NewProc("glUniform4i")
	glUniform4iv                          = libGLESv2.NewProc("glUniform4iv")
	glUniformMatrix2fv                    = libGLESv2.NewProc("glUniformMatrix2fv")
	glUniformMatrix3fv                    = libGLESv2.NewProc("glUniformMatrix3fv")
	glUniformMatrix4fv                    = libGLESv2.NewProc("glUniformMatrix4fv")
	glUseProgram                          = libGLESv2.NewProc("glUseProgram")
	glValidateProgram                     = libGLESv2.NewProc("glValidateProgram")
	glVertexAttrib1f                      = libGLESv2.NewProc("glVertexAttrib1f")
	glVertexAttrib1fv                     = libGLESv2.NewProc("glVertexAttrib1fv")
	glVertexAttrib2f                      = libGLESv2.NewProc("glVertexAttrib2f")
	glVertexAttrib2fv                     = libGLESv2.NewProc("glVertexAttrib2fv")
	glVertexAttrib3f                      = libGLESv2.NewProc("glVertexAttrib3f")
	glVertexAttrib3fv                     = libGLESv2.NewProc("glVertexAttrib3fv")
	glVertexAttrib4f                      = libGLESv2.NewProc("glVertexAttrib4f")
	glVertexAttrib4fv                     = libGLESv2.NewProc("glVertexAttrib4fv")
	glVertexAttribPointer                 = libGLESv2.NewProc("glVertexAttribPointer")
	glViewport                            = libGLESv2.NewProc("glViewport")
)
