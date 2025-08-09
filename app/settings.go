package app

import (
	"bytes"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/theme"
)

// SettingsSchema is used for loading and storing global settings
type SettingsSchema struct {
	// these items are used for global settings load
	ThemeName         string  `json:"theme"`
	Scale             float32 `json:"scale"`
	PrimaryColor      string  `json:"primary_color"`
	CloudName         string  `json:"cloud_name"`
	CloudConfig       string  `json:"cloud_config"`
	DisableAnimations bool    `json:"no_animations"`
}

// StoragePath returns the location of the settings storage
func (sc *SettingsSchema) StoragePath() string {
	return filepath.Join(app.RootConfigDir(), "settings.json")
}

// Declare conformity with Settings interface
var _ fyne.Settings = (*settings)(nil)

type settings struct {
	theme          fyne.Theme
	themeSpecified bool
	variant        fyne.ThemeVariant

	listeners       []func(fyne.Settings)
	changeListeners async.Map[chan fyne.Settings, bool]
	watcher         any // normally *fsnotify.Watcher or nil - avoid import in this file

	schema SettingsSchema
}

func (s *settings) BuildType() fyne.BuildType {
	return build.Mode
}

func (s *settings) PrimaryColor() string {
	return s.schema.PrimaryColor
}

// OverrideTheme allows the settings app to temporarily preview different theme details.
// Please make sure that you remember the original settings and call this again to revert the change.
//
// Deprecated: Use container.NewThemeOverride to change the appearance of part of your application.
func (s *settings) OverrideTheme(theme fyne.Theme, name string) {
	s.schema.PrimaryColor = name
	s.theme = theme
}

func (s *settings) Theme() fyne.Theme {
	if s == nil {
		fyne.LogError("Attempt to access current Fyne theme when no app is started", nil)
		return nil
	}
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.themeSpecified = true
	s.applyTheme(theme, s.variant)
}

func (s *settings) ShowAnimations() bool {
	return !s.schema.DisableAnimations && !build.NoAnimations
}

func (s *settings) ThemeVariant() fyne.ThemeVariant {
	return s.variant
}

func (s *settings) applyTheme(theme fyne.Theme, variant fyne.ThemeVariant) {
	s.variant = variant
	s.theme = theme
	s.apply()
}

func (s *settings) applyVariant(variant fyne.ThemeVariant) {
	s.variant = variant
	s.apply()
}

func (s *settings) Scale() float32 {
	if s.schema.Scale < 0.0 {
		return 1.0 // catching any really old data still using the `-1`  value for "auto" scale
	}
	return s.schema.Scale
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
	s.changeListeners.Store(listener, true) // the boolean is just a dummy value here.
}

func (s *settings) AddListener(listener func(fyne.Settings)) {
	s.listeners = append(s.listeners, listener)
}

func (s *settings) apply() {
	s.changeListeners.Range(func(listener chan fyne.Settings, _ bool) bool {
		select {
		case listener <- s:
		default:
			l := listener
			go func() { l <- s }()
		}
		return true
	})

	for _, l := range s.listeners {
		l(s)
	}
}

func (s *settings) fileChanged() {
	s.load()
	s.apply()
}

func (s *settings) loadSystemTheme() fyne.Theme {
	path := filepath.Join(app.RootConfigDir(), "theme.json")
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

	variant := app.DefaultVariant()
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
