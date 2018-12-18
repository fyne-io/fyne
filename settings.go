package fyne

// Settings describes the system configuration available and allows configurable
// items to be changed.
type Settings interface {
	Theme() Theme
	SetTheme(Theme)

	AddChangeListener(chan Settings)
}

type settings struct {
	theme Theme

	changeListeners []chan Settings
}

var settingsCache *settings

func (s *settings) Theme() Theme {
	return s.theme
}

func (s *settings) SetTheme(theme Theme) {
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

	return s
}

// GlobalSettings returns the system wide settings currently configured
func GlobalSettings() Settings {
	if settingsCache == nil {
		settingsCache = loadSettings()
	}

	return settingsCache
}
