package desktop

import (
	"fyne.io/fyne"
)

const (
	// KeyShift represents the left or right shift key
	KeyShift fyne.KeyName = "Shift"
	// KeyControl represents the left or right control key
	KeyControl fyne.KeyName = "Control"
	// KeyAlt represents the left or right alt key
	KeyAlt fyne.KeyName = "Alt"
	// KeySuper represents the left or right "Windows" key (or "Command" key on macOS)
	KeySuper fyne.KeyName = "Super"
	// KeyMenu represents the left or right menu / application key
	KeyMenu fyne.KeyName = "Menu"
)

// Modifier captures any key modifiers (shift etc) pressed during this key event
type Modifier int

const (
	// ShiftModifier represents a shift key being held
	ShiftModifier Modifier = 1 << iota
	// ControlModifier represents the ctrl key being held
	ControlModifier
	// AltModifier represents either alt keys being held
	AltModifier
	// SuperModifier represents either super keys being held
	SuperModifier
)
