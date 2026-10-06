//go:build wasm || test_web_driver

package common

const (
	ScrollAccelerateRate   = float64(10)
	ScrollAccelerateCutoff = float64(5)
	ScrollSpeed            = float32(0.2)
)
