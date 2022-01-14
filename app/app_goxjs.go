//go:build !ci && (!android || !ios || !mobile) && (js || wasm || web)
// +build !ci
// +build !android !ios !mobile
// +build js wasm web

package app

import (
	"errors"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	// TODO #2736
	return errors.New("OpenURL is not supported yet with GopherJS backend.")
}

func (app *fyneApp) SendNotification(_ *fyne.Notification) {
	// TODO #2735
	fyne.LogError("Sending notification is not supported yet.", nil)
}

func rootConfigDir() string {
	return "/data/"
}

func defaultTheme() fyne.Theme {
	return theme.DarkTheme()
}
