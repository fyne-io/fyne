package binding

import "sync"

type Binding interface {
	addListener(func())
}

type itemBinding struct {
	Binding
	sync.RWMutex
	listeners []func()
}

func (b *itemBinding) addListener(listener func()) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *itemBinding) notify() {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l()
	}
}

type SliceBinding struct {
	itemBinding
	values []Binding
}

func (b *SliceBinding) Length() int {
	return len(b.values)
}

func (b *SliceBinding) Get(index int) Binding {
	return b.values[index]
}

func (b *SliceBinding) Append(data ...Binding) {
	b.values = append(b.values, data...)
	b.notify()
}

func (b *SliceBinding) Set(index int, data Binding) {
	old := b.values[index]
	if old == data {
		return
	}
	b.values[index] = data
	b.notify()
}

func (b *SliceBinding) AddListener(listener func()) {
	b.addListener(listener)
}

type MapBinding struct {
	itemBinding
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
