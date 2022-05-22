//go:build js
// +build js

package glfw

func goroutineID() uint64 {
	return mainGoroutineID
}
