//go:build darwin

package glfw

const (
	// MacOS applies its own scroll accelerate curve, so set
	// scrollAccelerateRate to 1 for no acceleration effect
	scrollAccelerateRate   = float64(1)
	scrollAccelerateCutoff = float64(5)
	scrollSpeed            = float32(10)
)
