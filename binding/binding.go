package binding

import "sync"

type UnitBinding interface {
	addListener(func(interface{}))
	Set(interface{})
}

type BaseBinding struct {
	sync.RWMutex
	listeners []func(interface{})
}

func (b *BaseBinding) addListener(listener func(interface{})) {
	b.Lock()
	defer b.Unlock()
	b.listeners = append(b.listeners, listener)
}

func (b *BaseBinding) notify(value interface{}) {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		go l(value)
	}
}
