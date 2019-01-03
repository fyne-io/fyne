// +build !ci,efl

package app

import (
	"fyne.io/fyne"
	"fyne.io/fyne/driver/efl"
)

// New returns a new app instance using the EFL driver.
func New() fyne.App {
	return NewAppWithDriver(efl.NewEFLDriver())
}
