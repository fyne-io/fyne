//go:build !wayland || x11

package build

// IsWayland is true when compiling for the wayland windowing system, or auto-detecting and Wayland is loaded.
var IsWayland = false
