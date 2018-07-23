package test

import "github.com/fyne-io/fyne"

type testDriver struct {
}

func (d *testDriver) CreateWindow(string) fyne.Window {
	return NewTestWindow(nil)
}

func (d *testDriver) AllWindows() []fyne.Window {
	return windows
}

func (d *testDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	return fyne.NewSize(len(text)*size, size)
}

func (d *testDriver) Quit() {
	// no-op
}

// NewTestDriver sets up and registers a new dummy driver for test purpose
func NewTestDriver() fyne.Driver {
	driver := new(testDriver)
	// make a single dummy window for rendering tests
	NewTestWindow(nil)

	return driver
}
