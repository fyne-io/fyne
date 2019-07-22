// +build !ci

// +build linux,android

package app

import (
	"net/url"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/android"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("/system/bin/am", "start", "-a", "android.intent.action.VIEW", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(android.NewAndroidDriver())
}
