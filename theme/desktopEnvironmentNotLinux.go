//go:build !linux
// +build !linux

package theme

import (
	"fyne.io/fyne/v2"
)

// FromDesktopEnvironment returns the theme that is set in the desktop environment.
func fromDesktopEnvironment() fyne.Theme {
	return DefaultTheme()
}
