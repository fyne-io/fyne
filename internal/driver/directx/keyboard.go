//go:build windows && directx

// Keyboard event processing: key down/up, typed runes, shortcuts and tab focus
// traversal.

package directx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func (w *window) capturesTab(modifier fyne.KeyModifier) bool {
	if ent, ok := w.canvas.Focused().(fyne.Tabbable); ok && ent.AcceptsTab() {
		return true
	}
	switch modifier {
	case 0:
		w.canvas.FocusNext()
	case fyne.KeyModifierShift:
		w.canvas.FocusPrevious()
	}
	return false
}

func (w *window) processKeyPressed(keyName fyne.KeyName, scanCode int, act action, repeat bool) {
	if keyName == fyne.KeyUnknown {
		return
	}
	modifier := currentModifiers()
	keyEvent := &fyne.KeyEvent{Name: keyName, Physical: fyne.HardwareKey{ScanCode: scanCode}, Repeat: repeat}

	if act == release {
		if focused := w.canvas.Focused(); focused != nil {
			if keyable, ok := focused.(desktop.Keyable); ok {
				keyable.KeyUp(keyEvent)
			}
		} else if w.canvas.onKeyUp != nil {
			w.canvas.onKeyUp(keyEvent)
		}
		return
	}

	if focused := w.canvas.Focused(); focused != nil {
		if keyable, ok := focused.(desktop.Keyable); ok {
			keyable.KeyDown(keyEvent)
		}
	} else if w.canvas.onKeyDown != nil {
		w.canvas.onKeyDown(keyEvent)
	}

	modifierOtherThanShift := (modifier & fyne.KeyModifierControl) |
		(modifier & fyne.KeyModifierAlt) | (modifier & fyne.KeyModifierSuper)
	if keyName == fyne.KeyTab && modifierOtherThanShift == 0 && !w.capturesTab(modifier) {
		return
	}
	if shortcut := shortcutFor(keyName, modifier); shortcut != nil {
		if w.triggerMainMenuShortcut(shortcut) {
			return
		}
		if focused, ok := w.canvas.Focused().(fyne.Shortcutable); ok {
			// A disabled selectable widget still allows Copy, but nothing that would
			// change its content.
			type selectableText interface {
				fyne.Disableable
				SelectedText() string
			}
			if sel, ok := focused.(selectableText); ok && sel.Disabled() &&
				shortcut.ShortcutName() != "Copy" {
				return
			}
			focused.TypedShortcut(shortcut)
			return
		}
		w.canvas.TypedShortcut(shortcut)
		return
	}

	if focused := w.canvas.Focused(); focused != nil {
		focused.TypedKey(keyEvent)
	} else if w.canvas.onTypedKey != nil {
		w.canvas.onTypedKey(keyEvent)
	}
}

func (w *window) processCharInput(char rune) {
	if focused := w.canvas.Focused(); focused != nil {
		focused.TypedRune(char)
	} else if w.canvas.onTypedRune != nil {
		w.canvas.onTypedRune(char)
	}
}
