package async

var mainGoroutineID uint64

func init() {
	mainGoroutineID = goroutineID()
}

func IsMainGoroutine() bool {
	return goroutineID() == mainGoroutineID
}
