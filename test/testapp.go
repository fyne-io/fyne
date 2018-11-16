// Package test provides utility drivers for running UI tests without rendering
package test

import "github.com/fyne-io/fyne"

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
}

func (a *testApp) NewWindow(title string) fyne.Window {
	return &testWindow{title: title}
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
	for _, window := range fyne.GetDriver().AllWindows() {
		content := window.Content()

		switch themed := content.(type) {
		case fyne.ThemedObject:
			themed.ApplyTheme()
			window.Canvas().Refresh(content)
		}
	}
}

// NewApp returns a new dummy app used for testing..
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	test := &testApp{}
	fyne.SetDriver(NewTestDriver())

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
