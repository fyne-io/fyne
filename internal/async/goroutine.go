package async

// mainGoroutineID stores the main goroutine ID.
// This ID must be initialized in main.init because
// a main goroutine may not equal to 1 due to the
// influence of a garbage collector.
var mainGoroutineID uint64

func init() {
	mainGoroutineID = goroutineID()
}

func IsMainGoroutine() bool {
	return goroutineID() == mainGoroutineID
}
