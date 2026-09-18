//go:build windows && directx

// Virtual key code to fyne.KeyName mapping, plus modifier state and the
// standard editing shortcuts.

package directx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Virtual key codes we translate. Values come from winuser.h.
const (
	vkBack     = 0x08
	vkTab      = 0x09
	vkReturn   = 0x0D
	vkShift    = 0x10
	vkControl  = 0x11
	vkMenu     = 0x12
	vkCapital  = 0x14
	vkEscape   = 0x1B
	vkSpace    = 0x20
	vkPrior    = 0x21
	vkNext     = 0x22
	vkEnd      = 0x23
	vkHome     = 0x24
	vkLeft     = 0x25
	vkUp       = 0x26
	vkRight    = 0x27
	vkDown     = 0x28
	vkSnapshot = 0x2C
	vkInsert   = 0x2D
	vkDelete   = 0x2E
	vkLWin     = 0x5B
	vkRWin     = 0x5C
	vkApps     = 0x5D
	vkNumpad0  = 0x60
	vkMultiply = 0x6A
	vkAdd      = 0x6B
	vkSubtract = 0x6D
	vkDecimal  = 0x6E
	vkDivide   = 0x6F
	vkF1       = 0x70
	vkLShift   = 0xA0
	vkRShift   = 0xA1
	vkLControl = 0xA2
	vkRControl = 0xA3
	vkLMenu    = 0xA4
	vkRMenu    = 0xA5
	vkOEM1     = 0xBA
	vkOEMPlus  = 0xBB
	vkOEMComma = 0xBC
	vkOEMMinus = 0xBD
	vkOEMPerio = 0xBE
	vkOEM2     = 0xBF
	vkOEM3     = 0xC0
	vkOEM4     = 0xDB
	vkOEM5     = 0xDC
	vkOEM6     = 0xDD
	vkOEM7     = 0xDE
)

// keyNames maps virtual key codes that do not fall into a contiguous range.

// keyNames maps virtual key codes that do not fall into a contiguous range.
var keyNames = map[uintptr]fyne.KeyName{
	vkBack:     fyne.KeyBackspace,
	vkTab:      fyne.KeyTab,
	vkReturn:   fyne.KeyReturn,
	vkEscape:   fyne.KeyEscape,
	vkSpace:    fyne.KeySpace,
	vkPrior:    fyne.KeyPageUp,
	vkNext:     fyne.KeyPageDown,
	vkEnd:      fyne.KeyEnd,
	vkHome:     fyne.KeyHome,
	vkLeft:     fyne.KeyLeft,
	vkUp:       fyne.KeyUp,
	vkRight:    fyne.KeyRight,
	vkDown:     fyne.KeyDown,
	vkInsert:   fyne.KeyInsert,
	vkDelete:   fyne.KeyDelete,
	vkSnapshot: desktop.KeyPrintScreen,
	vkCapital:  desktop.KeyCapsLock,
	vkLShift:   desktop.KeyShiftLeft,
	vkRShift:   desktop.KeyShiftRight,
	vkLControl: desktop.KeyControlLeft,
	vkRControl: desktop.KeyControlRight,
	vkLMenu:    desktop.KeyAltLeft,
	vkRMenu:    desktop.KeyAltRight,
	vkLWin:     desktop.KeySuperLeft,
	vkRWin:     desktop.KeySuperRight,
	vkApps:     desktop.KeyMenu,
	vkShift:    desktop.KeyShiftLeft,
	vkControl:  desktop.KeyControlLeft,
	vkMenu:     desktop.KeyAltLeft,
	vkMultiply: fyne.KeyAsterisk,
	vkAdd:      fyne.KeyPlus,
	vkSubtract: fyne.KeyMinus,
	vkDecimal:  fyne.KeyPeriod,
	vkDivide:   fyne.KeySlash,
	vkOEM1:     fyne.KeySemicolon,
	vkOEMPlus:  fyne.KeyEqual,
	vkOEMComma: fyne.KeyComma,
	vkOEMMinus: fyne.KeyMinus,
	vkOEMPerio: fyne.KeyPeriod,
	vkOEM2:     fyne.KeySlash,
	vkOEM3:     fyne.KeyBackTick,
	vkOEM4:     fyne.KeyLeftBracket,
	vkOEM5:     fyne.KeyBackslash,
	vkOEM6:     fyne.KeyRightBracket,
	vkOEM7:     fyne.KeyApostrophe,
}

var digitKeys = [10]fyne.KeyName{
	fyne.Key0, fyne.Key1, fyne.Key2, fyne.Key3, fyne.Key4,
	fyne.Key5, fyne.Key6, fyne.Key7, fyne.Key8, fyne.Key9,
}

var letterKeys = [26]fyne.KeyName{
	fyne.KeyA, fyne.KeyB, fyne.KeyC, fyne.KeyD, fyne.KeyE, fyne.KeyF, fyne.KeyG,
	fyne.KeyH, fyne.KeyI, fyne.KeyJ, fyne.KeyK, fyne.KeyL, fyne.KeyM, fyne.KeyN,
	fyne.KeyO, fyne.KeyP, fyne.KeyQ, fyne.KeyR, fyne.KeyS, fyne.KeyT, fyne.KeyU,
	fyne.KeyV, fyne.KeyW, fyne.KeyX, fyne.KeyY, fyne.KeyZ,
}

var functionKeys = [12]fyne.KeyName{
	fyne.KeyF1, fyne.KeyF2, fyne.KeyF3, fyne.KeyF4, fyne.KeyF5, fyne.KeyF6,
	fyne.KeyF7, fyne.KeyF8, fyne.KeyF9, fyne.KeyF10, fyne.KeyF11, fyne.KeyF12,
}

// keyToName translates a virtual key code to a Fyne key name. The contiguous
// ASCII-aligned ranges are handled arithmetically, everything else by table.
func keyToName(vk uintptr) fyne.KeyName {
	switch {
	case vk >= '0' && vk <= '9':
		return digitKeys[vk-'0']
	case vk >= 'A' && vk <= 'Z':
		return letterKeys[vk-'A']
	case vk >= vkF1 && vk < vkF1+12:
		return functionKeys[vk-vkF1]
	case vk >= vkNumpad0 && vk <= vkNumpad0+9:
		return digitKeys[vk-vkNumpad0]
	}
	if name, ok := keyNames[vk]; ok {
		return name
	}
	return fyne.KeyUnknown
}

// keyNameFor resolves a full key message to a name. The numpad Enter shares
// VK_RETURN with the main Return key and is told apart by the extended-key bit.
func keyNameFor(vk, lParam uintptr) fyne.KeyName {
	if vk == vkReturn && lParam&(1<<24) != 0 {
		return fyne.KeyEnter
	}
	return keyToName(vk)
}

// scanCodeFor extracts the hardware scan code from a key message, keeping the
// extended bit so left/right key variants stay distinct.
func scanCodeFor(lParam uintptr) int {
	return int((lParam >> 16) & 0x1ff)
}

// keyDown reports whether a virtual key is currently held. GetKeyState sets the
// high bit of the returned int16, which is the sign bit.
func keyDown(vk int32) bool {
	return getKeyState(vk) < 0
}

// currentModifiers reads the live modifier state. Windows reports modifiers via
// key state rather than packing them into the message, so they are sampled when
// each key or click event is handled.
func currentModifiers() fyne.KeyModifier {
	var mod fyne.KeyModifier
	if keyDown(vkShift) {
		mod |= fyne.KeyModifierShift
	}
	if keyDown(vkControl) {
		mod |= fyne.KeyModifierControl
	}
	if keyDown(vkMenu) {
		mod |= fyne.KeyModifierAlt
	}
	if keyDown(vkLWin) || keyDown(vkRWin) {
		mod |= fyne.KeyModifierSuper
	}
	return mod
}

func isKeyModifier(keyName fyne.KeyName) bool {
	switch keyName {
	case desktop.KeyShiftLeft, desktop.KeyShiftRight,
		desktop.KeyControlLeft, desktop.KeyControlRight,
		desktop.KeyAltLeft, desktop.KeyAltRight,
		desktop.KeySuperLeft, desktop.KeySuperRight:
		return true
	}
	return false
}

// shortcutFor maps a key plus modifiers onto one of Fyne's standard editing
// shortcuts, falling back to a desktop.CustomShortcut so applications can
// register their own.
func shortcutFor(keyName fyne.KeyName, modifier fyne.KeyModifier) fyne.Shortcut {
	var shortcut fyne.Shortcut
	if modifier == fyne.KeyModifierControl {
		switch keyName {
		case fyne.KeyZ:
			shortcut = &fyne.ShortcutUndo{}
		case fyne.KeyY:
			shortcut = &fyne.ShortcutRedo{}
		case fyne.KeyV:
			shortcut = &fyne.ShortcutPaste{Clipboard: NewClipboard()}
		case fyne.KeyC:
			shortcut = &fyne.ShortcutCopy{Clipboard: NewClipboard()}
		case fyne.KeyInsert:
			shortcut = &fyne.ShortcutCopy{Clipboard: NewClipboard(), Secondary: true}
		case fyne.KeyX:
			shortcut = &fyne.ShortcutCut{Clipboard: NewClipboard()}
		case fyne.KeyA:
			shortcut = &fyne.ShortcutSelectAll{}
		}
	}

	if modifier == fyne.KeyModifierShift {
		switch keyName {
		case fyne.KeyInsert:
			shortcut = &fyne.ShortcutPaste{Clipboard: NewClipboard(), Secondary: true}
		case fyne.KeyDelete:
			shortcut = &fyne.ShortcutCut{Clipboard: NewClipboard(), Secondary: true}
		}
	}

	if shortcut == nil && modifier != 0 && !isKeyModifier(keyName) && modifier != fyne.KeyModifierShift {
		shortcut = &desktop.CustomShortcut{KeyName: keyName, Modifier: modifier}
	}
	return shortcut
}
