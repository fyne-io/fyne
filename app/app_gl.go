//go:build !ci && !android && !ios && !mobile
// +build !ci,!android,!ios,!mobile

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/glfw"
)

// NewWithID returns a new app instance using the appropriate runtime driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(glfw.NewGLDriver(), id)
}
