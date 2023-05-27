//go:build !(windows || linux || freebsd || netbsd || openbsd)
// +build !windows,!linux,!freebsd,!netbsd,!openbsd

package glfw

const (
	scrollAccelerateRate   = float64(5)
	scrollAccelerateCutoff = float64(5)
	scrollSpeed            = float32(10)
)
