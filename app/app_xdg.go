// +build !ci

// +build linux openbsd freebsd netbsd

package app

import (
	"net/url"
	"os"
	"os/exec"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := exec.Command("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
