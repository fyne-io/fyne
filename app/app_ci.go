// +build ci

package app

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/test"
)

// New returns a new app instance using the test (headless) driver.
func New() fyne.App {
	return NewAppWithDriver(test.NewTestDriver())
}
