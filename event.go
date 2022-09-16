package fyne

import "time"

// HardwareKey contains information associated with physical key events
// Most applications should use KeyName for cross-platform compatibility.
type HardwareKey struct {
	// ScanCode represents a hardware ID for (normally desktop) keyboard events.
	ScanCode int
}

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	// Name describes the keyboard event that is consistent across platforms.
	Name KeyName
	// Physical is a platform specific field that reports the hardware information of physical keyboard events.
	Physical HardwareKey
}

// PointEvent describes a pointer input event. The position is relative to the
// top-left of the CanvasObject this is triggered on.
type PointEvent struct {
	AbsolutePosition Position // The absolute position of the event
	Position         Position // The relative position of the event
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how large the scroll was in two dimensions.
type ScrollEvent struct {
	PointEvent
	Scrolled Delta
}

// DragEvent defines the parameters of a pointer or other drag event.
// The DraggedX and DraggedY fields show how far the item was dragged since the last event.
type DragEvent struct {
	PointEvent
	Dragged Delta
}

func NewRuneEvent(char rune) *RuneEvent {
	return &RuneEvent{
		char:      char,
		eventTime: time.Now(),
	}
}

type RuneEvent struct {
	char      rune
	eventTime time.Time
}

func (r *RuneEvent) Rune() rune {
	return r.char
}

func (r *RuneEvent) When() time.Time {
	return r.eventTime
}
