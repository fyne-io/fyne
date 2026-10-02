//go:build !darwin || test_web_driver

package glfw

import "time"

func startMainLoop() func(time.Duration) {
	return nil
}

func stopMainLoop() {
}

func wakeMainLoop() {
}
