//go:build !ci && wasm
// +build !ci,wasm

package app

import (
	"fmt"
	"net/url"
	"syscall/js"
)

func (app *fyneApp) OpenURL(url *url.URL) error {
	window := js.Global().Call("open", url.String(), "_blank", "")
	if window.Equal(js.Null()) {
		return fmt.Errorf("Unable to open a new window/tab for URL: %v.", url)
	}
	window.Call("focus")
	return nil
}
