package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	String string

	Name      string
	Modifiers Modifier
}

// MouseEvent describes a pointer input event. The position is relative to the top-left
// of the CanvasObject this is triggered on.
type MouseEvent struct {
	Position Position    // The position of the event
	Button   MouseButton // The mouse button which caused the event
}
