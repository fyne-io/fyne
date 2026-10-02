//go:build darwin && !test_web_driver

package glfw

import (
	"time"

	"github.com/go-gl/glfw/v3.4/glfw"
)

var nativeMainLoopWake mainLoopWake

func startMainLoop() func(time.Duration) {
	nativeMainLoopWake.start(glfw.PostEmptyEvent)
	return func(timeout time.Duration) {
		if timeout > 0 {
			glfw.WaitEventsTimeout(timeout.Seconds())
		}
	}
}

func stopMainLoop() {
	nativeMainLoopWake.stop()
}

func wakeMainLoop() {
	nativeMainLoopWake.notify()
}
