package desktop

import "fyne.io/fyne/v2"

// MouseButton represents a single button in a desktop MouseEvent
type MouseButton int

const (
	// MouseButtonPrimary is the most common mouse button - on some systems the only one.
	// This will normally be on the left side of a mouse.
	//
	// Since: 2.0
	MouseButtonPrimary MouseButton = 1 << iota

	// MouseButtonSecondary is the secondary button on most mouse input devices.
	// This will normally be on the right side of a mouse.
	//
	// Since: 2.0
	MouseButtonSecondary

	// MouseButtonTertiary is the middle button on the mouse, assuming it has one.
	//
	// Since: 2.0
	MouseButtonTertiary

	// LeftMouseButton is the most common mouse button - on some systems the only one.
	//
	// Deprecated: use MouseButtonPrimary which will adapt to mouse configuration.
	LeftMouseButton = MouseButtonPrimary

	// RightMouseButton is the secondary button on most mouse input devices.
	//
	// Deprecated: use MouseButtonSecondary which will adapt to mouse configuration.
	RightMouseButton = MouseButtonSecondary
)

// MouseEvent contains data relating to desktop mouse events
type MouseEvent struct {
	fyne.PointEvent
	Button   MouseButton
	Modifier fyne.KeyModifier
}

// Mouseable represents desktop mouse events that can be sent to CanvasObjects
type Mouseable interface {
	MouseDown(*MouseEvent)
	MouseUp(*MouseEvent)
}

// Hoverable is used when a canvas object wishes to know if a pointer device moves over it.
type Hoverable interface {
	// MouseIn is a hook that is called if the mouse pointer enters the element.
	MouseIn(*MouseEvent)
	// MouseMoved is a hook that is called if the mouse pointer moved over the element.
	MouseMoved(*MouseEvent)
	// MouseOut is a hook that is called if the mouse pointer leaves the element.
	MouseOut()
}
