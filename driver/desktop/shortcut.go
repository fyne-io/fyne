package desktop

import (
	"fyne.io/fyne"
)

// CustomShortcut describes a shortcut desktop event.
type CustomShortcut struct {
	fyne.KeyName
	Modifier
}

// Shortcut returns the shortcut name associated to the event
func (cs *CustomShortcut) Shortcut() string {
	return "CustomDesktop"
}
