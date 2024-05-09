//go:build ci || software

package app

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/software"
	painter "fyne.io/fyne/v2/internal/painter/software"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(software.NewDriverWithPainter(painter.NewPainter()), id)
}

// NewWithSoftwareDriver returns a new app instance using the Software (custom) driver.
// The ID string should be globally unique to this app.
func NewWithSoftwareDriver(id string, painter func(image.Image, []image.Rectangle), events chan any) fyne.App {
	return newAppWithDriver(software.NewDriver(painter, events), id)
}
