// +build !ci,gl

package gl

import (
	"github.com/fyne-io/fyne"
)

type gLDriver struct {
	windows []fyne.Window
	program uint32
	done    chan interface{}
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return fyne.NewSize(len(text)*10, size+4)
}

func (d *gLDriver) Quit() {
	close(d.done)
}

func (d *gLDriver) Run() {
	d.runGL()
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	driver := new(gLDriver)
	driver.done = make(chan interface{})

	return driver
}
