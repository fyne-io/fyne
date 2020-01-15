package dataapi

import (
	"sort"
	"strings"
	"sync"
)

// SliceDataSource is an implementation of a data source based on a simple []string
type SliceDataSource struct {
	sync.RWMutex
	Data      []DataItem
	mData     sync.RWMutex
	callbacks map[int]func(item DataSource)
	id        int
}

// NewSliceDataSource returns a new SliceDataSource
func NewSliceDataSource(data []DataItem) *SliceDataSource {
	return &SliceDataSource{
		Data: data,
		callbacks: make(map[int]func(DataSource)),
	}
}

// Count returns the size of the dataSource
func (b *SliceDataSource) Count() int {
	b.mData.RLock()
	defer b.mData.RUnlock()
	return len(b.Data)
}

// String returns a string representation of the whole data source
func (b *SliceDataSource) String() string {
	//return ""
	value := ""
	i := 0
	for _, v := range b.Data {
		if i > 0 {
			value = value + "\n"
		}
		value = value + v.String()
		i++
	}
	return value
}

// Get retuns the dataItem with the given index, and a flag denoting whether it was found
func (b *SliceDataSource) Get(idx int) (DataItem, bool) {
	b.mData.RLock()
	defer b.mData.RUnlock()
	if idx < 0 || idx >= len(b.Data) {
		return nil, false
	}
	return b.Data[idx], true
}

// GetStringSlice returns a copy of the underlying values
func (b *SliceDataSource) GetStringSlice() []string {
	b.mData.RLock()
	data := make([]string, 0, len(b.Data))
	for _, v := range b.Data {
		data = append(data, v.String())
	}
	b.mData.RUnlock()
	return data
}

// Append adds an item to the data and invokes the listeners
func (b *SliceDataSource) Append(data DataItem) *SliceDataSource {
	b.mData.Lock()
	b.Data = append(b.Data, data)
	b.mData.Unlock()
	b.update()
	return b
}

// SetItem sets a given item, and invokes the listeners n
func (b *SliceDataSource) SetItem(idx int, data DataItem) *SliceDataSource {
	b.mData.Lock()
	if idx < 0 || idx >= len(b.Data) {
		b.mData.Unlock()
		return b
	}
	b.Data[idx] = data
	b.mData.Unlock()
	return b
}

// SetFromStringSlice sets the underlying slice, and invokes the listeners
func (b *SliceDataSource) SetFromStringSlice(data []string) *SliceDataSource {
	b.mData.Lock()
	b.Data = make([]DataItem, 0, len(data))
	for _, v := range data {
		b.Data = append(b.Data, NewString(v))
	}
	b.mData.Unlock()
	b.update()
	return b
}

// SetFromMapKeys sets the underlying slice from maps keys, and invokes the listeners
func (b *SliceDataSource) SetFromMapKeys(m map[string]interface{}) *SliceDataSource {
	b.mData.Lock()
	b.Data = make([]DataItem, 0, len(m))
	for k := range m {
		b.Data = append(b.Data, NewString(k))
	}
	sort.Slice(b.Data, func(i, j int) bool {
		return strings.Compare(b.Data[i].String(), b.Data[j].String()) > 0
	})
	b.mData.Unlock()
	b.update()
	return b
}

// DeleteItem removes an item from the map, and invokes the listeners
func (b *SliceDataSource) DeleteItem(idx int) {
	b.mData.Lock()
	if idx < 0 || idx >= len(b.Data) {
		b.mData.Unlock()
		return
	}
	b.Data = append(b.Data[:idx], b.Data[idx+1:]...)
	b.mData.Unlock()
	b.update()
}

// AddListener adds a new listener callback to this BaseDataItem
func (b *SliceDataSource) AddListener(f func(data DataSource)) int {
	b.Lock()
	defer b.Unlock()
	b.id++
	b.callbacks[b.id] = f
	return b.id
}

// DeleteListener removes the listener with the given ID
func (b *SliceDataSource) DeleteListener(i int) {
	b.Lock()
	defer b.Unlock()
	delete(b.callbacks, i)
}

func (b *SliceDataSource) update() {
	b.RLock()
	for _, f := range b.callbacks {
		f(b)
	}
	b.RUnlock()
}
