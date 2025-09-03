//go:build !wasm && !test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func (d *gLDriver) initGLFW() {
	err := glfw.Init()
	if err != nil {
		fyne.LogError("failed to initialise GLFW", err)
		return
	}

	initCursors()
}

func (d *gLDriver) pollEvents() {
	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}
