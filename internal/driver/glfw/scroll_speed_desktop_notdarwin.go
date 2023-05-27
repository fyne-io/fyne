//go:build windows || linux || freebsd || netbsd || openbsd
// +build windows linux freebsd netbsd openbsd

package glfw

const (
	scrollAccelerateRate = float64(125)
	scrollAccelerateCutoff = float64(10)
	scrollSpeed            = float32(25)
)