package app

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SettingsSchema is used for loading and storing global settings
type SettingsSchema struct {
	// these items are used for global settings load
	ThemeName    string  `json:"theme"`
	Scale        float32 `json:"scale"`
	PrimaryColor string  `json:"primary_color"`
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
	variant        fyne.ThemeVariant

	changeListeners sync.Map    // map[chan fyne.Settings]bool
	watcher         interface{} // normally *fsnotify.Watcher or nil - avoid import in this file

	schema SettingsSchema
}

func (s *settings) BuildType() fyne.BuildType {
	return buildMode
}

func (s *settings) PrimaryColor() string {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.schema.PrimaryColor
}

// OverrideTheme allows the settings app to temporarily preview different theme details.
// Please make sure that you remember the original settings and call this again to revert the change.
func (s *settings) OverrideTheme(theme fyne.Theme, name string) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.schema.PrimaryColor = name
	s.theme = theme
}

func (s *settings) Theme() fyne.Theme {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.themeSpecified = true
	s.applyTheme(theme, s.variant)
}

func (s *settings) ThemeVariant() fyne.ThemeVariant {
	return s.variant
}

func (s *settings) applyTheme(theme fyne.Theme, variant fyne.ThemeVariant) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.variant = variant
	s.theme = theme
	s.apply()
}

func (s *settings) Scale() float32 {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	if s.schema.Scale < 0.0 {
		return 1.0 // catching any really old data still using the `-1`  value for "auto" scale
	}
	return s.schema.Scale
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
	s.changeListeners.Store(listener, true) // the boolean is just a dummy value here.
}

func (s *settings) apply() {
	s.changeListeners.Range(func(key, _ interface{}) bool {
		listener := key.(chan fyne.Settings)
		select {
		case listener <- s:
		default:
			l := listener
			go func() { l <- s }()
		}
		return true
	})
}

func (s *settings) fileChanged() {
	s.load()
	s.apply()
}

func (s *settings) loadSystemTheme() fyne.Theme {
	path := filepath.Join(rootConfigDir(), "theme.json")
	data, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fyne.LogError("Failed to load user theme file: "+path, err)
		}
		return theme.DefaultTheme()
	}
	if data != nil && data.Content() != nil {
		th, err := theme.FromJSONReader(bytes.NewReader(data.Content()))
		if err == nil {
			return th
		}
		fyne.LogError("Failed to parse user theme file: "+path, err)
	}
	return theme.DefaultTheme()
}

func (s *settings) setupTheme() {
	name := s.schema.ThemeName
	if env := os.Getenv("FYNE_THEME"); env != "" {
		name = env
	}

	variant := defaultVariant()
	effectiveTheme := s.theme
	if !s.themeSpecified {
		effectiveTheme = s.loadSystemTheme()
	}
	switch name {
	case "light":
		variant = theme.VariantLight
	case "dark":
		variant = theme.VariantDark
	}

	s.applyTheme(effectiveTheme, variant)
}

func loadSettings() *settings {
	s := &settings{}
	s.load()

	return s
}
