//go:build js && !wasm
// +build js,!wasm

package glfw

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

var contextWatcher ContextWatcher

func Init(cw ContextWatcher) error {
	contextWatcher = cw
	return nil
}

func Terminate() error {
	return nil
}

func CreateWindow(_, _ int, title string, monitor *Monitor, share *Window) (*Window, error) {
	// THINK: Consider https://developer.mozilla.org/en-US/docs/Web/API/Window.open?

	// HACK: Go fullscreen?
	width := dom.GetWindow().InnerWidth()
	height := dom.GetWindow().InnerHeight()

	canvas := document.CreateElement("canvas").(*dom.HTMLCanvasElement)

	devicePixelRatio := js.Global.Get("devicePixelRatio").Float()
	canvas.Width = int(float64(width)*devicePixelRatio + 0.5)   // Nearest non-negative int.
	canvas.Height = int(float64(height)*devicePixelRatio + 0.5) // Nearest non-negative int.
	canvas.Style().SetProperty("width", fmt.Sprintf("%vpx", width), "")
	canvas.Style().SetProperty("height", fmt.Sprintf("%vpx", height), "")

	if js.Global.Get("document").Get("body") == nil {
		body := js.Global.Get("document").Call("createElement", "body")
		js.Global.Get("document").Set("body", body)
	}
	document.Body().Style().SetProperty("margin", "0", "")
	document.Body().AppendChild(canvas)

	document.SetTitle(title)

	// DEBUG: Add framebuffer information div.
	if false {
		//canvas.Height -= 30
		text := document.CreateElement("div")
		textContent := fmt.Sprintf("%v %v (%v) @%v", dom.GetWindow().InnerWidth(), canvas.Width, float64(width)*devicePixelRatio, devicePixelRatio)
		text.SetTextContent(textContent)
		document.Body().AppendChild(text)
	}

	// Use glfw hints.
	attrs := defaultAttributes()
	attrs.Alpha = (hints[AlphaBits] > 0)
	if _, ok := hints[DepthBits]; ok {
		attrs.Depth = (hints[DepthBits] > 0)
	}
	attrs.Stencil = (hints[StencilBits] > 0)
	attrs.Antialias = (hints[Samples] > 0)
	attrs.PremultipliedAlpha = (hints[PremultipliedAlpha] > 0)
	attrs.PreserveDrawingBuffer = (hints[PreserveDrawingBuffer] > 0)
	attrs.PreferLowPowerToHighPerformance = (hints[PreferLowPowerToHighPerformance] > 0)
	attrs.FailIfMajorPerformanceCaveat = (hints[FailIfMajorPerformanceCaveat] > 0)

	// Create GL context.
	context, err := newContext(canvas.Underlying(), attrs)
	if err != nil {
		return nil, err
	}

	w := &Window{
		canvas:           canvas,
		context:          context,
		devicePixelRatio: devicePixelRatio,
	}

	if w.canvas.Underlying().Get("requestPointerLock") == js.Undefined ||
		document.Underlying().Get("exitPointerLock") == js.Undefined {

		w.missing.pointerLock = true
	}
	if w.canvas.Underlying().Get("webkitRequestFullscreen") == js.Undefined ||
		document.Underlying().Get("webkitExitFullscreen") == js.Undefined {

		w.missing.fullscreen = true
	}

	if monitor != nil {
		if w.missing.fullscreen {
			log.Println("warning: Fullscreen API unsupported")
		} else {
			w.requestFullscreen = true
		}
	}

	dom.GetWindow().AddEventListener("resize", false, func(event dom.Event) {
		// HACK: Go fullscreen?
		width := dom.GetWindow().InnerWidth()
		height := dom.GetWindow().InnerHeight()

		devicePixelRatio := js.Global.Get("devicePixelRatio").Float()
		w.canvas.Width = int(float64(width)*devicePixelRatio + 0.5)   // Nearest non-negative int.
		w.canvas.Height = int(float64(height)*devicePixelRatio + 0.5) // Nearest non-negative int.

		if w.framebufferSizeCallback != nil {
			// TODO: Callbacks may be blocking so they need to happen asyncronously. However,
			//       GLFW API promises the callbacks will occur from one thread (i.e., sequentially), so may want to do that.
			go w.framebufferSizeCallback(w, w.canvas.Width, w.canvas.Height)
		}
		if w.sizeCallback != nil {
			boundingW, boundingH := width, height
			go w.sizeCallback(w, boundingW, boundingH)
		}
	})

	document.AddEventListener("keydown", false, func(event dom.Event) {
		w.goFullscreenIfRequested()

		ke := event.(*dom.KeyboardEvent)

		action := Press
		if ke.Repeat {
			action = Repeat
		}

		key := toKey(ke)

		// Extend slice if needed.
		neededSize := int(key) + 1
		if neededSize > len(w.keys) {
			w.keys = append(w.keys, make([]Action, neededSize-len(w.keys))...)
		}
		w.keys[key] = action

		if w.keyCallback != nil {
			mods := toModifierKey(ke)

			go w.keyCallback(w, key, -1, action, mods)
		}

		if w.charCallback != nil {
			if len(ke.Key) == 1 {
				keyRune := []rune(ke.Key)
				go w.charCallback(w, keyRune[0])
			}
		}

		ke.PreventDefault()
	})
	document.AddEventListener("keyup", false, func(event dom.Event) {
		w.goFullscreenIfRequested()

		ke := event.(*dom.KeyboardEvent)

		key := toKey(ke)

		// Extend slice if needed.
		neededSize := int(key) + 1
		if neededSize > len(w.keys) {
			w.keys = append(w.keys, make([]Action, neededSize-len(w.keys))...)
		}
		w.keys[key] = Release

		if w.keyCallback != nil {
			mods := toModifierKey(ke)

			go w.keyCallback(w, key, -1, Release, mods)
		}

		ke.PreventDefault()
	})

	document.AddEventListener("mousedown", false, func(event dom.Event) {
		w.goFullscreenIfRequested()

		me := event.(*dom.MouseEvent)
		if !(me.Button >= 0 && me.Button <= 2) {
			return
		}

		w.mouseButton[me.Button] = Press
		if w.mouseButtonCallback != nil {
			go w.mouseButtonCallback(w, MouseButton(me.Button), Press, 0)
		}

		me.PreventDefault()
	})
	document.AddEventListener("mouseup", false, func(event dom.Event) {
		w.goFullscreenIfRequested()

		me := event.(*dom.MouseEvent)
		if !(me.Button >= 0 && me.Button <= 2) {
			return
		}

		w.mouseButton[me.Button] = Release
		if w.mouseButtonCallback != nil {
			go w.mouseButtonCallback(w, MouseButton(me.Button), Release, 0)
		}

		me.PreventDefault()
	})
	document.AddEventListener("contextmenu", false, func(event dom.Event) {
		event.PreventDefault()
	})

	document.AddEventListener("mousemove", false, func(event dom.Event) {
		me := event.(*dom.MouseEvent)

		var movementX, movementY float64
		if !w.missing.pointerLock {
			movementX = float64(me.MovementX)
			movementY = float64(me.MovementY)
		} else {
			movementX = float64(me.ClientX) - w.cursorPos[0]
			movementY = float64(me.ClientY) - w.cursorPos[1]
		}

		w.cursorPos[0], w.cursorPos[1] = float64(me.ClientX), float64(me.ClientY)
		if w.cursorPosCallback != nil {
			go w.cursorPosCallback(w, w.cursorPos[0], w.cursorPos[1])
		}
		if w.mouseMovementCallback != nil {
			go w.mouseMovementCallback(w, w.cursorPos[0], w.cursorPos[1], movementX, movementY)
		}

		me.PreventDefault()
	})
	document.AddEventListener("wheel", false, func(event dom.Event) {
		we := event.(*dom.WheelEvent)

		var multiplier float64
		switch we.DeltaMode {
		case dom.DeltaPixel:
			multiplier = 0.1
		case dom.DeltaLine:
			multiplier = 1
		default:
			log.Println("unsupported WheelEvent.DeltaMode:", we.DeltaMode)
			multiplier = 1
		}

		if w.scrollCallback != nil {
			go w.scrollCallback(w, -we.DeltaX*multiplier, -we.DeltaY*multiplier)
		}

		we.PreventDefault()
	})

	// Hacky mouse-emulation-via-touch.
	touchHandler := func(event dom.Event) {
		w.goFullscreenIfRequested()

		te := event.(*dom.TouchEvent)

		touches := te.Get("touches")
		if touches.Length() > 0 {
			t := touches.Index(0)

			if w.touches != nil && w.touches.Length() > 0 { // This event is a movement only if we previously had > 0 touch points.
				if w.mouseMovementCallback != nil {
					go w.mouseMovementCallback(w, t.Get("clientX").Float(), t.Get("clientY").Float(), t.Get("clientX").Float()-w.cursorPos[0], t.Get("clientY").Float()-w.cursorPos[1])
				}
			}

			w.cursorPos[0], w.cursorPos[1] = t.Get("clientX").Float(), t.Get("clientY").Float()
			if w.cursorPosCallback != nil {
				go w.cursorPosCallback(w, w.cursorPos[0], w.cursorPos[1])
			}
		}
		w.touches = touches

		te.PreventDefault()
	}
	document.AddEventListener("touchstart", false, touchHandler)
	document.AddEventListener("touchmove", false, touchHandler)
	document.AddEventListener("touchend", false, touchHandler)

	// Request first animation frame.
	js.Global.Call("requestAnimationFrame", animationFrame)

	return w, nil
}

func SwapInterval(interval int) error {
	// TODO: Implement.
	return nil
}

type Window struct {
	canvas            *dom.HTMLCanvasElement
	context           *js.Object
	requestFullscreen bool // requestFullscreen is set to true when fullscreen should be entered as soon as possible (in a user input handler).
	fullscreen        bool // fullscreen is true if we're currently in fullscreen mode.

	// Unavailable browser APIs.
	missing struct {
		pointerLock bool // Pointer Lock API.
		fullscreen  bool // Fullscreen API.
	}

	devicePixelRatio float64

	cursorMode  int
	cursorPos   [2]float64
	mouseButton [3]Action

	keys []Action

	cursorPosCallback       CursorPosCallback
	mouseMovementCallback   MouseMovementCallback
	mouseButtonCallback     MouseButtonCallback
	keyCallback             KeyCallback
	charCallback            CharCallback
	scrollCallback          ScrollCallback
	framebufferSizeCallback FramebufferSizeCallback
	sizeCallback            SizeCallback

	touches *js.Object // Hacky mouse-emulation-via-touch.
}

func (w *Window) SetPos(xpos, ypos int) {
	fmt.Println("not implemented: SetPos:", xpos, ypos)
}

func (w *Window) SetSize(width, height int) {
	fmt.Println("not implemented: SetSize:", width, height)
}

// goFullscreenIfRequested performs webkitRequestFullscreen if it was scheduled. It is called only from
// user events, because that API will fail if called at any other time.
func (w *Window) goFullscreenIfRequested() {
	if !w.requestFullscreen {
		return
	}
	w.requestFullscreen = false
	w.canvas.Underlying().Call("webkitRequestFullscreen")
	w.fullscreen = true
}

type Monitor struct{}

func (m *Monitor) GetVideoMode() *VidMode {
	return &VidMode{
		// HACK: Hardcoded sample values.
		// TODO: Try to get real values from browser via some API, if possible.
		Width:       1680,
		Height:      1050,
		RedBits:     8,
		GreenBits:   8,
		BlueBits:    8,
		RefreshRate: 60,
	}
}

func GetPrimaryMonitor() *Monitor {
	// TODO: Implement real functionality.
	return &Monitor{}
}

func PollEvents() error {
	return nil
}

func (w *Window) MakeContextCurrent() {
	contextWatcher.OnMakeCurrent(w.context)
}

func DetachCurrentContext() {
	contextWatcher.OnDetach()
}

func GetCurrentContext() *Window {
	panic("not implemented")
}

type CursorPosCallback func(w *Window, xpos float64, ypos float64)

func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	w.cursorPosCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type MouseMovementCallback func(w *Window, xpos float64, ypos float64, xdelta float64, ydelta float64)

func (w *Window) SetMouseMovementCallback(cbfun MouseMovementCallback) (previous MouseMovementCallback) {
	w.mouseMovementCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	w.keyCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type CharCallback func(w *Window, char rune)

func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	w.charCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type ScrollCallback func(w *Window, xoff float64, yoff float64)

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	w.scrollCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	w.mouseButtonCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type FramebufferSizeCallback func(w *Window, width int, height int)

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	w.framebufferSizeCallback = cbfun

	// TODO: Handle previous.
	return nil
}

func (w *Window) GetSize() (width, height int) {
	width = int(w.canvas.GetBoundingClientRect().Width*w.devicePixelRatio + 0.5)
	height = int(w.canvas.GetBoundingClientRect().Height*w.devicePixelRatio + 0.5)

	return width, height
}

func (w *Window) GetFramebufferSize() (width, height int) {
	return w.canvas.Width, w.canvas.Height
}

func (w *Window) GetPos() (x, y int) {
	// Not implemented.
	return
}

func (w *Window) ShouldClose() bool {
	return false
}

func (w *Window) SetShouldClose(value bool) {
	// TODO: Implement.
	// THINK: What should happen in the browser if we're told to "close" the window. Do we destroy/remove the canvas? Or nothing?
	//        Perhaps https://developer.mozilla.org/en-US/docs/Web/API/Window.close is relevant.
}

func (w *Window) SwapBuffers() error {
	<-animationFrameChan
	js.Global.Call("requestAnimationFrame", animationFrame)

	return nil
}

var animationFrameChan = make(chan struct{}, 1)

func animationFrame() {
	animationFrameChan <- struct{}{}
}

func (w *Window) GetCursorPos() (x, y float64) {
	return w.cursorPos[0], w.cursorPos[1]
}

var keyWarnings = 10

func (w *Window) GetKey(key Key) Action {
	if key == -1 && keyWarnings > 0 {
		// TODO: Implement all keys, get rid of this.
		keyWarnings--
		log.Println("GetKey: key not implemented.")
		return Release
	}
	if int(key) >= len(w.keys) {
		return Release
	}
	return w.keys[key]
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	if !(button >= 0 && button <= 2) {
		panic(fmt.Errorf("button is out of range: %v", button))
	}

	// Hacky mouse-emulation-via-touch.
	if w.touches != nil {
		switch button {
		case MouseButton1:
			if w.touches.Length() == 1 || w.touches.Length() == 3 {
				return Press
			}
		case MouseButton2:
			if w.touches.Length() == 2 || w.touches.Length() == 3 {
				return Press
			}
		}

		return Release
	}

	return w.mouseButton[button]
}

func (w *Window) GetInputMode(mode InputMode) int {
	switch mode {
	case CursorMode:
		return w.cursorMode
	default:
		panic(errors.New("not implemented"))
	}
}

var ErrInvalidParameter = errors.New("invalid parameter")
var ErrInvalidValue = errors.New("invalid value")

func (w *Window) SetInputMode(mode InputMode, value int) {
	switch mode {
	case CursorMode:
		// TODO; Make cursor API compatible with GLFW and Fyne use/expectation.
		/*
			if w.missing.pointerLock {
				log.Println("warning: Pointer Lock API unsupported")
				return
			}
			switch value {
			case CursorNormal:
				w.cursorMode = value
				document.Underlying().Call("exitPointerLock")
				w.canvas.Style().SetProperty("cursor", "initial", "")
				return
			case CursorHidden:
				w.cursorMode = value
				document.Underlying().Call("exitPointerLock")
				w.canvas.Style().SetProperty("cursor", "none", "")
				return
			case CursorDisabled:
				w.cursorMode = value
				w.canvas.Underlying().Call("requestPointerLock")
				return
			default:
				panic(ErrInvalidValue)
			}
		*/
		return
	case StickyKeysMode:
		panic(errors.New("not implemented"))
	case StickyMouseButtonsMode:
		panic(errors.New("not implemented"))
	default:
		panic(ErrInvalidParameter)
	}
}

type Key int

// TODO: Keys defined as -iota-1 need to be set to a valid positive value that matches the keyCode
//       generated by browsers. -iota-1 is used as a temporary solution to have unique but invalid values.
//       See https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/keyCode.
const (
	KeySpace        Key = 32
	KeyApostrophe   Key = 222
	KeyComma        Key = 188
	KeyMinus        Key = 189
	KeyPeriod       Key = 190
	KeySlash        Key = 191
	Key0            Key = 48
	Key1            Key = 49
	Key2            Key = 50
	Key3            Key = 51
	Key4            Key = 52
	Key5            Key = 53
	Key6            Key = 54
	Key7            Key = 55
	Key8            Key = 56
	Key9            Key = 57
	KeySemicolon    Key = 186
	KeyEqual        Key = 187
	KeyA            Key = 65
	KeyB            Key = 66
	KeyC            Key = 67
	KeyD            Key = 68
	KeyE            Key = 69
	KeyF            Key = 70
	KeyG            Key = 71
	KeyH            Key = 72
	KeyI            Key = 73
	KeyJ            Key = 74
	KeyK            Key = 75
	KeyL            Key = 76
	KeyM            Key = 77
	KeyN            Key = 78
	KeyO            Key = 79
	KeyP            Key = 80
	KeyQ            Key = 81
	KeyR            Key = 82
	KeyS            Key = 83
	KeyT            Key = 84
	KeyU            Key = 85
	KeyV            Key = 86
	KeyW            Key = 87
	KeyX            Key = 88
	KeyY            Key = 89
	KeyZ            Key = 90
	KeyLeftBracket  Key = 219
	KeyBackslash    Key = 220
	KeyRightBracket Key = 221
	KeyGraveAccent  Key = 192
	KeyWorld1       Key = -iota - 1
	KeyWorld2       Key = -iota - 1
	KeyEscape       Key = 27
	KeyEnter        Key = 13
	KeyTab          Key = 9
	KeyBackspace    Key = 8
	KeyInsert       Key = -iota - 1
	KeyDelete       Key = 46
	KeyRight        Key = 39
	KeyLeft         Key = 37
	KeyDown         Key = 40
	KeyUp           Key = 38
	KeyPageUp       Key = -iota - 1
	KeyPageDown     Key = -iota - 1
	KeyHome         Key = -iota - 1
	KeyEnd          Key = -iota - 1
	KeyCapsLock     Key = 20
	KeyScrollLock   Key = -iota - 1
	KeyNumLock      Key = -iota - 1
	KeyPrintScreen  Key = -iota - 1
	KeyPause        Key = -iota - 1
	KeyF1           Key = 112
	KeyF2           Key = 113
	KeyF3           Key = 114
	KeyF4           Key = 115
	KeyF5           Key = 116
	KeyF6           Key = 117
	KeyF7           Key = 118
	KeyF8           Key = 119
	KeyF9           Key = 120
	KeyF10          Key = 121
	KeyF11          Key = 122
	KeyF12          Key = 123
	KeyF13          Key = -iota - 1
	KeyF14          Key = -iota - 1
	KeyF15          Key = -iota - 1
	KeyF16          Key = -iota - 1
	KeyF17          Key = -iota - 1
	KeyF18          Key = -iota - 1
	KeyF19          Key = -iota - 1
	KeyF20          Key = -iota - 1
	KeyF21          Key = -iota - 1
	KeyF22          Key = -iota - 1
	KeyF23          Key = -iota - 1
	KeyF24          Key = -iota - 1
	KeyF25          Key = -iota - 1
	KeyKP0          Key = -iota - 1
	KeyKP1          Key = -iota - 1
	KeyKP2          Key = -iota - 1
	KeyKP3          Key = -iota - 1
	KeyKP4          Key = -iota - 1
	KeyKP5          Key = -iota - 1
	KeyKP6          Key = -iota - 1
	KeyKP7          Key = -iota - 1
	KeyKP8          Key = -iota - 1
	KeyKP9          Key = -iota - 1
	KeyKPDecimal    Key = -iota - 1
	KeyKPDivide     Key = -iota - 1
	KeyKPMultiply   Key = -iota - 1
	KeyKPSubtract   Key = -iota - 1
	KeyKPAdd        Key = -iota - 1
	KeyKPEnter      Key = -iota - 1
	KeyKPEqual      Key = -iota - 1
	KeyLeftShift    Key = 340
	KeyLeftControl  Key = 341
	KeyLeftAlt      Key = 342
	KeyLeftSuper    Key = 91
	KeyRightShift   Key = 344
	KeyRightControl Key = 345
	KeyRightAlt     Key = 346
	KeyRightSuper   Key = 93
	KeyMenu         Key = -iota - 1
)

// toKey extracts Key from given KeyboardEvent.
func toKey(ke *dom.KeyboardEvent) Key {
	key := Key(ke.KeyCode)
	switch {
	case key == 16 && ke.Location == dom.KeyLocationLeft:
		key = KeyLeftShift
	case key == 16 && ke.Location == dom.KeyLocationRight:
		key = KeyRightShift
	case key == 17 && ke.Location == dom.KeyLocationLeft:
		key = KeyLeftControl
	case key == 17 && ke.Location == dom.KeyLocationRight:
		key = KeyRightControl
	case key == 18 && ke.Location == dom.KeyLocationLeft:
		key = KeyLeftAlt
	case key == 18 && ke.Location == dom.KeyLocationRight:
		key = KeyRightAlt
	}
	return key
}

// toModifierKey extracts ModifierKey from given KeyboardEvent.
func toModifierKey(ke *dom.KeyboardEvent) ModifierKey {
	mods := ModifierKey(0)
	if ke.ShiftKey {
		mods += ModShift
	}
	if ke.CtrlKey {
		mods += ModControl
	}
	if ke.AltKey {
		mods += ModAlt
	}
	if ke.MetaKey {
		mods += ModSuper
	}
	return mods
}

type MouseButton int

const (
	MouseButton1 MouseButton = 0
	MouseButton2 MouseButton = 2 // Web MouseEvent has middle and right mouse buttons in reverse order.
	MouseButton3 MouseButton = 1 // Web MouseEvent has middle and right mouse buttons in reverse order.

	MouseButtonLeft   = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMiddle = MouseButton3
)

type Action int

const (
	Release Action = 0
	Press   Action = 1
	Repeat  Action = 2
)

type InputMode int

const (
	CursorMode InputMode = iota
	StickyKeysMode
	StickyMouseButtonsMode
	LockKeyMods
	RawMouseMotion
)

const (
	CursorNormal = iota
	CursorHidden
	CursorDisabled
)

type ModifierKey int

const (
	ModShift ModifierKey = (1 << iota)
	ModControl
	ModAlt
	ModSuper
)

// Open opens a named asset. It's the caller's responsibility to close it when done.
func Open(name string) (io.ReadCloser, error) {
	resp, err := http.Get(name)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status: %s", resp.Status)
	}
	return resp.Body, nil
}

// ---

func WaitEvents() {
	// TODO.

	runtime.Gosched()
}

func PostEmptyEvent() {
	// TODO: Implement.
}

func DefaultWindowHints() {
	// TODO: Implement.
}

func (w *Window) SetClipboardString(str string) {
	// TODO: Implement.
}
func (w *Window) GetClipboardString() (string, error) {
	// TODO: Implement.
	return "", errors.New("GetClipboardString not implemented")
}

func (w *Window) SetTitle(title string) {
	document.SetTitle(title)
}

func (w *Window) Show() {
	// TODO: Implement.
}

func (w *Window) Hide() {
	// TODO: Implement.
}

func (w *Window) Destroy() {
	document.Body().RemoveChild(w.canvas)
	if w.fullscreen {
		if w.missing.fullscreen {
			log.Println("warning: Fullscreen API unsupported")
		} else {
			document.Underlying().Call("webkitExitFullscreen")
			w.fullscreen = false
		}
	}
}

type CloseCallback func(w *Window)

func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type RefreshCallback func(w *Window)

func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type SizeCallback func(w *Window, width int, height int)

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.sizeCallback = cbfun

	// TODO: Handle previous.
	return nil
}

type CursorEnterCallback func(w *Window, entered bool)

func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type CharModsCallback func(w *Window, char rune, mods ModifierKey)

func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type PosCallback func(w *Window, xpos int, ypos int)

func (w *Window) SetPosCallback(cbfun PosCallback) (previous PosCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type FocusCallback func(w *Window, focused bool)

func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type IconifyCallback func(w *Window, iconified bool)

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	// TODO: Implement. Can use HTML5 file drag and drop API?

	// TODO: Handle previous.
	return nil
}
