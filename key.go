package fyne

// KeyName represents the name of a key that has been pressed
type KeyName string

const (
	// KeyEscape is the "esc" key
	KeyEscape KeyName = "Escape"
	// KeyReturn is the carriage return (main keyboard)
	KeyReturn KeyName = "Return"
	// KeyTab is the tab advance key
	KeyTab KeyName = "Tab"
	// KeyBackspace is the delete-before-cursor key
	KeyBackspace KeyName = "BackSpace"
	// KeyInsert is the insert mode key
	KeyInsert KeyName = "Insert"
	// KeyDelete is the delete-after-cursor key
	KeyDelete KeyName = "Delete"
	// KeyRight is the right arrow key
	KeyRight KeyName = "Right"
	// KeyLeft is the left arrow key
	KeyLeft KeyName = "Left"
	// KeyDown is the down arrow key
	KeyDown KeyName = "Down"
	// KeyUp is the up arrow key
	KeyUp KeyName = "Up"
	// KeyPageUp is the page up num-pad key
	KeyPageUp KeyName = "Prior"
	// KeyPageDown is the page down num-pad key
	KeyPageDown KeyName = "Next"
	// KeyHome is the line-home key
	KeyHome KeyName = "Home"
	// KeyEnd is the line-end key
	KeyEnd KeyName = "End"

	// KeyF1 is the first function key
	KeyF1 KeyName = "F1"
	// KeyF2 is the second function key
	KeyF2 KeyName = "F2"
	// KeyF3 is the third function key
	KeyF3 KeyName = "F3"
	// KeyF4 is the fourth function key
	KeyF4 KeyName = "F4"
	// KeyF5 is the fifth function key
	KeyF5 KeyName = "F5"
	// KeyF6 is the sixth function key
	KeyF6 KeyName = "F6"
	// KeyF7 is the seventh function key
	KeyF7 KeyName = "F7"
	// KeyF8 is the eighth function key
	KeyF8 KeyName = "F8"
	// KeyF9 is the ninth function key
	KeyF9 KeyName = "F9"
	// KeyF10 is the tenth function key
	KeyF10 KeyName = "F10"
	// KeyF11 is the eleventh function key
	KeyF11 KeyName = "F11"
	// KeyF12 is the twelfth function key
	KeyF12 KeyName = "F12"
	/*
		F13
		...
		F25
	*/

	// KeyEnter is the enter/ return key (keypad)
	KeyEnter KeyName = "KP_Enter"
)
