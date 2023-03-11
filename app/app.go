// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/v2/app"

import (
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"

	"golang.org/x/sys/execabs"
)

// Declare conformity with App interface
var _ fyne.App = (*fyneApp)(nil)

type fyneApp struct {
	driver   fyne.Driver
	icon     fyne.Resource
	uniqueID string

	cloud     fyne.CloudProvider
	lifecycle fyne.Lifecycle
	settings  *settings
	storage   fyne.Storage
	prefs     fyne.Preferences

	running uint32 // atomic, 1 == running, 0 == stopped
	exec    func(name string, arg ...string) *execabs.Cmd
}

func (a *fyneApp) CloudProvider() fyne.CloudProvider {
	return a.cloud
}

func (a *fyneApp) Icon() fyne.Resource {
	if a.icon != nil {
		return a.icon
	}

	return a.Metadata().Icon
}

func (a *fyneApp) SetIcon(icon fyne.Resource) {
	a.icon = icon
}

func (a *fyneApp) UniqueID() string {
	if a.uniqueID != "" {
		return a.uniqueID
	}
	if a.Metadata().ID != "" {
		return a.Metadata().ID
	}

	fyne.LogError("Preferences API requires a unique ID, use app.NewWithID() or the FyneApp.toml ID field", nil)
	a.uniqueID = "missing-id-" + strconv.FormatInt(time.Now().Unix(), 10) // This is a fake unique - it just has to not be reused...
	return a.uniqueID
}

func (a *fyneApp) NewWindow(title string) fyne.Window {
	return a.driver.CreateWindow(title)
}

func (a *fyneApp) Run() {
	if atomic.CompareAndSwapUint32(&a.running, 0, 1) {
		a.driver.Run()
		return
	}
}

func (a *fyneApp) Quit() {
	for _, window := range a.driver.AllWindows() {
		window.Close()
	}

	a.driver.Quit()
	a.settings.stopWatching()
	atomic.StoreUint32(&a.running, 0)
}

func (a *fyneApp) Driver() fyne.Driver {
	return a.driver
}

// Settings returns the application settings currently configured.
func (a *fyneApp) Settings() fyne.Settings {
	return a.settings
}

func (a *fyneApp) Storage() fyne.Storage {
	return a.storage
}

func (a *fyneApp) Preferences() fyne.Preferences {
	if a.UniqueID() == "" {
		fyne.LogError("Preferences API requires a unique ID, use app.NewWithID() or the FyneApp.toml ID field", nil)
	}
	return a.prefs
}

func (a *fyneApp) Lifecycle() fyne.Lifecycle {
	return a.lifecycle
}

func (a *fyneApp) newDefaultPreferences() fyne.Preferences {
	p := fyne.Preferences(newPreferences(a))
	if pref, ok := p.(interface{ load() }); ok && a.uniqueID != "" {
		pref.load()
	}
	return p
}

// New returns a new application instance with the default driver and no unique ID (unless specified in FyneApp.toml)
func New() fyne.App {
	if meta.ID == "" {
		internal.LogHint("Applications should be created with a unique ID using app.NewWithID()")
	}
	return NewWithID(meta.ID)
}

func makeStoreDocs(id string, p fyne.Preferences, s *store) *internal.Docs {
	if id != "" {
		if pref, ok := p.(interface{ load() }); ok {
			pref.load()
		}
		err := os.MkdirAll(s.a.storageRoot(), 0755) // make the space before anyone can use it
		if err != nil {
			fyne.LogError("Failed to create app storage space", err)
		}

		root, _ := s.docRootURI()
		return &internal.Docs{RootDocURI: root}
	} else {
		return &internal.Docs{} // an empty impl to avoid crashes
	}
}

func newAppWithDriver(d fyne.Driver, id string) fyne.App {
	newApp := &fyneApp{uniqueID: id, driver: d, exec: execabs.Command, lifecycle: &app.Lifecycle{}}
	fyne.SetCurrentApp(newApp)

	newApp.prefs = newApp.newDefaultPreferences()
	newApp.settings = loadSettings()
	store := &store{a: newApp}
	store.Docs = makeStoreDocs(id, newApp.prefs, store)
	newApp.storage = store

	if !d.Device().IsMobile() {
		newApp.settings.watchSettings()
	}

	repository.Register("http", intRepo.NewHTTPRepository())
	repository.Register("https", intRepo.NewHTTPRepository())

	return newApp
}

// marker interface to pass system tray to supporting drivers
type systrayDriver interface {
	SetSystemTrayMenu(*fyne.Menu)
	SetSystemTrayIcon(resource fyne.Resource)
}
