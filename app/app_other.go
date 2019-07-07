// +build ci !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package app

import (
	"errors"
	"net/url"
)

func rootConfigDir() string {
	return "/tmp/"
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("Unable to open url for unknown operating system")
}
