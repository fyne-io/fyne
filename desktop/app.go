// +build !ci,!gl

// Package desktop is a legacy package which provided a an EFL based driver.
// This supports Windows, Mac OS X and Linux using the EFL (Evas) render pipeline.
// The replacement app package should be used by calling app.New() to create a new desktop application.
package desktop

import "log"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/driver/efl"

// NewApp returns a new app instance using the desktop (EFL) driver.
func NewApp() fyne.App {
	log.Println("desktop.NewApp() is deprecated - please use app.New()")
	return app.NewAppWithDriver(efl.NewEFLDriver())
}
