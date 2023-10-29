//go:build !js && !wasm
// +build !js,!wasm

package glfw

// Checks if browser is running on Mac OSX. Returns false because non-browser builds are a no op.
func isMacOSBrowser() bool {
	return false
}
