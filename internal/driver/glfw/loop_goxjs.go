//go:build wasm || test_web_driver
// +build wasm test_web_driver

package glfw

import (
	"fmt"

	"fyne.io/fyne/v2"

	gl "github.com/fyne-io/gl-js"
	glfw "github.com/fyne-io/glfw-js"
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

func (*gLDriver) waitForEvents() {
	defer func() {
		// See https://github.com/glfw/glfw/issues/1785 and https://github.com/fyne-io/fyne/issues/1024.
		if r := recover(); r != nil {
			fyne.LogError(fmt.Sprint("GLFW poll event error: ", r), nil)
		}
	}()

	glfw.WaitEvents()
}

func postEmptyEvent() {
	glfw.PostEmptyEvent()
}

func (*gLDriver) terminate() {
	glfw.Terminate()
}
