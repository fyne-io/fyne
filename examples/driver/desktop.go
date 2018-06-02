// +build !ci

package driver

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/desktop"

// NewApp creates a new desktop app to run the examples
func NewApp() app.App {
	return desktop.NewApp()
}
