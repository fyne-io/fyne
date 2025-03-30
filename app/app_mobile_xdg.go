//go:build !ci && mobile && !android && !ios

package app

import (
	"errors"
	"net/url"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) OpenURL(_ *url.URL) error {
	return errors.New("mobile simulator does not support open URLs yet")
}

func (a *fyneApp) SendNotification(_ *fyne.Notification) {
	fyne.LogError("Notifications are not supported in the mobile simulator yet", nil)
}

func watchTheme(_ *settings) {
	// not implemented yet
}
