//go:build !go1.17
// +build !go1.17

package binding

import "sync"

// BasicBinder stores a DataItem and a function to be called when it changes.
// It provides a convenient way to replace data and callback independently.
type BasicBinder struct {
	callbackLock         sync.RWMutex
	callback             func(DataItem) // access guarded by callbackLock
	dataListenerPairLock sync.RWMutex
	dataListenerPair     annotatedListener // access guarded by dataListenerPairLock
}

// SetCallback replaces the function to be called when the data changes.
func (binder *BasicBinder) SetCallback(f func(data DataItem)) {
	binder.callbackLock.Lock()
	binder.callback = f
	binder.callbackLock.Unlock()
}

// Bind replaces the data item whose changes are tracked by the callback function.
func (binder *BasicBinder) Bind(data DataItem) {
	listener := NewDataListener(func() { // NB: listener captures `data` but always calls the up-to-date callback
		binder.callbackLock.RLock()
		f := binder.callback
		binder.callbackLock.RUnlock()
		if f != nil {
			f(data)
		}
	})
	data.AddListener(listener)
	listenerInfo := annotatedListener{
		data:     data,
		listener: listener,
	}

	binder.Unbind()
	nilPair := annotatedListener{nil, nil}
	for {
		binder.dataListenerPairLock.Lock()
		if binder.dataListenerPair != nilPair {
			binder.dataListenerPairLock.Unlock()
			binder.Unbind()
			continue // keep unbinding until dataListenerPair is a nil pair
		}
		binder.dataListenerPair = listenerInfo
		binder.dataListenerPairLock.Unlock()
		break
	}
}

// Unbind requests the callback to be no longer called when the previously bound
// data item changes.
func (binder *BasicBinder) Unbind() {
	nilPair := annotatedListener{nil, nil}

	binder.dataListenerPairLock.Lock()
	previousListener := binder.dataListenerPair
	binder.dataListenerPair = nilPair
	binder.dataListenerPairLock.Unlock()

	if previousListener.listener == nil || previousListener.data == nil {
		return
	}
	previousListener.data.RemoveListener(previousListener.listener)
}

type annotatedListener struct {
	data     DataItem
	listener DataListener
}
