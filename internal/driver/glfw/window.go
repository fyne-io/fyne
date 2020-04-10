package glfw

import (
	_ "image/png" // for the icon
	"runtime"
	"time"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"
)

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title

	w.runOnMainWhenCreated(func() {
		w.viewport.SetTitle(title)
	})
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) CenterOnScreen() {
	w.centered = true

	if w.visible {
		w.centerOnScreen()
	}
}

// minSizeOnScreen gets the padded minimum size of a window content in screen pixels
func (w *window) minSizeOnScreen() (int, int) {
	// get minimum size of content inside the window
	return w.screenSize(w.canvas.MinSize())
}

// screenSize computes the actual output size of the given content size in screen pixels
func (w *window) screenSize(canvasSize fyne.Size) (int, int) {
	return internal.ScaleInt(w.canvas, canvasSize.Width), internal.ScaleInt(w.canvas, canvasSize.Height)
}

func (w *window) Resize(size fyne.Size) {
	w.canvas.Resize(size)
	scaleSize := internal.ScaleSize(w.canvas, size)

	w.runOnMainWhenCreated(func() {
		w.viewport.SetSize(scaleSize.Width, scaleSize.Height)
		w.fitContent()
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
	w.runOnMainWhenCreated(w.fitContent)
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

func (w *window) calculatedScale() float32 {
	return calculateScale(userScale(), fyne.CurrentDevice().SystemScaleForWindow(w), w.detectScale())
}

func (w *window) Show() {
	go w.doShow()
}

func (w *window) doShow() {
	for !running() {
		time.Sleep(time.Millisecond * 10)
	}
	w.createLock.Do(w.create)

	runOnMain(func() {
		w.visible = true
		w.viewport.SetTitle(w.title)
		w.viewport.Show()

		// save coordinates
		w.xpos, w.ypos = w.viewport.GetPos()
	})

	if w.fullScreen { // this does not work if called before viewport.Show()...
		w.SetFullScreen(true)
	}

	// show top canvas element
	if w.canvas.content != nil {
		w.canvas.content.Show()
	}
}

func (w *window) Hide() {
	if w.viewport == nil {
		return
	}

	runOnMain(func() {
		w.viewport.Hide()
		w.visible = false

		// hide top canvas element
		if w.canvas.Content() != nil {
			w.canvas.Content().Hide()
		}
	})
}

func (w *window) Close() {
	if w.viewport == nil {
		return
	}
	w.closed(w.viewport)
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Driver().Run()
}

//Clipboard returns the system clipboard
func (w *window) Clipboard() fyne.Clipboard {
	if w.viewport == nil {
		return nil
	}

	if w.clipboard == nil {
		w.clipboard = &clipboard{window: w.viewport}
	}
	return w.clipboard
}

func (w *window) Content() fyne.CanvasObject {
	return w.canvas.content
}

func (w *window) resize(canvasSize fyne.Size) {
	if !w.fullScreen && !w.fixedSize {
		w.size = internal.ScaleSize(w.canvas, canvasSize)
	}

	w.canvas.Resize(canvasSize)
}

func (w *window) SetContent(content fyne.CanvasObject) {
	// hide old canvas element
	if w.visible && w.canvas.Content() != nil {
		w.canvas.Content().Hide()
	}

	log.Println("SetContent")
	w.canvas.SetContent(content)
	w.RescaleContext()
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

// destroy this window and, if it's the last window quit the app
func (w *window) destroy(d *gLDriver) {
	w.eventLock.RLock()
	queue := w.eventQueue
	w.eventLock.RUnlock()

	// finish serial event queue and nil it so we don't panic if window.closed() is called twice.
	if queue != nil {
		w.waitForEvents()

		w.eventLock.Lock()
		close(w.eventQueue)
		w.eventQueue = nil
		w.eventLock.Unlock()
	}

	if w.master || len(d.windowList()) == 0 {
		d.Quit()
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
	w.queueEvent(func() {
		if obj != nil {
			obj.MouseIn(ev)
		}
		w.mouseOver = obj
	})
}

func (w *window) mouseOut() {
	w.queueEvent(func() {
		if w.mouseOver != nil {
			w.mouseOver.MouseOut()
			w.mouseOver = nil
		}
	})
}

func (w *window) processMoved(x, y int) {
	if w.fullScreen { // don't save the move to top left when changint to fullscreen
		return
	}

	// save coordinates
	w.xpos, w.ypos = x, y

	if w.canvas.detectedScale == w.detectScale() {
		return
	}

	w.canvas.detectedScale = w.detectScale()
	go w.canvas.SetScale(fyne.SettingsScaleAuto) // scale is ignored
}

func (w *window) processResized(width, height int) {
	w.resize(fyne.NewSize(internal.UnscaleInt(w.canvas, width), internal.UnscaleInt(w.canvas, height)))
}

func (w *window) processFrameSized(width, height int) {
	if width == 0 || height == 0 {
		return
	}

	winWidth, _ := w.viewport.GetSize()
	texScale := float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	w.canvas.texScale = texScale
}

func (w *window) processRefresh() {
	refreshWindow(w)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse fyne.Position, matches func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(mouse, matches, canvas.Overlays().Top(), canvas.menu, canvas.content)
}

func (w *window) processClosed() {
	w.canvas.walkTrees(nil, func(node *renderCacheNode) {
		switch co := node.obj.(type) {
		case fyne.Widget:
			cache.DestroyRenderer(co)
		}
	})

	// trigger callbacks
	if w.onClosed != nil {
		w.queueEvent(w.onClosed)
	}
}

func (w *window) processMouseMoved(xpos, ypos float64) {
	w.mousePos = fyne.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))

	cursor := cursorMap[desktop.DefaultCursor]
	obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		if cursorable, ok := object.(desktop.Cursorable); ok {
			fyneCursor := cursorable.Cursor()
			cursor = fyneToNativeCursor(fyneCursor)
		}

		_, hover := object.(desktop.Hoverable)
		return hover
	})

	w.cursor = cursor
	w.SetCursor(cursor)
	if obj != nil && !w.objIsDragged(obj) {
		ev := new(desktop.MouseEvent)
		ev.Position = pos
		ev.Button = w.mouseButton

		if hovered, ok := obj.(desktop.Hoverable); ok {
			if hovered == w.mouseOver {
				w.queueEvent(func() { hovered.MouseMoved(ev) })
			} else {
				w.mouseOut()
				w.mouseIn(hovered, ev)
			}
		}
	} else if w.mouseOver != nil && !w.objIsDragged(w.mouseOver) {
		w.mouseOut()
	}

	if w.mouseDragged != nil {
		if w.mouseButton > 0 {
			draggedObjPos := w.mouseDragged.(fyne.CanvasObject).Position()
			ev := new(fyne.DragEvent)
			ev.Position = w.mousePos.Subtract(w.mouseDraggedOffset).Subtract(draggedObjPos)
			ev.DraggedX = w.mousePos.X - w.mouseDragPos.X
			ev.DraggedY = w.mousePos.Y - w.mouseDragPos.Y
			wd := w.mouseDragged
			w.queueEvent(func() { wd.Dragged(ev) })

			w.mouseDragStarted = true
			w.mouseDragPos = w.mousePos
		}
	}	
}

func (w *window) processMouseScrolled(xoff, yoff float64) {
	co, _, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Scrollable)
		return ok
	})
	switch wid := co.(type) {
	case fyne.Scrollable:
		ev := &fyne.ScrollEvent{}
		ev.DeltaX = int(xoff * scrollSpeed)
		ev.DeltaY = int(yoff * scrollSpeed)
		wid.Scrolled(ev)
	}
}

func (w *window) processMouseClicked(button desktop.MouseButton, action desktop.Action, modifiers desktop.Modifier) {
	co, pos, layer := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case fyne.Tappable, fyne.SecondaryTappable, fyne.Focusable, fyne.Draggable, desktop.Mouseable, desktop.Hoverable:
			return true
		}

		return false
	})
	ev := new(fyne.PointEvent)
	ev.Position = pos
	ev.AbsolutePosition = w.mousePos

	coMouse := co
	// Switch the mouse target to the dragging object if one is set
	if w.mouseDragged != nil && !w.objIsDragged(co) {
		co, _ = w.mouseDragged.(fyne.CanvasObject)
		ev.Position = w.mousePos.Subtract(w.mouseDraggedOffset).Subtract(co.Position())
	}

	if wid, ok := co.(desktop.Mouseable); ok {
		mev := new(desktop.MouseEvent)
		mev.Position = ev.Position
		mev.AbsolutePosition = w.mousePos
		mev.Button = button
		mev.Modifier = modifiers
		if action == desktop.Press {
			w.queueEvent(func() { wid.MouseDown(mev) })
		} else if action == desktop.Release {
			w.queueEvent(func() { wid.MouseUp(mev) })
		}
	}

	needsfocus := false
	if layer != 1 { // 0 - overlay, 1 - menu, 2 - content
		needsfocus = true

		if wid := w.canvas.Focused(); wid != nil {
			if wid.(fyne.CanvasObject) != co {
				w.canvas.Unfocus()
			} else {
				needsfocus = false
			}
		}
	}

	if action == desktop.Press {
		w.mouseButton = button
	} else if action == desktop.Release {
		w.mouseButton = 0
	}

	// we cannot switch here as objects may respond to multiple cases
	if wid, ok := co.(fyne.Focusable); ok && needsfocus {
		if dis, ok := wid.(fyne.Disableable); !ok || !dis.Disabled() {
			w.canvas.Focus(wid)
		}
	}

	// Check for double click/tap
	doubleTapped := false
	if action == desktop.Release && button == desktop.LeftMouseButton {
		now := time.Now()
		// we can safely subtract the first "zero" time as it'll be much larger than doubleClickDelay
		if now.Sub(w.mouseClickTime).Nanoseconds()/1e6 <= doubleClickDelay {
			if wid, ok := co.(fyne.DoubleTappable); ok {
				doubleTapped = true
				w.queueEvent(func() { wid.DoubleTapped(ev) })
			}
		}
		w.mouseClickTime = now
		w.mouseLastClick = co
	}

	_, tap := co.(fyne.Tappable)
	_, altTap := co.(fyne.SecondaryTappable)
	// Prevent Tapped from triggering if DoubleTapped has been sent
	if (tap || altTap) && doubleTapped == false {
		if action == desktop.Press {
			w.mousePressed = co
		} else if action == desktop.Release {
			if co == w.mousePressed {
				if button == desktop.RightMouseButton && altTap {
					w.queueEvent(func() { co.(fyne.SecondaryTappable).TappedSecondary(ev) })
				} else if button == desktop.LeftMouseButton && tap {
					w.queueEvent(func() { co.(fyne.Tappable).Tapped(ev) })
				}
			}
			w.mousePressed = nil
		}
	}
	if wid, ok := co.(fyne.Draggable); ok {
		if action == desktop.Press {
			w.mouseDragPos = w.mousePos
			w.mouseDragged = wid
			w.mouseDraggedOffset = w.mousePos.Subtract(co.Position()).Subtract(ev.Position)
		}
	}
	if action == desktop.Release && w.mouseDragged != nil {
		if w.mouseDragStarted {
			w.queueEvent(w.mouseDragged.DragEnd)
			w.mouseDragStarted = false
		}
		if w.objIsDragged(w.mouseOver) && !w.objIsDragged(coMouse) {
			w.mouseOut()
		}
		w.mouseDragged = nil
	}
}

func (w *window) processKeyPressed(keyName fyne.KeyName, scancode int, action desktop.Action, keyDesktopModifier desktop.Modifier) {
	if keyName == "" {
		return
	}

	keyEvent := &fyne.KeyEvent{Name: keyName}

	if keyName == fyne.KeyTab {
		if keyDesktopModifier == 0 {
			if action != desktop.Release {
				w.canvas.focusMgr.FocusNext(w.canvas.focused)
			}
			return
		} else if keyDesktopModifier == desktop.ShiftModifier {
			if action != desktop.Release {
				w.canvas.focusMgr.FocusPrevious(w.canvas.focused)
			}
			return
		}
	}
	if action == desktop.Press {
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.queueEvent(func() { focused.KeyDown(keyEvent) })
			}
		} else if w.canvas.onKeyDown != nil {
			w.queueEvent(func() { w.canvas.onKeyDown(keyEvent) })
		}
	} else if action == desktop.Release { // ignore key up in core events
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.queueEvent(func() { focused.KeyUp(keyEvent) })
			}
		} else if w.canvas.onKeyUp != nil {
			w.queueEvent(func() { w.canvas.onKeyUp(keyEvent) })
		}
		return
	} // key repeat will fall through to TypedKey and TypedShortcut

	var shortcut fyne.Shortcut
	ctrlMod := desktop.ControlModifier
	if runtime.GOOS == "darwin" {
		ctrlMod = desktop.SuperModifier
	}
	if keyDesktopModifier == ctrlMod {
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
	if keyDesktopModifier == desktop.ShiftModifier {
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

	if shortcut == nil && keyDesktopModifier != 0 && keyDesktopModifier != desktop.ShiftModifier {
		shortcut = &desktop.CustomShortcut{
			KeyName:  keyName,
			Modifier: keyDesktopModifier,
		}
	}

	if shortcut != nil {
		if focused, ok := w.canvas.Focused().(fyne.Shortcutable); ok {
			w.queueEvent(func() { focused.TypedShortcut(shortcut) })
			return
		}

		w.queueEvent(func() { w.canvas.shortcut.TypedShortcut(shortcut) })
		return
	}

	// No shortcut detected, pass down to TypedKey
	focused := w.canvas.Focused()
	if focused != nil {
		w.queueEvent(func() { focused.TypedKey(keyEvent) })
	} else if w.canvas.onTypedKey != nil {
		w.queueEvent(func() { w.canvas.onTypedKey(keyEvent) })
	}
}

func (w *window) processCharInput(char rune) {
	if w.canvas.Focused() == nil && w.canvas.onTypedRune == nil {
		return
	}

	focused := w.canvas.Focused()
	if focused != nil {
		w.queueEvent(func() { focused.TypedRune(char) })
	} else if w.canvas.onTypedRune != nil {
		w.queueEvent(func() { w.canvas.onTypedRune(char) })
	}
}

func (w *window) processFocused(focused bool) {
	if w.canvas.focused == nil {
		return
	}

	if focused {
		w.canvas.focused.FocusGained()
	} else {
		w.canvas.focused.FocusLost()
	}
}

func (w *window) RunWithContext(f func()) {
	w.viewport.MakeContextCurrent()

	f()

	w.DetachCurrentContext()
}

func (w *window) RescaleContext() {
	runOnMain(func() {
		w.rescaleOnMain()
	})
}

func (w *window) Context() interface{} {
	return nil
}

// Use this method to queue up a callback that handles an event. This ensures
// user interaction events for a given window are processed in order.
func (w *window) queueEvent(fn func()) {
	w.eventWait.Add(1)
	select {
	case w.eventQueue <- fn:
	default:
		fyne.LogError("EventQueue full", nil)
	}
}

func (w *window) runOnMainWhenCreated(fn func()) {
	if w.viewport != nil {
		runOnMain(fn)
	}

	w.pending = append(w.pending, fn)
}

func (w *window) runEventQueue() {
	w.eventLock.Lock()
	queue := w.eventQueue
	w.eventLock.Unlock()

	for fn := range queue {
		fn()
		w.eventWait.Done()
	}
}

func (w *window) waitForEvents() {
	w.eventWait.Wait()
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

		ret = &window{title: title, decorate: decorate}
		// This channel will be closed when the window is closed.
		ret.eventQueue = make(chan func(), 1024)
		go ret.runEventQueue()

		ret.canvas = newCanvas()
		ret.canvas.context = ret
		ret.SetIcon(ret.icon)
		d.addWindow(ret)
	})
	return ret
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
