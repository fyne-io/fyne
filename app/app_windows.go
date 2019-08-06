// +build !ci

package app

import (
	"filepath"
	"net/url"
	"os"
)

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming")
	return filepath.Join(desktopConfig, "fyne")
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
