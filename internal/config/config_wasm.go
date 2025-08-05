//go:build !ci && (!android || !ios || !mobile) && (wasm || test_web_driver)

package config

func rootConfigDir() string {
	return "/data/"
}
