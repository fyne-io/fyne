package widget

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// ===============================================================
// Cursor ticker
// ===============================================================

type cursorTicker interface {
	WaitTick() (reset bool)
	Stop()
	Reset()
	Start(d time.Duration)
	Started() bool
}

type realTicker struct {
	tim      *time.Timer
	rstCh    chan struct{}
	duration time.Duration
}

func (t *realTicker) Start(d time.Duration) {
	if t.Started() {
		return
	}
	t.duration = d
	t.rstCh = make(chan struct{}, 1)
	t.tim = time.NewTimer(t.duration)
}

func (t *realTicker) Reset() {
	if !t.Started() {
		return
	}
	select {
	case t.rstCh <- struct{}{}:
	default:
	}
}

// Stop must be called in the same go routine where WaitTick is used.
func (t *realTicker) Stop() {
	if !t.Started() {
		return
	}
	if !t.tim.Stop() {
		<-t.tim.C
	}
	t.tim = nil // TODO is it safe?
	t.rstCh = nil
}

func (t *realTicker) WaitTick() (reset bool) {
	if !t.Started() {
		reset = true // TODO what to do here?
		return
	}
	select {
	case <-t.tim.C:
		reset = false
		t.tim.Stop()
	case <-t.rstCh:
		reset = true
		if !t.tim.Stop() {
			<-t.tim.C
		}
	}
	t.tim.Reset(t.duration)
	return
}

func (t *realTicker) Started() bool { return t.tim != nil }

// ===============================================================
// Implementation
// ===============================================================

const cursorInterruptTime = 300 * time.Millisecond

type cursorState int

const (
	cursorStateRunning cursorState = iota
	cursorStateInterrupted
	cursorStateStopped
)

type entryCursorAnimation struct {
	mu       *sync.RWMutex
	inverted bool
	state    cursorState
	ticker   cursorTicker
	cursor   *canvas.Rectangle
	anim     *fyne.Animation
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor}
	a.ticker = &realTicker{}
	a.inverted = false
	a.state = cursorStateStopped
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
	if a.anim != nil || a.ticker.Started() || a.state != cursorStateStopped {
		return
	}
	a.anim = a.createAnim(false)
	a.state = cursorStateRunning
	a.ticker.Start(cursorInterruptTime)
	go func() {
		defer func() {
			a.mu.Lock()
			a.state = cursorStateStopped
			a.ticker.Stop()
			a.mu.Unlock()
		}()
		for {
			if reset := a.ticker.WaitTick(); reset {
				continue
			}
			a.mu.RLock()
			interrupted := a.state == cursorStateInterrupted
			cancel := a.anim == nil
			a.mu.RUnlock()
			if cancel {
				return
			}
			if !interrupted {
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

// TemporaryStop temporarily stops the cursor by "cursorInterruptTime".
func (a *entryCursorAnimation) temporaryStop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim == nil || !a.ticker.Started() {
		return
	}
	a.ticker.Reset()
	a.anim.Stop()
	if !a.inverted {
		a.anim = a.createAnim(true)
	}
	if a.state == cursorStateInterrupted {
		return
	}
	a.state = cursorStateInterrupted
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
