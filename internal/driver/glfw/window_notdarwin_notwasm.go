//go:build !darwin && !wasm && !test_web_driver

package glfw

import "github.com/go-gl/glfw/v3.4/glfw"

type monitor = glfw.Monitor

func (w *window) getSecondaryMonitor() *monitor {
	primary := glfw.GetPrimaryMonitor()
	for _, m := range glfw.GetMonitors() {
		if m.GetName() != primary.GetName() {
			return m
		}
	}

	return primary
}
