// Package test provides utility drivers for running UI tests without rendering
package test // import "fyne.io/fyne/test"

import (
	"net/url"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/theme"
)

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
	driver   *testDriver
	settings fyne.Settings
	prefs    fyne.Preferences
}

func (a *testApp) Icon() fyne.Resource {
	return nil
}

func (a *testApp) SetIcon(fyne.Resource) {
	// no-op
}

func (a *testApp) NewWindow(title string) fyne.Window {
	return a.driver.CreateWindow(title)
}

func (a *testApp) OpenURL(url *url.URL) error {
	// no-op
	return nil
}

func (a *testApp) Run() {
	// no-op
}

func (a *testApp) Quit() {
	// no-op
}

func (a *testApp) UniqueID() string {
	return "testApp" // TODO should this be randomised?
}

func (a *testApp) applyThemeTo(content fyne.CanvasObject, canv fyne.Canvas) {
	if content == nil {
		return
	}
	content.Refresh()

	if wid, ok := content.(fyne.Widget); ok {
		for _, o := range wid.CreateRenderer().Objects() {
			a.applyThemeTo(o, canv)
		}
	}
	if c, ok := content.(*fyne.Container); ok {
		for _, o := range c.Objects {
			a.applyThemeTo(o, canv)
		}
	}
}

func (a *testApp) applyTheme() {
	for _, window := range a.driver.AllWindows() {
		content := window.Content()

		a.applyThemeTo(content, window.Canvas())
	}
}

func (a *testApp) Driver() fyne.Driver {
	return a.driver
}

func (a *testApp) Settings() fyne.Settings {
	return a.settings
}

func (a *testApp) Preferences() fyne.Preferences {
	return a.prefs
}

// NewApp returns a new dummy app used for testing.
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	settings := &testSettings{}
	settings.listenerMutex = &sync.Mutex{}
	prefs := internal.NewInMemoryPreferences()
	test := &testApp{settings: settings, prefs: prefs}
	fyne.SetCurrentApp(test)
	test.driver = NewDriver().(*testDriver)

	listener := make(chan fyne.Settings)
	test.Settings().AddChangeListener(listener)
	go func() {
		for {
			_ = <-listener
			test.applyTheme()
		}
	}()

	return test
}

type testSettings struct {
	theme fyne.Theme
	scale float32

	changeListeners []chan fyne.Settings
	listenerMutex   *sync.Mutex
}

func (s *testSettings) AddChangeListener(listener chan fyne.Settings) {
	s.listenerMutex.Lock()
	defer s.listenerMutex.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *testSettings) SetTheme(theme fyne.Theme) {
	s.theme = theme

	s.apply()
}

func (s *testSettings) Theme() fyne.Theme {
	if s.theme == nil {
		return theme.DarkTheme()
	}

	return s.theme
}

func (s *testSettings) Scale() float32 {
	return s.scale
}

func (s *testSettings) apply() {
	s.listenerMutex.Lock()
	defer s.listenerMutex.Unlock()
	for _, listener := range s.changeListeners {
		listener <- s
	}
}
