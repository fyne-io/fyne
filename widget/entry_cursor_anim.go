package widget

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/theme"
)

const cursorInterruptTime = 300 * time.Millisecond

type entryCursorAnimation struct {
	mu                *sync.RWMutex
	cursor            *canvas.Rectangle
	anim              *fyne.Animation
	lastInterruptTime time.Time

	timeNow func() time.Time // useful for testing
}

func newEntryCursorAnimation(cursor *canvas.Rectangle) *entryCursorAnimation {
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, timeNow: time.Now}
	return a
}

// creates fyne animation
func (a *entryCursorAnimation) createAnim(inverted bool) *fyne.Animation {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := col.ToNRGBA(theme.PrimaryColor())
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}
	start, end := color.Color(cursorDim), cursorOpaque
	if inverted {
		start, end = cursorOpaque, color.Color(cursorDim)
	}
	interrupted := false
	anim := canvas.NewColorRGBAAnimation(start, end, time.Second/2, func(c color.Color) {
		a.mu.RLock()
		shouldInterrupt := a.timeNow().Sub(a.lastInterruptTime) <= cursorInterruptTime
		a.mu.RUnlock()
		if shouldInterrupt {
			if !interrupted {
				a.cursor.FillColor = cursorOpaque
				a.cursor.Refresh()
				interrupted = true
			}
			return
		}
		if interrupted {
			a.mu.Lock()
			a.anim.Stop()
			if !inverted {
				a.anim = a.createAnim(true)
			}
			interrupted = false
			a.mu.Unlock()
			go func() {
				a.mu.RLock()
				canStart := a.anim != nil
				a.mu.RUnlock()
				if canStart {
					a.anim.Start()
				}
			}()
			return
		}
		a.cursor.FillColor = c
		a.cursor.Refresh()
	})

	anim.RepeatCount = fyne.AnimationRepeatForever
	anim.AutoReverse = true
	return anim
}

// starts cursor animation.
func (a *entryCursorAnimation) start() {
	a.mu.Lock()
	isStopped := a.anim == nil
	if isStopped {
		a.anim = a.createAnim(false)
	}
	a.mu.Unlock()
	if isStopped {
		a.anim.Start()
	}
}

// temporarily stops the animation by "cursorInterruptTime".
func (a *entryCursorAnimation) interrupt() {
	a.mu.Lock()
	a.lastInterruptTime = a.timeNow()
	a.mu.Unlock()
}

// stops cursor animation.
func (a *entryCursorAnimation) stop() {
	a.mu.Lock()
	if a.anim != nil {
		a.anim.Stop()
		a.anim = nil
	}
	a.mu.Unlock()
}
