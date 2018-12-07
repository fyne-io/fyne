package fyne

import "os"

// Settings describes the system configuration available and allows configurable
// items to be changed.
type Settings interface {
	Theme() string
	SetTheme(string)

	AddChangeListener(chan Settings)
}

type settings struct {
	theme string

	changeListeners []chan Settings
}

var settingsCache *settings

func (s *settings) Theme() string {
	return s.theme
}

func (s *settings) SetTheme(theme string) {
	s.theme = theme
	s.apply()
}

func (s *settings) AddChangeListener(listener chan Settings) {
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *settings) apply() {
	for _, listener := range s.changeListeners {
		go func(listener chan Settings) {
			listener <- s
		}(listener)
	}
}

func loadSettings() *settings {
	s := &settings{}

	env := os.Getenv("FYNE_THEME")
	if env == "light" {
		s.theme = env
	} else {
		s.theme = "dark"
	}

	return s
}

// GlobalSettings returns the system wide settings currently configured
func GlobalSettings() Settings {
	if settingsCache == nil {
		settingsCache = loadSettings()
	}

	return settingsCache
}
