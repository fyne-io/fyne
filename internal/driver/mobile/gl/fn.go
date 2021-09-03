//go:build android || ios
// +build android ios

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
	glfnGetAttribLocation
	glfnGetError
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
		return TRUE
	}
	return FALSE
}
