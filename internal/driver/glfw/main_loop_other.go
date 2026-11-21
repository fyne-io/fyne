//go:build !darwin || test_web_driver

package glfw

func startMainLoop() *mainLoopWake {
	return nil
}

func stopMainLoop() {
}

func wakeMainLoop() {
}
