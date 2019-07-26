// +build !ci,!android

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/glfw"
)

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(glfw.NewGLDriver())
}
