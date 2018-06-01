package input

// KeyCode represents the numeric code of a key pressed without modifiers
type KeyCode int

// Modifier captures any key modifiers (shift etc) pressed during this key event
type Modifier int

const (
	// ShiftModifier represents a shift key being held
	ShiftModifier Modifier = 1 << iota
	// ControlModifier represents the ctrl key being held
	ControlModifier
	// AltModifier represents either alt keys being held
	AltModifier
)
