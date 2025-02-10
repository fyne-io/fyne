//go:generate go run gen.go

// Package binding provides support for binding data to widgets.
// All APIs in the binding package are safe to invoke directly from any goroutine.
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
// All APIs on bindable data items are safe to invoke directly fron any goroutine.
//
// Since: 2.0
type DataItem interface {
	// AddListener attaches a new change listener to this DataItem.
	// Listeners are called each time the data inside this DataItem changes.
	// Additionally, the listener will be triggered upon successful connection to get the current value.
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
	listeners []DataListener

	lock sync.RWMutex
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l DataListener) {
	fyne.Do(func() {
		b.listeners = append(b.listeners, l)
		l.DataChanged()
	})
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l DataListener) {
	fyne.Do(func() {
		for i, listener := range b.listeners {
			if listener == l {
				// Delete without preserving order:
				lastIndex := len(b.listeners) - 1
				b.listeners[i] = b.listeners[lastIndex]
				b.listeners[lastIndex] = nil
				b.listeners = b.listeners[:lastIndex]
				return
			}
		}
	})
}

func (b *base) trigger() {
	fyne.Do(b.triggerFromMain)
}

func (b *base) triggerFromMain() {
	for _, listen := range b.listeners {
		listen.DataChanged()
	}
}

// Untyped supports binding an any value.
//
// Since: 2.1
type Untyped = Item[any]

// NewUntyped returns a bindable any value that is managed internally.
//
// Since: 2.1
func NewUntyped() Untyped {
	return NewItem(func(a1, a2 any) bool { return a1 == a2 })
}

// ExternalUntyped supports binding a any value to an external value.
//
// Since: 2.1
type ExternalUntyped = ExternalItem[any]

// BindUntyped returns a bindable any value that is bound to an external type.
// The parameter must be a pointer to the type you wish to bind.
//
// Since: 2.1
func BindUntyped(v any) ExternalUntyped {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		fyne.LogError("Invalid type passed to BindUntyped, must be a pointer", nil)
		v = nil
	}

	if v == nil {
		v = new(any) // never allow a nil value pointer
	}

	b := &boundExternalUntyped{}
	b.val = reflect.ValueOf(v).Elem()
	b.old = b.val.Interface()
	return b
}

type boundExternalUntyped struct {
	base

	val reflect.Value
	old any
}

func (b *boundExternalUntyped) Get() (any, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.val.Interface(), nil
}

func (b *boundExternalUntyped) Set(val any) error {
	b.lock.Lock()
	if b.old == val {
		b.lock.Unlock()
		return nil
	}
	b.val.Set(reflect.ValueOf(val))
	b.old = val
	b.lock.Unlock()

	b.trigger()
	return nil
}

func (b *boundExternalUntyped) Reload() error {
	return b.Set(b.val.Interface())
}
