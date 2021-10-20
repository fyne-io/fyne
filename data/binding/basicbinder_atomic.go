//go:build go1.17
// +build go1.17

package binding

import (
	"sync/atomic"
)

// BasicBinder stores a DataItem and a function to be called when it changes.
// It provides a convenient way to replace data and callback independently.
type BasicBinder struct {
	callback         atomic.Value // of type wrappedCallback
	dataListenerPair atomic.Value // of type annotatedListener
}

// SetCallback replaces the function to be called when the data changes.
func (binder *BasicBinder) SetCallback(f func(data DataItem)) {
	binder.callback.Store(wrappedCallback{
		notNil: (f != nil),
		f:      f,
	})
}

// Bind replaces the data item whose changes are tracked by the callback function.
func (binder *BasicBinder) Bind(data DataItem) {
	listener := NewDataListener(func() { // NB: listener captures `data` but always calls the up-to-date callback
		callbackInterface := binder.callback.Load()
		if callbackInterface == nil {
			return
		}
		callback := callbackInterface.(wrappedCallback)
		if callback.notNil {
			callback.f(data)
		}
	})
	data.AddListener(listener)
	listenerInfo := annotatedListener{
		data:     data,
		listener: listener,
	}

	binder.Unbind()
	nilPair := annotatedListener{nil, nil}
	for !binder.dataListenerPair.CompareAndSwap(nilPair, listenerInfo) {
		binder.Unbind() // keep unbinding until binder.dataListenerPair is a nil pair
	}
}

// Unbind requests the callback to be no longer called when the previously bound
// data item changes.
func (binder *BasicBinder) Unbind() {
	nilPair := annotatedListener{nil, nil}
	state := binder.dataListenerPair.Swap(nilPair)
	if state == nil {
		return // this will happen only for the very first operation with dataListenerPair
	}
	previousListener := state.(annotatedListener)
	if previousListener.listener == nil || previousListener.data == nil {
		return
	}
	previousListener.data.RemoveListener(previousListener.listener)
}

type annotatedListener struct {
	data     DataItem
	listener DataListener
}

// wrappedCallback only exists since atomic.Value cannot store nil callback value, so
// we have to represent nil callback as {notNil: false, f: nil}.
type wrappedCallback struct {
	notNil bool // must always equal (f != nil)
	f      func(DataItem)
}
