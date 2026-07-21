//go:build (!linux || (x11 && !wayland) || (!x11 && wayland)) && !wasm && !test_web_driver

package glfw

// forcePlatform returns platformAuto when either x11 or wayland was specified
// or if we're not on a suitable OS.
func forcePlatform() string {
	return platformAuto
}
