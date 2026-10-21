//go:build linux && !wasm && !test_web_driver

package glfw

import "testing"

func TestDRMLostOutput(t *testing.T) {
	if drmLostOutput(nil, map[string]drmConn{"card0-eDP-1": {Status: "connected", Enabled: "enabled"}}) {
		t.Fatal("empty previous snapshot is startup, not a loss")
	}
	prev := map[string]drmConn{
		"card0-eDP-1": {Status: "connected", Enabled: "enabled"},
		"card0-DP-1":  {Status: "connected", Enabled: "enabled"},
	}
	curStatus := map[string]drmConn{
		"card0-eDP-1": {Status: "connected", Enabled: "enabled"},
		"card0-DP-1":  {Status: "disconnected", Enabled: "disabled"},
	}
	if !drmLostOutput(prev, curStatus) {
		t.Fatal("DP status connected→disconnected must count as a loss")
	}
	curEnabled := map[string]drmConn{
		"card0-eDP-1": {Status: "connected", Enabled: "enabled"},
		"card0-DP-1":  {Status: "connected", Enabled: "disabled"},
	}
	if !drmLostOutput(prev, curEnabled) {
		t.Fatal("DP enabled→disabled (KVM) must count as a loss")
	}
	if drmLostOutput(prev, prev) {
		t.Fatal("unchanged outputs are not a loss")
	}
}

func TestShowOnRestoreSkipsHiddenWindow(t *testing.T) {
	shown := &window{visible: true}
	hidden := &window{}
	closing := &window{visible: true, closing: true}

	shown.noteVisibilityForRestore()
	hidden.noteVisibilityForRestore()
	closing.noteVisibilityForRestore()

	if !shown.showOnRestore {
		t.Fatal("a visible window must be mapped again")
	}
	if hidden.showOnRestore {
		t.Fatal("a hidden window must stay unmapped")
	}
	if closing.showOnRestore {
		t.Fatal("a closing window must not be mapped again")
	}
}

func TestDRMStillMissing(t *testing.T) {
	lost := map[string]drmConn{
		"card1-HDMI-A-1": {Status: "connected", Enabled: "enabled"},
		"card1-DP-1":     {Status: "connected", Enabled: "enabled"},
	}
	away := map[string]drmConn{
		"card1-HDMI-A-1": {Status: "connected", Enabled: "disabled"},
		"card1-DP-1":     {Status: "connected", Enabled: "enabled"},
	}
	if !drmStillMissing(lost, away) {
		t.Fatal("enabled→disabled must keep the display down")
	}
	back := map[string]drmConn{
		"card1-HDMI-A-1": {Status: "connected", Enabled: "enabled"},
		"card1-DP-1":     {Status: "connected", Enabled: "enabled"},
	}
	if drmStillMissing(lost, back) {
		t.Fatal("restored connector should allow the window back")
	}
}
