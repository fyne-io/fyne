// +build !ci,gl

package gl

import (
	"runtime"
	"time"

	"github.com/fyne-io/fyne"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func (d *gLDriver) runGL() {
	fps := time.NewTicker(time.Second / 60)

	for {
		select {
		case <-d.done:
			fps.Stop()
			glfw.Terminate()
			return
		case object := <-refreshQueue:
			freeWalked := func(obj fyne.CanvasObject, _ fyne.Position) {
				texture := textures[obj]
				if texture != 0 {
					gl.DeleteTextures(1, &texture)
					delete(textures, obj)
				}
			}

			walkObjects(object, fyne.NewPos(0, 0), freeWalked)
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
					glfw.DetachCurrentContext()
					continue
				}

				if !canvas.isDirty() {
					glfw.DetachCurrentContext()
					continue
				}
				win.(*window).fitContent()

				size := canvas.Size()
				winWidth := scaleInt(canvas, size.Width)
				winHeight := scaleInt(canvas, size.Height)

				gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
				canvas.paint(size)

				win.(*window).viewport.SwapBuffers()
				glfw.DetachCurrentContext()
			}
		}
	}
}
