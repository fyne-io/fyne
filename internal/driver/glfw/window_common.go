package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var keyCodeMap = map[glfw.Key]fyne.KeyName{
	// non-printable
	glfw.KeyEscape:    fyne.KeyEscape,
	glfw.KeyEnter:     fyne.KeyReturn,
	glfw.KeyTab:       fyne.KeyTab,
	glfw.KeyBackspace: fyne.KeyBackspace,
	glfw.KeyInsert:    fyne.KeyInsert,
	glfw.KeyDelete:    fyne.KeyDelete,
	glfw.KeyRight:     fyne.KeyRight,
	glfw.KeyLeft:      fyne.KeyLeft,
	glfw.KeyDown:      fyne.KeyDown,
	glfw.KeyUp:        fyne.KeyUp,
	glfw.KeyPageUp:    fyne.KeyPageUp,
	glfw.KeyPageDown:  fyne.KeyPageDown,
	glfw.KeyHome:      fyne.KeyHome,
	glfw.KeyEnd:       fyne.KeyEnd,

	glfw.KeySpace:   fyne.KeySpace,
	glfw.KeyKPEnter: fyne.KeyEnter,

	// functions
	glfw.KeyF1:  fyne.KeyF1,
	glfw.KeyF2:  fyne.KeyF2,
	glfw.KeyF3:  fyne.KeyF3,
	glfw.KeyF4:  fyne.KeyF4,
	glfw.KeyF5:  fyne.KeyF5,
	glfw.KeyF6:  fyne.KeyF6,
	glfw.KeyF7:  fyne.KeyF7,
	glfw.KeyF8:  fyne.KeyF8,
	glfw.KeyF9:  fyne.KeyF9,
	glfw.KeyF10: fyne.KeyF10,
	glfw.KeyF11: fyne.KeyF11,
	glfw.KeyF12: fyne.KeyF12,

	// numbers - lookup by code to avoid AZERTY using the symbol name instead of number
	glfw.Key0:   fyne.Key0,
	glfw.KeyKP0: fyne.Key0,
	glfw.Key1:   fyne.Key1,
	glfw.KeyKP1: fyne.Key1,
	glfw.Key2:   fyne.Key2,
	glfw.KeyKP2: fyne.Key2,
	glfw.Key3:   fyne.Key3,
	glfw.KeyKP3: fyne.Key3,
	glfw.Key4:   fyne.Key4,
	glfw.KeyKP4: fyne.Key4,
	glfw.Key5:   fyne.Key5,
	glfw.KeyKP5: fyne.Key5,
	glfw.Key6:   fyne.Key6,
	glfw.KeyKP6: fyne.Key6,
	glfw.Key7:   fyne.Key7,
	glfw.KeyKP7: fyne.Key7,
	glfw.Key8:   fyne.Key8,
	glfw.KeyKP8: fyne.Key8,
	glfw.Key9:   fyne.Key9,
	glfw.KeyKP9: fyne.Key9,

	// desktop
	glfw.KeyLeftShift:    desktop.KeyShiftLeft,
	glfw.KeyRightShift:   desktop.KeyShiftRight,
	glfw.KeyLeftControl:  desktop.KeyControlLeft,
	glfw.KeyRightControl: desktop.KeyControlRight,
	glfw.KeyLeftAlt:      desktop.KeyAltLeft,
	glfw.KeyRightAlt:     desktop.KeyAltRight,
	glfw.KeyLeftSuper:    desktop.KeySuperLeft,
	glfw.KeyRightSuper:   desktop.KeySuperRight,
	glfw.KeyMenu:         desktop.KeyMenu,
	glfw.KeyPrintScreen:  desktop.KeyPrintScreen,
	glfw.KeyCapsLock:     desktop.KeyCapsLock,
}

func unknownKey(key string) bool {
	if len(key) != 1 {
		return true
	}

	keyRune := key[0]
	if keyRune >= 'a' || keyRune <= 'z' {
		return false
	}

	_, ok := keyKnownRuneMap[keyRune]
	return !ok
}

// keyKnownRuneMap does not contain 'a' to 'z'. They are handled outside.
var keyKnownRuneMap = map[byte]struct{}{
	'\'': {},
	',':  {},
	'-':  {},
	'.':  {},
	'/':  {},
	'*':  {},
	'`':  {},

	';': {},
	'+': {},
	'=': {},

	'[':  {},
	'\\': {},
	']':  {},
}

func convertASCII(key glfw.Key) fyne.KeyName {
	keyRune := rune('A' + key - glfw.KeyA)
	if keyRune < 'A' || keyRune > 'Z' {
		return fyne.KeyUnknown
	}

	return fyne.KeyName(keyRune)
}
