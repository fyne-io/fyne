package binding

import (
	"sync/atomic"
)

func newBaseItem[T bool | float64 | int | rune | string]() *baseItem[T] {
	return &baseItem[T]{val: new(T)}
}

type baseItem[T bool | float64 | int | rune | string] struct {
	base

	val *T
}

func (b *baseItem[T]) Get() (T, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return *new(T), nil
	}
	return *b.val, nil
}

func (b *baseItem[T]) Set(val T) error {
	b.lock.Lock()
	oldVal := *b.val
	*b.val = val
	b.lock.Unlock()

	if oldVal != val {
		b.trigger()
	}

	return nil
}

func baseBindExternal[T bool | float64 | int | rune | string](val *T) *baseExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &baseExternalItem[T]{}
	b.val = val
	b.old = *val
	return b
}

type baseExternalItem[T bool | float64 | int | rune | string] struct {
	baseItem[T]

	old T
}

func (b *baseExternalItem[T]) Set(val T) error {
	b.lock.Lock()
	if b.old == val {
		b.lock.Unlock()
		return nil
	}
	*b.val = val
	b.old = val
	b.lock.Unlock()

	b.trigger()
	return nil
}

func (b *baseExternalItem[T]) Reload() error {
	return b.Set(*b.val)
}

type prefBoundBase[T bool | float64 | int | string] struct {
	base
	key   string
	get   func(string) T
	set   func(string, T)
	cache atomic.Pointer[T]
}

func (b *prefBoundBase[T]) Get() (T, error) {
	cache := b.get(b.key)
	b.cache.Store(&cache)
	return cache, nil
}

func (b *prefBoundBase[T]) Set(v T) error {
	b.set(b.key, v)

	b.lock.RLock()
	defer b.lock.RUnlock()
	b.trigger()
	return nil
}

func (b *prefBoundBase[T]) setKey(key string) {
	b.key = key
}

func (b *prefBoundBase[T]) checkForChange() {
	val := b.cache.Load()
	if val != nil && b.get(b.key) == *val {
		return
	}
	b.trigger()
}
