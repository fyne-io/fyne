// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios
// +build darwin,!ios

package app

// Simple on-screen app debugging for OS X. Not an officially supported
// development target for apps, as screens with mice are very different
// than screens with touch panels.

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>
#include <pthread.h>

void runApp(void);
void stopApp(void);
void makeCurrentContext(GLintptr);
uint64 threadID();
*/
import "C"
import (
	"log"
	"runtime"

	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/lifecycle"
	"fyne.io/fyne/v2/internal/driver/mobile/event/paint"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
)

var initThreadID uint64

func init() {
	// Lock the goroutine responsible for initialization to an OS thread.
	// This means the goroutine running main (and calling runApp below)
	// is locked to the OS thread that started the program. This is
	// necessary for the correct delivery of Cocoa events to the process.
	//
	// A discussion on this topic:
	// https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()
	initThreadID = uint64(C.threadID())
}

func main(f func(App)) {
	if tid := uint64(C.threadID()); tid != initThreadID {
		log.Fatalf("app.Main called on thread %d, but app.init ran on %d", tid, initThreadID)
	}

	go func() {
		f(theApp)
		C.stopApp()
		// TODO(crawshaw): trigger runApp to return
	}()

	C.runApp()
}

func GoBack() {
	// When simulating mobile there are no other activities open (and we can't just force background)
}

// loop is the primary drawing loop.
//
// After Cocoa has captured the initial OS thread for processing Cocoa
// events in runApp, it starts loop on another goroutine. It is locked
// to an OS thread for its OpenGL context.
//
// The loop processes GL calls until a publish event appears.
// Then it runs any remaining GL calls and flushes the screen.
//
// As NSOpenGLCPSwapInterval is set to 1, the call to CGLFlushDrawable
// blocks until the screen refresh.
func (a *app) loop(ctx C.GLintptr) {
	runtime.LockOSThread()
	C.makeCurrentContext(ctx)

	workAvailable := a.worker.WorkAvailable()

	for {
		select {
		case <-workAvailable:
			a.worker.DoWork()
		case <-theApp.publish:
		loop1:
			for {
				select {
				case <-workAvailable:
					a.worker.DoWork()
				default:
					break loop1
				}
			}
			C.CGLFlushDrawable(C.CGLGetCurrentContext())
			theApp.publishResult <- PublishResult{}
			select {
			case drawDone <- struct{}{}:
			default:
			}
		}
	}
}

var drawDone = make(chan struct{})

// drawgl is used by Cocoa to occasionally request screen updates.
//
//export drawgl
func drawgl() {
	switch theApp.lifecycleStage {
	case lifecycle.StageFocused, lifecycle.StageVisible:
		theApp.Send(paint.Event{
			External: true,
		})
		<-drawDone
	}
}

//export startloop
func startloop(ctx C.GLintptr) {
	go theApp.loop(ctx)
}

var windowHeightPx float32

//export setGeom
func setGeom(pixelsPerPt float32, widthPx, heightPx int) {
	windowHeightPx = float32(heightPx)
	theApp.events.In() <- size.Event{
		WidthPx:     widthPx,
		HeightPx:    heightPx,
		WidthPt:     float32(widthPx) / pixelsPerPt,
		HeightPt:    float32(heightPx) / pixelsPerPt,
		PixelsPerPt: pixelsPerPt,
		Orientation: screenOrientation(widthPx, heightPx),
	}
}

func sendTouch(t touch.Type, x, y float32) {
	theApp.events.In() <- touch.Event{
		X:        x,
		Y:        windowHeightPx - y,
		Sequence: 0,
		Type:     t,
	}
}

//export eventMouseDown
func eventMouseDown(x, y float32) { sendTouch(touch.TypeBegin, x, y) }

//export eventMouseDragged
func eventMouseDragged(x, y float32) { sendTouch(touch.TypeMove, x, y) }

//export eventMouseEnd
func eventMouseEnd(x, y float32) { sendTouch(touch.TypeEnd, x, y) }

//export lifecycleDead
func lifecycleDead() { theApp.sendLifecycle(lifecycle.StageDead) }

//export eventKey
func eventKey(runeVal int32, direction uint8, code uint16, flags uint32) {
	var modifiers key.Modifiers
	for _, mod := range mods {
		if flags&mod.flags == mod.flags {
			modifiers |= mod.mod
		}
	}

	theApp.events.In() <- key.Event{
		Rune:      convRune(rune(runeVal)),
		Code:      convVirtualKeyCode(code),
		Modifiers: modifiers,
		Direction: key.Direction(direction),
	}
}

//export eventFlags
func eventFlags(flags uint32) {
	for _, mod := range mods {
		if flags&mod.flags == mod.flags && lastFlags&mod.flags != mod.flags {
			eventKey(-1, uint8(key.DirPress), mod.code, flags)
		}
		if lastFlags&mod.flags == mod.flags && flags&mod.flags != mod.flags {
			eventKey(-1, uint8(key.DirRelease), mod.code, flags)
		}
	}
	lastFlags = flags
}

var lastFlags uint32

var mods = [...]struct {
	flags uint32
	code  uint16
	mod   key.Modifiers
}{
	// Left and right variants of modifier keys have their own masks,
	// but they are not documented. These were determined empirically.
	{1<<17 | 0x102, C.kVK_Shift, key.ModShift},
	{1<<17 | 0x104, C.kVK_RightShift, key.ModShift},
	{1<<18 | 0x101, C.kVK_Control, key.ModControl},
	// TODO key.ControlRight
	{1<<19 | 0x120, C.kVK_Option, key.ModAlt},
	{1<<19 | 0x140, C.kVK_RightOption, key.ModAlt},
	{1<<20 | 0x108, C.kVK_Command, key.ModMeta},
	{1<<20 | 0x110, C.kVK_Command, key.ModMeta}, // TODO: missing kVK_RightCommand
}

//export lifecycleAlive
func lifecycleAlive() { theApp.sendLifecycle(lifecycle.StageAlive) }

//export lifecycleVisible
func lifecycleVisible() {
	theApp.sendLifecycle(lifecycle.StageVisible)
}

//export lifecycleFocused
func lifecycleFocused() { theApp.sendLifecycle(lifecycle.StageFocused) }

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

// convRune marks the Carbon/Cocoa private-range unicode rune representing
// a non-unicode key event to -1, used for Rune in the key package.
//
// http://www.unicode.org/Public/MAPPINGS/VENDORS/APPLE/CORPCHAR.TXT
func convRune(r rune) rune {
	if '\uE000' <= r && r <= '\uF8FF' {
		return -1
	}
	return r
}

var virtualKeyCodeMap = map[uint16]key.Code{
	C.kVK_ANSI_A: key.CodeA,
	C.kVK_ANSI_B: key.CodeB,
	C.kVK_ANSI_C: key.CodeC,
	C.kVK_ANSI_D: key.CodeD,
	C.kVK_ANSI_E: key.CodeE,
	C.kVK_ANSI_F: key.CodeF,
	C.kVK_ANSI_G: key.CodeG,
	C.kVK_ANSI_H: key.CodeH,
	C.kVK_ANSI_I: key.CodeI,
	C.kVK_ANSI_J: key.CodeJ,
	C.kVK_ANSI_K: key.CodeK,
	C.kVK_ANSI_L: key.CodeL,
	C.kVK_ANSI_M: key.CodeM,
	C.kVK_ANSI_N: key.CodeN,
	C.kVK_ANSI_O: key.CodeO,
	C.kVK_ANSI_P: key.CodeP,
	C.kVK_ANSI_Q: key.CodeQ,
	C.kVK_ANSI_R: key.CodeR,
	C.kVK_ANSI_S: key.CodeS,
	C.kVK_ANSI_T: key.CodeT,
	C.kVK_ANSI_U: key.CodeU,
	C.kVK_ANSI_V: key.CodeV,
	C.kVK_ANSI_W: key.CodeW,
	C.kVK_ANSI_X: key.CodeX,
	C.kVK_ANSI_Y: key.CodeY,
	C.kVK_ANSI_Z: key.CodeZ,
	C.kVK_ANSI_1: key.Code1,
	C.kVK_ANSI_2: key.Code2,
	C.kVK_ANSI_3: key.Code3,
	C.kVK_ANSI_4: key.Code4,
	C.kVK_ANSI_5: key.Code5,
	C.kVK_ANSI_6: key.Code6,
	C.kVK_ANSI_7: key.Code7,
	C.kVK_ANSI_8: key.Code8,
	C.kVK_ANSI_9: key.Code9,
	C.kVK_ANSI_0: key.Code0,
	// TODO: move the rest of these codes to constants in key.go
	// if we are happy with them.
	C.kVK_Return:            key.CodeReturnEnter,
	C.kVK_Escape:            key.CodeEscape,
	C.kVK_Delete:            key.CodeDeleteBackspace,
	C.kVK_Tab:               key.CodeTab,
	C.kVK_Space:             key.CodeSpacebar,
	C.kVK_ANSI_Minus:        key.CodeHyphenMinus,
	C.kVK_ANSI_Equal:        key.CodeEqualSign,
	C.kVK_ANSI_LeftBracket:  key.CodeLeftSquareBracket,
	C.kVK_ANSI_RightBracket: key.CodeRightSquareBracket,
	C.kVK_ANSI_Backslash:    key.CodeBackslash,
	// 50: Keyboard Non-US "#" and ~
	C.kVK_ANSI_Semicolon: key.CodeSemicolon,
	C.kVK_ANSI_Quote:     key.CodeApostrophe,
	C.kVK_ANSI_Grave:     key.CodeGraveAccent,
	C.kVK_ANSI_Comma:     key.CodeComma,
	C.kVK_ANSI_Period:    key.CodeFullStop,
	C.kVK_ANSI_Slash:     key.CodeSlash,
	C.kVK_CapsLock:       key.CodeCapsLock,
	C.kVK_F1:             key.CodeF1,
	C.kVK_F2:             key.CodeF2,
	C.kVK_F3:             key.CodeF3,
	C.kVK_F4:             key.CodeF4,
	C.kVK_F5:             key.CodeF5,
	C.kVK_F6:             key.CodeF6,
	C.kVK_F7:             key.CodeF7,
	C.kVK_F8:             key.CodeF8,
	C.kVK_F9:             key.CodeF9,
	C.kVK_F10:            key.CodeF10,
	C.kVK_F11:            key.CodeF11,
	C.kVK_F12:            key.CodeF12,
	// 70: PrintScreen
	// 71: Scroll Lock
	// 72: Pause
	// 73: Insert
	C.kVK_Home:                key.CodeHome,
	C.kVK_PageUp:              key.CodePageUp,
	C.kVK_ForwardDelete:       key.CodeDeleteForward,
	C.kVK_End:                 key.CodeEnd,
	C.kVK_PageDown:            key.CodePageDown,
	C.kVK_RightArrow:          key.CodeRightArrow,
	C.kVK_LeftArrow:           key.CodeLeftArrow,
	C.kVK_DownArrow:           key.CodeDownArrow,
	C.kVK_UpArrow:             key.CodeUpArrow,
	C.kVK_ANSI_KeypadClear:    key.CodeKeypadNumLock,
	C.kVK_ANSI_KeypadDivide:   key.CodeKeypadSlash,
	C.kVK_ANSI_KeypadMultiply: key.CodeKeypadAsterisk,
	C.kVK_ANSI_KeypadMinus:    key.CodeKeypadHyphenMinus,
	C.kVK_ANSI_KeypadPlus:     key.CodeKeypadPlusSign,
	C.kVK_ANSI_KeypadEnter:    key.CodeKeypadEnter,
	C.kVK_ANSI_Keypad1:        key.CodeKeypad1,
	C.kVK_ANSI_Keypad2:        key.CodeKeypad2,
	C.kVK_ANSI_Keypad3:        key.CodeKeypad3,
	C.kVK_ANSI_Keypad4:        key.CodeKeypad4,
	C.kVK_ANSI_Keypad5:        key.CodeKeypad5,
	C.kVK_ANSI_Keypad6:        key.CodeKeypad6,
	C.kVK_ANSI_Keypad7:        key.CodeKeypad7,
	C.kVK_ANSI_Keypad8:        key.CodeKeypad8,
	C.kVK_ANSI_Keypad9:        key.CodeKeypad9,
	C.kVK_ANSI_Keypad0:        key.CodeKeypad0,
	C.kVK_ANSI_KeypadDecimal:  key.CodeKeypadFullStop,
	C.kVK_ANSI_KeypadEquals:   key.CodeKeypadEqualSign,
	C.kVK_F13:                 key.CodeF13,
	C.kVK_F14:                 key.CodeF14,
	C.kVK_F15:                 key.CodeF15,
	C.kVK_F16:                 key.CodeF16,
	C.kVK_F17:                 key.CodeF17,
	C.kVK_F18:                 key.CodeF18,
	C.kVK_F19:                 key.CodeF19,
	C.kVK_F20:                 key.CodeF20,
	// 116: Keyboard Execute
	C.kVK_Help: key.CodeHelp,
	// 118: Keyboard Menu
	// 119: Keyboard Select
	// 120: Keyboard Stop
	// 121: Keyboard Again
	// 122: Keyboard Undo
	// 123: Keyboard Cut
	// 124: Keyboard Copy
	// 125: Keyboard Paste
	// 126: Keyboard Find
	C.kVK_Mute:       key.CodeMute,
	C.kVK_VolumeUp:   key.CodeVolumeUp,
	C.kVK_VolumeDown: key.CodeVolumeDown,
	// 130: Keyboard Locking Caps Lock
	// 131: Keyboard Locking Num Lock
	// 132: Keyboard Locking Scroll Lock
	// 133: Keyboard Comma
	// 134: Keyboard Equal Sign
	// ...: Bunch of stuff
	C.kVK_Control:      key.CodeLeftControl,
	C.kVK_Shift:        key.CodeLeftShift,
	C.kVK_Option:       key.CodeLeftAlt,
	C.kVK_Command:      key.CodeLeftGUI,
	C.kVK_RightControl: key.CodeRightControl,
	C.kVK_RightShift:   key.CodeRightShift,
	C.kVK_RightOption:  key.CodeRightAlt,
}

// convVirtualKeyCode converts a Carbon/Cocoa virtual key code number
// into the standard keycodes used by the key package.
//
// To get a sense of the key map, see the diagram on
//
//	http://boredzo.org/blog/archives/2007-05-22/virtual-key-codes
func convVirtualKeyCode(vkcode uint16) key.Code {
	if code, ok := virtualKeyCodeMap[vkcode]; ok {
		return code
	}

	// TODO key.CodeRightGUI
	return key.CodeUnknown
}
