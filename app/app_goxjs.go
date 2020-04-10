// +build !ci

// +build !android !ios !mobile
// +build js wasm web

package app

import (
	"net/url"
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("OpenURL is not supported yet with GopherJS backend.")
}

func (app *fyneApp) SendNotification(_ *fyne.Notification) {
	fyne.LogError("Sending notification is not supported yet.", nil)
}

func rootConfigDir() string {
	return "/data/"
}

func defaultTheme() fyne.Theme {
	return theme.DarkTheme()
}
