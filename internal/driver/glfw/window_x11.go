//go:build !wayland && (linux || freebsd || openbsd || netbsd)
// +build !wayland
// +build linux freebsd openbsd netbsd

package glfw

import "strconv"

// GetWindowHandle returns the window handle. Only implemented for X11 currently.
func (w *window) GetWindowHandle() string {
	xid := uint(w.viewport.GetX11Window())
	return "x11:" + strconv.FormatUint(uint64(xid), 16)
}
