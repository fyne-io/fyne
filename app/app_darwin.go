// +build !ci

// +build !mobile,!ios

package app

import (
	"net/url"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	// TODO read the macOS setting in Mojave onwards
	return theme.DarkTheme()
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "fyne")
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
