package embedded

import (
	"context"
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/test"
)

const (
	tapMoveDecay        = 0.92                   // how much should the scroll continue decay on each frame?
	tapMoveEndThreshold = 2.0                    // at what offset will we stop decaying?
	tapMoveThreshold    = 4.0                    // how far can we move before it is a drag
	tapSecondaryDelay   = 300 * time.Millisecond // how long before secondary tap
	tapDoubleDelay      = 500 * time.Millisecond // max duration between taps for a DoubleTap event
)

type touchCanvas struct {
	test.WindowlessCanvas

	lastTapDown    map[int]time.Time
	lastTapDownPos map[int]fyne.Position
	lastTapDelta   map[int]fyne.Delta

	dragOffset fyne.Position
	dragStart  fyne.Position
	dragging   fyne.Draggable

	touched         map[int]mobile.Touchable
	touchCancelFunc context.CancelFunc
	touchCancelLock sync.Mutex
	touchLastTapped fyne.CanvasObject
	touchTapCount   int
}

func newTouchCanvas() *touchCanvas {
	ret := &touchCanvas{
		WindowlessCanvas: software.NewCanvas(),
		lastTapDown:      make(map[int]time.Time),
		lastTapDownPos:   make(map[int]fyne.Position),
		lastTapDelta:     make(map[int]fyne.Delta),
		touched:          make(map[int]mobile.Touchable),
	}
	return ret
}

func (c *touchCanvas) tapDown(pos fyne.Position, tapID int) {
	c.lastTapDown[tapID] = time.Now()
	c.lastTapDownPos[tapID] = pos
	c.dragging = nil

	co, objPos, layer := driver.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case mobile.Touchable, fyne.Focusable:
			return true
		}

		return false
	}, nil, nil, c.Content())

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchDown(touchEv)
		c.touched[tapID] = wid
	}

	if layer != 1 { // 0 - overlay, 1 - window head / menu, 2 - content
		if wid, ok := co.(fyne.Focusable); !ok || wid != c.Focused() {
			c.Unfocus()
		}
	}
}

func (c *touchCanvas) tapMove(pos fyne.Position, tapID int,
	dragCallback func(fyne.Draggable, *fyne.DragEvent),
) {
	previousPos := c.lastTapDownPos[tapID]
	deltaX := pos.X - previousPos.X
	deltaY := pos.Y - previousPos.Y

	if c.dragging == nil && (math.Abs(float64(deltaX)) < tapMoveThreshold && math.Abs(float64(deltaY)) < tapMoveThreshold) {
		return
	}
	c.lastTapDownPos[tapID] = pos
	offset := fyne.Delta{DX: deltaX, DY: deltaY}
	c.lastTapDelta[tapID] = offset

	co, objPos, _ := driver.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Draggable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		}

		return false
	}, nil, nil, c.Content())

	if c.touched[tapID] != nil {
		if touch, ok := co.(mobile.Touchable); !ok || c.touched[tapID] != touch {
			touchEv := &mobile.TouchEvent{}
			touchEv.Position = objPos
			touchEv.AbsolutePosition = pos
			c.touched[tapID].TouchCancel(touchEv)
			c.touched[tapID] = nil
		}
	}

	if c.dragging == nil {
		if drag, ok := co.(fyne.Draggable); ok {
			c.dragging = drag
			c.dragOffset = previousPos.Subtract(objPos)
			c.dragStart = co.Position()
		} else {
			return
		}
	}

	ev := &fyne.DragEvent{}
	draggedObjDelta := c.dragStart.Subtract(c.dragging.(fyne.CanvasObject).Position())
	ev.Position = pos.Subtract(c.dragOffset).Add(draggedObjDelta)
	ev.Dragged = offset

	dragCallback(c.dragging, ev)
}

func (c *touchCanvas) tapUp(pos fyne.Position, tapID int,
	tapCallback func(fyne.Tappable, *fyne.PointEvent),
	tapAltCallback func(fyne.SecondaryTappable, *fyne.PointEvent),
	doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent),
	dragCallback func(fyne.Draggable, *fyne.DragEvent),
) {
	if c.dragging != nil {
		previousDelta := c.lastTapDelta[tapID]
		ev := &fyne.DragEvent{Dragged: previousDelta}
		draggedObjDelta := c.dragStart.Subtract(c.dragging.(fyne.CanvasObject).Position())
		ev.Position = pos.Subtract(c.dragOffset).Add(draggedObjDelta)
		ev.AbsolutePosition = pos
		dragCallback(c.dragging, ev)

		c.dragging = nil
		return
	}

	duration := time.Since(c.lastTapDown[tapID])

	co, objPos, _ := driver.FindObjectAtPositionMatching(pos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.SecondaryTappable); ok {
			return true
		} else if _, ok := object.(mobile.Touchable); ok {
			return true
		} else if _, ok := object.(fyne.DoubleTappable); ok {
			return true
		}

		return false
	}, nil, nil, c.Content())

	if wid, ok := co.(mobile.Touchable); ok {
		touchEv := &mobile.TouchEvent{}
		touchEv.Position = objPos
		touchEv.AbsolutePosition = pos
		wid.TouchUp(touchEv)
		c.touched[tapID] = nil
	}

	ev := &fyne.PointEvent{
		Position:         objPos,
		AbsolutePosition: pos,
	}

	if duration < tapSecondaryDelay {
		_, doubleTap := co.(fyne.DoubleTappable)
		if doubleTap {
			c.touchCancelLock.Lock()
			c.touchTapCount++
			c.touchLastTapped = co
			cancel := c.touchCancelFunc
			c.touchCancelLock.Unlock()
			if cancel != nil {
				cancel()
				return
			}
			go c.waitForDoubleTap(co, ev, tapCallback, doubleTapCallback)
		} else {
			if wid, ok := co.(fyne.Tappable); ok {
				tapCallback(wid, ev)
			}
		}
	} else {
		if wid, ok := co.(fyne.SecondaryTappable); ok {
			tapAltCallback(wid, ev)
		}
	}
}

func (c *touchCanvas) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent, tapCallback func(fyne.Tappable, *fyne.PointEvent), doubleTapCallback func(fyne.DoubleTappable, *fyne.PointEvent)) {
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(tapDoubleDelay))
	c.touchCancelLock.Lock()
	c.touchCancelFunc = cancel
	c.touchCancelLock.Unlock()
	defer cancel()

	<-ctx.Done()
	fyne.CurrentApp().Driver().DoFromGoroutine(func() {
		c.touchCancelLock.Lock()
		touchCount := c.touchTapCount
		touchLast := c.touchLastTapped
		c.touchCancelLock.Unlock()

		if touchCount == 2 && touchLast == co {
			if wid, ok := co.(fyne.DoubleTappable); ok {
				doubleTapCallback(wid, ev)
			}
		} else {
			if wid, ok := co.(fyne.Tappable); ok {
				tapCallback(wid, ev)
			}
		}

		c.touchCancelLock.Lock()
		c.touchTapCount = 0
		c.touchCancelFunc = nil
		c.touchLastTapped = nil
		c.touchCancelLock.Unlock()
	}, true)
}
