//go:build !ci && !software && !wasm && test_web_driver

package app

import (
	"errors"
	"net/url"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("OpenURL is not supported with the test web driver.")
}
