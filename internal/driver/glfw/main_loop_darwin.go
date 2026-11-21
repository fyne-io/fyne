//go:build darwin && !test_web_driver

package glfw

import (
	"time"

	"github.com/go-gl/glfw/v3.4/glfw"
)

var nativeMainLoopWake = mainLoopWake{wait: func(timeout time.Duration) {
	glfw.WaitEventsTimeout(timeout.Seconds())
}}

func startMainLoop() *mainLoopWake {
	nativeMainLoopWake.start(glfw.PostEmptyEvent)
	return &nativeMainLoopWake
}

func stopMainLoop() {
	nativeMainLoopWake.stop()
}

func wakeMainLoop() {
	nativeMainLoopWake.notify()
}
