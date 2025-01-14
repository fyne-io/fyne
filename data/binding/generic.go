package binding

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
