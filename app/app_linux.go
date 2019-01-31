// +build !ci

package app

import (
	"net/url"
	"os/exec"
)

func (app *fyneApp) OpenURL(url *url.URL) {
	exec.Command("xdg-open", url.String()).Run()
}
