//go:build wasm || test_web_driver

package glfw

import glfw "github.com/fyne-io/glfw-js"

func initWindowHints() {
	// Request an alpha channel so the WebGL canvas framebuffer is RGBA.
	glfw.WindowHint(glfw.AlphaBits, 8)
}
