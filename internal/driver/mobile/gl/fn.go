// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "unsafe"

type call struct {
	args     fnargs
	parg     unsafe.Pointer
	blocking bool
}

type fnargs struct {
	fn glfn

	a0 uintptr
	a1 uintptr
	a2 uintptr
	a3 uintptr
	a4 uintptr
	a5 uintptr
	a6 uintptr
	a7 uintptr
	a8 uintptr
	a9 uintptr
}

type glfn int

const (
	glfnUNDEFINED glfn = iota
	glfnActiveTexture
	glfnAttachShader
	glfnBindBuffer
	glfnBindTexture
	glfnBindVertexArray
	glfnBlendColor
	glfnBlendFunc
	glfnBufferData
	glfnClear
	glfnClearColor
	glfnCompileShader
	glfnCreateProgram
	glfnCreateShader
	glfnDeleteBuffer
	glfnDeleteTexture
	glfnDisable
	glfnDrawArrays
	glfnEnable
	glfnEnableVertexAttribArray
	glfnFlush
	glfnGenBuffer
	glfnGenTexture
	glfnGenVertexArray
	glfnGetAttribLocation
	glfnGetError
	glfnGetProgramInfoLog
	glfnGetProgramiv
	glfnGetShaderInfoLog
	glfnGetShaderSource
	glfnGetShaderiv
	glfnGetTexParameteriv
	glfnGetUniformLocation
	glfnLinkProgram
	glfnReadPixels
	glfnScissor
	glfnShaderSource
	glfnTexImage2D
	glfnTexParameteri
	glfnUniform1f
	glfnUniform2f
	glfnUniform4f
	glfnUniform4fv
	glfnUseProgram
	glfnVertexAttribPointer
	glfnViewport
)

func goString(buf []byte) string {
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	panic("buf is not NUL-terminated")
}

func glBoolean(b bool) uintptr {
	if b {
		return True
	}
	return False
}
