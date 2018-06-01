package ui

import "github.com/fyne-io/fyne/api/ui/input"

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	String string

	Name      string
	Code      input.KeyCode
	Modifiers input.Modifier
}

// MouseEvent describes a pointer input event. The position is relative to the top-left
// of the CanvasObject this is triggered on.
type MouseEvent struct {
	Position Position          // The position of the event
	Button   input.MouseButton // The mouse button which caused the event
}
