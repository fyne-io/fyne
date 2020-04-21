package binding

// Binding is the base interface of the Data Binding API.
type Binding interface {
	// Update notifies all listeners after a change.
	Update()
}

// List defines a data binding wrapping a list of bindings.
type List interface {
	Binding
	// Length returns the length of the list.
	Length() int
	// Get returns the binding at the given index.
	Get(int) Binding
	// Listen returns a channel through which list length updates will be published.
	Listen() <-chan int
}

// Map defines a data binding wrapping a list of bindings.
type Map interface {
	Binding
	// Length returns the length of the map.
	Length() int
	// Get returns the binding for the given key.
	Get(string) (Binding, bool)
	// Listen returns a channel through which map length updates will be published.
	Listen() <-chan int
}
