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

// X11 keysym values handled by the driver. The values match the
// <X11/keysymdef.h> definitions for the respective key names.
const (
	xkSpace        = 0x20 // XK_space
	xkApostrophe   = 0x27 // XK_apostrophe
	xkComma        = 0x2c // XK_comma
	xkMinus        = 0x2d // XK_minus
	xkPeriod       = 0x2e // XK_period
	xkSlash        = 0x2f // XK_slash
	xk0            = 0x30 // XK_0
	xk1            = 0x31 // XK_1
	xk9            = 0x39 // XK_9
	xkSemicolon    = 0x3b // XK_semicolon
	xkEqual        = 0x3d // XK_equal
	xkA            = 0x41 // XK_A
	xkZ            = 0x5a // XK_Z
	xkBracketLeft  = 0x5b // XK_bracketleft
	xkBackslash    = 0x5c // XK_backslash
	xkBracketRight = 0x5d // XK_bracketright
	xkGrave        = 0x60 // XK_grave
	xkLowerA       = 0x61 // XK_a
	xkLowerZ       = 0x7a // XK_z

	xkBackSpace  = 0xff08 // XK_BackSpace
	xkTab        = 0xff09 // XK_Tab
	xkReturn     = 0xff0d // XK_Return
	xkPause      = 0xff13 // XK_Pause
	xkEscape     = 0xff1b // XK_Escape
	xkHome       = 0xff50 // XK_Home
	xkLeft       = 0xff51 // XK_Left
	xkUp         = 0xff52 // XK_Up
	xkRight      = 0xff53 // XK_Right
	xkDown       = 0xff54 // XK_Down
	xkPageUp     = 0xff55 // XK_Page_Up
	xkPageDown   = 0xff56 // XK_Page_Down
	xkEnd        = 0xff57 // XK_End
	xkInsert     = 0xff63 // XK_Insert
	xkHelp       = 0xff6a // XK_Help
	xkNumLock    = 0xff7f // XK_Num_Lock
	xkKPEnter    = 0xff8d // XK_KP_Enter
	xkKPHome     = 0xff95 // XK_KP_Home
	xkKPLeft     = 0xff96 // XK_KP_Left
	xkKPUp       = 0xff97 // XK_KP_Up
	xkKPRight    = 0xff98 // XK_KP_Right
	xkKPDown     = 0xff99 // XK_KP_Down
	xkKPPageUp   = 0xff9a // XK_KP_Page_Up
	xkKPPageDown = 0xff9b // XK_KP_Page_Down
	xkKPEnd      = 0xff9c // XK_KP_End
	xkKPBegin    = 0xff9d // XK_KP_Begin
	xkKPInsert   = 0xff9e // XK_KP_Insert
	xkKPDelete   = 0xff9f // XK_KP_Delete
	xkKPMultiply = 0xffaa // XK_KP_Multiply
	xkKPAdd      = 0xffab // XK_KP_Add
	xkKPSubtract = 0xffad // XK_KP_Subtract
	xkKPDecimal  = 0xffae // XK_KP_Decimal
	xkKPDivide   = 0xffaf // XK_KP_Divide
	xkKP0        = 0xffb0 // XK_KP_0
	xkKP1        = 0xffb1 // XK_KP_1
	xkKP9        = 0xffb9 // XK_KP_9
	xkKPEqual    = 0xffbd // XK_KP_Equal
	xkF1         = 0xffbe // XK_F1
	xkF12        = 0xffc9 // XK_F12
	xkF13        = 0xffca // XK_F13
	xkF24        = 0xffd5 // XK_F24
	xkShiftL     = 0xffe1 // XK_Shift_L
	xkShiftR     = 0xffe2 // XK_Shift_R
	xkControlL   = 0xffe3 // XK_Control_L
	xkControlR   = 0xffe4 // XK_Control_R
	xkCapsLock   = 0xffe5 // XK_Caps_Lock
	xkAltL       = 0xffe9 // XK_Alt_L
	xkAltR       = 0xffea // XK_Alt_R
	xkSuperL     = 0xffeb // XK_Super_L
	xkSuperR     = 0xffec // XK_Super_R
	xkDelete     = 0xffff // XK_Delete

	xkAudioLowerVolume = 0x1008ff11 // XF86AudioLowerVolume
	xkAudioMute        = 0x1008ff12 // XF86AudioMute
	xkAudioRaiseVolume = 0x1008ff13 // XF86AudioRaiseVolume
)

// x11KeyToCode maps the X11 keysyms with a 1:1 association to their Fyne
// key codes. Contiguous runs (letters, digits, function keys and keypad
// digits) are decoded arithmetically in x11KeySymToFyneKeyCode.
var x11KeyToCode = map[int]key.Code{
	xkSpace:            key.CodeSpacebar,
	xkApostrophe:       key.CodeApostrophe,
	xkComma:            key.CodeComma,
	xkMinus:            key.CodeHyphenMinus,
	xkPeriod:           key.CodeFullStop,
	xkSlash:            key.CodeSlash,
	xkSemicolon:        key.CodeSemicolon,
	xkEqual:            key.CodeEqualSign,
	xkBracketLeft:      key.CodeLeftSquareBracket,
	xkBackslash:        key.CodeBackslash,
	xkBracketRight:     key.CodeRightSquareBracket,
	xkGrave:            key.CodeGraveAccent,
	xkBackSpace:        key.CodeDeleteBackspace,
	xkTab:              key.CodeTab,
	xkReturn:           key.CodeReturnEnter,
	xkPause:            key.CodePause,
	xkEscape:           key.CodeEscape,
	xkHome:             key.CodeHome,
	xkLeft:             key.CodeLeftArrow,
	xkUp:               key.CodeUpArrow,
	xkRight:            key.CodeRightArrow,
	xkDown:             key.CodeDownArrow,
	xkPageUp:           key.CodePageUp,
	xkPageDown:         key.CodePageDown,
	xkEnd:              key.CodeEnd,
	xkInsert:           key.CodeInsert,
	xkHelp:             key.CodeHelp,
	xkNumLock:          key.CodeKeypadNumLock,
	xkKPEnter:          key.CodeKeypadEnter,
	xkKPHome:           key.CodeHome,
	xkKPLeft:           key.CodeLeftArrow,
	xkKPUp:             key.CodeUpArrow,
	xkKPRight:          key.CodeRightArrow,
	xkKPDown:           key.CodeDownArrow,
	xkKPPageUp:         key.CodePageUp,
	xkKPPageDown:       key.CodePageDown,
	xkKPEnd:            key.CodeEnd,
	xkKPBegin:          key.CodeHome,
	xkKPInsert:         key.CodeInsert,
	xkKPDelete:         key.CodeDeleteForward,
	xkKPMultiply:       key.CodeKeypadAsterisk,
	xkKPAdd:            key.CodeKeypadPlusSign,
	xkKPSubtract:       key.CodeKeypadHyphenMinus,
	xkKPDecimal:        key.CodeKeypadFullStop,
	xkKPDivide:         key.CodeKeypadSlash,
	xkKPEqual:          key.CodeKeypadEqualSign,
	xkShiftL:           key.CodeLeftShift,
	xkShiftR:           key.CodeRightShift,
	xkControlL:         key.CodeLeftControl,
	xkControlR:         key.CodeRightControl,
	xkCapsLock:         key.CodeCapsLock,
	xkAltL:             key.CodeLeftAlt,
	xkAltR:             key.CodeRightAlt,
	xkSuperL:           key.CodeLeftGUI,
	xkSuperR:           key.CodeRightGUI,
	xkDelete:           key.CodeDeleteForward,
	xkAudioLowerVolume: key.CodeVolumeDown,
	xkAudioMute:        key.CodeMute,
	xkAudioRaiseVolume: key.CodeVolumeUp,
}

// x11KeySymToFyneKeyCode maps an X11 keysym (the unshifted base keysym of a
// key) to the corresponding USB HID key code used by the mobile event key
// package. It returns key.CodeUnknown for unmapped keys.
func x11KeySymToFyneKeyCode(keysym int) key.Code {
	if code, ok := x11KeyToCode[keysym]; ok {
		return code
	}

	// Letters: XK_a..XK_z and XK_A..XK_Z map onto the contiguous
	// HID codes CodeA..CodeZ.
	if keysym >= xkLowerA && keysym <= xkLowerZ {
		return key.Code(int(key.CodeA) + keysym - xkLowerA)
	}
	if keysym >= xkA && keysym <= xkZ {
		return key.Code(int(key.CodeA) + keysym - xkA)
	}
	// Digits: XK_1..XK_9 map onto the contiguous HID codes Code1..Code9,
	// followed by XK_0 as Code0.
	if keysym >= xk1 && keysym <= xk9 {
		return key.Code(int(key.Code1) + keysym - xk1)
	}
	if keysym == xk0 {
		return key.Code0
	}
	// Function keys: the HID codes are contiguous per segment F1..F12
	// and F13..F24.
	if keysym >= xkF1 && keysym <= xkF12 {
		return key.Code(int(key.CodeF1) + keysym - xkF1)
	}
	if keysym >= xkF13 && keysym <= xkF24 {
		return key.Code(int(key.CodeF13) + keysym - xkF13)
	}
	// Keypad digits: XK_KP_1..XK_KP_9 map onto the contiguous HID codes
	// CodeKeypad1..CodeKeypad9, followed by XK_KP_0 as CodeKeypad0.
	if keysym >= xkKP1 && keysym <= xkKP9 {
		return key.Code(int(key.CodeKeypad1) + keysym - xkKP1)
	}
	if keysym == xkKP0 {
		return key.CodeKeypad0
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
