//go:build wayland && (linux || freebsd || openbsd || netbsd)

package glfw

// GetWindowHandle returns the window handle. Only implemented for X11 currently.
func (w *window) GetWindowHandle() string {
	return "" // TODO: Find a way to get the Wayland handle for xdg_foreign protocol. Return "wayland:{id}".
}
