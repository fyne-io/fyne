//go:build js || wasm || web
// +build js wasm web

package glfw

import (
	"fmt"

	"fyne.io/fyne/v2"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
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

func (d *gLDriver) tryPollEvents() {
	defer func() {
		if r := recover(); r != nil {
			fyne.LogError(fmt.Sprint("GLFW poll event error: ", r), nil)
		}
	}()

	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}
