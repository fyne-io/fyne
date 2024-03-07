//go:build wasm

package glfw

func goroutineID() uint64 {
	return mainGoroutineID
}
