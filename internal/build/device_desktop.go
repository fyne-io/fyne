//go:build !mobile && !js && !wasm
// +build !mobile,!js,!wasm

package build

// IsMobile returns true if running on a mobile device.
func IsMobile() bool {
	return false
}
