package gl

import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/widget"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	scrollSpeed      = 10
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

var (
	defaultCursor, entryCursor, hyperlinkCursor *glfw.Cursor
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
	mouseButton        desktop.MouseButton
	mouseOver          desktop.Hoverable
	mouseClickTime     time.Time
	mousePressed       fyne.Tappable
	onClosed           func()

	xpos, ypos    int
	width, height int
	ignoreResize  bool
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

// minSizeOnScreen gets thpaddede minimum size of a window content in screen pixels
func (w *window) minSizeOnScreen() (int, int) {
	// get minimum size of content inside the window
	return w.screenSize(w.canvas.MinSize())
}

// screenSize computes the actual output size of the given content size in screen pixels
func (w *window) screenSize(canvasSize fyne.Size) (int, int) {
	return scaleInt(w.canvas, canvasSize.Width), scaleInt(w.canvas, canvasSize.Height)
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
	scale := w.canvas.Scale()
	w.width, w.height = int(float32(size.Width)*scale), int(float32(size.Height)*scale)
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
	for _, monitor := range glfw.GetMonitors() {
		x, y := monitor.GetPos()

		if x > w.xpos || y > w.ypos {
			continue
		}
		if x+monitor.GetVideoMode().Width <= w.xpos || y+monitor.GetVideoMode().Height <= w.ypos {
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

func (w *window) detectScale() float32 {
	env := os.Getenv("FYNE_SCALE")
	if env == "" {
		return 1.0
	}

	scale, err := strconv.ParseFloat(env, 32)
	if err == nil && scale != 0 {
		return float32(scale)
	}
	if err != nil && env != "auto" { // ignore auto, otherwise report error
		fyne.LogError("Error reading scale", err)
	}

	monitor := w.getMonitorForWindow()
	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	dpi := float32(widthPx) / (float32(widthMm) / 25.4)
	return float32(math.Round(float64(dpi)/144.0*10.0)) / 10.0
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
	if !w.fullScreen {
		w.width = scaleInt(w.canvas, canvasSize.Width)
		w.height = scaleInt(w.canvas, canvasSize.Height)
	}

	w.canvas.Resize(canvasSize)
}

func (w *window) SetContent(content fyne.CanvasObject) {
	// hide old canvas element
	if w.visible && w.canvas.Content() != nil {
		w.canvas.Content().Hide()
	}

	w.canvas.SetContent(content)
	// show top canvas element
	if w.visible {
		w.canvas.Content().Show()
	}

	runOnMain(func() {
		w.fitContent()
		w.resize(w.canvas.size)
	})
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

	driver.WalkCompleteObjectTree(w.canvas.content, nil, func(obj, _ fyne.CanvasObject) {
		switch co := obj.(type) {
		case fyne.Widget:
			widget.DestroyRenderer(co)
		}
	})

	// trigger callbacks
	if w.onClosed != nil {
		w.onClosed()
	}
}

func (w *window) moved(viewport *glfw.Window, x, y int) {
	// save coordinates
	w.xpos, w.ypos = x, y
	scale := w.canvas.scale
	newScale := w.detectScale()

	if scale == newScale {
		forceWindowRefresh(w)
		return
	}

	w.canvas.SetScale(newScale)

	// this can trigger resize events that we need to ignore
	w.fitContent()

	newWidth, newHeight := w.screenSize(w.canvas.size)
	w.viewport.SetSize(newWidth, newHeight)
}

func (w *window) resized(viewport *glfw.Window, width, height int) {
	if w.ignoreResize {
		return
	}
	w.resize(fyne.NewSize(unscaleInt(w.canvas, width), unscaleInt(w.canvas, height)))
}

func (w *window) frameSized(viewport *glfw.Window, width, height int) {
	if width == 0 || height == 0 {
		return
	}

	winWidth, _ := w.viewport.GetSize()
	w.canvas.texScale = float32(width) / float32(winWidth) // This will be > 1.0 on a HiDPI screen
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (w *window) refresh(viewport *glfw.Window) {
	forceWindowRefresh(w)
	w.canvas.setDirty(true)
}

func (w *window) findObjectAtPositionMatching(canvas *glCanvas, mouse fyne.Position,
	matches func(object fyne.CanvasObject) bool) (fyne.CanvasObject, int, int) {
	var found fyne.CanvasObject
	foundX, foundY := 0, 0

	findFunc := func(walked fyne.CanvasObject, pos fyne.Position, clipPos fyne.Position, clipSize fyne.Size) bool {
		if !walked.Visible() {
			return false
		}

		if mouse.X < clipPos.X || mouse.Y < clipPos.Y {
			return false
		}

		if mouse.X >= clipPos.X+clipSize.Width || mouse.Y >= clipPos.Y+clipSize.Height {
			return false
		}

		if mouse.X < pos.X || mouse.Y < pos.Y {
			return false
		}

		if mouse.X >= pos.X+walked.Size().Width || mouse.Y >= pos.Y+walked.Size().Height {
			return false
		}

		if matches(walked) {
			found = walked
			foundX, foundY = mouse.X-pos.X, mouse.Y-pos.Y
		}
		return false
	}

	if canvas.overlay != nil {
		driver.WalkVisibleObjectTree(canvas.overlay, findFunc, nil)
	} else {
		if canvas.menu != nil {
			driver.WalkVisibleObjectTree(canvas.menu, findFunc, nil)
		}
		if found == nil {
			driver.WalkVisibleObjectTree(canvas.content, findFunc, nil)
		}
	}

	return found, foundX, foundY
}

func (w *window) mouseMoved(viewport *glfw.Window, xpos float64, ypos float64) {
	w.mousePos = fyne.NewPos(unscaleInt(w.canvas, int(xpos)), unscaleInt(w.canvas, int(ypos)))

	cursor := defaultCursor
	obj, x, y := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
		if wid, ok := object.(*widget.Entry); ok {
			if !wid.ReadOnly {
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
		ev.Position = fyne.NewPos(x, y)
		ev.Button = w.mouseButton

		if hovered, ok := obj.(desktop.Hoverable); ok {
			if hovered == w.mouseOver {
				hovered.MouseMoved(ev)
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
			w.mouseDragged.Dragged(ev)

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

func (w *window) mouseClicked(viewport *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	co, x, y := w.findObjectAtPositionMatching(w.canvas, w.mousePos, func(object fyne.CanvasObject) bool {
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
	ev.Position = fyne.NewPos(x, y)

	if wid, ok := co.(desktop.Mouseable); ok {
		mev := new(desktop.MouseEvent)
		mev.Position = ev.Position
		mev.Button = convertMouseButton(button)
		if action == glfw.Press {
			go wid.MouseDown(mev)
		} else if action == glfw.Release {
			go wid.MouseUp(mev)
		}
	}

	needsfocus := true
	wid := w.canvas.Focused()
	if wid != nil {
		needsfocus = false
		if wid.(fyne.CanvasObject) != co {
			w.canvas.Unfocus()
			needsfocus = true
		}
	}

	if action == glfw.Press {
		w.mouseButton = convertMouseButton(button)
	} else if action == glfw.Release {
		w.mouseButton = 0
	}

	// we cannot switch here as objects may respond to multiple cases
	if wid, ok := co.(fyne.Focusable); ok {
		if needsfocus == true {
			w.canvas.Focus(wid)
		}
	}

	// Check for double click/tap
	doubleTapped := false
	if action == glfw.Release && button == glfw.MouseButtonLeft {
		now := time.Now()
		// we can safely subtract the first "zero" time as it'll be much larger than doubleClickDelay
		if now.Sub(w.mouseClickTime).Nanoseconds()/1e6 <= doubleClickDelay {
			if wid, ok := co.(fyne.DoubleTappable); ok {
				doubleTapped = true
				go wid.DoubleTapped(ev)
			}
		}
		w.mouseClickTime = now
	}

	// Prevent Tapped from triggering if DoubleTapped has been sent
	if wid, ok := co.(fyne.Tappable); ok && doubleTapped == false {
		if action == glfw.Press {
			w.mousePressed = wid
		} else if action == glfw.Release {
			if wid == w.mousePressed {
				switch button {
				case glfw.MouseButtonRight:
					go wid.TappedSecondary(ev)
				default:
					go wid.Tapped(ev)
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
		w.mouseDragged.DragEnd()
		if w.objIsDragged(w.mouseOver) && !w.objIsDragged(co) {
			w.mouseOut()
		}
		w.mouseDragged = nil
	}
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
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

func convertMouseButton(button glfw.MouseButton) desktop.MouseButton {
	switch button {
	case glfw.MouseButton1:
		return desktop.LeftMouseButton
	case glfw.MouseButton2:
		return desktop.RightMouseButton
	default:
		return 0
	}
}

func keyToName(key glfw.Key) fyne.KeyName {
	switch key {
	// non-printable
	case glfw.KeyEscape:
		return fyne.KeyEscape
	case glfw.KeyEnter:
		return fyne.KeyReturn
	case glfw.KeyTab:
		return fyne.KeyTab
	case glfw.KeyBackspace:
		return fyne.KeyBackspace
	case glfw.KeyInsert:
		return fyne.KeyInsert
	case glfw.KeyDelete:
		return fyne.KeyDelete
	case glfw.KeyRight:
		return fyne.KeyRight
	case glfw.KeyLeft:
		return fyne.KeyLeft
	case glfw.KeyDown:
		return fyne.KeyDown
	case glfw.KeyUp:
		return fyne.KeyUp
	case glfw.KeyPageUp:
		return fyne.KeyPageUp
	case glfw.KeyPageDown:
		return fyne.KeyPageDown
	case glfw.KeyHome:
		return fyne.KeyHome
	case glfw.KeyEnd:
		return fyne.KeyEnd

	case glfw.KeyF1:
		return fyne.KeyF1
	case glfw.KeyF2:
		return fyne.KeyF2
	case glfw.KeyF3:
		return fyne.KeyF3
	case glfw.KeyF4:
		return fyne.KeyF4
	case glfw.KeyF5:
		return fyne.KeyF5
	case glfw.KeyF6:
		return fyne.KeyF6
	case glfw.KeyF7:
		return fyne.KeyF7
	case glfw.KeyF8:
		return fyne.KeyF8
	case glfw.KeyF9:
		return fyne.KeyF9
	case glfw.KeyF10:
		return fyne.KeyF10
	case glfw.KeyF11:
		return fyne.KeyF11
	case glfw.KeyF12:
		return fyne.KeyF12

	case glfw.KeyKPEnter:
		return fyne.KeyEnter

	// printable
	case glfw.KeyA:
		return fyne.KeyA
	case glfw.KeyB:
		return fyne.KeyB
	case glfw.KeyC:
		return fyne.KeyC
	case glfw.KeyD:
		return fyne.KeyD
	case glfw.KeyE:
		return fyne.KeyE
	case glfw.KeyF:
		return fyne.KeyF
	case glfw.KeyG:
		return fyne.KeyG
	case glfw.KeyH:
		return fyne.KeyH
	case glfw.KeyI:
		return fyne.KeyI
	case glfw.KeyJ:
		return fyne.KeyJ
	case glfw.KeyK:
		return fyne.KeyK
	case glfw.KeyL:
		return fyne.KeyL
	case glfw.KeyM:
		return fyne.KeyM
	case glfw.KeyN:
		return fyne.KeyN
	case glfw.KeyO:
		return fyne.KeyO
	case glfw.KeyP:
		return fyne.KeyP
	case glfw.KeyQ:
		return fyne.KeyQ
	case glfw.KeyR:
		return fyne.KeyR
	case glfw.KeyS:
		return fyne.KeyS
	case glfw.KeyT:
		return fyne.KeyT
	case glfw.KeyU:
		return fyne.KeyU
	case glfw.KeyV:
		return fyne.KeyV
	case glfw.KeyW:
		return fyne.KeyW
	case glfw.KeyX:
		return fyne.KeyX
	case glfw.KeyY:
		return fyne.KeyY
	case glfw.KeyZ:
		return fyne.KeyZ
	case glfw.Key0:
		return fyne.Key0
	case glfw.Key1:
		return fyne.Key1
	case glfw.Key2:
		return fyne.Key2
	case glfw.Key3:
		return fyne.Key3
	case glfw.Key4:
		return fyne.Key4
	case glfw.Key5:
		return fyne.Key5
	case glfw.Key6:
		return fyne.Key6
	case glfw.Key7:
		return fyne.Key7
	case glfw.Key8:
		return fyne.Key8
	case glfw.Key9:
		return fyne.Key9

	// desktop
	case glfw.KeyLeftShift:
		return desktop.KeyShiftLeft
	case glfw.KeyRightShift:
		return desktop.KeyShiftRight
	case glfw.KeyLeftControl:
		return desktop.KeyControlLeft
	case glfw.KeyRightControl:
		return desktop.KeyControlRight
	case glfw.KeyLeftAlt:
		return desktop.KeyAltLeft
	case glfw.KeyRightAlt:
		return desktop.KeyAltRight
	case glfw.KeyLeftSuper:
		return desktop.KeySuperLeft
	case glfw.KeyRightSuper:
		return desktop.KeySuperRight
	case glfw.KeyMenu:
		return desktop.KeyMenu
	}
	return ""
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
				go focused.KeyDown(keyEvent)
			}
		} else if w.canvas.onKeyDown != nil {
			go w.canvas.onKeyDown(keyEvent)
		}
	} else if action == glfw.Release { // ignore key up in core events
		if w.canvas.Focused() != nil {
			if focused, ok := w.canvas.Focused().(desktop.Keyable); ok {
				go focused.KeyUp(keyEvent)
			}
		} else if w.canvas.onKeyUp != nil {
			go w.canvas.onKeyUp(keyEvent)
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
		}
	}
	if shortcut == nil && keyDesktopModifier != 0 {
		shortcut = &desktop.CustomShortcut{
			KeyName:  keyName,
			Modifier: keyDesktopModifier,
		}
	}

	if shortcutable, ok := w.canvas.Focused().(fyne.Shortcutable); ok {
		if shortcutable.TypedShortcut(shortcut) {
			return
		}
	} else if w.canvas.shortcut.TypedShortcut(shortcut) {
		return
	}

	// No shortcut detected, pass down to TypedKey
	if w.canvas.Focused() != nil {
		go w.canvas.Focused().TypedKey(keyEvent)
	} else if w.canvas.onTypedKey != nil {
		go w.canvas.onTypedKey(keyEvent)
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

	if w.canvas.Focused() != nil {
		w.canvas.Focused().TypedRune(char)
	} else if w.canvas.onTypedRune != nil {
		w.canvas.onTypedRune(char)
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

func (w *window) runWithContext(f func()) {
	w.viewport.MakeContextCurrent()

	f()

	glfw.DetachCurrentContext()
}

func (d *gLDriver) CreateWindow(title string) fyne.Window {
	var ret *window
	runOnMain(func() {
		master := len(d.windows) == 0
		if master {
			err := glfw.Init()
			if err != nil {
				fyne.LogError("failed to initialise GLFW", err)
				return
			}

			initCursors()
		}

		// make the window hidden, we will set it up and then show it later
		glfw.WindowHint(glfw.Visible, 0)

		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)

		win, err := glfw.CreateWindow(10, 10, title, nil, nil)
		if err != nil {
			fyne.LogError("window creation error", err)
			return
		}
		win.MakeContextCurrent()

		if master {
			err := gl.Init()
			if err != nil {
				fyne.LogError("failed to initialise OpenGL", err)
				return
			}

			gl.Disable(gl.DEPTH_TEST)
		}
		ret = &window{viewport: win, title: title, master: master}
		canvas := newCanvas()
		canvas.context = ret
		canvas.SetScale(ret.detectScale())
		canvas.texScale = 1.0
		ret.canvas = canvas
		ret.SetIcon(ret.icon) // if this is nil we will get the app icon
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
