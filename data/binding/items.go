package binding

import (
	"bytes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// Item supports binding any type T generically.
//
// Since: 2.6
type Item[T any] interface {
	DataItem
	Get() (T, error)
	Set(T) error
}

// ExternalItem supports binding any external value of type T.
//
// Since: 2.6
type ExternalItem[T any] interface {
	Item[T]
	Reload() error
}

// NewItem returns a bindable value of type T that is managed internally.
//
// Since: 2.6
func NewItem[T any](comparator func(T, T) bool) Item[T] {
	return &item[T]{val: new(T), comparator: comparator}
}

// BindItem returns a new bindable value that controls the contents of the provided variable of type T.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.6
func BindItem[T any](val *T, comparator func(T, T) bool) ExternalItem[T] {
	if val == nil {
		val = new(T) // never allow a nil value pointer
	}
	b := &externalItem[T]{}
	b.comparator = comparator
	b.val = val
	b.old = *val
	return b
}

// Bool supports binding a bool value.
//
// Since: 2.0
type Bool = Item[bool]

// ExternalBool supports binding a bool value to an external value.
//
// Since: 2.0
type ExternalBool = ExternalItem[bool]

// NewBool returns a bindable bool value that is managed internally.
//
// Since: 2.0
func NewBool() Bool {
	return newItemComparable[bool]()
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBool(v *bool) ExternalBool {
	return bindExternalComparable(v)
}

// Bytes supports binding a []byte value.
//
// Since: 2.2
type Bytes = Item[[]byte]

// ExternalBytes supports binding a []byte value to an external value.
//
// Since: 2.2
type ExternalBytes = ExternalItem[[]byte]

// NewBytes returns a bindable []byte value that is managed internally.
//
// Since: 2.2
func NewBytes() Bytes {
	return NewItem(bytes.Equal)
}

// BindBytes returns a new bindable value that controls the contents of the provided []byte variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytes(v *[]byte) ExternalBytes {
	return BindItem(v, bytes.Equal)
}

// Float supports binding a float64 value.
//
// Since: 2.0
type Float = Item[float64]

// ExternalFloat supports binding a float64 value to an external value.
//
// Since: 2.0
type ExternalFloat = ExternalItem[float64]

// NewFloat returns a bindable float64 value that is managed internally.
//
// Since: 2.0
func NewFloat() Float {
	return newItemComparable[float64]()
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloat(v *float64) ExternalFloat {
	return bindExternalComparable(v)
}

// Int supports binding a int value.
//
// Since: 2.0
type Int = Item[int]

// ExternalInt supports binding a int value to an external value.
//
// Since: 2.0
type ExternalInt = ExternalItem[int]

// NewInt returns a bindable int value that is managed internally.
//
// Since: 2.0
func NewInt() Int {
	return newItemComparable[int]()
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindInt(v *int) ExternalInt {
	return bindExternalComparable(v)
}

// Rune supports binding a rune value.
//
// Since: 2.0
type Rune = Item[rune]

// ExternalRune supports binding a rune value to an external value.
//
// Since: 2.0
type ExternalRune = ExternalItem[rune]

// NewRune returns a bindable rune value that is managed internally.
//
// Since: 2.0
func NewRune() Rune {
	return newItemComparable[rune]()
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRune(v *rune) ExternalRune {
	return bindExternalComparable(v)
}

// String supports binding a string value.
//
// Since: 2.0
type String = Item[string]

// ExternalString supports binding a string value to an external value.
//
// Since: 2.0
type ExternalString = ExternalItem[string]

// NewString returns a bindable string value that is managed internally.
//
// Since: 2.0
func NewString() String {
	return newItemComparable[string]()
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindString(v *string) ExternalString {
	return bindExternalComparable(v)
}

// URI supports binding a fyne.URI value.
//
// Since: 2.1
type URI = Item[fyne.URI]

// ExternalURI supports binding a fyne.URI value to an external value.
//
// Since: 2.1
type ExternalURI = ExternalItem[fyne.URI]

// NewURI returns a bindable fyne.URI value that is managed internally.
//
// Since: 2.1
func NewURI() URI {
	return NewItem(storage.EqualURI)
}

// BindURI returns a new bindable value that controls the contents of the provided fyne.URI variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURI(v *fyne.URI) ExternalURI {
	return BindItem(v, storage.EqualURI)
}

func newItemComparable[T bool | float64 | int | rune | string]() Item[T] {
	return NewItem[T](func(a, b T) bool { return a == b })
}

type item[T any] struct {
	base

	comparator func(T, T) bool
	val        *T
}

func (b *item[T]) Get() (T, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.val == nil {
		return *new(T), nil
	}
	return *b.val, nil
}

func (b *item[T]) Set(val T) error {
	b.lock.Lock()
	equal := b.comparator(*b.val, val)
	*b.val = val
	b.lock.Unlock()

	if !equal {
		b.trigger()
	}

	return nil
}

func bindExternalComparable[T bool | float64 | int | rune | string](val *T) ExternalItem[T] {
	return BindItem(val, func(t1, t2 T) bool { return t1 == t2 })
}

type externalItem[T any] struct {
	item[T]

	old T
}

func (b *externalItem[T]) Set(val T) error {
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

func (b *externalItem[T]) Reload() error {
	return b.Set(*b.val)
}
