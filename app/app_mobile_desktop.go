//go:build !ci && mobile && !android && !ios

package app

import (
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", url.String())
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	} else {
		cmd := exec.Command("xdg-open", url.String())
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Start()
	}
}

func (a *fyneApp) SendNotification(_ *fyne.Notification) {
	fyne.LogError("Notifications are not supported in the mobile simulator yet", nil)
}

func watchTheme(_ *settings) {
	// not implemented yet
}
