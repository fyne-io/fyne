// +build !ci,!gl

// Package desktop provides a full implementation of the Fyne APIs for desktop
// applications. This supports Windows, Mac OS X and Linux using the EFL (Evas)
// render pipeline.
package desktop

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/desktop/efl"

// NewApp returns a new app instance using the desktop (EFL) driver.
func NewApp() fyne.App {
	return fyne.NewAppWithDriver(efl.NewEFLDriver())
}
