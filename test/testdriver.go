package test

import "github.com/fyne-io/fyne/ui"

type testDriver struct {
}

func (d *testDriver) CreateWindow(string) ui.Window {
	return &testWindow{}
}

func (d *testDriver) AllWindows() []ui.Window {
	return windows
}

// NewTestDriver sets up and registers a new dummy driver for test purpose
func NewTestDriver() ui.Driver {
	driver := new(testDriver)

	ui.SetDriver(driver)
	return driver
}
