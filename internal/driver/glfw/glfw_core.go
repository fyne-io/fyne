// +build !gles,!arm,!arm64 darwin

package glfw

import "github.com/go-gl/glfw/v3.3/glfw"

func initWindowHints() {
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
}
