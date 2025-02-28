package common

// RingBuffer is a growable ring buffer supporting
// enqueue and dequeue operations. It is not thread safe.
type RingBuffer[T any] struct {
	buf  []T
	head int
	len  int
}

// NewRingBuffer initializes and returns a new RingBuffer.
func NewRingBuffer[T any](initialCap int) RingBuffer[T] {
	return RingBuffer[T]{
		buf: make([]T, initialCap),
	}
}

// Len returns the number of elements in the buffer.
func (r *RingBuffer[T]) Len() int {
	return r.len
}

// Push adds the value to the end of the buffer.
func (r *RingBuffer[T]) Push(value T) {
	r.checkGrow()

	pos := (r.head + r.len) % len(r.buf)
	r.buf[pos] = value
	r.len++
}

// Pull removes the first item from the buffer, if any.
func (r *RingBuffer[T]) Pull() (value T, ok bool) {
	if r.len == 0 {
		return value, false
	}
	return r.pullOne(), true
}

// PullN removes up to len(buf) items from the queue,
// copying them into the supplied buffer and returning
// the number of elements copied.
func (r *RingBuffer[T]) PullN(buf []T) int {
	l := len(buf)
	if r.len < l {
		l = r.len
	}
	if l == 0 {
		return 0
	}
	for i := 0; i < l; i++ {
		buf[i] = r.pullOne()
	}
	return l
}

func (r *RingBuffer[T]) pullOne() (value T) {
	emptyT := value
	value = r.buf[r.head]
	r.buf[r.head] = emptyT

	if r.len == 1 {
		r.head = 0
	} else {
		r.head = (r.head + 1) % len(r.buf)
	}
	r.len--

	return value
}

func (r *RingBuffer[T]) checkGrow() {
	if l := len(r.buf); r.len == l {
		newBuf := make([]T, l*2)
		for i := 0; i < r.len; i++ {
			newBuf[i] = r.buf[(r.head+i)%l]
		}
		r.head = 0
		r.buf = newBuf
	}
}
