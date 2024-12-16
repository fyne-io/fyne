package widget

import (
	"fyne.io/fyne/v2"
)

const (
	textModifierAlt     = "⌥"
	textModifierControl = "⌃"
	textModifierShift   = "⇧"
	textModifierSuper   = "⌘"
)

var (
	styleModifiers   = fyne.TextStyle{Symbol: true}
	defaultStyleKeys = fyne.TextStyle{Monospace: true}
)

var keyTexts = map[fyne.KeyName]string{
	fyne.KeyBackspace: "⌫",
	fyne.KeyDelete:    "⌦",
	fyne.KeyDown:      "↓",
	fyne.KeyEnd:       "↘",
	fyne.KeyEnter:     "↩",
	fyne.KeyEscape:    "⎋",
	fyne.KeyHome:      "↖",
	fyne.KeyLeft:      "←",
	fyne.KeyPageDown:  "⇟",
	fyne.KeyPageUp:    "⇞",
	fyne.KeyReturn:    "↩",
	fyne.KeyRight:     "→",
	fyne.KeySpace:     "␣",
	fyne.KeyTab:       "⇥",
	fyne.KeyUp:        "↑",
}
