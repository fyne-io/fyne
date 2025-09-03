// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
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

var glfnFuncs = [...]func(c call) (ret uintptr){
	glfnActiveTexture: func(c call) (ret uintptr) {
		syscall.SyscallN(glActiveTexture.Addr(), c.args.a0)
		return ret
	},
	glfnAttachShader: func(c call) (ret uintptr) {
		syscall.SyscallN(glAttachShader.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnBindBuffer: func(c call) (ret uintptr) {
		syscall.SyscallN(glBindBuffer.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnBindTexture: func(c call) (ret uintptr) {
		syscall.SyscallN(glBindTexture.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnBindVertexArray: func(c call) (ret uintptr) {
		syscall.SyscallN(glBindVertexArray.Addr(), c.args.a0)
		return ret
	},
	glfnBlendColor: func(c call) (ret uintptr) {
		syscall.SyscallN(glBlendColor.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3)
		return ret
	},
	glfnBlendFunc: func(c call) (ret uintptr) {
		syscall.SyscallN(glBlendFunc.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnBufferData: func(c call) (ret uintptr) {
		syscall.SyscallN(glBufferData.Addr(), c.args.a0, c.args.a1, uintptr(c.parg), c.args.a2)
		return ret
	},
	glfnBufferSubData: func(c call) (ret uintptr) {
		syscall.SyscallN(glBufferSubData.Addr(), c.args.a0, c.args.a1, c.args.a2, uintptr(c.parg))
		return ret
	},
	glfnClear: func(c call) (ret uintptr) {
		syscall.SyscallN(glClear.Addr(), c.args.a0)
		return ret
	},
	glfnClearColor: func(c call) (ret uintptr) {
		syscall.SyscallN(glClearColor.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3)
		return ret
	},
	glfnCompileShader: func(c call) (ret uintptr) {
		syscall.SyscallN(glCompileShader.Addr(), c.args.a0)
		return ret
	},
	glfnCreateProgram: func(c call) (ret uintptr) {
		ret, _, _ = syscall.SyscallN(glCreateProgram.Addr())
		return ret
	},
	glfnCreateShader: func(c call) (ret uintptr) {
		ret, _, _ = syscall.SyscallN(glCreateShader.Addr(), c.args.a0)
		return ret
	},
	glfnDeleteBuffer: func(c call) (ret uintptr) {
		syscall.SyscallN(glDeleteBuffers.Addr(), 1, uintptr(unsafe.Pointer(&c.args.a0)))
		return ret
	},
	glfnDeleteTexture: func(c call) (ret uintptr) {
		syscall.SyscallN(glDeleteTextures.Addr(), 1, uintptr(unsafe.Pointer(&c.args.a0)))
		return ret
	},
	glfnDisable: func(c call) (ret uintptr) {
		syscall.SyscallN(glDisable.Addr(), c.args.a0)
		return ret
	},
	glfnDrawArrays: func(c call) (ret uintptr) {
		syscall.SyscallN(glDrawArrays.Addr(), c.args.a0, c.args.a1, c.args.a2)
		return ret
	},
	glfnEnable: func(c call) (ret uintptr) {
		syscall.SyscallN(glEnable.Addr(), c.args.a0)
		return ret
	},
	glfnEnableVertexAttribArray: func(c call) (ret uintptr) {
		syscall.SyscallN(glEnableVertexAttribArray.Addr(), c.args.a0)
		return ret
	},
	glfnFlush: func(c call) (ret uintptr) {
		syscall.SyscallN(glFlush.Addr())
		return ret
	},
	glfnGenBuffer: func(c call) (ret uintptr) {
		syscall.SyscallN(glGenBuffers.Addr(), 1, uintptr(unsafe.Pointer(&ret)))
		return ret
	},
	glfnGenVertexArray: func(c call) (ret uintptr) {
		syscall.SyscallN(glGenVertexArrays.Addr(), 1, uintptr(unsafe.Pointer(&ret)))
		return ret
	},
	glfnGenTexture: func(c call) (ret uintptr) {
		syscall.SyscallN(glGenTextures.Addr(), 1, uintptr(unsafe.Pointer(&ret)))
		return ret
	},
	glfnGetAttribLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.SyscallN(glGetAttribLocation.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnGetError: func(c call) (ret uintptr) {
		ret, _, _ = syscall.SyscallN(glGetError.Addr())
		return ret
	},
	glfnGetProgramInfoLog: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetProgramInfoLog.Addr(), c.args.a0, c.args.a1, 0, uintptr(c.parg))
		return ret
	},
	glfnGetProgramiv: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetProgramiv.Addr(), c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return ret
	},
	glfnGetShaderInfoLog: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetShaderInfoLog.Addr(), c.args.a0, c.args.a1, 0, uintptr(c.parg))
		return ret
	},
	glfnGetShaderSource: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetShaderSource.Addr(), c.args.a0, c.args.a1, 0, uintptr(c.parg))
		return ret
	},
	glfnGetShaderiv: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetShaderiv.Addr(), c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return ret
	},
	glfnGetTexParameteriv: func(c call) (ret uintptr) {
		syscall.SyscallN(glGetTexParameteriv.Addr(), c.args.a0, c.args.a1, uintptr(c.parg))
		return ret
	},
	glfnGetUniformLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.SyscallN(glGetUniformLocation.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnLinkProgram: func(c call) (ret uintptr) {
		syscall.SyscallN(glLinkProgram.Addr(), c.args.a0)
		return ret
	},
	glfnReadPixels: func(c call) (ret uintptr) {
		syscall.SyscallN(glReadPixels.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, uintptr(c.parg))
		return ret
	},
	glfnScissor: func(c call) (ret uintptr) {
		syscall.SyscallN(glScissor.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3)
		return ret
	},
	glfnShaderSource: func(c call) (ret uintptr) {
		syscall.SyscallN(glShaderSource.Addr(), c.args.a0, c.args.a1, c.args.a2, 0)
		return ret
	},
	glfnTexImage2D: func(c call) (ret uintptr) {
		syscall.SyscallN(glTexImage2D.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, 0, c.args.a5, c.args.a6, uintptr(c.parg))
		return ret
	},
	glfnTexParameteri: func(c call) (ret uintptr) {
		syscall.SyscallN(glTexParameteri.Addr(), c.args.a0, c.args.a1, c.args.a2)
		return ret
	},
	glfnUniform1f: func(c call) (ret uintptr) {
		syscall.SyscallN(glUniform1f.Addr(), c.args.a0, c.args.a1)
		return ret
	},
	glfnUniform2f: func(c call) (ret uintptr) {
		syscall.SyscallN(glUniform2f.Addr(), c.args.a0, c.args.a1, c.args.a2)
		return ret
	},
	glfnUniform4f: func(c call) (ret uintptr) {
		syscall.SyscallN(glUniform4f.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4)
		return ret
	},
	glfnUniform4fv: func(c call) (ret uintptr) {
		syscall.SyscallN(glUniform4fv.Addr(), c.args.a0, c.args.a1, uintptr(c.parg))
		return ret
	},
	glfnUseProgram: func(c call) (ret uintptr) {
		syscall.SyscallN(glUseProgram.Addr(), c.args.a0)
		return ret
	},
	glfnVertexAttribPointer: func(c call) (ret uintptr) {
		syscall.SyscallN(glVertexAttribPointer.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return ret
	},
	glfnViewport: func(c call) (ret uintptr) {
		syscall.SyscallN(glViewport.Addr(), c.args.a0, c.args.a1, c.args.a2, c.args.a3)
		return ret
	},
}

func (ctx *context) doWork(c call) (ret uintptr) {
	if int(c.args.fn) < len(glfnFuncs) {
		return glfnFuncs[c.args.fn](c)
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
	libGLESv2                 = LibGLESv2
	glActiveTexture           = libGLESv2.NewProc("glActiveTexture")
	glAttachShader            = libGLESv2.NewProc("glAttachShader")
	glBindBuffer              = libGLESv2.NewProc("glBindBuffer")
	glBindTexture             = libGLESv2.NewProc("glBindTexture")
	glBindVertexArray         = libGLESv2.NewProc("glBindVertexArray")
	glBlendColor              = libGLESv2.NewProc("glBlendColor")
	glBlendFunc               = libGLESv2.NewProc("glBlendFunc")
	glBufferData              = libGLESv2.NewProc("glBufferData")
	glBufferSubData           = libGLESv2.NewProc("glBufferSubData")
	glClear                   = libGLESv2.NewProc("glClear")
	glClearColor              = libGLESv2.NewProc("glClearColor")
	glCompileShader           = libGLESv2.NewProc("glCompileShader")
	glCreateProgram           = libGLESv2.NewProc("glCreateProgram")
	glCreateShader            = libGLESv2.NewProc("glCreateShader")
	glDeleteBuffers           = libGLESv2.NewProc("glDeleteBuffers")
	glDeleteTextures          = libGLESv2.NewProc("glDeleteTextures")
	glDisable                 = libGLESv2.NewProc("glDisable")
	glDrawArrays              = libGLESv2.NewProc("glDrawArrays")
	glEnable                  = libGLESv2.NewProc("glEnable")
	glEnableVertexAttribArray = libGLESv2.NewProc("glEnableVertexAttribArray")
	glFlush                   = libGLESv2.NewProc("glFlush")
	glGenBuffers              = libGLESv2.NewProc("glGenBuffers")
	glGenTextures             = libGLESv2.NewProc("glGenTextures")
	glGenVertexArrays         = libGLESv2.NewProc("glGenVertexArrays")
	glGetAttribLocation       = libGLESv2.NewProc("glGetAttribLocation")
	glGetError                = libGLESv2.NewProc("glGetError")
	glGetProgramInfoLog       = libGLESv2.NewProc("glGetProgramInfoLog")
	glGetProgramiv            = libGLESv2.NewProc("glGetProgramiv")
	glGetShaderInfoLog        = libGLESv2.NewProc("glGetShaderInfoLog")
	glGetShaderSource         = libGLESv2.NewProc("glGetShaderSource")
	glGetShaderiv             = libGLESv2.NewProc("glGetShaderiv")
	glGetTexParameteriv       = libGLESv2.NewProc("glGetTexParameteriv")
	glGetUniformLocation      = libGLESv2.NewProc("glGetUniformLocation")
	glPixelStorei             = libGLESv2.NewProc("glPixelStorei")
	glLinkProgram             = libGLESv2.NewProc("glLinkProgram")
	glReadPixels              = libGLESv2.NewProc("glReadPixels")
	glScissor                 = libGLESv2.NewProc("glScissor")
	glShaderSource            = libGLESv2.NewProc("glShaderSource")
	glTexImage2D              = libGLESv2.NewProc("glTexImage2D")
	glTexParameteri           = libGLESv2.NewProc("glTexParameteri")
	glUniform1f               = libGLESv2.NewProc("glUniform1f")
	glUniform2f               = libGLESv2.NewProc("glUniform2f")
	glUniform4f               = libGLESv2.NewProc("glUniform4f")
	glUniform4fv              = libGLESv2.NewProc("glUniform4fv")
	glUseProgram              = libGLESv2.NewProc("glUseProgram")
	glVertexAttribPointer     = libGLESv2.NewProc("glVertexAttribPointer")
	glViewport                = libGLESv2.NewProc("glViewport")
)
