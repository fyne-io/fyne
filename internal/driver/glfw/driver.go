// Package glfw provides a full Fyne desktop driver that uses the system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package glfw

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

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
	drawDone   chan interface{}

	animations []*anim
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

func (d *gLDriver) StartAnimation(a *fyne.Animation) {
	wasStopped := len(d.animations) == 0

	d.animations = append(d.animations, &anim{a, time.Now(), time.Now().Add(a.Duration)})
	if wasStopped {
		d.runAnimations()
	}
}

func (d *gLDriver) addWindow(w *window) {
	d.windowLock.Lock()
	defer d.windowLock.Unlock()
	d.windows = append(d.windows, w)
}

// a trivial implementation of "focus previous" - return to the most recently opened, or master if set.
// This may not do the right thing if your app has 3 or more windows open, but it was agreed this was not much
// of an issue, and the added complexity to track focus was not needed at this time.
func (d *gLDriver) focusPreviousWindow() {
	d.windowLock.RLock()
	wins := d.windows
	d.windowLock.RUnlock()

	var chosen fyne.Window
	for _, w := range wins {
		chosen = w
		if w.(*window).master {
			break
		}
	}

	if chosen == nil || chosen.(*window).view() == nil {
		return
	}
	chosen.RequestFocus()
}

type anim struct {
	a     *fyne.Animation
	start time.Time
	end   time.Time
}

func (d *gLDriver) runAnimations() {
	draw := time.NewTicker(time.Second / 60)

	go func() {
		for len(d.animations) > 0 {

			<-draw.C
			for i := len(d.animations) - 1; i >= 0; i-- { // backwards so we can remove safely
				a := d.animations[i]

				if !d.tickAnimation(a) {
					if i == len(d.animations)-1 {
						d.animations = d.animations[:len(d.animations)-1]
					} else {
						d.animations = append(d.animations[:i], d.animations[i+1])
					}
				}
			}
		}
	}()
}

func (d *gLDriver) tickAnimation(a *anim) bool {
	if time.Now().After(a.end) {
		a.a.Tick(1.0)
		if !a.a.Repeat {
			return false
		}

		a.start = time.Now()
		a.end = a.start.Add(a.a.Duration)
	}

	total := a.end.Sub(a.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12
	delta := time.Since(a.start).Nanoseconds() / 1000000

	val := float32(delta) / float32(total)
	a.a.Tick(val)

	return true
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
	d.drawDone = make(chan interface{})

	return d
}
