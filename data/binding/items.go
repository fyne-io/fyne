package binding

import (
	"bytes"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// Bool supports binding a bool value.
//
// Since: 2.0
type Bool interface {
	DataItem
	Get() (bool, error)
	Set(bool) error
}

// ExternalBool supports binding a bool value to an external value.
//
// Since: 2.0
type ExternalBool interface {
	Bool
	Reload() error
}

// NewBool returns a bindable bool value that is managed internally.
//
// Since: 2.0
func NewBool() Bool {
	return newBaseItemComparable[bool]()
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindBool(v *bool) ExternalBool {
	return baseBindExternalComparable(v)
}

// Bytes supports binding a []byte value.
//
// Since: 2.2
type Bytes interface {
	DataItem
	Get() ([]byte, error)
	Set([]byte) error
}

// ExternalBytes supports binding a []byte value to an external value.
//
// Since: 2.2
type ExternalBytes interface {
	Bytes
	Reload() error
}

// NewBytes returns a bindable []byte value that is managed internally.
//
// Since: 2.2
func NewBytes() Bytes {
	return newBaseItem(bytes.Equal)
}

// BindBytes returns a new bindable value that controls the contents of the provided []byte variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.2
func BindBytes(v *[]byte) ExternalBytes {
	return baseBindExternal(v, bytes.Equal)
}

// Float supports binding a float64 value.
//
// Since: 2.0
type Float interface {
	DataItem
	Get() (float64, error)
	Set(float64) error
}

// ExternalFloat supports binding a float64 value to an external value.
//
// Since: 2.0
type ExternalFloat interface {
	Float
	Reload() error
}

// NewFloat returns a bindable float64 value that is managed internally.
//
// Since: 2.0
func NewFloat() Float {
	return newBaseItemComparable[float64]()
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindFloat(v *float64) ExternalFloat {
	return baseBindExternalComparable(v)
}

// Int supports binding a int value.
//
// Since: 2.0
type Int interface {
	DataItem
	Get() (int, error)
	Set(int) error
}

// ExternalInt supports binding a int value to an external value.
//
// Since: 2.0
type ExternalInt interface {
	Int
	Reload() error
}

// NewInt returns a bindable int value that is managed internally.
//
// Since: 2.0
func NewInt() Int {
	return newBaseItemComparable[int]()
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindInt(v *int) ExternalInt {
	return baseBindExternalComparable(v)
}

// Rune supports binding a rune value.
//
// Since: 2.0
type Rune interface {
	DataItem
	Get() (rune, error)
	Set(rune) error
}

// ExternalRune supports binding a rune value to an external value.
//
// Since: 2.0
type ExternalRune interface {
	Rune
	Reload() error
}

// NewRune returns a bindable rune value that is managed internally.
//
// Since: 2.0
func NewRune() Rune {
	return newBaseItemComparable[rune]()
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindRune(v *rune) ExternalRune {
	return baseBindExternalComparable(v)
}

// String supports binding a string value.
//
// Since: 2.0
type String interface {
	DataItem
	Get() (string, error)
	Set(string) error
}

// ExternalString supports binding a string value to an external value.
//
// Since: 2.0
type ExternalString interface {
	String
	Reload() error
}

// NewString returns a bindable string value that is managed internally.
//
// Since: 2.0
func NewString() String {
	return newBaseItemComparable[string]()
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.0
func BindString(v *string) ExternalString {
	return baseBindExternalComparable(v)
}

// URI supports binding a fyne.URI value.
//
// Since: 2.1
type URI interface {
	DataItem
	Get() (fyne.URI, error)
	Set(fyne.URI) error
}

// ExternalURI supports binding a fyne.URI value to an external value.
//
// Since: 2.1
type ExternalURI interface {
	URI
	Reload() error
}

// NewURI returns a bindable fyne.URI value that is managed internally.
//
// Since: 2.1
func NewURI() URI {
	return newBaseItem(storage.EqualURI)
}

// BindURI returns a new bindable value that controls the contents of the provided fyne.URI variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
//
// Since: 2.1
func BindURI(v *fyne.URI) ExternalURI {
	return baseBindExternal(v, storage.EqualURI)
}

type bindableItem[T any] interface {
	DataItem
	Get() (T, error)
	Set(T) error
}

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

func lookupExistingBinding[T any](key string, p fyne.Preferences) (bindableItem[T], bool) {
	binds := prefBinds.getBindings(p)
	if binds == nil {
		return nil, false
	}

	if listen, ok := binds.Load(key); listen != nil && ok {
		if l, ok := listen.(bindableItem[T]); ok {
			return l, ok
		}
		fyne.LogError(keyTypeMismatchError+key, nil)
	}

	return nil, false
}
