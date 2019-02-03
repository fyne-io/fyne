// +build !ci

package app

import (
	"net/url"
	"os"
	"os/exec"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
