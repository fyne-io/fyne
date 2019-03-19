package glfw

//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetJoystickCallbackCB();
//void glfwSetKeyCallbackCB(GLFWwindow *window);
//void glfwSetCharCallbackCB(GLFWwindow *window);
//void glfwSetCharModsCallbackCB(GLFWwindow *window);
//void glfwSetMouseButtonCallbackCB(GLFWwindow *window);
//void glfwSetCursorPosCallbackCB(GLFWwindow *window);
//void glfwSetCursorEnterCallbackCB(GLFWwindow *window);
//void glfwSetScrollCallbackCB(GLFWwindow *window);
//void glfwSetDropCallbackCB(GLFWwindow *window);
//float GetAxisAtIndex(float *axis, int i);
//unsigned char GetButtonsAtIndex(unsigned char *buttons, int i);
import "C"

import (
	"image"
	"image/draw"
	"unsafe"
)

var fJoystickHolder func(joy, event int)

// Joystick corresponds to a joystick.
type Joystick int

// Joystick IDs.
const (
	Joystick1    Joystick = C.GLFW_JOYSTICK_1
	Joystick2    Joystick = C.GLFW_JOYSTICK_2
	Joystick3    Joystick = C.GLFW_JOYSTICK_3
	Joystick4    Joystick = C.GLFW_JOYSTICK_4
	Joystick5    Joystick = C.GLFW_JOYSTICK_5
	Joystick6    Joystick = C.GLFW_JOYSTICK_6
	Joystick7    Joystick = C.GLFW_JOYSTICK_7
	Joystick8    Joystick = C.GLFW_JOYSTICK_8
	Joystick9    Joystick = C.GLFW_JOYSTICK_9
	Joystick10   Joystick = C.GLFW_JOYSTICK_10
	Joystick11   Joystick = C.GLFW_JOYSTICK_11
	Joystick12   Joystick = C.GLFW_JOYSTICK_12
	Joystick13   Joystick = C.GLFW_JOYSTICK_13
	Joystick14   Joystick = C.GLFW_JOYSTICK_14
	Joystick15   Joystick = C.GLFW_JOYSTICK_15
	Joystick16   Joystick = C.GLFW_JOYSTICK_16
	JoystickLast Joystick = C.GLFW_JOYSTICK_LAST
)

// Key corresponds to a keyboard key.
type Key int

// These key codes are inspired by the USB HID Usage Tables v1.12 (p. 53-60),
// but re-arranged to map to 7-bit ASCII for printable keys (function keys are
// put in the 256+ range).
const (
	KeyUnknown      Key = C.GLFW_KEY_UNKNOWN
	KeySpace        Key = C.GLFW_KEY_SPACE
	KeyApostrophe   Key = C.GLFW_KEY_APOSTROPHE
	KeyComma        Key = C.GLFW_KEY_COMMA
	KeyMinus        Key = C.GLFW_KEY_MINUS
	KeyPeriod       Key = C.GLFW_KEY_PERIOD
	KeySlash        Key = C.GLFW_KEY_SLASH
	Key0            Key = C.GLFW_KEY_0
	Key1            Key = C.GLFW_KEY_1
	Key2            Key = C.GLFW_KEY_2
	Key3            Key = C.GLFW_KEY_3
	Key4            Key = C.GLFW_KEY_4
	Key5            Key = C.GLFW_KEY_5
	Key6            Key = C.GLFW_KEY_6
	Key7            Key = C.GLFW_KEY_7
	Key8            Key = C.GLFW_KEY_8
	Key9            Key = C.GLFW_KEY_9
	KeySemicolon    Key = C.GLFW_KEY_SEMICOLON
	KeyEqual        Key = C.GLFW_KEY_EQUAL
	KeyA            Key = C.GLFW_KEY_A
	KeyB            Key = C.GLFW_KEY_B
	KeyC            Key = C.GLFW_KEY_C
	KeyD            Key = C.GLFW_KEY_D
	KeyE            Key = C.GLFW_KEY_E
	KeyF            Key = C.GLFW_KEY_F
	KeyG            Key = C.GLFW_KEY_G
	KeyH            Key = C.GLFW_KEY_H
	KeyI            Key = C.GLFW_KEY_I
	KeyJ            Key = C.GLFW_KEY_J
	KeyK            Key = C.GLFW_KEY_K
	KeyL            Key = C.GLFW_KEY_L
	KeyM            Key = C.GLFW_KEY_M
	KeyN            Key = C.GLFW_KEY_N
	KeyO            Key = C.GLFW_KEY_O
	KeyP            Key = C.GLFW_KEY_P
	KeyQ            Key = C.GLFW_KEY_Q
	KeyR            Key = C.GLFW_KEY_R
	KeyS            Key = C.GLFW_KEY_S
	KeyT            Key = C.GLFW_KEY_T
	KeyU            Key = C.GLFW_KEY_U
	KeyV            Key = C.GLFW_KEY_V
	KeyW            Key = C.GLFW_KEY_W
	KeyX            Key = C.GLFW_KEY_X
	KeyY            Key = C.GLFW_KEY_Y
	KeyZ            Key = C.GLFW_KEY_Z
	KeyLeftBracket  Key = C.GLFW_KEY_LEFT_BRACKET
	KeyBackslash    Key = C.GLFW_KEY_BACKSLASH
	KeyRightBracket Key = C.GLFW_KEY_RIGHT_BRACKET
	KeyGraveAccent  Key = C.GLFW_KEY_GRAVE_ACCENT
	KeyWorld1       Key = C.GLFW_KEY_WORLD_1
	KeyWorld2       Key = C.GLFW_KEY_WORLD_2
	KeyEscape       Key = C.GLFW_KEY_ESCAPE
	KeyEnter        Key = C.GLFW_KEY_ENTER
	KeyTab          Key = C.GLFW_KEY_TAB
	KeyBackspace    Key = C.GLFW_KEY_BACKSPACE
	KeyInsert       Key = C.GLFW_KEY_INSERT
	KeyDelete       Key = C.GLFW_KEY_DELETE
	KeyRight        Key = C.GLFW_KEY_RIGHT
	KeyLeft         Key = C.GLFW_KEY_LEFT
	KeyDown         Key = C.GLFW_KEY_DOWN
	KeyUp           Key = C.GLFW_KEY_UP
	KeyPageUp       Key = C.GLFW_KEY_PAGE_UP
	KeyPageDown     Key = C.GLFW_KEY_PAGE_DOWN
	KeyHome         Key = C.GLFW_KEY_HOME
	KeyEnd          Key = C.GLFW_KEY_END
	KeyCapsLock     Key = C.GLFW_KEY_CAPS_LOCK
	KeyScrollLock   Key = C.GLFW_KEY_SCROLL_LOCK
	KeyNumLock      Key = C.GLFW_KEY_NUM_LOCK
	KeyPrintScreen  Key = C.GLFW_KEY_PRINT_SCREEN
	KeyPause        Key = C.GLFW_KEY_PAUSE
	KeyF1           Key = C.GLFW_KEY_F1
	KeyF2           Key = C.GLFW_KEY_F2
	KeyF3           Key = C.GLFW_KEY_F3
	KeyF4           Key = C.GLFW_KEY_F4
	KeyF5           Key = C.GLFW_KEY_F5
	KeyF6           Key = C.GLFW_KEY_F6
	KeyF7           Key = C.GLFW_KEY_F7
	KeyF8           Key = C.GLFW_KEY_F8
	KeyF9           Key = C.GLFW_KEY_F9
	KeyF10          Key = C.GLFW_KEY_F10
	KeyF11          Key = C.GLFW_KEY_F11
	KeyF12          Key = C.GLFW_KEY_F12
	KeyF13          Key = C.GLFW_KEY_F13
	KeyF14          Key = C.GLFW_KEY_F14
	KeyF15          Key = C.GLFW_KEY_F15
	KeyF16          Key = C.GLFW_KEY_F16
	KeyF17          Key = C.GLFW_KEY_F17
	KeyF18          Key = C.GLFW_KEY_F18
	KeyF19          Key = C.GLFW_KEY_F19
	KeyF20          Key = C.GLFW_KEY_F20
	KeyF21          Key = C.GLFW_KEY_F21
	KeyF22          Key = C.GLFW_KEY_F22
	KeyF23          Key = C.GLFW_KEY_F23
	KeyF24          Key = C.GLFW_KEY_F24
	KeyF25          Key = C.GLFW_KEY_F25
	KeyKP0          Key = C.GLFW_KEY_KP_0
	KeyKP1          Key = C.GLFW_KEY_KP_1
	KeyKP2          Key = C.GLFW_KEY_KP_2
	KeyKP3          Key = C.GLFW_KEY_KP_3
	KeyKP4          Key = C.GLFW_KEY_KP_4
	KeyKP5          Key = C.GLFW_KEY_KP_5
	KeyKP6          Key = C.GLFW_KEY_KP_6
	KeyKP7          Key = C.GLFW_KEY_KP_7
	KeyKP8          Key = C.GLFW_KEY_KP_8
	KeyKP9          Key = C.GLFW_KEY_KP_9
	KeyKPDecimal    Key = C.GLFW_KEY_KP_DECIMAL
	KeyKPDivide     Key = C.GLFW_KEY_KP_DIVIDE
	KeyKPMultiply   Key = C.GLFW_KEY_KP_MULTIPLY
	KeyKPSubtract   Key = C.GLFW_KEY_KP_SUBTRACT
	KeyKPAdd        Key = C.GLFW_KEY_KP_ADD
	KeyKPEnter      Key = C.GLFW_KEY_KP_ENTER
	KeyKPEqual      Key = C.GLFW_KEY_KP_EQUAL
	KeyLeftShift    Key = C.GLFW_KEY_LEFT_SHIFT
	KeyLeftControl  Key = C.GLFW_KEY_LEFT_CONTROL
	KeyLeftAlt      Key = C.GLFW_KEY_LEFT_ALT
	KeyLeftSuper    Key = C.GLFW_KEY_LEFT_SUPER
	KeyRightShift   Key = C.GLFW_KEY_RIGHT_SHIFT
	KeyRightControl Key = C.GLFW_KEY_RIGHT_CONTROL
	KeyRightAlt     Key = C.GLFW_KEY_RIGHT_ALT
	KeyRightSuper   Key = C.GLFW_KEY_RIGHT_SUPER
	KeyMenu         Key = C.GLFW_KEY_MENU
	KeyLast         Key = C.GLFW_KEY_LAST
)

// ModifierKey corresponds to a modifier key.
type ModifierKey int

// Modifier keys.
const (
	ModShift   ModifierKey = C.GLFW_MOD_SHIFT
	ModControl ModifierKey = C.GLFW_MOD_CONTROL
	ModAlt     ModifierKey = C.GLFW_MOD_ALT
	ModSuper   ModifierKey = C.GLFW_MOD_SUPER
)

// MouseButton corresponds to a mouse button.
type MouseButton int

// Mouse buttons.
const (
	MouseButton1      MouseButton = C.GLFW_MOUSE_BUTTON_1
	MouseButton2      MouseButton = C.GLFW_MOUSE_BUTTON_2
	MouseButton3      MouseButton = C.GLFW_MOUSE_BUTTON_3
	MouseButton4      MouseButton = C.GLFW_MOUSE_BUTTON_4
	MouseButton5      MouseButton = C.GLFW_MOUSE_BUTTON_5
	MouseButton6      MouseButton = C.GLFW_MOUSE_BUTTON_6
	MouseButton7      MouseButton = C.GLFW_MOUSE_BUTTON_7
	MouseButton8      MouseButton = C.GLFW_MOUSE_BUTTON_8
	MouseButtonLast   MouseButton = C.GLFW_MOUSE_BUTTON_LAST
	MouseButtonLeft   MouseButton = C.GLFW_MOUSE_BUTTON_LEFT
	MouseButtonRight  MouseButton = C.GLFW_MOUSE_BUTTON_RIGHT
	MouseButtonMiddle MouseButton = C.GLFW_MOUSE_BUTTON_MIDDLE
)

// StandardCursor corresponds to a standard cursor icon.
type StandardCursor int

// Standard cursors
const (
	ArrowCursor     StandardCursor = C.GLFW_ARROW_CURSOR
	IBeamCursor     StandardCursor = C.GLFW_IBEAM_CURSOR
	CrosshairCursor StandardCursor = C.GLFW_CROSSHAIR_CURSOR
	HandCursor      StandardCursor = C.GLFW_HAND_CURSOR
	HResizeCursor   StandardCursor = C.GLFW_HRESIZE_CURSOR
	VResizeCursor   StandardCursor = C.GLFW_VRESIZE_CURSOR
)

// Action corresponds to a key or button action.
type Action int

// Action types.
const (
	Release Action = C.GLFW_RELEASE // The key or button was released.
	Press   Action = C.GLFW_PRESS   // The key or button was pressed.
	Repeat  Action = C.GLFW_REPEAT  // The key was held down until it repeated.
)

// InputMode corresponds to an input mode.
type InputMode int

// Input modes.
const (
	CursorMode             InputMode = C.GLFW_CURSOR               // See Cursor mode values
	StickyKeysMode         InputMode = C.GLFW_STICKY_KEYS          // Value can be either 1 or 0
	StickyMouseButtonsMode InputMode = C.GLFW_STICKY_MOUSE_BUTTONS // Value can be either 1 or 0
)

// Cursor mode values.
const (
	CursorNormal   int = C.GLFW_CURSOR_NORMAL
	CursorHidden   int = C.GLFW_CURSOR_HIDDEN
	CursorDisabled int = C.GLFW_CURSOR_DISABLED
)

// Cursor represents a cursor.
type Cursor struct {
	data *C.GLFWcursor
}

//export goJoystickCB
func goJoystickCB(joy, event C.int) {
	fJoystickHolder(int(joy), int(event))
}

//export goMouseButtonCB
func goMouseButtonCB(window unsafe.Pointer, button, action, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fMouseButtonHolder(w, MouseButton(button), Action(action), ModifierKey(mods))
}

//export goCursorPosCB
func goCursorPosCB(window unsafe.Pointer, xpos, ypos C.double) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCursorPosHolder(w, float64(xpos), float64(ypos))
}

//export goCursorEnterCB
func goCursorEnterCB(window unsafe.Pointer, entered C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	hasEntered := glfwbool(entered)
	w.fCursorEnterHolder(w, hasEntered)
}

//export goScrollCB
func goScrollCB(window unsafe.Pointer, xoff, yoff C.double) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fScrollHolder(w, float64(xoff), float64(yoff))
}

//export goKeyCB
func goKeyCB(window unsafe.Pointer, key, scancode, action, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fKeyHolder(w, Key(key), int(scancode), Action(action), ModifierKey(mods))
}

//export goCharCB
func goCharCB(window unsafe.Pointer, character C.uint) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCharHolder(w, rune(character))
}

//export goCharModsCB
func goCharModsCB(window unsafe.Pointer, character C.uint, mods C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCharModsHolder(w, rune(character), ModifierKey(mods))
}

//export goDropCB
func goDropCB(window unsafe.Pointer, count C.int, names **C.char) { // TODO: The types of name can be `**C.char` or `unsafe.Pointer`, use whichever is better.
	w := windows.get((*C.GLFWwindow)(window))
	namesSlice := make([]string, int(count)) // TODO: Make this better. This part is unfinished, hacky, probably not correct, and not idiomatic.
	for i := 0; i < int(count); i++ {        // TODO: Make this better. It should be cleaned up and vetted.
		var x *C.char                                                                                 // TODO: Make this better.
		p := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(names)) + uintptr(i)*unsafe.Sizeof(x))) // TODO: Make this better.
		namesSlice[i] = C.GoString(*p)                                                                // TODO: Make this better.
	}
	w.fDropHolder(w, namesSlice)
}

// GetInputMode returns the value of an input option of the window.
func (w *Window) GetInputMode(mode InputMode) int {
	ret := int(C.glfwGetInputMode(w.data, C.int(mode)))
	panicError()
	return ret
}

// SetInputMode sets an input option for the window.
func (w *Window) SetInputMode(mode InputMode, value int) {
	C.glfwSetInputMode(w.data, C.int(mode), C.int(value))
	panicError()
}

// GetKey returns the last reported state of a keyboard key. The returned state
// is one of Press or Release. The higher-level state Repeat is only reported to
// the key callback.
//
// If the StickyKeys input mode is enabled, this function returns Press the first
// time you call this function after a key has been pressed, even if the key has
// already been released.
//
// The key functions deal with physical keys, with key tokens named after their
// use on the standard US keyboard layout. If you want to input text, use the
// Unicode character callback instead.
func (w *Window) GetKey(key Key) Action {
	ret := Action(C.glfwGetKey(w.data, C.int(key)))
	panicError()
	return ret
}

// GetKeyName returns the localized name of the specified printable key.
//
// If the key is glfw.KeyUnknown, the scancode is used, otherwise the scancode is ignored.
func GetKeyName(key Key, scancode int) string {
	ret := C.glfwGetKeyName(C.int(key), C.int(scancode))
	panicError()
	return C.GoString(ret)
}

// GetMouseButton returns the last state reported for the specified mouse button.
//
// If the StickyMouseButtons input mode is enabled, this function returns Press
// the first time you call this function after a mouse button has been pressed,
// even if the mouse button has already been released.
func (w *Window) GetMouseButton(button MouseButton) Action {
	ret := Action(C.glfwGetMouseButton(w.data, C.int(button)))
	panicError()
	return ret
}

// GetCursorPos returns the last reported position of the cursor.
//
// If the cursor is disabled (with CursorDisabled) then the cursor position is
// unbounded and limited only by the minimum and maximum values of a double.
//
// The coordinate can be converted to their integer equivalents with the floor
// function. Casting directly to an integer type works for positive coordinates,
// but fails for negative ones.
func (w *Window) GetCursorPos() (x, y float64) {
	var xpos, ypos C.double
	C.glfwGetCursorPos(w.data, &xpos, &ypos)
	panicError()
	return float64(xpos), float64(ypos)
}

// SetCursorPos sets the position of the cursor. The specified window must
// be focused. If the window does not have focus when this function is called,
// it fails silently.
//
// If the cursor is disabled (with CursorDisabled) then the cursor position is
// unbounded and limited only by the minimum and maximum values of a double.
func (w *Window) SetCursorPos(xpos, ypos float64) {
	C.glfwSetCursorPos(w.data, C.double(xpos), C.double(ypos))
	panicError()
}

// CreateCursor creates a new custom cursor image that can be set for a window with SetCursor.
// The cursor can be destroyed with Destroy. Any remaining cursors are destroyed by Terminate.
//
// The image is ideally provided in the form of *image.NRGBA.
// The pixels are 32-bit, little-endian, non-premultiplied RGBA, i.e. eight
// bits per channel with the red channel first. They are arranged canonically
// as packed sequential rows, starting from the top-left corner. If the image
// type is not *image.NRGBA, it will be converted to it.
//
// The cursor hotspot is specified in pixels, relative to the upper-left corner of the cursor image.
// Like all other coordinate systems in GLFW, the X-axis points to the right and the Y-axis points down.
func CreateCursor(img image.Image, xhot, yhot int) *Cursor {
	var imgC C.GLFWimage
	var pixels []uint8
	b := img.Bounds()

	switch img := img.(type) {
	case *image.NRGBA:
		pixels = img.Pix
	default:
		m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
		pixels = m.Pix
	}

	pix, free := bytes(pixels)

	imgC.width = C.int(b.Dx())
	imgC.height = C.int(b.Dy())
	imgC.pixels = (*C.uchar)(pix)

	c := C.glfwCreateCursor(&imgC, C.int(xhot), C.int(yhot))

	free()
	panicError()

	return &Cursor{c}
}

// CreateStandardCursor returns a cursor with a standard shape,
// that can be set for a window with SetCursor.
func CreateStandardCursor(shape StandardCursor) *Cursor {
	c := C.glfwCreateStandardCursor(C.int(shape))
	panicError()
	return &Cursor{c}
}

// Destroy destroys a cursor previously created with CreateCursor.
// Any remaining cursors will be destroyed by Terminate.
func (c *Cursor) Destroy() {
	C.glfwDestroyCursor(c.data)
	panicError()
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
// of the specified window. The set cursor will only be visible when the cursor mode of the
// window is CursorNormal.
//
// On some platforms, the set cursor may not be visible unless the window also has input focus.
func (w *Window) SetCursor(c *Cursor) {
	if c == nil {
		C.glfwSetCursor(w.data, nil)
	} else {
		C.glfwSetCursor(w.data, c.data)
	}
	panicError()
}

// JoystickCallback is the joystick configuration callback.
type JoystickCallback func(joy, event int)

// SetJoystickCallback sets the joystick configuration callback, or removes the
// currently set callback. This is called when a joystick is connected to or
// disconnected from the system.
func SetJoystickCallback(cbfun JoystickCallback) (previous JoystickCallback) {
	previous = fJoystickHolder
	fJoystickHolder = cbfun
	if cbfun == nil {
		C.glfwSetJoystickCallback(nil)
	} else {
		C.glfwSetJoystickCallbackCB()
	}
	panicError()
	return previous
}

// KeyCallback is the key callback.
type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

// SetKeyCallback sets the key callback which is called when a key is pressed,
// repeated or released.
//
// The key functions deal with physical keys, with layout independent key tokens
// named after their values in the standard US keyboard layout. If you want to
// input text, use the SetCharCallback instead.
//
// When a window loses focus, it will generate synthetic key release events for
// all pressed keys. You can tell these events from user-generated events by the
// fact that the synthetic ones are generated after the window has lost focus,
// i.e. Focused will be false and the focus callback will have already been
// called.
func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	previous = w.fKeyHolder
	w.fKeyHolder = cbfun
	if cbfun == nil {
		C.glfwSetKeyCallback(w.data, nil)
	} else {
		C.glfwSetKeyCallbackCB(w.data)
	}
	panicError()
	return previous
}

// CharCallback is the character callback.
type CharCallback func(w *Window, char rune)

// SetCharCallback sets the character callback which is called when a
// Unicode character is input.
//
// The character callback is intended for Unicode text input. As it deals with
// characters, it is keyboard layout dependent, whereas the
// key callback is not. Characters do not map 1:1
// to physical keys, as a key may produce zero, one or more characters. If you
// want to know whether a specific physical key was pressed or released, see
// the key callback instead.
//
// The character callback behaves as system text input normally does and will
// not be called if modifier keys are held down that would prevent normal text
// input on that platform, for example a Super (Command) key on OS X or Alt key
// on Windows. There is a character with modifiers callback that receives these events.
func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	previous = w.fCharHolder
	w.fCharHolder = cbfun
	if cbfun == nil {
		C.glfwSetCharCallback(w.data, nil)
	} else {
		C.glfwSetCharCallbackCB(w.data)
	}
	panicError()
	return previous
}

// CharModsCallback is the character with modifiers callback.
type CharModsCallback func(w *Window, char rune, mods ModifierKey)

// SetCharModsCallback sets the character with modifiers callback which is called when a
// Unicode character is input regardless of what modifier keys are used.
//
// The character with modifiers callback is intended for implementing custom
// Unicode character input. For regular Unicode text input, see the
// character callback. Like the character callback, the character with modifiers callback
// deals with characters and is keyboard layout dependent. Characters do not
// map 1:1 to physical keys, as a key may produce zero, one or more characters.
// If you want to know whether a specific physical key was pressed or released,
// see the key callback instead.
func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	previous = w.fCharModsHolder
	w.fCharModsHolder = cbfun
	if cbfun == nil {
		C.glfwSetCharModsCallback(w.data, nil)
	} else {
		C.glfwSetCharModsCallbackCB(w.data)
	}
	panicError()
	return previous
}

// MouseButtonCallback is the mouse button callback.
type MouseButtonCallback func(w *Window, button MouseButton, action Action, mod ModifierKey)

// SetMouseButtonCallback sets the mouse button callback which is called when a
// mouse button is pressed or released.
//
// When a window loses focus, it will generate synthetic mouse button release
// events for all pressed mouse buttons. You can tell these events from
// user-generated events by the fact that the synthetic ones are generated after
// the window has lost focus, i.e. Focused will be false and the focus
// callback will have already been called.
func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	previous = w.fMouseButtonHolder
	w.fMouseButtonHolder = cbfun
	if cbfun == nil {
		C.glfwSetMouseButtonCallback(w.data, nil)
	} else {
		C.glfwSetMouseButtonCallbackCB(w.data)
	}
	panicError()
	return previous
}

// CursorPosCallback the cursor position callback.
type CursorPosCallback func(w *Window, xpos float64, ypos float64)

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the window.
func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	previous = w.fCursorPosHolder
	w.fCursorPosHolder = cbfun
	if cbfun == nil {
		C.glfwSetCursorPosCallback(w.data, nil)
	} else {
		C.glfwSetCursorPosCallbackCB(w.data)
	}
	panicError()
	return previous
}

// CursorEnterCallback is the cursor boundary crossing callback.
type CursorEnterCallback func(w *Window, entered bool)

// SetCursorEnterCallback the cursor boundary crossing callback which is called
// when the cursor enters or leaves the client area of the window.
func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {
	previous = w.fCursorEnterHolder
	w.fCursorEnterHolder = cbfun
	if cbfun == nil {
		C.glfwSetCursorEnterCallback(w.data, nil)
	} else {
		C.glfwSetCursorEnterCallbackCB(w.data)
	}
	panicError()
	return previous
}

// ScrollCallback is the scroll callback.
type ScrollCallback func(w *Window, xoff float64, yoff float64)

// SetScrollCallback sets the scroll callback which is called when a scrolling
// device is used, such as a mouse wheel or scrolling area of a touchpad.
func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	previous = w.fScrollHolder
	w.fScrollHolder = cbfun
	if cbfun == nil {
		C.glfwSetScrollCallback(w.data, nil)
	} else {
		C.glfwSetScrollCallbackCB(w.data)
	}
	panicError()
	return previous
}

// DropCallback is the drop callback.
type DropCallback func(w *Window, names []string)

// SetDropCallback sets the drop callback which is called when an object
// is dropped over the window.
func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	previous = w.fDropHolder
	w.fDropHolder = cbfun
	if cbfun == nil {
		C.glfwSetDropCallback(w.data, nil)
	} else {
		C.glfwSetDropCallbackCB(w.data)
	}
	panicError()
	return previous
}

// JoystickPresent reports whether the specified joystick is present.
func JoystickPresent(joy Joystick) bool {
	ret := glfwbool(C.glfwJoystickPresent(C.int(joy)))
	panicError()
	return ret
}

// GetJoystickAxes returns a slice of axis values.
func GetJoystickAxes(joy Joystick) []float32 {
	var length int

	axis := C.glfwGetJoystickAxes(C.int(joy), (*C.int)(unsafe.Pointer(&length)))
	panicError()
	if axis == nil {
		return nil
	}

	a := make([]float32, length)
	for i := 0; i < length; i++ {
		a[i] = float32(C.GetAxisAtIndex(axis, C.int(i)))
	}
	return a
}

// GetJoystickButtons returns a slice of button values.
func GetJoystickButtons(joy Joystick) []byte {
	var length int

	buttons := C.glfwGetJoystickButtons(C.int(joy), (*C.int)(unsafe.Pointer(&length)))
	panicError()
	if buttons == nil {
		return nil
	}

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte(C.GetButtonsAtIndex(buttons, C.int(i)))
	}
	return b
}

// GetJoystickName returns the name, encoded as UTF-8, of the specified joystick.
func GetJoystickName(joy Joystick) string {
	jn := C.glfwGetJoystickName(C.int(joy))
	panicError()
	return C.GoString(jn)
}
