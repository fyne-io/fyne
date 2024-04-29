//go:build ci || software

package app

import (
	"image"

	"fyne.io/fyne/v2"
	software2 "fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/test"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(test.NewDriverWithPainter(software.NewPainter()), id)
}

// NewWithSoftwareDriver returns a new app instance using the Software (custom) driver.
// The ID string should be globally unique to this app.
func NewWithSoftwareDriver(id string, painter func(image.Image), events chan any) fyne.App {
	return newAppWithDriver(software2.NewDriver(painter, events), id)
}
