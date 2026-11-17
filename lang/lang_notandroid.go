//go:build !android

package lang

func initRuntime() {}

var runtimeReady = func() bool {
	return true
}
