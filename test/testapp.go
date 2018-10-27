// Package test provides utility drivers for running UI tests without rendering
package test

import "github.com/fyne-io/fyne"

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
	delegate fyne.App
}

func (a *testApp) NewWindow(title string) fyne.Window {
	return a.delegate.NewWindow(title)
}

func (a *testApp) OpenURL(url string) {
	// no-op
}

func (a *testApp) Run() {
	// no-op
}

func (a *testApp) Quit() {
	a.delegate.Quit()
}

// NewApp returns a new dummy app used for testing..
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	realApp := fyne.NewAppWithDriver(NewTestDriver())

	return &testApp{delegate: realApp}
}
