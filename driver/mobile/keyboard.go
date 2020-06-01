package mobile

import (
	"fyne.io/fyne"
)

// KeyboardType represents a type of virtual keyboard
type KeyboardType int32

const (
	// DefaultKeyboard is the keyboard with default input style and "return" return key
	DefaultKeyboard KeyboardType = iota
	// SingleLineKeyboard is the keyboard with default input style and "Done" return key
	SingleLineKeyboard
	// NumberKeyboard is the keyboard with number input style and "Done" return key
	NumberKeyboard
)

// Keyboardable describes any CanvasObject that needs a keyboard
type Keyboardable interface {
	fyne.Focusable

	Keyboard() KeyboardType
}
