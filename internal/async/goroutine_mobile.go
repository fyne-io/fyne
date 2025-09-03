//go:build mobile

package async

// IsMainGoroutine returns true if it is called from the main goroutine, false otherwise.
func IsMainGoroutine() bool {
	routineID := mainGoroutineID.Load()
	return routineID == 0 || goroutineID() == routineID
}
