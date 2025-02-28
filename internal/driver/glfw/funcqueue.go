package glfw

import (
	"sync"

	"fyne.io/fyne/v2/internal/driver/common"
)

type funcData struct {
	f    func()
	done chan struct{} // Zero allocation signalling channel
}

// glfwFuncQueue is a queue to hold funcs posted to
// execute on the main thread
type glfwFuncQueue struct {
	mutex sync.Mutex

	buffer common.RingBuffer[funcData]
}

func newGlfwFuncQueue() *glfwFuncQueue {
	return &glfwFuncQueue{
		buffer: common.NewRingBuffer[funcData](1024),
	}
}

// Push adds a new func to the queue. It wakes up
// the driver's main thread if necessary to begin
// processing queued functions.
func (g *glfwFuncQueue) Push(f func(), wait bool) {
	g.mutex.Lock()
	data := funcData{f: f}
	if wait {
		done := common.DonePool.Get()
		defer common.DonePool.Put(done)

		data.done = done
	}
	g.buffer.Push(data)
	l := g.buffer.Len()
	g.mutex.Unlock()

	if l == 1 {
		// only need to wake up driver if there
		// were previously no funcs queued.
		// otherwise, it has already been woken.
		wakeUpDriver()
	}
	if wait {
		<-data.done
	}
}

// PullN pulls up to N funcs from the queue.
func (g *glfwFuncQueue) PullN(buf []funcData) int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.buffer.PullN(buf)
}
