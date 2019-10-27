package dataapi

import (
	"sync"
	"time"
)

// Clock - an implementation of a simple string DataItem that changes every second
type Clock struct {
	sync.RWMutex
	t    *time.Ticker
	l    map[int]func(item DataItem)
	i    int
	done chan bool
}

func NewClock() Clock {
	c := Clock{
		t:    time.NewTicker(time.Second),
		l:    make(map[int]func(item DataItem)),
		done: make(chan bool),
	}
	go c.run()
	return c
}

func (c Clock) run() {
	for {
		select {
		case <-c.done:
			return
		case <-c.t.C:
			println("tick tock")
			for _, f := range c.l {
				f(c)
			}
		}
	}
}

func (c Clock) stop() {
	c.Lock()
	defer c.Unlock()
	c.t.Stop()
	c.done <- true
}

func (c Clock) String() string {
	return time.Now().Format("03:04:05")
}

func (c Clock) AddListener(f func(data DataItem)) int {
	c.Lock()
	defer c.Unlock()
	c.i++
	c.l[c.i] = f
	return c.i
}

func (c Clock) DeleteListener(i int) {
	c.Lock()
	defer c.Unlock()
	delete(c.l, i)
}
