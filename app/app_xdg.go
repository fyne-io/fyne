// +build !ci

// +build linux openbsd freebsd netbsd
// +build !android,!mobile

package app

import (
	"net/url"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	return theme.DarkTheme()
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(homeDir, ".config")
	return filepath.Join(desktopConfig, "fyne")
}
