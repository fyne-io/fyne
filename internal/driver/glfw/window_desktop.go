// +build !js,!wasm,!web

package glfw

import (
	"bytes"
	"image"
	_ "image/png" // for the icon
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/painter/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	scrollSpeed      = 10
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

var (
	cursorMap    map[desktop.Cursor]*glfw.Cursor
	defaultTitle = "Fyne Application"
)

func initCursors() {
	cursorMap = map[desktop.Cursor]*glfw.Cursor{
		desktop.DefaultCursor:   glfw.CreateStandardCursor(glfw.ArrowCursor),
		desktop.TextCursor:      glfw.CreateStandardCursor(glfw.IBeamCursor),
		desktop.CrosshairCursor: glfw.CreateStandardCursor(glfw.CrosshairCursor),
		desktop.PointerCursor:   glfw.CreateStandardCursor(glfw.HandCursor),
		desktop.HResizeCursor:   glfw.CreateStandardCursor(glfw.HResizeCursor),
		desktop.VResizeCursor:   glfw.CreateStandardCursor(glfw.VResizeCursor),
	}
}

// Declare conformity to Window interface
var _ fyne.Window = (*window)(nil)

type window struct {
	viewport   *glfw.Window
	createLock sync.Once
	decorate   bool
	fixedSize  bool

	cursor   *glfw.Cursor
	painted  int // part of the macOS GL fix, updated GLFW should fix this
	canvas   *glCanvas
	title    string
	icon     fyne.Resource
	mainmenu *fyne.MainMenu

	clipboard fyne.Clipboard

	master     bool
	fullScreen bool
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
	mousePressed       fyne.CanvasObject
	onClosed           func()

	xpos, ypos    int
	size          fyne.PixelSize // This field should be considered read only

	eventLock  sync.RWMutex
	eventQueue chan func()
	eventWait  sync.WaitGroup
	pending    []func()
}

func (w *window) SetFullScreen(full bool) {
	if full {
		w.fullScreen = true
	}
	w.runOnMainWhenCreated(func() {
		monitor := w.getMonitorForWindow()
		mode := monitor.GetVideoMode()

		if full {
			w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		} else {
			w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.size.Width, w.size.Height, 0)
			w.fullScreen = false
		}
	})
}

// centerOnScreen handles the logic for centering a window
func (w *window) centerOnScreen() {
	w.runOnMainWhenCreated(func() {
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
	}) // end of runOnMain(){}
}

func (w *window) RequestFocus() {
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

func (w *window) fitContent() {
	w.canvas.RLock()
	content := w.canvas.content
	w.canvas.RUnlock()
	if content == nil {
		return
	}

	if w.viewport == nil {
		return
	}

	minWidth, minHeight := w.minSizeOnScreen()
	if w.fixedSize {
		limitSize := internal.ScaleSize(w.canvas, w.Canvas().Size())

		w.viewport.SetSizeLimits(limitSize.Width, limitSize.Height, limitSize.Width, limitSize.Height)
		w.viewport.SetSize(limitSize.Width, limitSize.Height)
	} else {
		// Adjust size to be at least minSize
		if w.size.Width < minWidth || w.size.Height < minHeight {
			targetWidth, targetHeight := w.size.Width, w.size.Height
			if targetWidth < minWidth {
				targetWidth = minWidth
			}
			if targetHeight < minHeight {
				targetHeight = minHeight
			}
			w.viewport.SetSize(targetWidth, targetHeight)
		}

		w.viewport.SetSizeLimits(minWidth, minHeight, glfw.DontCare, glfw.DontCare)
	}
}

func (w *window) getMonitorForWindow() *glfw.Monitor {
	xOff := w.xpos + (w.size.Width / 2)
	yOff := w.ypos + (w.size.Height / 2)

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

func (w *window) frameSized(_*glfw.Window, width, height int) {
	w.processFrameSized(width, height)
}

func (w *window) refresh(_ *glfw.Window) {
	w.processRefresh()
}

func (w *window) closed(viewport *glfw.Window) {
	viewport.SetShouldClose(true)

	w.processClosed()
}

func (w *window) rescaleOnMain() {
	if w.viewport == nil {
		return
	}
	w.fitContent()

	if w.fullScreen {
		w.size.Width, w.size.Height = w.viewport.GetSize()
		scaledFull := internal.UnscaleSize(w.canvas, w.size)
		w.canvas.Resize(scaledFull)
		return
	}

	size := w.canvas.size.Union(w.canvas.MinSize())
	newWidth, newHeight := w.screenSize(size)
	w.viewport.SetSize(newWidth, newHeight)
}

func fyneToNativeCursor(cursor desktop.Cursor) *glfw.Cursor {
	ret, ok := cursorMap[cursor]
	if !ok {
		return cursorMap[desktop.DefaultCursor]
	}
	return ret
}

func (w *window) SetCursor(cursor *glfw.Cursor) {
	w.viewport.SetCursor(cursor)
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
	glfw.Key0: fyne.Key0,
	glfw.Key1: fyne.Key1,
	glfw.Key2: fyne.Key2,
	glfw.Key3: fyne.Key3,
	glfw.Key4: fyne.Key4,
	glfw.Key5: fyne.Key5,
	glfw.Key6: fyne.Key6,
	glfw.Key7: fyne.Key7,
	glfw.Key8: fyne.Key8,
	glfw.Key9: fyne.Key9,

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

var keyNameMap = map[string]fyne.KeyName{
	"'": fyne.KeyApostrophe,
	",": fyne.KeyComma,
	"-": fyne.KeyMinus,
	".": fyne.KeyPeriod,
	"/": fyne.KeySlash,
	"`": fyne.KeyBackTick,

	";": fyne.KeySemicolon,
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

func convertAction(action glfw.Action) desktop.Action {
	switch action {
	case glfw.Press:
		return desktop.Press
	case glfw.Release:
		return desktop.Release
	case glfw.Repeat:
		return desktop.Repeat
	}
	panic("Could not convert glfw.Action.")
}

func (w *window) keyPressed(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key, scancode)
	keyDesktopModifier := desktopModifier(mods)
	keyAction := convertAction(action)

	w.processKeyPressed(keyName, scancode, keyAction, keyDesktopModifier)
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
func (w *window) charInput(viewport *glfw.Window, char rune) {
	w.processCharInput(char)
}

func (w *window) focused(_ *glfw.Window, focused bool) {
	w.processFocused(focused)
}

func (w *window) DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func (w *window) create() {
	runOnMain(func() {
		// make the window hidden, we will set it up and then show it later
		glfw.WindowHint(glfw.Visible, 0)
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
		if pixWidth == 0 {
			pixWidth = 10
		}
		if pixHeight == 0 {
			pixHeight = 10
		}

		win, err := glfw.CreateWindow(pixWidth, pixHeight, w.title, nil, nil)
		if err != nil {
			fyne.LogError("window creation error", err)
			return
		}
		w.viewport = win
	})

	// run the GL init on the draw thread
	runOnDraw(w, func() {
		w.canvas.painter = gl.NewPainter(w.canvas, w)
		w.canvas.painter.Init()
	})

	runOnMain(func() {
		win := w.viewport
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
		winWidth, _ := win.GetSize()
		texWidth, _ := win.GetFramebufferSize()
		w.canvas.texScale = float32(texWidth) / float32(winWidth)

		for _, fn := range w.pending {
			fn()
		}
		w.fitContent()
	})
}
