//go:generate go run gen.go

package binding

import (
	"errors"
	"sync"

	"fyne.io/fyne/v2"
)

var (
	errKeyNotFound = errors.New("key not found")
	errOutOfBounds = errors.New("index out of bounds")
	errParseFailed = errors.New("format did not match 1 value")

	// As an optimisation we connect any listeners asking for the same key, so that there is only 1 per preference item.
	prefBinds = make(map[fyne.Preferences]map[string]preferenceItem)
	prefLock  sync.RWMutex
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
	listeners []DataListener
	lock      sync.RWMutex
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.listeners = append(b.listeners, l)
	queueItem(l.DataChanged)
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l DataListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	for i, listen := range b.listeners {
		if listen != l {
			continue
		}

		if i == len(b.listeners)-1 {
			b.listeners = b.listeners[:len(b.listeners)-1]
		} else {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
		}
	}
}

func (b *base) trigger() {
	for _, listen := range b.listeners {
		queueItem(listen.DataChanged)
	}
}

type preferenceItem interface {
	checkForChange()
}

func ensurePreferencesAttached(p fyne.Preferences) {
	prefLock.Lock()
	defer prefLock.Unlock()
	if prefBinds[p] != nil {
		return
	}

	prefBinds[p] = make(map[string]preferenceItem)
	p.AddChangeListener(func() {
		preferencesChanged(p)
	})
}

func preferencesChanged(p fyne.Preferences) {
	prefLock.RLock()
	defer prefLock.RUnlock()

	for _, item := range prefBinds[p] {
		item.checkForChange()
	}
}
