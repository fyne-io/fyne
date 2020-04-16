package binding

import (
	"sync"

	"fyne.io/fyne"
)

// Binding is the base interface of the Data Binding API.
type Binding interface {
	// AddListener adds the given listener to the binding.
	AddListener(Notifiable)
	// DeleteListener removes the given listener from the binding.
	DeleteListener(Notifiable)
	// Update tells the binding that the data changed.
	Update()
}

// List defines a data binding wrapping a list of bindings.
type List interface {
	Binding
	// Length returns the length of the list.
	Length() int
	// Get returns the binding at the given index.
	Get(int) Binding
}

// Map defines a data binding wrapping a list of bindings.
type Map interface {
	Binding
	// Length returns the length of the map.
	Length() int
	// Get returns the binding for the given key.
	Get(string) (Binding, bool)
}

// Base is the base implementation of a data binding and handles adding, deleting, and notifying listeners.
type Base struct {
	Binding
	sync.RWMutex
	listeners []Notifiable
}

// AddListener adds the given listener to the binding.
func (b *Base) AddListener(listener Notifiable) {
	b.Lock()
	defer b.Unlock()
	for _, l := range b.listeners {
		if l == listener {
			fyne.LogError("Listener already added to this Binding", nil)
			return
		}
	}
	b.listeners = append(b.listeners, listener)
	// Call the listener with the current state
	go listener.Notify(b)
}

// DeleteListener removes the given listener from the binding.
func (b *Base) DeleteListener(listener Notifiable) {
	b.Lock()
	defer b.Unlock()
	var listeners []Notifiable
	for _, l := range b.listeners {
		if l != listener {
			listeners = append(listeners, l)
		}
	}
	b.listeners = listeners
}

// Update should be called whenever the data changed.
// This will enqueue the binding for the Updater.
func (b *Base) Update() {
	if updaterQueue == nil {
		panic("binding.Start() not called")
	}
	select {
	case updaterQueue <- b.notifyListeners:
	default:
		fyne.LogError("UpdateQueue full", nil)
	}
}

// notifyListeners should only be called by Updater
func (b *Base) notifyListeners() {
	b.RLock()
	defer b.RUnlock()
	for _, l := range b.listeners {
		l.Notify(b)
	}
}

var updaterQueue chan func()

// Start initializes the binding queue and starts the updater.
// Start should be called once when an app is created.
func Start() {
	if updaterQueue != nil {
		Stop()
	}
	updaterQueue = make(chan func(), 1024)
	go func() {
		for fn := range updaterQueue {
			fn()
		}
	}()
}

// Stop closes the binding queue causing the updater to quit.
// Stop should be called once when an app is quitting.
func Stop() {
	if updaterQueue == nil {
		return
	}
	close(updaterQueue)
	updaterQueue = nil
}
