//go:build !wasm && (linux || freebsd || openbsd || netbsd) && ((!x11 && !wayland) || wayland)

package glfw

/*
#cgo pkg-config: wayland-client
#include <stdlib.h>
#include <wayland-client.h>

typedef struct { int ready; struct wl_callback *cb; } frame_state;

static void frame_done(void *data, struct wl_callback *cb, uint32_t t) {
    (void)t;
    frame_state *s = (frame_state *)data;
    s->ready = 1;
    if (s->cb == cb) s->cb = NULL;
    wl_callback_destroy(cb);
}
static const struct wl_callback_listener frame_listener = { frame_done };

static frame_state *frame_state_new(void) {
    frame_state *s = calloc(1, sizeof(frame_state));
    s->ready = 1;
    return s;
}
static void frame_request(struct wl_surface *surface, frame_state *s) {
    s->ready = 0;
    if (s->cb) wl_callback_destroy(s->cb);
    s->cb = wl_surface_frame(surface);
    wl_callback_add_listener(s->cb, &frame_listener, s);
}
static int  frame_ready(frame_state *s) { return s->ready; }
static void frame_state_free(frame_state *s) {
    if (!s) return;
    if (s->cb) wl_callback_destroy(s->cb);
    free(s);
}
*/
import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2/internal/build"
)

type frameTracker struct {
	state  *C.frame_state
	window *window
}

func newPresentGate(w *window) presentGate {
	if !build.IsWayland {
		return noGate{}
	}
	return &frameTracker{state: C.frame_state_new(), window: w}
}

func (t *frameTracker) ready() bool { return t.state == nil || C.frame_ready(t.state) != 0 }

func (t *frameTracker) requestFrame() {
	if t.state == nil || t.window == nil || t.window.viewport == nil {
		return
	}
	C.frame_request((*C.struct_wl_surface)(unsafe.Pointer(t.window.viewport.GetWaylandWindow())), t.state)
}

func (t *frameTracker) markReady() {
	if t.state != nil {
		t.state.ready = 1
	}
}

func (t *frameTracker) free() {
	C.frame_state_free(t.state)
	t.state = nil
}
