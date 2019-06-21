package desktop

import (
	"runtime"
	"strings"

	"fyne.io/fyne"
)

// Declare conformity with Shortcut interface
var _ fyne.Shortcut = (*CustomShortcut)(nil)

// CustomShortcut describes a shortcut desktop event.
type CustomShortcut struct {
	fyne.KeyName
	Modifier
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

func modifierToString(mods Modifier) string {
	s := []string{}
	if (mods & ShiftModifier) != 0 {
		s = append(s, string("Shift"))
	}
	if (mods & ControlModifier) != 0 {
		s = append(s, string("Control"))
	}
	if (mods & AltModifier) != 0 {
		s = append(s, string("Alt"))
	}
	if (mods & SuperModifier) != 0 {
		if runtime.GOOS == "darwin" {
			s = append(s, string("Command"))
		} else {
			s = append(s, string("Super"))
		}
	}
	return strings.Join(s, "+")
}
