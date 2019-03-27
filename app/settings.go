package app

import (
	"fyne.io/fyne"
	"sync"
)

type settings struct {
	sync.RWMutex
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
	s.Lock()
	defer s.Unlock()
	s.changeListeners = append(s.changeListeners, listener)
}

func (s *settings) apply() {
	s.RLock()
	defer s.RUnlock()
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
