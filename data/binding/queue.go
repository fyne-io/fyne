package binding

var itemQueue = make(chan itemData, 1024)

type itemData struct {
	fn   func()
	done chan interface{}
}

func queueItem(f func()) {
	itemQueue <- itemData{fn: f}
}

func init() {
	go processItems()
}

func processItems() {
	for {
		i := <-itemQueue
		if i.fn != nil {
			i.fn()
		}
		if i.done != nil {
			i.done <- struct{}{}
		}
	}
}

func waitForItems() {
	done := make(chan interface{})
	itemQueue <- itemData{done: done}
	<-done
}
