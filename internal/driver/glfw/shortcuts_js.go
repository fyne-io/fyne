//go:build js && !wasm
// +build js,!wasm

package glfw

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

// Checks if running on Mac OSX
func isMacOSRuntime() bool {
	return strings.Contains(strings.ToLower(js.Global.Get("window").Get("navigator").Get("platform").String()), "mac")
}
