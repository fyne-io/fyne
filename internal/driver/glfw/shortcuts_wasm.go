//go:build wasm

package glfw

import (
	"strings"

	"syscall/js"
)

var isMacOS = strings.Contains(strings.ToLower(js.Global().Get("window").Get("navigator").Get("platform").String()), "mac")

// Checks if running on Mac OSX
func isMacOSRuntime() bool {
	return isMacOS
}
