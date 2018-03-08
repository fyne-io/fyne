// Package event defines the various types of events that widgets or apps can
// listen for
package event

import "github.com/fyne-io/fyne/ui"

// MouseEvent describes an input event. The position is relative to the top-left
// of the CanvasObject this is triggered on.
type MouseEvent struct {
	Position ui.Position // The position of the event
	Button   MouseButton // The mouse button which caused the event
}

// MouseButton represents a single button in a MouseEvent
type MouseButton int

const (
	// LeftMouseButton is the most common mouse button - on some systems the only one
	LeftMouseButton MouseButton = iota + 1
	// Don't currently expose this button
	middleMouseButton
	// RightMouseButton is the secondary button on most mouse input devices.
	// For a touch screen this may be refer to a tap-and-hold action.
	RightMouseButton
)
