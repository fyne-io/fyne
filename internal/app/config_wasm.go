//go:build !ci && (!android || !ios || !mobile) && (wasm || test_web_driver)

package app

func rootConfigDir() string {
	return "/data/"
}
