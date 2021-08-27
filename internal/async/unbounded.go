package async

import "fyne.io/fyne/v2"

// TODO: use type parameters when Go 1.18 is released.

type T = fyne.CanvasObject

// UnboundedCanvasObjectChan is a channel with an unbounded buffer for caching
// fyne.CanvasObject objects.
//
// Delicate dance: One must aware that an unbounded channel may lead
// to OOM when the consuming speed of the buffer is lower than the
// producing speed constantly. However, such a channel may be fairly
// used for event delivering if the consumer of the channel consumes
// the incoming forever, especially, rendering tasks.
type UnboundedCanvasObjectChan struct {
	in, out chan T
}

// NewUnboundedCanvasObject returns a unbounded channel with an interface
// buffer that has unlimited capacity.
func NewUnboundedCanvasObject() *UnboundedCanvasObjectChan {
	ch := &UnboundedCanvasObjectChan{
		// A fyne.CanvasObject is a 16-bit pointer, we use 128 that fits
		// a CPU cache line (L2, 256 Bytes) that may reduce cache misses.
		in:  make(chan T, 128),
		out: make(chan T, 128),
	}
	go func() {
		// This is a preallocation of the internal unbounded buffer.
		// The size is randomly picked. Further, there is no memory leak
		// since the queue is garbage collected.
		q := make([]T, 0, 1<<10)
		for {
			e, ok := <-ch.in
			if !ok {
				close(ch.out)
				return
			}
			q = append(q, e)
			for len(q) > 0 {
				select {
				case ch.out <- q[0]:
					q = q[1:]
				case e, ok := <-ch.in:
					if ok {
						q = append(q, e)
						break
					}
					for _, e := range q {
						ch.out <- e
					}
					close(ch.out)
					return
				}
			}
			// If the remaining capacity is too small, we prefer to
			// reallocate the entire buffer.
			if cap(q) < 1<<5 {
				q = make([]T, 0, 1<<10)
			}
		}
	}()
	return ch
}

// In returns a send-only channel that can be used to send values
// to the channel.
func (ch *UnboundedCanvasObjectChan) In() chan<- T { return ch.in }

// Out returns a receive-only channel that can be used to receive
// values from the channel.
func (ch *UnboundedCanvasObjectChan) Out() <-chan T { return ch.out }
