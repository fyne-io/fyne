//go:build !ci && mobile && !android && !ios && !windows

package app

import (
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

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

func (a *fyneApp) ScheduleNotification(n *fyne.Notification, when time.Time) (*fyne.ScheduledNotification, error) {
	return a.scheduleViaScheduler(n, when)
}

func (a *fyneApp) CancelScheduledNotification(id string) error {
	return a.cancelViaScheduler(id)
}

func watchTheme(_ *settings) {
	// not implemented yet
}
