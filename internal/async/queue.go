package async

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Queue implements lock-free FIFO freelist based queue.
//
// Reference: https://dl.acm.org/citation.cfm?doid=248052.248106
type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	len  uint64
}

// NewQueue returns a queue for caching values.
func NewQueue() *Queue {
	head := &item{next: nil, v: nil}
	return &Queue{
		tail: unsafe.Pointer(head),
		head: unsafe.Pointer(head),
	}
}

var itemPool = sync.Pool{
	New: func() interface{} { return &item{next: nil, v: nil} },
}

// In puts the given value at the tail of the queue.
func (q *Queue) In(v interface{}) {
	i := itemPool.Get().(*item)
	i.next = nil
	i.v = v

	var last, lastnext *item
	for {
		last = load(&q.tail)
		lastnext = load(&last.next)
		if load(&q.tail) == last {
			if lastnext == nil {
				if cas(&last.next, lastnext, i) {
					cas(&q.tail, last, i)
					atomic.AddUint64(&q.len, 1)
					return
				}
			} else {
				cas(&q.tail, last, lastnext)
			}
		}
	}
}

// Out removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *Queue) Out() interface{} {
	var first, last, firstnext *item
	for {
		first = load(&q.head)
		last = load(&q.tail)
		firstnext = load(&first.next)
		if first == load(&q.head) {
			if first == last {
				if firstnext == nil {
					return nil
				}
				cas(&q.tail, last, firstnext)
			} else {
				v := firstnext.v
				if cas(&q.head, first, firstnext) {
					atomic.AddUint64(&q.len, ^uint64(0))
					itemPool.Put(first)
					return v
				}
			}
		}
	}
}

// Len returns the length of the queue.
func (q *Queue) Len() uint64 { return atomic.LoadUint64(&q.len) }

type item struct {
	next unsafe.Pointer
	v    interface{}
}

func load(p *unsafe.Pointer) *item { return (*item)(atomic.LoadPointer(p)) }
func cas(p *unsafe.Pointer, old, new *item) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
