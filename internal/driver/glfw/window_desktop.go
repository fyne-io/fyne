//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"bytes"
	"context"
	"image"
	_ "image/png" // for the icon
	"runtime"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/storage"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const defaultTitle = "Fyne Application"

// Input modes.
const (
	CursorMode             glfw.InputMode = glfw.CursorMode
	StickyKeysMode         glfw.InputMode = glfw.StickyKeysMode
	StickyMouseButtonsMode glfw.InputMode = glfw.StickyMouseButtonsMode
	LockKeyMods            glfw.InputMode = glfw.LockKeyMods
	RawMouseMotion         glfw.InputMode = glfw.RawMouseMotion
)

// Cursor mode values.
const (
	CursorNormal   int = glfw.CursorNormal
	CursorHidden   int = glfw.CursorHidden
	CursorDisabled int = glfw.CursorDisabled
)

var cursorMap map[desktop.StandardCursor]*glfw.Cursor

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
	common.Window

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

	mouseLock            sync.RWMutex
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

	onClosed           func()
	onCloseIntercepted func()

	menuTogglePending       fyne.KeyName
	menuDeactivationPending fyne.KeyName

	xpos, ypos                      int
	width, height                   int
	requestedWidth, requestedHeight int
	shouldWidth, shouldHeight       int
	shouldExpand                    bool

	pending []func()
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
			if w.width == 0 && w.height == 0 { // if we were fullscreen on creation...
				w.width, w.height = w.screenSize(w.canvas.Size())
			}
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

func (w *window) SetOnDropped(dropped func(pos fyne.Position, items []fyne.URI)) {
	w.runOnMainWhenCreated(func() {
		w.viewport.SetDropCallback(func(win *glfw.Window, names []string) {
			if dropped == nil {
				return
			}

			uris := make([]fyne.URI, len(names))
			for i, name := range names {
				uris[i] = storage.NewFileURI(name)
			}

			dropped(w.mousePos, uris)
		})
	})
}

func (w *window) doCenterOnScreen() {
	viewWidth, viewHeight := w.screenSize(w.canvas.size)
	if w.width > viewWidth { // in case our window has not called back to canvas size yet
		viewWidth = w.width
	}
	if w.height > viewHeight {
		viewHeight = w.height
	}

	// get window dimensions in pixels
	monitor := w.getMonitorForWindow()
	monMode := monitor.GetVideoMode()

	// these come into play when dealing with multiple monitors
	monX, monY := monitor.GetPos()

	// math them to the middle
	newX := (monMode.Width-viewWidth)/2 + monX
	newY := (monMode.Height-viewHeight)/2 + monY

	// set new window coordinates
	w.viewport.SetPos(newX, newY)
}

func (w *window) RequestFocus() {
	if isWayland || w.view() == nil {
		return
	}

	w.runOnMainWhenCreated(w.viewport.Focus)
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

	w.runOnMainWhenCreated(func() {
		if w.icon == nil {
			w.viewport.SetIcon(nil)
			return
		}

		var img image.Image
		if svg.IsResourceSVG(w.icon) {
			img = painter.PaintImage(&canvas.Image{Resource: w.icon}, nil, windowIconSize, windowIconSize)
		} else {
			pix, _, err := image.Decode(bytes.NewReader(w.icon.Content()))
			if err != nil {
				fyne.LogError("Failed to decode image for window icon", err)
				return
			}
			img = pix
		}

		w.viewport.SetIcon([]image.Image{img})
	})
}

func (w *window) SetMaster() {
	w.master = true
}

func (w *window) fitContent() {
	if w.canvas.Content() == nil || (w.fullScreen && w.visible) {
		return
	}

	if w.isClosing() {
		return
	}

	minWidth, minHeight := w.minSizeOnScreen()
	w.viewLock.RLock()
	view := w.viewport
	w.viewLock.RUnlock()
	w.shouldWidth, w.shouldHeight = w.width, w.height
	if w.width < minWidth || w.height < minHeight {
		if w.width < minWidth {
			w.shouldWidth = minWidth
		}
		if w.height < minHeight {
			w.shouldHeight = minHeight
		}
		w.viewLock.Lock()
		w.shouldExpand = true // queue the resize to happen on main
		w.viewLock.Unlock()
	}
	if w.fixedSize {
		if w.shouldWidth > w.requestedWidth {
			w.requestedWidth = w.shouldWidth
		}
		if w.shouldHeight > w.requestedHeight {
			w.requestedHeight = w.shouldHeight
		}
		view.SetSizeLimits(w.requestedWidth, w.requestedHeight, w.requestedWidth, w.requestedHeight)
	} else {
		view.SetSizeLimits(minWidth, minHeight, glfw.DontCare, glfw.DontCare)
	}
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
		if videoMode := monitor.GetVideoMode(); x+videoMode.Width <= xOff || y+videoMode.Height <= yOff {
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
	if isWayland { // Wayland controls scale through content scaling
		return 1.0
	}
	monitor := w.getMonitorForWindow()
	if monitor == nil {
		return 1.0
	}

	widthMm, _ := monitor.GetPhysicalSize()
	widthPx := monitor.GetVideoMode().Width

	return calculateDetectedScale(widthMm, widthPx)
}

func (w *window) moved(_ *glfw.Window, x, y int) {
	w.processMoved(x, y)
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	w.processResized(width, height)
}

func (w *window) scaled(_ *glfw.Window, x float32, y float32) {
	if !isWayland { // other platforms handle this using older APIs
		return
	}

	w.canvas.texScale = x
	w.canvas.Refresh(w.canvas.content)
}

func (w *window) frameSized(_ *glfw.Window, width, height int) {
	w.processFrameSized(width, height)
}

func (w *window) refresh(_ *glfw.Window) {
	w.processRefresh()
}

func (w *window) closed(viewport *glfw.Window) {
	if viewport != nil {
		viewport.SetShouldClose(false) // reset the closed flag until we check the veto in processClosed
	}

	w.processClosed()
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

func (w *window) SetCursor(cursor *glfw.Cursor) {
	w.viewport.SetCursor(cursor)
}

func (w *window) setCustomCursor(rawCursor *glfw.Cursor, isCustomCursor bool) {
	if w.customCursor != nil {
		w.customCursor.Destroy()
		w.customCursor = nil
	}
	if isCustomCursor {
		w.customCursor = rawCursor
	}

}

func (w *window) mouseMoved(_ *glfw.Window, xpos, ypos float64) {
	w.processMouseMoved(xpos, ypos)
}

func (w *window) mouseClicked(_ *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	button, modifiers := convertMouseButton(btn, mods)
	mouseAction := convertAction(action)

	w.processMouseClicked(button, mouseAction, modifiers)
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff float64, yoff float64) {
	if runtime.GOOS != "darwin" && xoff == 0 &&
		(viewport.GetKey(glfw.KeyLeftShift) == glfw.Press ||
			viewport.GetKey(glfw.KeyRightShift) == glfw.Press) {
		xoff, yoff = yoff, xoff
	}

	w.processMouseScrolled(xoff, yoff)
}

func convertMouseButton(btn glfw.MouseButton, mods glfw.ModifierKey) (desktop.MouseButton, fyne.KeyModifier) {
	modifier := desktopModifier(mods)
	var button desktop.MouseButton
	rightClick := false
	if runtime.GOOS == "darwin" {
		if modifier&fyne.KeyModifierControl != 0 {
			rightClick = true
			modifier &^= fyne.KeyModifierControl
		}
		if modifier&fyne.KeyModifierSuper != 0 {
			modifier |= fyne.KeyModifierControl
			modifier &^= fyne.KeyModifierSuper
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

//gocyclo:ignore
func glfwKeyToKeyName(key glfw.Key) fyne.KeyName {
	switch key {
	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	case glfw.Key0, glfw.KeyKP0:
		return fyne.Key0
	case glfw.Key1, glfw.KeyKP1:
		return fyne.Key1
	case glfw.Key2, glfw.KeyKP2:
		return fyne.Key2
	case glfw.Key3, glfw.KeyKP3:
		return fyne.Key3
	case glfw.Key4, glfw.KeyKP4:
		return fyne.Key4
	case glfw.Key5, glfw.KeyKP5:
		return fyne.Key5
	case glfw.Key6, glfw.KeyKP6:
		return fyne.Key6
	case glfw.Key7, glfw.KeyKP7:
		return fyne.Key7
	case glfw.Key8, glfw.KeyKP8:
		return fyne.Key8
	case glfw.Key9, glfw.KeyKP9:
		return fyne.Key9

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

	case glfw.KeySpace:
		return fyne.KeySpace
	case glfw.KeyKPEnter:
		return fyne.KeyEnter

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
	case glfw.KeyPrintScreen:
		return desktop.KeyPrintScreen
	case glfw.KeyCapsLock:
		return desktop.KeyCapsLock

	// functions
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
	}

	return fyne.KeyUnknown
}

func keyCodeToKeyName(code string) fyne.KeyName {
	if len(code) != 1 {
		return fyne.KeyUnknown
	}

	char := code[0]
	if char >= 'a' && char <= 'z' {
		// Our alphabetical keys are all upper case characters.
		return fyne.KeyName('A' + char - 'a')
	}

	switch char {
	case '[':
		return fyne.KeyLeftBracket
	case '\\':
		return fyne.KeyBackslash
	case ']':
		return fyne.KeyRightBracket
	case '\'':
		return fyne.KeyApostrophe
	case ',':
		return fyne.KeyComma
	case '-':
		return fyne.KeyMinus
	case '.':
		return fyne.KeyPeriod
	case '/':
		return fyne.KeySlash
	case '*':
		return fyne.KeyAsterisk
	case '`':
		return fyne.KeyBackTick
	case ';':
		return fyne.KeySemicolon
	case '+':
		return fyne.KeyPlus
	case '=':
		return fyne.KeyEqual
	}

	return fyne.KeyUnknown
}

func keyToName(code glfw.Key, scancode int) fyne.KeyName {
	if runtime.GOOS == "darwin" && scancode == 0x69 { // TODO remove once fixed upstream glfw/glfw#1786
		code = glfw.KeyPrintScreen
	}

	ret := glfwKeyToKeyName(code)
	if ret != fyne.KeyUnknown {
		return ret
	}

	keyName := glfw.GetKeyName(code, scancode)
	return keyCodeToKeyName(keyName)
}

func convertAction(action glfw.Action) action {
	switch action {
	case glfw.Press:
		return press
	case glfw.Release:
		return release
	case glfw.Repeat:
		return repeat
	}
	panic("Could not convert glfw.Action.")
}

func convertASCII(key glfw.Key) fyne.KeyName {
	if key < glfw.KeyA || key > glfw.KeyZ {
		return fyne.KeyUnknown
	}

	return fyne.KeyName(rune(key))
}

func (w *window) keyPressed(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	keyDesktopModifier := desktopModifier(mods)
	w.driver.currentKeyModifiers = desktopModifierCorrected(mods, key, action)
	keyAction := convertAction(action)
	keyASCII := convertASCII(key)

	w.processKeyPressed(keyName, keyASCII, scancode, keyAction, keyDesktopModifier)
}

func desktopModifier(mods glfw.ModifierKey) fyne.KeyModifier {
	var m fyne.KeyModifier
	if (mods & glfw.ModShift) != 0 {
		m |= fyne.KeyModifierShift
	}
	if (mods & glfw.ModControl) != 0 {
		m |= fyne.KeyModifierControl
	}
	if (mods & glfw.ModAlt) != 0 {
		m |= fyne.KeyModifierAlt
	}
	if (mods & glfw.ModSuper) != 0 {
		m |= fyne.KeyModifierSuper
	}
	return m
}

func desktopModifierCorrected(mods glfw.ModifierKey, key glfw.Key, action glfw.Action) fyne.KeyModifier {
	// On X11, pressing/releasing modifier keys does not include newly pressed/released keys in 'mod' mask.
	// https://github.com/glfw/glfw/issues/1630
	if action == glfw.Press {
		mods |= glfwKeyToModifier(key)
	} else {
		mods &= ^glfwKeyToModifier(key)
	}
	return desktopModifier(mods)
}

func glfwKeyToModifier(key glfw.Key) glfw.ModifierKey {
	var m glfw.ModifierKey
	switch key {
	case glfw.KeyLeftControl, glfw.KeyRightControl:
		m = glfw.ModControl
	case glfw.KeyLeftAlt, glfw.KeyRightAlt:
		m = glfw.ModAlt
	case glfw.KeyLeftShift, glfw.KeyRightShift:
		m = glfw.ModShift
	case glfw.KeyLeftSuper, glfw.KeyRightSuper:
		m = glfw.ModSuper
	}
	return m
}

// charInput defines the character with modifiers callback which is called when a
// Unicode character is input.
//
// Characters do not map 1:1 to physical keys, as a key may produce zero, one or more characters.
func (w *window) charInput(viewport *glfw.Window, char rune) {
	w.processCharInput(char)
}

func (w *window) focused(_ *glfw.Window, focused bool) {
	w.processFocused(focused)
}

func (w *window) DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func (w *window) rescaleOnMain() {
	if w.isClosing() {
		return
	}
	w.fitContent()

	if w.fullScreen {
		w.width, w.height = w.viewport.GetSize()
		scaledFull := fyne.NewSize(
			scale.ToFyneCoordinate(w.canvas, w.width),
			scale.ToFyneCoordinate(w.canvas, w.height))
		w.canvas.Resize(scaledFull)
		return
	}

	size := w.canvas.size.Max(w.canvas.MinSize())
	newWidth, newHeight := w.screenSize(size)
	w.viewport.SetSize(newWidth, newHeight)
}

func (w *window) create() {
	runOnMain(func() {
		if !isWayland {
			// make the window hidden, we will set it up and then show it later
			glfw.WindowHint(glfw.Visible, glfw.False)
		}
		if w.decorate {
			glfw.WindowHint(glfw.Decorated, glfw.True)
		} else {
			glfw.WindowHint(glfw.Decorated, glfw.False)
		}
		if w.fixedSize {
			glfw.WindowHint(glfw.Resizable, glfw.False)
		} else {
			glfw.WindowHint(glfw.Resizable, glfw.True)
		}
		glfw.WindowHint(glfw.AutoIconify, glfw.False)
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
		w.canvas.SetPainter(gl.NewPainter(w.canvas, w))
		w.canvas.Painter().Init()
	})

	runOnMain(func() {
		w.setDarkMode()

		win := w.view()
		win.SetCloseCallback(w.closed)
		win.SetPosCallback(w.moved)
		win.SetSizeCallback(w.resized)
		win.SetFramebufferSizeCallback(w.frameSized)
		win.SetRefreshCallback(w.refresh)
		win.SetContentScaleCallback(w.scaled)
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

		if w.FixedSize() && (w.requestedWidth == 0 || w.requestedHeight == 0) {
			bigEnough := w.canvas.canvasSize(w.canvas.Content().MinSize())
			w.width, w.height = scale.ToScreenCoordinate(w.canvas, bigEnough.Width), scale.ToScreenCoordinate(w.canvas, bigEnough.Height)
			w.shouldWidth, w.shouldHeight = w.width, w.height
		}

		w.requestedWidth, w.requestedHeight = w.width, w.height
		// order of operation matters so we do these last items in order
		w.viewport.SetSize(w.shouldWidth, w.shouldHeight) // ensure we requested latest size
	})
}

func (w *window) view() *glfw.Window {
	w.viewLock.RLock()
	defer w.viewLock.RUnlock()

	if w.closing {
		return nil
	}
	return w.viewport
}
