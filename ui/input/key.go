package input

// KeyCode represents the numeric code of a key pressed without modifiers
type KeyCode int

// Modifier captures any key modifiers (shift etc) pressed during this key event
type Modifier int

const (
	ShiftModifier Modifier = 1 << iota
	ControlModifier
	AltModifier
)
