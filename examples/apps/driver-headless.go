// +build ci

package apps

import "github.com/fyne-io/fyne/api"
import "github.com/fyne-io/fyne/test"

// NewApp creates a new headless app to test the examples code
func NewApp() fyne.App {
	return test.NewTestApp()
}
