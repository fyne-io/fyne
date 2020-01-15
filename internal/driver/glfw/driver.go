// Package glfw provides a full Fyne desktop driver that uses the system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package glfw

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
)

var canvasMutex sync.RWMutex
var canvases = make(map[fyne.CanvasObject]fyne.Canvas)

// Declare conformity with Driver
var _ fyne.Driver = (*gLDriver)(nil)

type gLDriver struct {
	windows []fyne.Window
	device  *glDevice
	done    chan interface{}
}

func (d *gLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return painter.RenderedTextSize(text, size, style)
}

func (d *gLDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	canvasMutex.RLock()
	defer canvasMutex.RUnlock()
	return canvases[obj]
}

func (d *gLDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	var pos fyne.Position
	c := fyne.CurrentApp().Driver().CanvasForObject(co).(*glCanvas)

	driver.WalkCompleteObjectTree(c.content, func(o fyne.CanvasObject, p fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if o == co {
			pos = p
			return true
		}
		return false
	}, nil)

	return pos
}

func (d *gLDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &glDevice{}
	}

	return d.device
}

func (d *gLDriver) Quit() {
	defer func() {
		recover() // we could be called twice - no safe way to check if d.done is closed
	}()
	close(d.done)
}

func (d *gLDriver) Run() {
	d.runGL()
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	d := new(gLDriver)
	d.done = make(chan interface{})

	return d
}
