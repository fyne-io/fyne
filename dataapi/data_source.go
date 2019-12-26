package dataapi

import "sync"

// BaseDataSource is a null implementation of a data source
type BaseDataSource struct {
	sync.RWMutex
	data      []DataItem
	mData     sync.RWMutex
	callbacks map[int]func(item DataItem)
	id        int
}

// NewNullDataSource returns a new BaseDataSource
func NewBaseDataSource() *BaseDataSource {
	return &BaseDataSource{
		callbacks: make(map[int]func(DataItem)),
	}
}

// Count returns the size of the dataSource
func (b *BaseDataSource) Count() int {
	b.mData.RLock()
	defer b.mData.RUnlock()
	return len(b.data)
}

// String returns a string representation of the whole data source
func (b *BaseDataSource) String() string {
	//return ""
	value := ""
	i := 0
	for _, v := range b.data {
		if i > 0 {
			value = value + "\n"
		}
		value = value + v.String()
		i++
	}
	return value
}

// Get retuns the dataItem with the given index, and a flag denoting whether it was found
func (b *BaseDataSource) Get(idx int) (DataItem, bool) {
	b.mData.RLock()
	defer b.mData.RUnlock()
	if idx < 0 || idx >= len(b.data) {
		return nil, false
	}
	return b.data[idx], true
}

// AppendItem adds an item to the data and invokes the listeners
func (b *BaseDataSource) AppendItem(data DataItem) {
	b.mData.Lock()
	b.data = append(b.data, data)
	b.update()
	b.mData.Unlock()
}

// SetItem sets a given item, and invokes the listeners n
func (b *BaseDataSource) SetItem(idx int, data DataItem) {
	b.mData.Lock()
	if idx < 0 || idx >= len(b.data) {
		b.mData.Unlock()
		return
	}
	b.data[idx] = data
	b.mData.Unlock()
}

// DeleteItem removes an item from the map, and invokes the listeners
func (b *BaseDataSource) DeleteItem(idx int) {
	b.mData.Lock()
	if idx < 0 || idx >= len(b.data) {
		b.mData.Unlock()
		return
	}
	b.data = append(b.data[:idx], b.data[idx+1:]...)
	b.update()
	b.mData.Unlock()
}

// AddListener adds a new listener callback to this BaseDataItem
func (b *BaseDataSource) AddListener(f func(data DataItem)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

// DeleteListener removes the listener with the given ID
func (b *BaseDataSource) DeleteListener(i int) {
	b.Lock()
	defer b.Unlock()
	delete(b.callbacks, i)
}

func (b *BaseDataSource) update() {
	b.RLock()
	for _, f := range b.callbacks {
		f(b)
	}
	b.RUnlock()
}
