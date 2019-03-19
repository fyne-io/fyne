package app

import "fyne.io/fyne"

type settings struct {
	theme fyne.Theme

	changeListeners []chan fyne.Settings
}

func (s *settings) Theme() fyne.Theme {
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.theme = theme
	s.apply()
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *settings) apply() {
	for _, listener := range s.changeListeners {
		go func(listener chan fyne.Settings) {
			listener <- s
		}(listener)
	}
}

func loadSettings() *settings {
	s := &settings{}

	return s
}
