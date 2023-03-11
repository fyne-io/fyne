package widget

import (
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2/data/binding"
)

// basicBinder stores a DataItem and a function to be called when it changes.
// It provides a convenient way to replace data and callback independently.
type basicBinder struct {
	callback atomic.Value // func(binding.DataItem)

	dataListenerPairLock sync.RWMutex
	dataListenerPair     annotatedListener // access guarded by dataListenerPairLock
}

// Bind replaces the data item whose changes are tracked by the callback function.
func (binder *basicBinder) Bind(data binding.DataItem) {
	listener := binding.NewDataListener(func() { // NB: listener captures `data` but always calls the up-to-date callback
		f := binder.callback.Load()
		if f != nil {
			f.(func(binding.DataItem))(data)
		}
	})
	data.AddListener(listener)
	listenerInfo := annotatedListener{
		data:     data,
		listener: listener,
	}

	binder.dataListenerPairLock.Lock()
	binder.unbindLocked()
	binder.dataListenerPair = listenerInfo
	binder.dataListenerPairLock.Unlock()
}

// CallWithData passes the currently bound data item as an argument to the
// provided function.
func (binder *basicBinder) CallWithData(f func(data binding.DataItem)) {
	binder.dataListenerPairLock.RLock()
	data := binder.dataListenerPair.data
	binder.dataListenerPairLock.RUnlock()
	f(data)
}

// SetCallback replaces the function to be called when the data changes.
func (binder *basicBinder) SetCallback(f func(data binding.DataItem)) {
	binder.callback.Store(f)
}

// Unbind requests the callback to be no longer called when the previously bound
// data item changes.
func (binder *basicBinder) Unbind() {
	binder.dataListenerPairLock.Lock()
	binder.unbindLocked()
	binder.dataListenerPairLock.Unlock()
}

// unbindLocked expects the caller to hold dataListenerPairLock.
func (binder *basicBinder) unbindLocked() {
	previousListener := binder.dataListenerPair
	binder.dataListenerPair = annotatedListener{nil, nil}

	if previousListener.listener == nil || previousListener.data == nil {
		return
	}
	previousListener.data.RemoveListener(previousListener.listener)
}

type annotatedListener struct {
	data     binding.DataItem
	listener binding.DataListener
}
