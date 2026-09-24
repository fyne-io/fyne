// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || freebsd || openbsd

package app

/*
Simple on-screen app debugging for X11. Not an officially supported
development target for apps, as screens with mice are very different
than screens with touch panels.
*/

/*
#cgo LDFLAGS: -lEGL -lGLESv2 -lX11
#cgo freebsd CFLAGS: -I/usr/local/include/
#cgo openbsd CFLAGS: -I/usr/X11R6/include/

void createWindow(void);
void processEvents(void);
void swapBuffers(void);
*/
import "C"

import (
	"runtime"
	"time"

	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/lifecycle"
	"fyne.io/fyne/v2/internal/driver/mobile/event/paint"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
)

func init() {
	theApp.registerGLViewportFilter()
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

	workAvailable := theApp.worker.WorkAvailable()
	heartbeat := time.NewTicker(time.Second / 60)

	C.createWindow()

	// TODO: send lifecycle events when e.g. the X11 window is iconified or moved off-screen.
	theApp.sendLifecycle(lifecycle.StageFocused)

	// TODO: translate X11 expose events to shiny paint events, instead of
	// sending this synthetic paint event as a hack.
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
			C.processEvents()
		case <-workAvailable:
			theApp.worker.DoWork()
		case <-theApp.publish:
			C.swapBuffers()
			tc = ticker.C
		case <-tc:
			tc = nil
			theApp.publishResult <- PublishResult{}
		}
	}
}

//export onResize
func onResize(w, h int) {
	// TODO(nigeltao): don't assume 72 DPI. DisplayWidth and DisplayWidthMM
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

//export onTouchBegin
func onTouchBegin(x, y float32) { sendTouch(touch.TypeBegin, x, y) }

//export onTouchMove
func onTouchMove(x, y float32) { sendTouch(touch.TypeMove, x, y) }

//export onTouchEnd
func onTouchEnd(x, y float32) { sendTouch(touch.TypeEnd, x, y) }

//export onKeyPress
func onKeyPress(keysym C.long, r C.int, mods C.long) {
	sendKey(keysym, r, mods, key.DirPress)
}

//export onKeyRelease
func onKeyRelease(keysym C.long, r C.int, mods C.long) {
	sendKey(keysym, r, mods, key.DirRelease)
}

func sendKey(keysym C.long, r C.int, mods C.long, dir key.Direction) {
	theApp.events.In() <- key.Event{
		Rune:      rune(r),
		Code:      x11KeySymToFyneKeyCode(int(keysym)),
		Modifiers: key.Modifiers(mods),
		Direction: dir,
	}
}

// x11KeySymToFyneKeyCode maps an X11 keysym (the unshifted base keysym of a
// key) to the corresponding USB HID key code used by the mobile event key
// package. It returns key.CodeUnknown for unmapped keys.
func x11KeySymToFyneKeyCode(keysym int) key.Code {
	// Letters: XK_a..XK_z and XK_A..XK_Z.
	if keysym >= 0x0061 && keysym <= 0x007a {
		return key.Code(0x04 + keysym - 0x0061)
	}
	if keysym >= 0x0041 && keysym <= 0x005a {
		return key.Code(0x04 + keysym - 0x0041)
	}
	// Digits: XK_0..XK_9.
	if keysym == 0x0030 { // XK_0
		return key.Code0
	}
	if keysym >= 0x0031 && keysym <= 0x0039 { // XK_1..XK_9
		return key.Code(30 + keysym - 0x0031)
	}

	switch keysym {
	case 0x0020: // XK_space
		return key.CodeSpacebar
	case 0x0027: // XK_apostrophe
		return key.CodeApostrophe
	case 0x002c: // XK_comma
		return key.CodeComma
	case 0x002d: // XK_minus
		return key.CodeHyphenMinus
	case 0x002e: // XK_period
		return key.CodeFullStop
	case 0x002f: // XK_slash
		return key.CodeSlash
	case 0x003b: // XK_semicolon
		return key.CodeSemicolon
	case 0x003d: // XK_equal
		return key.CodeEqualSign
	case 0x005b: // XK_bracketleft
		return key.CodeLeftSquareBracket
	case 0x005c: // XK_backslash
		return key.CodeBackslash
	case 0x005d: // XK_bracketright
		return key.CodeRightSquareBracket
	case 0x0060: // XK_grave
		return key.CodeGraveAccent
	case 0xff08: // XK_BackSpace
		return key.CodeDeleteBackspace
	case 0xff09: // XK_Tab
		return key.CodeTab
	case 0xff0d: // XK_Return
		return key.CodeReturnEnter
	case 0xff13: // XK_Pause
		return key.CodePause
	case 0xff1b: // XK_Escape
		return key.CodeEscape
	case 0xff50: // XK_Home
		return key.CodeHome
	case 0xff51: // XK_Left
		return key.CodeLeftArrow
	case 0xff52: // XK_Up
		return key.CodeUpArrow
	case 0xff53: // XK_Right
		return key.CodeRightArrow
	case 0xff54: // XK_Down
		return key.CodeDownArrow
	case 0xff55: // XK_Page_Up
		return key.CodePageUp
	case 0xff56: // XK_Page_Down
		return key.CodePageDown
	case 0xff57: // XK_End
		return key.CodeEnd
	case 0xff63: // XK_Insert
		return key.CodeInsert
	case 0xff6a: // XK_Help
		return key.CodeHelp
	case 0xff7f: // XK_Num_Lock
		return key.CodeKeypadNumLock
	case 0xff8d: // XK_KP_Enter
		return key.CodeKeypadEnter
	case 0xff95: // XK_KP_Home
		return key.CodeHome
	case 0xff96: // XK_KP_Left
		return key.CodeLeftArrow
	case 0xff97: // XK_KP_Up
		return key.CodeUpArrow
	case 0xff98: // XK_KP_Right
		return key.CodeRightArrow
	case 0xff99: // XK_KP_Down
		return key.CodeDownArrow
	case 0xff9a: // XK_KP_Page_Up
		return key.CodePageUp
	case 0xff9b: // XK_KP_Page_Down
		return key.CodePageDown
	case 0xff9c: // XK_KP_End
		return key.CodeEnd
	case 0xff9d: // XK_KP_Begin
		return key.CodeHome
	case 0xff9e: // XK_KP_Insert
		return key.CodeInsert
	case 0xff9f: // XK_KP_Delete
		return key.CodeDeleteForward
	case 0xffaa: // XK_KP_Multiply
		return key.CodeKeypadAsterisk
	case 0xffab: // XK_KP_Add
		return key.CodeKeypadPlusSign
	case 0xffad: // XK_KP_Subtract
		return key.CodeKeypadHyphenMinus
	case 0xffae: // XK_KP_Decimal
		return key.CodeKeypadFullStop
	case 0xffaf: // XK_KP_Divide
		return key.CodeKeypadSlash
	case 0xffb0: // XK_KP_0
		return key.CodeKeypad0
	case 0xffb1: // XK_KP_1
		return key.CodeKeypad1
	case 0xffb2: // XK_KP_2
		return key.CodeKeypad2
	case 0xffb3: // XK_KP_3
		return key.CodeKeypad3
	case 0xffb4: // XK_KP_4
		return key.CodeKeypad4
	case 0xffb5: // XK_KP_5
		return key.CodeKeypad5
	case 0xffb6: // XK_KP_6
		return key.CodeKeypad6
	case 0xffb7: // XK_KP_7
		return key.CodeKeypad7
	case 0xffb8: // XK_KP_8
		return key.CodeKeypad8
	case 0xffb9: // XK_KP_9
		return key.CodeKeypad9
	case 0xffbd: // XK_KP_Equal
		return key.CodeKeypadEqualSign
	case 0xffbe: // XK_F1
		return key.CodeF1
	case 0xffbf: // XK_F2
		return key.CodeF2
	case 0xffc0: // XK_F3
		return key.CodeF3
	case 0xffc1: // XK_F4
		return key.CodeF4
	case 0xffc2: // XK_F5
		return key.CodeF5
	case 0xffc3: // XK_F6
		return key.CodeF6
	case 0xffc4: // XK_F7
		return key.CodeF7
	case 0xffc5: // XK_F8
		return key.CodeF8
	case 0xffc6: // XK_F9
		return key.CodeF9
	case 0xffc7: // XK_F10
		return key.CodeF10
	case 0xffc8: // XK_F11
		return key.CodeF11
	case 0xffc9: // XK_F12
		return key.CodeF12
	case 0xffca: // XK_F13
		return key.CodeF13
	case 0xffcb: // XK_F14
		return key.CodeF14
	case 0xffcc: // XK_F15
		return key.CodeF15
	case 0xffcd: // XK_F16
		return key.CodeF16
	case 0xffce: // XK_F17
		return key.CodeF17
	case 0xffcf: // XK_F18
		return key.CodeF18
	case 0xffd0: // XK_F19
		return key.CodeF19
	case 0xffd1: // XK_F20
		return key.CodeF20
	case 0xffd2: // XK_F21
		return key.CodeF21
	case 0xffd3: // XK_F22
		return key.CodeF22
	case 0xffd4: // XK_F23
		return key.CodeF23
	case 0xffd5: // XK_F24
		return key.CodeF24
	case 0xffe1: // XK_Shift_L
		return key.CodeLeftShift
	case 0xffe2: // XK_Shift_R
		return key.CodeRightShift
	case 0xffe3: // XK_Control_L
		return key.CodeLeftControl
	case 0xffe4: // XK_Control_R
		return key.CodeRightControl
	case 0xffe5: // XK_Caps_Lock
		return key.CodeCapsLock
	case 0xffe9: // XK_Alt_L
		return key.CodeLeftAlt
	case 0xffea: // XK_Alt_R
		return key.CodeRightAlt
	case 0xffeb: // XK_Super_L
		return key.CodeLeftGUI
	case 0xffec: // XK_Super_R
		return key.CodeRightGUI
	case 0xffff: // XK_Delete
		return key.CodeDeleteForward
	case 0x1008ff11: // XF86AudioLowerVolume
		return key.CodeVolumeDown
	case 0x1008ff12: // XF86AudioMute
		return key.CodeMute
	case 0x1008ff13: // XF86AudioRaiseVolume
		return key.CodeVolumeUp
	}

	return key.CodeUnknown
}

var stopped bool

//export onStop
func onStop() {
	if stopped {
		return
	}
	stopped = true
	theApp.sendLifecycle(lifecycle.StageDead)
	theApp.events.Close()
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
