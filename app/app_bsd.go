// +build !ci

// +build openbsd freebsd netbsd

package app

import (
	"os/exec"
	"net/url"
)

func (app *fyneApp) OpenURL(url *url.URL) {
	exec.Command("xdg-open", url.String()).Run()
}
