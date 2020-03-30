// +build !ci

// +build android

package app

import (
	"log"
	"net/url"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	return theme.LightTheme()
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("/system/bin/am", "start", "-a", "android.intent.action.VIEW", "--user", "0",
		"-d", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func rootConfigDir() string {
	filesDir := os.Getenv("FILESDIR")
	if filesDir == "" {
		log.Println("FILESDIR env was not set by android native code")
		return "/data/data" // probably won't work, but we can't make a better guess
	}

	return filesDir
}
