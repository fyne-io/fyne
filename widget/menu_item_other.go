//go:build !darwin

package widget

import (
	"fyne.io/fyne/v2"
)

const (
	textModifierAlt     = "Alt+"
	textModifierControl = "Ctrl+"
	textModifierShift   = "Shift+"
	textModifierSuper   = "Super+"
)

var (
	styleModifiers   = fyne.TextStyle{}
	defaultStyleKeys = fyne.TextStyle{}
)

var keyTexts = map[fyne.KeyName]string{
	fyne.KeyBackspace: "Backspace",
	fyne.KeyDelete:    "Del",
	fyne.KeyDown:      "↓",
	fyne.KeyEnd:       "End",
	fyne.KeyEnter:     "Enter",
	fyne.KeyEscape:    "Esc",
	fyne.KeyHome:      "Home",
	fyne.KeyLeft:      "←",
	fyne.KeyPageDown:  "PgDn",
	fyne.KeyPageUp:    "PgUp",
	fyne.KeyReturn:    "Return",
	fyne.KeyRight:     "→",
	fyne.KeySpace:     "Space",
	fyne.KeyTab:       "Tab",
	fyne.KeyUp:        "↑",
}
