// +build !ci

package app

import "os/exec"

func (app *fyneApp) OpenURL(url string) {
	exec.Command("open", url).Run()
}
