package glfw

// Action represents the change of state of a key or mouse button event
type action int

const (
	// Release Keyboard button was released
	release action = 0
	// Press Keyboard button was pressed
	press action = 1
	// Repeat Keyboard button was hold pressed for long enough that it trigger a repeat
	repeat action = 2
)
