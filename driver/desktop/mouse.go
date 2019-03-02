package desktop

import "fyne.io/fyne"

// MouseButton represents a single button in a desktop MouseEvent
type MouseButton int

const (
	// LeftMouseButton is the most common mouse button - on some systems the only one
	LeftMouseButton MouseButton = iota + 1
	// Don't currently expose this button
	middleMouseButton
	// RightMouseButton is the secondary button on most mouse input devices.
	RightMouseButton
)

// MouseEvent contains data relating to desktop mouse events
type MouseEvent struct {
	fyne.PointEvent
	Button MouseButton
}

// Mouseable represents desktop mouse events that can be sent to CanvasObjects
type Mouseable interface {
	MouseDown(*MouseEvent)
	MouseUp(*MouseEvent)
}
