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

const (
	cursorInterruptTime = 300 * time.Millisecond
	cursorFadeAlpha     = uint8(0x16)
	cursorFadeRatio     = 0.1
)

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
	ri, gi, bi, ai := col.ToNRGBA(theme.PrimaryColor())
	r := uint8(ri >> 8)
	g := uint8(gi >> 8)
	b := uint8(bi >> 8)
	endA := uint8(ai >> 8)
	startA := cursorFadeAlpha
	cursorDim := color.NRGBA{R: r, G: g, B: b, A: cursorFadeAlpha}
	if inverted {
		a.cursor.FillColor = cursorOpaque
		startA, endA = endA, startA
	} else {
		a.cursor.FillColor = cursorDim
	}

	deltaA := endA - startA
	fadeStart := float32(0.5 - cursorFadeRatio)
	fadeStop := float32(0.5 + cursorFadeRatio)

	interrupted := false
	anim := fyne.NewAnimation(time.Second/2, func(f float32) {
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

		alpha := uint8(0)
		if f < fadeStart {
			if _, _, _, al := a.cursor.FillColor.RGBA(); uint8(al>>8) == cursorFadeAlpha {
				return
			}

			a.cursor.FillColor = cursorDim
		} else if f >= fadeStop {
			if _, _, _, al := a.cursor.FillColor.RGBA(); al == 0xffff {
				return
			}

			a.cursor.FillColor = cursorOpaque
		} else {
			fade := (f + cursorFadeRatio - 0.5) * (1 / (cursorFadeRatio * 2))
			alpha = uint8(float32(deltaA) * fade)
			a.cursor.FillColor = color.NRGBA{R: r, G: g, B: b, A: alpha}
		}

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
