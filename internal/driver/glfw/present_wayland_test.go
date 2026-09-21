//go:build !wasm && (linux || freebsd || openbsd || netbsd) && ((!x11 && !wayland) || wayland)

package glfw

import "testing"

func TestFrameTrackerAbandonMarksReady(t *testing.T) {
	g := newPresentGate(nil)
	defer g.free()
	ft, ok := g.(*frameTracker)
	if !ok {
		t.Skip("not a frameTracker (Wayland not active at runtime)")
	}
	if ft.state == nil {
		t.Fatal("frameTracker.state is nil")
	}
	ft.state.ready = 0
	g.resetAfterDisplayLoss()
	if !g.ready() {
		t.Fatal("abandon of a pending callback must mark ready")
	}
}
