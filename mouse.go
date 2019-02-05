package fyne

type MouseButton int
// MouseButton represents a single button in a PointerEvent

const (
	// LeftMouseButton is the most common mouse button - on some systems the only one
	LeftMouseButton MouseButton = iota + 1
	// Don't currently expose this button
	middleMouseButton
	// RightMouseButton is the secondary button on most mouse input devices.
	// For a touch screen this may be refer to a tap-and-hold action.
	RightMouseButton
)
