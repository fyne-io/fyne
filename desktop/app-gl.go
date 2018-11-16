// +build !ci,gl

package desktop

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/desktop/gl"

// NewApp returns a new app instance using the desktop (GL) driver.
func NewApp() fyne.App {
	log.Println("desktop.NewApp() is deprecated - please use app.New()")
	return fyne.NewAppWithDriver(gl.NewGLDriver())
}
