package binding

import (
	"sync"
)

var (
	itemQueueIn      chan<- *itemData
	itemQueueOut     <-chan *itemData
	processOnce      sync.Once
	queueInitialized = make(chan int)
)

type itemData struct {
	fn   func()
	done chan interface{}
}

func queueItem(f func()) {
	if itemQueueIn == nil || itemQueueOut == nil {
		processOnce.Do(func() {
			itemQueueIn, itemQueueOut = makeInfiniteQueue()
			go processItems()
			close(queueInitialized)
		})
	}
	<-queueInitialized
	itemQueueIn <- &itemData{fn: f}
}

func makeInfiniteQueue() (chan<- *itemData, <-chan *itemData) {
	in := make(chan *itemData)
	out := make(chan *itemData)
	go func() {
		queued := make([]*itemData, 0, 1024)
		pending := func() chan *itemData {
			if len(queued) == 0 {
				return nil
			}
			return out
		}
		next := func() *itemData {
			if len(queued) == 0 {
				return nil
			}
			return queued[0]
		}
		for len(queued) > 0 || in != nil {
			select {
			case val, ok := <-in:
				if !ok {
					in = nil
				} else {
					queued = append(queued, val)
				}
			case pending() <- next():
				queued = queued[1:]
			}
		}
		close(out)
	}()
	return in, out
}

func processItems() {
	for {
		i := <-itemQueueOut
		if i.fn != nil {
			i.fn()
		}
		if i.done != nil {
			i.done <- struct{}{}
		}
	}
}

func waitForItems() {
	<-queueInitialized
	done := make(chan interface{})
	itemQueueIn <- &itemData{done: done}
	<-done
}
