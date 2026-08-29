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

func main(f func(App)) {
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

func GoBack() {
	// When simulating mobile there are no other activities open (and we can't just force background)
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
func onKeyPress(keySym int) {
	keyCode := X11KeySymToFyneKeyCode(keySym)
	theApp.events.In() <- key.Event{
		Direction: key.DirPress,
		Code:      keyCode,
		Rune:      codeToRune(keyCode),
	}
}

//export onKeyRelease
func onKeyRelease(keySym int) {
	keyCode := X11KeySymToFyneKeyCode(keySym)
	theApp.events.In() <- key.Event{
		Direction: key.DirRelease,
		Code:      keyCode,
		Rune:      codeToRune(keyCode),
	}
}

func codeToRune(c key.Code) rune {
	switch c {
	case key.CodeA:
		return 'a'
	case key.CodeB:
		return 'b'
	case key.CodeC:
		return 'c'
	case key.CodeD:
		return 'd'
	case key.CodeE:
		return 'e'
	case key.CodeF:
		return 'f'
	case key.CodeG:
		return 'g'
	case key.CodeH:
		return 'h'
	case key.CodeI:
		return 'i'
	case key.CodeJ:
		return 'j'
	case key.CodeK:
		return 'k'
	case key.CodeL:
		return 'l'
	case key.CodeM:
		return 'm'
	case key.CodeN:
		return 'n'
	case key.CodeO:
		return 'o'
	case key.CodeP:
		return 'p'
	case key.CodeQ:
		return 'q'
	case key.CodeR:
		return 'r'
	case key.CodeS:
		return 's'
	case key.CodeT:
		return 't'
	case key.CodeU:
		return 'u'
	case key.CodeV:
		return 'v'
	case key.CodeW:
		return 'w'
	case key.CodeX:
		return 'x'
	case key.CodeY:
		return 'y'
	case key.CodeZ:
		return 'z'

	case key.Code1:
		return '1'
	case key.Code2:
		return '2'
	case key.Code3:
		return '3'
	case key.Code4:
		return '4'
	case key.Code5:
		return '5'
	case key.Code6:
		return '6'
	case key.Code7:
		return '7'
	case key.Code8:
		return '8'
	case key.Code9:
		return '9'
	case key.Code0:
		return '0'

	case key.CodeSpacebar:
		return ' '
	case key.CodeHyphenMinus:
		return '-'
	case key.CodeEqualSign:
		return '='
	case key.CodeLeftSquareBracket:
		return '['
	case key.CodeRightSquareBracket:
		return ']'
	case key.CodeBackslash:
		return '\\'
	case key.CodeSemicolon:
		return ';'
	case key.CodeApostrophe:
		return '\''
	case key.CodeGraveAccent:
		return '`'
	case key.CodeComma:
		return ','
	case key.CodeFullStop:
		return '.'
	case key.CodeSlash:
		return '/'

	case key.CodeReturnEnter:
		return '\n'
	case key.CodeTab:
		return '\t'

	default:
		return 0 // Unsupported code or non-printable character
	}
}

func X11KeySymToFyneKeyCode(keysym int) key.Code {
	// Handle alphabetic characters (a-z, A-Z)
	if keysym >= 0x0061 && keysym <= 0x007A { // Lowercase a-z
		return key.Code(0x04 + (keysym - 0x0061))
	}
	if keysym >= 0x0041 && keysym <= 0x005A { // Uppercase A-Z
		return key.Code(0x04 + (keysym - 0x0041))
	}
	if keysym == 0x0030 {
		return key.Code0
	}
	if keysym >= 0x0031 && keysym <= 0x0039 { // 1-9
		return key.Code(30 + (keysym - 0x0031))
	}

	switch keysym {
	case 0xFF0D:
		return key.CodeReturnEnter // XK_Return
	case 0xFF1B:
		return key.CodeEscape // XK_Escape
	case 0xFF08:
		return key.CodeDeleteBackspace // XK_BackSpace
	case 0xFF09:
		return key.CodeTab // XK_Tab
	case 0x0020:
		return key.CodeSpacebar // XK_space
	case 0x002D:
		return key.CodeHyphenMinus // XK_minus
	case 0x003D:
		return key.CodeEqualSign // XK_equal
	case 0x005B:
		return key.CodeLeftSquareBracket // XK_bracketleft
	case 0x005D:
		return key.CodeRightSquareBracket // XK_bracketright
	case 0x005C:
		return key.CodeBackslash // XK_backslash
	case 0x003B:
		return key.CodeSemicolon // XK_semicolon
	case 0x0027:
		return key.CodeApostrophe // XK_apostrophe
	case 0x0060:
		return key.CodeGraveAccent // XK_grave
	case 0x002C:
		return key.CodeComma // XK_comma
	case 0x002E:
		return key.CodeFullStop // XK_period
	case 0x002F:
		return key.CodeSlash // XK_slash
	case 0xFFE5:
		return key.CodeCapsLock // XK_Caps_Lock
	case 0xFFBE:
		return key.CodeF1 // XK_F1
	case 0xFFBF:
		return key.CodeF2 // XK_F2
	case 0xFFC0:
		return key.CodeF3 // XK_F3
	case 0xFFC1:
		return key.CodeF4 // XK_F4
	case 0xFFC2:
		return key.CodeF5 // XK_F5
	case 0xFFC3:
		return key.CodeF6 // XK_F6
	case 0xFFC4:
		return key.CodeF7 // XK_F7
	case 0xFFC5:
		return key.CodeF8 // XK_F8
	case 0xFFC6:
		return key.CodeF9 // XK_F9
	case 0xFFC7:
		return key.CodeF10 // XK_F10
	case 0xFFC8:
		return key.CodeF11 // XK_F11
	case 0xFFC9:
		return key.CodeF12 // XK_F12
	case 0xFF13:
		return key.CodePause // XK_Pause
	case 0xFF63:
		return key.CodeInsert // XK_Insert
	case 0xFF50:
		return key.CodeHome // XK_Home
	case 0xFF55:
		return key.CodePageUp // XK_Page_Up
	case 0xFFFF:
		return key.CodeDeleteForward // XK_Delete
	case 0xFF57:
		return key.CodeEnd // XK_End
	case 0xFF56:
		return key.CodePageDown // XK_Page_Down
	case 0xFF53:
		return key.CodeRightArrow // XK_Right
	case 0xFF51:
		return key.CodeLeftArrow // XK_Left
	case 0xFF54:
		return key.CodeDownArrow // XK_Down
	case 0xFF52:
		return key.CodeUpArrow // XK_Up
	case 0xFF7F:
		return key.CodeKeypadNumLock // XK_Num_Lock
	case 0xFFAF:
		return key.CodeKeypadSlash // XK_KP_Divide
	case 0xFFAA:
		return key.CodeKeypadAsterisk // XK_KP_Multiply
	case 0xFFAD:
		return key.CodeKeypadHyphenMinus // XK_KP_Subtract
	case 0xFFAB:
		return key.CodeKeypadPlusSign // XK_KP_Add
	case 0xFF8D:
		return key.CodeKeypadEnter // XK_KP_Enter
	case 0xFF9C:
		return key.CodeKeypad1 // XK_KP_1
	case 0xFF99:
		return key.CodeKeypad2 // XK_KP_2
	case 0xFF9B:
		return key.CodeKeypad3 // XK_KP_3
	case 0xFF96:
		return key.CodeKeypad4 // XK_KP_4
	case 0xFF9D:
		return key.CodeKeypad5 // XK_KP_5
	case 0xFF98:
		return key.CodeKeypad6 // XK_KP_6
	case 0xFF95:
		return key.CodeKeypad7 // XK_KP_7
	case 0xFF97:
		return key.CodeKeypad8 // XK_KP_8
	case 0xFF9A:
		return key.CodeKeypad9 // XK_KP_9
	case 0xFF9E:
		return key.CodeKeypad0 // XK_KP_0
	case 0xFF9F:
		return key.CodeKeypadFullStop // XK_KP_Decimal
	case 0xFFBD:
		return key.CodeKeypadEqualSign // XK_KP_Equal
	case 0xFFCA:
		return key.CodeF13 // XK_F13
	case 0xFFCB:
		return key.CodeF14 // XK_F14
	case 0xFFCC:
		return key.CodeF15 // XK_F15
	case 0xFFCD:
		return key.CodeF16 // XK_F16
	case 0xFFCE:
		return key.CodeF17 // XK_F17
	case 0xFFCF:
		return key.CodeF18 // XK_F18
	case 0xFFD0:
		return key.CodeF19 // XK_F19
	case 0xFFD1:
		return key.CodeF20 // XK_F20
	case 0xFFD2:
		return key.CodeF21 // XK_F21
	case 0xFFD3:
		return key.CodeF22 // XK_F22
	case 0xFFD4:
		return key.CodeF23 // XK_F23
	case 0xFFD5:
		return key.CodeF24 // XK_F24
	case 0xFF6A:
		return key.CodeHelp // XK_Help
	case 0x1008FF12:
		return key.CodeMute // XF86AudioMute
	case 0x1008FF13:
		return key.CodeVolumeUp // XF86AudioRaiseVolume
	case 0x1008FF11:
		return key.CodeVolumeDown // XF86AudioLowerVolume
	case 0xFFE3:
		return key.CodeLeftControl // XK_Control_L
	case 0xFFE1:
		return key.CodeLeftShift // XK_Shift_L
	case 0xFFE9:
		return key.CodeLeftAlt // XK_Alt_L
	case 0xFFEB:
		return key.CodeLeftGUI // XK_Super_L
	case 0xFFE4:
		return key.CodeRightControl // XK_Control_R
	case 0xFFE2:
		return key.CodeRightShift // XK_Shift_R
	case 0xFFEA:
		return key.CodeRightAlt // XK_Alt_R
	case 0xFFEC:
		return key.CodeRightGUI // XK_Super_R
	}

	return 0 // Unknown or unmapped key
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
