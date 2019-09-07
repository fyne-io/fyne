package test

import (
	"image"

	"fyne.io/fyne"
)

// SoftwarePainter describes a simple type that can render canvases
type SoftwarePainter interface {
	Paint(fyne.Canvas) image.Image
}

type testDriver struct {
	device  *device
	painter SoftwarePainter
}

// Declare conformity with Driver
var _ fyne.Driver = (*testDriver)(nil)

func (d *testDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	windowsMutex.RLock()
	defer windowsMutex.RUnlock()
	// cheating as we only have a single test window
	return windows[0].Canvas()
}

func (d *testDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	return co.Position() // TODO get the real answer
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	return NewWindow(nil)
}

func (d *testDriver) AllWindows() []fyne.Window {
	windowsMutex.RLock()
	defer windowsMutex.RUnlock()
	return windows
}

func (d *testDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	// The real text height differs from the requested text size.
	// We simulate this behaviour here.
	return fyne.NewSize(len(text)*size, size*13/10+1)
}

func (d *testDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}
	return d.device
}

func (d *testDriver) Run() {
	// no-op
}

func (d *testDriver) Quit() {
	// no-op
}

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver() fyne.Driver {
	driver := new(testDriver)
	// make a single dummy window for rendering tests
	NewWindow(nil)

	return driver
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter SoftwarePainter) fyne.Driver {
	driver := new(testDriver)
	driver.painter = painter

	return driver
}
