package glfw

// Action represents the change of state of a key or mouse button event
type Action int

const (
	// Release Keyboard button was released
	Release Action = 0
	// Press Keyboard button was pressed
	Press Action = 1
	// Repeat Keyboard button was hold pressed for long enough that it trigger a repeat
	Repeat Action = 2
)
