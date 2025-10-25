package desktop

import (
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
)

// Declare conformity with Shortcut interface
var (
	_ fyne.Shortcut         = (*CustomShortcut)(nil)
	_ fyne.KeyboardShortcut = (*CustomShortcut)(nil)
)

// CustomShortcut describes a shortcut desktop event.
type CustomShortcut struct {
	fyne.KeyName
	Modifier fyne.KeyModifier
}

// Key returns the key name of this shortcut.
func (cs *CustomShortcut) Key() fyne.KeyName {
	return cs.KeyName
}

// Mod returns the modifier of this shortcut.
func (cs *CustomShortcut) Mod() fyne.KeyModifier {
	return cs.Modifier
}

// ShortcutName returns the shortcut name associated to the event.
func (cs *CustomShortcut) ShortcutName() string {
	id := &strings.Builder{}
	id.WriteString("CustomDesktop:")
	writeModifiers(id, cs.Modifier)
	id.WriteString(string(cs.KeyName))
	return id.String()
}

func writeModifiers(w *strings.Builder, mods fyne.KeyModifier) {
	if (mods & fyne.KeyModifierShift) != 0 {
		w.WriteString("Shift+")
	}
	if (mods & fyne.KeyModifierControl) != 0 {
		w.WriteString("Control+")
	}
	if (mods & fyne.KeyModifierAlt) != 0 {
		w.WriteString("Alt+")
	}
	if (mods & fyne.KeyModifierSuper) != 0 {
		if runtime.GOOS == "darwin" {
			w.WriteString("Command+")
		} else {
			w.WriteString("Super+")
		}
	}
}
