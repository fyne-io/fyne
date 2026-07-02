//go:build ((x11 && wayland) || (!x11 && !wayland)) && linux && !wasm && !test_web_driver

package glfw

/*
#cgo pkg-config: wayland-client
#include <stdlib.h>
#include <string.h>
#include <wayland-client.h>

static void fyne_registry_global(void *data, struct wl_registry *registry,
		uint32_t name, const char *interface, uint32_t version) {
	if (strcmp(interface, "zxdg_decoration_manager_v1") == 0) {
		*(int *)data = 1;
	}
}

static void fyne_registry_global_remove(void *data, struct wl_registry *registry,
		uint32_t name) {
}

static const struct wl_registry_listener fyne_registry_listener = {
	fyne_registry_global,
	fyne_registry_global_remove,
};

// fyne_wayland_has_ssd returns 1 if the running Wayland compositor advertises the
// xdg-decoration manager (server-side decorations are available), 0 if it does
// not (client-side decorations are forced), and -1 if no Wayland display could
// be reached.
static int fyne_wayland_has_ssd() {
	struct wl_display *display = wl_display_connect(NULL);
	if (display == NULL) {
		return -1;
	}

	int hasSSD = 0;
	struct wl_registry *registry = wl_display_get_registry(display);
	wl_registry_add_listener(registry, &fyne_registry_listener, &hasSSD);
	wl_display_roundtrip(display);

	wl_registry_destroy(registry);
	wl_display_disconnect(display);
	return hasSSD;
}
*/
import "C"

import "os"

// shouldForceX11 reports whether GLFW's auto-detected backend is built with X11
// and Wayland, and is running a Wayland session with Client Side Decorations (CSD).
// If so we should force X11 (XWayland) platform before glfw.Init() is called.
func shouldForceX11() bool {
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		return false // not a Wayland session, leave the choice to GLFW
	}
	if os.Getenv("DISPLAY") == "" {
		return false // no XWayland available to fall back to, stay on Wayland
	}

	// Zero means the detected compositor forces client-side decorations (libdecor).
	return C.fyne_wayland_has_ssd() == 0
}
