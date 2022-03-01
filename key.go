package fyne

// KeyName represents the name of a key that has been pressed.
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

	// Key0 represents the key 0
	Key0 KeyName = "0"
	// Key1 represents the key 1
	Key1 KeyName = "1"
	// Key2 represents the key 2
	Key2 KeyName = "2"
	// Key3 represents the key 3
	Key3 KeyName = "3"
	// Key4 represents the key 4
	Key4 KeyName = "4"
	// Key5 represents the key 5
	Key5 KeyName = "5"
	// Key6 represents the key 6
	Key6 KeyName = "6"
	// Key7 represents the key 7
	Key7 KeyName = "7"
	// Key8 represents the key 8
	Key8 KeyName = "8"
	// Key9 represents the key 9
	Key9 KeyName = "9"
	// KeyA represents the key A
	KeyA KeyName = "A"
	// KeyB represents the key B
	KeyB KeyName = "B"
	// KeyC represents the key C
	KeyC KeyName = "C"
	// KeyD represents the key D
	KeyD KeyName = "D"
	// KeyE represents the key E
	KeyE KeyName = "E"
	// KeyF represents the key F
	KeyF KeyName = "F"
	// KeyG represents the key G
	KeyG KeyName = "G"
	// KeyH represents the key H
	KeyH KeyName = "H"
	// KeyI represents the key I
	KeyI KeyName = "I"
	// KeyJ represents the key J
	KeyJ KeyName = "J"
	// KeyK represents the key K
	KeyK KeyName = "K"
	// KeyL represents the key L
	KeyL KeyName = "L"
	// KeyM represents the key M
	KeyM KeyName = "M"
	// KeyN represents the key N
	KeyN KeyName = "N"
	// KeyO represents the key O
	KeyO KeyName = "O"
	// KeyP represents the key P
	KeyP KeyName = "P"
	// KeyQ represents the key Q
	KeyQ KeyName = "Q"
	// KeyR represents the key R
	KeyR KeyName = "R"
	// KeyS represents the key S
	KeyS KeyName = "S"
	// KeyT represents the key T
	KeyT KeyName = "T"
	// KeyU represents the key U
	KeyU KeyName = "U"
	// KeyV represents the key V
	KeyV KeyName = "V"
	// KeyW represents the key W
	KeyW KeyName = "W"
	// KeyX represents the key X
	KeyX KeyName = "X"
	// KeyY represents the key Y
	KeyY KeyName = "Y"
	// KeyZ represents the key Z
	KeyZ KeyName = "Z"

	// KeySpace is the space key
	KeySpace KeyName = "Space"
	// KeyApostrophe is the key "'"
	KeyApostrophe KeyName = "'"
	// KeyComma is the key ","
	KeyComma KeyName = ","
	// KeyMinus is the key "-"
	KeyMinus KeyName = "-"
	// KeyPeriod is the key "." (full stop)
	KeyPeriod KeyName = "."
	// KeySlash is the key "/"
	KeySlash KeyName = "/"
	// KeyBackslash is the key "\"
	KeyBackslash KeyName = "\\"
	// KeyLeftBracket is the key "["
	KeyLeftBracket KeyName = "["
	// KeyRightBracket is the key "]"
	KeyRightBracket KeyName = "]"
	// KeySemicolon is the key ";"
	KeySemicolon KeyName = ";"
	// KeyEqual is the key "="
	KeyEqual KeyName = "="
	// KeyAsterisk is the keypad key "*"
	KeyAsterisk KeyName = "*"
	// KeyPlus is the keypad key "+"
	KeyPlus KeyName = "+"
	// KeyBackTick is the key "`" on a US keyboard
	KeyBackTick KeyName = "`"

	// KeyUnknown is used for key events where the underlying hardware generated an
	// event that Fyne could not decode.
	//
	// Since: 2.1
	KeyUnknown KeyName = ""
)

// KeyModifier represents any modifier key (shift etc.) that is being pressed together with a key.
//
// Since: 2.2
type KeyModifier int

const (
	// KeyModifierShift represents a shift key being held
	//
	// Since: 2.2
	KeyModifierShift KeyModifier = 1 << iota
	// KeyModifierControl represents the ctrl key being held
	//
	// Since: 2.2
	KeyModifierControl
	// KeyModifierAlt represents either alt keys being held
	//
	// Since: 2.2
	KeyModifierAlt
	// KeyModifierSuper represents either super keys being held
	//
	// Since: 2.2
	KeyModifierSuper
)
