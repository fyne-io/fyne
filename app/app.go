// Package app provides app implementations for working with Fyne graphical interfaces.
// The fastest way to get started is to call app.New() which will normally load a new desktop application.
// If the "ci" tag is passed to go (go run -tags ci myapp.go) it will run an in-memory application.
package app // import "fyne.io/fyne/v2/app"

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/app"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

// Declare conformity with App interface
var _ fyne.App = (*fyneApp)(nil)

type fyneApp struct {
	driver    fyne.Driver
	clipboard fyne.Clipboard
	icon      fyne.Resource
	uniqueID  string

	cloud     fyne.CloudProvider
	lifecycle app.Lifecycle
	settings  *settings
	storage   fyne.Storage
	prefs     fyne.Preferences
}

func (a *fyneApp) CloudProvider() fyne.CloudProvider {
	return a.cloud
}

func (a *fyneApp) Icon() fyne.Resource {
	if a.icon != nil {
		return a.icon
	}

	if a.Metadata().Icon == nil || len(a.Metadata().Icon.Content()) == 0 {
		return nil
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
	go a.lifecycle.RunEventQueue(a.driver.DoFromGoroutine)

	if !a.driver.Device().IsMobile() {
		a.settings.watchSettings()
	}

	a.driver.Run()
}

func (a *fyneApp) Quit() {
	for _, window := range a.driver.AllWindows() {
		window.Close()
	}

	a.driver.Quit()
	a.settings.stopWatching()
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
	return &a.lifecycle
}

func (a *fyneApp) newDefaultPreferences() *preferences {
	p := newPreferences(a)
	if a.uniqueID != "" {
		p.load()
	}
	return p
}

func (a *fyneApp) Clipboard() fyne.Clipboard {
	return a.clipboard
}

// New returns a new application instance with the default driver and no unique ID (unless specified in FyneApp.toml)
func New() fyne.App {
	if meta.ID == "" {
		checkLocalMetadata() // if no ID passed, check if it was in toml
		if meta.ID == "" {
			internal.LogHint("Applications should be created with a unique ID using app.NewWithID()")
		}
	}
	return NewWithID(meta.ID)
}

func makeStoreDocs(id string, s *store) *internal.Docs {
	if id == "" {
		return &internal.Docs{} // an empty impl to avoid crashes
	}
	if root := s.a.storageRoot(); root != "" {
		uri, err := storage.ParseURI(root)
		if err != nil {
			uri = storage.NewFileURI(root)
		}

		exists, err := storage.Exists(uri)
		if !exists || err != nil {
			err = storage.CreateListable(uri)
			if err != nil {
				fyne.LogError("Failed to create app storage space", err)
			}
		}

		root, _ := s.docRootURI()
		return &internal.Docs{RootDocURI: root}
	} else {
		return &internal.Docs{} // an empty impl to avoid crashes
	}
}

func newAppWithDriver(d fyne.Driver, clipboard fyne.Clipboard, id string) fyne.App {
	newApp := &fyneApp{uniqueID: id, clipboard: clipboard, driver: d}
	fyne.SetCurrentApp(newApp)

	newApp.prefs = newApp.newDefaultPreferences()
	newApp.lifecycle.InitEventQueue()
	newApp.lifecycle.SetOnStoppedHookExecuted(func() {
		if prefs, ok := newApp.prefs.(*preferences); ok {
			prefs.forceImmediateSave()
		}
	})

	newApp.registerRepositories() // for web this may provide docs / settings
	newApp.settings = loadSettings()
	store := &store{a: newApp}
	store.Docs = makeStoreDocs(id, store)
	newApp.storage = store

	httpHandler := intRepo.NewHTTPRepository()
	repository.Register("http", httpHandler)
	repository.Register("https", httpHandler)
	return newApp
}

// marker interface to pass system tray to supporting drivers
type systrayDriver interface {
	SetSystemTrayMenu(*fyne.Menu)
	SetSystemTrayIcon(fyne.Resource)
	SetSystemTrayWindow(fyne.Window)
}
