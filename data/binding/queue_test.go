package binding

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/async"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testQueueLazyInit()
	os.Exit(m.Run())
}

func TestQueueItem(t *testing.T) {
	called := 0
	queueItem(func() {
		called++
	})
	queueItem(func() {
		called++
	})

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
	close(queue.In())

	wg.Wait()
	assert.Equal(t, 2048, c)
}

func testQueueLazyInit() {
	initialGoRoutines := runtime.NumGoroutine()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go queueItem(func() { wg.Done() })
	}
	wg.Wait()
	if runtime.NumGoroutine() != initialGoRoutines+2 {
		fmt.Println("--- FAIL: testQueueLazyInit")
		os.Exit(1)
	}
}
