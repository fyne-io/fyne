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
