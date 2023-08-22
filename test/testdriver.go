package test

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/software"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"
)

// SoftwarePainter describes a simple type that can render canvases
type SoftwarePainter interface {
	Paint(fyne.Canvas) image.Image
}

type testDriver struct {
	device       *device
	painter      SoftwarePainter
	windows      []fyne.Window
	windowsMutex sync.RWMutex
}

// Declare conformity with Driver
var _ fyne.Driver = (*testDriver)(nil)

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() fyne.Driver {
	drv := &testDriver{windowsMutex: sync.RWMutex{}}
	repository.Register("file", intRepo.NewFileRepository())

	// make a single dummy window for rendering tests
	drv.CreateWindow("")

	return drv
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) fyne.Driver {
	return &testDriver{
		painter:      painter,
		windowsMutex: sync.RWMutex{},
	}
}

func (d *testDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	tc := c.(*testCanvas)
	return driver.AbsolutePositionForObject(co, tc.objectTrees())
}

func (d *testDriver) AllWindows() []fyne.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

func (d *testDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating: probably the last created window is meant
	return d.windows[len(d.windows)-1].Canvas()
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	canvas := NewCanvas().(*testCanvas)
	if d.painter != nil {
		canvas.painter = d.painter
	} else {
		canvas.painter = software.NewPainter()
	}

	window := &testWindow{canvas: canvas, driver: d}
	window.clipboard = &testClipboard{}

	d.windowsMutex.Lock()
	d.windows = append(d.windows, window)
	d.windowsMutex.Unlock()
	return window
}

func (d *testDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}
	return d.device
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *testDriver) RenderedTextSize(text string, size float32, style fyne.TextStyle) (fyne.Size, float32) {
	return painter.RenderedTextSize(text, size, style)
}

func (d *testDriver) Run() {
	// no-op
}

func (d *testDriver) StartAnimation(a *fyne.Animation) {
	// currently no animations in test app, we just initialise it and leave
	a.Tick(1.0)
}

func (d *testDriver) StopAnimation(a *fyne.Animation) {
	// currently no animations in test app, do nothing
}

func (d *testDriver) Quit() {
	// no-op
}

func (d *testDriver) removeWindow(w *testWindow) {
	d.windowsMutex.Lock()
	i := 0
	for _, window := range d.windows {
		if window == w {
			break
		}
		i++
	}

	d.windows = append(d.windows[:i], d.windows[i+1:]...)
	d.windowsMutex.Unlock()
}
