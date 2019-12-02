package dataapi

import (
	"time"
)

// Clock - an implementation of a simple string DataItem that changes every second
type Clock struct {
	*BaseDataItem
	t    *time.Ticker
	done chan bool
}

func NewClock() *Clock {
	c := &Clock{
		BaseDataItem: NewBaseDataItem(),
		t:            time.NewTicker(time.Second),
		done:         make(chan bool),
	}
	go c.run()
	return c
}

func (c *Clock) run() {
	for {
		select {
		case <-c.done:
			return
		case <-c.t.C:
			for _, f := range c.callbacks {
				f(c)
			}
		}
	}
}

func (c *Clock) stop() {
	c.Lock()
	defer c.Unlock()
	c.t.Stop()
	c.done <- true
}

func (c *Clock) String() string {
	return time.Now().Format("03:04:05")
}
