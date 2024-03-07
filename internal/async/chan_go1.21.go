//go:build go1.21

package async

import "fyne.io/fyne/v2"

// UnboundedFuncChan is a channel with an unbounded buffer for caching
// Func objects. A channel must be closed via Close method
type UnboundedFuncChan = UnboundedChan[func()]

// NewUnboundedFuncChan returns a unbounded channel, of func(), with unlimited capacity.
func NewUnboundedFuncChan() *UnboundedFuncChan {
	return NewUnboundedChan[func()]()
}

// UnboundedInterfaceChan is a channel with an unbounded buffer for caching
// Interface objects. A channel must be closed via Close method.
type UnboundedInterfaceChan = UnboundedChan[any]

// NewUnboundedInterfaceChan returns a unbounded channel, of any type, with unlimited capacity.
func NewUnboundedInterfaceChan() *UnboundedInterfaceChan {
	return NewUnboundedChan[any]()
}

// UnboundedCanvasObjectChan is a channel with an unbounded buffer for caching
// CanvasObject objects. A channel must be closed via Close method.
type UnboundedCanvasObjectChan = UnboundedChan[fyne.CanvasObject]

// NewUnboundedCanvasObjectChan returns a unbounded channel, of canvas objects, with unlimited capacity.
func NewUnboundedCanvasObjectChan() *UnboundedChan[fyne.CanvasObject] {
	return NewUnboundedChan[fyne.CanvasObject]()
}

// UnboundedChan is a channel with an unbounded buffer for caching
// Func objects. A channel must be closed via Close method.
type UnboundedChan[T any] struct {
	in, out chan T
	close   chan struct{}
	q       []T
}

// NewUnboundedChan returns a unbounded channel with unlimited capacity.
func NewUnboundedChan[T any]() *UnboundedChan[T] {
	ch := &UnboundedChan[T]{
		// The size of Func, Interface, and CanvasObject are all less than 16 bytes, we use 16 to fit
		// a CPU cache line (L2, 256 Bytes), which may reduce cache misses.
		in:    make(chan T, 16),
		out:   make(chan T, 16),
		close: make(chan struct{}),
	}
	go ch.processing()
	return ch
}

// In returns the send channel of the given channel, which can be used to
// send values to the channel.
func (ch *UnboundedChan[T]) In() chan<- T { return ch.in }

// Out returns the receive channel of the given channel, which can be used
// to receive values from the channel.
func (ch *UnboundedChan[T]) Out() <-chan T { return ch.out }

// Close closes the channel.
func (ch *UnboundedChan[T]) Close() { ch.close <- struct{}{} }

func (ch *UnboundedChan[T]) processing() {
	// This is a preallocation of the internal unbounded buffer.
	// The size is randomly picked. But if one changes the size, the
	// reallocation size at the subsequent for loop should also be
	// changed too. Furthermore, there is no memory leak since the
	// queue is garbage collected.
	ch.q = make([]T, 0, 1<<10)
	for {
		select {
		case e, ok := <-ch.in:
			if !ok {
				// We don't want the input channel be accidentally closed
				// via close() instead of Close(). If that happens, it is
				// a misuse, do a panic as warning.
				panic("async: misuse of unbounded channel, In() was closed")
			}
			ch.q = append(ch.q, e)
		case <-ch.close:
			ch.closed()
			return
		}
		for len(ch.q) > 0 {
			select {
			case ch.out <- ch.q[0]:
				clear(ch.q[:1]) // de-reference earlier to help GC
				ch.q = ch.q[1:]
			case e, ok := <-ch.in:
				if !ok {
					// We don't want the input channel be accidentally closed
					// via close() instead of Close(). If that happens, it is
					// a misuse, do a panic as warning.
					panic("async: misuse of unbounded channel, In() was closed")
				}
				ch.q = append(ch.q, e)
			case <-ch.close:
				ch.closed()
				return
			}
		}
		// If the remaining capacity is too small, we prefer to
		// reallocate the entire buffer.
		if cap(ch.q) < 1<<5 {
			ch.q = make([]T, 0, 1<<10)
		}
	}
}

func (ch *UnboundedChan[T]) closed() {
	close(ch.in)
	for e := range ch.in {
		ch.q = append(ch.q, e)
	}
	for len(ch.q) > 0 {
		select {
		case ch.out <- ch.q[0]:
			clear(ch.q[:1]) // de-reference earlier to help GC
			ch.q = ch.q[1:]
		default:
		}
	}
	close(ch.out)
	close(ch.close)
}
