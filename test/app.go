// Package test provides utility drivers for running UI tests without rendering to a screen.
package test // import "fyne.io/fyne/v2/test"

import (
	"net/url"
	"sync"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/theme"
)

// ensure we have a dummy app loaded and ready to test
func init() {
	NewApp()
}

type app struct {
	driver       *driver
	settings     *testSettings
	prefs        fyne.Preferences
	propertyLock sync.RWMutex
	storage      fyne.Storage
	lifecycle    intapp.Lifecycle
	clip         fyne.Clipboard
	cloud        fyne.CloudProvider

	// user action variables
	appliedTheme     fyne.Theme
	lastNotification *fyne.Notification
}

func (a *app) CloudProvider() fyne.CloudProvider {
	return a.cloud
}

func (a *app) Icon() fyne.Resource {
	return nil
}

func (a *app) SetIcon(fyne.Resource) {
	// no-op
}

func (a *app) NewWindow(title string) fyne.Window {
	return a.driver.CreateWindow(title)
}

func (a *app) OpenURL(url *url.URL) error {
	// no-op
	return nil
}

func (a *app) Run() {
	// no-op
}

func (a *app) Quit() {
	// no-op
}

func (a *app) Clipboard() fyne.Clipboard {
	return a.clip
}

func (a *app) UniqueID() string {
	return "testApp" // TODO should this be randomised?
}

func (a *app) Driver() fyne.Driver {
	return a.driver
}

func (a *app) SendNotification(notify *fyne.Notification) {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	a.lastNotification = notify
}

func (a *app) SetCloudProvider(p fyne.CloudProvider) {
	if p == nil {
		a.cloud = nil
		return
	}

	a.transitionCloud(p)
}

func (a *app) Settings() fyne.Settings {
	return a.settings
}

func (a *app) Preferences() fyne.Preferences {
	return a.prefs
}

func (a *app) Storage() fyne.Storage {
	return a.storage
}

func (a *app) Lifecycle() fyne.Lifecycle {
	return &a.lifecycle
}

func (a *app) Metadata() fyne.AppMetadata {
	return fyne.AppMetadata{} // just dummy data
}

func (a *app) lastAppliedTheme() fyne.Theme {
	a.propertyLock.Lock()
	defer a.propertyLock.Unlock()

	return a.appliedTheme
}

func (a *app) transitionCloud(p fyne.CloudProvider) {
	if a.cloud != nil {
		a.cloud.Cleanup(a)
	}

	err := p.Setup(a)
	if err != nil {
		fyne.LogError("Failed to set up cloud provider "+p.ProviderName(), err)
		return
	}
	a.cloud = p

	listeners := a.prefs.ChangeListeners()
	if pp, ok := p.(fyne.CloudProviderPreferences); ok {
		a.prefs = pp.CloudPreferences(a)
	} else {
		a.prefs = internal.NewInMemoryPreferences()
	}
	if store, ok := p.(fyne.CloudProviderStorage); ok {
		a.storage = store.CloudStorage(a)
	} else {
		a.storage = &testStorage{}
	}

	for _, l := range listeners {
		a.prefs.AddChangeListener(l)
		l() // assume that preferences have changed because we replaced the provider
	}

	// after transition ensure settings listener is fired
	a.settings.apply()
}

// NewTempApp returns a new dummy app and tears it down at the end of the test.
//
// Since: 2.5
func NewTempApp(t testing.TB) fyne.App {
	app := NewApp()
	t.Cleanup(func() { NewApp() })
	return app
}

// NewApp returns a new dummy app used for testing.
// It loads a test driver which creates a virtual window in memory for testing.
func NewApp() fyne.App {
	settings := &testSettings{scale: 1.0, theme: Theme()}
	prefs := internal.NewInMemoryPreferences()
	store := &testStorage{}
	test := &app{settings: settings, prefs: prefs, storage: store, driver: NewDriver().(*driver), clip: NewClipboard()}
	settings.app = test
	root, _ := store.docRootURI()
	store.Docs = &internal.Docs{RootDocURI: root}
	painter.ClearFontCache()
	cache.ResetThemeCaches()
	fyne.SetCurrentApp(test)

	return test
}

type testSettings struct {
	primaryColor string
	scale        float32
	theme        fyne.Theme

	listeners       []func(fyne.Settings)
	changeListeners []chan fyne.Settings
	propertyLock    sync.RWMutex
	app             *app
}

func (s *testSettings) AddChangeListener(listener chan fyne.Settings) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *testSettings) AddListener(listener func(fyne.Settings)) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.listeners = append(s.listeners, listener)
}

func (s *testSettings) BuildType() fyne.BuildType {
	return fyne.BuildStandard
}

func (s *testSettings) PrimaryColor() string {
	if s.primaryColor != "" {
		return s.primaryColor
	}

	return theme.ColorBlue
}

func (s *testSettings) SetTheme(theme fyne.Theme) {
	s.propertyLock.Lock()
	s.theme = theme
	s.propertyLock.Unlock()

	s.apply()
}

func (s *testSettings) ShowAnimations() bool {
	return true
}

func (s *testSettings) Theme() fyne.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if s.theme == nil {
		return test.DarkTheme(theme.DefaultTheme())
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
	listenersFns := s.listeners
	s.propertyLock.RUnlock()

	for _, listener := range listeners {
		listener <- s
	}

	s.app.driver.DoFromGoroutine(func() {
		s.app.propertyLock.Lock()
		painter.ClearFontCache()
		cache.ResetThemeCaches()
		intapp.ApplySettings(s, s.app)
		s.app.propertyLock.Unlock()

		for _, l := range listenersFns {
			l(s)
		}
	}, false)

	s.app.propertyLock.Lock()
	s.app.appliedTheme = s.Theme()
	s.app.propertyLock.Unlock()
}
