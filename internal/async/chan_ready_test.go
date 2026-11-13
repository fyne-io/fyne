package async

import (
	"slices"
	"testing"
)

func TestUnboundedChanReadyNotificationOnClose(t *testing.T) {
	ch := &UnboundedChan[int]{
		in:    make(chan int, 1),
		out:   make(chan int, 1),
		close: make(chan struct{}),
		q:     []int{1, 2},
	}
	ch.in <- 3
	var notified []int
	ch.onReady = []func(){func() {
		select {
		case value := <-ch.Out():
			notified = append(notified, value)
		default:
			t.Fatal("close notification ran before output became available")
		}
	}}
	ch.closed()
	if !slices.Equal(notified, []int{1, 2, 3}) {
		t.Fatalf("close notifications lost or reordered values: %v", notified)
	}
	if _, open := <-ch.Out(); open {
		t.Fatal("output stayed open after draining")
	}
}
