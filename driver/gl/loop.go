// +build !ci,gl

package gl

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
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
			for i, win := range d.windows {
				viewport := win.(*window).viewport
				viewport.MakeContextCurrent()

				canvas := win.(*window).canvas
				gl.UseProgram(canvas.program)

				if viewport.ShouldClose() {
					if win.(*window).master {
						close(d.done)
					}

					// remove window from window list
					d.windows = append(d.windows[:i], d.windows[i+1:]...)
					continue
				}

				canvas.refresh()
				win.(*window).viewport.SwapBuffers()
			}
		}
	}
}
