//go:build windows && directx

// Mouse event processing - hover, drag, tap and double tap.

package directx

import (
	"context"
	"math"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/scale"
)

func (w *window) findObjectAtPositionMatching(mouse fyne.Position,
	matches func(fyne.CanvasObject) bool,
) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(mouse, matches, w.canvas.Overlays().Top(), nil, w.canvas.Content())
}

func (w *window) objIsDragged(obj any) bool {
	if w.mouseDragged != nil && obj != nil {
		draggedObj, _ := obj.(fyne.Draggable)
		return draggedObj == w.mouseDragged
	}
	return false
}

func (w *window) mouseIn(obj desktop.Hoverable, ev *desktop.MouseEvent) {
	if obj != nil {
		obj.MouseIn(ev)
	}
	w.mouseOver = obj
}

func (w *window) mouseOut() {
	if w.mouseOver != nil {
		w.mouseOver.MouseOut()
		w.mouseOver = nil
	}
}

func (w *window) processMouseMoved(xpos, ypos float64) {
	previousPos := w.mousePos
	w.mousePos = fyne.NewPos(scale.ToFyneCoordinate(w.canvas, int(xpos)),
		scale.ToFyneCoordinate(w.canvas, int(ypos)))
	mousePos := w.mousePos
	mouseButton := w.mouseButton
	mouseDragPos := w.mouseDragPos
	mouseOver := w.mouseOver

	cursor := desktop.Cursor(desktop.DefaultCursor)
	obj, pos, _ := w.findObjectAtPositionMatching(mousePos, func(object fyne.CanvasObject) bool {
		if cursorable, ok := object.(desktop.Cursorable); ok {
			cursor = cursorable.Cursor()
		}
		_, hover := object.(desktop.Hoverable)
		return hover
	})

	if w.cursor != cursor {
		w.cursor = cursor
		w.applyCursor(cursor)
	}

	if w.mouseButton != 0 && w.mouseButton != desktop.MouseButtonSecondary && !w.mouseDragStarted {
		obj, pos, _ := w.findObjectAtPositionMatching(previousPos, func(object fyne.CanvasObject) bool {
			_, ok := object.(fyne.Draggable)
			return ok
		})

		deltaX := mousePos.X - mouseDragPos.X
		deltaY := mousePos.Y - mouseDragPos.Y
		overThreshold := math.Abs(float64(deltaX)) >= dragMoveThreshold ||
			math.Abs(float64(deltaY)) >= dragMoveThreshold

		if wid, ok := obj.(fyne.Draggable); ok && overThreshold {
			w.mouseDragged = wid
			w.mouseDraggedOffset = previousPos.Subtract(pos)
			w.mouseDraggedObjStart = obj.Position()
			w.mouseDragStarted = true
		}
	}

	if obj != nil && !w.objIsDragged(obj) {
		ev := &desktop.MouseEvent{Button: mouseButton}
		ev.AbsolutePosition = mousePos
		ev.Position = pos

		if hovered, ok := obj.(desktop.Hoverable); ok {
			if hovered == mouseOver {
				hovered.MouseMoved(ev)
			} else {
				w.mouseOut()
				w.mouseIn(hovered, ev)
			}
		} else if mouseOver != nil {
			isChild := false
			driver.WalkCompleteObjectTree(mouseOver.(fyne.CanvasObject),
				func(co fyne.CanvasObject, _, _ fyne.Position, _ fyne.Size) bool {
					if co == obj {
						isChild = true
						return true
					}
					return false
				}, nil)
			if !isChild {
				w.mouseOut()
			}
		}
	} else if mouseOver != nil && !w.objIsDragged(mouseOver) {
		w.mouseOut()
	}

	mouseDragged := w.mouseDragged
	mouseDragPos = w.mouseDragPos
	if mouseDragged != nil && w.mouseButton != desktop.MouseButtonSecondary {
		if w.mouseButton > 0 {
			draggedObjDelta := w.mouseDraggedObjStart.Subtract(mouseDragged.(fyne.CanvasObject).Position())
			ev := &fyne.DragEvent{}
			ev.AbsolutePosition = mousePos
			ev.Position = mousePos.Subtract(w.mouseDraggedOffset).Add(draggedObjDelta)
			ev.Dragged = fyne.NewDelta(mousePos.X-mouseDragPos.X, mousePos.Y-mouseDragPos.Y)
			mouseDragged.Dragged(ev)
		}
		w.mouseDragStarted = true
		w.mouseDragPos = mousePos
	}
}

func (w *window) mouseClicked(button desktop.MouseButton, act action) {
	if act == press {
		procSetCapture.Call(uintptr(w.hwnd))
	} else {
		procReleaseCapture.Call()
	}
	w.processMouseClicked(button, act, currentModifiers())
}

func (w *window) processMouseClicked(button desktop.MouseButton, act action, modifiers fyne.KeyModifier) {
	w.mouseDragPos = w.mousePos
	mousePos := w.mousePos
	mouseDragStarted := w.mouseDragStarted

	co, pos, _ := w.findObjectAtPositionMatching(mousePos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case fyne.Tappable, fyne.SecondaryTappable, fyne.DoubleTappable, fyne.Focusable, desktop.Mouseable:
			return true
		case fyne.Draggable:
			if mouseDragStarted {
				return true
			}
		}
		return false
	})
	ev := &fyne.PointEvent{Position: pos, AbsolutePosition: mousePos}

	coMouse := co
	if wid, ok := co.(desktop.Mouseable); ok {
		mev := &desktop.MouseEvent{Button: button, Modifier: modifiers}
		mev.Position = ev.Position
		mev.AbsolutePosition = mousePos
		w.mouseClickedHandleMouseable(mev, act, wid)
	}

	focused := w.canvas.Focused()
	if wid, ok := co.(fyne.Focusable); !ok || wid != focused {
		ignore := false
		if focusedObj, ok := focused.(fyne.CanvasObject); ok {
			found, _, _ := w.findObjectAtPositionMatching(mousePos, func(object fyne.CanvasObject) bool {
				return object == focusedObj
			})
			ignore = found != nil
		}
		if !ignore {
			w.canvas.Unfocus()
		}
	}

	switch act {
	case press:
		w.mouseButton |= button
	case release:
		w.mouseButton &= ^button
	}

	mouseDragged := w.mouseDragged
	mouseDragStarted = w.mouseDragStarted
	mouseOver := w.mouseOver
	shouldMouseOut := w.objIsDragged(mouseOver) && !w.objIsDragged(coMouse)
	mousePressed := w.mousePressed

	if act == release && mouseDragged != nil {
		if mouseDragStarted {
			mouseDragged.DragEnd()
			w.mouseDragStarted = false
		}
		if shouldMouseOut {
			w.mouseOut()
		}
		w.mouseDragged = nil
	}

	_, tap := co.(fyne.Tappable)
	secondary, altTap := co.(fyne.SecondaryTappable)
	if tap || altTap {
		switch act {
		case press:
			w.mousePressed = co
		case release:
			if co == mousePressed && button == desktop.MouseButtonSecondary && altTap {
				prevOverlay := w.canvas.Overlays().Top()
				secondary.TappedSecondary(ev)

				// If the secondary tap dismissed an overlay rather than opening one,
				// forward the event to whatever is now underneath.
				if prevOverlay != nil && !slices.Contains(w.canvas.Overlays().List(), prevOverlay) {
					co2, pos2, _ := w.findObjectAtPositionMatching(mousePos, func(object fyne.CanvasObject) bool {
						_, ok := object.(fyne.SecondaryTappable)
						return ok
					})
					if sec2, ok := co2.(fyne.SecondaryTappable); ok {
						sec2.TappedSecondary(&fyne.PointEvent{Position: pos2, AbsolutePosition: mousePos})
					}
				}
			}
		}
	}

	if act == release && button == desktop.MouseButtonPrimary && !mouseDragStarted {
		w.mouseClickedHandleTapDoubleTap(co, ev)
	}
}

func (w *window) mouseClickedHandleMouseable(mev *desktop.MouseEvent, act action, wid desktop.Mouseable) {
	switch act {
	case press:
		wid.MouseDown(mev)
	case release:
		mouseDragged := w.mouseDragged
		mouseDraggedOffset := w.mouseDraggedOffset
		if mouseDragged == nil {
			wid.MouseUp(mev)
		} else if dragged, ok := mouseDragged.(desktop.Mouseable); ok {
			mev.Position = mev.AbsolutePosition.Subtract(mouseDraggedOffset)
			dragged.MouseUp(mev)
		} else {
			wid.MouseUp(mev)
		}
	}
}

func (w *window) mouseClickedHandleTapDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent) {
	if _, doubleTap := co.(fyne.DoubleTappable); doubleTap {
		w.mouseClickCount++
		w.mouseLastClick = co

		if w.mouseCancelFunc != nil {
			w.mouseCancelFunc()
			return
		}
		go w.waitForDoubleTap(co, ev)
		return
	}

	if wid, ok := co.(fyne.Tappable); ok && co == w.mousePressed {
		wid.Tapped(ev)
	}
	w.mousePressed = nil
}

func (w *window) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent) {
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(w.driver.DoubleTapDelay()))
	defer runOnMain(cancel)
	runOnMain(func() { w.mouseCancelFunc = cancel })

	<-ctx.Done()

	runOnMain(func() { w.waitForDoubleTapEnded(co, ev) })
}

func (w *window) waitForDoubleTapEnded(co fyne.CanvasObject, ev *fyne.PointEvent) {
	if w.mouseClickCount == 2 && w.mouseLastClick == co {
		if wid, ok := co.(fyne.DoubleTappable); ok {
			wid.DoubleTapped(ev)
		}
	} else if co == w.mousePressed {
		if wid, ok := co.(fyne.Tappable); ok {
			wid.Tapped(ev)
		}
	}

	w.mouseClickCount = 0
	w.mousePressed = nil
	w.mouseCancelFunc = nil
	w.mouseLastClick = nil
}

func (w *window) processMouseScrolled(xoff, yoff float64) {
	co, pos, _ := w.findObjectAtPositionMatching(w.mousePos, func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Scrollable)
		return ok
	})
	wid, ok := co.(fyne.Scrollable)
	if !ok {
		return
	}

	// A high-resolution wheel or trackpad reports large offsets; accelerate those
	// so a hard flick covers real distance instead of creeping.
	if math.Abs(xoff) >= scrollAccelerateCutoff {
		xoff *= scrollAccelerateRate
	}
	if math.Abs(yoff) >= scrollAccelerateCutoff {
		yoff *= scrollAccelerateRate
	}

	ev := &fyne.ScrollEvent{}
	ev.Scrolled = fyne.NewDelta(float32(xoff)*scrollSpeed, float32(yoff)*scrollSpeed)
	ev.Position = pos
	ev.AbsolutePosition = w.mousePos
	wid.Scrolled(ev)
}
