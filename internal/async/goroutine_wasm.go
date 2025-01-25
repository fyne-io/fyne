//go:build wasm

package async

func goroutineID() uint64 {
	return mainGoroutineID.Load()
}
