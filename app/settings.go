package app

import (
	"sync"

	"fyne.io/fyne"
)

type settings struct {
	theme fyne.Theme

	changeListeners []chan fyne.Settings
	listenerMutex   *sync.RWMutex
}

func (s *settings) Theme() fyne.Theme {
	return s.theme
}

func (s *settings) SetTheme(theme fyne.Theme) {
	s.theme = theme
	s.apply()
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
	s.listenerMutex.Lock()
	s.changeListeners = append(s.changeListeners, listener)
	s.listenerMutex.Unlock()
}

func (s *settings) apply() {
	s.listenerMutex.RLock()
	for _, listener := range s.changeListeners {
		go func(listener chan fyne.Settings) {
			listener <- s
		}(listener)
	}
	s.listenerMutex.RUnlock()
}

func loadSettings() *settings {
	s := &settings{
		listenerMutex: &sync.RWMutex{},
	}

	return s
}
