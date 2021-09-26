package async

// UnboundedStructChan is a channel with an unbounded buffer for caching
// struct{} objects. This implementation is a specialized version that
// optimizes for struct{} objects than other types.
//
// Delicate dance: One must aware that an unbounded channel may lead
// to OOM when the consuming speed of the buffer is lower than the
// producing speed constantly. However, such a channel may be fairly
// used for event delivering if the consumer of the channel consumes
// the incoming forever. Especially, rendering and even processing tasks.
type UnboundedStructChan struct {
	in, out, close chan struct{}
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
	go func() {
		var n uint64
		for {
			select {
			case _, ok := <-ch.in:
				if !ok {
					close(ch.out)
					return
				}
				n++
			case <-ch.close:
				goto closed
			}
			for n > 0 {
				select {
				case ch.out <- struct{}{}:
					n--
				case _, ok := <-ch.in:
					if ok {
						n++
						break
					}
					for ; n > 0; n-- {
						ch.out <- struct{}{}
					}
					return
				case <-ch.close:
					goto closed
				}
			}
		}

	closed:
		close(ch.in)
		for range ch.in {
			n++
		}
		for ; n > 0; n-- {
			select {
			case ch.out <- struct{}{}:
			default:
			}
		}
		close(ch.out)
		close(ch.close)
	}()
	return ch
}

// In returns a send-only channel that can be used to send values
// to the channel.
func (ch *UnboundedStructChan) In() chan<- struct{} { return ch.in }

// Out returns a receive-only channel that can be used to receive
// values from the channel.
func (ch *UnboundedStructChan) Out() <-chan struct{} { return ch.out }

// Close closes the channel.
func (ch *UnboundedStructChan) Close() {
	ch.close <- struct{}{}
}
