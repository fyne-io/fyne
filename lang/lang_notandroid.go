//go:build !android

package lang

func initRuntime() {}

func platformLocaleReady() bool {
	return true
}
