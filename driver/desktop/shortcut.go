package desktop

import (
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
)

// Declare conformity with Shortcut interface
var _ fyne.Shortcut = (*CustomShortcut)(nil)
var _ fyne.KeyboardShortcut = (*CustomShortcut)(nil)

// CustomShortcut describes a shortcut desktop event.
type CustomShortcut struct {
	fyne.KeyName
	Modifier fyne.KeyModifier
}

// Key returns the key name of this shortcut.
// @implements KeyboardShortcut
func (cs *CustomShortcut) Key() fyne.KeyName {
	return cs.KeyName
}

// Mod returns the modifier of this shortcut.
// @implements KeyboardShortcut
func (cs *CustomShortcut) Mod() fyne.KeyModifier {
	return cs.Modifier
}

// ShortcutName returns the shortcut name associated to the event
func (cs *CustomShortcut) ShortcutName() string {
	id := &strings.Builder{}
	id.WriteString("CustomDesktop:")
	id.WriteString(modifierToString(cs.Modifier))
	id.WriteString("+")
	id.WriteString(string(cs.KeyName))
	return id.String()
}

func modifierToString(mods fyne.KeyModifier) string {
	s := []string{}
	if (mods & fyne.KeyModifierShift) != 0 {
		s = append(s, string("Shift"))
	}
	if (mods & fyne.KeyModifierControl) != 0 {
		s = append(s, string("Control"))
	}
	if (mods & fyne.KeyModifierAlt) != 0 {
		s = append(s, string("Alt"))
	}
	if (mods & fyne.KeyModifierSuper) != 0 {
		if runtime.GOOS == "darwin" {
			s = append(s, string("Command"))
		} else {
			s = append(s, string("Super"))
		}
	}
	return strings.Join(s, "+")
}
