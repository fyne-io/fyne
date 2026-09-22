//go:build !linux || wasm || test_web_driver

package glfw

func (*gLDriver) skipPollForDRM() bool { return false }
