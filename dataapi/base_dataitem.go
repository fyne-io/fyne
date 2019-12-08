package dataapi

import (
	"sync"
)

// BaseDataItem provides a base implementation to inherit from, that provides listener functionality
type BaseDataItem struct {
	sync.RWMutex
	callbacks map[int]func(item DataItem)
	id        int
}

// NewBaseDataItem returns a new BaseDataItem
func NewBaseDataItem() *BaseDataItem {
	return &BaseDataItem{
		callbacks: make(map[int]func(DataItem)),
	}
}

// String returns an empty string, needed in order to implement stringer
func (b *BaseDataItem) String() string {
	return ""
}

// AddListener adds a new listener callback to this BaseDataItem
func (b *BaseDataItem) AddListener(f func(data DataItem)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

// DeleteListener removes the listener with the given ID
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
