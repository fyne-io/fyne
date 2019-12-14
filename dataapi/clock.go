package dataapi

// This Clock package will be removed shortly - it belongs in a fyne/x type package

import (
	"time"
)

// Clock - an implementation of a simple string DataItem that changes every second
type Clock struct {
	*BaseDataItem
	t    *time.Ticker
	done chan bool
}

// NewClock returns a new clock dataItem
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
			c.RLock()
			for _, f := range c.callbacks {
				f(c)
			}
			c.RUnlock()
		}
	}
}

func (c *Clock) stop() {
	c.Lock()
	defer c.Unlock()
	c.t.Stop()
	c.done <- true
}

// String returns the value of the clock as a string
func (c *Clock) String() string {
	return time.Now().Format("03:04:05")
}
