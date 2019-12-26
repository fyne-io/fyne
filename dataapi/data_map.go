package dataapi

import "sync"

// DataItemMap is a basic DataMap of strings
type DataItemMap struct {
	sync.RWMutex
	data      map[string]DataItem
	mData     sync.RWMutex
	callbacks map[int]func(item DataItem)
	id        int
}

// NewDataItemMap returns a new DataItemMap
func NewDataItemMap() *DataItemMap {
	return &DataItemMap{
		callbacks: make(map[int]func(DataItem)),
		data:      make(map[string]DataItem),
	}
}

// String returns a string representation of the whole map !
func (b *DataItemMap) String() string {
	value := ""
	i := 0
	for k, v := range b.data {
		if i > 0 {
			value = value + "\n"
		}
		value = value + k + ": " + v.String()
		i++
	}
	return value
}

// Get retuns the dataItem with the given key, and a flag denoting whether it was found
func (b *DataItemMap) Get(key string) (DataItem, bool) {
	b.mData.RLock()
	i, ok := b.data[key]
	b.mData.RUnlock()
	return i, ok
}

// UpdateItem changes an item in the map and updates listeners
func (b *DataItemMap) UpdateItem(key string, data DataItem) {
	b.Lock()
	b.data[key] = data
	b.Unlock()
	b.update()
}

// SetString sets a map value given a string
func (b *DataItemMap) SetString(key, value string) {
	b.Lock()
	b.data[key] = NewString(value)
	b.Unlock()
	b.update()
}

// DeleteItem removes an item from the map, and updates listeners
func (b *DataItemMap) DeleteItem(key string) {
	b.Lock()
	delete(b.data, key)
	b.Unlock()
	b.update()
}

// AddListener adds a new listener callback to this BaseDataItem
func (b *DataItemMap) AddListener(f func(data DataItem)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

// DeleteListener removes the listener with the given ID
func (b *DataItemMap) DeleteListener(i int) {
	b.Lock()
	defer b.Unlock()
	delete(b.callbacks, i)
}

func (b *DataItemMap) update() {
	b.RLock()
	defer b.RUnlock()
	for _, f := range b.callbacks {
		f(b)
	}
}
