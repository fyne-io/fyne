//go:build !wasm

package glfw

import "runtime"

// Checks if running on Mac OSX
func isMacOSRuntime() bool {
	return runtime.GOOS == "darwin"
}
