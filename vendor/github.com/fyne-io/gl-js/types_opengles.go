// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin || linux) && (arm || arm64)
// +build darwin linux
// +build arm arm64

package gl

/*
#cgo ios     LDFLAGS: -framework OpenGLES
#cgo android LDFLAGS: -lGLESv2

#cgo ios     CFLAGS: -Dos_ios
#cgo android CFLAGS: -Dos_android

#ifdef os_ios
#include <OpenGLES/ES2/gl.h>
#endif
#ifdef os_android
#include <GLES2/gl2.h>
#endif

void blendColor(GLfloat r, GLfloat g, GLfloat b, GLfloat a) { glBlendColor(r, g, b, a); }
void clearColor(GLfloat r, GLfloat g, GLfloat b, GLfloat a) { glClearColor(r, g, b, a); }
void clearDepthf(GLfloat d)                                 { glClearDepthf(d); }
void depthRangef(GLfloat n, GLfloat f)                      { glDepthRangef(n, f); }
void sampleCoverage(GLfloat v, GLboolean invert)            { glSampleCoverage(v, invert); }
*/
import "C"

type Enum uint32

type Attrib struct {
	Value uint
}

type Program struct {
	Value uint32
}

type Shader struct {
	Value uint32
}

type Buffer struct {
	Value uint32
}

type Framebuffer struct {
	Value uint32
}

type Renderbuffer struct {
	Value uint32
}

type Texture struct {
	Value uint32
}

type Uniform struct {
	Value int32
}

var NoAttrib Attrib
var NoProgram Program
var NoShader Shader
var NoBuffer Buffer
var NoFramebuffer Framebuffer
var NoRenderbuffer Renderbuffer
var NoTexture Texture
var NoUniform Uniform

func (v Attrib) c() C.GLuint       { return C.GLuint(v.Value) }
func (v Enum) c() C.GLenum         { return C.GLenum(v) }
func (v Program) c() C.GLuint      { return C.GLuint(v.Value) }
func (v Shader) c() C.GLuint       { return C.GLuint(v.Value) }
func (v Buffer) c() C.GLuint       { return C.GLuint(v.Value) }
func (v Framebuffer) c() C.GLuint  { return C.GLuint(v.Value) }
func (v Renderbuffer) c() C.GLuint { return C.GLuint(v.Value) }
func (v Texture) c() C.GLuint      { return C.GLuint(v.Value) }
func (v Uniform) c() C.GLint       { return C.GLint(v.Value) }

func glBoolean(b bool) C.GLboolean {
	if b {
		return TRUE
	}
	return FALSE
}

// Desktop OpenGL and the ES 2/3 APIs have a very slight difference
// that is imperceptible to C programmers: some function parameters
// use the type Glclampf and some use GLfloat. These two types are
// equivalent in size and bit layout (both are single-precision
// floats), but it plays havoc with cgo. We adjust the types by using
// C wrappers for the problematic functions.

func blendColor(r, g, b, a float32) {
	C.blendColor(C.GLfloat(r), C.GLfloat(g), C.GLfloat(b), C.GLfloat(a))
}
func clearColor(r, g, b, a float32) {
	C.clearColor(C.GLfloat(r), C.GLfloat(g), C.GLfloat(b), C.GLfloat(a))
}
func clearDepthf(d float32)            { C.clearDepthf(C.GLfloat(d)) }
func depthRangef(n, f float32)         { C.depthRangef(C.GLfloat(n), C.GLfloat(f)) }
func sampleCoverage(v float32, i bool) { C.sampleCoverage(C.GLfloat(v), glBoolean(i)) }
