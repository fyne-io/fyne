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
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
)

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
	newX := (monMode.Width / 2) - (viewWidth / 2) + monX
	newY := (monMode.Height / 2) - (viewHeight / 2) + monY

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
		if painter.IsResourceSVG(w.icon) {
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

func (w *window) detectScale() float32 {
	monitor := w.getMonitorForWindow()
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

func (w *window) frameSized(_ *glfw.Window, width, height int) {
	w.processFrameSized(width, height)
}

func (w *window) refresh(_ *glfw.Window) {
	w.processRefresh()
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(false) // reset the closed flag until we check the veto in processClosed

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

var keyCodeMapASCII = map[glfw.Key]fyne.KeyName{
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
		return fyne.KeyUnknown
	}

	return ret
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
	ret, ok := keyCodeMapASCII[key]
	if !ok {
		return fyne.KeyUnknown
	}
	return ret
}

func (w *window) keyPressed(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	keyDesktopModifier := desktopModifier(mods)
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
			internal.UnscaleInt(w.canvas, w.width),
			internal.UnscaleInt(w.canvas, w.height))
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
			w.width, w.height = internal.ScaleInt(w.canvas, bigEnough.Width), internal.ScaleInt(w.canvas, bigEnough.Height)
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
