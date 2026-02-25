package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// List supports binding a list of values with type T.
//
// Since: 2.7
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

// ExternalList supports binding a list of values, with type T, from an external variable.
//
// Since: 2.7
type ExternalList[T any] interface {
	List[T]

	Reload() error
}

// NewList returns a bindable list of values with type T.
//
// Since: 2.7
func NewList[T any](comparator func(T, T) bool) List[T] {
	return newList[T](comparator)
}

// BindList returns a bound list of values with type T, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.7
func BindList[T any](v *[]T, comparator func(T, T) bool) ExternalList[T] {
	return bindList(v, comparator)
}

// DataList is the base interface for all bindable data lists.
//
// Since: 2.0
type DataList interface {
	DataItem
	GetItem(index int) (DataItem, error)
	Length() int
}

// BoolList supports binding a list of bool values.
//
// Since: 2.0
type BoolList = List[bool]

// ExternalBoolList supports binding a list of bool values from an external variable.
//
// Since: 2.0
type ExternalBoolList = ExternalList[bool]

// NewBoolList returns a bindable list of bool values.
//
// Since: 2.0
func NewBoolList() List[bool] {
	return newListComparable[bool]()
}

// BindBoolList returns a bound list of bool values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBoolList(v *[]bool) ExternalList[bool] {
	return bindListComparable(v)
}

// BytesList supports binding a list of []byte values.
//
// Since: 2.2
type BytesList = List[[]byte]

// ExternalBytesList supports binding a list of []byte values from an external variable.
//
// Since: 2.2
type ExternalBytesList = ExternalList[[]byte]

// NewBytesList returns a bindable list of []byte values.
//
// Since: 2.2
func NewBytesList() List[[]byte] {
	return newList(bytes.Equal)
}

// BindBytesList returns a bound list of []byte values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytesList(v *[][]byte) ExternalList[[]byte] {
	return bindList(v, bytes.Equal)
}

// FloatList supports binding a list of float64 values.
//
// Since: 2.0
type FloatList = List[float64]

// ExternalFloatList supports binding a list of float64 values from an external variable.
//
// Since: 2.0
type ExternalFloatList = ExternalList[float64]

// NewFloatList returns a bindable list of float64 values.
//
// Since: 2.0
func NewFloatList() List[float64] {
	return newListComparable[float64]()
}

// BindFloatList returns a bound list of float64 values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloatList(v *[]float64) ExternalList[float64] {
	return bindListComparable(v)
}

// IntList supports binding a list of int values.
//
// Since: 2.0
type IntList = List[int]

// ExternalIntList supports binding a list of int values from an external variable.
//
// Since: 2.0
type ExternalIntList = ExternalList[int]

// NewIntList returns a bindable list of int values.
//
// Since: 2.0
func NewIntList() List[int] {
	return newListComparable[int]()
}

// BindIntList returns a bound list of int values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindIntList(v *[]int) ExternalList[int] {
	return bindListComparable(v)
}

// RuneList supports binding a list of rune values.
//
// Since: 2.0
type RuneList = List[rune]

// ExternalRuneList supports binding a list of rune values from an external variable.
//
// Since: 2.0
type ExternalRuneList = ExternalList[rune]

// NewRuneList returns a bindable list of rune values.
//
// Since: 2.0
func NewRuneList() List[rune] {
	return newListComparable[rune]()
}

// BindRuneList returns a bound list of rune values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRuneList(v *[]rune) ExternalList[rune] {
	return bindListComparable(v)
}

// StringList supports binding a list of string values.
//
// Since: 2.0
type StringList = List[string]

// ExternalStringList supports binding a list of string values from an external variable.
//
// Since: 2.0
type ExternalStringList = ExternalList[string]

// NewStringList returns a bindable list of string values.
//
// Since: 2.0
func NewStringList() List[string] {
	return newListComparable[string]()
}

// BindStringList returns a bound list of string values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindStringList(v *[]string) ExternalList[string] {
	return bindListComparable(v)
}

// UntypedList supports binding a list of any values.
//
// Since: 2.1
type UntypedList = List[any]

// ExternalUntypedList supports binding a list of any values from an external variable.
//
// Since: 2.1
type ExternalUntypedList = ExternalList[any]

// NewUntypedList returns a bindable list of any values.
//
// Since: 2.1
func NewUntypedList() List[any] {
	return newList(func(t1, t2 any) bool { return t1 == t2 })
}

// BindUntypedList returns a bound list of any values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindUntypedList(v *[]any) ExternalList[any] {
	return bindList(v, func(t1, t2 any) bool { return t1 == t2 })
}

// URIList supports binding a list of fyne.URI values.
//
// Since: 2.1
type URIList = List[fyne.URI]

// ExternalURIList supports binding a list of fyne.URI values from an external variable.
//
// Since: 2.1
type ExternalURIList = ExternalList[fyne.URI]

// NewURIList returns a bindable list of fyne.URI values.
//
// Since: 2.1
func NewURIList() List[fyne.URI] {
	return newList(storage.EqualURI)
}

// BindURIList returns a bound list of fyne.URI values, based on the contents of the passed slice.
// If your code changes the content of the slice this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURIList(v *[]fyne.URI) ExternalList[fyne.URI] {
	return bindList(v, storage.EqualURI)
}

type listBase struct {
	base
	items []DataItem
}

// GetItem returns the DataItem at the specified index.
func (b *listBase) GetItem(i int) (DataItem, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if i < 0 || i >= len(b.items) {
		return nil, errOutOfBounds
	}

	return b.items[i], nil
}

// Length returns the number of items in this data list.
func (b *listBase) Length() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.items)
}

func (b *listBase) appendItem(i DataItem) {
	b.items = append(b.items, i)
}

func (b *listBase) deleteItem(i int) {
	b.items = append(b.items[:i], b.items[i+1:]...)
}

func newList[T any](comparator func(T, T) bool) *boundList[T] {
	return &boundList[T]{val: new([]T), comparator: comparator}
}

func newListComparable[T comparable]() *boundList[T] {
	return newList(func(t1, t2 T) bool { return t1 == t2 })
}

func newExternalList[T any](v *[]T, comparator func(T, T) bool) *boundList[T] {
	return &boundList[T]{val: v, comparator: comparator, updateExternal: true}
}

func bindList[T any](v *[]T, comparator func(T, T) bool) *boundList[T] {
	if v == nil {
		return newList(comparator)
	}

	l := newExternalList(v, comparator)
	for i := range *v {
		l.appendItem(bindListItem(v, i, l.updateExternal, comparator))
	}

	return l
}

func bindListComparable[T comparable](v *[]T) *boundList[T] {
	return bindList(v, func(t1, t2 T) bool { return t1 == t2 })
}

type boundList[T any] struct {
	listBase

	comparator     func(T, T) bool
	updateExternal bool
	val            *[]T

	parentListener func(int)
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
			item := bindListItem(l.val, i, l.updateExternal, l.comparator)

			if l.parentListener != nil {
				index := i
				item.AddListener(NewDataListener(func() {
					l.parentListener(index)
				}))
			}

			l.appendItem(item)
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
	return trigger, retErr
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
	return item.(Item[T]).Set(v)
}

func bindListItem[T any](v *[]T, i int, external bool, comparator func(T, T) bool) Item[T] {
	if external {
		ret := &boundExternalListItem[T]{old: (*v)[i]}
		ret.val = v
		ret.index = i
		ret.comparator = comparator
		return ret
	}

	return &boundListItem[T]{val: v, index: i, comparator: comparator}
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
