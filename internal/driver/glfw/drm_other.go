//go:build !linux || wasm || test_web_driver

package glfw

func (*gLDriver) hideIfDRMOutputLost() {}

func (*gLDriver) shouldSkipPoll() bool { return false }
