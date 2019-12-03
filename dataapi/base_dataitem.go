package dataapi

import (
	"sync"
)

// BaseDataItem
type BaseDataItem struct {
	sync.RWMutex
	callbacks map[int]func(item DataItem)
	id        int
}

func NewBaseDataItem() *BaseDataItem {
	return &BaseDataItem{
		callbacks: make(map[int]func(DataItem)),
	}
}

func (b *BaseDataItem) String() string {
	return ""
}

func (b *BaseDataItem) AddListener(f func(data DataItem)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

func (b *BaseDataItem) DeleteListener(i int) {
	b.Lock()
	defer b.Unlock()
	delete(b.callbacks, i)
}

func (b *BaseDataItem) update() {
	b.RLock()
	defer b.RUnlock()
	for _, f := range b.callbacks {
		f(b)
	}
}
