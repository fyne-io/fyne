// Package test provides utility drivers for running UI tests without rendering
package test

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
)

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
	driver *testDriver
}

func (a *testApp) Icon() fyne.Resource {
	return theme.FyneLogo()
}

func (a *testApp) SetIcon(fyne.Resource) {
	// no-op
}

func (a *testApp) NewWindow(title string) fyne.Window {
	return NewTestWindow(nil)
}

func (a *testApp) OpenURL(url string) {
	// no-op
}

func (a *testApp) Run() {
	// no-op
}

func (a *testApp) Quit() {
	// no-op
}

func (a *testApp) applyTheme(fyne.Settings) {
	for _, window := range a.driver.AllWindows() {
		content := window.Content()

		switch themed := content.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
			window.Canvas().Refresh(content)
		}
	}
}

func (a *testApp) Driver() fyne.Driver {
	return a.driver
}

// NewApp returns a new dummy app used for testing..
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	test := &testApp{driver: NewTestDriver().(*testDriver)}
	fyne.SetCurrentApp(test)

	listener := make(chan fyne.Settings)
	fyne.GetSettings().AddChangeListener(listener)
	go func() {
		for {
			settings := <-listener
			test.applyTheme(settings)
		}
	}()

	return test
}
