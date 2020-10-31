// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gl implements Go bindings for OpenGL ES 2.0 and ES 3.0.

The GL functions are defined on a Context object that is responsible for
tracking a GL context. Typically a windowing system package (such as
golang.org/x/exp/shiny/screen) will call NewContext and provide
a gl.Context for a user application.

If the gl package is compiled on a platform capable of supporting ES 3.0,
the gl.Context object also implements gl.Context3.

The bindings are deliberately minimal, staying as close the C API as
possible. The semantics of each function maps onto functions
described in the Khronos documentation:

https://www.khronos.org/opengles/sdk/docs/man/

One notable departure from the C API is the introduction of types
to represent common uses of GLint: Texture, Surface, Buffer, etc.

Debug Logging

A tracing version of the OpenGL bindings is behind the `gldebug` build
tag. It acts as a simplified version of apitrace. Build your Go binary
with

	-tags gldebug

and each call to a GL function will log its input, output, and any
error messages. For example,

	I/Fyne   (27668): gl.GenBuffers(1) [Buffer(70001)]
	I/Fyne   (27668): gl.BindBuffer(ARRAY_BUFFER, Buffer(70001))
	I/Fyne   (27668): gl.BufferData(ARRAY_BUFFER, 36, len(36), STATIC_DRAW)
	I/Fyne   (27668): gl.BindBuffer(ARRAY_BUFFER, Buffer(70001))
	I/Fyne   (27668): gl.VertexAttribPointer(Attrib(0), 6, FLOAT, false, 0, 0)  error: [INVALID_VALUE]

The gldebug tracing has very high overhead, so make sure to remove
the build tag before deploying any binaries.
*/
package gl // import "github.com/fyne-io/mobile/gl"

/*
Implementation details.

All GL function calls fill out a C.struct_fnargs and drop it on the work
queue. The Start function drains the work queue and hands over a batch
of calls to C.process which runs them. This allows multiple GL calls to
be executed in a single cgo call.

A GL call is marked as blocking if it returns a value, or if it takes a
Go pointer. In this case the call will not return until C.process sends a
value on the retvalue channel.

This implementation ensures any goroutine can make GL calls, but it does
not make the GL interface safe for simultaneous use by multiple goroutines.
For the purpose of analyzing this code for race conditions, picture two
separate goroutines: one blocked on gl.Start, and another making calls to
the gl package exported functions.
*/

//go:generate go run gendebug.go -o gldebug.go
