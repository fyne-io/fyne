// +build !ci,gl

package desktop

import "log"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/driver/gl"

// NewApp returns a new app instance using the desktop (GL) driver.
func NewApp() fyne.App {
	log.Println("desktop.NewApp() is deprecated - please use app.New()")
	return app.NewAppWithDriver(gl.NewGLDriver())
}
