package async

// The size of Func, Interface, and CanvasObject are all less than 16 bytes.
// A CPU cache line (L2) is 256 bytes, thus it can hold 16 entity references.
const maxEntitiesPerCPUCacheLine = 16

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
		// We make the channels fit into CPU cache lines, which may reduce cache misses.
		in:    make(chan T, maxEntitiesPerCPUCacheLine),
		out:   make(chan T, maxEntitiesPerCPUCacheLine),
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
				ch.q[0] = *new(T) // de-reference earlier to help GC (use clear() when Go 1.21 is base)
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
			ch.q[0] = *new(T) // de-reference earlier to help GC (use clear() when Go 1.21 is base)
			ch.q = ch.q[1:]
		default:
		}
	}
	close(ch.out)
	close(ch.close)
}
