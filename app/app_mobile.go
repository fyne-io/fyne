// +build !ci

// +build android ios mobile

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/gomobile"
)

var systemTheme fyne.ThemeVariant

// NewWithID returns a new app instance using the appropriate runtime driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	d := gomobile.NewGoMobileDriver()
	a := newAppWithDriver(d, id)
	d.(gomobile.ConfiguredDriver).SetOnConfigurationChanged(func(c *gomobile.Configuration) {
		systemTheme = c.SystemTheme

		a.Settings().(*settings).setupTheme()
	})
	return a
}
