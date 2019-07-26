// +build !ci

// +build linux,!android openbsd freebsd netbsd

package app

import (
	"net/url"
	"os"
	"path/filepath"
)

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
