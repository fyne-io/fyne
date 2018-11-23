// +build !ci,gl

package gl

import (
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func (d *gLDriver) runGL() {
	fps := time.NewTicker(time.Second / 6) /// 60)

	for {
		select {
		case <-d.done:
			fps.Stop()
			glfw.Terminate()
			return
		case <-fps.C:
			glfw.PollEvents()
			for _, win := range d.windows { // TODO per window?
				if win.(*window).viewport.ShouldClose() && win.(*window).master {
					close(d.done)
					continue
				}

				win.(*window).canvas.refresh()
				win.(*window).viewport.SwapBuffers()
			}
		}
	}
}
