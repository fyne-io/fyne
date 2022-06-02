//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package glfw

import (
	"context"
	_ "image/png" // for the icon
	"runtime"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter/gl"

	"github.com/fyne-io/glfw-js"
)

type Cursor struct {
}

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
	cursorMap    map[desktop.Cursor]*Cursor
	defaultTitle = "Fyne Application"
)

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

	cursor   desktop.Cursor
	canvas   *glCanvas
	driver   *gLDriver
	title    string
	icon     fyne.Resource
	mainmenu *fyne.MainMenu

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
	w.fullScreen = true
}

// centerOnScreen handles the logic for centering a window
func (w *window) CenterOnScreen() {
	// FIXME: not supported with WebGL
	w.centered = true
}

func (w *window) doCenterOnScreen() {
	// FIXME: no meaning for defining center on screen in WebGL
}

func (w *window) RequestFocus() {
	// FIXME: no meaning for defining focus in WebGL
}

func (w *window) SetIcon(icon fyne.Resource) {
	// FIXME: no support for SetIcon yet
}

func (w *window) SetMaster() {
	// FIXME: there could really only be one window
}

func (w *window) fitContent() {
	w.shouldWidth, w.shouldHeight = w.requestedWidth, w.requestedHeight
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	return glfw.GetPrimaryMonitor()
}

func scaleForDpi(xdpi int) float32 {
	switch {
	case xdpi > 1000:
		// assume that this is a mistake and bail
		return float32(1.0)
	case xdpi > 192:
		return float32(1.5)
	case xdpi > 144:
		return float32(1.35)
	case xdpi > 120:
		return float32(1.2)
	default:
		return float32(1.0)
	}
}

func (w *window) detectScale() float32 {
	return scaleForDpi(int(96))
}

func (w *window) moved(_ *glfw.Window, x, y int) {
	w.processMoved(x, y)
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	w.canvas.scale = w.calculatedScale()
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

func fyneToNativeCursor(cursor desktop.Cursor) (*Cursor, bool) {
	return nil, false
}

func (w *window) SetCursor(_ *Cursor) {
}

func (w *window) setCustomCursor(rawCursor *Cursor, isCustomCursor bool) {
}

func (w *window) mouseMoved(_ *glfw.Window, xpos, ypos float64) {
	w.processMouseMoved(xpos, ypos)
}

func (w *window) mouseClicked(viewport *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	button, modifiers := convertMouseButton(btn, mods)
	mouseAction := convertAction(action)

	w.processMouseClicked(button, mouseAction, modifiers)
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff, yoff float64) {
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

	//	keyName := glfw.GetKeyName(code, scancode)
	//	ret, ok = keyNameMap[keyName]
	//	if !ok {
	return fyne.KeyUnknown
	//	}

	//	return ret
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

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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
// Unicode character is input regardless of what modifier keys are used.
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
	if w.viewport == nil {
		return
	}

	//	if w.fullScreen {
	w.width, w.height = w.viewport.GetSize()
	scaledFull := fyne.NewSize(
		internal.UnscaleInt(w.canvas, w.width),
		internal.UnscaleInt(w.canvas, w.height))
	w.canvas.Resize(scaledFull)
	return
	//	}

	//	size := w.canvas.size.Union(w.canvas.MinSize())
	//	newWidth, newHeight := w.screenSize(size)
	//	w.viewport.SetSize(newWidth, newHeight)
}

func (w *window) create() {
	runOnMain(func() {
		// we can't hide the window in webgl, so there might be some artifact
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

		w.requestedWidth, w.requestedHeight = w.width, w.height

		width, height := win.GetSize()
		w.processFrameSized(width, height)
		w.processResized(width, height)
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
