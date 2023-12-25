//go:build !js && !wasm
// +build !js,!wasm

package glfw

import "runtime"

// Checks if running on Mac OSX
func isMacOSRuntime() bool {
	return runtime.GOOS == "darwin"
}
