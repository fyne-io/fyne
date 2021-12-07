//go:generate go run gen.go

// Package binding provides support for binding data to widgets.
package binding

import (
	"errors"
	"reflect"
	"sync"

	"fyne.io/fyne/v2"
)

var (
	errKeyNotFound = errors.New("key not found")
	errOutOfBounds = errors.New("index out of bounds")
	errParseFailed = errors.New("format did not match 1 value")

	// As an optimisation we connect any listeners asking for the same key, so that there is only 1 per preference item.
	prefBinds = newPreferencesMap()
)

// DataItem is the base interface for all bindable data items.
//
// Since: 2.0
type DataItem interface {
	// AddListener attaches a new change listener to this DataItem.
	// Listeners are called each time the data inside this DataItem changes.
	// Additionally the listener will be triggered upon successful connection to get the current value.
	AddListener(DataListener)
	// RemoveListener will detach the specified change listener from the DataItem.
	// Disconnected listener will no longer be triggered when changes occur.
	RemoveListener(DataListener)
}

// DataListener is any object that can register for changes in a bindable DataItem.
// See NewDataListener to define a new listener using just an inline function.
//
// Since: 2.0
type DataListener interface {
	DataChanged()
}

// NewDataListener is a helper function that creates a new listener type from a simple callback function.
//
// Since: 2.0
func NewDataListener(fn func()) DataListener {
	return &listener{fn}
}

type listener struct {
	callback func()
}

func (l *listener) DataChanged() {
	l.callback()
}

type base struct {
	listeners sync.Map // map[DataListener]bool

	lock sync.RWMutex
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l DataListener) {
	b.listeners.Store(l, true)
	queueItem(l.DataChanged)
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l DataListener) {
	b.listeners.Delete(l)
}

func (b *base) trigger() {
	b.listeners.Range(func(key, _ interface{}) bool {
		queueItem(key.(DataListener).DataChanged)
		return true
	})
}

// Untyped supports binding a interface{} value.
//
// Since: 2.1
type Untyped interface {
	DataItem
	Get() (interface{}, error)
	Set(interface{}) error
}

// NewUntyped returns a bindable interface{} value that is managed internally.
//
// Since: 2.1
func NewUntyped() Untyped {
	var blank interface{} = nil
	v := &blank
	return &boundUntyped{val: reflect.ValueOf(v).Elem()}
}

type boundUntyped struct {
	base

	val reflect.Value
}

func (b *boundUntyped) Get() (interface{}, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.val.Interface(), nil
}

func (b *boundUntyped) Set(val interface{}) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.val.Interface() == val {
		return nil
	}

	b.val.Set(reflect.ValueOf(val))

	b.trigger()
	return nil
}

// ExternalUntyped supports binding a interface{} value to an external value.
//
// Since: 2.1
type ExternalUntyped interface {
	Untyped
	Reload() error
}

// BindUntyped returns a bindable interface{} value that is bound to an external type.
// The parameter must be a pointer to the type you wish to bind.
//
// Since: 2.1
func BindUntyped(v interface{}) ExternalUntyped {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		fyne.LogError("Invalid type passed to BindUntyped, must be a pointer", nil)
		v = nil
	}

	if v == nil {
		var blank interface{}
		v = &blank // never allow a nil value pointer
	}

	b := &boundExternalUntyped{}
	b.val = reflect.ValueOf(v).Elem()
	b.old = b.val.Interface()
	return b
}

type boundExternalUntyped struct {
	boundUntyped

	old interface{}
}

func (b *boundExternalUntyped) Set(val interface{}) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return nil
	}
	b.val.Set(reflect.ValueOf(val))
	b.old = val

	b.trigger()
	return nil
}

func (b *boundExternalUntyped) Reload() error {
	return b.Set(b.val.Interface())
}
