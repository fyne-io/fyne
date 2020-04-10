// +build js wasm web

package glfw

import (
	_ "image/png" // for the icon
	"runtime"
	"sync"
	"time"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal"
	"fyne.io/fyne/internal/painter/gl"

	"github.com/goxjs/glfw"
)

const (
	scrollSpeed      = 10
	doubleClickDelay = 500 // ms (maximum interval between clicks for double click detection)
)

type Cursor struct {
}

var (
        cursorMap    map[desktop.Cursor]*Cursor
        defaultTitle = "Fyne Application"
)

// Declare conformity to Window interface
var _ fyne.Window = (*window)(nil)

type window struct {
	viewport   *glfw.Window
	createLock sync.Once
	decorate   bool
	fixedSize  bool

	painted  int // part of the macOS GL fix, updated GLFW should fix this
	canvas   *glCanvas
	title    string
	icon     fyne.Resource
	mainmenu *fyne.MainMenu
	cursor   *Cursor

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
	size          fyne.PixelSize

	eventLock  sync.RWMutex
	eventQueue chan func()
	eventWait  sync.WaitGroup
	pending    []func()
}

func (w *window) SetFullScreen(full bool) {
//	if full {
		w.fullScreen = true
//	} else {
//		w.fullScreen = false
//	}
}

// centerOnScreen handles the logic for centering a window
func (w *window) centerOnScreen() {
	// exit immediately if window is not supposed to be centered
	if !w.centered {
		return
	}

	// FIXME: not supported with WebGL
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

func (w *window) getMonitorForWindow() *glfw.Monitor {
	return glfw.GetPrimaryMonitor()
}

func (w *window) detectScale() float32 {
	return scaleForDpi(int(96))
}

func (w *window) moved(_ *glfw.Window, x, y int) {
	w.processMoved(x, y)
}

func (w *window) resized(_ *glfw.Window, width, height int) {
	log.Println("resized : ", width, height, w.calculatedScale())
	w.canvas.scale = w.calculatedScale()
	w.processResized(width, height)
}

func (w *window) frameSized(_ *glfw.Window, width, height int) {
	log.Println("framesize : ", width, height)
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

//	if w.fullScreen {
		w.size.Width, w.size.Height = w.viewport.GetSize()
		scaledFull := internal.UnscaleSize(w.canvas, w.size)
		w.canvas.Resize(scaledFull)
		return
//	}

//	size := w.canvas.size.Union(w.canvas.MinSize())
//	newWidth, newHeight := w.screenSize(size)
//	w.viewport.SetSize(newWidth, newHeight)
}

func fyneToNativeCursor(cursor desktop.Cursor) *Cursor {
	return nil
}

func (w *window) SetCursor(_ *Cursor) {
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

func convertAction(action glfw.Action) desktop.Action {
	switch action {
	case glfw.Press:
		return desktop.Press;
	case glfw.Release:
		return desktop.Release;
	case glfw.Repeat:
		return desktop.Repeat;
	}
	panic("Could not convert glfw.Action.")
}

func (w *window) keyPressed(viewport *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyName := keyToName(key)
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
// Unicode character is input regardless of what modifier keys are used.
//
// The character with modifiers callback is intended for implementing custom
// Unicode character input. Characters do not map 1:1 to physical keys,
// as a key may produce zero, one or more characters.
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
		// we can't hide the window in webgl, so there might be some artifact
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
		glfw.DetachCurrentContext()

		w.canvas.detectedScale = w.detectScale()
		w.canvas.SetScale(w.calculatedScale())
		w.canvas.scale = w.calculatedScale()
		for _, fn := range w.pending {
			fn()
		}

		width, height := win.GetSize()
		log.Println("created: ", width, height)
		w.processFrameSized(width, height)
		w.processResized(width, height)
	})
}
