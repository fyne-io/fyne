package theme

import "fyne.io/fyne/v2"

// FromDesktopEnvironment returns a new WindowManagerTheme instance for the current desktop session.
// If the desktop manager is not supported or if it is not found, return the default theme
func FromDesktopEnvironment() fyne.Theme {
	return fromDesktopEnvironment()
}
