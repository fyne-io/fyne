// +build !ci,!nacl,!android,!ios,!mobile

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/glfw"
)

// NewWithID returns a new app instance using the OpenGL driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return NewAppWithDriver(glfw.NewGLDriver(), id)
}
