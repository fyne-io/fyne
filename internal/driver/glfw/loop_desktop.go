//go:build !wasm && !test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var glfwInited bool

func (d *gLDriver) initGLFW() {
	if !glfwInited {
		err := glfw.Init()
		if err != nil {
			fyne.LogError("failed to initialise GLFW", err)
			return
		}

		initCursors()
		glfwInited = true
	}

}

func (d *gLDriver) pollEvents() {
	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}
