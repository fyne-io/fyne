package widget

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type safeCounter struct {
	mu  sync.RWMutex
	val int
}

func (c *safeCounter) Inc(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val += n
}

func (c *safeCounter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.val
}

func (c *safeCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val = 0
}

const cursorStopTimex10ms = 35 // stop time in multiple of 10 ms

type cursorState int

const (
	cursorStateRunning cursorState = iota
	cursorStateTemporarilyStopped
	cursorStateStopped
)

type entryCursorAnimation struct {
	mu       *sync.RWMutex
	counter  *safeCounter
	inverted bool
	state    cursorState
	cursor   *canvas.Rectangle
	anim     *fyne.Animation
	sleepFn  func(time.Duration) // added for tests, do not change it outside of tests!!
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, counter: &safeCounter{}}
	a.inverted = false
	a.state = cursorStateStopped
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
func (a *entryCursorAnimation) start() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim != nil || a.state != cursorStateStopped {
		return
	}
	a.anim = a.createAnim(false)
	a.state = cursorStateRunning
	go func() {
		defer func() {
			a.mu.Lock()
			a.state = cursorStateStopped
			a.mu.Unlock()
		}()
		for {
			a.sleepFn(10 * time.Millisecond)
			a.mu.RLock()
			tempStop := a.state == cursorStateTemporarilyStopped
			cancel := a.anim == nil
			a.mu.RUnlock()
			if cancel {
				return
			}
			if !tempStop {
				continue
			}
			a.counter.Inc(1)
			if a.counter.Value() != cursorStopTimex10ms {
				continue
			}
			a.mu.Lock()
			if a.anim != nil {
				a.anim.Start()
			}
			a.state = cursorStateRunning
			a.mu.Unlock()
		}
	}()
	a.anim.Start()
}

// TemporaryStop temporarily stops the cursor by "cursorStopTimex10ms".
func (a *entryCursorAnimation) temporaryStop() {
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
	if a.state == cursorStateTemporarilyStopped {
		return
	}
	a.state = cursorStateTemporarilyStopped
	a.cursor.FillColor = theme.PrimaryColor()
	a.cursor.Refresh()
}

// Stop stops cursor animation.
func (a *entryCursorAnimation) stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim != nil {
		a.anim.Stop()
		a.anim = nil
	}
}
