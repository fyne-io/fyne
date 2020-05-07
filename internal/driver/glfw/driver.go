// Package glfw provides a full Fyne desktop driver that uses the system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package glfw

import (
	"runtime"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
)

const mainGoroutineID = 1

var canvasMutex sync.RWMutex
var canvases = make(map[fyne.CanvasObject]fyne.Canvas)

// Declare conformity with Driver
var _ fyne.Driver = (*gLDriver)(nil)

type gLDriver struct {
	windowLock sync.RWMutex
	windows    []fyne.Window
	device     *glDevice
	done       chan interface{}
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
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	glc := c.(*glCanvas)
	return driver.AbsolutePositionForObject(co, glc.objectTrees())
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
	if goroutineID() != mainGoroutineID {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}
	d.runGL()
}

func (d *gLDriver) addWindow(w *window) {
	d.windowLock.Lock()
	defer d.windowLock.Unlock()
	d.windows = append(d.windows, w)
}

func (d *gLDriver) windowList() []fyne.Window {
	d.windowLock.RLock()
	defer d.windowLock.RUnlock()
	return d.windows
}

func goroutineID() int {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// string format expects "goroutine X [running..."
	id := strings.Split(strings.TrimSpace(string(b)), " ")[1]

	num, _ := strconv.Atoi(id)
	return num
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	d := new(gLDriver)
	d.done = make(chan interface{})

	return d
}
