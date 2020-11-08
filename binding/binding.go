package binding

// DataItem is the base interface for all bindable data items.
type DataItem interface {
	AddListener(DataItemListener)
	RemoveListener(DataItemListener)
}

// DataItemListener is any object that can register for changes in a bindable DataItem.
// See NewDataItemListener to define a new listener using just an inline function.
type DataItemListener interface {
	DataChanged(DataItem)
}

// NewDataItemListener is a helper function that creates a new listener type from a simple callback function.
func NewDataItemListener(fn func(DataItem)) DataItemListener {
	return &listener{fn}
}

type listener struct {
	callback func(DataItem)
}

func (l *listener) DataChanged(i DataItem) {
	l.callback(i)
}

type base struct {
	listeners []DataItemListener
}

// AddListener allows a data listener to be informed of changes to this item.
func (b *base) AddListener(l DataItemListener) {
	b.listeners = append(b.listeners, l)
}

// RemoveListener should be called if the listener is no longer interested in being informed of data change events.
func (b *base) RemoveListener(l DataItemListener) {
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

func (b *base) trigger(i DataItem) {
	for _, listen := range b.listeners {
		listen.DataChanged(i)
	}
}
