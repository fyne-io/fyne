package widget

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type safeCounter struct {
	mu  sync.RWMutex
	cnt int
}

func (c *safeCounter) Inc(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cnt += n
}

func (c *safeCounter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cnt
}

func (c *safeCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cnt = 0
}

type entryCursorAnimation struct {
	mu *sync.RWMutex
	// timeout *time.Ticker
	counter *safeCounter
	paused  bool
	stopped bool
	cursor  *canvas.Rectangle
	anim    *fyne.Animation
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, counter: &safeCounter{}}
	a.stopped = true
	return a
}

func (a *entryCursorAnimation) createAnim() {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := theme.PrimaryColor().RGBA()
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}
	anim := canvas.NewColorRGBAAnimation(cursorDim, cursorOpaque, time.Second/2, func(c color.Color) {
		a.cursor.FillColor = c
		a.cursor.Refresh()
	})
	anim.RepeatCount = fyne.AnimationRepeatForever
	anim.AutoReverse = true
	a.anim = anim
}

// Start must only be called once.
func (a *entryCursorAnimation) Start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	// anim should be nil and current state should be stopped
	if a.anim != nil || !a.stopped {
		return
	}
	fmt.Printf("anim: %p\n", a)
	a.createAnim()
	a.stopped = false
	// a.timeout = time.NewTicker(10 * time.Millisecond)
	go func() {
		fmt.Println("start goroutine")
		defer func() {
			fmt.Println("===> stop goroutine")
			a.mu.Lock()
			a.stopped = true
			a.mu.Unlock()
		}()
		for {
			// <-a.timeout.C
			time.Sleep(10 * time.Millisecond)
			// >>>> lock
			a.mu.RLock()
			paused := a.paused
			cancel := a.anim == nil
			a.mu.RUnlock()
			// <<<< unlock
			if cancel {
				return
			}
			if !paused {
				continue
			}
			a.counter.Inc(1)
			if a.counter.Value() == 50 {
				// >>>> lock
				a.mu.Lock()
				if a.anim != nil {
					a.anim.Start()
				}
				a.paused = false
				a.mu.Unlock()
				// <<<< unlock
			}
		}
	}()
	a.anim.Start()
}

func (a *entryCursorAnimation) TemporaryPause() {
	a.counter.Reset()
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim == nil {
		return
	}
	a.anim.Stop()
	if a.paused {
		return
	}
	a.paused = true
	a.cursor.FillColor = theme.PrimaryColor()
	a.cursor.Refresh()
}

// Stop must only be called once when cursor is destroyed.
func (a *entryCursorAnimation) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	// if a.timeout != nil {
	// 	a.timeout.Stop()
	// }
	if a.anim != nil {
		a.anim.Stop()
		a.anim = nil
	}
}
