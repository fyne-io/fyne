// +build !ci

package app

import "os/exec"

func (app *fyneApp) OpenURL(url string) {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
}
