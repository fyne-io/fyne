// Package test provides utility drivers for running UI tests without rendering
package test

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/api/ui"

// ensure we have a dummy app loaded and ready to test
func init() {
	NewTestApp()
}

type testApp struct {
}

func (a *testApp) NewWindow(title string) ui.Window {
	return NewTestWindow()
}

func (a *testApp) OpenURL(url string) {
}

func (a *testApp) Quit() {
}

// NewTestApp returns a new dummy app used for testing..
// It laods a test driver and creates a single window in memory for testing UI.
func NewTestApp() fyne.App {
	NewTestDriver()
	NewTestWindow()

	return &testApp{}
}
