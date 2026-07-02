//go:build (!linux || (x11 && !wayland) || (!x11 && wayland)) && !wasm && !test_web_driver

package glfw

// shouldForceX11 returns false when either x11 or wayland was specified
// or if we're not on a suitable OS.
func shouldForceX11() bool {
	return false
}
