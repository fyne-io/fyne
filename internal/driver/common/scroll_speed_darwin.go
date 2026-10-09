//go:build darwin

package common

const (
	// MacOS applies its own scroll accelerate curve, so set
	// ScrollAccelerateRate to 1 for no acceleration effect
	ScrollAccelerateRate   = float64(1)
	ScrollAccelerateCutoff = float64(5)
	ScrollSpeed            = float32(10)
)
