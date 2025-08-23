//go:build android || ios || mobile

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func (s *settings) loadSystemTheme() fyne.Theme {
	return theme.DefaultTheme()
}

func (s *settings) watchSettings() {
	// no-op on mobile
}

func (s *settings) stopWatching() {
	// no-op on mobile
}
