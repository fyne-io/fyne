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
	glfnBindBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glBindBuffer.Addr(), 2, c.args.a0, c.args.a1, 0)
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
	glfnBlendFunc: func(c call) (ret uintptr) {
		syscall.Syscall(glBlendFunc.Addr(), 2, c.args.a0, c.args.a1, 0)
		return
	},
	glfnBufferData: func(c call) (ret uintptr) {
		syscall.Syscall6(glBufferData.Addr(), 4, c.args.a0, c.args.a1, uintptr(c.parg), c.args.a2, 0, 0)
		return
	},
	glfnClear: func(c call) (ret uintptr) {
		syscall.Syscall(glClear.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnClearColor: func(c call) (ret uintptr) {
		syscall.Syscall6(glClearColor.Addr(), 4, c.args.a0, c.args.a1, c.args.a2, c.args.a3, 0, 0)
		return
	},
	glfnCompileShader: func(c call) (ret uintptr) {
		syscall.Syscall(glCompileShader.Addr(), 1, c.args.a0, 0, 0)
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
	glfnDeleteBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDeleteTexture: func(c call) (ret uintptr) {
		syscall.Syscall(glDeleteTextures.Addr(), 2, 1, uintptr(unsafe.Pointer(&c.args.a0)), 0)
		return
	},
	glfnDisable: func(c call) (ret uintptr) {
		syscall.Syscall(glDisable.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnDrawArrays: func(c call) (ret uintptr) {
		syscall.Syscall(glDrawArrays.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
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
	glfnFlush: func(c call) (ret uintptr) {
		syscall.Syscall(glFlush.Addr(), 0, 0, 0, 0)
		return
	},
	glfnGenBuffer: func(c call) (ret uintptr) {
		syscall.Syscall(glGenBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&ret)), 0)
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
	glfnGetAttribLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetAttribLocation.Addr(), 2, c.args.a0, c.args.a1, 0)
		return ret
	},
	glfnGetError: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetError.Addr(), 0, 0, 0, 0)
		return ret
	},
	glfnGetProgramInfoLog: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetProgramInfoLog.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
		return
	},
	glfnGetProgramiv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetProgramiv.Addr(), 3, c.args.a0, c.args.a1, uintptr(unsafe.Pointer(&ret)))
		return
	},
	glfnGetShaderInfoLog: func(c call) (ret uintptr) {
		syscall.Syscall6(glGetShaderInfoLog.Addr(), 4, c.args.a0, c.args.a1, 0, uintptr(c.parg), 0, 0)
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
	glfnGetTexParameteriv: func(c call) (ret uintptr) {
		syscall.Syscall(glGetTexParameteriv.Addr(), 3, c.args.a0, c.args.a1, uintptr(c.parg))
		return
	},
	glfnGetUniformLocation: func(c call) (ret uintptr) {
		ret, _, _ = syscall.Syscall(glGetUniformLocation.Addr(), 2, c.args.a0, c.args.a1, 0)
		return ret
	},
	glfnLinkProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glLinkProgram.Addr(), 1, c.args.a0, 0, 0)
		return
	},
	glfnReadPixels: func(c call) (ret uintptr) {
		syscall.Syscall9(glReadPixels.Addr(), 7, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5, uintptr(c.parg), 0, 0)
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
	glfnTexImage2D: func(c call) (ret uintptr) {
		syscall.Syscall9(glTexImage2D.Addr(), 9, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, 0, c.args.a5, c.args.a6, uintptr(c.parg))
		return
	},
	glfnTexParameteri: func(c call) (ret uintptr) {
		syscall.Syscall(glTexParameteri.Addr(), 3, c.args.a0, c.args.a1, c.args.a2)
		return
	},
	glfnUniform1f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform1f.Addr(), 2, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
		return
	},
	glfnUniform2f: func(c call) (ret uintptr) {
		syscall.Syscall6(glUniform2f.Addr(), 3, c.args.a0, c.args.a1, c.args.a2, c.args.a3, c.args.a4, c.args.a5)
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
	glfnUseProgram: func(c call) (ret uintptr) {
		syscall.Syscall(glUseProgram.Addr(), 1, c.args.a0, 0, 0)
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
	libGLESv2                 = LibGLESv2
	glActiveTexture           = libGLESv2.NewProc("glActiveTexture")
	glAttachShader            = libGLESv2.NewProc("glAttachShader")
	glBindBuffer              = libGLESv2.NewProc("glBindBuffer")
	glBindTexture             = libGLESv2.NewProc("glBindTexture")
	glBindVertexArray         = libGLESv2.NewProc("glBindVertexArray")
	glBlendColor              = libGLESv2.NewProc("glBlendColor")
	glBlendFunc               = libGLESv2.NewProc("glBlendFunc")
	glBufferData              = libGLESv2.NewProc("glBufferData")
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
