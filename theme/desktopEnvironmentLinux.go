//go:build linux
// +build linux

package theme

import (
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

// FromDesktopEnvironment returns a new WindowManagerTheme instance for the current desktop session.
func fromDesktopEnvironment() fyne.Theme {
	wm := os.Getenv("XDG_CURRENT_DESKTOP")
	if wm == "" {
		wm = os.Getenv("DESKTOP_SESSION")
	}
	wm = strings.ToLower(wm)

	switch wm {
	case "gnome", "xfce", "unity", "gnome-shell", "gnome-classic", "mate", "gnome-mate":
		return NewGnomeTheme(-1, GnomeFlagAutoReload)
	case "kde", "kde-plasma", "plasma":
		return NewKDETheme()

	}

	log.Println("Window manager not supported:", wm, "using default theme")
	return DefaultTheme()
}
