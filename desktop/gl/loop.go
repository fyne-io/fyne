// +build !ci,gl

package gl

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func (d *gLDriver) runGL() {
	/*
		select {
		case <-d.done:
			glfw.Terminate()
			return
		}
	*/

	run := true
	for run { // TODO one per window (started on Show() stopped on Hide())
		for _, win := range d.windows {
			if win.(*window).viewport.ShouldClose() && win.(*window).master {
				run = false
				close(d.done)
			}

			win.(*window).canvas.refresh()
			win.(*window).viewport.SwapBuffers()
			glfw.PollEvents()
		}
	}
	glfw.Terminate()
}
