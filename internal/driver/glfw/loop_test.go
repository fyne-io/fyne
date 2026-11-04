//go:build !no_glfw && !mobile

package glfw

import (
	"sync"
	"testing"
	"time"
)

// BenchmarkRunOnMain measures the cost of calling a function
// on the main thread.
func BenchmarkRunOnMain(b *testing.B) {
	f := func() {}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runOnMain(f)
	}
}

// BenchmarkRunOnDraw measures the cost of calling a function
// on the draw thread.
func BenchmarkRunOnDraw(b *testing.B) {
	f := func() {}
	w := createWindow("Test")
	w.create()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.RunWithContext(f)
	}
}

// Goroutines that queue work for the main thread while the app shuts down
// must neither panic nor wait forever.
func TestMainQueueAddDuringClose(t *testing.T) {
	q := newMainQueue()

	stop := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 10; i++ {
		workers.Add(1)
		go func(wait bool) {
			defer workers.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}

				data := funcData{f: func() {}}
				if wait {
					data.done = make(chan struct{}, 1)
				}
				if q.add(data) && wait {
					<-data.done
				}
			}
		}(i%2 == 0)
	}

	// stand in for the main loop: run queued functions until shutdown
	timeout := time.After(50 * time.Millisecond)
	for running := true; running; {
		select {
		case f := <-q.queue.Out():
			f.f()
			if f.done != nil {
				f.done <- struct{}{}
			}
		case <-timeout:
			running = false
		}
	}
	q.close()

	// the workers keep going after shutdown
	time.Sleep(100 * time.Millisecond)
	close(stop)

	finished := make(chan struct{})
	go func() {
		workers.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("goroutines queueing work are stuck after shutdown")
	}
}
