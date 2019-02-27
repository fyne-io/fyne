package desktop

import (
	"strings"

	"fyne.io/fyne"
)

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
		s = append(s, string(KeyShift))
	}
	if (mods & ControlModifier) != 0 {
		s = append(s, string(KeyControl))
	}
	if (mods & AltModifier) != 0 {
		s = append(s, string(KeyAlt))
	}
	if (mods & SuperModifier) != 0 {
		s = append(s, string(KeySuper))
	}
	return strings.Join(s, "+")
}
