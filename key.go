package fyne

// KeyName represents the name of a key that has been pressed
type KeyName string

const (
	// KeyEscape is the "esc" key
	KeyEscape KeyName = "Escape"
	// KeyReturn is the carriage return (main keyboard)
	KeyReturn = "Return"
	// KeyTab is the tab advance key
	KeyTab = "Tab"
	// KeyBackspace is the delete-before-cursor key
	KeyBackspace = "BackSpace"
	// KeyInsert is the insert mode key
	KeyInsert = "Insert"
	// KeyDelete is the delete-after-cursor key
	KeyDelete = "Delete"
	// KeyRight is the right arrow key
	KeyRight = "Right"
	// KeyLeft is the left arrow key
	KeyLeft = "Left"
	// KeyDown is the down arrow key
	KeyDown = "Down"
	// KeyUp is the up arrow key
	KeyUp = "Up"
	// KeyPageUp is the page up num-pad key
	KeyPageUp = "Prior"
	// KeyPageDown is the page down num-pad key
	KeyPageDown = "Next"
	// KeyHome is the line-home key
	KeyHome = "Home"
	// KeyEnd is the line-end key
	KeyEnd = "End"

	// KeyF1 is the first function key
	KeyF1 = "F1"
	// KeyF2 is the second function key
	KeyF2 = "F2"
	// KeyF3 is the third function key
	KeyF3 = "F3"
	// KeyF4 is the fourth function key
	KeyF4 = "F4"
	// KeyF5 is the fifth function key
	KeyF5 = "F5"
	// KeyF6 is the sixth function key
	KeyF6 = "F6"
	// KeyF7 is the seventh function key
	KeyF7 = "F7"
	// KeyF8 is the eighth function key
	KeyF8 = "F8"
	// KeyF9 is the ninth function key
	KeyF9 = "F9"
	// KeyF10 is the tenth function key
	KeyF10 = "F10"
	// KeyF11 is the eleventh function key
	KeyF11 = "F11"
	// KeyF12 is the twelfth function key
	KeyF12 = "F12"
	/*
		F13
		...
		F25
	*/

	// KeyEnter is the enter/ return key (keypad)
	KeyEnter = "KP_Enter"
)
