// +build ci

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/test"
)

// New returns a new app instance using the test (headless) driver.
func New() fyne.App {
	return NewAppWithDriver(test.NewDriver())
}
