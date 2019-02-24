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

// Shortcut returns the shortcut name associated to the event
func (cs *CustomShortcut) Shortcut() string {
	id := &strings.Builder{}
	id.WriteString("CustomDesktop:")
	id.WriteString(modifierToString(cs.Modifier))
	id.WriteString("+")
	id.WriteString(string(cs.KeyName))
	return id.String()
}
