// +build ci

package driver

import "github.com/fyne-io/fyne/api/app"
import "github.com/fyne-io/fyne/test"

func NewApp() app.App {
	return test.NewTestApp()
}
