package test

import "github.com/fyne-io/fyne/api/ui"

type testDriver struct {
}

func (d *testDriver) CreateWindow(string) ui.Window {
	return &testWindow{}
}

func (d *testDriver) AllWindows() []ui.Window {
	return windows
}

func (d *testDriver) RenderedTextSize(text string, size int) ui.Size {
	return ui.NewSize(len(text)*size, size)
}

func (d *testDriver) Quit() {
	// no-op
}

// NewTestDriver sets up and registers a new dummy driver for test purpose
func NewTestDriver() ui.Driver {
	driver := new(testDriver)
	// make a single dummy window for rendering tests
	NewTestWindow()

	ui.SetDriver(driver)
	return driver
}
