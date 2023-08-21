//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"fmt"

	"fyne.io/fyne/v2"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func (d *gLDriver) initGLFW() {
	initOnce.Do(func() {
		err := glfw.Init()
		if err != nil {
			fyne.LogError("failed to initialise GLFW", err)
			return
		}

		initCursors()
		d.startDrawThread()
	})
}

func (d *gLDriver) waitForEvents() {
	defer func() {
		if r := recover(); r != nil {
			fyne.LogError(fmt.Sprint("GLFW poll event error: ", r), nil)
		}
	}()

	glfw.WaitEvents()
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}

// PostEmptyEvent posts an event to the GLFW event queue.
// This is used to signal execution to continue from glfw.WaitEvents().
func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}
