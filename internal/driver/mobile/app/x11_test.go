//go:build (linux && !android) || freebsd || openbsd

package app

import (
	"testing"

	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
)

func TestX11KeySymToFyneKeyCode(t *testing.T) {
	tests := map[int]key.Code{
		// letters and digits
		0x0061: key.CodeA, // XK_a
		0x007a: key.CodeZ, // XK_z
		0x0041: key.CodeA, // XK_A
		0x005a: key.CodeZ, // XK_Z
		0x0030: key.Code0, // XK_0
		0x0039: key.Code9, // XK_9
		// punctuation
		0x0020: key.CodeSpacebar,    // XK_space
		0x0027: key.CodeApostrophe,  // XK_apostrophe
		0x002c: key.CodeComma,       // XK_comma
		0x002d: key.CodeHyphenMinus, // XK_minus
		0x002e: key.CodeFullStop,    // XK_period
		0x002f: key.CodeSlash,       // XK_slash
		0x003b: key.CodeSemicolon,   // XK_semicolon
		0x003d: key.CodeEqualSign,   // XK_equal
		0x005b: key.CodeLeftSquareBracket,
		0x005c: key.CodeBackslash,
		0x005d: key.CodeRightSquareBracket,
		0x0060: key.CodeGraveAccent, // XK_grave
		// control and navigation
		0xff08: key.CodeDeleteBackspace, // XK_BackSpace
		0xff09: key.CodeTab,             // XK_Tab
		0xff0d: key.CodeReturnEnter,     // XK_Return
		0xff1b: key.CodeEscape,          // XK_Escape
		0xff50: key.CodeHome,            // XK_Home
		0xff51: key.CodeLeftArrow,       // XK_Left
		0xff52: key.CodeUpArrow,         // XK_Up
		0xff53: key.CodeRightArrow,      // XK_Right
		0xff54: key.CodeDownArrow,       // XK_Down
		0xff55: key.CodePageUp,          // XK_Page_Up
		0xff56: key.CodePageDown,        // XK_Page_Down
		0xff57: key.CodeEnd,             // XK_End
		0xff63: key.CodeInsert,          // XK_Insert
		0xffff: key.CodeDeleteForward,   // XK_Delete
		// function keys
		0xffbe: key.CodeF1,
		0xffc9: key.CodeF12,
		0xffd5: key.CodeF24,
		// keypad (number keys with num lock on)
		0xff7f: key.CodeKeypadNumLock,     // XK_Num_Lock
		0xff8d: key.CodeKeypadEnter,       // XK_KP_Enter
		0xffaa: key.CodeKeypadAsterisk,    // XK_KP_Multiply
		0xffab: key.CodeKeypadPlusSign,    // XK_KP_Add
		0xffad: key.CodeKeypadHyphenMinus, // XK_KP_Subtract
		0xffae: key.CodeKeypadFullStop,    // XK_KP_Decimal
		0xffaf: key.CodeKeypadSlash,       // XK_KP_Divide
		0xffb0: key.CodeKeypad0,           // XK_KP_0
		0xffb1: key.CodeKeypad1,           // XK_KP_1
		0xffb2: key.CodeKeypad2,           // XK_KP_2
		0xffb3: key.CodeKeypad3,           // XK_KP_3
		0xffb4: key.CodeKeypad4,           // XK_KP_4
		0xffb5: key.CodeKeypad5,           // XK_KP_5
		0xffb6: key.CodeKeypad6,           // XK_KP_6
		0xffb7: key.CodeKeypad7,           // XK_KP_7
		0xffb8: key.CodeKeypad8,           // XK_KP_8
		0xffb9: key.CodeKeypad9,           // XK_KP_9
		0xffbd: key.CodeKeypadEqualSign,   // XK_KP_Equal
		// keypad (navigation with num lock off)
		0xff95: key.CodeHome,          // XK_KP_Home
		0xff96: key.CodeLeftArrow,     // XK_KP_Left
		0xff97: key.CodeUpArrow,       // XK_KP_Up
		0xff98: key.CodeRightArrow,    // XK_KP_Right
		0xff99: key.CodeDownArrow,     // XK_KP_Down
		0xff9a: key.CodePageUp,        // XK_KP_Page_Up
		0xff9b: key.CodePageDown,      // XK_KP_Page_Down
		0xff9c: key.CodeEnd,           // XK_KP_End
		0xff9d: key.CodeHome,          // XK_KP_Begin
		0xff9e: key.CodeInsert,        // XK_KP_Insert
		0xff9f: key.CodeDeleteForward, // XK_KP_Delete
		// modifiers
		0xffe1: key.CodeLeftShift,
		0xffe2: key.CodeRightShift,
		0xffe3: key.CodeLeftControl,
		0xffe4: key.CodeRightControl,
		0xffe5: key.CodeCapsLock,
		0xffe9: key.CodeLeftAlt,
		0xffea: key.CodeRightAlt,
		0xffeb: key.CodeLeftGUI,
		0xffec: key.CodeRightGUI,
		// media keys
		0x1008ff11: key.CodeVolumeDown,
		0x1008ff12: key.CodeMute,
		0x1008ff13: key.CodeVolumeUp,
	}

	for keysym, want := range tests {
		if got := x11KeySymToFyneKeyCode(keysym); got != want {
			t.Errorf("x11KeySymToFyneKeyCode(0x%x) = %v, want %v", keysym, got, want)
		}
	}

	if got := x11KeySymToFyneKeyCode(0x12345); got != key.CodeUnknown {
		t.Errorf("unmapped keysym should be CodeUnknown, got %v", got)
	}
}
