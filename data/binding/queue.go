package binding

import (
	"runtime"

	"fyne.io/fyne/v2"
)

var mainGoroutineID uint64

func init() {
	runtime.LockOSThread()
	mainGoroutineID = goroutineID()
}

func goroutineID() (id uint64) {
	var buf [30]byte
	runtime.Stack(buf[:], false)
	for i := 10; buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}

	return id
}

func queueItem(f func()) {
	if goroutineID() == mainGoroutineID {
		f()
		return
	}

	fyne.CurrentApp().Driver().CallFromGoroutine(f)
}
