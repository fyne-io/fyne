package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// SettingsSchema is used for loading and storing global settings
type SettingsSchema struct {
	// these items are used for global settings load
	ThemeName string  `json:"theme"`
	Scale     float32 `json:"scale"`
}

// StoragePath returns the location of the settings storage
func (sc *SettingsSchema) StoragePath() string {
	return filepath.Join(rootConfigDir(), "settings.json")
}

// Declare conformity with Settings interface
var _ fyne.Settings = (*settings)(nil)

type settings struct {
	themeLock sync.RWMutex
	theme     fyne.Theme

	listenerLock    sync.Mutex
	changeListeners []chan fyne.Settings
	watcher         interface{} // normally *fsnotify.Watcher or nil - avoid import in this file

	schema SettingsSchema
}

func (s *settings) Theme() fyne.Theme {
	s.themeLock.RLock()
	defer s.themeLock.RUnlock()
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.themeLock.Lock()
	defer s.themeLock.Unlock()
	s.theme = theme
	s.apply()
}

func (s *settings) Scale() float32 {
	s.themeLock.RLock()
	defer s.themeLock.RUnlock()
	return s.schema.Scale
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *settings) apply() {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()

	for _, listener := range s.changeListeners {
		select {
		case listener <- s:
		default:
			l := listener
			go func() { l <- s }()
		}
	}
}

func (s *settings) load() {
	err := s.loadFromFile(s.schema.StoragePath())
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}

	s.setupTheme()
}

func (s *settings) loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.schema)
}

func (s *settings) fileChanged() {
	s.load()
	s.apply()
}

func (s *settings) setupTheme() {
	name := s.schema.ThemeName
	if env := os.Getenv("FYNE_THEME"); env != "" {
		name = env
	}

	if name == "light" {
		s.SetTheme(theme.LightTheme())
	} else if name == "dark" {
		s.SetTheme(theme.DarkTheme())
	} else {
		s.SetTheme(defaultTheme())
	}
}

func loadSettings() *settings {
	s := &settings{}
	s.load()

	return s
}
