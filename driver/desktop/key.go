package desktop

import (
	"fyne.io/fyne"
)

const (
	// KeyControl represents the left or right control key
	KeyControl fyne.KeyName = "Control"
	// KeyAlt represents the left or right alt key
	KeyAlt fyne.KeyName = "Alt"
	// KeySuper represents the left or right "Windows" key (or "Command" key on macOS)
	KeySuper fyne.KeyName = "Super"
	// KeyMenu represents the left or right menu / application key
	KeyMenu fyne.KeyName = "Menu"
	// KeyEscape is the "esc" key
	KeyEscape fyne.KeyName = "Escape"
	// KeyTab is the tab advance key
	KeyTab fyne.KeyName = "Tab"
	// KeyInsert is the insert mode key
	KeyInsert fyne.KeyName = "Insert"
	// KeyDelete is the delete-after-cursor key
	KeyDelete fyne.KeyName = "Delete"
	// KeyRight is the right arrow key
	KeyRight fyne.KeyName = "Right"
	// KeyLeft is the left arrow key
	KeyLeft fyne.KeyName = "Left"
	// KeyDown is the down arrow key
	KeyDown fyne.KeyName = "Down"
	// KeyUp is the up arrow key
	KeyUp fyne.KeyName = "Up"
	// KeyPageUp is the page up num-pad key
	KeyPageUp fyne.KeyName = "Prior"
	// KeyPageDown is the page down num-pad key
	KeyPageDown fyne.KeyName = "Next"
	// KeyHome is the line-home key
	KeyHome fyne.KeyName = "Home"
	// KeyEnd is the line-end key
	KeyEnd fyne.KeyName = "End"

	// KeyF1 is the first function key
	KeyF1 fyne.KeyName = "F1"
	// KeyF2 is the second function key
	KeyF2 fyne.KeyName = "F2"
	// KeyF3 is the third function key
	KeyF3 fyne.KeyName = "F3"
	// KeyF4 is the fourth function key
	KeyF4 fyne.KeyName = "F4"
	// KeyF5 is the fifth function key
	KeyF5 fyne.KeyName = "F5"
	// KeyF6 is the sixth function key
	KeyF6 fyne.KeyName = "F6"
	// KeyF7 is the seventh function key
	KeyF7 fyne.KeyName = "F7"
	// KeyF8 is the eighth function key
	KeyF8 fyne.KeyName = "F8"
	// KeyF9 is the ninth function key
	KeyF9 fyne.KeyName = "F9"
	// KeyF10 is the tenth function key
	KeyF10 fyne.KeyName = "F10"
	// KeyF11 is the eleventh function key
	KeyF11 fyne.KeyName = "F11"
	// KeyF12 is the twelfth function key
	KeyF12 fyne.KeyName = "F12"
	/*
		F13
		...
		F25
	*/

	// KeyEnter is the enter/ return key (keypad)
	KeyEnter fyne.KeyName = "KP_Enter"
)

// Modifier captures any key modifiers (shift etc) pressed during this key event
type Modifier int

const (
	// NoModifier represents a no modifier key has being held
	NoModifier Modifier = 1 << iota
	// ShiftModifier represents a shift key being held
	ShiftModifier
	// ControlModifier represents the ctrl key being held
	ControlModifier
	// AltModifier represents either alt keys being held
	AltModifier
	// SuperModifier represents either super keys being held
	SuperModifier
)
