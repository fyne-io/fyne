package binding

import "sync"

type Binding interface {
}

type itemBinding struct {
	Binding
	sync.RWMutex
	listeners []func(interface{})
}

func (b *itemBinding) addListener(listener func(interface{})) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *itemBinding) notify(value interface{}) {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l(value)
	}
}

type sliceBinding struct {
	Binding
	sync.RWMutex
	listeners []func(int, Binding)
}

func (b *sliceBinding) addListener(listener func(int, Binding)) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *sliceBinding) notify(index int, value Binding) {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l(index, value)
	}
}

type mapBinding struct {
	Binding
	sync.RWMutex
	listeners []func(string, Binding)
}

func (b *mapBinding) addListener(listener func(string, Binding)) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *mapBinding) notify(key string, value Binding) {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l(key, value)
	}
}
