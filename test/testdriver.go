package test

import "fyne.io/fyne"

type testDriver struct {
}

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
