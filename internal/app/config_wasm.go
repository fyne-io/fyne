//go:build !ci && !android && !ios && !mobile && (wasm || test_web_driver) && !noos && !tinygo

package app

func RootConfigDir() string {
	return "/data/"
}
