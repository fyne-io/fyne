//go:build ci || (mobile && !android && !ios) || (!linux && !darwin && !windows && !freebsd && !openbsd && !netbsd && !wasm && !test_web_driver)

package app

import (
	"errors"
	"net/url"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) OpenURL(_ *url.URL) error {
	return errors.New("Unable to open url for unknown operating system")
}

func (a *fyneApp) SendNotification(_ *fyne.Notification) {
	fyne.LogError("Refusing to show notification for unknown operating system", nil)
}

func watchTheme(_ *settings) {
	// no-op
}
