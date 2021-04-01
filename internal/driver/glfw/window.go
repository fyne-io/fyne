package glfw

import "C"
import (
	"bytes"
	"context"
	"image"
	_ "image/png" // for the icon
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/widget"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	scrollSpeed      = float32(10)
	doubleClickDelay = 300 // ms (maximum interval between clicks for double click detection)
)

var (
	cursorMap    map[desktop.StandardCursor]*glfw.Cursor
	defaultTitle = "Fyne Application"
)

func initCursors() {
	cursorMap = map[desktop.StandardCursor]*glfw.Cursor{
		desktop.DefaultCursor:   glfw.CreateStandardCursor(glfw.ArrowCursor),
		desktop.TextCursor:      glfw.CreateStandardCursor(glfw.IBeamCursor),
		desktop.CrosshairCursor: glfw.CreateStandardCursor(glfw.CrosshairCursor),
		desktop.PointerCursor:   glfw.CreateStandardCursor(glfw.HandCursor),
		desktop.HResizeCursor:   glfw.CreateStandardCursor(glfw.HResizeCursor),
		desktop.VResizeCursor:   glfw.CreateStandardCursor(glfw.VResizeCursor),
		desktop.HiddenCursor:    nil,
	}
}

// Declare conformity to Window interface
var _ fyne.Window = (*window)(nil)

type window struct {
	viewport   *glfw.Window
	viewLock   sync.RWMutex
	createLock sync.Once
	decorate   bool
	closing    bool
	fixedSize  bool

	cursor       desktop.Cursor
	customCursor *glfw.Cursor
	canvas       *glCanvas
	driver       *gLDriver
	title        string
	icon         fyne.Resource
	mainmenu     *fyne.MainMenu

	clipboard fyne.Clipboard

	master     bool
	fullScreen bool
	centered   bool
	visible    bool

	mousePos             fyne.Position
	mouseDragged         fyne.Draggable
	mouseDraggedObjStart fyne.Position
	mouseDraggedOffset   fyne.Position
	mouseDragPos         fyne.Position
	mouseDragStarted     bool
	mouseButton          desktop.MouseButton
	mouseOver            desktop.Hoverable
	mouseLastClick       fyne.CanvasObject
	mousePressed         fyne.CanvasObject
	mouseClickCount      int
	mouseCancelFunc      context.CancelFunc
	onClosed             func()
	onCloseIntercepted   func()

	menuTogglePending       fyne.KeyName
	menuDeactivationPending fyne.KeyName

	xpos, ypos    int
	width, height int
	shouldExpand  bool

	eventLock  sync.RWMutex
	eventQueue chan func()
	eventWait  sync.WaitGroup
	pending    []func()
}

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

func (w *window) SetFullScreen(full bool) {
	w.fullScreen = full
	if !w.visible {
		return
	}

	runOnMain(func() {
		monitor := w.getMonitorForWindow()
		mode := monitor.GetVideoMode()

		if full {
			w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		} else {
			w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.width, w.height, 0)
		}
	})
}

func (w *window) CenterOnScreen() {
	w.centered = true

	if w.view() != nil {
		runOnMain(w.doCenterOnScreen)
	}
}

func (w *window) doCenterOnScreen() {
	viewWidth, viewHeight := w.screenSize(w.canvas.size)

	// get window dimensions in pixels
	monitor := w.getMonitorForWindow()
	monMode := monitor.GetVideoMode()

	// these come into play when dealing with multiple monitors
	monX, monY := monitor.GetPos()

	// math them to the middle
	newX := (monMode.Width / 2) - (viewWidth / 2) + monX
	newY := (monMode.Height / 2) - (viewHeight / 2) + monY

	// set new window coordinates
	w.viewport.SetPos(newX, newY)
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

func (w *window) RequestFocus() {
	if isWayland {
		return
	}

	w.runOnMainWhenCreated(w.viewport.Focus)
}

func (w *window) Resize(size fyne.Size) {
	// we cannot perform this until window is prepared as we don't know it's scale!

	w.runOnMainWhenCreated(func() {
		w.canvas.Resize(size)
		w.viewLock.Lock()

		width, height := internal.ScaleInt(w.canvas, size.Width), internal.ScaleInt(w.canvas, size.Height)
		if w.fixedSize || !w.visible { // fixed size ignores future `resized` and if not visible we may not get the event
			w.width, w.height = width, height
		}
		w.viewLock.Unlock()

		w.viewport.SetSize(width, height)
		w.fitContent()
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed

	if w.view() != nil {
		w.fitContent()
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

func (w *window) SetIcon(icon fyne.Resource) {
	w.icon = icon
	if icon == nil {
		appIcon := fyne.CurrentApp().Icon()
		if appIcon != nil {
			w.SetIcon(appIcon)
		}
		return
	}

	if string(icon.Content()[:4]) == "<svg" {
		fyne.LogError("Window icon does not support vector images", nil)
		return
	}

	w.runOnMainWhenCreated(func() {
		if w.icon == nil {
			w.viewport.SetIcon(nil)
			return
		}

		pix, _, err := image.Decode(bytes.NewReader(w.icon.Content()))
		if err != nil {
			fyne.LogError("Failed to decode image for window icon", err)
			return
		}

		w.viewport.SetIcon([]image.Image{pix})
	})
}

func (w *window) SetMaster() {
	w.master = true
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

func (w *window) fitContent() {
	if w.canvas.Content() == nil {
		return
	}

	if w.isClosing() {
		return
	}

	minWidth, minHeight := w.minSizeOnScreen()
	w.viewLock.RLock()
	view := w.viewport
	w.viewLock.RUnlock()
	if w.width < minWidth || w.height < minHeight {
		if w.width < minWidth {
			w.width = minWidth
		}
		if w.height < minHeight {
			w.height = minHeight
		}
		w.viewLock.Lock()
		w.shouldExpand = true // queue the resize to happen on main
		w.viewLock.Unlock()
	}
	if w.fixedSize {
		w.width = internal.ScaleInt(w.canvas, w.Canvas().Size().Width)
		w.height = internal.ScaleInt(w.canvas, w.Canvas().Size().Height)

		view.SetSizeLimits(w.width, w.height, w.width, w.height)
	} else {
		view.SetSizeLimits(minWidth, minHeight, glfw.DontCare, glfw.DontCare)
	}
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) SetCloseIntercept(callback func()) {
	w.onCloseIntercepted = callback
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	x, y := w.xpos, w.ypos
	if w.fullScreen {
		x, y = w.viewport.GetPos()
	}
	xOff := x + (w.width / 2)
	yOff := y + (w.height / 2)

	for _, monitor := range glfw.GetMonitors() {
		x, y := monitor.GetPos()

		if x > xOff || y > yOff {
			continue
		}
		if x+monitor.GetVideoMode().Width <= xOff || y+monitor.GetVideoMode().Height <= yOff {
			continue
		}

		return monitor
	}

	// try built-in function to detect monitor if above logic didn't succeed
	// if it doesn't work then return primary monitor as default
	monitor := w.viewport.GetMonitor()
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}
	return monitor
}

func (w *window) calculatedScale() float32 {
	return calculateScale(userScale(), fyne.CurrentDevice().SystemScaleForWindow(w), w.detectScale())
}

func (w *window) detectScale() float32 {
	monitor := w.getMonitorForWindow()
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	return calculateDetectedScale(widthMm, widthPx)
}

func (w *window) detectTextureScale() float32 {
	winWidth, _ := w.viewport.GetSize()
	texWidth, _ := w.viewport.GetFramebufferSize()
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

	for !running() {
		time.Sleep(time.Millisecond * 10)
	}
	w.createLock.Do(w.create)
	if w.view() == nil {
		return
	}

	runOnMain(func() {
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
		w.viewport.SetTitle(w.title)

		if w.centered {
			w.doCenterOnScreen() // lastly center if that was requested
		}
		w.viewport.Show()

		// save coordinates
		w.xpos, w.ypos = w.viewport.GetPos()

		if w.fullScreen { // this does not work if called before viewport.Show()
			go func() {
				time.Sleep(time.Millisecond * 100)
				w.SetFullScreen(true)
			}()
		}
	})

	// show top canvas element
	if w.canvas.Content() != nil {
		w.canvas.Content().Show()
	}
}

func (w *window) Hide() {
	if w.isClosing() {
		return
	}

	runOnMain(func() {
		w.viewLock.Lock()
		w.visible = false
		w.viewport.Hide()
		w.viewLock.Unlock()

		// hide top canvas element
		if w.canvas.Content() != nil {
			w.canvas.Content().Hide()
		}
	})
}

func (w *window) Close() {
	if w.isClosing() {
		return
	}

	w.closing = true
	w.viewport.SetShouldClose(true)

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
	w.RescaleContext()
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(false)

	if w.onCloseIntercepted != nil {
		w.queueEvent(w.onCloseIntercepted)
		return
	}

	w.Close()
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

	if w.master {
		d.Quit()
	} else if runtime.GOOS == "darwin" {
		go d.focusPreviousWindow()
	}
}

func (w *window) moved(_ *glfw.Window, x, y int) {
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

func (w *window) resized(_ *glfw.Window, width, height int) {
	if w.fixedSize {
		return
	}

	canvasSize := fyne.NewSize(internal.UnscaleInt(w.canvas, width), internal.UnscaleInt(w.canvas, height))
	if !w.fullScreen {
		w.width = internal.ScaleInt(w.canvas, canvasSize.Width)
		w.height = internal.ScaleInt(w.canvas, canvasSize.Height)
	}

	if !w.visible { // don't redraw if hidden
		w.canvas.Resize(canvasSize)
		return
	}

	w.platformResize(canvasSize)
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	if width == 0 || height == 0 || runtime.GOOS != "darwin" {
		return
	}

	winWidth, _ := viewport.GetSize()
	newTexScale := float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	if w.canvas.texScale != newTexScale {
		w.canvas.texScale = newTexScale
		w.canvas.Refresh(w.canvas.Content()) // reset graphics to apply texture scale
	}
}

func (w *window) refresh(_ *glfw.Window) {
	refreshWindow(w)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse fyne.Position, matches func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position, int) {
	return driver.FindObjectAtPositionMatching(mouse, matches, canvas.Overlays().Top(), canvas.menu, canvas.Content())
}

func fyneToNativeCursor(cursor desktop.Cursor) (*glfw.Cursor, bool) {
	switch v := cursor.(type) {
	case desktop.StandardCursor:
		ret, ok := cursorMap[v]
		if !ok {
			return cursorMap[desktop.DefaultCursor], false
		}
		return ret, false
	default:
		img, x, y := cursor.Image()
		if img == nil {
			return nil, true
		}
		return glfw.CreateCursor(img, x, y), true
	}
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	previousPos := w.mousePos
	w.mousePos = fyne.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))

	cursor := desktop.Cursor(desktop.DefaultCursor)

	obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		if cursorable, ok := object.(desktop.Cursorable); ok {
			cursor = cursorable.Cursor()
		}
		if _, ok := object.(fyne.Draggable); ok {
			return true
		}

		_, hover := object.(desktop.Hoverable)
		return hover
	})

	if w.cursor != cursor {
		// cursor has changed, store new cursor and apply change via glfw
		rawCursor, isCustomCursor := fyneToNativeCursor(cursor)
		w.cursor = cursor

		if rawCursor == nil {
			viewport.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		} else {
			viewport.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
			viewport.SetCursor(rawCursor)
		}
		if w.customCursor != nil {
			w.customCursor.Destroy()
			w.customCursor = nil
		}
		if isCustomCursor {
			w.customCursor = rawCursor
		}
	}

	if w.mouseButton != 0 && !w.mouseDragStarted {
		obj, pos, _ := w.findObjectAtPositionMatching(w.canvas, previousPos, func(object fyne.CanvasObject) bool {
			_, ok := object.(fyne.Draggable)
			return ok
		})

		if wid, ok := obj.(fyne.Draggable); ok {
			w.mouseDragPos = previousPos
			w.mouseDragged = wid
			w.mouseDraggedOffset = previousPos.Subtract(pos)
			w.mouseDraggedObjStart = obj.Position()
			w.mouseDragStarted = true
		}
	}

	if obj != nil && !w.objIsDragged(obj) {
		ev := new(desktop.MouseEvent)
		ev.AbsolutePosition = w.mousePos
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
			draggedObjDelta := w.mouseDraggedObjStart.Subtract(w.mouseDragged.(fyne.CanvasObject).Position())
			ev := new(fyne.DragEvent)
			ev.AbsolutePosition = w.mousePos
			ev.Position = w.mousePos.Subtract(w.mouseDraggedOffset).Add(draggedObjDelta)
			ev.Dragged = fyne.NewDelta(w.mousePos.X-w.mouseDragPos.X, w.mousePos.Y-w.mouseDragPos.Y)
			wd := w.mouseDragged
			w.queueEvent(func() { wd.Dragged(ev) })

			w.mouseDragStarted = true
			w.mouseDragPos = w.mousePos
		}
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

func (w *window) mouseClicked(_ *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if w.mousePos.IsZero() { // window may not be focused (darwin mostly) and so position callbacks not happening
		xpos, ypos := w.viewport.GetCursorPos()
		w.mousePos = fyne.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))
	}

	co, pos, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		switch object.(type) {
		case fyne.Tappable, fyne.SecondaryTappable, fyne.DoubleTappable, fyne.Focusable, desktop.Mouseable, desktop.Hoverable:
			return true
		case fyne.Draggable:
			if w.mouseDragStarted {
				return true
			}
		}

		return false
	})
	ev := new(fyne.PointEvent)
	ev.Position = pos
	ev.AbsolutePosition = w.mousePos

	coMouse := co
	button, modifiers := convertMouseButton(btn, mods)
	if wid, ok := co.(desktop.Mouseable); ok {
		mev := new(desktop.MouseEvent)
		mev.Position = ev.Position
		mev.AbsolutePosition = w.mousePos
		mev.Button = button
		mev.Modifier = modifiers
		if action == glfw.Press {
			w.queueEvent(func() { wid.MouseDown(mev) })
		} else if action == glfw.Release {
			if w.mouseDragged == nil {
				w.queueEvent(func() { wid.MouseUp(mev) })
			} else {
				if dragged, ok := w.mouseDragged.(desktop.Mouseable); ok {
					mev.Position = w.mousePos.Subtract(w.mouseDraggedOffset)
					w.queueEvent(func() { dragged.MouseUp(mev) })
				} else {
					w.queueEvent(func() { wid.MouseUp(mev) })
				}
			}
		}
	}

	if wid, ok := co.(fyne.Focusable); ok {
		w.canvas.Focus(wid)
	} else {
		w.canvas.Unfocus()
	}

	if action == glfw.Press {
		w.mouseButton |= button
	} else if action == glfw.Release {
		w.mouseButton &= ^button
	}

	if action == glfw.Release && w.mouseDragged != nil {
		if w.mouseDragStarted {
			w.queueEvent(w.mouseDragged.DragEnd)
			w.mouseDragStarted = false
		}
		if w.objIsDragged(w.mouseOver) && !w.objIsDragged(coMouse) {
			w.mouseOut()
		}
		w.mouseDragged = nil
	}
	_, tap := co.(fyne.Tappable)
	_, altTap := co.(fyne.SecondaryTappable)
	if tap || altTap {
		if action == glfw.Press {
			w.mousePressed = co
		} else if action == glfw.Release {
			if co == w.mousePressed {
				if button == desktop.MouseButtonSecondary && altTap {
					w.queueEvent(func() { co.(fyne.SecondaryTappable).TappedSecondary(ev) })
				}
			}
		}
	}

	// Check for double click/tap on left mouse button
	if action == glfw.Release && button == desktop.MouseButtonPrimary {
		_, doubleTap := co.(fyne.DoubleTappable)
		if doubleTap {
			w.mouseClickCount++
			w.mouseLastClick = co
			if w.mouseCancelFunc != nil {
				w.mouseCancelFunc()
				return
			}
			go w.waitForDoubleTap(co, ev)
		} else {
			if wid, ok := co.(fyne.Tappable); ok && co == w.mousePressed {
				w.queueEvent(func() { wid.Tapped(ev) })
			}
			w.mousePressed = nil
		}
	}
}

func (w *window) waitForDoubleTap(co fyne.CanvasObject, ev *fyne.PointEvent) {
	var ctx context.Context
	ctx, w.mouseCancelFunc = context.WithDeadline(context.TODO(), time.Now().Add(time.Millisecond*doubleClickDelay))
	defer w.mouseCancelFunc()

	<-ctx.Done()
	if w.mouseClickCount == 2 && w.mouseLastClick == co {
		if wid, ok := co.(fyne.DoubleTappable); ok {
			w.queueEvent(func() { wid.DoubleTapped(ev) })
		}
	} else if co == w.mousePressed {
		if wid, ok := co.(fyne.Tappable); ok {
			w.queueEvent(func() { wid.Tapped(ev) })
		}
	}
	w.mouseClickCount = 0
	w.mousePressed = nil
	w.mouseCancelFunc = nil
	w.mouseLastClick = nil
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
	co, _, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		_, ok := object.(fyne.Scrollable)
		return ok
	})
	switch wid := co.(type) {
	case fyne.Scrollable:
		if runtime.GOOS != "darwin" && xoff == 0 &&
			(viewport.GetKey(glfw.KeyLeftShift) == glfw.Press ||
				viewport.GetKey(glfw.KeyRightShift) == glfw.Press) {
			xoff, yoff = yoff, xoff
		}
		ev := &fyne.ScrollEvent{}
		ev.Scrolled = fyne.NewDelta(float32(xoff)*scrollSpeed, float32(yoff)*scrollSpeed)
		wid.Scrolled(ev)
	}
}

func convertMouseButton(btn glfw.MouseButton, mods glfw.ModifierKey) (desktop.MouseButton, desktop.Modifier) {
	modifier := desktopModifier(mods)
	var button desktop.MouseButton
	rightClick := false
	if runtime.GOOS == "darwin" {
		if modifier&desktop.ControlModifier != 0 {
			rightClick = true
			modifier &^= desktop.ControlModifier
		}
		if modifier&desktop.SuperModifier != 0 {
			modifier |= desktop.ControlModifier
			modifier &^= desktop.SuperModifier
		}
	}
	switch btn {
	case glfw.MouseButton1:
		if rightClick {
			button = desktop.MouseButtonSecondary
		} else {
			button = desktop.MouseButtonPrimary
		}
	case glfw.MouseButton2:
		button = desktop.MouseButtonSecondary
	case glfw.MouseButton3:
		button = desktop.MouseButtonTertiary
	}
	return button, modifier
}

var keyCodeMap = map[glfw.Key]fyne.KeyName{
	// non-printable
	glfw.KeyEscape:    fyne.KeyEscape,
	glfw.KeyEnter:     fyne.KeyReturn,
	glfw.KeyTab:       fyne.KeyTab,
	glfw.KeyBackspace: fyne.KeyBackspace,
	glfw.KeyInsert:    fyne.KeyInsert,
	glfw.KeyDelete:    fyne.KeyDelete,
	glfw.KeyRight:     fyne.KeyRight,
	glfw.KeyLeft:      fyne.KeyLeft,
	glfw.KeyDown:      fyne.KeyDown,
	glfw.KeyUp:        fyne.KeyUp,
	glfw.KeyPageUp:    fyne.KeyPageUp,
	glfw.KeyPageDown:  fyne.KeyPageDown,
	glfw.KeyHome:      fyne.KeyHome,
	glfw.KeyEnd:       fyne.KeyEnd,

	glfw.KeySpace:   fyne.KeySpace,
	glfw.KeyKPEnter: fyne.KeyEnter,

	// functions
	glfw.KeyF1:  fyne.KeyF1,
	glfw.KeyF2:  fyne.KeyF2,
	glfw.KeyF3:  fyne.KeyF3,
	glfw.KeyF4:  fyne.KeyF4,
	glfw.KeyF5:  fyne.KeyF5,
	glfw.KeyF6:  fyne.KeyF6,
	glfw.KeyF7:  fyne.KeyF7,
	glfw.KeyF8:  fyne.KeyF8,
	glfw.KeyF9:  fyne.KeyF9,
	glfw.KeyF10: fyne.KeyF10,
	glfw.KeyF11: fyne.KeyF11,
	glfw.KeyF12: fyne.KeyF12,

	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	glfw.Key0:   fyne.Key0,
	glfw.KeyKP0: fyne.Key0,
	glfw.Key1:   fyne.Key1,
	glfw.KeyKP1: fyne.Key1,
	glfw.Key2:   fyne.Key2,
	glfw.KeyKP2: fyne.Key2,
	glfw.Key3:   fyne.Key3,
	glfw.KeyKP3: fyne.Key3,
	glfw.Key4:   fyne.Key4,
	glfw.KeyKP4: fyne.Key4,
	glfw.Key5:   fyne.Key5,
	glfw.KeyKP5: fyne.Key5,
	glfw.Key6:   fyne.Key6,
	glfw.KeyKP6: fyne.Key6,
	glfw.Key7:   fyne.Key7,
	glfw.KeyKP7: fyne.Key7,
	glfw.Key8:   fyne.Key8,
	glfw.KeyKP8: fyne.Key8,
	glfw.Key9:   fyne.Key9,
	glfw.KeyKP9: fyne.Key9,

	// desktop
	glfw.KeyLeftShift:    desktop.KeyShiftLeft,
	glfw.KeyRightShift:   desktop.KeyShiftRight,
	glfw.KeyLeftControl:  desktop.KeyControlLeft,
	glfw.KeyRightControl: desktop.KeyControlRight,
	glfw.KeyLeftAlt:      desktop.KeyAltLeft,
	glfw.KeyRightAlt:     desktop.KeyAltRight,
	glfw.KeyLeftSuper:    desktop.KeySuperLeft,
	glfw.KeyRightSuper:   desktop.KeySuperRight,
	glfw.KeyMenu:         desktop.KeyMenu,
	glfw.KeyPrintScreen:  desktop.KeyPrintScreen,
	glfw.KeyCapsLock:     desktop.KeyCapsLock,
}

var keyNameMap = map[string]fyne.KeyName{
	"'": fyne.KeyApostrophe,
	",": fyne.KeyComma,
	"-": fyne.KeyMinus,
	".": fyne.KeyPeriod,
	"/": fyne.KeySlash,
	"*": fyne.KeyAsterisk,
	"`": fyne.KeyBackTick,

	";": fyne.KeySemicolon,
	"+": fyne.KeyPlus,
	"=": fyne.KeyEqual,

	"a": fyne.KeyA,
	"b": fyne.KeyB,
	"c": fyne.KeyC,
	"d": fyne.KeyD,
	"e": fyne.KeyE,
	"f": fyne.KeyF,
	"g": fyne.KeyG,
	"h": fyne.KeyH,
	"i": fyne.KeyI,
	"j": fyne.KeyJ,
	"k": fyne.KeyK,
	"l": fyne.KeyL,
	"m": fyne.KeyM,
	"n": fyne.KeyN,
	"o": fyne.KeyO,
	"p": fyne.KeyP,
	"q": fyne.KeyQ,
	"r": fyne.KeyR,
	"s": fyne.KeyS,
	"t": fyne.KeyT,
	"u": fyne.KeyU,
	"v": fyne.KeyV,
	"w": fyne.KeyW,
	"x": fyne.KeyX,
	"y": fyne.KeyY,
	"z": fyne.KeyZ,

	"[":  fyne.KeyLeftBracket,
	"\\": fyne.KeyBackslash,
	"]":  fyne.KeyRightBracket,
}

func keyToName(code glfw.Key, scancode int) fyne.KeyName {
	if runtime.GOOS == "darwin" && scancode == 0x69 { // TODO remove once fixed upstream glfw/glfw#1786
		code = glfw.KeyPrintScreen
	}

	ret, ok := keyCodeMap[code]
	if ok {
		return ret
	}

	keyName := glfw.GetKeyName(code, scancode)
	ret, ok = keyNameMap[keyName]
	if !ok {
		return ""
	}

	return ret
}

func (w *window) keyPressed(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	if keyName == "" {
		return
	}

	keyEvent := &fyne.KeyEvent{Name: keyName}
	keyDesktopModifier := desktopModifier(mods)
	pendingMenuToggle := w.menuTogglePending
	pendingMenuDeactivation := w.menuDeactivationPending
	w.menuTogglePending = desktop.KeyNone
	w.menuDeactivationPending = desktop.KeyNone
	switch action {
	case glfw.Release:
		if action == glfw.Release {
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
				w.queueEvent(func() { focused.KeyUp(keyEvent) })
			}
		} else if w.canvas.onKeyUp != nil {
			w.queueEvent(func() { w.canvas.onKeyUp(keyEvent) })
		}
		return // ignore key up in other core events
	case glfw.Press:
		switch keyName {
		case desktop.KeyAltLeft, desktop.KeyAltRight:
			if (keyName == desktop.KeyAltLeft || keyName == desktop.KeyAltRight) && keyDesktopModifier == desktop.AltModifier {
				w.menuTogglePending = keyName
			}
		case fyne.KeyEscape:
			w.menuDeactivationPending = keyName
		}
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.queueEvent(func() { focused.KeyDown(keyEvent) })
			}
		} else if w.canvas.onKeyDown != nil {
			w.queueEvent(func() { w.canvas.onKeyDown(keyEvent) })
		}
	default:
		// key repeat will fall through to TypedKey and TypedShortcut
	}

	switch keyName {
	case fyne.KeyTab:
		capture := false
		// TODO at some point allow widgets to mark as capturing
		if ent, ok := w.canvas.Focused().(*widget.Entry); ok && ent.MultiLine {
			capture = true
		}
		if !capture {
			switch keyDesktopModifier {
			case 0:
				w.canvas.FocusNext()
				return
			case desktop.ShiftModifier:
				w.canvas.FocusPrevious()
				return
			}
		}
	}

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
			shouldRunShortcut := true
			type selectableText interface {
				fyne.Disableable
				SelectedText() string
			}
			if selectableTextWid, ok := focused.(selectableText); ok && selectableTextWid.Disabled() {
				shouldRunShortcut = shortcut.ShortcutName() == "Copy"
			}
			if shouldRunShortcut {
				w.queueEvent(func() { focused.TypedShortcut(shortcut) })
			}
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

func desktopModifier(mods glfw.ModifierKey) desktop.Modifier {
	var m desktop.Modifier
	if (mods & glfw.ModShift) != 0 {
		m |= desktop.ShiftModifier
	}
	if (mods & glfw.ModControl) != 0 {
		m |= desktop.ControlModifier
	}
	if (mods & glfw.ModAlt) != 0 {
		m |= desktop.AltModifier
	}
	if (mods & glfw.ModSuper) != 0 {
		m |= desktop.SuperModifier
	}
	return m
}

// charInput defines the character with modifiers callback which is called when a
// Unicode character is input.
//
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or more characters.
func (w *window) charInput(_ *glfw.Window, char rune) {
	if focused := w.canvas.Focused(); focused != nil {
		w.queueEvent(func() { focused.TypedRune(char) })
	} else if w.canvas.onTypedRune != nil {
		w.queueEvent(func() { w.canvas.onTypedRune(char) })
	}
}

func (w *window) focused(_ *glfw.Window, isFocused bool) {
	if isFocused {
		w.canvas.FocusGained()
	} else {
		w.canvas.FocusLost()
		w.mousePos = fyne.Position{}
	}
}

func (w *window) RunWithContext(f func()) {
	if w.isClosing() {
		return
	}
	w.viewport.MakeContextCurrent()

	f()

	glfw.DetachCurrentContext()
}

func (w *window) RescaleContext() {
	runOnMain(func() {
		w.rescaleOnMain()
	})
}

func (w *window) rescaleOnMain() {
	if w.isClosing() {
		return
	}
	w.fitContent()

	if w.fullScreen {
		w.width, w.height = w.viewport.GetSize()
		scaledFull := fyne.NewSize(
			internal.UnscaleInt(w.canvas, w.width),
			internal.UnscaleInt(w.canvas, w.height))
		w.canvas.Resize(scaledFull)
		return
	}

	size := w.canvas.size.Max(w.canvas.MinSize())
	newWidth, newHeight := w.screenSize(size)
	w.viewport.SetSize(newWidth, newHeight)
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
		fyne.LogError("EventQueue full, perhaps a callback blocked the event handler", nil)
	}
}

func (w *window) runOnMainWhenCreated(fn func()) {
	if w.viewport != nil {
		runOnMain(fn)
		return
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

		ret = &window{title: title, decorate: decorate, driver: d}
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

func (w *window) create() {
	runOnMain(func() {
		if !isWayland {
			// make the window hidden, we will set it up and then show it later
			glfw.WindowHint(glfw.Visible, 0)
		}
		if w.decorate {
			glfw.WindowHint(glfw.Decorated, 1)
		} else {
			glfw.WindowHint(glfw.Decorated, 0)
		}
		if w.fixedSize {
			glfw.WindowHint(glfw.Resizable, 0)
		} else {
			glfw.WindowHint(glfw.Resizable, 1)
		}
		initWindowHints()

		pixWidth, pixHeight := w.screenSize(w.canvas.size)
		pixWidth = int(fyne.Max(float32(pixWidth), float32(w.width)))
		if pixWidth == 0 {
			pixWidth = 10
		}
		pixHeight = int(fyne.Max(float32(pixHeight), float32(w.height)))
		if pixHeight == 0 {
			pixHeight = 10
		}

		win, err := glfw.CreateWindow(pixWidth, pixHeight, w.title, nil, nil)
		if err != nil {
			w.driver.initFailed("window creation error", err)
			return
		}

		w.viewLock.Lock()
		w.viewport = win
		w.viewLock.Unlock()
	})
	if w.view() == nil { // something went wrong above, it will have been logged
		return
	}

	// run the GL init on the draw thread
	runOnDraw(w, func() {
		w.canvas.painter = gl.NewPainter(w.canvas, w)
		w.canvas.painter.Init()
	})

	runOnMain(func() {
		win := w.view()
		win.SetCloseCallback(w.closed)
		win.SetPosCallback(w.moved)
		win.SetSizeCallback(w.resized)
		win.SetFramebufferSizeCallback(w.frameSized)
		win.SetRefreshCallback(w.refresh)
		win.SetCursorPosCallback(w.mouseMoved)
		win.SetMouseButtonCallback(w.mouseClicked)
		win.SetScrollCallback(w.mouseScrolled)
		win.SetKeyCallback(w.keyPressed)
		win.SetCharCallback(w.charInput)
		win.SetFocusCallback(w.focused)

		w.canvas.detectedScale = w.detectScale()
		w.canvas.scale = w.calculatedScale()
		w.canvas.texScale = w.detectTextureScale()
		// update window size now we have scaled detected
		w.fitContent()

		for _, fn := range w.pending {
			fn()
		}

		if w.fixedSize { // as the window will not be sized later we may need to pack menus etc
			w.canvas.Resize(w.canvas.Size())
		}
		// order of operation matters so we do these last items in order
		w.viewport.SetSize(w.width, w.height) // ensure we requested latest size
	})
}

func (w *window) doShowAgain() {
	if w.isClosing() {
		return
	}

	runOnMain(func() {
		// show top canvas element
		if w.canvas.Content() != nil {
			w.canvas.Content().Show()
		}

		w.viewport.SetPos(w.xpos, w.ypos)
		w.viewport.Show()
		w.viewLock.Lock()
		w.visible = true
		w.viewLock.Unlock()
	})
}

func (w *window) isClosing() bool {
	return w.closing || w.viewport == nil
}

func (w *window) view() *glfw.Window {
	w.viewLock.RLock()
	defer w.viewLock.RUnlock()

	if w.closing {
		return nil
	}
	return w.viewport
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
