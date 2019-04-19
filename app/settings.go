package app

import (
	"sync"

	"fyne.io/fyne"
)

type settings struct {
	themeLock sync.RWMutex
	theme     fyne.Theme

	listenerLock    sync.Mutex
	changeListeners []chan fyne.Settings
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

func loadSettings() *settings {
	s := &settings{}

	return s
}
