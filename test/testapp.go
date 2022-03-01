// Package test provides utility drivers for running UI tests without rendering to a screen.
package test // import "fyne.io/fyne/v2/test"

import (
	"net/url"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/theme"
)

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type testApp struct {
	driver       *testDriver
	settings     fyne.Settings
	prefs        fyne.Preferences
	propertyLock sync.RWMutex
	storage      fyne.Storage
	lifecycle    fyne.Lifecycle

	// user action variables
	appliedTheme     fyne.Theme
	lastNotification *fyne.Notification
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

func (a *testApp) Driver() fyne.Driver {
	return a.driver
}

func (a *testApp) SendNotification(notify *fyne.Notification) {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	a.lastNotification = notify
}

func (a *testApp) Settings() fyne.Settings {
	return a.settings
}

func (a *testApp) Preferences() fyne.Preferences {
	return a.prefs
}

func (a *testApp) Storage() fyne.Storage {
	return a.storage
}

func (a *testApp) Lifecycle() fyne.Lifecycle {
	return a.lifecycle
}

func (a *testApp) Metadata() fyne.AppMetadata {
	return fyne.AppMetadata{} // just dummy data
}

func (a *testApp) lastAppliedTheme() fyne.Theme {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	return a.appliedTheme
}

// NewApp returns a new dummy app used for testing.
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	settings := &testSettings{scale: 1.0, theme: Theme()}
	prefs := internal.NewInMemoryPreferences()
	store := &testStorage{}
	test := &testApp{settings: settings, prefs: prefs, storage: store, driver: NewDriver().(*testDriver),
		lifecycle: &app.Lifecycle{}}
	root, _ := store.docRootURI()
	store.Docs = &internal.Docs{RootDocURI: root}
	cache.ResetThemeCaches()
	fyne.SetCurrentApp(test)

	listener := make(chan fyne.Settings)
	test.Settings().AddChangeListener(listener)
	go func() {
		for {
			<-listener
			cache.ResetThemeCaches()
			app.ApplySettings(test.Settings(), test)

			test.propertyLock.Lock()
			test.appliedTheme = test.Settings().Theme()
			test.propertyLock.Unlock()
		}
	}()

	return test
}

type testSettings struct {
	theme fyne.Theme
	scale float32

	changeListeners []chan fyne.Settings
	propertyLock    sync.RWMutex
}

func (s *testSettings) AddChangeListener(listener chan fyne.Settings) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *testSettings) BuildType() fyne.BuildType {
	return fyne.BuildStandard
}

func (s *testSettings) PrimaryColor() string {
	return theme.ColorBlue
}

func (s *testSettings) SetTheme(theme fyne.Theme) {
	s.propertyLock.Lock()
	s.theme = theme
	s.propertyLock.Unlock()

	s.apply()
}

func (s *testSettings) Theme() fyne.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if s.theme == nil {
		return theme.DarkTheme()
	}

	return s.theme
}

func (s *testSettings) ThemeVariant() fyne.ThemeVariant {
	return 2 // not a preference
}

func (s *testSettings) Scale() float32 {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.scale
}

func (s *testSettings) apply() {
	s.propertyLock.RLock()
	listeners := s.changeListeners
	s.propertyLock.RUnlock()

	for _, listener := range listeners {
		listener <- s
	}
}
