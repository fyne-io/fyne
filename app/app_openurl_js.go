//go:build !ci && js && !wasm
// +build !ci,js,!wasm

package app

import (
	"fmt"
	"net/url"

	"honnef.co/go/js/dom"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	window := dom.GetWindow().Open(url.String(), "_blank", "")
	if window == nil {
		return fmt.Errorf("Unable to open a new window/tab for URL: %v.", url)
	}
	window.Focus()
	return nil
}
