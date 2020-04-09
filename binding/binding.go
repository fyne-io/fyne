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

func (b *ListBinding) Length() int {
	return len(b.values)
}

func (b *ListBinding) Get(index int) Binding {
	return b.values[index]
}

func (b *ListBinding) Append(data ...Binding) {
	b.values = append(b.values, data...)
	b.notify()
}

func (b *ListBinding) Set(index int, data Binding) {
	old := b.values[index]
	if old == data {
		return
	}
	b.values[index] = data
	b.notify()
}

func (b *ListBinding) AddListener(listener func()) {
	b.addListener(listener)
}

type MapBinding struct {
	ItemBinding
	values map[string]Binding
}

func (b *MapBinding) Length() int {
	return len(b.values)
}

func (b *MapBinding) Get(key string) (Binding, bool) {
	v, ok := b.values[key]
	return v, ok
}

func (b *MapBinding) Set(key string, data Binding) {
	old, ok := b.values[key]
	if ok && old == data {
		return
	}
	b.values[key] = data
	b.notify()
}

func (b *MapBinding) AddListener(listener func()) {
	b.addListener(listener)
}
