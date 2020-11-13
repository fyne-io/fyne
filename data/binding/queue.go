package binding

var itemQueue = make(chan itemData, 1024)

type itemData struct {
	fn   func(DataItem)
	item DataItem
	done chan interface{}
}

func queueItem(f func(DataItem), i DataItem) {
	itemQueue <- itemData{fn: f, item: i}
}

func init() {
	go processItems()
}

func processItems() {
	for {
		i := <-itemQueue
		if i.fn != nil {
			i.fn(i.item)
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
