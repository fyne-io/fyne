package binding

import (
	"os"
	"runtime"
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/async"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestQueueLazyInit resets the current unbounded func queue, and tests
// if the queue is lazy initialized.
//
// Note that this test may fail, if any of other tests in this package
// calls t.Parallel().
func TestQueueLazyInit(t *testing.T) {
	if queue != nil { // Reset queues
		queue.Close()
		queue = nil
		once = sync.Once{}
	}

	initialGoRoutines := runtime.NumGoroutine()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		queueItem(func() { wg.Done() })
	}
	wg.Wait()

	n := runtime.NumGoroutine()
	if n > initialGoRoutines+2 {
		t.Fatalf("unexpected number of goroutines after initialization, probably leaking: got %v want %v", n, initialGoRoutines+2)
	}
}

func TestQueueItem(t *testing.T) {
	called := 0
	queueItem(func() { called++ })
	queueItem(func() { called++ })
	waitForItems()
	assert.Equal(t, 2, called)
}

func TestMakeInfiniteQueue(t *testing.T) {
	var wg sync.WaitGroup
	queue := async.NewUnboundedFuncChan()

	wg.Add(1)
	c := 0
	go func() {
		for range queue.Out() {
			c++
		}
		wg.Done()
	}()

	for i := 0; i < 2048; i++ {
		queue.In() <- func() {}
	}
	queue.Close()

	wg.Wait()
	assert.Equal(t, 2048, c)
}
