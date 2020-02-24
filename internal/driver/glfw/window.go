package glfw

import "C"
import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/widget"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	scrollSpeed      = 10
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

var (
	defaultCursor, entryCursor, hyperlinkCursor *glfw.Cursor
	initOnce                                    = &sync.Once{}
	defaultTitle                                = "Fyne Application"
)

func initCursors() {
	defaultCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	entryCursor = glfw.CreateStandardCursor(glfw.IBeamCursor)
	hyperlinkCursor = glfw.CreateStandardCursor(glfw.HandCursor)
}

// Declare conformity to Window interface
var _ fyne.Window = (*window)(nil)

type window struct {
	viewport *glfw.Window
	painted  int // part of the macOS GL fix, updated GLFW should fix this
	canvas   *glCanvas
	title    string
	icon     fyne.Resource
	mainmenu *fyne.MainMenu

	clipboard fyne.Clipboard

	master     bool
	fullScreen bool
	fixedSize  bool
	centered   bool
	visible    bool

	mousePos           fyne.Position
	mouseDragged       fyne.Draggable
	mouseDraggedOffset fyne.Position
	mouseDragPos       fyne.Position
	mouseDragStarted   bool
	mouseButton        desktop.MouseButton
	mouseOver          desktop.Hoverable
	mouseClickTime     time.Time
	mouseLastClick     fyne.CanvasObject
	mousePressed       fyne.Tappable
	onClosed           func()

	xpos, ypos    int
	width, height int
	ignoreResize  bool

	eventQueue chan func()
	eventWait  sync.WaitGroup
}

func (w *window) Title() string {
	return w.title
}

func (w *window) SetTitle(title string) {
	w.title = title
	runOnMain(func() {
		w.viewport.SetTitle(title)
	})
}

func (w *window) FullScreen() bool {
	return w.fullScreen
}

func (w *window) SetFullScreen(full bool) {
	if full {
		w.fullScreen = true
	}
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
			w.fullScreen = false
		}
	})
}

func (w *window) CenterOnScreen() {
	w.centered = true
	// if window is currently visible, make it centered
	if w.visible {
		w.centerOnScreen()
	}
}

// centerOnScreen handles the logic for centering a window
func (w *window) centerOnScreen() {
	runOnMain(func() {
		viewWidth, viewHeight := w.viewport.GetSize()

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
	}) // end of runOnMain(){}
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
	runOnMain(func() {
		err := w.viewport.Focus()
		if err != nil {
			fyne.LogError("Error requesting focus", err)
		}
	})
}

func (w *window) Resize(size fyne.Size) {
	w.canvas.Resize(size)
	w.width, w.height = internal.ScaleInt(w.canvas, size.Width), internal.ScaleInt(w.canvas, size.Height)
	runOnMain(func() {
		w.ignoreResize = true
		w.viewport.SetSize(w.width, w.height)
		w.ignoreResize = false
		w.fitContent()
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
	runOnMain(w.fitContent)
}

func (w *window) Padded() bool {
	return w.canvas.padded
}

func (w *window) SetPadded(padded bool) {
	w.canvas.SetPadded(padded)

	runOnMain(w.fitContent)
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

	pix, _, err := image.Decode(bytes.NewReader(icon.Content()))
	if err != nil {
		fyne.LogError("Failed to decode image for window icon", err)
		return
	}

	w.viewport.SetIcon([]image.Image{pix})
}

func (w *window) SetMaster() {
	w.master = true
}

func (w *window) MainMenu() *fyne.MainMenu {
	return w.mainmenu
}

func (w *window) SetMainMenu(menu *fyne.MainMenu) {
	w.mainmenu = menu
	w.canvas.buildMenuBar(menu)
}

func (w *window) fitContent() {
	w.canvas.RLock()
	content := w.canvas.content
	w.canvas.RUnlock()
	if content == nil {
		return
	}

	w.ignoreResize = true
	minWidth, minHeight := w.minSizeOnScreen()
	if w.width < minWidth || w.height < minHeight {
		if w.width < minWidth {
			w.width = minWidth
		}
		if w.height < minHeight {
			w.height = minHeight
		}
		w.viewport.SetSize(w.width, w.height)
	}
	if w.fixedSize {
		w.width = internal.ScaleInt(w.canvas, w.Canvas().Size().Width)
		w.height = internal.ScaleInt(w.canvas, w.Canvas().Size().Height)

		w.viewport.SetSizeLimits(w.width, w.height, w.width, w.height)
	} else {
		w.viewport.SetSizeLimits(minWidth, minHeight, glfw.DontCare, glfw.DontCare)
	}
	w.ignoreResize = false
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	xOff := w.xpos + (w.width / 2)
	yOff := w.ypos + (w.height / 2)

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

func (w *window) userScale() float32 {
	env := os.Getenv("FYNE_SCALE")

	if env != "" && env != "auto" {
		scale, err := strconv.ParseFloat(env, 32)
		if err == nil && scale != 0 {
			return float32(scale)
		}
		fyne.LogError("Error reading scale", err)
	}

	if env != "auto" {
		setting := fyne.CurrentApp().Settings().Scale()
		if setting != fyne.SettingsScaleAuto && setting != 0.0 {
			return setting
		}
	}

	return 1.0 // user preference for auto is now passed as 1 so the system auto is picked up
}

func calculateScale(user, system, detected float32) float32 {
	if user == fyne.SettingsScaleAuto {
		user = 1.0
	}

	if system == fyne.SettingsScaleAuto {
		system = detected
	}

	return system * user
}
func (w *window) calculatedScale() float32 {
	val := calculateScale(w.userScale(), fyne.CurrentDevice().SystemScale(), w.detectScale())
	val = float32(math.Round(float64(val*10.0))) / 10.0

	return val
}

func (w *window) detectScale() float32 {
	monitor := w.getMonitorForWindow()
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	dpi := float32(widthPx) / (float32(widthMm) / 25.4)
	if dpi > 1000 || dpi < 10 {
		dpi = 96
	}
	return float32(float64(dpi) / 96.0)
}

func (w *window) Show() {
	if w.centered {
		w.centerOnScreen()
	}

	runOnMain(func() {
		w.visible = true
		w.viewport.Show()
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
	w.closed(w.viewport)
}

func (w *window) ShowAndRun() {
	w.Show()
	fyne.CurrentApp().Driver().Run()
}

//Clipboard returns the system clipboard
func (w *window) Clipboard() fyne.Clipboard {
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
		w.width = internal.ScaleInt(w.canvas, canvasSize.Width)
		w.height = internal.ScaleInt(w.canvas, canvasSize.Height)
	}

	w.canvas.Resize(canvasSize)
}

func (w *window) SetContent(content fyne.CanvasObject) {
	// hide old canvas element
	if w.visible && w.canvas.Content() != nil {
		w.canvas.Content().Hide()
	}

	w.canvas.SetContent(content)
	w.RescaleContext()
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

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

// destroy this window and, if it's the last window quit the app
func (w *window) destroy(d *gLDriver) {
	// finish serial event queue and nil it so we don't panic if window.closed() is called twice.
	if w.eventQueue != nil {
		w.waitForEvents()
		close(w.eventQueue)
		w.eventQueue = nil
	}

	if w.master || len(d.windows) == 0 {
		d.Quit()
	}
}

func (w *window) moved(viewport *glfw.Window, x, y int) {
	// save coordinates
	w.xpos, w.ypos = x, y

	newDetected := w.detectScale()
	if w.canvas.detectedScale == newDetected {
		return
	}

	w.canvas.detectedScale = newDetected
	go w.canvas.SetScale(fyne.SettingsScaleAuto) // scale value is ignored
}

func (w *window) resized(viewport *glfw.Window, width, height int) {
	if w.ignoreResize {
		return
	}
	w.resize(fyne.NewSize(internal.UnscaleInt(w.canvas, width), internal.UnscaleInt(w.canvas, height)))
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	if width == 0 || height == 0 {
		return
	}

	winWidth, _ := w.viewport.GetSize()
	texScale := float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	w.canvas.setTextureScale(texScale)
	w.canvas.painter.SetOutputSize(width, height)
}

func (w *window) refresh(viewport *glfw.Window) {
	forceWindowRefresh(w)
	w.canvas.setDirty(true)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse fyne.Position,
	matches func(object fyne.CanvasObject) bool) (fyne.CanvasObject, fyne.Position) {
	roots := []fyne.CanvasObject{canvas.content}

	if canvas.menu != nil {
		roots = []fyne.CanvasObject{canvas.menu, canvas.content}
	}

	return driver.FindObjectAtPositionMatching(mouse, matches, canvas.overlay, roots...)
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mousePos = fyne.NewPos(internal.UnscaleInt(w.canvas, int(xpos)), internal.UnscaleInt(w.canvas, int(ypos)))

	cursor := defaultCursor
	obj, pos := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		if wid, ok := object.(*widget.Entry); ok {
			if !wid.Disabled() {
				cursor = entryCursor
			}
		} else if _, ok := object.(*widget.Hyperlink); ok {
			cursor = hyperlinkCursor
		}

		_, hover := object.(desktop.Hoverable)
		return hover
	})

	viewport.SetCursor(cursor)
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

func (w *window) mouseClicked(viewport *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	co, pos := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		if _, ok := object.(fyne.Tappable); ok {
			return true
		} else if _, ok := object.(fyne.Focusable); ok {
			return true
		} else if _, ok := object.(fyne.Draggable); ok {
			return true
		} else if _, ok := object.(desktop.Mouseable); ok {
			return true
		} else if _, ok := object.(desktop.Hoverable); ok {
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
			w.queueEvent(func() { wid.MouseUp(mev) })
		}
	}

	needsfocus := true
	wid := w.canvas.Focused()
	if wid != nil {
		if wid.(fyne.CanvasObject) != co {
			w.canvas.Unfocus()
		} else {
			needsfocus = false
		}
	}

	if action == glfw.Press {
		w.mouseButton = button
	} else if action == glfw.Release {
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
	if action == glfw.Release && button == desktop.LeftMouseButton {
		now := time.Now()
		// we can safely subtract the first "zero" time as it'll be much larger than doubleClickDelay
		if now.Sub(w.mouseClickTime).Nanoseconds()/1e6 <= doubleClickDelay && w.mouseLastClick == co {
			if wid, ok := co.(fyne.DoubleTappable); ok {
				doubleTapped = true
				w.queueEvent(func() { wid.DoubleTapped(ev) })
			}
		}
		w.mouseClickTime = now
		w.mouseLastClick = co
	}

	// Prevent Tapped from triggering if DoubleTapped has been sent
	if wid, ok := co.(fyne.Tappable); ok && doubleTapped == false {
		if action == glfw.Press {
			w.mousePressed = wid
		} else if action == glfw.Release {
			if wid == w.mousePressed {
				switch button {
				case desktop.RightMouseButton:
					w.queueEvent(func() { wid.TappedSecondary(ev) })
				default:
					w.queueEvent(func() { wid.Tapped(ev) })
				}
			}
			w.mousePressed = nil
		}
	}
	if wid, ok := co.(fyne.Draggable); ok {
		if action == glfw.Press {
			w.mouseDragPos = w.mousePos
			w.mouseDragged = wid
			w.mouseDraggedOffset = w.mousePos.Subtract(co.Position()).Subtract(ev.Position)
		}
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
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
	co, _ := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
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
		ev.DeltaX = int(xoff * scrollSpeed)
		ev.DeltaY = int(yoff * scrollSpeed)
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
			button = desktop.RightMouseButton
		} else {
			button = desktop.LeftMouseButton
		}
	case glfw.MouseButton2:
		button = desktop.RightMouseButton
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

	glfw.KeyKPEnter: fyne.KeyEnter,

	// printable
	glfw.KeySpace:      fyne.KeySpace,
	glfw.KeyApostrophe: fyne.KeyApostrophe,
	glfw.KeyComma:      fyne.KeyComma,
	glfw.KeyMinus:      fyne.KeyMinus,
	glfw.KeyPeriod:     fyne.KeyPeriod,
	glfw.KeySlash:      fyne.KeySlash,

	glfw.Key0:         fyne.Key0,
	glfw.Key1:         fyne.Key1,
	glfw.Key2:         fyne.Key2,
	glfw.Key3:         fyne.Key3,
	glfw.Key4:         fyne.Key4,
	glfw.Key5:         fyne.Key5,
	glfw.Key6:         fyne.Key6,
	glfw.Key7:         fyne.Key7,
	glfw.Key8:         fyne.Key8,
	glfw.Key9:         fyne.Key9,
	glfw.KeySemicolon: fyne.KeySemicolon,
	glfw.KeyEqual:     fyne.KeyEqual,

	glfw.KeyA: fyne.KeyA,
	glfw.KeyB: fyne.KeyB,
	glfw.KeyC: fyne.KeyC,
	glfw.KeyD: fyne.KeyD,
	glfw.KeyE: fyne.KeyE,
	glfw.KeyF: fyne.KeyF,
	glfw.KeyG: fyne.KeyG,
	glfw.KeyH: fyne.KeyH,
	glfw.KeyI: fyne.KeyI,
	glfw.KeyJ: fyne.KeyJ,
	glfw.KeyK: fyne.KeyK,
	glfw.KeyL: fyne.KeyL,
	glfw.KeyM: fyne.KeyM,
	glfw.KeyN: fyne.KeyN,
	glfw.KeyO: fyne.KeyO,
	glfw.KeyP: fyne.KeyP,
	glfw.KeyQ: fyne.KeyQ,
	glfw.KeyR: fyne.KeyR,
	glfw.KeyS: fyne.KeyS,
	glfw.KeyT: fyne.KeyT,
	glfw.KeyU: fyne.KeyU,
	glfw.KeyV: fyne.KeyV,
	glfw.KeyW: fyne.KeyW,
	glfw.KeyX: fyne.KeyX,
	glfw.KeyY: fyne.KeyY,
	glfw.KeyZ: fyne.KeyZ,

	glfw.KeyLeftBracket:  fyne.KeyLeftBracket,
	glfw.KeyBackslash:    fyne.KeyBackslash,
	glfw.KeyRightBracket: fyne.KeyRightBracket,

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
}

func keyToName(code glfw.Key) fyne.KeyName {
	ret, ok := keyCodeMap[code]
	if !ok {
		return ""
	}

	return ret
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key)
	if keyName == "" {
		return
	}
	keyEvent := &fyne.KeyEvent{Name: keyName}
	keyDesktopModifier := desktopModifier(mods)

	if keyName == fyne.KeyTab {
		if keyDesktopModifier == 0 {
			if action != glfw.Release {
				w.canvas.focusMgr.FocusNext(w.canvas.focused)
			}
			return
		} else if keyDesktopModifier == desktop.ShiftModifier {
			if action != glfw.Release {
				w.canvas.focusMgr.FocusPrevious(w.canvas.focused)
			}
			return
		}
	}
	if action == glfw.Press {
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				w.queueEvent(func() { focused.KeyDown(keyEvent) })
			}
		} else if w.canvas.onKeyDown != nil {
			w.queueEvent(func() { w.canvas.onKeyDown(keyEvent) })
		}
	} else if action == glfw.Release { // ignore key up in core events
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
		case fyne.KeyC:
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

// charModInput defines the character with modifiers callback which is called when a
// Unicode character is input regardless of what modifier keys are used.
//
// The character with modifiers callback is intended for implementing custom
// Unicode character input. Characters do not map 1:1 to physical keys,
// as a key may produce zero, one or more characters.
func (w *window) charModInput(viewport *glfw.Window, char rune, mods glfw.ModifierKey) {
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

func (w *window) focused(viewport *glfw.Window, focused bool) {
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

	glfw.DetachCurrentContext()
}

func (w *window) RescaleContext() {
	runOnMain(func() {
		w.rescaleOnMain()
	})
}

func (w *window) rescaleOnMain() {
	w.fitContent()
	if w.fullScreen {
		w.width, w.height = w.viewport.GetSize()
		scaledFull := fyne.NewSize(
			internal.UnscaleInt(w.canvas, w.width),
			internal.UnscaleInt(w.canvas, w.height))
		w.canvas.Resize(scaledFull)
		return
	}

	size := w.canvas.size.Union(w.canvas.MinSize())
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
		fyne.LogError("EventQueue full", nil)
	}
}

func (w *window) runEventQueue() {
	for fn := range w.eventQueue {
		fn()
		w.eventWait.Done()
	}
}

func (w *window) waitForEvents() {
	w.eventWait.Wait()
}

func (d *gLDriver) CreateWindow(title string) fyne.Window {
	var ret *window
	if title == "" {
		title = defaultTitle
	}
	runOnMain(func() {
		initOnce.Do(d.initGLFW)

		// make the window hidden, we will set it up and then show it later
		glfw.WindowHint(glfw.Visible, 0)
		initWindowHints()

		win, err := glfw.CreateWindow(10, 10, title, nil, nil)
		if err != nil {
			fyne.LogError("window creation error", err)
			return
		}
		win.MakeContextCurrent()

		ret = &window{viewport: win, title: title}

		// This channel will be closed when the window is closed.
		ret.eventQueue = make(chan func(), 1024)
		go ret.runEventQueue()

		ret.canvas = newCanvas()
		ret.canvas.painter = gl.NewPainter(ret.canvas, ret)
		ret.canvas.painter.Init()
		ret.canvas.context = ret
		ret.canvas.detectedScale = ret.detectScale()
		ret.canvas.scale = ret.calculatedScale()
		ret.SetIcon(ret.icon)
		d.windows = append(d.windows, ret)

		win.SetCloseCallback(ret.closed)
		win.SetPosCallback(ret.moved)
		win.SetSizeCallback(ret.resized)
		win.SetFramebufferSizeCallback(ret.frameSized)
		win.SetRefreshCallback(ret.refresh)
		win.SetCursorPosCallback(ret.mouseMoved)
		win.SetMouseButtonCallback(ret.mouseClicked)
		win.SetScrollCallback(ret.mouseScrolled)
		win.SetKeyCallback(ret.keyPressed)
		win.SetCharModsCallback(ret.charModInput)
		win.SetFocusCallback(ret.focused)
		glfw.DetachCurrentContext()
	})
	return ret
}

func (d *gLDriver) AllWindows() []fyne.Window {
	return d.windows
}
