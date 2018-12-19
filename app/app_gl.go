// +build !ci,!efl

package app

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/driver/gl"
)

// New returns a new app instance using the OpenGL driver.
func New() fyne.App {
	return NewAppWithDriver(gl.NewGLDriver())
}
