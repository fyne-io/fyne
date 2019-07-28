// +build !ci

// +build android ios mobile

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/gomobile"
)

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(gomobile.NewGoMobileDriver())
}
