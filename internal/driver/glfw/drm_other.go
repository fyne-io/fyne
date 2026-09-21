//go:build !linux || wasm || test_web_driver

package glfw

func (*gLDriver) guardDRM() {}

func (*gLDriver) shouldSkipPoll() bool { return false }
