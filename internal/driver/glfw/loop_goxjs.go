//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"

	"github.com/fyne-io/gl-js"
	"github.com/fyne-io/glfw-js"
)

func (d *gLDriver) initGLFW() {
	initOnce.Do(func() {
		err := glfw.Init(gl.ContextWatcher)
		if err != nil {
			fyne.LogError("failed to initialise GLFW", err)
			return
		}

		d.startDrawThread()
	})
}

func (d *gLDriver) pollEvents() {
	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}
