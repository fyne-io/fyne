//go:build js && wasm
// +build js,wasm

package glfw

import (
	"strings"

	"syscall/js"
)

// Checks if the browser is runnning on Mac OSX
func isMacOSBrowser() bool {
	if strings.Contains(strings.ToLower(js.Global().Get("window").Get("navigator").Get("platform").String()), "mac") {
		return true
	} else {
		return false
	}
}
