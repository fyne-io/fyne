// +build !ci,!nacl

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/gl"
)

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(gl.NewGLDriver())
}
