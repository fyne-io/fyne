package binding

import "sync"

// Binding is the base interface of the Data Binding API.
type Binding interface {
	addListener(func())
}

// ItemBinding implements a data binding for a single item.
type ItemBinding struct {
	Binding
	sync.RWMutex
	listeners []func()
}

func (b *ItemBinding) addListener(listener func()) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *ItemBinding) notify() {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l()
	}
}

// ListBinding implements a data binding for a list of bindings.
type ListBinding struct {
	ItemBinding
	values []Binding
}

// Length returns the length of the bound list.
func (b *ListBinding) Length() int {
	return len(b.values)
}

// Get returns the binding at the given index.
func (b *ListBinding) Get(index int) Binding {
	return b.values[index]
}

// Append adds the given binding(s) to the list.
func (b *ListBinding) Append(data ...Binding) {
	b.values = append(b.values, data...)
	b.notify()
}

// Set puts the given binding into the list at the given index.
func (b *ListBinding) Set(index int, data Binding) {
	old := b.values[index]
	if old == data {
		return
	}
	b.values[index] = data
	b.notify()
}

// AddListener adds the given listener to the binding.
func (b *ListBinding) AddListener(listener func()) {
	b.addListener(listener)
}

// MapBinding implements a data binding for a map string to binding.
type MapBinding struct {
	ItemBinding
	values map[string]Binding
}

// Length returns the length of the bound map.
func (b *MapBinding) Length() int {
	return len(b.values)
}

// Get returns the binding for the given key.
func (b *MapBinding) Get(key string) (Binding, bool) {
	v, ok := b.values[key]
	return v, ok
}

// Set puts the given binding into the map at the given key.
func (b *MapBinding) Set(key string, data Binding) {
	old, ok := b.values[key]
	if ok && old == data {
		return
	}
	b.values[key] = data
	b.notify()
}

// AddListener adds the given listener to the binding.
func (b *MapBinding) AddListener(listener func()) {
	b.addListener(listener)
}
