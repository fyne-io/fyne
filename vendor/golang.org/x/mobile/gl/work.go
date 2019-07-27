// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux openbsd

package gl

/*
#cgo ios                LDFLAGS: -framework OpenGLES
#cgo darwin,amd64,!ios  LDFLAGS: -framework OpenGL
#cgo darwin,arm         LDFLAGS: -framework OpenGLES
#cgo darwin,arm64       LDFLAGS: -framework OpenGLES
#cgo linux              LDFLAGS: -lGLESv2
#cgo openbsd            LDFLAGS: -L/usr/X11R6/lib/ -lGLESv2

#cgo android            CFLAGS: -Dos_android
#cgo ios                CFLAGS: -Dos_ios
#cgo darwin,amd64,!ios  CFLAGS: -Dos_osx
#cgo darwin,arm         CFLAGS: -Dos_ios
#cgo darwin,arm64       CFLAGS: -Dos_ios
#cgo darwin             CFLAGS: -DGL_SILENCE_DEPRECATION
#cgo linux              CFLAGS: -Dos_linux
#cgo openbsd            CFLAGS: -Dos_openbsd

#cgo openbsd            CFLAGS: -I/usr/X11R6/include/

#include <stdint.h>
#include "work.h"

uintptr_t process(struct fnargs* cargs, char* parg0, char* parg1, char* parg2, int count) {
	uintptr_t ret;

	ret = processFn(&cargs[0], parg0);
	if (count > 1) {
		ret = processFn(&cargs[1], parg1);
	}
	if (count > 2) {
		ret = processFn(&cargs[2], parg2);
	}

	return ret;
}
*/
import "C"

import "unsafe"

const workbufLen = 3

type context struct {
	cptr  uintptr
	debug int32

	workAvailable chan struct{}

	// work is a queue of calls to execute.
	work chan call

	// retvalue is sent a return value when blocking calls complete.
	// It is safe to use a global unbuffered channel here as calls
	// cannot currently be made concurrently.
	//
	// TODO: the comment above about concurrent calls isn't actually true: package
	// app calls package gl, but it has to do so in a separate goroutine, which
	// means that its gl calls (which may be blocking) can race with other gl calls
	// in the main program. We should make it safe to issue blocking gl calls
	// concurrently, or get the gl calls out of package app, or both.
	retvalue chan C.uintptr_t

	cargs [workbufLen]C.struct_fnargs
	parg  [workbufLen]*C.char
}

func (ctx *context) WorkAvailable() <-chan struct{} { return ctx.workAvailable }

type context3 struct {
	*context
}

// NewContext creates a cgo OpenGL context.
//
// See the Worker interface for more details on how it is used.
func NewContext() (Context, Worker) {
	glctx := &context{
		workAvailable: make(chan struct{}, 1),
		work:          make(chan call, workbufLen),
		retvalue:      make(chan C.uintptr_t),
	}
	if C.GLES_VERSION == "GL_ES_2_0" {
		return glctx, glctx
	}
	return context3{glctx}, glctx
}

// Version returns a GL ES version string, either "GL_ES_2_0" or "GL_ES_3_0".
// Future versions of the gl package may return "GL_ES_3_1".
func Version() string {
	return C.GLES_VERSION
}

func (ctx *context) enqueue(c call) uintptr {
	ctx.work <- c

	select {
	case ctx.workAvailable <- struct{}{}:
	default:
	}

	if c.blocking {
		return uintptr(<-ctx.retvalue)
	}
	return 0
}

func (ctx *context) DoWork() {
	queue := make([]call, 0, workbufLen)
	for {
		// Wait until at least one piece of work is ready.
		// Accumulate work until a piece is marked as blocking.
		select {
		case w := <-ctx.work:
			queue = append(queue, w)
		default:
			return
		}
		blocking := queue[len(queue)-1].blocking
	enqueue:
		for len(queue) < cap(queue) && !blocking {
			select {
			case w := <-ctx.work:
				queue = append(queue, w)
				blocking = queue[len(queue)-1].blocking
			default:
				break enqueue
			}
		}

		// Process the queued GL functions.
		for i, q := range queue {
			ctx.cargs[i] = *(*C.struct_fnargs)(unsafe.Pointer(&q.args))
			ctx.parg[i] = (*C.char)(q.parg)
		}
		ret := C.process(&ctx.cargs[0], ctx.parg[0], ctx.parg[1], ctx.parg[2], C.int(len(queue)))

		// Cleanup and signal.
		queue = queue[:0]
		if blocking {
			ctx.retvalue <- ret
		}
	}
}

func init() {
	if unsafe.Sizeof(C.GLint(0)) != unsafe.Sizeof(int32(0)) {
		panic("GLint is not an int32")
	}
}

// cString creates C string off the Go heap.
// ret is a *char.
func (ctx *context) cString(str string) (uintptr, func()) {
	ptr := unsafe.Pointer(C.CString(str))
	return uintptr(ptr), func() { C.free(ptr) }
}

// cString creates a pointer to a C string off the Go heap.
// ret is a **char.
func (ctx *context) cStringPtr(str string) (uintptr, func()) {
	s, free := ctx.cString(str)
	ptr := C.malloc(C.size_t(unsafe.Sizeof((*int)(nil))))
	*(*uintptr)(ptr) = s
	return uintptr(ptr), func() {
		free()
		C.free(ptr)
	}
}
