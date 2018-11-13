// +build !ci,gl

package gl

import (
	"github.com/fyne-io/fyne"
)

type gLDriver struct {
	windows []fyne.Window
	done    chan interface{}
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return fyne.NewSize(len(text)*5, size)
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
