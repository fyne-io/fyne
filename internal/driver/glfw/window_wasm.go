//go:build wasm || test_web_driver

package glfw

import (
	"context"
	_ "image/png" // for the icon
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/internal/scale"

	"github.com/fyne-io/glfw-js"
)

type Cursor struct {
	JSName string
}

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

// Declare conformity to Window interface
var _ fyne.Window = (*window)(nil)

type window struct {
	viewport  *glfw.Window
	created   bool
	decorate  bool
	closing   bool
	fixedSize bool

	cursor   desktop.Cursor
	canvas   *glCanvas
	driver   *gLDriver
	title    string
	icon     fyne.Resource
	mainmenu *fyne.MainMenu

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

	lastWalkedTime time.Time
}

func (w *window) SetFullScreen(full bool) {
	w.fullScreen = true
}

// centerOnScreen handles the logic for centering a window
func (w *window) CenterOnScreen() {
	// FIXME: not supported with WebGL
	w.centered = true
}

func (w *window) SetOnDropped(dropped func(pos fyne.Position, items []fyne.URI)) {
	// FIXME: not implemented yet
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
	runOnMain(func() {
		w.processMoved(x, y)
	})
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	runOnMain(func() {
		w.canvas.scale = w.calculatedScale()
		w.processResized(width, height)
	})
}

func (w *window) frameSized(_ *glfw.Window, width, height int) {
	runOnMain(func() {
		w.processFrameSized(width, height)
	})
}

func (w *window) refresh(_ *glfw.Window) {
	runOnMain(w.processRefresh)
}

func (w *window) closed(viewport *glfw.Window) {
	runOnMain(func() {
		viewport.SetShouldClose(false) // reset the closed flag until we check the veto in processClosed

		w.processClosed()
	})
}

func fyneToNativeCursor(cursor desktop.Cursor) (*Cursor, bool) {
	if _, ok := cursor.(desktop.StandardCursor); !ok {
		return nil, false // Custom cursors not implemented yet.
	}

	name := "default"
	switch cursor {
	case desktop.TextCursor:
		name = "text"
	case desktop.CrosshairCursor:
		name = "crosshair"
	case desktop.DefaultCursor:
		name = "default"
	case desktop.PointerCursor:
		name = "pointer"
	case desktop.HResizeCursor:
		name = "ew-resize"
	case desktop.VResizeCursor:
		name = "ns-resize"
	case desktop.HiddenCursor:
		name = "none"
	}

	return &Cursor{JSName: name}, false
}

func (w *window) SetCursor(cursor *Cursor) {
	setCursor(cursor.JSName)
}

func (w *window) setCustomCursor(rawCursor *Cursor, isCustomCursor bool) {
}

func (w *window) mouseMoved(_ *glfw.Window, xpos, ypos float64) {
	runOnMain(func() {
		w.processMouseMoved(w.scaleInput(xpos), w.scaleInput(ypos))
	})
}

func (w *window) mouseClicked(viewport *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	runOnMain(func() {
		button, modifiers := convertMouseButton(btn, mods)
		mouseAction := convertAction(action)

		w.processMouseClicked(button, mouseAction, modifiers)
	})
}

func (w *window) mouseScrolled(viewport *glfw.Window, xoff, yoff float64) {
	runOnMain(func() {
		if xoff == 0 &&
			(viewport.GetKey(glfw.KeyLeftShift) == glfw.Press ||
				viewport.GetKey(glfw.KeyRightShift) == glfw.Press) {
			xoff, yoff = yoff, xoff
		}

		w.processMouseScrolled(xoff, yoff)
	})
}

func convertMouseButton(btn glfw.MouseButton, mods glfw.ModifierKey) (desktop.MouseButton, fyne.KeyModifier) {
	modifier := desktopModifier(mods)
	rightClick := false
	if isMacOSRuntime() {
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
			return desktop.MouseButtonSecondary, modifier
		}
		return desktop.MouseButtonPrimary, modifier
	case glfw.MouseButton2:
		return desktop.MouseButtonSecondary, modifier
	case glfw.MouseButton3:
		return desktop.MouseButtonTertiary, modifier
	default:
		return 0, modifier
	}
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
	ret := glfwKeyToKeyName(code)
	if ret != fyne.KeyUnknown {
		return ret
	}

	//	keyName := glfw.GetKeyName(code, scancode)
	//	return keyCodeToKeyName(keyName)
	return fyne.KeyUnknown
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

func (w *window) RescaleContext() {
	if w.viewport == nil {
		return
	}

	w.width, w.height = w.viewport.GetSize()
	scaledFull := fyne.NewSize(
		scale.ToFyneCoordinate(w.canvas, w.width),
		scale.ToFyneCoordinate(w.canvas, w.height))
	w.canvas.Resize(scaledFull)

	// Ensure textures re-rasterize at the new scale
	cache.DeleteTextTexturesFor(w.canvas)
	w.canvas.content.Refresh()
}

func (w *window) create() {
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

	w.viewport = win

	if w.view() == nil { // something went wrong above, it will have been logged
		return
	}

	// run the GL init on the draw thread
	w.RunWithContext(func() {
		w.canvas.SetPainter(gl.NewPainter(w.canvas, w))
		w.canvas.Painter().Init()
	})

	w.setDarkMode()

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

	w.drainPendingEvents()

	w.requestedWidth, w.requestedHeight = w.width, w.height

	width, height := win.GetSize()
	w.processFrameSized(width, height)
	w.processResized(width, height)
}

func (w *window) view() *glfw.Window {
	if w.closing {
		return nil
	}
	return w.viewport
}

// wrapInner represents a window that is provided by an InnerWindow container in the canvas.
type wrapInner struct {
	fyne.Window
	inner *container.InnerWindow
	d     *gLDriver

	centered bool
	onClosed func()
}

func wrapInnerWindow(w *container.InnerWindow, root fyne.Window, d *gLDriver) fyne.Window {
	wrapped := &wrapInner{inner: w, d: d}
	wrapped.Window = root
	w.CloseIntercept = wrapped.doClose
	return wrapped
}

func (w *wrapInner) CenterOnScreen() {
	w.centered = true

	w.doCenter()
}

func (w *wrapInner) Close() {
	w.inner.Close()
}

func (w *wrapInner) Hide() {
	w.inner.Hide()
	w.updateVisibility()
}

func (w *wrapInner) Move(p fyne.Position) {
	w.inner.Move(p)
}

func (w *wrapInner) Resize(s fyne.Size) {
	w.inner.Resize(s)
}

func (w *wrapInner) SetContent(o fyne.CanvasObject) {
	w.inner.SetContent(o)
}

func (w *wrapInner) SetOnClosed(fn func()) {
	w.onClosed = fn
}

func (w *wrapInner) Show() {
	c := w.Window.Canvas().(*glCanvas)
	multi := c.webExtraWindows
	multi.Show()
	w.inner.Show()

	c.Overlays().Add(multi)

	if w.centered {
		w.doCenter()
	}
}

func (w *wrapInner) doCenter() {
	c := w.Window.Canvas().(*glCanvas)
	multi := c.webExtraWindows

	min := w.inner.MinSize()
	min = min.Max(w.inner.Size())

	x := (multi.Size().Width - min.Width) / 2
	y := (multi.Size().Height - min.Height) / 2

	w.inner.Move(fyne.NewPos(x, y))
}

func (w *wrapInner) doClose() {
	c := w.Window.Canvas().(*glCanvas)
	multi := c.webExtraWindows

	pos := -1
	for i, child := range multi.Windows {
		if child == w.inner {
			pos = i
			w.inner.Hide()
			break
		}
	}
	if pos != -1 {
		count := len(multi.Windows)
		copy(multi.Windows[pos:], multi.Windows[pos+1:])
		multi.Windows[count-1] = nil
		multi.Windows = multi.Windows[:count-1]
	}

	if w.onClosed != nil {
		w.onClosed()
	}
	w.updateVisibility()
}

func (w *wrapInner) updateVisibility() {
	c := w.Window.Canvas().(*glCanvas)
	multi := c.webExtraWindows

	visible := 0
	for _, win := range multi.Windows {
		if win.Visible() {
			visible++
		}
	}

	if visible > 0 {
		multi.Refresh()
	} else {
		multi.Hide()
		c.Overlays().Remove(multi)
	}
}

func (w *window) scaleInput(in float64) float64 {
	return in
}
