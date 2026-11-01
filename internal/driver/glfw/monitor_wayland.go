//go:build !wasm && !test_web_driver && (linux || freebsd || openbsd || netbsd) && ((!x11 && !wayland) || wayland)

package glfw

import "github.com/go-gl/glfw/v3.4/glfw"

func (w *window) waylandMonitor() *glfw.Monitor {
	return w.viewport.GetWaylandWindowMonitor()
}
