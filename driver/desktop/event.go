package desktop

import (
	"fyne.io/fyne"
)

// ShortcutDesktopEvent describes a shortcut clipoard event.
type ShortcutDesktopEvent struct {
	fyne.ShortcutName
	fyne.KeyName
	Modifier
}

// Shortcut returns the shortcut name associated to the event
func (se *ShortcutDesktopEvent) Shortcut() fyne.ShortcutName {
	return se.ShortcutName
}
