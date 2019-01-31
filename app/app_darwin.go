// +build !ci

package app

import (
	"os/exec"
	"net/url"
)

func (app *fyneApp) OpenURL(url *url.URL) {
	exec.Command("open", url.String()).Run()
}
