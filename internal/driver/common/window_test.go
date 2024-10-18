package common

import (
	"runtime"
	"testing"
)

func TestWindow(t *testing.T) {
	w := &Window{}
	w.InitEventQueue()
	w.DestroyEventQueue()
	runtime.Gosched()

	// checking if it will panic (it should not)
	w.QueueEvent(func() {})
}
