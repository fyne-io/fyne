//go:build !js
// +build !js

package glfw

import (
	"io"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

var contextWatcher ContextWatcher

// Init initializes the library.
//
// A valid ContextWatcher must be provided. It gets notified when context becomes current or detached.
// It should be provided by the GL bindings you are using, so you can do glfw.Init(gl.ContextWatcher).
func Init(cw ContextWatcher) error {
	contextWatcher = cw
	return glfw.Init()
}

func Terminate() {
	glfw.Terminate()
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var m *glfw.Monitor
	if monitor != nil {
		m = monitor.Monitor
	}
	var s *glfw.Window
	if share != nil {
		s = share.Window
	}

	w, err := glfw.CreateWindow(width, height, title, m, s)
	if err != nil {
		return nil, err
	}

	window := &Window{Window: w}

	return window, err
}

func SwapInterval(interval int) {
	glfw.SwapInterval(interval)
}

func (w *Window) MakeContextCurrent() {
	w.Window.MakeContextCurrent()
	// In reality, context is available on each platform via GetGLXContext, GetWGLContext, GetNSGLContext, etc.
	// Pretend it is not available and pass nil, since it's not actually needed at this time.
	contextWatcher.OnMakeCurrent(nil)
}

func DetachCurrentContext() {
	glfw.DetachCurrentContext()
	contextWatcher.OnDetach()
}

type Window struct {
	*glfw.Window
}

type Monitor struct {
	*glfw.Monitor
}

func GetPrimaryMonitor() *Monitor {
	m := glfw.GetPrimaryMonitor()
	return &Monitor{Monitor: m}
}

func PollEvents() {
	glfw.PollEvents()
}

type CursorPosCallback func(w *Window, xpos float64, ypos float64)

func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos float64, ypos float64) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetCursorPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type MouseMovementCallback func(w *Window, xpos float64, ypos float64, xdelta float64, ydelta float64)

var lastMousePos [2]float64 // HACK.

// TODO: For now, this overrides SetCursorPosCallback; should support both.
func (w *Window) SetMouseMovementCallback(cbfun MouseMovementCallback) (previous MouseMovementCallback) {
	lastMousePos[0], lastMousePos[1] = w.Window.GetCursorPos()
	wrappedCbfun := func(_ *glfw.Window, xpos float64, ypos float64) {
		xdelta, ydelta := xpos-lastMousePos[0], ypos-lastMousePos[1]
		lastMousePos[0], lastMousePos[1] = xpos, ypos
		cbfun(w, xpos, ypos, xdelta, ydelta)
	}

	p := w.Window.SetCursorPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	wrappedCbfun := func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, Key(key), scancode, Action(action), ModifierKey(mods))
	}

	p := w.Window.SetKeyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CharCallback func(w *Window, char rune)

func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	wrappedCbfun := func(_ *glfw.Window, char rune) {
		cbfun(w, char)
	}

	p := w.Window.SetCharCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type ScrollCallback func(w *Window, xoff float64, yoff float64)

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	wrappedCbfun := func(_ *glfw.Window, xoff float64, yoff float64) {
		cbfun(w, xoff, yoff)
	}

	p := w.Window.SetScrollCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	wrappedCbfun := func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		cbfun(w, MouseButton(button), Action(action), ModifierKey(mods))
	}

	p := w.Window.SetMouseButtonCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FramebufferSizeCallback func(w *Window, width int, height int)

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetFramebufferSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

func (w *Window) GetKey(key Key) Action {
	a := w.Window.GetKey(glfw.Key(key))
	return Action(a)
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	a := w.Window.GetMouseButton(glfw.MouseButton(button))
	return Action(a)
}

func (w *Window) GetInputMode(mode InputMode) int {
	return w.Window.GetInputMode(glfw.InputMode(mode))
}

func (w *Window) SetInputMode(mode InputMode, value int) {
	w.Window.SetInputMode(glfw.InputMode(mode), value)
}

type Key glfw.Key

const (
	KeySpace        = Key(glfw.KeySpace)
	KeyApostrophe   = Key(glfw.KeyApostrophe)
	KeyComma        = Key(glfw.KeyComma)
	KeyMinus        = Key(glfw.KeyMinus)
	KeyPeriod       = Key(glfw.KeyPeriod)
	KeySlash        = Key(glfw.KeySlash)
	Key0            = Key(glfw.Key0)
	Key1            = Key(glfw.Key1)
	Key2            = Key(glfw.Key2)
	Key3            = Key(glfw.Key3)
	Key4            = Key(glfw.Key4)
	Key5            = Key(glfw.Key5)
	Key6            = Key(glfw.Key6)
	Key7            = Key(glfw.Key7)
	Key8            = Key(glfw.Key8)
	Key9            = Key(glfw.Key9)
	KeySemicolon    = Key(glfw.KeySemicolon)
	KeyEqual        = Key(glfw.KeyEqual)
	KeyA            = Key(glfw.KeyA)
	KeyB            = Key(glfw.KeyB)
	KeyC            = Key(glfw.KeyC)
	KeyD            = Key(glfw.KeyD)
	KeyE            = Key(glfw.KeyE)
	KeyF            = Key(glfw.KeyF)
	KeyG            = Key(glfw.KeyG)
	KeyH            = Key(glfw.KeyH)
	KeyI            = Key(glfw.KeyI)
	KeyJ            = Key(glfw.KeyJ)
	KeyK            = Key(glfw.KeyK)
	KeyL            = Key(glfw.KeyL)
	KeyM            = Key(glfw.KeyM)
	KeyN            = Key(glfw.KeyN)
	KeyO            = Key(glfw.KeyO)
	KeyP            = Key(glfw.KeyP)
	KeyQ            = Key(glfw.KeyQ)
	KeyR            = Key(glfw.KeyR)
	KeyS            = Key(glfw.KeyS)
	KeyT            = Key(glfw.KeyT)
	KeyU            = Key(glfw.KeyU)
	KeyV            = Key(glfw.KeyV)
	KeyW            = Key(glfw.KeyW)
	KeyX            = Key(glfw.KeyX)
	KeyY            = Key(glfw.KeyY)
	KeyZ            = Key(glfw.KeyZ)
	KeyLeftBracket  = Key(glfw.KeyLeftBracket)
	KeyBackslash    = Key(glfw.KeyBackslash)
	KeyRightBracket = Key(glfw.KeyRightBracket)
	KeyGraveAccent  = Key(glfw.KeyGraveAccent)
	KeyWorld1       = Key(glfw.KeyWorld1)
	KeyWorld2       = Key(glfw.KeyWorld2)
	KeyEscape       = Key(glfw.KeyEscape)
	KeyEnter        = Key(glfw.KeyEnter)
	KeyTab          = Key(glfw.KeyTab)
	KeyBackspace    = Key(glfw.KeyBackspace)
	KeyInsert       = Key(glfw.KeyInsert)
	KeyDelete       = Key(glfw.KeyDelete)
	KeyRight        = Key(glfw.KeyRight)
	KeyLeft         = Key(glfw.KeyLeft)
	KeyDown         = Key(glfw.KeyDown)
	KeyUp           = Key(glfw.KeyUp)
	KeyPageUp       = Key(glfw.KeyPageUp)
	KeyPageDown     = Key(glfw.KeyPageDown)
	KeyHome         = Key(glfw.KeyHome)
	KeyEnd          = Key(glfw.KeyEnd)
	KeyCapsLock     = Key(glfw.KeyCapsLock)
	KeyScrollLock   = Key(glfw.KeyScrollLock)
	KeyNumLock      = Key(glfw.KeyNumLock)
	KeyPrintScreen  = Key(glfw.KeyPrintScreen)
	KeyPause        = Key(glfw.KeyPause)
	KeyF1           = Key(glfw.KeyF1)
	KeyF2           = Key(glfw.KeyF2)
	KeyF3           = Key(glfw.KeyF3)
	KeyF4           = Key(glfw.KeyF4)
	KeyF5           = Key(glfw.KeyF5)
	KeyF6           = Key(glfw.KeyF6)
	KeyF7           = Key(glfw.KeyF7)
	KeyF8           = Key(glfw.KeyF8)
	KeyF9           = Key(glfw.KeyF9)
	KeyF10          = Key(glfw.KeyF10)
	KeyF11          = Key(glfw.KeyF11)
	KeyF12          = Key(glfw.KeyF12)
	KeyF13          = Key(glfw.KeyF13)
	KeyF14          = Key(glfw.KeyF14)
	KeyF15          = Key(glfw.KeyF15)
	KeyF16          = Key(glfw.KeyF16)
	KeyF17          = Key(glfw.KeyF17)
	KeyF18          = Key(glfw.KeyF18)
	KeyF19          = Key(glfw.KeyF19)
	KeyF20          = Key(glfw.KeyF20)
	KeyF21          = Key(glfw.KeyF21)
	KeyF22          = Key(glfw.KeyF22)
	KeyF23          = Key(glfw.KeyF23)
	KeyF24          = Key(glfw.KeyF24)
	KeyF25          = Key(glfw.KeyF25)
	KeyKP0          = Key(glfw.KeyKP0)
	KeyKP1          = Key(glfw.KeyKP1)
	KeyKP2          = Key(glfw.KeyKP2)
	KeyKP3          = Key(glfw.KeyKP3)
	KeyKP4          = Key(glfw.KeyKP4)
	KeyKP5          = Key(glfw.KeyKP5)
	KeyKP6          = Key(glfw.KeyKP6)
	KeyKP7          = Key(glfw.KeyKP7)
	KeyKP8          = Key(glfw.KeyKP8)
	KeyKP9          = Key(glfw.KeyKP9)
	KeyKPDecimal    = Key(glfw.KeyKPDecimal)
	KeyKPDivide     = Key(glfw.KeyKPDivide)
	KeyKPMultiply   = Key(glfw.KeyKPMultiply)
	KeyKPSubtract   = Key(glfw.KeyKPSubtract)
	KeyKPAdd        = Key(glfw.KeyKPAdd)
	KeyKPEnter      = Key(glfw.KeyKPEnter)
	KeyKPEqual      = Key(glfw.KeyKPEqual)
	KeyLeftShift    = Key(glfw.KeyLeftShift)
	KeyLeftControl  = Key(glfw.KeyLeftControl)
	KeyLeftAlt      = Key(glfw.KeyLeftAlt)
	KeyLeftSuper    = Key(glfw.KeyLeftSuper)
	KeyRightShift   = Key(glfw.KeyRightShift)
	KeyRightControl = Key(glfw.KeyRightControl)
	KeyRightAlt     = Key(glfw.KeyRightAlt)
	KeyRightSuper   = Key(glfw.KeyRightSuper)
	KeyMenu         = Key(glfw.KeyMenu)
)

type MouseButton glfw.MouseButton

const (
	MouseButton1 = MouseButton(glfw.MouseButton1)
	MouseButton2 = MouseButton(glfw.MouseButton2)
	MouseButton3 = MouseButton(glfw.MouseButton3)

	MouseButtonLeft   = MouseButton(glfw.MouseButtonLeft)
	MouseButtonRight  = MouseButton(glfw.MouseButtonRight)
	MouseButtonMiddle = MouseButton(glfw.MouseButtonMiddle)
)

type Action glfw.Action

const (
	Release = Action(glfw.Release)
	Press   = Action(glfw.Press)
	Repeat  = Action(glfw.Repeat)
)

type InputMode int

const (
	CursorMode             = InputMode(glfw.CursorMode)
	StickyKeysMode         = InputMode(glfw.StickyKeysMode)
	StickyMouseButtonsMode = InputMode(glfw.StickyMouseButtonsMode)
	LockKeyMods            = InputMode(glfw.LockKeyMods)
	RawMouseMotion         = InputMode(glfw.RawMouseMotion)
)

const (
	CursorNormal   = int(glfw.CursorNormal)
	CursorHidden   = int(glfw.CursorHidden)
	CursorDisabled = int(glfw.CursorDisabled)
)

type ModifierKey int

const (
	ModShift   = ModifierKey(glfw.ModShift)
	ModControl = ModifierKey(glfw.ModControl)
	ModAlt     = ModifierKey(glfw.ModAlt)
	ModSuper   = ModifierKey(glfw.ModSuper)
)

// Open opens a named asset. It's the caller's responsibility to close it when done.
//
// For now, assets are read directly from the current working directory.
func Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// ---

func WaitEvents() {
	glfw.WaitEvents()
}

func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func DefaultWindowHints() {
	glfw.DefaultWindowHints()
}

func (w *Window) SetClipboardString(str string) {
	glfw.SetClipboardString(str)
}

func (w *Window) GetClipboardString() (string, error) {
	return glfw.GetClipboardString(), nil
}

type CloseCallback func(w *Window)

func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetCloseCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type RefreshCallback func(w *Window)

func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	wrappedCbfun := func(_ *glfw.Window) {
		cbfun(w)
	}

	p := w.Window.SetRefreshCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type SizeCallback func(w *Window, width int, height int)

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	wrappedCbfun := func(_ *glfw.Window, width int, height int) {
		cbfun(w, width, height)
	}

	p := w.Window.SetSizeCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CursorEnterCallback func(w *Window, entered bool)

func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {
	wrappedCbfun := func(_ *glfw.Window, entered bool) {
		cbfun(w, entered)
	}

	p := w.Window.SetCursorEnterCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type CharModsCallback func(w *Window, char rune, mods ModifierKey)

func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	wrappedCbfun := func(_ *glfw.Window, char rune, mods glfw.ModifierKey) {
		cbfun(w, char, ModifierKey(mods))
	}

	p := w.Window.SetCharModsCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type PosCallback func(w *Window, xpos int, ypos int)

func (w *Window) SetPosCallback(cbfun PosCallback) (previous PosCallback) {
	wrappedCbfun := func(_ *glfw.Window, xpos int, ypos int) {
		cbfun(w, xpos, ypos)
	}

	p := w.Window.SetPosCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type FocusCallback func(w *Window, focused bool)

func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	wrappedCbfun := func(_ *glfw.Window, focused bool) {
		cbfun(w, focused)
	}

	p := w.Window.SetFocusCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type IconifyCallback func(w *Window, iconified bool)

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	wrappedCbfun := func(_ *glfw.Window, iconified bool) {
		cbfun(w, iconified)
	}

	p := w.Window.SetIconifyCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	wrappedCbfun := func(_ *glfw.Window, names []string) {
		cbfun(w, names)
	}

	p := w.Window.SetDropCallback(wrappedCbfun)
	_ = p

	// TODO: Handle previous.
	return nil
}
