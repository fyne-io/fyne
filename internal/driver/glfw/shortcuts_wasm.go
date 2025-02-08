//go:build wasm

package glfw

import (
	"strings"

	"syscall/js"
)

var isMacOS = strings.Contains(js.Global().Get("window").Get("navigator").Get("platform").String(), "Mac")

// Checks if running on Mac OSX
func isMacOSRuntime() bool {
	return isMacOS
}
