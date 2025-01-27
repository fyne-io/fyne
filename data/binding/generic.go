package binding

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
)

func newBaseItem[T any](comparator func(T, T) bool) *baseItem[T] {
	return &baseItem[T]{val: new(T), comparator: comparator}
}

func newBaseItemComparable[T bool | float64 | int | rune | string]() *baseItem[T] {
	return newBaseItem[T](func(a, b T) bool { return a == b })
}

type baseItem[T any] struct {
	base

	comparator func(T, T) bool
	val        *T
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
	equal := b.comparator(*b.val, val)
	*b.val = val
	b.lock.Unlock()

	if !equal {
		b.trigger()
	}

	return nil
}

func baseBindExternal[T any](val *T, comparator func(T, T) bool) *baseExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &baseExternalItem[T]{}
	b.comparator = comparator
	b.val = val
	b.old = *val
	return b
}

func baseBindExternalComparable[T bool | float64 | int | rune | string](val *T) *baseExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &baseExternalItem[T]{}
	b.comparator = func(a, b T) bool { return a == b }
	b.val = val
	b.old = *val
	return b
}

type baseExternalItem[T any] struct {
	baseItem[T]

	old T
}

func (b *baseExternalItem[T]) Set(val T) error {
	b.lock.Lock()
	if b.comparator(b.old, val) {
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

type genericItem[T any] interface {
	DataItem
	Get() (T, error)
	Set(T) error
}

func lookupExistingBinding[T any](key string, p fyne.Preferences) (genericItem[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(genericItem[T]); ok {
			return l, ok
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	return nil, false
}

type List[T any] interface {
	DataList

	Append(value T) error
	Get() ([]T, error)
	GetValue(index int) (T, error)
	Prepend(value T) error
	Remove(value T) error
	Set(list []T) error
	SetValue(index int, value T) error
}

func newList[T any](comparator func(T, T) bool) *boundList[T] {
	return &boundList[T]{val: new([]T), comparator: comparator}
}

func newListComparable[T bool | float64 | int | rune | string]() *boundList[T] {
	return &boundList[T]{val: new([]T), comparator: func(t1, t2 T) bool { return t1 == t2 }}
}

type boundList[T any] struct {
	listBase

	comparator     func(T, T) bool
	updateExternal bool
	val            *[]T
}

func (l *boundList[T]) Append(val T) error {
	l.lock.Lock()
	*l.val = append(*l.val, val)

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Get() ([]T, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return *l.val, nil
}

func (l *boundList[T]) GetValue(i int) (T, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if i < 0 || i >= l.Length() {
		return *new(T), errOutOfBounds
	}

	return (*l.val)[i], nil
}

func (l *boundList[T]) Prepend(val T) error {
	l.lock.Lock()
	*l.val = append([]T{val}, *l.val...)

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Reload() error {
	l.lock.Lock()
	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Remove(val T) error {
	l.lock.Lock()

	v := *l.val
	if len(v) == 0 {
		l.lock.Unlock()
		return nil
	}
	if l.comparator(v[0], val) {
		*l.val = v[1:]
	} else if l.comparator(v[len(v)-1], val) {
		*l.val = v[:len(v)-1]
	} else {
		id := -1
		for i, v := range v {
			if l.comparator(v, val) {
				id = i
				break
			}
		}

		if id == -1 {
			l.lock.Unlock()
			return nil
		}
		*l.val = append(v[:id], v[id+1:]...)
	}

	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) Set(v []T) error {
	l.lock.Lock()
	*l.val = v
	trigger, err := l.doReload()
	l.lock.Unlock()

	if trigger {
		l.trigger()
	}

	return err
}

func (l *boundList[T]) doReload() (trigger bool, retErr error) {
	oldLen := len(l.items)
	newLen := len(*l.val)
	if oldLen > newLen {
		for i := oldLen - 1; i >= newLen; i-- {
			l.deleteItem(i)
		}
		trigger = true
	} else if oldLen < newLen {
		for i := oldLen; i < newLen; i++ {
			l.appendItem(bindListItem(l.val, i, l.updateExternal, l.comparator))
		}
		trigger = true
	}

	for i, item := range l.items {
		if i > oldLen || i > newLen {
			break
		}

		var err error
		if l.updateExternal {
			err = item.(*boundExternalListItem[T]).setIfChanged((*l.val)[i])
		} else {
			err = item.(*boundListItem[T]).doSet((*l.val)[i])
		}
		if err != nil {
			retErr = err
		}
	}
	return
}

func (l *boundList[T]) SetValue(i int, v T) error {
	l.lock.RLock()
	len := l.Length()
	l.lock.RUnlock()

	if i < 0 || i >= len {
		return errOutOfBounds
	}

	l.lock.Lock()
	(*l.val)[i] = v
	l.lock.Unlock()

	item, err := l.GetItem(i)
	if err != nil {
		return err
	}
	return item.(genericItem[T]).Set(v)
}

func bindListItem[T any](v *[]T, i int, external bool, comparator func(T, T) bool) genericItem[T] {
	if external {
		ret := &boundExternalListItem[T]{old: (*v)[i]}
		ret.val = v
		ret.index = i
		ret.comparator = comparator
		return ret
	}

	return &boundListItem[T]{val: v, index: i, comparator: comparator}
}

func bindListItemComparable[T bool | float64 | int | rune | string](v *[]T, i int, external bool) genericItem[T] {
	return bindListItem(v, i, external, func(t1, t2 T) bool { return t1 == t2 })
}

type boundListItem[T any] struct {
	base

	comparator func(T, T) bool
	val        *[]T
	index      int
}

func (b *boundListItem[T]) Get() (T, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.index < 0 || b.index >= len(*b.val) {
		return *new(T), errOutOfBounds
	}

	return (*b.val)[b.index], nil
}

func (b *boundListItem[T]) Set(val T) error {
	return b.doSet(val)
}

func (b *boundListItem[T]) doSet(val T) error {
	b.lock.Lock()
	(*b.val)[b.index] = val
	b.lock.Unlock()

	b.trigger()
	return nil
}

type boundExternalListItem[T any] struct {
	boundListItem[T]

	old T
}

func (b *boundExternalListItem[T]) setIfChanged(val T) error {
	b.lock.Lock()
	if b.comparator(val, b.old) {
		b.lock.Unlock()
		return nil
	}
	(*b.val)[b.index] = val
	b.old = val

	b.lock.Unlock()
	b.trigger()
	return nil
}
