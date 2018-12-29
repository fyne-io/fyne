// +build !ci,efl

package app

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/driver/efl"
)

// New returns a new app instance using the EFL driver.
func New() fyne.App {
	return NewAppWithDriver(efl.NewEFLDriver())
}
