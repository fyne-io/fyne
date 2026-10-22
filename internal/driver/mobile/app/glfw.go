//go:build (linux && !android) || freebsd || openbsd

package app

/*
Simple on-screen app debugging for X11 and Wayland, using the same GLFW
bindings as the desktop driver. Not an officially supported development target
for apps, as screens with mice are very different than screens with touch
panels.
*/

import (
	"runtime"
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/lifecycle"
	"fyne.io/fyne/v2/internal/driver/mobile/event/paint"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// The simulated device is a tall, narrow screen so that layouts are checked in
// the orientation that they are most likely to be broken in.
const (
	windowWidth  = 500
	windowHeight = 1000
)

const (
	colorBits = 8
	depthBits = 16
)

var stopped bool

func init() {
	theApp.registerGLViewportFilter()
}

// glfwKeyCode maps the keys that GLFW does not number by their character to the
// USB HID key codes that the mobile key package uses.
//
// Printable keys are not listed, glfwKeyToCode resolves them with getCodeFromRune
// so that they share the mapping of the soft keyboard.
var glfwKeyCode = map[glfw.Key]key.Code{
	// editing and navigation
	glfw.KeyBackspace: key.CodeDeleteBackspace,
	glfw.KeyDelete:    key.CodeDeleteForward,
	glfw.KeyDown:      key.CodeDownArrow,
	glfw.KeyEnd:       key.CodeEnd,
	glfw.KeyEnter:     key.CodeReturnEnter,
	glfw.KeyEscape:    key.CodeEscape,
	glfw.KeyHome:      key.CodeHome,
	glfw.KeyInsert:    key.CodeInsert,
	glfw.KeyLeft:      key.CodeLeftArrow,
	glfw.KeyPageDown:  key.CodePageDown,
	glfw.KeyPageUp:    key.CodePageUp,
	glfw.KeyRight:     key.CodeRightArrow,
	glfw.KeyTab:       key.CodeTab,
	glfw.KeyUp:        key.CodeUpArrow,

	// keypad
	glfw.KeyKP0:        key.CodeKeypad0,
	glfw.KeyKP1:        key.CodeKeypad1,
	glfw.KeyKP2:        key.CodeKeypad2,
	glfw.KeyKP3:        key.CodeKeypad3,
	glfw.KeyKP4:        key.CodeKeypad4,
	glfw.KeyKP5:        key.CodeKeypad5,
	glfw.KeyKP6:        key.CodeKeypad6,
	glfw.KeyKP7:        key.CodeKeypad7,
	glfw.KeyKP8:        key.CodeKeypad8,
	glfw.KeyKP9:        key.CodeKeypad9,
	glfw.KeyKPDecimal:  key.CodeKeypadFullStop,
	glfw.KeyKPEnter:    key.CodeKeypadEnter,
	glfw.KeyKPSubtract: key.CodeKeypadHyphenMinus,

	// modifiers
	glfw.KeyLeftAlt:      key.CodeLeftAlt,
	glfw.KeyLeftControl:  key.CodeLeftControl,
	glfw.KeyLeftShift:    key.CodeLeftShift,
	glfw.KeyLeftSuper:    key.CodeLeftGUI,
	glfw.KeyRightAlt:     key.CodeRightAlt,
	glfw.KeyRightControl: key.CodeRightControl,
	glfw.KeyRightShift:   key.CodeRightShift,
	glfw.KeyRightSuper:   key.CodeRightGUI,
	glfw.KeyMenu:         key.CodeCompose,

	// function keys, fyne has no names beyond F12
	glfw.KeyF1:  key.CodeF1,
	glfw.KeyF2:  key.CodeF2,
	glfw.KeyF3:  key.CodeF3,
	glfw.KeyF4:  key.CodeF4,
	glfw.KeyF5:  key.CodeF5,
	glfw.KeyF6:  key.CodeF6,
	glfw.KeyF7:  key.CodeF7,
	glfw.KeyF8:  key.CodeF8,
	glfw.KeyF9:  key.CodeF9,
	glfw.KeyF10: key.CodeF10,
	glfw.KeyF11: key.CodeF11,
	glfw.KeyF12: key.CodeF12,
}

// glfwKeyToCode returns the USB HID code of a key, or key.CodeUnknown if the key
// package has no equivalent for it.
func glfwKeyToCode(k glfw.Key) key.Code {
	if code, ok := glfwKeyCode[k]; ok {
		return code
	}

	// GLFW numbers the printable keys with their ASCII code point, and the codes
	// of those are in the table that the soft keyboard uses too.
	if k > 0 && k < utf8.RuneSelf {
		return getCodeFromRune(rune(k))
	}

	return key.CodeUnknown
}

func GoBack() {
	// When simulating mobile there are no other activities open (and we can't just force background)
}

// Main is called by the main.main function to run the mobile application.
//
// It calls f on the App, in a separate goroutine, as some OS-specific
// libraries require being on 'the main thread'.
func Main(f func(App)) {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		fyne.LogError("failed to initialise GLFW", err)
		return
	}
	defer glfw.Terminate()

	window, err := createGLFWWindow()
	if err != nil {
		fyne.LogError("failed to create window for the mobile simulation", err)
		return
	}
	defer window.Destroy()

	workAvailable := theApp.worker.WorkAvailable()
	heartbeat := time.NewTicker(time.Second / 60)

	// TODO: send lifecycle events when e.g. the window is iconified or moved off-screen.
	theApp.sendLifecycle(lifecycle.StageFocused)

	// TODO: send a paint event from the GLFW refresh callback instead of
	// sending this synthetic one as a hack.
	theApp.events.In() <- paint.Event{}

	donec := make(chan struct{})
	go func() {
		f(theApp)
		close(donec)
	}()

	// TODO: can we get the actual vsync signal?
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	var tc <-chan time.Time

	for {
		select {
		case <-donec:
			return
		case <-heartbeat.C:
			glfw.PollEvents()
			if window.ShouldClose() {
				return
			}
		case <-workAvailable:
			theApp.worker.DoWork()
		case <-theApp.publish:
			window.SwapBuffers()
			tc = ticker.C
		case <-tc:
			tc = nil
			theApp.publishResult <- PublishResult{}
		}
	}
}

// createGLFWWindow opens the window that the simulated device is displayed in.
// It must be called on the main thread with GLFW already initialised, the
// returned window holds the OpenGL context that the app will draw into.
func createGLFWWindow() (*glfw.Window, error) {
	// The mobile driver draws with OpenGL ES, so an OpenGL ES context is
	// required here too, a desktop GL context would not have the entry points
	// that internal/driver/mobile/gl is compiled against.
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.AutoIconify, glfw.False)
	glfw.WindowHint(glfw.RedBits, colorBits)
	glfw.WindowHint(glfw.GreenBits, colorBits)
	glfw.WindowHint(glfw.BlueBits, colorBits)
	glfw.WindowHint(glfw.AlphaBits, colorBits)
	glfw.WindowHint(glfw.DepthBits, depthBits)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Fyne", nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	// The framebuffer can be a different size to the window, for example on a
	// scaled display, and it is the framebuffer that we draw into.
	window.SetFramebufferSizeCallback(onResize)
	window.SetCursorPosCallback(mouseMoved)
	window.SetMouseButtonCallback(mouseClicked)
	window.SetKeyCallback(keyPressed)
	window.SetCharCallback(charInput)
	window.SetCloseCallback(closed)

	width, height := window.GetFramebufferSize()
	onResize(window, width, height)

	return window, nil
}

func onResize(_ *glfw.Window, w, h int) {
	// TODO: don't assume 72 DPI. DisplayWidth and DisplayWidthMM
	// is probably the best place to start looking.
	pixelsPerPt := float32(1)
	theApp.events.In() <- size.Event{
		WidthPx:     w,
		HeightPx:    h,
		WidthPt:     float32(w),
		HeightPt:    float32(h),
		PixelsPerPt: pixelsPerPt,
		Orientation: screenOrientation(w, h),
	}
}

func sendTouch(t touch.Type, x, y float32) {
	theApp.events.In() <- touch.Event{
		X:        x,
		Y:        y,
		Sequence: 0, // TODO: button??
		Type:     t,
	}
}

// mouseDown tracks whether a button is held, so that moving the pointer without
// pressing does not start a drag.
var mouseDown bool

func mouseClicked(w *glfw.Window, _ glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	x, y := w.GetCursorPos()

	switch action {
	case glfw.Press:
		mouseDown = true
		sendTouch(touch.TypeBegin, float32(x), float32(y))
	case glfw.Release:
		mouseDown = false
		sendTouch(touch.TypeEnd, float32(x), float32(y))
	}
}

func mouseMoved(_ *glfw.Window, x, y float64) {
	if !mouseDown {
		return
	}
	sendTouch(touch.TypeMove, float32(x), float32(y))
}

// currentModifiers is the modifier state reported to the last key event, it is
// used for the characters that follow it because GLFW deprecated the callback
// that reports the modifiers of a character input.
var currentModifiers key.Modifiers

func keyPressed(_ *glfw.Window, k glfw.Key, _ int, action glfw.Action, mods glfw.ModifierKey) {
	if stopped {
		return
	}
	currentModifiers = modifiersCorrected(mods, k, action)

	// A repeat is reported as a press so that holding a key, such as
	// backspace, will keep working.
	dir := key.DirPress
	if action == glfw.Release {
		dir = key.DirRelease
	}

	// A physical key does not always produce a character, and the character
	// callback reports the ones that it does, so the rune stays unset here.
	theApp.events.In() <- key.Event{
		Rune:      -1,
		Code:      glfwKeyToCode(k),
		Modifiers: currentModifiers,
		Direction: dir,
	}
}

// charInput is called when a key produces a Unicode character, the modifiers of
// the key press that started it are used as GLFW no longer reports them here.
func charInput(_ *glfw.Window, char rune) {
	if stopped {
		return
	}

	theApp.events.In() <- key.Event{
		Rune:      char,
		Code:      key.CodeUnknown,
		Modifiers: currentModifiers,
		Direction: key.DirPress,
	}
}

func closed(window *glfw.Window) {
	window.SetShouldClose(true)
	onStop()
}

func onStop() {
	if stopped {
		return
	}
	stopped = true
	theApp.sendLifecycle(lifecycle.StageDead)
	theApp.events.Close()
}

func modifiersCorrected(mods glfw.ModifierKey, k glfw.Key, action glfw.Action) key.Modifiers {
	// On X11 pressing/releasing a modifier key does not include the newly
	// pressed/released key in the modifier mask.
	// See https://github.com/glfw/glfw/issues/1630
	switch action {
	case glfw.Press, glfw.Repeat:
		mods |= glfwKeyToModifier(k)
	case glfw.Release:
		mods &= ^glfwKeyToModifier(k)
	}

	return modifiers(mods)
}

func modifiers(mods glfw.ModifierKey) key.Modifiers {
	var m key.Modifiers
	if mods&glfw.ModShift != 0 {
		m |= key.ModShift
	}
	if mods&glfw.ModControl != 0 {
		m |= key.ModControl
	}
	if mods&glfw.ModAlt != 0 {
		m |= key.ModAlt
	}
	if mods&glfw.ModSuper != 0 {
		m |= key.ModMeta
	}
	return m
}

func glfwKeyToModifier(k glfw.Key) glfw.ModifierKey {
	switch k {
	case glfw.KeyLeftControl, glfw.KeyRightControl:
		return glfw.ModControl
	case glfw.KeyLeftAlt, glfw.KeyRightAlt:
		return glfw.ModAlt
	case glfw.KeyLeftShift, glfw.KeyRightShift:
		return glfw.ModShift
	case glfw.KeyLeftSuper, glfw.KeyRightSuper:
		return glfw.ModSuper
	}

	return 0
}

// driverShowVirtualKeyboard does nothing on desktop
func driverShowVirtualKeyboard(KeyboardType) {
}

// driverHideVirtualKeyboard does nothing on desktop
func driverHideVirtualKeyboard() {
}

// driverShowFileOpenPicker does nothing on desktop
func driverShowFileOpenPicker(func(string, func()), *FileFilter) {
}

// driverShowFileSavePicker does nothing on desktop
func driverShowFileSavePicker(func(string, func()), *FileFilter, string) {
}
