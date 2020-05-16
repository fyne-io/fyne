package app

import (
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
	propertyLock   sync.RWMutex
	theme          fyne.Theme
	themeSpecified bool

	listenerLock    sync.Mutex
	changeListeners []chan fyne.Settings
	watcher         interface{} // normally *fsnotify.Watcher or nil - avoid import in this file

	schema SettingsSchema
}

func (s *settings) Theme() fyne.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.themeSpecified = true
	s.applyTheme(theme)
}

func (s *settings) applyTheme(theme fyne.Theme) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.theme = theme
	s.apply()
}

func (s *settings) Scale() float32 {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
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

func (s *settings) setupTheme() {
	if s.themeSpecified {
		return
	}
	name := s.schema.ThemeName
	if env := os.Getenv("FYNE_THEME"); env != "" {
		s.themeSpecified = true
		name = env
	}

	if name == "light" {
		s.applyTheme(theme.LightTheme())
	} else if name == "dark" {
		s.applyTheme(theme.DarkTheme())
	} else {
		s.applyTheme(defaultTheme())
	}
}

func loadSettings() *settings {
	s := &settings{}
	s.load()

	return s
}

func (s *settings) fileChanged() {
	s.load()
	s.apply()
}
