// +build !ci

// +build linux,android mobile

package app

import (
	"net/url"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/gomobile"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("/system/bin/am", "start", "-a", "android.intent.action.VIEW", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func rootConfigDir() string {
	homeDir := "/data" //, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(homeDir, "data")
	return filepath.Join(desktopConfig, "fyne")
}

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(gomobile.NewGoMobileDriver())
}
