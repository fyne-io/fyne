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

type cursorTicker struct {
	timer     *time.Timer
	resetChan chan struct{}
	duration  time.Duration
	mockWait  func() (reset bool)
}

func (t *cursorTicker) reset() {
	if !t.started() {
		return
	}
	select {
	case t.resetChan <- struct{}{}:
	default:
	}
}

func (t *cursorTicker) start() {
	if t.started() {
		return
	}
	t.resetChan = make(chan struct{}, 1)
	t.timer = time.NewTimer(t.duration)
}

func (t *cursorTicker) started() bool { return t.timer != nil }

// stop must be called in the same go routine where WaitTick is used.
func (t *cursorTicker) stop() {
	if !t.started() {
		return
	}
	if !t.timer.Stop() {
		<-t.timer.C
	}
	t.timer = nil
	t.resetChan = nil
}

func (t *cursorTicker) waitTick() (reset bool) {
	if !t.started() {
		return
	}
	if t.mockWait != nil {
		return t.mockWait()
	}
	select {
	case <-t.timer.C:
		reset = false
		t.timer.Stop()
	case <-t.resetChan:
		reset = true
		if !t.timer.Stop() {
			<-t.timer.C
		}
	}
	t.timer.Reset(t.duration)
	return
}

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
	ticker   *cursorTicker
	cursor   *canvas.Rectangle
	anim     *fyne.Animation
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor}
	a.ticker = &cursorTicker{duration: cursorInterruptTime}
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
	if a.anim != nil || a.ticker.started() || a.state != cursorStateStopped {
		return
	}
	a.anim = a.createAnim(false)
	a.state = cursorStateRunning
	a.ticker.start()
	go func() {
		for {
			if reset := a.ticker.waitTick(); reset {
				continue
			}
			a.mu.RLock()
			interrupted := a.state == cursorStateInterrupted
			cancel := a.anim == nil
			a.mu.RUnlock()
			if cancel {
				break
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
		a.mu.Lock()
		a.state = cursorStateStopped
		a.ticker.stop()
		a.mu.Unlock()
	}()
	a.anim.Start()
}

// temporaryStop temporarily stops the cursor by "cursorInterruptTime".
func (a *entryCursorAnimation) temporaryStop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim == nil || !a.ticker.started() {
		return
	}
	a.ticker.reset()
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

// stop stops cursor animation.
func (a *entryCursorAnimation) stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.anim != nil {
		a.anim.Stop()
		a.anim = nil
	}
}
