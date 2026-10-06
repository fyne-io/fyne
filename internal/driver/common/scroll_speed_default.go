//go:build !darwin && !wasm && !test_web_driver

package common

const (
	ScrollAccelerateRate   = float64(125)
	ScrollAccelerateCutoff = float64(10)
	ScrollSpeed            = float32(25)
)
