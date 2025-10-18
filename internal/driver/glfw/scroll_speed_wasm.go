//go:build wasm || test_web_driver

package glfw

const (
	scrollAccelerateRate   = float64(10)
	scrollAccelerateCutoff = float64(5)
	scrollSpeed            = float32(0.2)
)
