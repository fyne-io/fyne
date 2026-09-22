package glfw

import "testing"

func TestNoGateAlwaysReady(t *testing.T) {
	var g presentGate = noGate{}
	if !g.ready() {
		t.Fatal("noGate.ready() = false, want true (off-Wayland must always present)")
	}
	g.requestFrame()
	g.markReady()
	g.free()
}

func TestDecideRepaint(t *testing.T) {
	cases := []struct {
		name                        string
		visible, ready, dirty       bool
		wantRepaint, wantDirtyCheck bool
	}{
		{"visible+ready+dirty", true, true, true, true, true},
		{"visible+ready+clean", true, true, false, false, true},
		{"not ready defers, keeps dirty", true, false, true, false, false},
		{"hidden defers, keeps dirty", false, true, true, false, false},
		{"hidden+not ready", false, false, true, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checked := false
			got := decideRepaint(c.visible, c.ready, func() bool {
				checked = true
				return c.dirty
			})
			if got != c.wantRepaint {
				t.Errorf("decideRepaint = %v, want %v", got, c.wantRepaint)
			}
			if checked != c.wantDirtyCheck {
				t.Errorf("dirty checked = %v, want %v (dirty must be preserved when not presentable)", checked, c.wantDirtyCheck)
			}
		})
	}
}

func TestNewPresentGateReadyByDefault(t *testing.T) {
	g := newPresentGate(nil)
	defer g.free()
	if !g.ready() {
		t.Fatal("newPresentGate().ready() = false, want true")
	}
}
