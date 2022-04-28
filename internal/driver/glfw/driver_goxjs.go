//go:build js
// +build js

package glfw

func goroutineID() int {
	return mainGoroutineID
}
