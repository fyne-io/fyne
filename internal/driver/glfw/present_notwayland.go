//go:build wasm || (!linux && !freebsd && !openbsd && !netbsd) || (x11 && !wayland)

package glfw

func newPresentGate(_ *window) presentGate { return noGate{} }
