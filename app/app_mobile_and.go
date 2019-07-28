// +build !ci

// +build android mobile

package app

import (
	"net/url"
	"os"
	"path/filepath"
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
