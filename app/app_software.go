// +build ci nacl

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return NewAppWithDriver(test.NewDriverWithPainter(software.NewPainter()), id)
}
