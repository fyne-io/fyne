// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build !ios

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
	"sync"

	"github.com/fyne-io/mobile/event/key"
	"github.com/fyne-io/mobile/event/lifecycle"
	"github.com/fyne-io/mobile/event/paint"
	"github.com/fyne-io/mobile/event/size"
	"github.com/fyne-io/mobile/event/touch"
	"github.com/fyne-io/mobile/geom"
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
	theApp.eventsIn <- size.Event{
		WidthPx:     widthPx,
		HeightPx:    heightPx,
		WidthPt:     geom.Pt(float32(widthPx) / pixelsPerPt),
		HeightPt:    geom.Pt(float32(heightPx) / pixelsPerPt),
		PixelsPerPt: pixelsPerPt,
		Orientation: screenOrientation(widthPx, heightPx),
	}
}

var touchEvents struct {
	sync.Mutex
	pending []touch.Event
}

func sendTouch(t touch.Type, x, y float32) {
	theApp.eventsIn <- touch.Event{
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

	theApp.eventsIn <- key.Event{
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

// convVirtualKeyCode converts a Carbon/Cocoa virtual key code number
// into the standard keycodes used by the key package.
//
// To get a sense of the key map, see the diagram on
//	http://boredzo.org/blog/archives/2007-05-22/virtual-key-codes
func convVirtualKeyCode(vkcode uint16) key.Code {
	switch vkcode {
	case C.kVK_ANSI_A:
		return key.CodeA
	case C.kVK_ANSI_B:
		return key.CodeB
	case C.kVK_ANSI_C:
		return key.CodeC
	case C.kVK_ANSI_D:
		return key.CodeD
	case C.kVK_ANSI_E:
		return key.CodeE
	case C.kVK_ANSI_F:
		return key.CodeF
	case C.kVK_ANSI_G:
		return key.CodeG
	case C.kVK_ANSI_H:
		return key.CodeH
	case C.kVK_ANSI_I:
		return key.CodeI
	case C.kVK_ANSI_J:
		return key.CodeJ
	case C.kVK_ANSI_K:
		return key.CodeK
	case C.kVK_ANSI_L:
		return key.CodeL
	case C.kVK_ANSI_M:
		return key.CodeM
	case C.kVK_ANSI_N:
		return key.CodeN
	case C.kVK_ANSI_O:
		return key.CodeO
	case C.kVK_ANSI_P:
		return key.CodeP
	case C.kVK_ANSI_Q:
		return key.CodeQ
	case C.kVK_ANSI_R:
		return key.CodeR
	case C.kVK_ANSI_S:
		return key.CodeS
	case C.kVK_ANSI_T:
		return key.CodeT
	case C.kVK_ANSI_U:
		return key.CodeU
	case C.kVK_ANSI_V:
		return key.CodeV
	case C.kVK_ANSI_W:
		return key.CodeW
	case C.kVK_ANSI_X:
		return key.CodeX
	case C.kVK_ANSI_Y:
		return key.CodeY
	case C.kVK_ANSI_Z:
		return key.CodeZ
	case C.kVK_ANSI_1:
		return key.Code1
	case C.kVK_ANSI_2:
		return key.Code2
	case C.kVK_ANSI_3:
		return key.Code3
	case C.kVK_ANSI_4:
		return key.Code4
	case C.kVK_ANSI_5:
		return key.Code5
	case C.kVK_ANSI_6:
		return key.Code6
	case C.kVK_ANSI_7:
		return key.Code7
	case C.kVK_ANSI_8:
		return key.Code8
	case C.kVK_ANSI_9:
		return key.Code9
	case C.kVK_ANSI_0:
		return key.Code0
	// TODO: move the rest of these codes to constants in key.go
	// if we are happy with them.
	case C.kVK_Return:
		return key.CodeReturnEnter
	case C.kVK_Escape:
		return key.CodeEscape
	case C.kVK_Delete:
		return key.CodeDeleteBackspace
	case C.kVK_Tab:
		return key.CodeTab
	case C.kVK_Space:
		return key.CodeSpacebar
	case C.kVK_ANSI_Minus:
		return key.CodeHyphenMinus
	case C.kVK_ANSI_Equal:
		return key.CodeEqualSign
	case C.kVK_ANSI_LeftBracket:
		return key.CodeLeftSquareBracket
	case C.kVK_ANSI_RightBracket:
		return key.CodeRightSquareBracket
	case C.kVK_ANSI_Backslash:
		return key.CodeBackslash
	// 50: Keyboard Non-US "#" and ~
	case C.kVK_ANSI_Semicolon:
		return key.CodeSemicolon
	case C.kVK_ANSI_Quote:
		return key.CodeApostrophe
	case C.kVK_ANSI_Grave:
		return key.CodeGraveAccent
	case C.kVK_ANSI_Comma:
		return key.CodeComma
	case C.kVK_ANSI_Period:
		return key.CodeFullStop
	case C.kVK_ANSI_Slash:
		return key.CodeSlash
	case C.kVK_CapsLock:
		return key.CodeCapsLock
	case C.kVK_F1:
		return key.CodeF1
	case C.kVK_F2:
		return key.CodeF2
	case C.kVK_F3:
		return key.CodeF3
	case C.kVK_F4:
		return key.CodeF4
	case C.kVK_F5:
		return key.CodeF5
	case C.kVK_F6:
		return key.CodeF6
	case C.kVK_F7:
		return key.CodeF7
	case C.kVK_F8:
		return key.CodeF8
	case C.kVK_F9:
		return key.CodeF9
	case C.kVK_F10:
		return key.CodeF10
	case C.kVK_F11:
		return key.CodeF11
	case C.kVK_F12:
		return key.CodeF12
	// 70: PrintScreen
	// 71: Scroll Lock
	// 72: Pause
	// 73: Insert
	case C.kVK_Home:
		return key.CodeHome
	case C.kVK_PageUp:
		return key.CodePageUp
	case C.kVK_ForwardDelete:
		return key.CodeDeleteForward
	case C.kVK_End:
		return key.CodeEnd
	case C.kVK_PageDown:
		return key.CodePageDown
	case C.kVK_RightArrow:
		return key.CodeRightArrow
	case C.kVK_LeftArrow:
		return key.CodeLeftArrow
	case C.kVK_DownArrow:
		return key.CodeDownArrow
	case C.kVK_UpArrow:
		return key.CodeUpArrow
	case C.kVK_ANSI_KeypadClear:
		return key.CodeKeypadNumLock
	case C.kVK_ANSI_KeypadDivide:
		return key.CodeKeypadSlash
	case C.kVK_ANSI_KeypadMultiply:
		return key.CodeKeypadAsterisk
	case C.kVK_ANSI_KeypadMinus:
		return key.CodeKeypadHyphenMinus
	case C.kVK_ANSI_KeypadPlus:
		return key.CodeKeypadPlusSign
	case C.kVK_ANSI_KeypadEnter:
		return key.CodeKeypadEnter
	case C.kVK_ANSI_Keypad1:
		return key.CodeKeypad1
	case C.kVK_ANSI_Keypad2:
		return key.CodeKeypad2
	case C.kVK_ANSI_Keypad3:
		return key.CodeKeypad3
	case C.kVK_ANSI_Keypad4:
		return key.CodeKeypad4
	case C.kVK_ANSI_Keypad5:
		return key.CodeKeypad5
	case C.kVK_ANSI_Keypad6:
		return key.CodeKeypad6
	case C.kVK_ANSI_Keypad7:
		return key.CodeKeypad7
	case C.kVK_ANSI_Keypad8:
		return key.CodeKeypad8
	case C.kVK_ANSI_Keypad9:
		return key.CodeKeypad9
	case C.kVK_ANSI_Keypad0:
		return key.CodeKeypad0
	case C.kVK_ANSI_KeypadDecimal:
		return key.CodeKeypadFullStop
	case C.kVK_ANSI_KeypadEquals:
		return key.CodeKeypadEqualSign
	case C.kVK_F13:
		return key.CodeF13
	case C.kVK_F14:
		return key.CodeF14
	case C.kVK_F15:
		return key.CodeF15
	case C.kVK_F16:
		return key.CodeF16
	case C.kVK_F17:
		return key.CodeF17
	case C.kVK_F18:
		return key.CodeF18
	case C.kVK_F19:
		return key.CodeF19
	case C.kVK_F20:
		return key.CodeF20
	// 116: Keyboard Execute
	case C.kVK_Help:
		return key.CodeHelp
	// 118: Keyboard Menu
	// 119: Keyboard Select
	// 120: Keyboard Stop
	// 121: Keyboard Again
	// 122: Keyboard Undo
	// 123: Keyboard Cut
	// 124: Keyboard Copy
	// 125: Keyboard Paste
	// 126: Keyboard Find
	case C.kVK_Mute:
		return key.CodeMute
	case C.kVK_VolumeUp:
		return key.CodeVolumeUp
	case C.kVK_VolumeDown:
		return key.CodeVolumeDown
	// 130: Keyboard Locking Caps Lock
	// 131: Keyboard Locking Num Lock
	// 132: Keyboard Locking Scroll Lock
	// 133: Keyboard Comma
	// 134: Keyboard Equal Sign
	// ...: Bunch of stuff
	case C.kVK_Control:
		return key.CodeLeftControl
	case C.kVK_Shift:
		return key.CodeLeftShift
	case C.kVK_Option:
		return key.CodeLeftAlt
	case C.kVK_Command:
		return key.CodeLeftGUI
	case C.kVK_RightControl:
		return key.CodeRightControl
	case C.kVK_RightShift:
		return key.CodeRightShift
	case C.kVK_RightOption:
		return key.CodeRightAlt
	// TODO key.CodeRightGUI
	default:
		return key.CodeUnknown
	}
}
