// +build ci

package driver

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/test"

// NewApp creates a new headless app to test the examples code
func NewApp() app.App {
	return test.NewTestApp()
}
