package glfw

import (
	"context"
	_ "image/png" // for the icon
	"math"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
)

const (
	doubleClickDelay  = 300 // ms (maximum interval between clicks for double click detection)
	dragMoveThreshold = 2   // how far can we move before it is a drag
	windowIconSize    = 256
)

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title

	w.runOnMainWhenCreated(func() {
		w.view().SetTitle(title)
	})
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

// minSizeOnScreen gets the padded minimum size of a window content in screen pixels
func (w *window) minSizeOnScreen() (int, int) {
	// get minimum size of content inside the window
	return w.screenSize(w.canvas.MinSize())
}

// screenSize computes the actual output size of the given content size in screen pixels
func (w *window) screenSize(canvasSize fyne.Size) (int, int) {
	return scale.ToScreenCoordinate(w.canvas, canvasSize.Width), scale.ToScreenCoordinate(w.canvas, canvasSize.Height)
}

func (w *window) Resize(size fyne.Size) {
	// we cannot perform this until window is prepared as we don't know it's scale!
	bigEnough := size.Max(w.canvas.canvasSize(w.canvas.Content().MinSize()))
	w.runOnMainWhenCreated(func() {
		w.viewLock.Lock()

		width, height := scale.ToScreenCoordinate(w.canvas, bigEnough.Width), scale.ToScreenCoordinate(w.canvas, bigEnough.Height)
		if w.fixedSize || !w.visible { // fixed size ignores future `resized` and if not visible we may not get the event
			w.shouldWidth, w.shouldHeight = width, height
			w.width, w.height = width, height
		}
		w.viewLock.Unlock()
		w.requestedWidth, w.requestedHeight = width, height
		w.view().SetSize(width, height)
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
	if w.view() != nil {
		w.runOnMainWhenCreated(w.fitContent)
	}
}

func (w *window) Padded() bool {
	return w.canvas.padded
}

func (w *window) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)

	w.runOnMainWhenCreated(w.fitContent)
}

func (w *window) Icon() fyne.Resource {
	if w.icon == nil {
		return fyne.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) MainMenu() *fyne.MainMenu {
	return w.mainmenu
}

func (w *window) SetMainMenu(menu *fyne.MainMenu) {
	w.mainmenu = menu
	w.runOnMainWhenCreated(func() {
		w.canvas.buildMenu(w, menu)
	})
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) calculatedScale() float32 {
	return calculateScale(userScale(), fyne.CurrentDevice().SystemScaleForWindow(w), w.detectScale())
}

func (w *window) detectTextureScale() float32 {
	view := w.view()
	winWidth, _ := view.GetSize()
	texWidth, _ := view.GetFramebufferSize()
	return float32(texWidth) / float32(winWidth)
}

func (w *window) Show() {
	go w.doShow()
}

func (w *window) doShow() {
	if w.view() != nil {
		w.doShowAgain()
		return
	}

	run.L.Lock()
	for !run.flag {
		run.Wait()
	}
	run.L.Unlock()

	w.createLock.Do(w.create)
	if w.view() == nil {
		return
	}

	runOnMain(func() {
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
		view := w.view()
		view.SetTitle(w.title)

		if w.centered {
			w.doCenterOnScreen() // lastly center if that was requested
		}
		view.Show()

		// save coordinates
		w.xpos, w.ypos = view.GetPos()

		if w.fullScreen { // this does not work if called before viewport.Show()
			go func() {
				time.Sleep(time.Millisecond * 100)
				w.SetFullScreen(true)
			}()
		}
	})

	// show top canvas element
	if content := w.canvas.Content(); content != nil {
		content.Show()

		runOnDraw(w, func() {
			w.driver.repaintWindow(w)
		})
	}
}

func (w *window) Hide() {
	runOnMain(func() {
		w.viewLock.Lock()
		if w.closing || w.viewport == nil {
			w.viewLock.Unlock()
			return
		}

		w.visible = false
		v := w.viewport
		w.viewLock.Unlock()

		v.Hide()

		// hide top canvas element
		if content := w.canvas.Content(); content != nil {
			content.Hide()
		}
	})
}

func (w *window) Close() {
	if w.isClosing() {
		return
	}

	// trigger callbacks - early so window still exists
	if w.onClosed != nil {
		w.QueueEvent(w.onClosed)
	}

	// set w.closing flag inside draw thread to ensure we can free textures
	runOnDraw(w, func() {
		w.viewLock.Lock()
		w.closing = true
		w.viewLock.Unlock()
		w.viewport.SetShouldClose(true)

		cache.RangeTexturesFor(w.canvas, w.canvas.Painter().Free)
	})

	w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode, _ fyne.Position) {
		if wid, ok := node.Obj().(fyne.Widget); ok {
			cache.DestroyRenderer(wid)
		}
	})
}

func (w *window) ShowAndRun() {
	w.Show()
	w.driver.Run()
}

// Clipboard returns the system clipboard
func (w *window) Clipboard() fyne.Clipboard {
	if w.view() == nil {
		return nil
	}

	if w.clipboard == nil {
		w.clipboard = &clipboard{window: w.viewport}
	}
	return w.clipboard
}

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.Content()
}

func (w *window) SetContent(content fyne.CanvasObject) {
	w.viewLock.RLock()
	visible := w.visible
	w.viewLock.RUnlock()
	// hide old canvas element
	if visible && w.canvas.Content() != nil {
		w.canvas.Content().Hide()
	}

	w.canvas.SetContent(content)

	// show new canvas element
	if content != nil {
		content.Show()
	}
	w.RescaleContext()
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) processClosed() {
	if w.onCloseIntercepted != nil {
		w.QueueEvent(w.onCloseIntercepted)
		return
	}

	go w.Close() // unsure which thread this comes from, so don't block
}

// destroy this window and, if it's the last window quit the app
func (w *window) destroy(d *gLDriver) {
	w.DestroyEventQueue()
	cache.CleanCanvas(w.canvas)

	if w.master {
		d.Quit()
	} else if runtime.GOOS == "darwin" {
		go d.focusPreviousWindow()
	}
}

func (w *window) processMoved(x, y int) {
	if !w.fullScreen { // don't save the move to top left when changing to fullscreen
		// save coordinates
		w.xpos, w.ypos = x, y
	}

	if w.canvas.detectedScale == w.detectScale() {
		return
	}

	w.canvas.detectedScale = w.detectScale()
	go w.canvas.reloadScale()
}

func (w *window) processResized(width, height int) {
	canvasSize := w.computeCanvasSize(width, height)
	if !w.fullScreen {
		w.width = scale.ToScreenCoordinate(w.canvas, canvasSize.Width)
		w.height = scale.ToScreenCoordinate(w.canvas, canvasSize.Height)
	}

	if !w.visible { // don't redraw if hidden
		w.canvas.Resize(canvasSize)
		return
	}

	if w.fixedSize {
		w.canvas.Resize(canvasSize)
		w.fitContent()
		return
	}

	w.platformResize(canvasSize)
}

func (w *window) processFrameSized(width, height int) {
	if width == 0 || height == 0 || runtime.GOOS != "darwin" {
		return
	}

	winWidth, _ := w.view().GetSize()
	newTexScale := float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	w.canvas.RLock()
	texScale := w.canvas.texScale
	w.canvas.RUnlock()
	if texScale != newTexScale {
		w.canvas.Lock()
		w.canvas.texScale = newTexScale
		w.canvas.Unlock()
		w.canvas.Refresh(w.canvas.Content()) // reset graphics to apply texture scale
	}
}

func (w *window) processRefresh() {
	refreshWindow(w)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse fyne.Position, matches func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(mouse, matches, canvas.Overlays().Top(), canvas.menu, canvas.Content())
}

func (w *window) processMouseMoved(xpos float64, ypos float64) {
	w.mouseLock.Lock()
	previousPos := w.mousePos
	w.mousePos = fyne.NewPos(scale.ToFyneCoordinate(w.canvas, int(xpos)), scale.ToFyneCoordinate(w.canvas, int(ypos)))
	mousePos := w.mousePos
	mouseButton := w.mouseButton
	mouseDragPos := w.mouseDragPos
	mouseOver := w.mouseOver
	w.mouseLock.Unlock()

	cursor := desktop.Cursor(desktop.DefaultCursor)

	obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object fyne.CanvasObject) bool {
		if cursorable, ok := object.(desktop.Cursorable); ok {
			cursor = cursorable.Cursor()
		}

		_, hover := object.(desktop.Hoverable)
		return hover
	})

	if w.cursor != cursor {
		// cursor has changed, store new cursor and apply change via glfw
		rawCursor, isCustomCursor := fyneToNativeCursor(cursor)
		w.cursor = cursor

		if rawCursor == nil {
			w.view().SetInputMode(CursorMode, CursorHidden)
		} else {
			w.view().SetInputMode(CursorMode, CursorNormal)
			w.SetCursor(rawCursor)
		}
		w.setCustomCursor(rawCursor, isCustomCursor)
	}

	if w.mouseButton != 0 && w.mouseButton != desktop.MouseButtonSecondary && !w.mouseDragStarted {
		obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, previousPos, func(object fyne.CanvasObject) bool {
			_, ok := object.(fyne.Draggable)
			return ok
		})

		deltaX := mousePos.X - mouseDragPos.X
		deltaY := mousePos.Y - mouseDragPos.Y
		overThreshold := math.Abs(float64(deltaX)) >= dragMoveThreshold || math.Abs(float64(deltaY)) >= dragMoveThreshold

		if wid, ok := obj.(fyne.Draggable); ok && overThreshold {
			w.mouseLock.Lock()
			w.mouseDragged = wid
			w.mouseDraggedOffset = previousPos.Subtract(pos)
			w.mouseDraggedObjStart = obj.Position()
			w.mouseDragStarted = true
			w.mouseLock.Unlock()
		}
	}

	w.mouseLock.RLock()
	isObjDragged := w.objIsDragged(obj)
	isMouseOverDragged := w.objIsDragged(mouseOver)
	w.mouseLock.RUnlock()
	if obj != nil && !isObjDragged {
		ev := &desktop.MouseEvent{Button: mouseButton}
		ev.AbsolutePosition = mousePos
		ev.Position = pos

		if hovered, ok := obj.(desktop.Hoverable); ok {
			if hovered == mouseOver {
				w.QueueEvent(func() { hovered.MouseMoved(ev) })
			} else {
				w.mouseOut()
				w.mouseIn(hovered, ev)
			}
		} else if mouseOver != nil {
			isChild := false
			driver.WalkCompleteObjectTree(mouseOver.(fyne.CanvasObject),
				func(co fyne.CanvasObject, p1, p2 fyne.Position, s fyne.Size) bool {
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
	} else if mouseOver != nil && !isMouseOverDragged {
		w.mouseOut()
	}

	w.mouseLock.RLock()
	mouseButton = w.mouseButton
	mouseDragged := w.mouseDragged
	mouseDraggedObjStart := w.mouseDraggedObjStart
	mouseDraggedOffset := w.mouseDraggedOffset
	mouseDragPos = w.mouseDragPos
	w.mouseLock.RUnlock()
	if mouseDragged != nil && mouseButton != desktop.MouseButtonSecondary {
		if w.mouseButton > 0 {
			draggedObjDelta := mouseDraggedObjStart.Subtract(mouseDragged.(fyne.CanvasObject).Position())
			ev := &fyne.DragEvent{}
			ev.AbsolutePosition = mousePos
			ev.Position = mousePos.Subtract(mouseDraggedOffset).Add(draggedObjDelta)
			ev.Dragged = fyne.NewDelta(mousePos.X-mouseDragPos.X, mousePos.Y-mouseDragPos.Y)
			wd := mouseDragged
			w.QueueEvent(func() { wd.Dragged(ev) })
		}

		w.mouseLock.Lock()
		w.mouseDragStarted = true
		w.mouseDragPos = mousePos
		w.mouseLock.Unlock()
	}
}

func (w *window) objIsDragged(obj interface{}) bool {
	if w.mouseDragged != nil && obj != nil {
		draggedObj, _ := obj.(fyne.Draggable)
		return draggedObj == w.mouseDragged
	}
	return false
}

func (w *window) mouseIn(obj desktop.Hoverable, ev *desktop.MouseEvent) {
	w.QueueEvent(func() {
		if obj != nil {
			obj.MouseIn(ev)
		}
		w.mouseLock.Lock()
		w.mouseOver = obj
		w.mouseLock.Unlock()
	})
}

func (w *window) mouseOut() {
	w.QueueEvent(func() {
		w.mouseLock.RLock()
		mouseOver := w.mouseOver
		w.mouseLock.RUnlock()
		if mouseOver != nil {
			mouseOver.MouseOut()
			w.mouseLock.Lock()
			w.mouseOver = nil
			w.mouseLock.Unlock()
		}
	})
}

func (w *window) processMouseClicked(button desktop.MouseButton, action action, modifiers fyne.KeyModifier) {
	w.mouseLock.RLock()
	w.mouseDragPos = w.mousePos
	mousePos := w.mousePos
	mouseDragStarted := w.mouseDragStarted
	w.mouseLock.RUnlock()
	if mousePos.IsZero() { // window may not be focused (darwin mostly) and so position callbacks not happening
		xpos, ypos := w.view().GetCursorPos()
		w.mouseLock.Lock()
		w.mousePos = fyne.NewPos(scale.ToFyneCoordinate(w.canvas, int(xpos)), scale.ToFyneCoordinate(w.canvas, int(ypos)))
		mousePos = w.mousePos
		w.mouseLock.Unlock()
	}

	co, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case fyne.Tappable, fyne.SecondaryTappable, fyne.DoubleTappable, fyne.Focusable, desktop.Mouseable, desktop.Hoverable:
			return true
		case fyne.Draggable:
			if mouseDragStarted {
				return true
			}
		}

		return false
	})
	ev := &fyne.PointEvent{
		Position:         pos,
		AbsolutePosition: mousePos,
	}

	coMouse := co
	if wid, ok := co.(desktop.Mouseable); ok {
		mev := &desktop.MouseEvent{
			Button:   button,
			Modifier: modifiers,
		}
		mev.Position = ev.Position
		mev.AbsolutePosition = mousePos
		w.mouseClickedHandleMouseable(mev, action, wid)
	}

	if wid, ok := co.(fyne.Focusable); !ok || wid != w.canvas.Focused() {
		ignore := false
		_, _, _ = w.findObjectAtPositionMatching(w.canvas, mousePos, func(object fyne.CanvasObject) bool {
			switch object.(type) {
			case fyne.Focusable:
				ignore = true
				return true
			}

			return false
		})

		if !ignore { // if a parent item under the mouse has focus then ignore this tap unfocus
			w.canvas.Unfocus()
		}
	}

	w.mouseLock.Lock()
	if action == press {
		w.mouseButton |= button
	} else if action == release {
		w.mouseButton &= ^button
	}

	mouseDragged := w.mouseDragged
	mouseDragStarted = w.mouseDragStarted
	mouseOver := w.mouseOver
	shouldMouseOut := w.objIsDragged(mouseOver) && !w.objIsDragged(coMouse)
	mousePressed := w.mousePressed
	w.mouseLock.Unlock()

	if action == release && mouseDragged != nil {
		if mouseDragStarted {
			w.QueueEvent(mouseDragged.DragEnd)
			w.mouseLock.Lock()
			w.mouseDragStarted = false
			w.mouseLock.Unlock()
		}
		if shouldMouseOut {
			w.mouseOut()
		}
		w.mouseLock.Lock()
		w.mouseDragged = nil
		w.mouseLock.Unlock()
	}

	_, tap := co.(fyne.Tappable)
	secondary, altTap := co.(fyne.SecondaryTappable)
	if tap || altTap {
		if action == press {
			w.mouseLock.Lock()
			w.mousePressed = co
			w.mouseLock.Unlock()
		} else if action == release {
			if co == mousePressed {
				if button == desktop.MouseButtonSecondary && altTap {
					w.QueueEvent(func() { secondary.TappedSecondary(ev) })
				}
			}
		}
	}

	// Check for double click/tap on left mouse button
	if action == release && button == desktop.MouseButtonPrimary && !mouseDragStarted {
		w.mouseClickedHandleTapDoubleTap(co, ev)
	}
}

func (w *window) mouseClickedHandleMouseable(mev *desktop.MouseEvent, action action, wid desktop.Mouseable) {
	mousePos := mev.AbsolutePosition
	if action == press {
		w.QueueEvent(func() { wid.MouseDown(mev) })
	} else if action == release {
		w.mouseLock.RLock()
		mouseDragged := w.mouseDragged
		mouseDraggedOffset := w.mouseDraggedOffset
		w.mouseLock.RUnlock()
		if mouseDragged == nil {
			w.QueueEvent(func() { wid.MouseUp(mev) })
		} else {
			if dragged, ok := mouseDragged.(desktop.Mouseable); ok {
				mev.Position = mousePos.Subtract(mouseDraggedOffset)
				w.QueueEvent(func() { dragged.MouseUp(mev) })
			} else {
				w.QueueEvent(func() { wid.MouseUp(mev) })
			}
		}
	}
}

func (w *window) mouseClickedHandleTapDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent) {
	_, doubleTap := co.(fyne.DoubleTappable)
	if doubleTap {
		w.mouseLock.Lock()
		w.mouseClickCount++
		w.mouseLastClick = co
		mouseCancelFunc := w.mouseCancelFunc
		w.mouseLock.Unlock()
		if mouseCancelFunc != nil {
			mouseCancelFunc()
			return
		}
		go w.waitForDoubleTap(co, ev)
	} else {
		w.mouseLock.Lock()
		if wid, ok := co.(fyne.Tappable); ok && co == w.mousePressed {
			w.QueueEvent(func() { wid.Tapped(ev) })
		}
		w.mousePressed = nil
		w.mouseLock.Unlock()
	}
}

func (w *window) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent) {
	var ctx context.Context
	w.mouseLock.Lock()
	ctx, w.mouseCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(time.Millisecond*doubleClickDelay))
	defer w.mouseCancelFunc()
	w.mouseLock.Unlock()

	<-ctx.Done()

	w.mouseLock.Lock()
	defer w.mouseLock.Unlock()

	if w.mouseClickCount == 2 && w.mouseLastClick == co {
		if wid, ok := co.(fyne.DoubleTappable); ok {
			w.QueueEvent(func() { wid.DoubleTapped(ev) })
		}
	} else if co == w.mousePressed {
		if wid, ok := co.(fyne.Tappable); ok {
			w.QueueEvent(func() { wid.Tapped(ev) })
		}
	}

	w.mouseClickCount = 0
	w.mousePressed = nil
	w.mouseCancelFunc = nil
	w.mouseLastClick = nil
}

func (w *window) processMouseScrolled(xoff float64, yoff float64) {
	w.mouseLock.RLock()
	mousePos := w.mousePos
	w.mouseLock.RUnlock()
	co, pos, _ := w.findObjectAtPositionMatching(w.canvas, mousePos, func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Scrollable)
		return ok
	})
	switch wid := co.(type) {
	case fyne.Scrollable:
		if math.Abs(xoff) >= scrollAccelerateCutoff {
			xoff *= scrollAccelerateRate
		}
		if math.Abs(yoff) >= scrollAccelerateCutoff {
			yoff *= scrollAccelerateRate
		}

		ev := &fyne.ScrollEvent{}
		ev.Scrolled = fyne.NewDelta(float32(xoff)*scrollSpeed, float32(yoff)*scrollSpeed)
		ev.Position = pos
		ev.AbsolutePosition = mousePos
		wid.Scrolled(ev)
	}
}

func (w *window) capturesTab(modifier fyne.KeyModifier) bool {
	captures := false

	if ent, ok := w.canvas.Focused().(fyne.Tabbable); ok {
		captures = ent.AcceptsTab()
	}
	if !captures {
		switch modifier {
		case 0:
			w.QueueEvent(w.canvas.FocusNext)
			return false
		case fyne.KeyModifierShift:
			w.QueueEvent(w.canvas.FocusPrevious)
			return false
		}
	}

	return captures
}

func (w *window) processKeyPressed(keyName fyne.KeyName, keyASCII fyne.KeyName, scancode int, action action, keyDesktopModifier fyne.KeyModifier) {
	keyEvent := &fyne.KeyEvent{Name: keyName, Physical: fyne.HardwareKey{ScanCode: scancode}}

	pendingMenuToggle := w.menuTogglePending
	pendingMenuDeactivation := w.menuDeactivationPending
	w.menuTogglePending = desktop.KeyNone
	w.menuDeactivationPending = desktop.KeyNone
	switch action {
	case release:
		if action == release && keyName != "" {
			switch keyName {
			case pendingMenuToggle:
				w.canvas.ToggleMenu()
			case pendingMenuDeactivation:
				if w.canvas.DismissMenu() {
					return
				}
			}
		}

		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.QueueEvent(func() { focused.KeyUp(keyEvent) })
			}
		} else if w.canvas.onKeyUp != nil {
			w.QueueEvent(func() { w.canvas.onKeyUp(keyEvent) })
		}
		return // ignore key up in other core events
	case press:
		switch keyName {
		case desktop.KeyAltLeft, desktop.KeyAltRight:
			// compensate for GLFW modifiers bug https://github.com/glfw/glfw/issues/1630
			if (runtime.GOOS == "linux" && keyDesktopModifier == 0) || (runtime.GOOS != "linux" && keyDesktopModifier == fyne.KeyModifierAlt) {
				w.menuTogglePending = keyName
			}
		case fyne.KeyEscape:
			w.menuDeactivationPending = keyName
		}
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.QueueEvent(func() { focused.KeyDown(keyEvent) })
			}
		} else if w.canvas.onKeyDown != nil {
			w.QueueEvent(func() { w.canvas.onKeyDown(keyEvent) })
		}
	default:
		// key repeat will fall through to TypedKey and TypedShortcut
	}

	modifierOtherThanShift := (keyDesktopModifier & fyne.KeyModifierControl) |
		(keyDesktopModifier & fyne.KeyModifierAlt) |
		(keyDesktopModifier & fyne.KeyModifierSuper)
	if (keyName == fyne.KeyTab && modifierOtherThanShift == 0 && !w.capturesTab(keyDesktopModifier)) ||
		w.triggersShortcut(keyName, keyASCII, keyDesktopModifier) {
		return
	}

	// No shortcut detected, pass down to TypedKey
	focused := w.canvas.Focused()
	if focused != nil {
		w.QueueEvent(func() { focused.TypedKey(keyEvent) })
	} else if w.canvas.onTypedKey != nil {
		w.QueueEvent(func() { w.canvas.onTypedKey(keyEvent) })
	}
}

// charInput defines the character with modifiers callback which is called when a
// Unicode character is input.
//
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or more characters.
func (w *window) processCharInput(char rune) {
	if focused := w.canvas.Focused(); focused != nil {
		w.QueueEvent(func() { focused.TypedRune(char) })
	} else if w.canvas.onTypedRune != nil {
		w.QueueEvent(func() { w.canvas.onTypedRune(char) })
	}
}

func (w *window) processFocused(focus bool) {
	if focus {
		if curWindow == nil {
			fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerEnteredForeground()
		}
		curWindow = w
		w.canvas.FocusGained()
	} else {
		w.canvas.FocusLost()
		w.mouseLock.Lock()
		w.mousePos = fyne.Position{}
		w.mouseLock.Unlock()

		go func() { // check whether another window was focused or not
			time.Sleep(time.Millisecond * 100)
			if curWindow != w {
				return
			}

			curWindow = nil
			fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerExitedForeground()
		}()
	}
}

func (w *window) triggersShortcut(localizedKeyName fyne.KeyName, key fyne.KeyName, modifier fyne.KeyModifier) bool {
	var shortcut fyne.Shortcut
	ctrlMod := fyne.KeyModifierControl
	if runtime.GOOS == "darwin" {
		ctrlMod = fyne.KeyModifierSuper
	}
	// User pressing physical keys Ctrl+V while using a Russian (or any non-ASCII) keyboard layout
	// is reported as a fyne.KeyUnknown key with Control modifier. We should still consider this
	// as a "Paste" shortcut.
	// See https://github.com/fyne-io/fyne/pull/2587 for discussion.
	keyName := localizedKeyName
	resemblesShortcut := (modifier&(fyne.KeyModifierControl|fyne.KeyModifierSuper) != 0)
	if (localizedKeyName == fyne.KeyUnknown) && resemblesShortcut {
		if key != fyne.KeyUnknown {
			keyName = key
		}
	}
	if modifier == ctrlMod {
		switch keyName {
		case fyne.KeyV:
			// detect paste shortcut
			shortcut = &fyne.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyC, fyne.KeyInsert:
			// detect copy shortcut
			shortcut = &fyne.ShortcutCopy{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyX:
			// detect cut shortcut
			shortcut = &fyne.ShortcutCut{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyA:
			// detect selectAll shortcut
			shortcut = &fyne.ShortcutSelectAll{}
		}
	}

	if modifier == fyne.KeyModifierShift {
		switch keyName {
		case fyne.KeyInsert:
			// detect paste shortcut
			shortcut = &fyne.ShortcutPaste{
				Clipboard: w.Clipboard(),
			}
		case fyne.KeyDelete:
			// detect cut shortcut
			shortcut = &fyne.ShortcutCut{
				Clipboard: w.Clipboard(),
			}
		}
	}

	if shortcut == nil && modifier != 0 && !isKeyModifier(keyName) && modifier != fyne.KeyModifierShift {
		shortcut = &desktop.CustomShortcut{
			KeyName:  keyName,
			Modifier: modifier,
		}
	}

	if shortcut != nil {
		if focused, ok := w.canvas.Focused().(fyne.Shortcutable); ok {
			shouldRunShortcut := true
			type selectableText interface {
				fyne.Disableable
				SelectedText() string
			}
			if selectableTextWid, ok := focused.(selectableText); ok && selectableTextWid.Disabled() {
				shouldRunShortcut = shortcut.ShortcutName() == "Copy"
			}
			if shouldRunShortcut {
				w.QueueEvent(func() { focused.TypedShortcut(shortcut) })
			}
			return shouldRunShortcut
		}
		w.QueueEvent(func() { w.canvas.TypedShortcut(shortcut) })
		return true
	}

	return false
}

func (w *window) RunWithContext(f func()) {
	if w.isClosing() {
		return
	}
	w.view().MakeContextCurrent()

	f()

	w.DetachCurrentContext()
}

func (w *window) RescaleContext() {
	runOnMain(w.rescaleOnMain)
}

func (w *window) Context() interface{} {
	return nil
}

func (w *window) runOnMainWhenCreated(fn func()) {
	if w.view() != nil {
		runOnMain(fn)
		return
	}

	w.pending = append(w.pending, fn)
}

func (d *gLDriver) CreateWindow(title string) fyne.Window {
	return d.createWindow(title, true)
}

func (d *gLDriver) createWindow(title string, decorate bool) fyne.Window {
	var ret *window
	if title == "" {
		title = defaultTitle
	}
	runOnMain(func() {
		d.initGLFW()

		ret = &window{title: title, decorate: decorate, driver: d}
		// This queue is destroyed when the window is closed.
		ret.InitEventQueue()
		go ret.RunEventQueue()

		ret.canvas = newCanvas()
		ret.canvas.context = ret
		ret.SetIcon(ret.icon)
		d.addWindow(ret)
	})
	return ret
}

func (w *window) doShowAgain() {
	if w.isClosing() {
		return
	}

	runOnMain(func() {
		// show top canvas element
		if content := w.canvas.Content(); content != nil {
			content.Show()
		}

		view := w.view()
		view.SetPos(w.xpos, w.ypos)
		view.Show()
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
	})
}

func (w *window) isClosing() bool {
	w.viewLock.RLock()
	closing := w.closing || w.viewport == nil
	w.viewLock.RUnlock()
	return closing
}

func (d *gLDriver) CreateSplashWindow() fyne.Window {
	win := d.createWindow("", false)
	win.SetPadded(false)
	win.CenterOnScreen()
	return win
}

func (d *gLDriver) AllWindows() []fyne.Window {
	return d.windows
}

func isKeyModifier(keyName fyne.KeyName) bool {
	return keyName == desktop.KeyShiftLeft || keyName == desktop.KeyShiftRight ||
		keyName == desktop.KeyControlLeft || keyName == desktop.KeyControlRight ||
		keyName == desktop.KeyAltLeft || keyName == desktop.KeyAltRight ||
		keyName == desktop.KeySuperLeft || keyName == desktop.KeySuperRight
}
