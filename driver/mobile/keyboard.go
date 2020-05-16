package mobile

import (
	"fyne.io/fyne"
)

// Keyboard represents a standard fyne keyboard
type Keyboard int32

const (
	// DefaultKeyboard is the keyboard with default input style and "Done" return key
	DefaultKeyboard Keyboard = iota
	// MultiLineKeyboard is the keyboard with default input style and "return" return key
	MultiLineKeyboard
	// NumberKeyboard is the keyboard with number input style and "Done" return key
	NumberKeyboard
	// WebKeyboard is the keyboard with web input style and "Done" return key
	WebKeyboard
)

// Keyboardable describes any CanvasObject that needs a keyboard
type Keyboardable interface {
	fyne.Focusable

	Keyboard() Keyboard
}
