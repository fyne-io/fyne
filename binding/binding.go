package binding

// Binding is the base interface of the Data Binding API.
type Binding interface {
	// Update notifies all listeners after a change.
	Update()
}

// Collection is a binding holding a collection of multiple bindings.
type Collection interface {
	Binding
	// Length returns the length of the collection.
	Length() int
	// Listen returns a channel through which collection updates will be published.
	Listen() <-chan int
	// OnUpdate calls the given function whenever the collection updates.
	OnUpdate(func(int))
}

// List defines a data binding wrapping a list of bindings.
type List interface {
	Collection
	// Get returns the binding at the given index.
	Get(int) Binding
}

// Map defines a data binding wrapping a list of bindings.
type Map interface {
	Collection
	// Get returns the binding for the given key.
	Get(string) (Binding, bool)
}
