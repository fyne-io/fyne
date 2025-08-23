//go:build wasm || test_web_driver

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func (s *settings) load() {
	s.setupTheme()
	s.schema.Scale = 1
}

func (s *settings) loadFromFile(path string) error {
	return nil
}

func (s *settings) loadSystemTheme() fyne.Theme {
	return theme.DefaultTheme()
}

func watchFile(path string, callback func()) {
}

func (s *settings) watchSettings() {
	watchTheme(s)
}

func (s *settings) stopWatching() {
	stopWatchingTheme()
}
