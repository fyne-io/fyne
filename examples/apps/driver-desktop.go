// +build !ci

package apps

import "github.com/fyne-io/fyne/api"
import "github.com/fyne-io/fyne/desktop"

// NewApp creates a new desktop app to run the examples
func NewApp() fyne.App {
	return desktop.NewApp()
}
