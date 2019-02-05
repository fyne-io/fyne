package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	String string

	Name      KeyName
	Modifiers Modifier
}

// PointerEvent describes a pointer input event. The position is relative to the top-left
// of the CanvasObject this is triggered on.
type PointerEvent struct {
	Position Position    // The position of the event
	Button   MouseButton // The pointer button which caused the event
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how many large the scroll was in two dimensions.
type ScrollEvent struct {
	DeltaX, DeltaY int
}
