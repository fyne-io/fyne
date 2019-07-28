// +build !ci

// +build !mobile

package app

import (
	"net/url"
	"os"
	"path/filepath"
)

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
