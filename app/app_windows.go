// +build !ci

package app

import (
	"net/url"
	"os/exec"
)

func (app *fyneApp) OpenURL(url *url.URL) {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url.String()).Run()
}
