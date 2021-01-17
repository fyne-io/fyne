package widget

import (
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

const cursorPauseTimex10ms = 35 // pause time in multiple of 10 ms

type entryCursorAnimation struct {
	mu       *sync.RWMutex
	counter  *safeCounter
	inverted bool
	paused   bool
	stopped  bool
	cursor   *canvas.Rectangle
	anim     *fyne.Animation
	sleepFn  func(time.Duration) // added for tests, do not change it outside of tests!!
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, counter: &safeCounter{}}
	a.inverted = false
	a.stopped = true
	a.sleepFn = time.Sleep
	return a
}

func (a *entryCursorAnimation) createAnim(inverted bool) *fyne.Animation {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := theme.PrimaryColor().RGBA()
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}
	a.inverted = inverted
	start, end := color.Color(cursorDim), color.Color(cursorOpaque)
	if inverted {
		start, end = end, start
	}
	anim := canvas.NewColorRGBAAnimation(start, end, time.Second/2, func(c color.Color) {
		a.cursor.FillColor = c
		a.cursor.Refresh()
	})
	anim.RepeatCount = fyne.AnimationRepeatForever
	anim.AutoReverse = true
	return anim
}

// Start starts cursor animation.
func (a *entryCursorAnimation) Start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	// anim should be nil and current state should be stopped
	if a.anim != nil || !a.stopped {
		return
	}
	a.anim = a.createAnim(false)
	a.stopped = false
	go func() {
		defer func() {
			a.mu.Lock()
			a.stopped = true
			a.mu.Unlock()
		}()
		for {
			a.sleepFn(10 * time.Millisecond)
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
			if a.counter.Value() == cursorPauseTimex10ms {
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

// TemporaryPause pauses temporarily the cursor by `cursorPauseTimex10ms`.
func (a *entryCursorAnimation) TemporaryPause() {
	a.counter.Reset()
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim == nil {
		return
	}
	a.anim.Stop()
	if !a.inverted {
		a.anim = a.createAnim(true)
	}
	if a.paused {
		return
	}
	a.paused = true
	a.cursor.FillColor = theme.PrimaryColor()
	a.cursor.Refresh()
}

// Stop stops cursor animation.
func (a *entryCursorAnimation) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim != nil {
		a.anim.Stop()
		a.paused = false
		a.anim = nil
	}
}
