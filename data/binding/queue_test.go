package binding

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	in, out := makeInfiniteQueue()

	wg.Add(1)
	c := 0
	go func() {
		for range out {
			c++
		}
		wg.Done()
	}()

	for i := 0; i < 2048; i++ {
		in <- &itemData{}
	}
	close(in)

	wg.Wait()
	assert.Equal(t, 2048, c)
}
