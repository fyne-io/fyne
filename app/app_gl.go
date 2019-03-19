// +build !ci

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/gl"
)

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(gl.NewGLDriver())
}
