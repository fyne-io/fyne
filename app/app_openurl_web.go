//go:build !ci && !js && !wasm && test_web_driver
// +build !ci,!js,!wasm,test_web_driver

package app

import (
	"errors"
	"net/url"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("OpenURL is not supported with the test web driver.")
}
