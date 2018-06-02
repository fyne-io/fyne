// +build !ci

package driver

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/desktop"

func NewApp() app.App {
	return desktop.NewApp()
}
