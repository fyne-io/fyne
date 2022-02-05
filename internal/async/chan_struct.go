package async

// UnboundedStructChan is a channel with an unbounded buffer for caching
// struct{} objects. This implementation is a specialized version that
// optimizes for struct{} objects than other types. A channel must be
// closed via Close method.
type UnboundedStructChan struct {
	in, out, close chan struct{}
	n              uint64
}

// NewUnboundedStructChan returns a unbounded channel with unlimited capacity.
func NewUnboundedStructChan() *UnboundedStructChan {
	ch := &UnboundedStructChan{
		// The size of Struct is less than 16 bytes, we use 16 to fit
		// a CPU cache line (L2, 256 Bytes), which may reduce cache misses.
		in:    make(chan struct{}, 16),
		out:   make(chan struct{}, 16),
		close: make(chan struct{}),
	}
	go ch.processing()
	return ch
}

// In returns a send-only channel that can be used to send values
// to the channel.
func (ch *UnboundedStructChan) In() chan<- struct{} { return ch.in }

// Out returns a receive-only channel that can be used to receive
// values from the channel.
func (ch *UnboundedStructChan) Out() <-chan struct{} { return ch.out }

// Close closes the channel.
func (ch *UnboundedStructChan) Close() { ch.close <- struct{}{} }

func (ch *UnboundedStructChan) processing() {
	for {
		select {
		case _, ok := <-ch.in:
			if !ok {
				// We don't want the input channel be accidentally closed
				// via close() instead of Close(). If that happens, it is
				// a misuse, do a panic as warning.
				panic("async: misuse of unbounded channel, In() was closed")
			}
			ch.n++
		case <-ch.close:
			ch.closed()
			return
		}
		for ch.n > 0 {
			select {
			case ch.out <- struct{}{}:
				ch.n--
			case _, ok := <-ch.in:
				if !ok {
					// We don't want the input channel be accidentally closed
					// via close() instead of Close(). If that happens, it is
					// a misuse, do a panic as warning.
					panic("async: misuse of unbounded channel, In() was closed")
				}
				ch.n++
			case <-ch.close:
				ch.closed()
				return
			}
		}
	}
}

func (ch *UnboundedStructChan) closed() {
	close(ch.in)
	for range ch.in {
		ch.n++
	}
	for ; ch.n > 0; ch.n-- {
		select {
		case ch.out <- struct{}{}:
		default:
		}
	}
	close(ch.out)
	close(ch.close)
}
