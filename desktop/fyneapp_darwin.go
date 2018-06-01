// +build !ci

package fyneapp

import "os/exec"

import _ "github.com/fyne-io/fyne/desktop/efl" // import our Fyne driver for Mac OS X

func (app *fyneApp) OpenURL(url string) {
	exec.Command("open", url).Run()
}
